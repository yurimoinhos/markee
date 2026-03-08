package automation

import (
	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/platform/httpauth"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	auth    gin.HandlerFunc
}

func NewModule(client *ent.Client, validateJWT func(string) (string, string, error)) *Module {
	repo := NewRepository(client)
	svc := NewService(repo)
	return &Module{handler: NewHandler(svc), auth: httpauth.BearerAuthMiddleware(validateJWT)}
}

func (m *Module) Routes() router.Routes {
	return router.Routes{
		router.POST("/automation/run", m.handler.Run, m.auth),
		router.GET("/automation/runs", m.handler.ListRuns, m.auth),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}
