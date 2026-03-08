package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aggi-tech/aggipay/platform/common"
	"github.com/aggi-tech/aggipay/platform/config"
)

type AsaasProvider interface {
	CreateCharge(ctx context.Context, customerID string, amountCents uint64, method string, dueDate time.Time, description string) (externalID, paymentLink, qrCode string, err error)
}

type asaasHTTPProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewAsaasProvider(cfg *config.Config) AsaasProvider {
	return &asaasHTTPProvider{
		baseURL: cfg.Asaas.BaseURL,
		apiKey:  cfg.Asaas.APIKey,
		client:  &http.Client{Timeout: 12 * time.Second},
	}
}

func (p *asaasHTTPProvider) CreateCharge(ctx context.Context, customerID string, amountCents uint64, method string, dueDate time.Time, description string) (string, string, string, error) {
	if p.apiKey == "" {
		fakeID := common.GenID().Value
		return fakeID,
			"https://sandbox.asaas.com/i/" + fakeID,
			"PIX|" + fakeID + "|" + fmt.Sprintf("%.2f", float64(amountCents)/100),
			nil
	}

	payload := map[string]any{
		"customer":          customerID,
		"billingType":       asaasBillingType(method),
		"value":             float64(amountCents) / 100,
		"dueDate":           dueDate.Format("2006-01-02"),
		"description":       description,
		"externalReference": common.GenID().Value,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/payments", bytes.NewReader(body))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", "", "", fmt.Errorf("asaas retornou status %d", resp.StatusCode)
	}

	var out struct {
		ID          string `json:"id"`
		InvoiceURL  string `json:"invoiceUrl"`
		BankSlipURL string `json:"bankSlipUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", "", err
	}
	if out.ID == "" {
		return "", "", "", fmt.Errorf("asaas não retornou ID de cobrança")
	}

	paymentLink := out.InvoiceURL
	if paymentLink == "" {
		paymentLink = out.BankSlipURL
	}
	qr := ""
	if method == "pix" {
		qr = "PIX|" + out.ID
	}
	return out.ID, paymentLink, qr, nil
}

func asaasBillingType(method string) string {
	switch method {
	case "boleto":
		return "BOLETO"
	case "credit_card":
		return "CREDIT_CARD"
	case "bank_transfer":
		return "UNDEFINED"
	default:
		return "PIX"
	}
}
