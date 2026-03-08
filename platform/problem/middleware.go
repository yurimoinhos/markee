package problem

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler é um middleware Gin que intercepta todos os erros registrados via c.Error()
// e os converte em respostas RFC 9457. Deve ser registrado globalmente no router.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Resposta já foi escrita por algum handler — não sobrescrever.
		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var prob *BaseErr
		if !errors.As(err, &prob) {
			prob = Internal(err.Error())
		}

		// Preenche instance com o path da requisição se ainda não definido.
		if prob.Instance == "" {
			prob.Instance = c.Request.RequestURI
		}

		c.JSON(prob.StatusCode, prob)
	}
}
