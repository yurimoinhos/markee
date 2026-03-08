package contracts

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
//	@Summary		Criar contrato
//	@Description	Cadastra contrato de serviço de software
//	@Tags			Contracts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateContractInput	true	"Dados do contrato"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/contracts [post]
func (h *Handler) Create(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req CreateContractInput
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

// List godoc
//
//	@Summary		Listar contratos
//	@Tags			Contracts
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/contracts [get]
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

// Get godoc
//
//	@Summary		Buscar contrato
//	@Tags			Contracts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID do contrato"
//	@Success		200	{object}	map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Failure		404	{object}	problem.BaseErr
//	@Router			/contracts/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	item, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// Generate godoc
//
//	@Summary		Gerar versão de contrato
//	@Description	Gera artefatos PDF/editável para assinatura
//	@Tags			Contracts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"ID do contrato"
//	@Param			body	body		GenerateContractInput	true	"Template e conteúdo"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/contracts/{id}/generate [post]
func (h *Handler) Generate(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req GenerateContractInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	version, err := h.svc.Generate(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, version)
}

// SendSignature godoc
//
//	@Summary		Enviar para assinatura
//	@Description	Cria documento no Clicksign e retorna URL de assinatura
//	@Tags			Contracts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID do contrato"
//	@Param			body	body		SendSignatureInput	true	"Dados do signatário"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/contracts/{id}/send-signature [post]
func (h *Handler) SendSignature(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req SendSignatureInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	signature, signURL, err := h.svc.SendSignature(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"signature": signature, "sign_url": signURL})
}
