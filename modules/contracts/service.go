package contracts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/platform/problem"
)

type Service struct {
	repo     *Repository
	provider ClicksignProvider
}

func NewService(repo *Repository, provider ClicksignProvider) *Service {
	return &Service{repo: repo, provider: provider}
}

func (s *Service) Create(ctx context.Context, ownerUserID string, in CreateContractInput) (*ent.ServiceContract, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateContract(ctx, ownerUserID, in)
}

func (s *Service) Get(ctx context.Context, ownerUserID, id string) (*ent.ServiceContract, error) {
	return s.repo.GetContract(ctx, ownerUserID, id)
}

func (s *Service) List(ctx context.Context, ownerUserID string) ([]*ent.ServiceContract, error) {
	return s.repo.ListContracts(ctx, ownerUserID)
}

func (s *Service) Generate(ctx context.Context, ownerUserID, contractID string, in GenerateContractInput) (*ent.ContractVersion, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	contract, err := s.repo.GetContract(ctx, ownerUserID, contractID)
	if err != nil {
		return nil, err
	}

	nextVersion := 1
	last, err := s.repo.LastVersion(ctx, contract.ID)
	if err == nil {
		nextVersion = last.Version + 1
	}

	version, err := s.repo.CreateVersion(ctx, contract.ID, nextVersion, in)
	if err != nil {
		return nil, err
	}

	if contract.Status == "draft" {
		_ = s.repo.ActivateContract(ctx, contract.ID)
	}
	return version, nil
}

func (s *Service) SendSignature(ctx context.Context, ownerUserID, contractID string, in SendSignatureInput) (*ent.ContractSignature, string, error) {
	if err := in.Validate(); err != nil {
		return nil, "", err
	}

	contract, err := s.repo.GetContract(ctx, ownerUserID, contractID)
	if err != nil {
		return nil, "", err
	}

	version, err := s.repo.LastVersion(ctx, contract.ID)
	if err != nil {
		return nil, "", problem.BadRequest("gere uma versão do contrato antes de enviar para assinatura")
	}

	providerDocID, signURL, err := s.provider.CreateDocument(ctx, contract.ID, in.SignerName, in.SignerEmail, version.EditableContent)
	if err != nil {
		return nil, "", fmt.Errorf("falha na integração Clicksign: %w", err)
	}

	payload, _ := json.Marshal(map[string]any{
		"provider_doc_id": providerDocID,
		"sign_url":        signURL,
		"template":        version.TemplateName,
	})
	versionID := version.ID
	signature, err := s.repo.CreateSignature(ctx, contract.ID, &versionID, in.SignerName, in.SignerEmail, providerDocID, string(payload))
	if err != nil {
		return nil, "", err
	}

	return signature, signURL, nil
}

func (s *Service) MarkSignedByProvider(ctx context.Context, providerDocID, acceptedIP, payload string, acceptedAt int64) error {
	sig, err := s.repo.FindSignatureByProviderDocID(ctx, providerDocID)
	if err != nil {
		return err
	}
	if sig.Status == "signed" {
		return nil
	}
	at := sig.UpdatedAt
	if acceptedAt > 0 {
		at = time.Unix(acceptedAt, 0).UTC()
	}
	return s.repo.MarkSignatureSigned(ctx, sig.ID, acceptedIP, at, payload)
}
