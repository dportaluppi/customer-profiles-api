package profile

import (
	"context"

	"github.com/pkg/errors"
)

// getter implements the entity retrieval service.
type getter struct {
	repo Repository
	attr AttributesRepository
}

func NewGetter(repo Repository, attr AttributesRepository) Getter {
	return &getter{
		repo: repo,
		attr: attr,
	}
}

func (s *getter) GetByID(ctx context.Context, accountID, id string) (*Entity, error) {
	if id == "" {
		return nil, ErrIDMissing
	}
	if accountID == "" {
		return nil, ErrAccountIDMissing
	}

	p, err := s.repo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if p.AccountID != accountID {
		return nil, ErrInvalid
	}

	return p, nil
}

func (s *getter) GetAll(ctx context.Context, accountId string, page, limit int) ([]*Entity, int, error) {
	if page < 1 || limit < 1 {
		return nil, 0, ErrInvalidPaginationParameters
	}

	entities, count, err := s.repo.GetAll(ctx, accountId, page, limit)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return entities, count, nil
}

func (s *getter) Query(ctx context.Context, accountId string, query map[string]any, currentPage, perPage int) ([]*Entity, int, error) {
	// TODO: business logic to query entities, e.g. check semantic and syntactic validity of query
	return s.repo.ExecuteQuery(ctx, accountId, query, currentPage, perPage)
}

func (s *getter) Pipeline(ctx context.Context, accountId string, pipeline map[string]any, currentPage, perPage int) ([]*Entity, int, error) {
	// TODO: business logic to query entities, e.g. check semantic and syntactic validity of query
	return s.repo.ExecutePipeline(ctx, accountId, pipeline, currentPage, perPage)
}

func (s *getter) GetKeys(ctx context.Context, accountId string) (Attributes, error) {
	return s.attr.GetAll(ctx)
}
