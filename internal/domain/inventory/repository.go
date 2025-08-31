package inventory

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Repository defines the interface for inventory persistence (Repository Pattern)
type Repository interface {
	Save(ctx context.Context, item *InventoryItem) error
	FindByProductID(ctx context.Context, productID order.ProductID) (*InventoryItem, error)
	FindLowStockItems(ctx context.Context) ([]*InventoryItem, error)
	FindItemsWithExpiredReservations(ctx context.Context) ([]*InventoryItem, error)
	FindAll(ctx context.Context) ([]*InventoryItem, error)
}

// ReservationRepository defines the interface for reservation persistence
type ReservationRepository interface {
	Save(ctx context.Context, reservation *StockReservation) error
	FindByID(ctx context.Context, id ReservationID) (*StockReservation, error)
	FindByCustomerID(ctx context.Context, customerID order.CustomerID) ([]*StockReservation, error)
	FindExpiredReservations(ctx context.Context) ([]*StockReservation, error)
	Delete(ctx context.Context, id ReservationID) error
}

// EventPublisher defines the interface for publishing inventory domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}