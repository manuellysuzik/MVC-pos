package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware protege todas as rotas com autenticação via header X-API-Key.
// Recebe a chave por parâmetro — a validação de que ela não está vazia
// deve ser feita no startup (main.go) antes de chamar SetupRouter.
func APIKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
