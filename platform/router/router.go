// Package router fornece utilitários para registro declarativo de rotas Gin.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route descreve um endpoint da API.
type Route struct {
	Method      string            // HTTP method (GET, POST, PUT, DELETE…)
	Path        string            // Caminho relativo ao RouterGroup pai
	Handler     gin.HandlerFunc   // Handler principal
	Middlewares []gin.HandlerFunc // Middlewares aplicados apenas a esta rota
}

// Routes é uma lista de Route.
type Routes []Route

// Module é a interface que todo módulo de feature deve implementar
// para que suas rotas sejam registradas pelo bootstrapper central.
type Module interface {
	Routes() Routes
}

// Register registra todas as rotas de um slice no RouterGroup informado.
func Register(group *gin.RouterGroup, routes Routes) {
	for _, r := range routes {
		handlers := make([]gin.HandlerFunc, 0, len(r.Middlewares)+1)
		handlers = append(handlers, r.Middlewares...)
		handlers = append(handlers, r.Handler)
		group.Handle(r.Method, r.Path, handlers...)
	}
}

// GET é um helper para definir uma rota GET.
func GET(path string, handler gin.HandlerFunc, mws ...gin.HandlerFunc) Route {
	return Route{Method: http.MethodGet, Path: path, Handler: handler, Middlewares: mws}
}

// POST é um helper para definir uma rota POST.
func POST(path string, handler gin.HandlerFunc, mws ...gin.HandlerFunc) Route {
	return Route{Method: http.MethodPost, Path: path, Handler: handler, Middlewares: mws}
}

// PUT é um helper para definir uma rota PUT.
func PUT(path string, handler gin.HandlerFunc, mws ...gin.HandlerFunc) Route {
	return Route{Method: http.MethodPut, Path: path, Handler: handler, Middlewares: mws}
}

// DELETE é um helper para definir uma rota DELETE.
func DELETE(path string, handler gin.HandlerFunc, mws ...gin.HandlerFunc) Route {
	return Route{Method: http.MethodDelete, Path: path, Handler: handler, Middlewares: mws}
}
