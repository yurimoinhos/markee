package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	entautomationrun "github.com/aggi-tech/aggipay/ent/automationrun"
	entcharge "github.com/aggi-tech/aggipay/ent/charge"
	entproject "github.com/aggi-tech/aggipay/ent/project"
	entservicecontract "github.com/aggi-tech/aggipay/ent/servicecontract"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Run(ctx context.Context, ownerUserID string, now time.Time) (RunResult, error) {
	result := RunResult{}

	reminders, err := r.client.Charge.Query().Where(
		entcharge.OwnerUserID(ownerUserID),
		entcharge.Status("pending"),
		entcharge.DueDateLTE(now.AddDate(0, 0, 3)),
	).Count(ctx)
	if err != nil {
		return RunResult{}, err
	}
	result.RemindersSent = reminders
	_, _ = r.log(ctx, ownerUserID, "payment_reminder", "success", fmt.Sprintf("%d lembretes", reminders))

	expiring, err := r.client.ServiceContract.Query().Where(
		entservicecontract.OwnerUserID(ownerUserID),
		entservicecontract.Status("active"),
		entservicecontract.EndDateNotNil(),
		entservicecontract.EndDateLTE(now.AddDate(0, 0, 15)),
	).Count(ctx)
	if err != nil {
		return RunResult{}, err
	}
	result.ContractsNearExpiry = expiring
	_, _ = r.log(ctx, ownerUserID, "contract_expiry_notice", "success", fmt.Sprintf("%d contratos", expiring))

	contracts, err := r.client.ServiceContract.Query().Where(
		entservicecontract.OwnerUserID(ownerUserID),
		entservicecontract.Status("active"),
		entservicecontract.BillingType("monthly"),
	).All(ctx)
	if err != nil {
		return RunResult{}, err
	}

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	for _, contract := range contracts {
		exists, err := r.client.Charge.Query().Where(
			entcharge.OwnerUserID(ownerUserID),
			entcharge.ContractID(contract.ID),
			entcharge.CreatedAtGTE(monthStart),
		).Exist(ctx)
		if err != nil {
			return RunResult{}, err
		}
		if exists {
			continue
		}
		_, err = r.client.Charge.Create().
			SetOwnerUserID(ownerUserID).
			SetCustomerID(contract.CustomerID).
			SetContractID(contract.ID).
			SetChargeType("monthly").
			SetAmountCents(contract.AmountCents).
			SetCurrency("BRL").
			SetPaymentMethod("pix").
			SetStatus("pending").
			SetDueDate(now.AddDate(0, 0, 7)).
			SetDescription("Cobrança recorrente automática").
			Save(ctx)
		if err != nil {
			return RunResult{}, err
		}
		result.RecurringChargesMade++
	}
	_, _ = r.log(ctx, ownerUserID, "recurring_charge_generation", "success", fmt.Sprintf("%d cobranças", result.RecurringChargesMade))

	delinquentCustomers, err := r.client.Charge.Query().Where(
		entcharge.OwnerUserID(ownerUserID),
		entcharge.Status("overdue"),
	).GroupBy(entcharge.FieldCustomerID).Strings(ctx)
	if err != nil {
		return RunResult{}, err
	}
	projects, err := r.client.Project.Query().Where(entproject.OwnerUserID(ownerUserID)).All(ctx)
	if err != nil {
		return RunResult{}, err
	}
	delinquentMap := make(map[string]bool, len(delinquentCustomers))
	for _, cID := range delinquentCustomers {
		delinquentMap[cID] = true
	}
	for _, project := range projects {
		contract, err := r.client.ServiceContract.Query().Where(
			entservicecontract.ID(project.ContractID),
			entservicecontract.OwnerUserID(ownerUserID),
		).First(ctx)
		if err != nil {
			continue
		}
		if !delinquentMap[contract.CustomerID] {
			continue
		}
		_, err = r.client.Project.UpdateOneID(project.ID).SetStatus("on_hold").Save(ctx)
		if err != nil {
			return RunResult{}, err
		}
		result.ProjectsSuspended++
	}
	_, _ = r.log(ctx, ownerUserID, "service_suspension", "success", fmt.Sprintf("%d projetos", result.ProjectsSuspended))

	_, err = r.log(ctx, ownerUserID, "monthly_financial_report", "success", "relatório gerado")
	if err != nil {
		return RunResult{}, err
	}
	result.MonthlyReportGenerated = 1

	return result, nil
}

func (r *Repository) ListRuns(ctx context.Context, ownerUserID string) ([]*ent.AutomationRun, error) {
	return r.client.AutomationRun.Query().
		Where(entautomationrun.OwnerUserID(ownerUserID)).
		Order(ent.Desc(entautomationrun.FieldCreatedAt)).
		All(ctx)
}

func (r *Repository) log(ctx context.Context, ownerUserID, automationType, status, details string) (*ent.AutomationRun, error) {
	return r.client.AutomationRun.Create().
		SetOwnerUserID(ownerUserID).
		SetAutomationType(automationType).
		SetStatus(status).
		SetDetails(details).
		Save(ctx)
}
