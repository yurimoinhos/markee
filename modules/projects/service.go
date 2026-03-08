package projects

import (
	"context"

	"github.com/aggi-tech/aggipay/ent"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateProject(ctx context.Context, ownerUserID string, in CreateProjectInput) (*ent.Project, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateProject(ctx, ownerUserID, in)
}

func (s *Service) ListProjects(ctx context.Context, ownerUserID string) ([]*ent.Project, error) {
	return s.repo.ListProjects(ctx, ownerUserID)
}

func (s *Service) CreateMilestone(ctx context.Context, ownerUserID, projectID string, in CreateMilestoneInput) (*ent.Milestone, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetProject(ctx, ownerUserID, projectID); err != nil {
		return nil, err
	}
	return s.repo.CreateMilestone(ctx, projectID, in)
}

func (s *Service) CreateWorklog(ctx context.Context, ownerUserID, projectID string, in CreateWorklogInput) (*ent.Worklog, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetProject(ctx, ownerUserID, projectID); err != nil {
		return nil, err
	}
	return s.repo.CreateWorklog(ctx, projectID, ownerUserID, in)
}
