package profile

import (
	"context"
	"time"
)

type Profile struct {
	ID         string             `json:"profileId,omitempty" bson:"_id,omitempty"`
	Channels   map[string]Channel `json:"channels" bson:"channels"`
	Attributes map[string]any     `json:"attributes" bson:"attributes"`
	CreatedAt  *time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt  *time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type Channel struct {
	ID         string `json:"id" bson:"id" binding:"required"`
	Attributes any    `json:"attributes" bson:"attributes"`
}

type Upserter interface {
	Create(ctx context.Context, profile *Profile) (*Profile, error)
	Update(ctx context.Context, id string, profile *Profile) (*Profile, error)
}

type Deleter interface {
	Delete(ctx context.Context, profileID string) error
}

type Getter interface {
	GetByID(ctx context.Context, profileID string) (*Profile, error)
	GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error)
	Query(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Profile, int, error)
}

type Repository interface {
	Updater(ctx context.Context, profile *Profile) (*Profile, error)
	GetByID(ctx context.Context, profileID string) (*Profile, error)
	Delete(ctx context.Context, profileID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Profile, int, error)
}
