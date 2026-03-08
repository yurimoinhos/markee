package payments

import (
	"context"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entcharge "github.com/aggi-tech/aggipay/ent/charge"
	entpaymentrecord "github.com/aggi-tech/aggipay/ent/paymentrecord"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Confirm(ctx context.Context, ownerUserID string, in ConfirmPaymentInput) (*ent.PaymentRecord, error) {
	chargeRow, err := r.client.Charge.Query().Where(entcharge.ID(in.ChargeID), entcharge.OwnerUserID(ownerUserID)).First(ctx)
	if err != nil {
		return nil, err
	}

	paidAt := time.Now().UTC()
	if in.PaidAt != nil {
		paidAt = *in.PaidAt
	}

	row, err := r.client.PaymentRecord.Create().
		SetChargeID(chargeRow.ID).
		SetAmountCents(in.AmountCents).
		SetMethod(in.Method).
		SetStatus("paid").
		SetPaidAt(paidAt).
		SetNillableDueDate(&chargeRow.DueDate).
		SetNillableReceiptURL(in.ReceiptURL).
		SetNillableTxHash(in.TxHash).
		SetNillableExternalID(chargeRow.ExternalID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.client.Charge.UpdateOneID(chargeRow.ID).
		SetStatus("paid").
		SetPaidAt(paidAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *Repository) AddEvidence(ctx context.Context, paymentID string, in EvidenceInput) (*ent.PaymentEvidence, error) {
	return r.client.PaymentEvidence.Create().
		SetPaymentID(paymentID).
		SetNillableFileURL(in.FileURL).
		SetNillableNote(in.Note).
		SetNillableTxHash(in.TxHash).
		Save(ctx)
}

func (r *Repository) ListByOwner(ctx context.Context, ownerUserID string) ([]*ent.PaymentRecord, error) {
	charges, err := r.client.Charge.Query().Where(entcharge.OwnerUserID(ownerUserID)).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(charges) == 0 {
		return []*ent.PaymentRecord{}, nil
	}
	chargeIDs := make([]string, 0, len(charges))
	for _, c := range charges {
		chargeIDs = append(chargeIDs, c.ID)
	}
	return r.client.PaymentRecord.Query().
		Where(entpaymentrecord.ChargeIDIn(chargeIDs...)).
		Order(ent.Desc(entpaymentrecord.FieldCreatedAt)).
		All(ctx)
}
