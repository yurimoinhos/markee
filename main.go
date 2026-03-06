// Package main é o entrypoint da aplicação AggiPay.
//
//	@title					AggiPay API
//	@version				1.0.0
//	@description			API modular de pagamentos com autenticação Bearer Token.
//	@host					localhost:8000
//	@BasePath				/api/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in						header
//	@name					Authorization
//	@description			Informe: Bearer {token}
//
//go:generate go run github.com/swaggo/swag/cmd/swag init --generalInfo main.go --dir . --output docs --parseDependency --parseInternal
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aggi-tech/aggipay/db"
	_ "github.com/aggi-tech/aggipay/docs"
	authmodule "github.com/aggi-tech/aggipay/internal/modules/auth"
	"github.com/aggi-tech/aggipay/internal/platform/config"
	"github.com/aggi-tech/aggipay/internal/platform/problem"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env nao carregado, seguindo com variaveis de ambiente do sistema")
	}

	cfg := config.Load()

	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("erro ao inicializar banco: %v", err)
	}
	defer client.Close()

	router := gin.Default()
	router.Use(problem.ErrorHandler())
	router.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	router.GET("/swagger/*all", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	auth := authmodule.NewModule(client, cfg)
	auth.RegisterRoutes(api)

	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
