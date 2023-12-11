package profile

import (
	"context"
)

type Attributes map[string][]string

type AttributesRepository interface {
	GetAll(ctx context.Context, accountID string) (Attributes, error)
	Updater(ctx context.Context, accountID string, e *Entity) error
	Delete(ctx context.Context, accountID string, e *Entity) error
}
