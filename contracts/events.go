// Package contracts define os eventos, tópicos e interfaces de barramento
// compartilhados entre todos os módulos do AggiPay.
//
// Regra: este pacote não pode importar nenhum outro módulo interno.
// É a única dependência cruzada permitida entre módulos.
package contracts

import "time"

// ── Eventos de Auth ──────────────────────────────────────────────────────────

// UserRegistered é publicado pelo módulo auth após o registro bem-sucedido de um usuário.
type UserRegistered struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewUserRegistered(userID, email, firstName, lastName string) UserRegistered {
	return UserRegistered{
		EventID:    newEventID(),
		UserID:     userID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		OccurredAt: time.Now().UTC(),
	}
}

// ── Eventos de Order ─────────────────────────────────────────────────────────

// PaymentReceived é publicado pelo módulo order quando um pedido é criado e aguarda pagamento.
// É o evento que dispara a saga de pagamento.
type PaymentReceived struct {
	EventID     string    `json:"event_id"`
	SagaID      string    `json:"saga_id"`
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	AmountCents uint64    `json:"amount_cents"`
	Description string    `json:"description,omitempty"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func NewPaymentReceived(sagaID, orderID, userID string, amountCents uint64, description string) PaymentReceived {
	return PaymentReceived{
		EventID:     newEventID(),
		SagaID:      sagaID,
		OrderID:     orderID,
		UserID:      userID,
		AmountCents: amountCents,
		Description: description,
		OccurredAt:  time.Now().UTC(),
	}
}

// ── Eventos de Saga → Payment ────────────────────────────────────────────────

// BalanceReserved é publicado pela saga após reservar o saldo do usuário.
// O módulo payment deve consumi-lo para processar o pagamento.
type BalanceReserved struct {
	EventID        string    `json:"event_id"`
	SagaID         string    `json:"saga_id"`
	OrderID        string    `json:"order_id"`
	UserID         string    `json:"user_id"`
	AmountCents    uint64    `json:"amount_cents"`
	IdempotencyKey string    `json:"idempotency_key"`
	OccurredAt     time.Time `json:"occurred_at"`
}

func NewBalanceReserved(sagaID, orderID, userID, idempotencyKey string, amountCents uint64) BalanceReserved {
	return BalanceReserved{
		EventID:        newEventID(),
		SagaID:         sagaID,
		OrderID:        orderID,
		UserID:         userID,
		AmountCents:    amountCents,
		IdempotencyKey: idempotencyKey,
		OccurredAt:     time.Now().UTC(),
	}
}

// BalanceFailed é publicado pela saga quando a reserva de saldo falha (compensação).
type BalanceFailed struct {
	EventID    string    `json:"event_id"`
	SagaID     string    `json:"saga_id"`
	OrderID    string    `json:"order_id"`
	UserID     string    `json:"user_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewBalanceFailed(sagaID, orderID, userID, reason string) BalanceFailed {
	return BalanceFailed{
		EventID:    newEventID(),
		SagaID:     sagaID,
		OrderID:    orderID,
		UserID:     userID,
		Reason:     reason,
		OccurredAt: time.Now().UTC(),
	}
}

// ── Eventos de Payment → Saga ────────────────────────────────────────────────

// PaymentProcessed é publicado pelo módulo payment após processar o pagamento com sucesso.
type PaymentProcessed struct {
	EventID     string    `json:"event_id"`
	SagaID      string    `json:"saga_id"`
	OrderID     string    `json:"order_id"`
	ProviderRef string    `json:"provider_ref"`
	AmountCents uint64    `json:"amount_cents"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func NewPaymentProcessed(sagaID, orderID, providerRef string, amountCents uint64) PaymentProcessed {
	return PaymentProcessed{
		EventID:     newEventID(),
		SagaID:      sagaID,
		OrderID:     orderID,
		ProviderRef: providerRef,
		AmountCents: amountCents,
		OccurredAt:  time.Now().UTC(),
	}
}

// PaymentFailed é publicado pelo módulo payment quando o pagamento falha.
type PaymentFailed struct {
	EventID    string    `json:"event_id"`
	SagaID     string    `json:"saga_id"`
	OrderID    string    `json:"order_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewPaymentFailed(sagaID, orderID, reason string) PaymentFailed {
	return PaymentFailed{
		EventID:    newEventID(),
		SagaID:     sagaID,
		OrderID:    orderID,
		Reason:     reason,
		OccurredAt: time.Now().UTC(),
	}
}

// ── Eventos de Saga → Order ──────────────────────────────────────────────────

// OrderConfirmed é publicado pela saga após o pagamento ser processado com sucesso.
type OrderConfirmed struct {
	EventID     string    `json:"event_id"`
	SagaID      string    `json:"saga_id"`
	OrderID     string    `json:"order_id"`
	ProviderRef string    `json:"provider_ref"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func NewOrderConfirmed(sagaID, orderID, providerRef string) OrderConfirmed {
	return OrderConfirmed{
		EventID:     newEventID(),
		SagaID:      sagaID,
		OrderID:     orderID,
		ProviderRef: providerRef,
		OccurredAt:  time.Now().UTC(),
	}
}

// OrderCancelled é publicado pela saga quando a saga é compensada (rollback).
type OrderCancelled struct {
	EventID    string    `json:"event_id"`
	SagaID     string    `json:"saga_id"`
	OrderID    string    `json:"order_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewOrderCancelled(sagaID, orderID, reason string) OrderCancelled {
	return OrderCancelled{
		EventID:    newEventID(),
		SagaID:     sagaID,
		OrderID:    orderID,
		Reason:     reason,
		OccurredAt: time.Now().UTC(),
	}
}
