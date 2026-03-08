package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// PaymentTx representa uma transação de pagamento processada por este módulo.
type PaymentTx struct {
	ID             string
	SagaID         string
	OrderID        string
	AmountCents    uint64
	Status         PaymentStatus
	ProviderRef    string
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Repository define o contrato de persistência do módulo payment.
type Repository interface {
	Create(ctx context.Context, tx PaymentTx) (*PaymentTx, error)
	FindBySagaID(ctx context.Context, sagaID string) (*PaymentTx, error)
	UpdateStatus(ctx context.Context, id string, status PaymentStatus, providerRef string) error
	ExistsByIdempotencyKey(ctx context.Context, key string) (bool, error)
}

var ErrPaymentNotFound = problem.New(
	http.StatusNotFound,
	"/problems/payment/not-found",
	"Pagamento não encontrado",
	"A transação de pagamento solicitada não existe.",
)

var ErrDuplicatePayment = problem.New(
	http.StatusConflict,
	"/problems/payment/duplicate",
	"Pagamento duplicado",
	"Já existe uma transação com esta chave de idempotência.",
)
