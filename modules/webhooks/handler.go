package webhooks

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Asaas godoc
//
//	@Summary		Webhook Asaas
//	@Description	Recebe eventos de cobrança/pagamento do Asaas
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			body	body		map[string]any	true	"Evento Asaas"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		401		{object}	map[string]any
//	@Router			/webhooks/asaas [post]
func (h *Handler) Asaas(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	status, payload, svcErr := h.svc.ProcessAsaas(c.Request.Context(), c.GetHeader("asaas-access-token"), body)
	if svcErr != nil {
		_ = c.Error(svcErr)
		return
	}
	c.JSON(status, payload)
}

// Clicksign godoc
//
//	@Summary		Webhook Clicksign
//	@Description	Recebe eventos de assinatura digital do Clicksign
//	@Tags			Webhooks
//	@Accept			json
//	@Produce		json
//	@Param			body	body		map[string]any	true	"Evento Clicksign"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		401		{object}	map[string]any
//	@Router			/webhooks/clicksign [post]
func (h *Handler) Clicksign(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	status, payload, svcErr := h.svc.ProcessClicksign(c.Request.Context(), c.GetHeader("X-Hub-Signature"), body)
	if svcErr != nil {
		_ = c.Error(svcErr)
		return
	}
	c.JSON(status, payload)
}
