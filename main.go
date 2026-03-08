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
	authmodule "github.com/aggi-tech/aggipay/modules/auth"
	automationmodule "github.com/aggi-tech/aggipay/modules/automation"
	billingmodule "github.com/aggi-tech/aggipay/modules/billing"
	contractsmodule "github.com/aggi-tech/aggipay/modules/contracts"
	customersmodule "github.com/aggi-tech/aggipay/modules/customers"
	financemodule "github.com/aggi-tech/aggipay/modules/finance"
	paymentsmodule "github.com/aggi-tech/aggipay/modules/payments"
	projectsmodule "github.com/aggi-tech/aggipay/modules/projects"
	webhooksmodule "github.com/aggi-tech/aggipay/modules/webhooks"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/password"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env nao carregado, seguindo com variaveis de ambiente do sistema")
	}

	cfg := config.Load()
	password.GlobalPepper = cfg.Pepper

	client, err := db.NewClient(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("erro ao inicializar banco: %v", err)
	}
	defer client.Close()

	router := gin.Default()
	router.Use(problem.ErrorHandler())

	swag.FullGoSearchPath()

	router.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	router.GET("/swagger/*all", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	auth := authmodule.NewModule(client, cfg, nil)
	auth.RegisterRoutes(api)

	validateJWT := func(token string) (string, string, error) {
		claims, err := auth.ValidateToken(token)
		if err != nil {
			return "", "", err
		}
		return claims.UserID, claims.Email, nil
	}

	customersmodule.NewModule(client, validateJWT).RegisterRoutes(api)
	contractsmodule.NewModule(client, cfg, validateJWT).RegisterRoutes(api)
	billingmodule.NewModule(client, cfg, validateJWT).RegisterRoutes(api)
	paymentsmodule.NewModule(client, validateJWT).RegisterRoutes(api)
	projectsmodule.NewModule(client, validateJWT).RegisterRoutes(api)
	financemodule.NewModule(client, validateJWT).RegisterRoutes(api)
	automationmodule.NewModule(client, validateJWT).RegisterRoutes(api)
	webhooksmodule.NewModule(client, cfg).RegisterRoutes(api)

	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
