package items

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/jwt-basic-app/backend/internal/database"
)

// Store wraps the sqlc-generated queries. The Item type and the ListItems /
// CreateItem / GetItem methods are generated from sql/queries/items.sql.
type Store struct {
	q *database.Queries
}

func NewStore(db *database.DB) *Store { return &Store{q: database.New(db.Pool)} }

func (s *Store) List(ctx context.Context) ([]database.Item, error) {
	return s.q.ListItems(ctx)
}

func (s *Store) Create(ctx context.Context, name string) (database.Item, error) {
	return s.q.CreateItem(ctx, name)
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (database.Item, error) {
	return s.q.GetItem(ctx, id)
}
