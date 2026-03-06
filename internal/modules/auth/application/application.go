package application

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aggi-tech/aggipay/internal/modules/auth/domain"
	"github.com/aggi-tech/aggipay/internal/modules/auth/infrastructure"
	"github.com/aggi-tech/aggipay/internal/platform/password"
	"github.com/aggi-tech/aggipay/internal/platform/problem"
	"github.com/go-playground/validator/v10"
)

// ── Erros ────────────────────────────────────────────────────────────────────

var ErrInvalidInput = problem.New(
	http.StatusBadRequest,
	"/problems/auth/invalid-input",
	"Dados de entrada inválidos",
	"",
)

var ErrInvalidCredentials = problem.New(
	http.StatusUnauthorized,
	"/problems/auth/invalid-credentials",
	"Credenciais inválidas",
	"O e-mail informado não possui conta ativa.",
)

// ── Interfaces ───────────────────────────────────────────────────────────────

// TokenProvider define o contrato para geração de JWTs.
type TokenProvider interface {
	GenerateToken(user domain.User) (string, error)
}

// OIDCProvider define o contrato do provider OIDC.
type OIDCProvider interface {
	AuthURL() (infrastructure.OIDCAuthParams, error)
	Exchange(ctx context.Context, code, stateJWT, codeVerifier string) (*infrastructure.OIDCUserInfo, error)
}

// Service define o contrato de casos de uso do módulo de autenticação.
type Service interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
	OIDCAuthURL() (infrastructure.OIDCAuthParams, error)
	OIDCExchange(ctx context.Context, input OIDCExchangeInput) (string, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*domain.User, error)
	DeactivateUser(ctx context.Context, id string) error
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type RegisterInput struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber *string
}

type LoginInput struct {
	Email    string
	Password string
}

type OIDCExchangeInput struct {
	Code         string
	State        string
	CodeVerifier string
}

type UpdateUserInput struct {
	FirstName   string
	LastName    string
	PhoneNumber *string
}

// ── UseCase ──────────────────────────────────────────────────────────────────

type UseCase struct {
	repo          domain.Repository
	tokenProvider TokenProvider
	oidcProvider  OIDCProvider
}

func NewUseCase(repo domain.Repository, tokenProvider TokenProvider, oidcProvider OIDCProvider) *UseCase {
	return &UseCase{repo: repo, tokenProvider: tokenProvider, oidcProvider: oidcProvider}
}

// Register cria um novo usuário com email e senha (login local).
func (u *UseCase) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	hash, err := password.Hash(input.Password)
	if err != nil {
		return nil, problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, err.Error())
	}

	entity := domain.User{
		FirstName:    strings.TrimSpace(input.FirstName),
		LastName:     strings.TrimSpace(input.LastName),
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		PhoneNumber:  cleanOptionalString(input.PhoneNumber),
		PasswordHash: &hash,
		Balance:      0,
		Active:       true,
	}

	return u.repo.Create(ctx, entity)
}

// Login autentica com email e senha e retorna um JWT.
func (u *UseCase) Login(ctx context.Context, input LoginInput) (string, error) {
	normalizedEmail := strings.TrimSpace(strings.ToLower(input.Email))
	if err := validator.New().Var(normalizedEmail, "required,email"); err != nil {
		return "", problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "email inválido")
	}

	found, err := u.repo.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !found.Active {
		return "", ErrInvalidCredentials
	}

	if found.PasswordHash == nil {
		return "", problem.New(ErrInvalidCredentials.StatusCode, ErrInvalidCredentials.Type,
			ErrInvalidCredentials.Title, "esta conta foi criada via provedor externo. Use o login social.")
	}

	if err := password.Verify(input.Password, *found.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	return u.tokenProvider.GenerateToken(*found)
}

// OIDCAuthURL retorna os parâmetros para iniciar o fluxo OIDC com PKCE.
func (u *UseCase) OIDCAuthURL() (infrastructure.OIDCAuthParams, error) {
	if u.oidcProvider == nil {
		return infrastructure.OIDCAuthParams{}, problem.Internal("provedor OIDC não configurado")
	}
	return u.oidcProvider.AuthURL()
}

// OIDCExchange troca o authorization code por um JWT da aplicação.
func (u *UseCase) OIDCExchange(ctx context.Context, input OIDCExchangeInput) (string, error) {
	if u.oidcProvider == nil {
		return "", problem.Internal("provedor OIDC não configurado")
	}

	info, err := u.oidcProvider.Exchange(ctx, input.Code, input.State, input.CodeVerifier)
	if err != nil {
		return "", problem.New(ErrInvalidCredentials.StatusCode, ErrInvalidCredentials.Type,
			ErrInvalidCredentials.Title, err.Error())
	}

	entity := domain.User{
		FirstName:     strings.TrimSpace(info.FirstName),
		LastName:      strings.TrimSpace(info.LastName),
		Email:         strings.TrimSpace(strings.ToLower(info.Email)),
		OAuthProvider: &info.Provider,
		OAuthSub:      &info.Sub,
		Active:        true,
	}

	found, err := u.repo.FindOrCreateOAuthUser(ctx, entity)
	if err != nil {
		return "", err
	}

	return u.tokenProvider.GenerateToken(*found)
}

func (u *UseCase) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "id é obrigatório")
	}
	if err := validateUpdateInput(input); err != nil {
		return nil, err
	}

	entity := domain.User{
		FirstName:   strings.TrimSpace(input.FirstName),
		LastName:    strings.TrimSpace(input.LastName),
		PhoneNumber: cleanOptionalString(input.PhoneNumber),
	}

	return u.repo.Update(ctx, id, entity)
}

func (u *UseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "id é obrigatório")
	}
	return u.repo.FindByID(ctx, id)
}

func (u *UseCase) ListUsers(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return u.repo.List(ctx, page, pageSize)
}

func (u *UseCase) DeactivateUser(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "id é obrigatório")
	}
	return u.repo.SoftDelete(ctx, id)
}

// IssueToken mantido para compatibilidade retroativa.
// Deprecated: use Login com email e senha.
func (u *UseCase) IssueToken(ctx context.Context, email string) (string, error) {
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))
	if err := validator.New().Var(normalizedEmail, "required,email"); err != nil {
		return "", problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "email inválido")
	}

	found, err := u.repo.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !found.Active {
		return "", ErrInvalidCredentials
	}

	return u.tokenProvider.GenerateToken(*found)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func validateRegisterInput(input RegisterInput) error {
	if err := validator.New().Var(strings.TrimSpace(input.Email), "required,email"); err != nil {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "email inválido")
	}
	if len(strings.TrimSpace(input.FirstName)) < 3 {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "first_name deve ter ao menos 3 caracteres")
	}
	if len(strings.TrimSpace(input.LastName)) < 3 {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "last_name deve ter ao menos 3 caracteres")
	}
	if strings.TrimSpace(input.Password) == "" {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "senha é obrigatória")
	}
	if input.PhoneNumber != nil && strings.TrimSpace(*input.PhoneNumber) != "" {
		if err := validator.New().Var(strings.TrimSpace(*input.PhoneNumber), "min=8,max=15"); err != nil {
			return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "phone_number deve ter entre 8 e 15 caracteres")
		}
	}
	return nil
}

func validateUpdateInput(input UpdateUserInput) error {
	if len(strings.TrimSpace(input.FirstName)) < 3 {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "first_name deve ter ao menos 3 caracteres")
	}
	if len(strings.TrimSpace(input.LastName)) < 3 {
		return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "last_name deve ter ao menos 3 caracteres")
	}
	if input.PhoneNumber != nil && strings.TrimSpace(*input.PhoneNumber) != "" {
		if err := validator.New().Var(strings.TrimSpace(*input.PhoneNumber), "min=8,max=15"); err != nil {
			return problem.New(ErrInvalidInput.StatusCode, ErrInvalidInput.Type, ErrInvalidInput.Title, "phone_number deve ter entre 8 e 15 caracteres")
		}
	}
	return nil
}

func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
