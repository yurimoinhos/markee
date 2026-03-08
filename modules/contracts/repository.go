package contracts

import (
	"context"
	"strconv"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entcontractsignature "github.com/aggi-tech/aggipay/ent/contractsignature"
	entcontractversion "github.com/aggi-tech/aggipay/ent/contractversion"
	entcontract "github.com/aggi-tech/aggipay/ent/servicecontract"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) CreateContract(ctx context.Context, ownerUserID string, in CreateContractInput) (*ent.ServiceContract, error) {
	builder := r.client.ServiceContract.Create().
		SetOwnerUserID(ownerUserID).
		SetCustomerID(in.CustomerID).
		SetContractType(in.ContractType).
		SetTitle(in.Title).
		SetAmountCents(in.AmountCents).
		SetBillingType(in.BillingType).
		SetAutoRenew(in.AutoRenew)

	if in.StartDate != nil {
		builder.SetStartDate(*in.StartDate)
	}
	builder.SetNillableEndDate(in.EndDate)
	builder.SetNillableDurationMonths(in.DurationMonths)
	builder.SetNillableDeliverables(in.Deliverables)
	builder.SetNillableSLA(in.SLA)
	builder.SetNillablePenalties(in.Penalties)
	builder.SetNillablePaymentTerms(in.PaymentTerms)

	return builder.Save(ctx)
}

func (r *Repository) ActivateContract(ctx context.Context, contractID string) error {
	_, err := r.client.ServiceContract.UpdateOneID(contractID).SetStatus("active").Save(ctx)
	return err
}

func (r *Repository) GetContract(ctx context.Context, ownerUserID, contractID string) (*ent.ServiceContract, error) {
	return r.client.ServiceContract.Query().
		Where(entcontract.ID(contractID), entcontract.OwnerUserID(ownerUserID)).
		First(ctx)
}

func (r *Repository) ListContracts(ctx context.Context, ownerUserID string) ([]*ent.ServiceContract, error) {
	return r.client.ServiceContract.Query().
		Where(entcontract.OwnerUserID(ownerUserID)).
		Order(ent.Desc(entcontract.FieldCreatedAt)).
		All(ctx)
}

func (r *Repository) CreateVersion(ctx context.Context, contractID string, nextVersion int, in GenerateContractInput) (*ent.ContractVersion, error) {
	v := strconv.Itoa(nextVersion)
	pdfURL := "https://files.markee.local/contracts/" + contractID + "/v" + v + ".pdf"
	editableURL := "https://files.markee.local/contracts/" + contractID + "/v" + v + ".md"

	return r.client.ContractVersion.Create().
		SetContractID(contractID).
		SetVersion(nextVersion).
		SetTemplateName(in.TemplateName).
		SetEditableContent(in.EditableContent).
		SetPdfURL(pdfURL).
		SetEditableURL(editableURL).
		Save(ctx)
}

func (r *Repository) LastVersion(ctx context.Context, contractID string) (*ent.ContractVersion, error) {
	return r.client.ContractVersion.Query().
		Where(entcontractversion.ContractID(contractID)).
		Order(ent.Desc(entcontractversion.FieldVersion)).
		First(ctx)
}

func (r *Repository) CreateSignature(ctx context.Context, contractID string, versionID *string, signerName, signerEmail, providerDocID string, payload string) (*ent.ContractSignature, error) {
	return r.client.ContractSignature.Create().
		SetContractID(contractID).
		SetNillableContractVersionID(versionID).
		SetSignerName(signerName).
		SetSignerEmail(signerEmail).
		SetProvider("clicksign").
		SetProviderDocID(providerDocID).
		SetPayloadJSON(payload).
		SetStatus("pending").
		Save(ctx)
}

func (r *Repository) FindSignatureByProviderDocID(ctx context.Context, providerDocID string) (*ent.ContractSignature, error) {
	return r.client.ContractSignature.Query().
		Where(entcontractsignature.ProviderDocID(providerDocID)).
		First(ctx)
}

func (r *Repository) MarkSignatureSigned(ctx context.Context, signatureID, ip string, acceptedAt time.Time, payload string) error {
	_, err := r.client.ContractSignature.UpdateOneID(signatureID).
		SetStatus("signed").
		SetAcceptedAt(acceptedAt).
		SetAcceptedIP(ip).
		SetPayloadJSON(payload).
		Save(ctx)
	return err
}
