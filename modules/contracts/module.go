package contracts

import (
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

func NewModule(client *ent.Client, cfg *config.Config, validateJWT func(string) (string, string, error)) *Module {
	repo := NewRepository(client)
	svc := NewService(repo, NewClicksignProvider(cfg))
	return &Module{handler: NewHandler(svc), auth: httpauth.BearerAuthMiddleware(validateJWT)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/contracts", m.handler.Create, m.auth),
		router.GET("/contracts", m.handler.List, m.auth),
		router.GET("/contracts/:id", m.handler.Get, m.auth),
		router.POST("/contracts/:id/generate", m.handler.Generate, m.auth),
		router.POST("/contracts/:id/send-signature", m.handler.SendSignature, m.auth),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
