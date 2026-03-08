package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/aggi-tech/aggipay/modules/auth/application"
	"github.com/aggi-tech/aggipay/modules/auth/domain"
	"github.com/aggi-tech/aggipay/platform/problem"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase application.Service
}

// registerRequest representa o body de criação de usuário com senha.
type registerRequest struct {
	FirstName   string `form:"firstName"   example:"João"`
	LastName    string `form:"lastName"    example:"Silva"`
	Email       string `form:"email"       example:"joao.silva@example.com"`
	PhoneNumber string `form:"phoneNumber" example:"11999999999"`
	Password    string `form:"password"    example:"minha-senha-segura-123"`
}

// loginRequest representa o body de login local.
type loginRequest struct {
	Email    string `form:"email"    example:"joao.silva@example.com"`
	Password string `form:"password" example:"minha-senha-segura-123"`
}

// UpdateUserRequest representa o body de atualização de usuário.
type updateUserRequest struct {
	FirstName   string `json:"firstName"   example:"João"`
	LastName    string `json:"lastName"    example:"Silva"`
	PhoneNumber string `json:"phoneNumber" example:"11999999999"`
}

// tokenResponse representa a resposta com o access token.
type tokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGci..."`
	TokenType   string `json:"token_type"   example:"Bearer"`
}

// oidcAuthResponse representa a resposta do início do fluxo OIDC.
type oidcAuthResponse struct {
	AuthURL      string `json:"auth_url"      example:"https://accounts.google.com/o/oauth2/auth?..."`
	State        string `json:"state"         example:"eyJhbGci..."`
	CodeVerifier string `json:"code_verifier" example:"dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"`
}

// oidcCallbackRequest representa o body do callback OIDC.
type oidcCallbackRequest struct {
	Code         string `json:"code"          example:"4/0AX4XfWi..."`
	State        string `json:"state"         example:"eyJhbGci..."`
	CodeVerifier string `json:"code_verifier" example:"dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"`
}

// userPageResponse is a concrete type used for swagger documentation of paginated user responses.
// It mirrors common.Page[domain.User], which swaggo cannot parse due to Go generics.
type userPageResponse struct {
	Items      []domain.User `json:"items"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	Total      int           `json:"total"`
	TotalPages int           `json:"totalPages"`
}

func NewHandler(useCase application.Service) *Handler {
	return &Handler{useCase: useCase}
}

// Register godoc
//
//	@Summary		Registrar usuário
//	@Description	Cria um novo usuário com email e senha. Usa criptografia Argon2ID com alto uso de processamento.
//	@Tags			Auth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			firstName	formData	string	true	"Primeiro nome"
//	@Param			lastName	formData	string	true	"Sobrenome"
//	@Param			email		formData	string	true	"E-mail"
//	@Param			phoneNumber	formData	string	false	"Telefone"
//	@Param			password	formData	string	true	"Senha"
//	@Success		201			{object}	domain.User
//	@Failure		400			{object}	problem.BaseErr
//	@Failure		409			{object}	problem.BaseErr
//	@Failure		500			{object}	problem.BaseErr
//	@Router			/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	user, err := h.useCase.Register(c.Request.Context(), application.RegisterInput{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: optionalString(req.PhoneNumber),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
//
//	@Summary		Login local
//	@Description	Autentica com email e senha, retorna JWT de acesso
//	@Tags			Auth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			email		formData	string	true	"E-mail"
//	@Param			password	formData	string	true	"Senha"
//	@Success		200			{object}	tokenResponse
//	@Failure		400			{object}	problem.BaseErr
//	@Failure		401			{object}	problem.BaseErr
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	token, err := h.useCase.Login(c.Request.Context(), application.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{AccessToken: token, TokenType: "Bearer"})
}

// GoogleLogin godoc
//
//	@Summary		Iniciar login com Google
//	@Description	Retorna a URL de autorização do Google e os parâmetros PKCE. O cliente deve armazenar state e code_verifier, depois redirecionar o usuário para auth_url.
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	oidcAuthResponse
//	@Failure		500	{object}	problem.BaseErr
//	@Router			/auth/google/login [get]
func (h *Handler) GoogleLogin(c *gin.Context) {
	params, err := h.useCase.OIDCAuthURL()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, oidcAuthResponse{
		AuthURL:      params.AuthURL,
		State:        params.State,
		CodeVerifier: params.CodeVerifier,
	})
}

// GoogleCallback godoc
//
//	@Summary		Finalizar login com Google
//	@Description	Troca o authorization code pelo JWT da aplicação usando PKCE. O cliente envia o code recebido do Google, o state e o code_verifier gerados no passo anterior.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		oidcCallbackRequest	true	"Parâmetros do callback"
//	@Success		200		{object}	tokenResponse
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Router			/auth/google/callback [post]
func (h *Handler) GoogleCallback(c *gin.Context) {
	var req oidcCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	token, err := h.useCase.OIDCExchange(c.Request.Context(), application.OIDCExchangeInput{
		Code:         req.Code,
		State:        req.State,
		CodeVerifier: req.CodeVerifier,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{AccessToken: token, TokenType: "Bearer"})
}

// UpdateUser godoc
//
//	@Summary		Atualizar usuário
//	@Description	Atualiza os dados de um usuário pelo ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"ID do usuário"
//	@Param			body	body		updateUserRequest	true	"Dados a atualizar"
//	@Success		200		{object}	domain.User
//	@Failure		400		{object}	problem.BaseErr
//	@Failure		401		{object}	problem.BaseErr
//	@Failure		404		{object}	problem.BaseErr
//	@Router			/auth/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(problem.BadRequest("payload inválido"))
		return
	}

	user, err := h.useCase.UpdateUser(c.Request.Context(), c.Param("id"), application.UpdateUserInput{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: optionalString(req.PhoneNumber),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeactivateUser godoc
//
//	@Summary		Desativar usuário
//	@Description	Realiza soft-delete de um usuário pelo ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID do usuário"
//	@Success		200	{object}	map[string]string
//	@Failure		401	{object}	problem.BaseErr
//	@Failure		404	{object}	problem.BaseErr
//	@Router			/auth/users/{id} [delete]
func (h *Handler) DeactivateUser(c *gin.Context) {
	if err := h.useCase.DeactivateUser(c.Request.Context(), c.Param("id")); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário desativado com sucesso."})
}

// ListUsers godoc
//
//	@Summary		Listar usuários
//	@Description	Retorna lista paginada de usuários
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"Página (padrão: 1)"
//	@Param			page_size	query		int	false	"Tamanho da página (padrão: 10, máx: 100)"
//	@Success		200			{object}	userPageResponse
//	@Failure		401			{object}	problem.BaseErr
//	@Router			/auth/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	page := parseIntOrDefault(c.Query("page"), 1)
	pageSize := parseIntOrDefault(c.Query("page_size"), 10)

	result, err := h.useCase.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUser godoc
//
//	@Summary		Buscar usuário
//	@Description	Retorna um usuário pelo ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"ID do usuário"
//	@Success		200	{object}	domain.User
//	@Failure		401	{object}	problem.BaseErr
//	@Failure		404	{object}	problem.BaseErr
//	@Router			/auth/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	user, err := h.useCase.GetUser(c.Request.Context(), c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func parseIntOrDefault(value string, defaultValue int) int {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// BearerAuthMiddleware valida o token Bearer do header Authorization.
func BearerAuthMiddleware(validate func(tokenString string) (userID string, email string, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		if authorization == "" {
			_ = c.Error(problem.Unauthorized("token não informado"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			_ = c.Error(problem.Unauthorized("formato de token inválido, use: Bearer {token}"))
			c.Abort()
			return
		}

		userID, email, err := validate(strings.TrimSpace(parts[1]))
		if err != nil {
			_ = c.Error(problem.Unauthorized("token inválido ou expirado"))
			c.Abort()
			return
		}

		c.Set("auth.user_id", userID)
		c.Set("auth.email", email)
		c.Next()
	}
}
