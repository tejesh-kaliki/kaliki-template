package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// authenticate validates the Bearer token. On success it stores user_id/role on
// the context and returns the claims; on failure it writes a 401 and returns
// false. Use it from a handler (see GetCurrentUser) or via Middleware.
func (s *Service) authenticate(c *gin.Context) (*Claims, bool) {
	const prefix = "Bearer "
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, prefix) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return nil, false
	}
	claims, err := s.tokens.Parse(strings.TrimPrefix(header, prefix))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return nil, false
	}
	c.Set("user_id", claims.Subject)
	c.Set("role", claims.Role)
	return claims, true
}

// Middleware validates the Bearer token and sets user_id and role on the
// context. Apply it to route groups that require authentication.
func (s *Service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := s.authenticate(c); ok {
			c.Next()
		}
	}
}
