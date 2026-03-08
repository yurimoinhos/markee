// Package problem implementa o padrão RFC 9457 (Problem Details for HTTP APIs).
package problem

import (
	"net/http"
	"time"
)

// BaseErr representa um erro padronizado conforme RFC 9457.
// Todos os erros da API devem ser (ou embrulhar) um BaseErr.
type BaseErr struct {
	// Type é um URI que identifica o tipo do problema (único por categoria de erro).
	Type string `json:"type"`
	// Title é o título legível do tipo de problema.
	Title string `json:"error"`
	// ErrorDescription é uma explicação específica desta ocorrência.
	ErrorDescription string `json:"errorDescription,omitempty"`
	// StatusCode é o código HTTP associado ao problema.
	StatusCode int `json:"statusCode"`
	// Timestamp indica quando o erro ocorreu (UTC, ISO 8601).
	Timestamp time.Time `json:"timestamp"`
	// Instance é o URI da requisição que originou o erro.
	Instance string `json:"instance,omitempty"`
	// Details contém informações adicionais opcionais sobre o erro.
	Details any `json:"details,omitempty"`
}

func (e *BaseErr) Error() string {
	if e.ErrorDescription != "" {
		return e.ErrorDescription
	}
	return e.Title
}

// Is permite comparação via errors.Is usando o Type como discriminador.
func (e *BaseErr) Is(target error) bool {
	t, ok := target.(*BaseErr)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// New cria um BaseErr com timestamp preenchido automaticamente.
func New(status int, errType, title, description string) *BaseErr {
	return &BaseErr{
		Type:             errType,
		Title:            title,
		ErrorDescription: description,
		StatusCode:       status,
		Timestamp:        time.Now().UTC(),
	}
}

// WithDetails adiciona informações extras ao BaseErr (fluent).
func (e *BaseErr) WithDetails(details any) *BaseErr {
	e.Details = details
	return e
}

// Construtores padrão por categoria HTTP.

func BadRequest(description string) *BaseErr {
	return New(http.StatusBadRequest, "/problems/bad-request", "Requisição inválida", description)
}

func Unauthorized(description string) *BaseErr {
	return New(http.StatusUnauthorized, "/problems/unauthorized", "Não autorizado", description)
}

func Forbidden(description string) *BaseErr {
	return New(http.StatusForbidden, "/problems/forbidden", "Acesso proibido", description)
}

func NotFound(description string) *BaseErr {
	return New(http.StatusNotFound, "/problems/not-found", "Recurso não encontrado", description)
}

func Conflict(description string) *BaseErr {
	return New(http.StatusConflict, "/problems/conflict", "Conflito", description)
}

func UnprocessableEntity(description string) *BaseErr {
	return New(http.StatusUnprocessableEntity, "/problems/unprocessable-entity", "Entidade não processável", description)
}

func Internal(description string) *BaseErr {
	return New(http.StatusInternalServerError, "/problems/internal-server-error", "Erro interno do servidor", description)
}

func NotImplemented(description string) *BaseErr {
	return New(http.StatusNotImplemented, "/problems/not-implemented", "Funcionalidade não disponível", description)
}
