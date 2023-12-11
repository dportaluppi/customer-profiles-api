package profile

import (
	"context"

	"github.com/pkg/errors"
)

// deleter implements the entity deletion service.
type deleter struct {
	repo Repository
	attr AttributesRepository
}

func NewDeleter(repo Repository, attr AttributesRepository) *deleter {
	return &deleter{
		repo: repo,
		attr: attr,
	}
}

func (s *deleter) Delete(ctx context.Context, accountID, id string) error {
	if id == "" {
		return ErrIDMissing
	}
	e, err := s.repo.GetByID(ctx, accountID, id)
	if e == nil {
		return errors.WithStack(err)
	}

	if e.AccountID != accountID {
		return ErrInvalid
	}

	if err = s.repo.Delete(ctx, accountID, id); err != nil {
		return errors.WithStack(err)
	}

	if err = s.attr.Delete(ctx, accountID, e); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
