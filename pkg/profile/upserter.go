package profile

import (
	"context"
	errstack "github.com/pkg/errors"
)

// upserter implements the profile upserter service.
type upserter struct {
	repo Repository
}

func NewUpserter(repo Repository) *upserter {
	return &upserter{repo: repo}
}

func (s *upserter) Create(ctx context.Context, profile *Profile) (*Profile, error) {
	// TODO: business logic to create a profile
	if profile == nil {
		return nil, ErrProfileInvalid
	}

	p, err := s.repo.Updater(ctx, profile)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}

func (s *upserter) Update(ctx context.Context, id string, profile *Profile) (*Profile, error) {
	// TODO: business logic to create a profile
	if id == "" {
		return nil, ErrProfileIDMissing
	}

	oldProfile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	profile.ID = oldProfile.ID
	profile.CreatedAt = oldProfile.CreatedAt
	profile.UpdatedAt = oldProfile.UpdatedAt

	p, err := s.repo.Updater(ctx, profile)
	if err != nil {
		return nil, errstack.WithStack(err)
	}
	return p, nil
}
