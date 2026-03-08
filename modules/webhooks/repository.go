package webhooks

import (
	"context"
	"errors"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entcharge "github.com/aggi-tech/aggipay/ent/charge"
	entcontractsignature "github.com/aggi-tech/aggipay/ent/contractsignature"
	entpaymentrecord "github.com/aggi-tech/aggipay/ent/paymentrecord"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) SaveWebhookEvent(ctx context.Context, provider, eventID, signature, payload, status string) (bool, error) {
	_, err := r.client.IntegrationWebhookEvent.Create().
		SetProvider(provider).
		SetEventID(eventID).
		SetNillableSignature(optionalString(signature)).
		SetPayload(payload).
		SetStatus(status).
		SetProcessedAt(time.Now().UTC()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) MarkChargePaidByExternal(ctx context.Context, externalID, method string, paidAt time.Time, txHash *string) error {
	chargeRow, err := r.client.Charge.Query().Where(entcharge.ExternalID(externalID)).First(ctx)
	if err != nil {
		return err
	}

	_, err = r.client.Charge.UpdateOneID(chargeRow.ID).
		SetStatus("paid").
		SetPaidAt(paidAt).
		Save(ctx)
	if err != nil {
		return err
	}

	existing, err := r.client.PaymentRecord.Query().Where(entpaymentrecord.ExternalID(externalID)).First(ctx)
	if err == nil {
		_, err = r.client.PaymentRecord.UpdateOneID(existing.ID).
			SetStatus("paid").
			SetPaidAt(paidAt).
			SetNillableTxHash(txHash).
			Save(ctx)
		return err
	}
	var notFound *ent.NotFoundError
	if !errors.As(err, &notFound) {
		return err
	}

	_, err = r.client.PaymentRecord.Create().
		SetChargeID(chargeRow.ID).
		SetAmountCents(chargeRow.AmountCents).
		SetMethod(method).
		SetStatus("paid").
		SetPaidAt(paidAt).
		SetNillableDueDate(&chargeRow.DueDate).
		SetNillableExternalID(&externalID).
		SetNillableTxHash(txHash).
		Save(ctx)
	return err
}

func (r *Repository) MarkSignatureSigned(ctx context.Context, providerDocID, ip, payload string, acceptedAt time.Time) error {
	sig, err := r.client.ContractSignature.Query().Where(entcontractsignature.ProviderDocID(providerDocID)).First(ctx)
	if err != nil {
		return err
	}
	_, err = r.client.ContractSignature.UpdateOneID(sig.ID).
		SetStatus("signed").
		SetAcceptedAt(acceptedAt).
		SetAcceptedIP(ip).
		SetPayloadJSON(payload).
		Save(ctx)
	return err
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
