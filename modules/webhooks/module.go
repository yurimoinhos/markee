package webhooks

import (
	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule(client *ent.Client, cfg *config.Config) *Module {
	repo := NewRepository(client)
	svc := NewService(repo, cfg)
	return &Module{handler: NewHandler(svc)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/webhooks/asaas", m.handler.Asaas),
		router.POST("/webhooks/clicksign", m.handler.Clicksign),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
