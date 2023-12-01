package user

import (
	"context"
	"github.com/pkg/errors"
)

// deleter implements the profile deletion service.
type deleter struct {
	repo Repository
}

func NewDeleter(repo Repository) *deleter {
	return &deleter{repo: repo}
}

func (s *deleter) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrProfileIDMissing
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
