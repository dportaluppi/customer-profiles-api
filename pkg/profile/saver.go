package profile

import (
	"context"
	errstack "github.com/pkg/errors"
)

// saver implements the entity saver service.
type saver struct {
	repo Repository
}

func NewSaver(repo Repository) *saver {
	return &saver{repo: repo}
}

func (s *saver) Create(ctx context.Context, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entities
	if entity == nil {
		return nil, ErrInvalid
	}

	p, err := s.repo.Upsert(ctx, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}

func (s *saver) Update(ctx context.Context, id string, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entity
	if id == "" {
		return nil, ErrIDMissing
	}

	oldEntity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	entity.ID = oldEntity.ID
	entity.CreatedAt = oldEntity.CreatedAt
	entity.UpdatedAt = oldEntity.UpdatedAt

	p, err := s.repo.Upsert(ctx, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}
