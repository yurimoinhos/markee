package contracts

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

type ClicksignProvider interface {
	CreateDocument(ctx context.Context, contractID, signerName, signerEmail, content string) (providerDocID, signURL string, err error)
}

type clicksignHTTPProvider struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewClicksignProvider(cfg *config.Config) ClicksignProvider {
	return &clicksignHTTPProvider{
		baseURL: cfg.Clicksign.BaseURL,
		token:   cfg.Clicksign.Token,
		client:  &http.Client{Timeout: 12 * time.Second},
	}
}

func (p *clicksignHTTPProvider) CreateDocument(ctx context.Context, contractID, signerName, signerEmail, content string) (string, string, error) {
	if p.token == "" {
		fake := common.GenID().Value
		return fake, "https://sandbox.clicksign.com/document/" + fake, nil
	}

	payload := map[string]any{
		"document": map[string]any{
			"path":       "/contracts/" + contractID,
			"content":    content,
			"signers":    []map[string]string{{"name": signerName, "email": signerEmail}},
			"auto_close": true,
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/documents?access_token="+p.token, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("clicksign retornou status %d", resp.StatusCode)
	}

	// Minimal parsing to keep provider integration resilient across payload versions.
	var out struct {
		Document struct {
			Key     string `json:"key"`
			SignURL string `json:"sign_url"`
		} `json:"document"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	if out.Document.Key == "" {
		return "", "", fmt.Errorf("clicksign sem key no response")
	}
	return out.Document.Key, out.Document.SignURL, nil
}
