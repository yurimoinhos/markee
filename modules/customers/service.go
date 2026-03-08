package customers

import (
	"context"
	"fmt"

	"github.com/aggi-tech/aggipay/ent"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, ownerUserID string, in CreateCustomerInput) (*ent.Customer, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, ownerUserID, in)
}

func (s *Service) Update(ctx context.Context, ownerUserID, id string, in UpdateCustomerInput) (*ent.Customer, error) {
	entity, err := s.repo.Update(ctx, ownerUserID, id, in)
	if err != nil {
		return nil, err
	}
	return s.refreshFinancialStatus(ctx, entity)
}

func (s *Service) Get(ctx context.Context, ownerUserID, id string) (*ent.Customer, CustomerFinancialSummary, error) {
	entity, err := s.repo.GetByID(ctx, ownerUserID, id)
	if err != nil {
		return nil, CustomerFinancialSummary{}, err
	}
	summary, err := s.repo.Summary(ctx, ownerUserID, entity.ID)
	if err != nil {
		return nil, CustomerFinancialSummary{}, err
	}
	return entity, summary, nil
}

func (s *Service) List(ctx context.Context, ownerUserID string) ([]*ent.Customer, error) {
	return s.repo.ListByOwner(ctx, ownerUserID)
}

func (s *Service) refreshFinancialStatus(ctx context.Context, entity *ent.Customer) (*ent.Customer, error) {
	summary, err := s.repo.Summary(ctx, entity.OwnerUserID, entity.ID)
	if err != nil {
		return nil, err
	}

	status := "regular"
	if summary.OverdueCharges > 0 {
		status = "delinquent"
	} else if summary.PendingCharges > 3 {
		status = "watch"
	}

	updated, err := s.repo.client.Customer.UpdateOneID(entity.ID).SetFinancialStatus(status).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao atualizar financial_status: %w", err)
	}
	return updated, nil
}
