package projects

import (
	"context"

	"github.com/aggi-tech/aggipay/ent"
	entproject "github.com/aggi-tech/aggipay/ent/project"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) CreateProject(ctx context.Context, ownerUserID string, in CreateProjectInput) (*ent.Project, error) {
	return r.client.Project.Create().
		SetOwnerUserID(ownerUserID).
		SetContractID(in.ContractID).
		SetName(in.Name).
		SetStatus("active").
		Save(ctx)
}

func (r *Repository) ListProjects(ctx context.Context, ownerUserID string) ([]*ent.Project, error) {
	return r.client.Project.Query().
		Where(entproject.OwnerUserID(ownerUserID)).
		Order(ent.Desc(entproject.FieldCreatedAt)).
		All(ctx)
}

func (r *Repository) GetProject(ctx context.Context, ownerUserID, projectID string) (*ent.Project, error) {
	return r.client.Project.Query().
		Where(entproject.ID(projectID), entproject.OwnerUserID(ownerUserID)).
		First(ctx)
}

func (r *Repository) CreateMilestone(ctx context.Context, projectID string, in CreateMilestoneInput) (*ent.Milestone, error) {
	b := r.client.Milestone.Create().
		SetProjectID(projectID).
		SetTitle(in.Title).
		SetNillableContractID(in.ContractID).
		SetNillableDeliverables(in.Deliverables).
		SetNillableAmountCents(in.AmountCents)
	if in.DueDate != nil {
		b.SetDueDate(*in.DueDate)
	}
	return b.Save(ctx)
}

func (r *Repository) CreateWorklog(ctx context.Context, projectID, userID string, in CreateWorklogInput) (*ent.Worklog, error) {
	b := r.client.Worklog.Create().
		SetProjectID(projectID).
		SetUserID(userID).
		SetHours(in.Hours).
		SetNillableMilestoneID(in.MilestoneID).
		SetNillableDescription(in.Description)
	if in.WorkedAt != nil {
		b.SetWorkedAt(*in.WorkedAt)
	}
	return b.Save(ctx)
}
