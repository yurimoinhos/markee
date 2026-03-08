package auth

import (
	"context"
	"errors"
	"log"
	"slices"
	"strings"

	"github.com/aggi-tech/aggipay/contracts"
	"github.com/aggi-tech/aggipay/ent"
	"github.com/aggi-tech/aggipay/modules/auth/application"
	"github.com/aggi-tech/aggipay/modules/auth/domain"
	"github.com/aggi-tech/aggipay/modules/auth/infrastructure"
	httpapi "github.com/aggi-tech/aggipay/modules/auth/interfaces/http"
	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/aggi-tech/aggipay/platform/httpauth"
	"github.com/aggi-tech/aggipay/platform/router"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler     *httpapi.Handler
	jwtProvider *infrastructure.JWTProvider
	repo        domain.Repository
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

	return &Module{handler: handler, jwtProvider: jwtProvider, repo: repo}
}

// Routes retorna as rotas do módulo de autenticação.
func (m *Module) Routes() router.Routes {
	auth := httpauth.BearerAuthMiddleware(m.ValidateToken)
	canReadUsers := httpauth.RequireAnyPermission(domain.PermUsersRead)
	canWriteUsers := httpauth.RequireAnyPermission(domain.PermUsersWrite)

	return router.Routes{
		router.POST("/auth/register", m.handler.Register),
		router.POST("/auth/login", m.handler.Login),
		router.GET("/auth/google/login", m.handler.GoogleLogin),
		router.POST("/auth/google/callback", m.handler.GoogleCallback),
		router.GET("/auth/me", m.handler.Me, auth),
		router.GET("/auth/users", m.handler.ListUsers, auth, canReadUsers),
		router.GET("/auth/users/:id", m.handler.GetUser, auth, canReadUsers),
		router.PUT("/auth/users/:id", m.handler.UpdateUser, auth, canWriteUsers),
		router.DELETE("/auth/users/:id", m.handler.DeactivateUser, auth, canWriteUsers),
	}
}

// RegisterRoutes registra as rotas no RouterGroup informado.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	router.Register(api, m.Routes())
}

// ValidateToken valida assinatura/expiração do JWT e confirma os dados no banco.
func (m *Module) ValidateToken(ctx context.Context, tokenString string) (*httpauth.AuthContext, error) {
	claims, err := m.jwtProvider.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := m.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if !user.Active {
		return nil, errors.New("usuário inativo")
	}

	expectedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	claimsEmail := strings.ToLower(strings.TrimSpace(claims.Email))
	if claimsEmail != expectedEmail {
		return nil, errors.New("email do token não confere com o banco")
	}

	expectedRole := domain.NormalizeRole(user.Role)
	claimsRole := strings.ToLower(strings.TrimSpace(claims.Role))
	if claimsRole != expectedRole {
		return nil, errors.New("role do token não confere com o banco")
	}

	expectedRoles := domain.NormalizeRoles(user.Roles)
	claimsRoles := domain.NormalizeRoles(claims.Roles)
	if !slices.Equal(claimsRoles, expectedRoles) {
		return nil, errors.New("roles do token não conferem com o banco")
	}

	expectedPermissions := domain.NormalizePermissions(domain.EnsurePermissionsForRoles(expectedRoles, user.Permissions))
	claimsPermissions := domain.NormalizePermissions(claims.Permissions)
	if !slices.Equal(claimsPermissions, expectedPermissions) {
		return nil, errors.New("permissions do token não conferem com o banco")
	}

	return &httpauth.AuthContext{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        expectedRole,
		Roles:       expectedRoles,
		Permissions: expectedPermissions,
	}, nil
}
