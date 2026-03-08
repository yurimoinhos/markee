package finance

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

// Dashboard godoc
//
//	@Summary		Dashboard financeiro
//	@Tags			Finance
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	Dashboard
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/finance/dashboard [get]
func (h *Handler) Dashboard(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	data, err := h.svc.Dashboard(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// CashFlow godoc
//
//	@Summary		Fluxo de caixa
//	@Tags			Finance
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		CashFlowPoint
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/finance/cashflow [get]
func (h *Handler) CashFlow(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	data, err := h.svc.CashFlow(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// Defaults godoc
//
//	@Summary		Métricas de inadimplência e crescimento
//	@Tags			Finance
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	DefaultMetrics
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/finance/defaults [get]
func (h *Handler) Defaults(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	data, err := h.svc.Defaults(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}
