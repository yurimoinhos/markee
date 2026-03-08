package contracts

// Tópicos RabbitMQ usados como routing keys no exchange "aggipay" (type: topic).
// Todos os módulos devem usar estas constantes para publicar e consumir eventos.
const (
	// Auth
	TopicUserRegistered = "user.registered"

	// Order → Saga
	TopicPaymentReceived = "payment.received"

	// Saga → Payment
	TopicBalanceReserved = "balance.reserved"
	TopicBalanceFailed   = "balance.failed"

	// Payment → Saga
	TopicPaymentProcessed = "payment.processed"
	TopicPaymentFailed    = "payment.failed"

	// Saga → Order
	TopicOrderConfirmed = "order.confirmed"
	TopicOrderCancelled = "order.cancelled"

	// Exchange e DLQ
	Exchange    = "aggipay"
	DLQExchange = "aggipay.dlq"
)
