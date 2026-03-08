package auth

import (
	"context"
	"log"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/modules/auth/application"
	"github.com/aggi-tech/aggipay/modules/auth/infrastructure"
	httpapi "github.com/aggi-tech/aggipay/modules/auth/interfaces/http"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler     *httpapi.Handler
	jwtProvider *infrastructure.JWTProvider
}

func NewModule(client *ent.Client, cfg *config.Config, publisher contracts.EventPublisher) *Module {
	repo := infrastructure.NewEntRepository(client)
	jwtProvider := infrastructure.NewJWTProvider(cfg.JWT)

	var oidcProvider application.OIDCProvider
	if cfg.OIDC.Google.ClientID != "" {
		p, err := infrastructure.NewGoogleOIDCProvider(context.Background(), cfg.OIDC.Google, cfg.JWT.Secret)
		if err != nil {
			log.Printf("[WARN] OIDC Google provider não inicializado: %v", err)
		} else {
			oidcProvider = p
		}
	}

	useCase := application.NewUseCase(repo, jwtProvider, oidcProvider, publisher)
	handler := httpapi.NewHandler(useCase)

	return &Module{handler: handler, jwtProvider: jwtProvider}
}

// Routes retorna as rotas do módulo de autenticação.
func (m *Module) Routes() router.Routes {
	auth := httpapi.BearerAuthMiddleware(func(tokenString string) (string, string, error) {
		claims, err := m.jwtProvider.ValidateToken(tokenString)
		if err != nil {
			return "", "", err
		}
		return claims.UserID, claims.Email, nil
	})

	return router.Routes{
		router.POST("/auth/register", m.handler.Register),
		router.POST("/auth/login", m.handler.Login),
		router.GET("/auth/google/login", m.handler.GoogleLogin),
		router.POST("/auth/google/callback", m.handler.GoogleCallback),
		router.GET("/auth/users", m.handler.ListUsers, auth),
		router.GET("/auth/users/:id", m.handler.GetUser, auth),
		router.PUT("/auth/users/:id", m.handler.UpdateUser, auth),
		router.DELETE("/auth/users/:id", m.handler.DeactivateUser, auth),
	}
}

// RegisterRoutes registra as rotas no RouterGroup informado.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}

// ValidateToken valida um Bearer token JWT e retorna as claims do usuário.
// Exportado para que outros módulos possam reutilizar o provider de JWT.
func (m *Module) ValidateToken(tokenString string) (*infrastructure.Claims, error) {
	return m.jwtProvider.ValidateToken(tokenString)
}
