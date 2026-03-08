package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/problem"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) ProcessAsaas(ctx context.Context, signature string, body []byte) (int, any, error) {
	if s.cfg.Webhooks.AsaasSecret != "" && signature != s.cfg.Webhooks.AsaasSecret {
		return 0, nil, problem.Unauthorized("assinatura do webhook Asaas inválida")
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0, nil, problem.BadRequest("payload inválido")
	}

	eventID := readString(payload, "event_id")
	if eventID == "" {
		eventID = readString(payload, "id")
	}
	if eventID == "" {
		eventID = "asaas-" + fmt.Sprint(time.Now().UnixNano())
	}

	stored, err := s.repo.SaveWebhookEvent(ctx, "asaas", eventID, signature, string(body), "processed")
	if err != nil {
		return 0, nil, err
	}
	if !stored {
		return http.StatusOK, map[string]any{"status": "ignored", "reason": "duplicate event"}, nil
	}

	paymentID, status, method, txHash := extractAsaasPayment(payload)
	if paymentID == "" {
		return http.StatusAccepted, map[string]any{"status": "accepted", "reason": "sem payment id"}, nil
	}
	if isPaidStatus(status) {
		if err := s.repo.MarkChargePaidByExternal(ctx, paymentID, method, time.Now().UTC(), optional(txHash)); err != nil {
			return 0, nil, err
		}
	}

	return http.StatusOK, map[string]any{"status": "processed", "event_id": eventID}, nil
}

func (s *Service) ProcessClicksign(ctx context.Context, signature string, body []byte) (int, any, error) {
	if s.cfg.Webhooks.ClicksignSecret != "" && signature != s.cfg.Webhooks.ClicksignSecret {
		return 0, nil, problem.Unauthorized("assinatura do webhook Clicksign inválida")
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0, nil, problem.BadRequest("payload inválido")
	}

	eventID := readString(payload, "event_id")
	if eventID == "" {
		eventID = readString(payload, "event")
	}
	if eventID == "" {
		eventID = "clicksign-" + fmt.Sprint(time.Now().UnixNano())
	}

	stored, err := s.repo.SaveWebhookEvent(ctx, "clicksign", eventID, signature, string(body), "processed")
	if err != nil {
		return 0, nil, err
	}
	if !stored {
		return http.StatusOK, map[string]any{"status": "ignored", "reason": "duplicate event"}, nil
	}

	providerDocID, signed, ip, acceptedAt := extractClicksign(payload)
	if providerDocID == "" {
		return http.StatusAccepted, map[string]any{"status": "accepted", "reason": "sem provider_doc_id"}, nil
	}
	if signed {
		if err := s.repo.MarkSignatureSigned(ctx, providerDocID, ip, string(body), acceptedAt); err != nil {
			return 0, nil, err
		}
	}
	return http.StatusOK, map[string]any{"status": "processed", "event_id": eventID}, nil
}

func extractAsaasPayment(payload map[string]any) (paymentID, status, method, txHash string) {
	if p, ok := payload["payment"].(map[string]any); ok {
		paymentID = readString(p, "id")
		status = strings.ToUpper(readString(p, "status"))
		method = strings.ToLower(readString(p, "billingType"))
		txHash = readString(p, "txHash")
		return
	}
	paymentID = readString(payload, "id")
	status = strings.ToUpper(readString(payload, "status"))
	method = strings.ToLower(readString(payload, "billingType"))
	txHash = readString(payload, "txHash")
	return
}

func extractClicksign(payload map[string]any) (providerDocID string, signed bool, ip string, acceptedAt time.Time) {
	acceptedAt = time.Now().UTC()
	if doc, ok := payload["document"].(map[string]any); ok {
		providerDocID = readString(doc, "key")
	}
	if providerDocID == "" {
		providerDocID = readString(payload, "document_key")
	}
	status := strings.ToLower(readString(payload, "status"))
	event := strings.ToLower(readString(payload, "event"))
	signed = status == "signed" || strings.Contains(event, "signed")
	ip = readString(payload, "ip")
	if ts := readString(payload, "occurred_at"); ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			acceptedAt = t
		}
	}
	return
}

func readString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func optional(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	v = strings.TrimSpace(v)
	return &v
}

func isPaidStatus(status string) bool {
	switch status {
	case "RECEIVED", "CONFIRMED", "RECEIVED_IN_CASH":
		return true
	default:
		return false
	}
}
