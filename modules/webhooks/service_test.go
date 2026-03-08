package webhooks

import "testing"

func TestIsPaidStatus(t *testing.T) {
	for _, status := range []string{"RECEIVED", "CONFIRMED", "RECEIVED_IN_CASH"} {
		if !isPaidStatus(status) {
			t.Fatalf("expected paid status for %s", status)
		}
	}
	for _, status := range []string{"PENDING", "OVERDUE", "FAILED"} {
		if isPaidStatus(status) {
			t.Fatalf("expected non-paid status for %s", status)
		}
	}
}

func TestExtractAsaasPayment(t *testing.T) {
	payload := map[string]any{
		"payment": map[string]any{
			"id":          "pay_123",
			"status":      "received",
			"billingType": "PIX",
			"txHash":      "abc",
		},
	}
	id, status, method, hash := extractAsaasPayment(payload)
	if id != "pay_123" || status != "RECEIVED" || method != "pix" || hash != "abc" {
		t.Fatalf("unexpected extraction: id=%s status=%s method=%s hash=%s", id, status, method, hash)
	}
}
