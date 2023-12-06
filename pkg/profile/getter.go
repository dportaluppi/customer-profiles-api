package profile

import (
	"context"

	"github.com/pkg/errors"
)

// getter implements the profile retrieval service.
type getter struct {
	repo Repository
}

func NewGetter(repo Repository) Getter {
	return &getter{repo: repo}
}

func (s *getter) GetByID(ctx context.Context, id string) (*Profile, error) {
	if id == "" {
		return nil, ErrProfileIDMissing
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return p, nil
}

func (s *getter) GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error) {
	if page < 1 || limit < 1 {
		return nil, 0, ErrInvalidPaginationParameters
	}

	profiles, count, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return profiles, count, nil
}

func (s *getter) Query(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Profile, int, error) {
	// TODO: business logic to query profiles, e.g. check semantic and syntactic validity of query
	return s.repo.ExecuteQuery(ctx, query, currentPage, perPage)
}

func (s *getter) GetKeys(ctx context.Context) (map[string][]any, error) {
	return s.repo.GetKeys(ctx)
}
