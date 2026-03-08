package billing

import (
	"context"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entcharge "github.com/aggi-tech/aggipay/ent/charge"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) CreateCharge(ctx context.Context, ownerUserID string, in CreateChargeInput) (*ent.Charge, error) {
	builder := r.client.Charge.Create().
		SetOwnerUserID(ownerUserID).
		SetCustomerID(in.CustomerID).
		SetNillableContractID(in.ContractID).
		SetNillableMilestoneID(in.MilestoneID).
		SetChargeType(in.ChargeType).
		SetAmountCents(in.AmountCents).
		SetPaymentMethod(in.PaymentMethod).
		SetCurrency("BRL").
		SetStatus("pending").
		SetNillableDescription(in.Description)

	if in.DueDate != nil {
		builder.SetDueDate(*in.DueDate)
	}
	return builder.Save(ctx)
}

func (r *Repository) UpdateExternalData(ctx context.Context, chargeID, externalID, link, qr string) (*ent.Charge, error) {
	return r.client.Charge.UpdateOneID(chargeID).
		SetExternalID(externalID).
		SetNillablePaymentLink(optionalString(link)).
		SetNillableQrCode(optionalString(qr)).
		Save(ctx)
}

func (r *Repository) ListByOwner(ctx context.Context, ownerUserID string) ([]*ent.Charge, error) {
	return r.client.Charge.Query().
		Where(entcharge.OwnerUserID(ownerUserID)).
		Order(ent.Desc(entcharge.FieldCreatedAt)).
		All(ctx)
}

func (r *Repository) GetByID(ctx context.Context, ownerUserID, chargeID string) (*ent.Charge, error) {
	return r.client.Charge.Query().
		Where(entcharge.ID(chargeID), entcharge.OwnerUserID(ownerUserID)).
		First(ctx)
}

func (r *Repository) MarkPaid(ctx context.Context, chargeID string, paidAt time.Time) error {
	_, err := r.client.Charge.UpdateOneID(chargeID).
		SetStatus("paid").
		SetPaidAt(paidAt).
		Save(ctx)
	return err
}

func (r *Repository) SetOverdueByDate(ctx context.Context, now time.Time) (int, error) {
	return r.client.Charge.Update().
		Where(entcharge.Status("pending"), entcharge.DueDateLT(now)).
		SetStatus("overdue").
		Save(ctx)
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
