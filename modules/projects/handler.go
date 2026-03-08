package projects

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

// CreateProject godoc
//
//	@Summary		Criar projeto
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateProjectInput	true	"Dados do projeto"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/projects [post]
func (h *Handler) CreateProject(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req CreateProjectInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.CreateProject(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// ListProjects godoc
//
//	@Summary		Listar projetos
//	@Tags			Projects
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		map[string]any
//	@Failure		401	{object}	problem.BaseErr
//	@Router			/projects [get]
func (h *Handler) ListProjects(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	items, err := h.svc.ListProjects(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateMilestone godoc
//
//	@Summary		Criar milestone
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID do projeto"
//	@Param			body	body		CreateMilestoneInput	true	"Dados do milestone"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/projects/{id}/milestones [post]
func (h *Handler) CreateMilestone(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req CreateMilestoneInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.CreateMilestone(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// CreateWorklog godoc
//
//	@Summary		Criar worklog
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID do projeto"
//	@Param			body	body		CreateWorklogInput	true	"Dados do worklog"
//	@Success		201		{object}	map[string]any
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/projects/{id}/worklogs [post]
func (h *Handler) CreateWorklog(c *gin.Context) {
	userID, ok := httpauth.UserID(c)
	if !ok {
		_ = c.Error(problem.Unauthorized("usuário não autenticado"))
		return
	}
	var req CreateWorklogInput
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}
	item, err := h.svc.CreateWorklog(c.Request.Context(), userID, c.Param("id"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
