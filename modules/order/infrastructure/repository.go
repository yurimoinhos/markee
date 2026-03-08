package infrastructure

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
	entorder "github.com/aggi-tech/aggipay/ent/order"
	"github.com/aggi-tech/aggipay/modules/order/domain"
	"github.com/aggi-tech/aggipay/platform/common"
)

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) *EntRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Create(ctx context.Context, o domain.Order) (*domain.Order, error) {
	row, err := r.client.Order.Create().
		SetID(common.GenID().Value).
		SetSagaID(o.SagaID).
		SetUserID(o.UserID).
		SetAmount(o.AmountCents).
		SetStatus(string(o.Status)).
		SetNillableDescription(nilIfEmpty(o.Description)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("order repo: falha ao criar: %w", err)
	}
	return toOrder(row), nil
}

func (r *EntRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	row, err := r.client.Order.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("order repo: %w", err)
	}
	return toOrder(row), nil
}

func (r *EntRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, error) {
	row, err := r.client.Order.UpdateOneID(id).
		SetStatus(string(status)).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("order repo: falha ao atualizar status: %w", err)
	}
	return toOrder(row), nil
}

func toOrder(row *ent.Order) *domain.Order {
	return &domain.Order{
		ID:          row.ID,
		SagaID:      row.SagaID,
		UserID:      row.UserID,
		AmountCents: row.Amount,
		Status:      domain.OrderStatus(row.Status),
		Description: stringVal(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var _ domain.Repository = (*EntRepository)(nil)

// Garante import do pacote order (usado para predicados de query futuros)
var _ = entorder.FieldID
