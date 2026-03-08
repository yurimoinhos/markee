package customers

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
		router.POST("/customers", m.handler.Create, m.auth, httpauth.RequireAnyPermission("customers.write")),
		router.PUT("/customers/:id", m.handler.Update, m.auth, httpauth.RequireAnyPermission("customers.write")),
		router.GET("/customers", m.handler.List, m.auth, httpauth.RequireAnyPermission("customers.read")),
		router.GET("/customers/:id", m.handler.Get, m.auth, httpauth.RequireAnyPermission("customers.read")),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
