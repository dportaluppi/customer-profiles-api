package store

import (
	"context"
	"github.com/pkg/errors"
)

// getter implements the store retrieval service.
type getter struct {
	repo Repository
}

func NewGetter(repo Repository) Getter {
	return &getter{repo: repo}
}

func (s *getter) GetByID(ctx context.Context, id string) (*Store, error) {
	if id == "" {
		return nil, ErrIDMissing
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return p, nil
}

func (s *getter) GetAll(ctx context.Context, page, limit int) ([]*Store, int, error) {
	if page < 1 || limit < 1 {
		return nil, 0, ErrInvalidPaginationParameters
	}

	stores, count, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return stores, count, nil
}

func (s *getter) Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*Store, int, error) {
	// TODO: business logic to query stores, e.g. check semantic and syntactic validity of query
	return s.repo.ExecuteQuery(ctx, query, currentPage, perPage)
}

func (s *getter) Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Store, int, error) {
	// TODO: business logic to query stores, e.g. check semantic and syntactic validity of query
	return s.repo.ExecutePipeline(ctx, pipeline, currentPage, perPage)
}
