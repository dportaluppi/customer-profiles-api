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

func (s *saver) Create(ctx context.Context, accountId string, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entities
	if accountId == "" {
		return nil, ErrAccountIDMissing
	}
	if entity == nil {
		return nil, ErrInvalid
	}
	entity.AccountID = accountId
	p, err := s.repo.Upsert(ctx, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}

func (s *saver) Update(ctx context.Context, accountId, id string, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entity
	if accountId == "" {
		return nil, ErrAccountIDMissing
	}
	if id == "" {
		return nil, ErrIDMissing
	}

	oldEntity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	if oldEntity.AccountID != accountId {
		return nil, ErrInvalid
	}

	entity.ID = oldEntity.ID
	entity.AccountID = oldEntity.AccountID
	entity.CreatedAt = oldEntity.CreatedAt
	entity.UpdatedAt = oldEntity.UpdatedAt

	p, err := s.repo.Upsert(ctx, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}
