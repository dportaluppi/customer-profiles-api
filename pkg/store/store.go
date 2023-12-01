package store

import "time"

type Store struct {
	ID         string              `json:"storeId,omitempty" bson:"_id,omitempty"`
	Providers  map[string]Provider `json:"providers" bson:"providers"`
	Attributes map[string]any      `json:"attributes" bson:"attributes"`
	CreatedAt  *time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  *time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Provider struct {
	ID         string         `json:"id" bson:"id" binding:"required"`
	Attributes map[string]any `json:"attributes" bson:"attributes"`
}
