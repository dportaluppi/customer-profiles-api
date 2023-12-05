package profile

import (
	"context"
	errstack "github.com/pkg/errors"
)

// saver implements the profile saver service.
type saver struct {
	repo Repository
}

func NewSaver(repo Repository) *saver {
	return &saver{repo: repo}
}

func (s *saver) Create(ctx context.Context, profile *Profile) (*Profile, error) {
	// TODO: business logic to create a profile
	if profile == nil {
		return nil, ErrInvalid
	}

	p, err := s.repo.Upsert(ctx, profile)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}

func (s *saver) Update(ctx context.Context, id string, profile *Profile) (*Profile, error) {
	// TODO: business logic to create a profile
	if id == "" {
		return nil, ErrIDMissing
	}

	oldProfile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	profile.ID = oldProfile.ID
	profile.CreatedAt = oldProfile.CreatedAt
	profile.UpdatedAt = oldProfile.UpdatedAt

	p, err := s.repo.Upsert(ctx, profile)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}
