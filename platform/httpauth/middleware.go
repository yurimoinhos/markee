package httpauth

import (
	"context"
	"slices"
	"strings"

	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
)

type AuthContext struct {
	UserID      string
	Email       string
	Role        string
	Roles       []string
	Permissions []string
}

// BearerAuthMiddleware validates the bearer token and stores auth context.
func BearerAuthMiddleware(validate func(ctx context.Context, token string) (*AuthContext, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		if authorization == "" {
			_ = c.Error(problem.Unauthorized("token não informado"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			_ = c.Error(problem.Unauthorized("formato de token inválido, use: Bearer {token}"))
			c.Abort()
			return
		}

		authCtx, err := validate(c.Request.Context(), strings.TrimSpace(parts[1]))
		if err != nil {
			_ = c.Error(problem.Unauthorized("token inválido, expirado ou desatualizado"))
			c.Abort()
			return
		}

		c.Set("auth.user_id", authCtx.UserID)
		c.Set("auth.email", authCtx.Email)
		c.Set("auth.role", authCtx.Role)
		c.Set("auth.roles", authCtx.Roles)
		c.Set("auth.permissions", authCtx.Permissions)
		c.Next()
	}
}

func UserID(c *gin.Context) (string, bool) {
	v, ok := c.Get("auth.user_id")
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

func Role(c *gin.Context) (string, bool) {
	v, ok := c.Get("auth.role")
	if !ok {
		return "", false
	}
	role, ok := v.(string)
	return role, ok
}

func Roles(c *gin.Context) ([]string, bool) {
	v, ok := c.Get("auth.roles")
	if !ok {
		return nil, false
	}
	roles, ok := v.([]string)
	return roles, ok
}

func Permissions(c *gin.Context) ([]string, bool) {
	v, ok := c.Get("auth.permissions")
	if !ok {
		return nil, false
	}
	perms, ok := v.([]string)
	return perms, ok
}

func HasPermission(c *gin.Context, permission string) bool {
	perms, ok := Permissions(c)
	if !ok {
		return false
	}
	return slices.Contains(perms, strings.ToLower(strings.TrimSpace(permission)))
}

func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	required := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		normalized := strings.ToLower(strings.TrimSpace(permission))
		if normalized != "" {
			required = append(required, normalized)
		}
	}

	return func(c *gin.Context) {
		if len(required) == 0 {
			c.Next()
			return
		}

		for _, permission := range required {
			if HasPermission(c, permission) {
				c.Next()
				return
			}
		}

		_ = c.Error(problem.Forbidden("você não possui permissão para executar esta ação"))
		c.Abort()
	}
}
