package repository

import (
	"context"
	"time"
)

// Entity is an interface that all entities must implement.
// Each entity should be capable of handling its own ID.
type Entity interface {
	GetID() string
	SetID(id string)
	GetCreatedAt() *time.Time
	SetCreatedAt(t time.Time)
	GetUpdatedAt() *time.Time
	SetUpdatedAt(t time.Time)
}

// Repository is a generic interface for a repository.
type Repository[T any] interface {
	Upsert(ctx context.Context, entity T) (T, error)
	GetByID(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, page, limit int) ([]T, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]T, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]interface{}, currentPage, perPage int) ([]T, int, error)
}
