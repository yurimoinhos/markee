package payments

import (
	"strings"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type ConfirmPaymentInput struct {
	ChargeID    string     `json:"charge_id"`
	AmountCents uint64     `json:"amount_cents"`
	Method      string     `json:"method"`
	PaidAt      *time.Time `json:"paid_at"`
	ReceiptURL  *string    `json:"receipt_url"`
	TxHash      *string    `json:"tx_hash"`
}

type EvidenceInput struct {
	FileURL *string `json:"file_url"`
	Note    *string `json:"note"`
	TxHash  *string `json:"tx_hash"`
}

func (in ConfirmPaymentInput) Validate() error {
	if strings.TrimSpace(in.ChargeID) == "" {
		return problem.BadRequest("charge_id é obrigatório")
	}
	if in.AmountCents == 0 {
		return problem.BadRequest("amount_cents deve ser maior que zero")
	}
	if strings.TrimSpace(in.Method) == "" {
		in.Method = "pix"
	}
	return nil
}
