package billing

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

// Create godoc
//
//	@Summary		Criar cobrança
//	@Description	Gera cobrança única, mensal ou por milestone
//	@Tags			Billing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateChargeInput	true	"Dados da cobrança"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/charges [post]
func (h *Handler) Create(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req CreateChargeInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// List godoc
//
//	@Summary		Listar cobranças
//	@Tags			Billing
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/charges [get]
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

// PaymentLink godoc
//
//	@Summary		Obter link de pagamento
//	@Tags			Billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID da cobrança"
//	@Success		200	{object}	PaymentLinkResponse
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/charges/{id}/pay-link [post]
func (h *Handler) PaymentLink(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	resp, err := h.svc.PaymentLink(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// PaymentQR godoc
//
//	@Summary		Obter QR Code de pagamento
//	@Tags			Billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID da cobrança"
//	@Success		200	{object}	PaymentQRResponse
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/charges/{id}/pay-qr [post]
func (h *Handler) PaymentQR(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	resp, err := h.svc.PaymentQR(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
