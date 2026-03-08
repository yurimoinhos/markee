package application

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/modules/order/domain"
	"github.com/aggi-tech/aggipay/platform/common"
)

// Service define os casos de uso do módulo order.
type Service interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
	ConfirmOrder(ctx context.Context, id, providerRef string) error
	CancelOrder(ctx context.Context, id, reason string) error
}

type CreateOrderInput struct {
	UserID      string `json:"user_id" validate:"required"`
	AmountCents uint64 `json:"amount_cents" validate:"min=1"`
	Description string `json:"description"`
}

type UseCase struct {
	repo      domain.Repository
	publisher contracts.EventPublisher
}

func NewUseCase(repo domain.Repository, publisher contracts.EventPublisher) *UseCase {
	return &UseCase{repo: repo, publisher: publisher}
}

// CreateOrder persiste um novo pedido e publica PaymentReceived para iniciar a saga.
func (uc *UseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	if input.AmountCents == 0 {
		return nil, domain.ErrInvalidAmount
	}

	sagaID := common.GenID().Value

	order, err := uc.repo.Create(ctx, domain.Order{
		SagaID:      sagaID,
		UserID:      input.UserID,
		AmountCents: input.AmountCents,
		Status:      domain.Pending,
		Description: input.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("order: falha ao criar pedido: %w", err)
	}

	event := contracts.NewPaymentReceived(sagaID, order.ID, input.UserID, input.AmountCents, input.Description)
	if err := uc.publisher.Publish(ctx, contracts.TopicPaymentReceived, event); err != nil {
		return nil, fmt.Errorf("order: falha ao publicar PaymentReceived: %w", err)
	}

	return order, nil
}

func (uc *UseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return uc.repo.FindByID(ctx, id)
}

// ConfirmOrder é chamado quando o evento OrderConfirmed chega do saga.
func (uc *UseCase) ConfirmOrder(ctx context.Context, id, _ string) error {
	_, err := uc.repo.UpdateStatus(ctx, id, domain.Confirmed)
	return err
}

// CancelOrder é chamado quando o evento OrderCancelled chega do saga.
func (uc *UseCase) CancelOrder(ctx context.Context, id, _ string) error {
	_, err := uc.repo.UpdateStatus(ctx, id, domain.Cancelled)
	return err
}
