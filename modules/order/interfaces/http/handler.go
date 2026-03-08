package http

import (
	"net/http"

	"github.com/aggi-tech/aggipay/modules/order/application"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
)

// Handler expõe os endpoints HTTP do módulo order.
type Handler struct {
	svc application.Service
}

func NewHandler(svc application.Service) *Handler {
	return &Handler{svc: svc}
}

type createOrderRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	AmountCents uint64 `json:"amount_cents" binding:"min=1"`
	Description string `json:"description"`
}

// CreateOrder godoc
// @Summary  Cria um pedido de pagamento
// @Tags     Orders
// @Accept   json
// @Produce  json
// @Param    body  body      createOrderRequest  true  "Dados do pedido"
// @Success  202   {object}  map[string]any
// @Router   /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest(err.Error()))
		return
	}

	order, err := h.svc.CreateOrder(c.Request.Context(), application.CreateOrderInput{
		UserID:      req.UserID,
		AmountCents: req.AmountCents,
		Description: req.Description,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id":           order.ID,
		"saga_id":      order.SagaID,
		"status":       order.Status,
		"amount_cents": order.AmountCents,
	})
}

// GetOrder godoc
// @Summary  Consulta um pedido pelo ID
// @Tags     Orders
// @Produce  json
// @Param    id   path      string  true  "ID do pedido"
// @Success  200  {object}  map[string]any
// @Router   /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, order)
}
