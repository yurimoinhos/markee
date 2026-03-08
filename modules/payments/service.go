package payments

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

func (s *Service) Confirm(ctx context.Context, ownerUserID string, in ConfirmPaymentInput) (*ent.PaymentRecord, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Confirm(ctx, ownerUserID, in)
}

func (s *Service) AddEvidence(ctx context.Context, paymentID string, in EvidenceInput) (*ent.PaymentEvidence, error) {
	return s.repo.AddEvidence(ctx, paymentID, in)
}

func (s *Service) List(ctx context.Context, ownerUserID string) ([]*ent.PaymentRecord, error) {
	return s.repo.ListByOwner(ctx, ownerUserID)
}
