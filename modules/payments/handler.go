package payments

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

// Confirm godoc
//
//	@Summary		Confirmar pagamento
//	@Description	Confirma pagamento manual ou por conciliação externa
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		ConfirmPaymentInput	true	"Dados do pagamento"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/payments/confirm [post]
func (h *Handler) Confirm(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req ConfirmPaymentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.Confirm(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// List godoc
//
//	@Summary		Listar pagamentos
//	@Tags			Payments
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/payments [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	items, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// AddEvidence godoc
//
//	@Summary		Adicionar comprovante
//	@Description	Anexa evidência/comprovante e hash de transação cripto
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string			true	"ID do pagamento"
//	@Param			body	body		EvidenceInput	true	"Dados do comprovante"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/payments/{id}/evidence [post]
func (h *Handler) AddEvidence(c *gin.Context) {
	_, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req EvidenceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.AddEvidence(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
