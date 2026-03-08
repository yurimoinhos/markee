package httpauth

import (
	"strings"

	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
)

// BearerAuthMiddleware validates the bearer token and stores auth context.
func BearerAuthMiddleware(validate func(token string) (userID string, email string, err error)) gin.HandlerFunc {
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

		userID, email, err := validate(strings.TrimSpace(parts[1]))
		if err != nil {
			_ = c.Error(problem.Unauthorized("token inválido ou expirado"))
			c.Abort()
			return
		}

		c.Set("auth.user_id", userID)
		c.Set("auth.email", email)
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
