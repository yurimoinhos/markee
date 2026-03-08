package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/db"
	"github.com/aggi-tech/aggipay/modules/payment/application"
	"github.com/aggi-tech/aggipay/modules/payment/infrastructure"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/eventbus/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[payment] .env não carregado, usando variáveis do sistema")
	}

	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ── Banco de dados
	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[payment] falha ao conectar ao banco: %v", err)
	}
	defer client.Close()

	// ── RabbitMQ
	amqpConn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("[payment] falha ao conectar ao RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("[payment] falha ao abrir canal: %v", err)
	}
	if err := rabbitmq.SetupTopology(ch, contracts.Exchange, contracts.DLQExchange, "aggipay.dlq"); err != nil {
		log.Fatalf("[payment] falha ao configurar topologia: %v", err)
	}
	ch.Close()

	publisher := rabbitmq.NewPublisher(amqpConn, contracts.Exchange)
	consumer := rabbitmq.NewTopicConsumer(amqpConn, contracts.Exchange, contracts.DLQExchange, "payment-service", 3)

	// ── Wiring
	repo := infrastructure.NewEntRepository(client)
	svc := application.NewUseCase(repo, publisher)

	// Consome BalanceReserved → processa pagamento
	consumer.Subscribe(contracts.TopicBalanceReserved, func(ctx context.Context, raw []byte) error {
		var evt contracts.BalanceReserved
		if err := json.Unmarshal(raw, &evt); err != nil {
			return err
		}
		log.Printf("[payment] BalanceReserved recebido: saga=%s order=%s amount=%d", evt.SagaID, evt.OrderID, evt.AmountCents)
		return svc.ProcessPayment(ctx, evt)
	})

	log.Println("[payment] worker iniciado, aguardando eventos...")
	if err := consumer.Start(ctx); err != nil && ctx.Err() == nil {
		log.Fatalf("[payment] consumer parou com erro: %v", err)
	}

	log.Println("[payment] encerrado.")
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
