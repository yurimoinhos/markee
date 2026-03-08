package application

import (
	"context"
	"fmt"
	"log"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/modules/payment/domain"
	"github.com/aggi-tech/aggipay/platform/common"
)

// UseCase processa pagamentos ao receber BalanceReserved e publica o resultado.
type UseCase struct {
	repo      domain.Repository
	publisher contracts.EventPublisher
}

func NewUseCase(repo domain.Repository, publisher contracts.EventPublisher) *UseCase {
	return &UseCase{repo: repo, publisher: publisher}
}

// ProcessPayment é chamado quando um evento BalanceReserved é recebido.
// Garante idempotência via IdempotencyKey antes de processar.
func (uc *UseCase) ProcessPayment(ctx context.Context, evt contracts.BalanceReserved) error {
	// Idempotência: não reprocessar pagamentos já realizados
	exists, err := uc.repo.ExistsByIdempotencyKey(ctx, evt.IdempotencyKey)
	if err != nil {
		return fmt.Errorf("payment: erro ao verificar idempotência: %w", err)
	}
	if exists {
		log.Printf("[payment] pagamento já processado (idempotency_key=%s), ignorando", evt.IdempotencyKey)
		return nil
	}

	tx, err := uc.repo.Create(ctx, domain.PaymentTx{
		SagaID:         evt.SagaID,
		OrderID:        evt.OrderID,
		AmountCents:    evt.AmountCents,
		Status:         domain.PaymentStatusPending,
		IdempotencyKey: evt.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("payment: falha ao registrar transação: %w", err)
	}

	// Simula processamento com provider externo
	providerRef, err := uc.callProvider(ctx, tx)
	if err != nil {
		log.Printf("[payment] provider falhou para saga=%s: %v", evt.SagaID, err)

		_ = uc.repo.UpdateStatus(ctx, tx.ID, domain.PaymentStatusFailed, "")

		failedEvt := contracts.NewPaymentFailed(evt.SagaID, evt.OrderID, err.Error())
		return uc.publisher.Publish(ctx, contracts.TopicPaymentFailed, failedEvt)
	}

	if err := uc.repo.UpdateStatus(ctx, tx.ID, domain.PaymentStatusCompleted, providerRef); err != nil {
		return fmt.Errorf("payment: falha ao atualizar status: %w", err)
	}

	processedEvt := contracts.NewPaymentProcessed(evt.SagaID, evt.OrderID, providerRef, evt.AmountCents)
	if err := uc.publisher.Publish(ctx, contracts.TopicPaymentProcessed, processedEvt); err != nil {
		return fmt.Errorf("payment: falha ao publicar PaymentProcessed: %w", err)
	}

	log.Printf("[payment] pagamento processado com sucesso: saga=%s provider_ref=%s", evt.SagaID, providerRef)
	return nil
}

// GetBySagaID consulta o status de uma transação pelo saga_id.
func (uc *UseCase) GetBySagaID(ctx context.Context, sagaID string) (*domain.PaymentTx, error) {
	return uc.repo.FindBySagaID(ctx, sagaID)
}

// callProvider simula a chamada a um gateway de pagamento externo.
// Em produção, substituir pela integração real (Stripe, Adyen, etc.).
func (uc *UseCase) callProvider(_ context.Context, tx *domain.PaymentTx) (string, error) {
	// Simula processamento bem-sucedido gerando uma referência de provider
	providerRef := fmt.Sprintf("prov_%s", common.GenID().Value)
	log.Printf("[payment] provider processou transação: tx=%s ref=%s amount=%d", tx.ID, providerRef, tx.AmountCents)
	return providerRef, nil
}
