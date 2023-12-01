package profile

import (
	"context"
	"time"
)

type Profile struct {
	ID         string             `json:"profileId,omitempty"`
	Channels   map[string]Channel `json:"channels" bson:"channels"`
	Attributes map[string]any     `json:"attributes" bson:"attributes"`
	CreatedAt  *time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt  *time.Time         `json:"updatedAt" bson:"updatedAt"`
}

func (u *Profile) GetCreatedAt() *time.Time {
	return u.CreatedAt
}

func (u *Profile) SetCreatedAt(t time.Time) {
	u.CreatedAt = &t
}

func (u *Profile) GetUpdatedAt() *time.Time {
	return u.UpdatedAt
}

func (u *Profile) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = &t
}

// GetID returns the ID of the profile.
func (u *Profile) GetID() string {
	return u.ID
}

// SetID sets the ID of the profile.
func (u *Profile) SetID(id string) {
	u.ID = id
}

type Channel struct {
	ID         string         `json:"id" bson:"id" binding:"required"`
	Attributes map[string]any `json:"attributes" bson:"attributes"`
}

type Upserter interface {
	Create(ctx context.Context, user *Profile) (*Profile, error)
	Update(ctx context.Context, id string, user *Profile) (*Profile, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID string) error
}

type Getter interface {
	GetByID(ctx context.Context, userID string) (*Profile, error)
	GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error)
	Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*Profile, int, error)
	Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Profile, int, error)
}

type Repository interface {
	Upsert(ctx context.Context, user *Profile) (*Profile, error)
	GetByID(ctx context.Context, userID string) (*Profile, error)
	Delete(ctx context.Context, userID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Profile, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Profile, int, error)
}
