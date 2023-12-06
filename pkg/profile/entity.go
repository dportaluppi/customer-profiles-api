package profile

import (
	"context"
	"time"
)

// Metadata represents any additional metadata information.
type Metadata map[string]any

// Attribute represents specific attributes of an entity.
type Attribute map[string]any

// Relationship defines a connection between entities.
type Relationship struct {
	Type     string `json:"type" bson:"type"`         // Type of relationship, e.g., 'buysFrom', 'sellsFor'
	TargetID string `json:"targetId" bson:"targetId"` // ID of the target entity in the relationship
}

// Entity represents a generic structure for Contact or Store, associated with a specific account.
type Entity struct {
	ID            string         `json:"id"`                                 // Unique identifier for the entity
	AccountID     string         `json:"accountId" bson:"accountId"`         // ID of the associated account
	Metadata      Metadata       `json:"metadata" bson:"metadata"`           // Additional metadata for the entity
	Type          string         `json:"type" bson:"type"`                   // Type of the entity, e.g., 'Contact', 'Store'
	Attributes    Attribute      `json:"attributes" bson:"attributes"`       // Specific attributes of the entity
	Relationships []Relationship `json:"relationships" bson:"relationships"` // Relationships with other entities

	CreatedAt *time.Time `json:"createdAt" bson:"createdAt"` // Timestamp of entity creation
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"` // Timestamp of last entity update
}

// GetID returns the entity's unique identifier.
func (e *Entity) GetID() string {
	return e.ID
}

// SetID sets the entity's unique identifier.
func (e *Entity) SetID(id string) {
	e.ID = id
}

// GetCreatedAt returns the timestamp of when the entity was created.
func (e *Entity) GetCreatedAt() *time.Time {
	return e.CreatedAt
}

// SetCreatedAt sets the timestamp of when the entity was created.
func (e *Entity) SetCreatedAt(t time.Time) {
	e.CreatedAt = &t
}

// GetUpdatedAt returns the timestamp of the last update to the entity.
func (e *Entity) GetUpdatedAt() *time.Time {
	return e.UpdatedAt
}

// SetUpdatedAt sets the timestamp of the last update to the entity.
func (e *Entity) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = &t
}

type Upserter interface {
	Create(ctx context.Context, Entity *Entity) (*Entity, error)
	Update(ctx context.Context, id string, Entity *Entity) (*Entity, error)
}

type Deleter interface {
	Delete(ctx context.Context, EntityID string) error
}

type Getter interface {
	GetByID(ctx context.Context, EntityID string) (*Entity, error)
	GetAll(ctx context.Context, page, limit int) ([]*Entity, int, error)
	Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*Entity, int, error)
	Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Entity, int, error)
}

type Repository interface {
	Upsert(ctx context.Context, Entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, EntityID string) (*Entity, error)
	Delete(ctx context.Context, EntityID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Entity, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Entity, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Entity, int, error)
}
