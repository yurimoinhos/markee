package contracts

import (
	"context"

	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/httpauth"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	auth    gin.HandlerFunc
}

func NewModule(client *ent.Client, cfg *config.Config, validateJWT func(context.Context, string) (*httpauth.AuthContext, error)) *Module {
	repo := NewRepository(client)
	svc := NewService(repo, NewClicksignProvider(cfg))
	return &Module{handler: NewHandler(svc), auth: httpauth.BearerAuthMiddleware(validateJWT)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/contracts", m.handler.Create, m.auth, httpauth.RequireAnyPermission("contracts.write")),
		router.GET("/contracts", m.handler.List, m.auth, httpauth.RequireAnyPermission("contracts.read")),
		router.GET("/contracts/:id", m.handler.Get, m.auth, httpauth.RequireAnyPermission("contracts.read")),
		router.POST("/contracts/:id/generate", m.handler.Generate, m.auth, httpauth.RequireAnyPermission("contracts.write")),
		router.POST("/contracts/:id/send-signature", m.handler.SendSignature, m.auth, httpauth.RequireAnyPermission("contracts.write")),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
