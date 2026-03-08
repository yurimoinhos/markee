package projects

import (
	"strings"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type CreateProjectInput struct {
	ContractID string `json:"contract_id"`
	Name       string `json:"name"`
}

type CreateMilestoneInput struct {
	ContractID   *string    `json:"contract_id"`
	Title        string     `json:"title"`
	Deliverables *string    `json:"deliverables"`
	AmountCents  *uint64    `json:"amount_cents"`
	DueDate      *time.Time `json:"due_date"`
}

type CreateWorklogInput struct {
	MilestoneID *string    `json:"milestone_id"`
	Hours       float64    `json:"hours"`
	Description *string    `json:"description"`
	WorkedAt    *time.Time `json:"worked_at"`
}

func (in CreateProjectInput) Validate() error {
	if strings.TrimSpace(in.ContractID) == "" {
		return problem.BadRequest("contract_id é obrigatório")
	}
	if strings.TrimSpace(in.Name) == "" {
		return problem.BadRequest("name é obrigatório")
	}
	return nil
}

func (in CreateMilestoneInput) Validate() error {
	if strings.TrimSpace(in.Title) == "" {
		return problem.BadRequest("title é obrigatório")
	}
	return nil
}

func (in CreateWorklogInput) Validate() error {
	if in.Hours <= 0 {
		return problem.BadRequest("hours deve ser maior que zero")
	}
	return nil
}
