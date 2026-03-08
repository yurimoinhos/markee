package payments

import (
	"context"

	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/platform/httpauth"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	auth    gin.HandlerFunc
}

func NewModule(client *ent.Client, validateJWT func(context.Context, string) (*httpauth.AuthContext, error)) *Module {
	repo := NewRepository(client)
	svc := NewService(repo)
	return &Module{handler: NewHandler(svc), auth: httpauth.BearerAuthMiddleware(validateJWT)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/payments/confirm", m.handler.Confirm, m.auth, httpauth.RequireAnyPermission("payments.write")),
		router.GET("/payments", m.handler.List, m.auth, httpauth.RequireAnyPermission("payments.read")),
		router.POST("/payments/:id/evidence", m.handler.AddEvidence, m.auth, httpauth.RequireAnyPermission("payments.write")),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
