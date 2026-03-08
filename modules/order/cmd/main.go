package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/db"
	"github.com/aggi-tech/aggipay/modules/order/application"
	"github.com/aggi-tech/aggipay/modules/order/infrastructure"
	orderhttp "github.com/aggi-tech/aggipay/modules/order/interfaces/http"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/eventbus/rabbitmq"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[order] .env não carregado, usando variáveis do sistema")
	}

	cfg := config.Load()
	addr := envOrDefault("ORDER_SERVER_ADDR", ":8002")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ── Banco de dados
	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[order] falha ao conectar ao banco: %v", err)
	}
	defer client.Close()

	// ── RabbitMQ
	amqpConn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("[order] falha ao conectar ao RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("[order] falha ao abrir canal RabbitMQ: %v", err)
	}
	if err := rabbitmq.SetupTopology(ch, contracts.Exchange, contracts.DLQExchange, "aggipay.dlq"); err != nil {
		log.Fatalf("[order] falha ao configurar topologia: %v", err)
	}
	ch.Close()

	publisher := rabbitmq.NewPublisher(amqpConn, contracts.Exchange)
	consumer := rabbitmq.NewTopicConsumer(amqpConn, contracts.Exchange, contracts.DLQExchange, "order-service", 3)

	// ── Wiring
	repo := infrastructure.NewEntRepository(client)
	svc := application.NewUseCase(repo, publisher)
	handler := orderhttp.NewHandler(svc)

	// Registra consumers de eventos
	orderhttp.RegisterConsumers(consumer, svc)

	// ── HTTP
	r := gin.Default()
	r.Use(problem.ErrorHandler())

	api := r.Group("/api/v1")
	routes := router.Routes{
		router.POST("/orders", handler.CreateOrder),
		router.GET("/orders/:id", handler.GetOrder),
	}
	router.Register(api, routes)

	// ── Consumer em background
	go func() {
		if err := consumer.Start(ctx); err != nil && ctx.Err() == nil {
			log.Printf("[order] consumer parou com erro: %v", err)
		}
	}()

	// ── Servidor HTTP
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		log.Printf("[order] servidor HTTP em %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[order] falha no servidor HTTP: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[order] encerrando...")
	_ = srv.Shutdown(context.Background())
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
