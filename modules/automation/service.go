package automation

import (
	"context"
	"time"

	"github.com/aggi-tech/aggipay/ent"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Run(ctx context.Context, ownerUserID string) (RunResult, error) {
	return s.repo.Run(ctx, ownerUserID, time.Now().UTC())
}

func (s *Service) ListRuns(ctx context.Context, ownerUserID string) ([]*ent.AutomationRun, error) {
	return s.repo.ListRuns(ctx, ownerUserID)
}
