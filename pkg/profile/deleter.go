package profile

import (
	"context"
	"github.com/pkg/errors"
)

// deleter implements the entity deletion service.
type deleter struct {
	repo Repository
}

func NewDeleter(repo Repository) *deleter {
	return &deleter{repo: repo}
}

func (s *deleter) Delete(ctx context.Context, accountId, id string) error {
	if id == "" {
		return ErrIDMissing
	}
	e, err := s.repo.GetByID(ctx, id)
	if e == nil {
		return errors.WithStack(err)
	}

	if e.AccountID != accountId {
		return ErrInvalid
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
