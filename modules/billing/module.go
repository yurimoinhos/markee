package billing

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
	svc := NewService(repo, NewAsaasProvider(cfg))
	return &Module{handler: NewHandler(svc), auth: httpauth.BearerAuthMiddleware(validateJWT)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/charges", m.handler.Create, m.auth, httpauth.RequireAnyPermission("charges.write")),
		router.GET("/charges", m.handler.List, m.auth, httpauth.RequireAnyPermission("charges.read")),
		router.POST("/charges/:id/pay-link", m.handler.PaymentLink, m.auth, httpauth.RequireAnyPermission("charges.write")),
		router.POST("/charges/:id/pay-qr", m.handler.PaymentQR, m.auth, httpauth.RequireAnyPermission("charges.write")),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
