package billing

import (
	"strings"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type CreateChargeInput struct {
	CustomerID    string     `json:"customer_id"`
	ContractID    *string    `json:"contract_id"`
	MilestoneID   *string    `json:"milestone_id"`
	ChargeType    string     `json:"charge_type"`
	AmountCents   uint64     `json:"amount_cents"`
	PaymentMethod string     `json:"payment_method"`
	DueDate       *time.Time `json:"due_date"`
	Description   *string    `json:"description"`
}

func (in CreateChargeInput) Validate() error {
	if strings.TrimSpace(in.CustomerID) == "" {
		return problem.BadRequest("customer_id é obrigatório")
	}
	if in.AmountCents == 0 {
		return problem.BadRequest("amount_cents deve ser maior que zero")
	}
	if strings.TrimSpace(in.ChargeType) == "" {
		in.ChargeType = "monthly"
	}
	if strings.TrimSpace(in.PaymentMethod) == "" {
		in.PaymentMethod = "pix"
	}
	return nil
}

type PaymentLinkResponse struct {
	ChargeID string `json:"charge_id"`
	Link     string `json:"link"`
}

type PaymentQRResponse struct {
	ChargeID string `json:"charge_id"`
	QRCode   string `json:"qr_code"`
}
