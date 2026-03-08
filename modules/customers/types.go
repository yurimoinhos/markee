package customers

import (
	"strings"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type CreateCustomerInput struct {
	Name                   string  `json:"name"`
	Company                *string `json:"company"`
	CpfCnpj                string  `json:"cpf_cnpj"`
	Email                  string  `json:"email"`
	Phone                  *string `json:"phone"`
	Address                *string `json:"address"`
	PreferredPaymentMethod string  `json:"preferred_payment_method"`
}

type UpdateCustomerInput struct {
	Name                   *string `json:"name"`
	Company                *string `json:"company"`
	Email                  *string `json:"email"`
	Phone                  *string `json:"phone"`
	Address                *string `json:"address"`
	PreferredPaymentMethod *string `json:"preferred_payment_method"`
}

type CustomerFinancialSummary struct {
	ContractsCount int `json:"contracts_count"`
	PendingCharges int `json:"pending_charges"`
	PaidCharges    int `json:"paid_charges"`
	OverdueCharges int `json:"overdue_charges"`
}

func (in CreateCustomerInput) Validate() error {
	if strings.TrimSpace(in.Name) == "" {
		return problem.BadRequest("name é obrigatório")
	}
	if strings.TrimSpace(in.CpfCnpj) == "" {
		return problem.BadRequest("cpf_cnpj é obrigatório")
	}
	if strings.TrimSpace(in.Email) == "" {
		return problem.BadRequest("email é obrigatório")
	}
	if strings.TrimSpace(in.PreferredPaymentMethod) == "" {
		in.PreferredPaymentMethod = "pix"
	}
	return nil
}
