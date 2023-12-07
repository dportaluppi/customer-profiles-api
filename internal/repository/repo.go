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
	Upsert(ctx context.Context, accountId string, entity T) (T, error)
	GetByID(ctx context.Context, accountId, id string) (T, error)
	Delete(ctx context.Context, accountId, id string) error
	GetAll(ctx context.Context, accountId string, page, limit int) ([]T, int, error)
	ExecuteQuery(ctx context.Context, accountId string, query map[string]interface{}, currentPage, perPage int) ([]T, int, error)
	ExecutePipeline(ctx context.Context, accountId string, pipeline map[string]interface{}, currentPage, perPage int) ([]T, int, error)
}
