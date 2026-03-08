package finance

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Dashboard(ctx context.Context, ownerUserID string) (Dashboard, error) {
	return s.repo.Dashboard(ctx, ownerUserID, time.Now().UTC())
}

func (s *Service) CashFlow(ctx context.Context, ownerUserID string) ([]CashFlowPoint, error) {
	return s.repo.CashFlow(ctx, ownerUserID, time.Now().UTC())
}

func (s *Service) Defaults(ctx context.Context, ownerUserID string) (DefaultMetrics, error) {
	return s.repo.Defaults(ctx, ownerUserID, time.Now().UTC())
}
