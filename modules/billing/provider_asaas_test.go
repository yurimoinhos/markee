package billing

import (
	"context"
	"testing"
	"time"

	"github.com/aggi-tech/aggipay/platform/config"
)

func TestAsaasBillingTypeMapping(t *testing.T) {
	cases := map[string]string{
		"pix":           "PIX",
		"boleto":        "BOLETO",
		"credit_card":   "CREDIT_CARD",
		"bank_transfer": "UNDEFINED",
	}
	for in, want := range cases {
		got := asaasBillingType(in)
		if got != want {
			t.Fatalf("asaasBillingType(%q)=%q want=%q", in, got, want)
		}
	}
}

func TestAsaasProviderFallbackWithoutAPIKey(t *testing.T) {
	p := NewAsaasProvider(&config.Config{Asaas: config.AsaasConfig{BaseURL: "https://api.asaas.com/v3", APIKey: ""}})
	externalID, link, qr, err := p.CreateCharge(context.Background(), "cus_1", 12345, "pix", time.Now(), "teste")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if externalID == "" || link == "" || qr == "" {
		t.Fatalf("expected non-empty fallback values, got id=%q link=%q qr=%q", externalID, link, qr)
	}
}
