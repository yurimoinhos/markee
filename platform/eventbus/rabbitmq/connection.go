// Package rabbitmq fornece conexão, publisher e consumer AMQP para RabbitMQ.
package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection gerencia a conexão AMQP com reconexão automática.
type Connection struct {
	url  string
	conn *amqp.Connection
}

// NewConnection cria e valida uma conexão AMQP. Retorna erro se o broker não estiver acessível.
func NewConnection(url string) (*Connection, error) {
	c := &Connection{url: url}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Connection) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("rabbitmq: falha ao conectar em %s: %w", c.url, err)
	}
	c.conn = conn
	return nil
}

// Channel abre um novo canal AMQP, reconectando se necessário.
func (c *Connection) Channel() (*amqp.Channel, error) {
	if c.conn == nil || c.conn.IsClosed() {
		log.Println("rabbitmq: reconectando...")
		for attempt := range 5 {
			if err := c.connect(); err == nil {
				break
			}
			wait := time.Duration(attempt+1) * time.Second
			log.Printf("rabbitmq: tentativa %d falhou, aguardando %s...", attempt+1, wait)
			time.Sleep(wait)
		}
		if c.conn == nil || c.conn.IsClosed() {
			return nil, fmt.Errorf("rabbitmq: não foi possível reconectar após 5 tentativas")
		}
	}
	return c.conn.Channel()
}

// Close encerra a conexão AMQP.
func (c *Connection) Close() error {
	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn.Close()
	}
	return nil
}

// SetupTopology declara o exchange topic, o DLQ exchange e a fila DLQ.
// Cada módulo é responsável por declarar suas próprias filas via DeclareQueue.
func SetupTopology(ch *amqp.Channel, exchange, dlqExchange, dlqQueue string) error {
	// Exchange principal tipo topic — suporta routing keys com wildcards
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: falha ao declarar exchange %q: %w", exchange, err)
	}

	// Exchange DLQ (direct) e fila de dead letters
	if err := ch.ExchangeDeclare(dlqExchange, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: falha ao declarar dlq exchange %q: %w", dlqExchange, err)
	}
	if _, err := ch.QueueDeclare(dlqQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: falha ao declarar dlq %q: %w", dlqQueue, err)
	}
	if err := ch.QueueBind(dlqQueue, "#", dlqExchange, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: falha ao bindar dlq: %w", err)
	}

	return nil
}

// DeclareQueue declara uma fila durável e a binda ao exchange com a routing key informada.
// dlqExchange é usado como dead-letter-exchange para mensagens não processadas.
func DeclareQueue(ch *amqp.Channel, exchange, dlqExchange, queueName, routingKey string) error {
	args := amqp.Table{
		"x-dead-letter-exchange": dlqExchange,
	}
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, args); err != nil {
		return fmt.Errorf("rabbitmq: falha ao declarar fila %q: %w", queueName, err)
	}
	if err := ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: falha ao bindar fila %q → %q: %w", queueName, routingKey, err)
	}
	return nil
}

