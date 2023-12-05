package store

import (
	"context"
	"time"
)

type Store struct {
	ID         string              `json:"storeId,omitempty"`
	Providers  map[string]Provider `json:"providers" bson:"providers"`
	Attributes map[string]any      `json:"attributes" bson:"attributes"`
	CreatedAt  *time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  *time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Provider struct {
	ID         string         `json:"id" bson:"id" binding:"required"`
	Attributes map[string]any `json:"attributes" bson:"attributes"`
}

func (s *Store) GetID() string {
	return s.ID
}

func (s *Store) SetID(id string) {
	s.ID = id
}

func (s *Store) GetCreatedAt() *time.Time {
	return s.CreatedAt
}

func (s *Store) SetCreatedAt(t time.Time) {
	s.CreatedAt = &t
}

func (s *Store) GetUpdatedAt() *time.Time {
	return s.UpdatedAt
}

func (s *Store) SetUpdatedAt(t time.Time) {
	s.UpdatedAt = &t
}

type Upserter interface {
	Create(ctx context.Context, store *Store) (*Store, error)
	Update(ctx context.Context, id string, store *Store) (*Store, error)
}

type Deleter interface {
	Delete(ctx context.Context, storeID string) error
}

type Getter interface {
	GetByID(ctx context.Context, storeID string) (*Store, error)
	GetAll(ctx context.Context, page, limit int) ([]*Store, int, error)
	Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*Store, int, error)
	Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Store, int, error)
}

type Repository interface {
	Upsert(ctx context.Context, store *Store) (*Store, error)
	GetByID(ctx context.Context, storeID string) (*Store, error)
	Delete(ctx context.Context, storeID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Store, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Store, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Store, int, error)
}
