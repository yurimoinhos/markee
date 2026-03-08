package customers

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/ent/charge"
	"github.com/aggi-tech/aggipay/ent/customer"
	"github.com/aggi-tech/aggipay/ent/servicecontract"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Create(ctx context.Context, ownerUserID string, in CreateCustomerInput) (*ent.Customer, error) {
	return r.client.Customer.Create().
		SetOwnerUserID(ownerUserID).
		SetName(in.Name).
		SetCpfCnpj(in.CpfCnpj).
		SetEmail(in.Email).
		SetPreferredPaymentMethod(in.PreferredPaymentMethod).
		SetNillableCompany(in.Company).
		SetNillablePhone(in.Phone).
		SetNillableAddress(in.Address).
		Save(ctx)
}

func (r *Repository) Update(ctx context.Context, ownerUserID, id string, in UpdateCustomerInput) (*ent.Customer, error) {
	builder := r.client.Customer.UpdateOneID(id)
	if in.Name != nil {
		builder.SetName(*in.Name)
	}
	if in.Email != nil {
		builder.SetEmail(*in.Email)
	}
	if in.Company != nil {
		builder.SetCompany(*in.Company)
	}
	if in.Phone != nil {
		builder.SetPhone(*in.Phone)
	}
	if in.Address != nil {
		builder.SetAddress(*in.Address)
	}
	if in.PreferredPaymentMethod != nil {
		builder.SetPreferredPaymentMethod(*in.PreferredPaymentMethod)
	}
	entity, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	if entity.OwnerUserID != ownerUserID {
		return nil, fmt.Errorf("customer não pertence ao usuário autenticado")
	}
	return entity, nil
}

func (r *Repository) GetByID(ctx context.Context, ownerUserID, id string) (*ent.Customer, error) {
	return r.client.Customer.Query().
		Where(customer.ID(id), customer.OwnerUserID(ownerUserID)).
		First(ctx)
}

func (r *Repository) ListByOwner(ctx context.Context, ownerUserID string) ([]*ent.Customer, error) {
	return r.client.Customer.Query().
		Where(customer.OwnerUserID(ownerUserID)).
		Order(ent.Desc(customer.FieldCreatedAt)).
		All(ctx)
}

func (r *Repository) Summary(ctx context.Context, ownerUserID, customerID string) (CustomerFinancialSummary, error) {
	contractsCount, err := r.client.ServiceContract.Query().Where(
		servicecontract.OwnerUserID(ownerUserID),
		servicecontract.CustomerID(customerID),
	).Count(ctx)
	if err != nil {
		return CustomerFinancialSummary{}, err
	}

	pendingCount, err := r.client.Charge.Query().Where(
		charge.OwnerUserID(ownerUserID),
		charge.CustomerID(customerID),
		charge.StatusIn("pending"),
	).Count(ctx)
	if err != nil {
		return CustomerFinancialSummary{}, err
	}

	paidCount, err := r.client.Charge.Query().Where(
		charge.OwnerUserID(ownerUserID),
		charge.CustomerID(customerID),
		charge.StatusIn("paid"),
	).Count(ctx)
	if err != nil {
		return CustomerFinancialSummary{}, err
	}

	overdueCount, err := r.client.Charge.Query().Where(
		charge.OwnerUserID(ownerUserID),
		charge.CustomerID(customerID),
		charge.StatusIn("overdue"),
	).Count(ctx)
	if err != nil {
		return CustomerFinancialSummary{}, err
	}

	return CustomerFinancialSummary{
		ContractsCount: contractsCount,
		PendingCharges: pendingCount,
		PaidCharges:    paidCount,
		OverdueCharges: overdueCount,
	}, nil
}
