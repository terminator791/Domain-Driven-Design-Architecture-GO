package order

import "context"

// Repository defines the interface for order persistence (Repository Pattern)
// This is part of the domain layer but implementation will be in infrastructure
type Repository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id OrderID) (*Order, error)
	FindByCustomerID(ctx context.Context, customerID CustomerID) ([]*Order, error)
	Delete(ctx context.Context, id OrderID) error
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}