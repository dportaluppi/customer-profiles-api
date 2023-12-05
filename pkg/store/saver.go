package store

import (
	"context"
	errstack "github.com/pkg/errors"
)

// saver implements the store saver service.
type saver struct {
	repo Repository
}

func NewSaver(repo Repository) *saver {
	return &saver{repo: repo}
}

func (s *saver) Create(ctx context.Context, store *Store) (*Store, error) {
	// TODO: business logic to create a store
	if store == nil {
		return nil, ErrInvalid
	}

	p, err := s.repo.Upsert(ctx, store)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}

func (s *saver) Update(ctx context.Context, id string, store *Store) (*Store, error) {
	// TODO: business logic to create a store
	if id == "" {
		return nil, ErrIDMissing
	}

	oldstore, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	store.ID = oldstore.ID
	store.CreatedAt = oldstore.CreatedAt
	store.UpdatedAt = oldstore.UpdatedAt

	p, err := s.repo.Upsert(ctx, store)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}
