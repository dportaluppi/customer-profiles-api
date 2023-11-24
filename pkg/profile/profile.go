package profile

import (
	"context"
	"time"
)

type Profile struct {
	ID        string     `json:"profileId"`
	Name      string     `json:"name" binding:"required"`
	Gender    string     `json:"gender" binding:"required,oneof=Male Female Other"`
	Birthday  time.Time  `json:"birthday" binding:"required"`
	Location  string     `json:"location" binding:"required"`
	Contact   Contact    `json:"contact" binding:"required"`
	Loyalty   Loyalty    `json:"loyalty" binding:"required"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type Contact struct {
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type Loyalty struct {
	Level          string    `json:"level" binding:"required"`
	EnrolledAt     time.Time `json:"enrolledAt" binding:"required"`
	LastActivityAt time.Time `json:"lastActivityAt"`
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
}

type Repository interface {
	Updater(ctx context.Context, profile *Profile) (*Profile, error)
	GetByID(ctx context.Context, profileID string) (*Profile, error)
	Delete(ctx context.Context, profileID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Profile, int, error)
}
