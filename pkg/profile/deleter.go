package profile

import (
	"context"

	"github.com/pkg/errors"
)

// deleter implements the profile deletion service.
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

func (s *deleter) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrProfileIDMissing
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	if err = s.attr.Delete(ctx, p); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
