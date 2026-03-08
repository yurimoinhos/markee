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
	authmodule "github.com/aggi-tech/aggipay/modules/auth"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/eventbus/rabbitmq"
	"github.com/aggi-tech/aggipay/platform/password"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[auth] .env não carregado, usando variáveis do sistema")
	}

	cfg := config.Load()
	password.GlobalPepper = cfg.Pepper

	addr := envOrDefault("AUTH_SERVER_ADDR", ":8001")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ── Banco de dados
	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[auth] falha ao conectar ao banco: %v", err)
	}
	defer client.Close()

	// ── RabbitMQ
	amqpConn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("[auth] falha ao conectar ao RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("[auth] falha ao abrir canal: %v", err)
	}
	if err := rabbitmq.SetupTopology(ch, contracts.Exchange, contracts.DLQExchange, "aggipay.dlq"); err != nil {
		log.Fatalf("[auth] falha ao configurar topologia: %v", err)
	}
	ch.Close()

	publisher := rabbitmq.NewPublisher(amqpConn, contracts.Exchange)

	// ── Módulo auth (passa publisher para emitir UserRegistered)
	auth := authmodule.NewModule(client, cfg, publisher)

	// ── HTTP
	r := gin.Default()
	r.Use(problem.ErrorHandler())

	api := r.Group("/api/v1")
	router.Register(api, auth.Routes())

	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		log.Printf("[auth] servidor HTTP em %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[auth] falha no servidor HTTP: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[auth] encerrando...")
	_ = srv.Shutdown(context.Background())
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
