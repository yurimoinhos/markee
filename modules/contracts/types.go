package contracts

import (
	"strings"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type CreateContractInput struct {
	CustomerID     string     `json:"customer_id"`
	ContractType   string     `json:"contract_type"`
	Title          string     `json:"title"`
	AmountCents    uint64     `json:"amount_cents"`
	BillingType    string     `json:"billing_type"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	DurationMonths *int       `json:"duration_months"`
	Deliverables   *string    `json:"deliverables"`
	SLA            *string    `json:"sla"`
	Penalties      *string    `json:"penalties"`
	PaymentTerms   *string    `json:"payment_terms"`
	AutoRenew      bool       `json:"auto_renew"`
}

type GenerateContractInput struct {
	TemplateName    string `json:"template_name"`
	EditableContent string `json:"editable_content"`
}

type SendSignatureInput struct {
	SignerName  string `json:"signer_name"`
	SignerEmail string `json:"signer_email"`
}

func (in CreateContractInput) Validate() error {
	if strings.TrimSpace(in.CustomerID) == "" {
		return problem.BadRequest("customer_id é obrigatório")
	}
	if strings.TrimSpace(in.ContractType) == "" {
		return problem.BadRequest("contract_type é obrigatório")
	}
	if strings.TrimSpace(in.Title) == "" {
		return problem.BadRequest("title é obrigatório")
	}
	if in.AmountCents == 0 {
		return problem.BadRequest("amount_cents deve ser maior que zero")
	}
	if strings.TrimSpace(in.BillingType) == "" {
		in.BillingType = "monthly"
	}
	return nil
}

func (in GenerateContractInput) Validate() error {
	if strings.TrimSpace(in.TemplateName) == "" {
		return problem.BadRequest("template_name é obrigatório")
	}
	if strings.TrimSpace(in.EditableContent) == "" {
		return problem.BadRequest("editable_content é obrigatório")
	}
	return nil
}

func (in SendSignatureInput) Validate() error {
	if strings.TrimSpace(in.SignerName) == "" {
		return problem.BadRequest("signer_name é obrigatório")
	}
	if strings.TrimSpace(in.SignerEmail) == "" {
		return problem.BadRequest("signer_email é obrigatório")
	}
	return nil
}
