package billing

import (
	"context"
	"time"

	"github.com/aggi-tech/aggipay/ent"
)

type Service struct {
	repo     *Repository
	provider AsaasProvider
}

func NewService(repo *Repository, provider AsaasProvider) *Service {
	return &Service{repo: repo, provider: provider}
}

func (s *Service) Create(ctx context.Context, ownerUserID string, in CreateChargeInput) (*ent.Charge, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	charge, err := s.repo.CreateCharge(ctx, ownerUserID, in)
	if err != nil {
		return nil, err
	}
	due := charge.DueDate
	desc := "Cobrança " + charge.ID
	if charge.Description != nil {
		desc = *charge.Description
	}
	externalID, link, qr, err := s.provider.CreateCharge(ctx, charge.CustomerID, charge.AmountCents, charge.PaymentMethod, due, desc)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateExternalData(ctx, charge.ID, externalID, link, qr)
}

func (s *Service) List(ctx context.Context, ownerUserID string) ([]*ent.Charge, error) {
	_, _ = s.repo.SetOverdueByDate(ctx, time.Now())
	return s.repo.ListByOwner(ctx, ownerUserID)
}

func (s *Service) PaymentLink(ctx context.Context, ownerUserID, chargeID string) (*PaymentLinkResponse, error) {
	charge, err := s.repo.GetByID(ctx, ownerUserID, chargeID)
	if err != nil {
		return nil, err
	}
	link := ""
	if charge.PaymentLink != nil {
		link = *charge.PaymentLink
	}
	return &PaymentLinkResponse{ChargeID: charge.ID, Link: link}, nil
}

func (s *Service) PaymentQR(ctx context.Context, ownerUserID, chargeID string) (*PaymentQRResponse, error) {
	charge, err := s.repo.GetByID(ctx, ownerUserID, chargeID)
	if err != nil {
		return nil, err
	}
	qr := ""
	if charge.QrCode != nil {
		qr = *charge.QrCode
	}
	return &PaymentQRResponse{ChargeID: charge.ID, QRCode: qr}, nil
}
