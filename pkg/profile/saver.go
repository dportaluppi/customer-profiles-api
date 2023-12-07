package profile

import (
	"context"

	errstack "github.com/pkg/errors"
)

// saver implements the entity saver service.
type saver struct {
	repo Repository
	attr AttributesRepository
}

func NewSaver(repo Repository, attr AttributesRepository) *saver {
	return &saver{
		repo: repo,
		attr: attr,
	}
}

func (s *saver) Create(ctx context.Context, accountID string, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entities
	if accountID == "" {
		return nil, ErrAccountIDMissing
	}

	if entity == nil {
		return nil, ErrInvalid
	}

	entity.AccountID = accountID

	p, err := s.repo.Upsert(ctx, accountID, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	err = s.attr.Updater(ctx, p)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	return p, nil
}

func (s *saver) Update(ctx context.Context, accountID, id string, entity *Entity) (*Entity, error) {
	// TODO: business logic to create a entity
	if accountID == "" {
		return nil, ErrAccountIDMissing
	}
	if id == "" {
		return nil, ErrIDMissing
	}

	oldEntity, err := s.repo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	if oldEntity.AccountID != accountID {
		return nil, ErrInvalid
	}

	err = s.attr.Delete(ctx, oldEntity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	entity.ID = oldEntity.ID
	entity.AccountID = oldEntity.AccountID
	entity.CreatedAt = oldEntity.CreatedAt
	entity.UpdatedAt = oldEntity.UpdatedAt

	p, err := s.repo.Upsert(ctx, accountID, entity)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	err = s.attr.Updater(ctx, p)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	return p, nil
}

func (s *saver) AddRelationship(ctx context.Context, accountId, id string, relationship Relationship) (*Entity, error) {
	e, err := s.repo.GetByID(ctx, accountId, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	if !e.Add(relationship) {
		return e, nil
	}

	err = s.attr.Delete(ctx, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	e, err = s.repo.Upsert(ctx, accountId, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	err = s.attr.Updater(ctx, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	return e, nil
}

func (s *saver) ReplaceRelationships(ctx context.Context, accountId, id string, relationships []Relationship) (*Entity, error) {
	e, err := s.repo.GetByID(ctx, accountId, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	err = s.attr.Delete(ctx, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	e.Relationships = relationships

	e, err = s.repo.Upsert(ctx, accountId, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	err = s.attr.Updater(ctx, e)
	if err != nil {
		return nil, errstack.WithStack(err)
	}

	return e, nil
}
