package state

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
	entstate "github.com/aggi-tech/aggipay/ent/sagastate"
	"github.com/aggi-tech/aggipay/platform/common"
)

// SagaStatus espelha os valores válidos do campo status no ent.
type SagaStatus = string

const (
	StatusPending         SagaStatus = "pending"
	StatusBalanceReserved SagaStatus = "balance_reserved"
	StatusPaymentSent     SagaStatus = "payment_sent"
	StatusConfirmed       SagaStatus = "confirmed"
	StatusCompensating    SagaStatus = "compensating"
	StatusCancelled       SagaStatus = "cancelled"
)

// Store persiste e consulta o estado das sagas no Postgres.
type Store struct {
	client *ent.Client
}

func NewStore(client *ent.Client) *Store {
	return &Store{client: client}
}

// Create cria um novo registro de saga no estado inicial "pending".
func (s *Store) Create(ctx context.Context, sagaType, orderID, userID, idempotencyKey string, amountCents uint64) (*ent.SagaState, error) {
	return s.client.SagaState.Create().
		SetID(common.GenID().Value).
		SetSagaType(sagaType).
		SetOrderID(orderID).
		SetUserID(userID).
		SetAmountCents(amountCents).
		SetCurrentStep("create_order").
		SetStatus(StatusPending).
		SetIdempotencyKey(idempotencyKey).
		Save(ctx)
}

// Advance atualiza o step e status corrente de uma saga.
func (s *Store) Advance(ctx context.Context, id, step, status string) error {
	_, err := s.client.SagaState.UpdateOneID(id).
		SetCurrentStep(step).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("saga state: falha ao avançar saga %s: %w", id, err)
	}
	return nil
}

// Complete marca a saga como confirmada com a referência do provider.
func (s *Store) Complete(ctx context.Context, id, providerRef string) error {
	_, err := s.client.SagaState.UpdateOneID(id).
		SetStatus(StatusConfirmed).
		SetCurrentStep("confirmed").
		SetProviderRef(providerRef).
		Save(ctx)
	return err
}

// Compensate marca a saga como cancelada registrando o motivo.
func (s *Store) Compensate(ctx context.Context, id, reason string) error {
	_, err := s.client.SagaState.UpdateOneID(id).
		SetStatus(StatusCancelled).
		SetCurrentStep("cancelled").
		SetFailureReason(reason).
		Save(ctx)
	return err
}

// FindByOrderID localiza a saga de pagamento pelo order_id.
func (s *Store) FindByOrderID(ctx context.Context, orderID string) (*ent.SagaState, error) {
	return s.client.SagaState.Query().
		Where(
			entstate.OrderID(orderID),
			entstate.SagaType("payment_saga"),
		).
		First(ctx)
}
