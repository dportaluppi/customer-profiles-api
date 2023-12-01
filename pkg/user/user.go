package user

import (
	"context"
	"time"
)

type User struct {
	ID         string             `json:"userId,omitempty"`
	Channels   map[string]Channel `json:"channels" bson:"channels"`
	Attributes map[string]any     `json:"attributes" bson:"attributes"`
	CreatedAt  *time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt  *time.Time         `json:"updatedAt" bson:"updatedAt"`
}

func (u *User) GetCreatedAt() *time.Time {
	return u.CreatedAt
}

func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = &t
}

func (u *User) GetUpdatedAt() *time.Time {
	return u.UpdatedAt
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = &t
}

// GetID returns the ID of the user.
func (u *User) GetID() string {
	return u.ID
}

// SetID sets the ID of the user.
func (u *User) SetID(id string) {
	u.ID = id
}

type Channel struct {
	ID         string         `json:"id" bson:"id" binding:"required"`
	Attributes map[string]any `json:"attributes" bson:"attributes"`
}

type Upserter interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, id string, user *User) (*User, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID string) error
}

type Getter interface {
	GetByID(ctx context.Context, userID string) (*User, error)
	GetAll(ctx context.Context, page, limit int) ([]*User, int, error)
	Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*User, int, error)
	Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*User, int, error)
}

type Repository interface {
	Upsert(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, userID string) (*User, error)
	Delete(ctx context.Context, userID string) error
	GetAll(ctx context.Context, page, limit int) ([]*User, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*User, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*User, int, error)
}
