package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aggi-tech/aggipay/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TopicConsumer implementa contracts.EventConsumer consumindo um exchange topic.
// Cada tópico subscrito recebe uma fila dedicada (service_name.topic).
type TopicConsumer struct {
	conn        *Connection
	exchange    string
	dlqExchange string
	serviceName string
	maxRetries  int
	handlers    map[string]contracts.HandlerFunc
	mu          sync.RWMutex
}

// NewTopicConsumer cria um consumer orientado a tópicos.
// serviceName é usado para nomear as filas (ex: "order-service").
func NewTopicConsumer(conn *Connection, exchange, dlqExchange, serviceName string, maxRetries int) *TopicConsumer {
	return &TopicConsumer{
		conn:        conn,
		exchange:    exchange,
		dlqExchange: dlqExchange,
		serviceName: serviceName,
		maxRetries:  maxRetries,
		handlers:    make(map[string]contracts.HandlerFunc),
	}
}

// Subscribe registra um handler para o tópico informado.
// Deve ser chamado antes de Start.
func (c *TopicConsumer) Subscribe(topic string, handler contracts.HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[topic] = handler
}

// Start declara as filas para cada tópico subscrito e inicia os consumers.
// Bloqueia até ctx ser cancelado.
func (c *TopicConsumer) Start(ctx context.Context) error {
	c.mu.RLock()
	topics := make(map[string]contracts.HandlerFunc, len(c.handlers))
	for t, h := range c.handlers {
		topics[t] = h
	}
	c.mu.RUnlock()

	if len(topics) == 0 {
		return fmt.Errorf("rabbitmq consumer: nenhum tópico subscrito para %q", c.serviceName)
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq consumer: %w", err)
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("rabbitmq consumer: falha ao configurar QoS: %w", err)
	}

	var wg sync.WaitGroup
	for topic, handler := range topics {
		queueName := fmt.Sprintf("%s.%s", c.serviceName, topic)
		if err := DeclareQueue(ch, c.exchange, c.dlqExchange, queueName, topic); err != nil {
			return err
		}

		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("rabbitmq consumer: falha ao consumir %q: %w", queueName, err)
		}

		log.Printf("rabbitmq consumer [%s]: inscrito em tópico %q (fila: %s)", c.serviceName, topic, queueName)

		wg.Add(1)
		go func(topic string, handler contracts.HandlerFunc, msgs <-chan amqp.Delivery) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-msgs:
					if !ok {
						return
					}
					c.handle(ctx, msg, handler)
				}
			}
		}(topic, handler, msgs)
	}

	wg.Wait()
	return nil
}

func (c *TopicConsumer) handle(ctx context.Context, msg amqp.Delivery, handler contracts.HandlerFunc) {
	retryCount := 0
	if v, ok := msg.Headers["x-retry-count"]; ok {
		if n, ok := v.(int32); ok {
			retryCount = int(n)
		}
	}

	if err := handler(ctx, msg.Body); err != nil {
		log.Printf("rabbitmq consumer [%s]: erro (tentativa %d/%d) tópico %q: %v",
			c.serviceName, retryCount+1, c.maxRetries, msg.RoutingKey, err)

		if retryCount >= c.maxRetries {
			log.Printf("rabbitmq consumer [%s]: enviando para DLQ (max retries atingido)", c.serviceName)
			_ = msg.Nack(false, false)
			return
		}

		time.Sleep(time.Duration(retryCount+1) * 500 * time.Millisecond)
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}

// DecodeJSON deserializa o body da mensagem para dst.
func DecodeJSON(body []byte, dst any) error {
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("rabbitmq: falha ao deserializar mensagem: %w", err)
	}
	return nil
}

