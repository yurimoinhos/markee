package automation

import (
	"net/http"

	"github.com/aggi-tech/aggipay/platform/httpauth"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Run godoc
//
//	@Summary		Executar automações financeiras
//	@Tags			Automation
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	RunResult
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/automation/run [post]
func (h *Handler) Run(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	result, err := h.svc.Run(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListRuns godoc
//
//	@Summary		Histórico de automações
//	@Tags			Automation
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/automation/runs [get]
func (h *Handler) ListRuns(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	rows, err := h.svc.ListRuns(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rows)
}
