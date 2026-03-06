package swaggerui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OpenAPIHandler() gin.HandlerFunc {
	spec := `{
  "openapi": "3.0.3",
  "info": {
    "title": "AggiPay API",
    "version": "1.0.0",
    "description": "Documentacao da API modular com autenticacao Bearer Token"
  },
  "servers": [
    {
      "url": "http://localhost:8000"
    }
  ],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "RegisterRequest": {
        "type": "object",
        "required": ["firstName", "lastName", "email"],
        "properties": {
          "firstName": {"type": "string", "example": "Yuri"},
          "lastName": {"type": "string", "example": "Moinhos"},
          "email": {"type": "string", "format": "email", "example": "yuri@example.com"},
          "phoneNumber": {"type": "string", "example": "11999999999"}
        }
      },
      "UpdateUserRequest": {
        "type": "object",
        "required": ["firstName", "lastName"],
        "properties": {
          "firstName": {"type": "string", "example": "Yuri"},
          "lastName": {"type": "string", "example": "Moinhos"},
          "phoneNumber": {"type": "string", "example": "11999999999"}
        }
      },
      "TokenRequest": {
        "type": "object",
        "required": ["email"],
        "properties": {
          "email": {"type": "string", "format": "email", "example": "yuri@example.com"}
        }
      },
      "TokenResponse": {
        "type": "object",
        "properties": {
          "access_token": {"type": "string"},
          "token_type": {"type": "string", "example": "Bearer"}
        }
      },
      "User": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "first_name": {"type": "string"},
          "last_name": {"type": "string"},
          "email": {"type": "string", "format": "email"},
          "phone_number": {"type": "string"},
          "balance": {"type": "integer", "format": "uint64"},
          "active": {"type": "boolean"}
        }
      }
    }
  },
  "paths": {
    "/api/v1/auth/register": {
      "post": {
        "tags": ["Auth"],
        "summary": "Registra usuario",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/RegisterRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Usuario criado",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/User"}
              }
            }
          }
        }
      }
    },
    "/api/v1/auth/token": {
      "post": {
        "tags": ["Auth"],
        "summary": "Gera access token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/TokenRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Token gerado",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/TokenResponse"}
              }
            }
          },
          "401": {"description": "Credenciais invalidas"}
        }
      }
    },
    "/api/v1/auth/users": {
      "get": {
        "tags": ["Auth"],
        "summary": "Lista usuarios",
        "security": [{"bearerAuth": []}],
        "responses": {
          "200": {
            "description": "Lista retornada",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"$ref": "#/components/schemas/User"}
                }
              }
            }
          },
          "401": {"description": "Nao autorizado"}
        }
      }
    },
    "/api/v1/auth/users/{id}": {
      "get": {
        "tags": ["Auth"],
        "summary": "Busca usuario por id",
        "security": [{"bearerAuth": []}],
        "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
        "responses": {
          "200": {"description": "Usuario encontrado"},
          "401": {"description": "Nao autorizado"},
          "404": {"description": "Nao encontrado"}
        }
      },
      "put": {
        "tags": ["Auth"],
        "summary": "Atualiza usuario",
        "security": [{"bearerAuth": []}],
        "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/UpdateUserRequest"}
            }
          }
        },
        "responses": {
          "200": {"description": "Usuario atualizado"},
          "401": {"description": "Nao autorizado"},
          "404": {"description": "Nao encontrado"}
        }
      },
      "delete": {
        "tags": ["Auth"],
        "summary": "Desativa usuario",
        "security": [{"bearerAuth": []}],
        "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
        "responses": {
          "200": {"description": "Usuario desativado"},
          "401": {"description": "Nao autorizado"},
          "404": {"description": "Nao encontrado"}
        }
      }
    }
  }
}`

	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(spec))
	}
}

func UIHandler(specPath string) gin.HandlerFunc {
	html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>AggiPay Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "` + specPath + `",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis],
      layout: 'BaseLayout'
    });
  </script>
</body>
</html>`

	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}
