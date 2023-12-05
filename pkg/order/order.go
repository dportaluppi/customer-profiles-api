package order

import (
	"context"
	"time"
)

type Order struct {
	UserID   *string `json:"userId,omitempty" bson:"userId,omitempty"`
	StoreID  string  `json:"storeId" bson:"storeId"`
	Provider string  `json:"provider" bson:"provider"` // FEMSA, Mondelez, etc

	OrderID          string    `json:"orderId" bson:"orderId"`
	OrderDate        time.Time `json:"orderDate" bson:"order_date"`
	OrderChannel     string    `json:"channel" bson:"order_channel"` // yalo, offline, etc
	OrderTotalAmount float64   `json:"total_amount" bson:"total_amount"`
}

type Upserter interface {
	Create(ctx context.Context, Order *Order) (*Order, error)
	Update(ctx context.Context, id string, Order *Order) (*Order, error)
}

type Deleter interface {
	Delete(ctx context.Context, OrderID string) error
}

type Getter interface {
	GetByID(ctx context.Context, OrderID string) (*Order, error)
	GetAll(ctx context.Context, page, limit int) ([]*Order, int, error)
	Query(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*Order, int, error)
	Pipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Order, int, error)
}

type Repository interface {
	Upsert(ctx context.Context, Order *Order) (*Order, error)
	GetByID(ctx context.Context, OrderID string) (*Order, error)
	Delete(ctx context.Context, OrderID string) error
	GetAll(ctx context.Context, page, limit int) ([]*Order, int, error)
	ExecuteQuery(ctx context.Context, query map[string]interface{}, currentPage, perPage int) ([]*Order, int, error)
	ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*Order, int, error)
}
