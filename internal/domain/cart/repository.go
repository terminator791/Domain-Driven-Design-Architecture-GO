package cart

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Repository defines the interface for cart persistence (Repository Pattern)
type Repository interface {
	Save(ctx context.Context, cart *Cart) error
	FindByID(ctx context.Context, id CartID) (*Cart, error)
	FindByCustomerID(ctx context.Context, customerID order.CustomerID) (*Cart, error)
	FindActiveByCustomerID(ctx context.Context, customerID order.CustomerID) (*Cart, error)
	FindExpiredCarts(ctx context.Context) ([]*Cart, error)
	Delete(ctx context.Context, id CartID) error
}

// EventPublisher defines the interface for publishing cart domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}