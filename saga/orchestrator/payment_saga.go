package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/saga/state"
)

// PaymentSaga orquestra o fluxo de pagamento consumindo e publicando eventos.
// Não faz chamadas diretas a outros módulos — comunica-se exclusivamente via RabbitMQ.
type PaymentSaga struct {
	store     *state.Store
	publisher contracts.EventPublisher
}

func NewPaymentSaga(store *state.Store, publisher contracts.EventPublisher) *PaymentSaga {
	return &PaymentSaga{store: store, publisher: publisher}
}

// RegisterConsumers registra todos os handlers de eventos que a saga precisa consumir.
func (s *PaymentSaga) RegisterConsumers(consumer contracts.EventConsumer) {
	consumer.Subscribe(contracts.TopicPaymentReceived, s.onPaymentReceived)
	consumer.Subscribe(contracts.TopicPaymentProcessed, s.onPaymentProcessed)
	consumer.Subscribe(contracts.TopicPaymentFailed, s.onPaymentFailed)
}

// onPaymentReceived: order publicou PaymentReceived → saga cria estado e reserva saldo.
func (s *PaymentSaga) onPaymentReceived(ctx context.Context, raw []byte) error {
	var evt contracts.PaymentReceived
	if err := json.Unmarshal(raw, &evt); err != nil {
		return fmt.Errorf("saga: falha ao deserializar PaymentReceived: %w", err)
	}

	log.Printf("[saga] PaymentReceived: saga=%s order=%s user=%s amount=%d",
		evt.SagaID, evt.OrderID, evt.UserID, evt.AmountCents)

	// Cria o registro de estado da saga
	sagaRow, err := s.store.Create(ctx,
		"payment_saga",
		evt.OrderID,
		evt.UserID,
		evt.EventID, // usa event_id como chave de idempotência
		evt.AmountCents,
	)
	if err != nil {
		return fmt.Errorf("saga: falha ao criar estado: %w", err)
	}

	// Avança para "reserva de saldo" e notifica o payment-service
	if err := s.store.Advance(ctx, sagaRow.ID, "reserve_balance", state.StatusBalanceReserved); err != nil {
		return err
	}

	balanceEvt := contracts.NewBalanceReserved(
		evt.SagaID,
		evt.OrderID,
		evt.UserID,
		evt.EventID, // idempotency_key
		evt.AmountCents,
	)

	if err := s.publisher.Publish(ctx, contracts.TopicBalanceReserved, balanceEvt); err != nil {
		return fmt.Errorf("saga: falha ao publicar BalanceReserved: %w", err)
	}

	log.Printf("[saga] BalanceReserved publicado: saga=%s", evt.SagaID)
	return nil
}

// onPaymentProcessed: payment publicou PaymentProcessed → saga confirma a ordem.
func (s *PaymentSaga) onPaymentProcessed(ctx context.Context, raw []byte) error {
	var evt contracts.PaymentProcessed
	if err := json.Unmarshal(raw, &evt); err != nil {
		return fmt.Errorf("saga: falha ao deserializar PaymentProcessed: %w", err)
	}

	log.Printf("[saga] PaymentProcessed: saga=%s order=%s provider=%s",
		evt.SagaID, evt.OrderID, evt.ProviderRef)

	sagaRow, err := s.store.FindByOrderID(ctx, evt.OrderID)
	if err != nil {
		return fmt.Errorf("saga: estado não encontrado para order=%s: %w", evt.OrderID, err)
	}

	if err := s.store.Complete(ctx, sagaRow.ID, evt.ProviderRef); err != nil {
		return err
	}

	confirmedEvt := contracts.NewOrderConfirmed(evt.SagaID, evt.OrderID, evt.ProviderRef)
	if err := s.publisher.Publish(ctx, contracts.TopicOrderConfirmed, confirmedEvt); err != nil {
		return fmt.Errorf("saga: falha ao publicar OrderConfirmed: %w", err)
	}

	log.Printf("[saga] OrderConfirmed publicado: saga=%s", evt.SagaID)
	return nil
}

// onPaymentFailed: payment publicou PaymentFailed → saga compensa (rollback).
func (s *PaymentSaga) onPaymentFailed(ctx context.Context, raw []byte) error {
	var evt contracts.PaymentFailed
	if err := json.Unmarshal(raw, &evt); err != nil {
		return fmt.Errorf("saga: falha ao deserializar PaymentFailed: %w", err)
	}

	log.Printf("[saga] PaymentFailed: saga=%s order=%s reason=%s",
		evt.SagaID, evt.OrderID, evt.Reason)

	sagaRow, err := s.store.FindByOrderID(ctx, evt.OrderID)
	if err != nil {
		return fmt.Errorf("saga: estado não encontrado para order=%s: %w", evt.OrderID, err)
	}

	if err := s.store.Compensate(ctx, sagaRow.ID, evt.Reason); err != nil {
		return err
	}

	// Publica compensações em paralelo
	balanceFailedEvt := contracts.NewBalanceFailed(evt.SagaID, evt.OrderID, sagaRow.UserID, evt.Reason)
	if err := s.publisher.Publish(ctx, contracts.TopicBalanceFailed, balanceFailedEvt); err != nil {
		log.Printf("[saga] WARN: falha ao publicar BalanceFailed: %v", err)
	}

	cancelledEvt := contracts.NewOrderCancelled(evt.SagaID, evt.OrderID, evt.Reason)
	if err := s.publisher.Publish(ctx, contracts.TopicOrderCancelled, cancelledEvt); err != nil {
		return fmt.Errorf("saga: falha ao publicar OrderCancelled: %w", err)
	}

	log.Printf("[saga] compensação concluída: saga=%s", evt.SagaID)
	return nil
}
