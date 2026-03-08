package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher implementa contracts.EventPublisher usando RabbitMQ com publisher confirms.
type Publisher struct {
	conn     *Connection
	exchange string
}

// NewPublisher cria um Publisher para o exchange informado.
func NewPublisher(conn *Connection, exchange string) *Publisher {
	return &Publisher{conn: conn, exchange: exchange}
}

// Publish serializa event como JSON e publica no tópico (routing key) informado.
// Usa publisher confirms para garantir entrega ao broker antes de retornar.
func (p *Publisher) Publish(ctx context.Context, topic string, event any) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq publisher: %w", err)
	}
	defer ch.Close()

	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("rabbitmq publisher: falha ao habilitar confirms: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("rabbitmq publisher: falha ao serializar evento %T: %w", event, err)
	}

	if err := ch.PublishWithContext(ctx, p.exchange, topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Body:         body,
	}); err != nil {
		return fmt.Errorf("rabbitmq publisher: falha ao publicar tópico %q: %w", topic, err)
	}

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return fmt.Errorf("rabbitmq publisher: broker nack para tópico %q", topic)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("rabbitmq publisher: timeout aguardando confirm: %w", ctx.Err())
	}
}

