// Package auth is a kept (non-example) module: concrete users schema, JWT
// signup/login, backed by sqlc queries.
package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/example/jwt-basic-app/backend/gen/api/auth"
	"github.com/example/jwt-basic-app/backend/internal/config"
	"github.com/example/jwt-basic-app/backend/internal/database"
)

type Service struct {
	store      *Store
	tokens     *TokenIssuer
	refreshTTL time.Duration
}

func refreshTTL(cfg config.TokenConfig) time.Duration {
	h := cfg.RefreshExpiryHours
	if h == 0 {
		h = 720 // 30 days
	}
	return time.Duration(h) * time.Hour
}

func New(db *database.DB, cfg config.TokenConfig) *Service {
	return &Service{
		store:      NewStore(db),
		tokens:     NewTokenIssuer(cfg),
		refreshTTL: refreshTTL(cfg),
	}
}

// Register mounts the generated routes under the given router group.
func (s *Service) Register(r gin.IRouter) {
	gen.RegisterHandlers(r, s)
}

// Signup implements gen.ServerInterface.
func (s *Service) Signup(c *gin.Context) {
	var body gen.SignupJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Email == "" || len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password (min 8 chars) are required"})
		return
	}

	hash, err := HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	name := ""
	if body.Name != nil {
		name = *body.Name
	}

	user, err := s.store.CreateUser(c.Request.Context(), string(body.Email), hash, name, "user")
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	s.respondWithSession(c, http.StatusCreated, user)
}

// Login implements gen.ServerInterface.
func (s *Service) Login(c *gin.Context) {
	var body gen.LoginJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, err := s.store.GetUserByEmail(c.Request.Context(), string(body.Email))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !CheckPassword(user.PasswordHash, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	s.respondWithSession(c, http.StatusOK, user)
}

func toAPIUser(u database.User) gen.User {
	created := u.CreatedAt
	return gen.User{
		Id:        u.ID,
		Email:     openapi_types.Email(u.Email),
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: &created,
	}
}
