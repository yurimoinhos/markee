package infrastructure

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
	entstate "github.com/aggi-tech/aggipay/ent/sagastate"
	"github.com/aggi-tech/aggipay/modules/payment/domain"
	"github.com/aggi-tech/aggipay/platform/common"
)

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) *EntRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Create(ctx context.Context, tx domain.PaymentTx) (*domain.PaymentTx, error) {
	row, err := r.client.SagaState.Create().
		SetID(common.GenID().Value).
		SetSagaType("payment_tx").
		SetOrderID(tx.OrderID).
		SetUserID(tx.SagaID).
		SetAmountCents(tx.AmountCents).
		SetCurrentStep("pending").
		SetStatus(string(tx.Status)).
		SetIdempotencyKey(tx.IdempotencyKey).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("payment repo: falha ao criar tx: %w", err)
	}

	return toTx(row), nil
}

func (r *EntRepository) FindBySagaID(ctx context.Context, sagaID string) (*domain.PaymentTx, error) {
	row, err := r.client.SagaState.Query().
		Where(
			entstate.UserID(sagaID),
			entstate.SagaType("payment_tx"),
		).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("payment repo: %w", err)
	}
	return toTx(row), nil
}

func (r *EntRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, providerRef string) error {
	q := r.client.SagaState.UpdateOneID(id).SetStatus(string(status))
	if providerRef != "" {
		q = q.SetProviderRef(providerRef)
	}
	_, err := q.Save(ctx)
	return err
}

func (r *EntRepository) ExistsByIdempotencyKey(ctx context.Context, key string) (bool, error) {
	return r.client.SagaState.Query().
		Where(entstate.IdempotencyKey(key)).
		Exist(ctx)
}

func toTx(row *ent.SagaState) *domain.PaymentTx {
	providerRef := ""
	if row.ProviderRef != nil {
		providerRef = *row.ProviderRef
	}
	return &domain.PaymentTx{
		ID:             row.ID,
		SagaID:         row.UserID,
		OrderID:        row.OrderID,
		AmountCents:    row.AmountCents,
		Status:         domain.PaymentStatus(row.Status),
		ProviderRef:    providerRef,
		IdempotencyKey: row.IdempotencyKey,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}

var _ domain.Repository = (*EntRepository)(nil)
