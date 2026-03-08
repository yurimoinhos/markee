package domain

import (
	"context"
	"time"

	"net/http"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Confirmed OrderStatus = "confirmed"
	Cancelled OrderStatus = "cancelled"
)

// Order representa um pedido de pagamento.
type Order struct {
	ID          string
	SagaID      string
	UserID      string
	AmountCents uint64
	Status      OrderStatus
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository define o contrato de persistência do módulo order.
type Repository interface {
	Create(ctx context.Context, o Order) (*Order, error)
	FindByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) (*Order, error)
}

var ErrOrderNotFound = problem.New(
	http.StatusNotFound,
	"/problems/order/not-found",
	"Pedido não encontrado",
	"O pedido solicitado não existe.",
)

var ErrInvalidAmount = problem.New(
	http.StatusBadRequest,
	"/problems/order/invalid-amount",
	"Valor inválido",
	"O valor do pedido deve ser maior que zero.",
)
