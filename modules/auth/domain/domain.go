package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/aggi-tech/aggipay/platform/problem"
)

// User representa a entidade de dominio do modulo de autenticacao.
type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	PhoneNumber   *string   `json:"phoneNumber,omitempty"`
	Balance       uint64    `json:"balance"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PasswordHash  *string   `json:"-"`
	OAuthProvider *string   `json:"-"`
	OAuthSub      *string   `json:"-"`
}

// Repository define contrato de persistencia do modulo de autenticacao.
type Repository interface {
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, id string, user User) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByOAuthSub(ctx context.Context, provider, sub string) (*User, error)
	FindOrCreateOAuthUser(ctx context.Context, user User) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]User, error)
	Count(ctx context.Context) (int, error)
	SoftDelete(ctx context.Context, id string) error
}

// Sentinelas de domínio do módulo auth.
// Cada erro tem um Type URI único — use errors.Is() para comparar.

var ErrUserNotFound = problem.New(
	http.StatusNotFound,
	"/problems/auth/user-not-found",
	"Usuário não encontrado",
	"O usuário solicitado não existe ou foi removido.",
)

var ErrEmailAlreadyInUse = problem.New(
	http.StatusConflict,
	"/problems/auth/email-already-in-use",
	"E-mail já cadastrado",
	"Já existe um usuário com este endereço de e-mail.",
)
