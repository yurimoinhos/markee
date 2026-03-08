package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/db"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/eventbus/rabbitmq"
	"github.com/aggi-tech/aggipay/saga/orchestrator"
	"github.com/aggi-tech/aggipay/saga/state"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[saga] .env não carregado, usando variáveis do sistema")
	}

	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ── Banco de dados
	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[saga] falha ao conectar ao banco: %v", err)
	}
	defer client.Close()

	// ── RabbitMQ
	amqpConn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("[saga] falha ao conectar ao RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("[saga] falha ao abrir canal: %v", err)
	}
	if err := rabbitmq.SetupTopology(ch, contracts.Exchange, contracts.DLQExchange, "aggipay.dlq"); err != nil {
		log.Fatalf("[saga] falha ao configurar topologia: %v", err)
	}
	ch.Close()

	publisher := rabbitmq.NewPublisher(amqpConn, contracts.Exchange)
	consumer := rabbitmq.NewTopicConsumer(amqpConn, contracts.Exchange, contracts.DLQExchange, "saga-worker", 3)

	// ── Wiring
	sagaStore := state.NewStore(client)
	saga := orchestrator.NewPaymentSaga(sagaStore, publisher)
	saga.RegisterConsumers(consumer)

	log.Println("[saga] orchestrator iniciado, aguardando eventos...")
	if err := consumer.Start(ctx); err != nil && ctx.Err() == nil {
		log.Fatalf("[saga] consumer parou com erro: %v", err)
	}

	log.Println("[saga] encerrado.")
}
