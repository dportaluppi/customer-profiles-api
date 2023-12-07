package profile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attributes map[string][]string

type Attribute struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Attribute string             `bson:"attribute"`
	Value     string             `bson:"value"`
	Count     int                `bson:"count"`
}

type AttributesRepository interface {
	Get(context.Context) (Attribute, error)
	Updater(context.Context, *Profile) (Attribute, error)
	Delete(context.Context, *Profile) (Attribute, error)
}
