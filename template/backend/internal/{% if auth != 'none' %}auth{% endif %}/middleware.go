package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware validates the Bearer token and sets user_id and role on the
// context. Apply it to route groups that require authentication.
func (s *Service) Middleware() gin.HandlerFunc {
	const prefix = "Bearer "
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := s.tokens.Parse(strings.TrimPrefix(header, prefix))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.Subject)
		c.Set("role", claims.Role)
		c.Next()
	}
}
