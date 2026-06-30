package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/kitchen-sink-app/backend/internal/database"
)

// Store wraps the sqlc-generated user queries (sql/queries/auth.sql).
type Store struct {
	q *database.Queries
}

func NewStore(db *database.DB) *Store { return &Store{q: database.New(db.Pool)} }

func (s *Store) CreateUser(ctx context.Context, email, hash, name, role string) (database.User, error) {
	return s.q.CreateUser(ctx, database.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		Name:         name,
		Role:         role,
	})
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return s.q.GetUserByEmail(ctx, email)
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return s.q.GetUserByID(ctx, id)
}
