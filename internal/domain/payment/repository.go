package payment

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// Repository defines the interface for payment persistence (Repository Pattern)
type Repository interface {
	Save(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id PaymentID) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID order.OrderID) ([]*Payment, error)
	FindByCustomerID(ctx context.Context, customerID order.CustomerID) ([]*Payment, error)
	FindPendingPayments(ctx context.Context) ([]*Payment, error)
	Delete(ctx context.Context, id PaymentID) error
}

// PaymentMethodRepository defines the interface for payment method persistence
type PaymentMethodRepository interface {
	Save(ctx context.Context, method *PaymentMethod) error
	FindByID(ctx context.Context, id PaymentMethodID) (*PaymentMethod, error)
	FindByCustomerID(ctx context.Context, customerID order.CustomerID) ([]*PaymentMethod, error)
	FindDefaultByCustomerID(ctx context.Context, customerID order.CustomerID) (*PaymentMethod, error)
	Delete(ctx context.Context, id PaymentMethodID) error
}

// PaymentProcessor defines the interface for external payment processing
type PaymentProcessor interface {
	Authorize(ctx context.Context, payment *Payment, method *PaymentMethod) (string, order.Money, map[string]string, error)
	Capture(ctx context.Context, payment *Payment, amount order.Money) (string, map[string]string, error)
	Refund(ctx context.Context, payment *Payment, amount order.Money, reason RefundReason) (string, map[string]string, error)
	Void(ctx context.Context, payment *Payment) error
}

// EventPublisher defines the interface for publishing payment domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}