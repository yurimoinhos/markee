package customers

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
//	@Summary		Criar cliente
//	@Description	Cadastra um cliente para gestão financeira/contratos
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateCustomerInput	true	"Dados do cliente"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/customers [post]
func (h *Handler) Create(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}

	var req CreateCustomerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	entity, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

// Update godoc
//
//	@Summary		Atualizar cliente
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID do cliente"
//	@Param			body	body		UpdateCustomerInput	true	"Dados para atualização"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/customers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}

	var req UpdateCustomerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	entity, err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, entity)
}

// List godoc
//
//	@Summary		Listar clientes
//	@Tags			Customers
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/customers [get]
func (h *Handler) List(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}

	list, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// Get godoc
//
//	@Summary		Detalhar cliente
//	@Description	Retorna cadastro e resumo financeiro do cliente
//	@Tags			Customers
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID do cliente"
//	@Success		200	{object}	map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Failure		404	{object}	problem.BaseErr
//	@Router			/customers/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}

	entity, summary, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customer":          entity,
		"financial_summary": summary,
	})
}
