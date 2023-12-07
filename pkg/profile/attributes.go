package profile

import (
	"context"
)

type Attributes map[string][]string

type AttributesRepository interface {
	GetAll(context.Context) (Attributes, error)
	Updater(context.Context, *Profile) error
	Delete(context.Context, *Profile) error
}
