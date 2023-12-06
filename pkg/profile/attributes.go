package profile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attribute struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name"`
	Domain map[string]int     `bson:"domain"`
}

type Attributes map[string]Counter

type Counter map[string]int

type AttributesRepository interface {
	Get(context.Context) (Attribute, error)
	Updater(context.Context, *Profile) (Attribute, error)
	Delete(context.Context, *Profile) (Attribute, error)
}
