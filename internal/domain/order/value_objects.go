package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderID represents a unique identifier for an order (Value Object)
type OrderID struct {
	value string
}

func NewOrderID() OrderID {
	return OrderID{value: uuid.New().String()}
}

func NewOrderIDFromString(id string) (OrderID, error) {
	if id == "" {
		return OrderID{}, errors.New("order ID cannot be empty")
	}
	return OrderID{value: id}, nil
}

func (id OrderID) String() string {
	return id.value
}

func (id OrderID) Equals(other OrderID) bool {
	return id.value == other.value
}

// ProductID represents a unique identifier for a product (Value Object)
type ProductID struct {
	value string
}

func NewProductID() ProductID {
	return ProductID{value: uuid.New().String()}
}

func NewProductIDFromString(id string) (ProductID, error) {
	if id == "" {
		return ProductID{}, errors.New("product ID cannot be empty")
	}
	return ProductID{value: id}, nil
}

func (id ProductID) String() string {
	return id.value
}

func (id ProductID) Equals(other ProductID) bool {
	return id.value == other.value
}

// CustomerID represents a unique identifier for a customer (Value Object)
type CustomerID struct {
	value string
}

func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

func NewCustomerIDFromString(id string) (CustomerID, error) {
	if id == "" {
		return CustomerID{}, errors.New("customer ID cannot be empty")
	}
	return CustomerID{value: id}, nil
}

func (id CustomerID) String() string {
	return id.value
}

func (id CustomerID) Equals(other CustomerID) bool {
	return id.value == other.value
}

// Money represents monetary value (Value Object)
type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Multiply(factor int) Money {
	return Money{amount: m.amount * float64(factor), currency: m.currency}
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Quantity represents a quantity of items (Value Object)
type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value <= 0 {
		return Quantity{}, errors.New("quantity must be positive")
	}
	return Quantity{value: value}, nil
}

func (q Quantity) Value() int {
	return q.value
}

func (q Quantity) Add(other Quantity) Quantity {
	return Quantity{value: q.value + other.value}
}

func (q Quantity) Equals(other Quantity) bool {
	return q.value == other.value
}

// OrderStatus represents the status of an order (Value Object)
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// OrderCreatedEvent represents a domain event when an order is created
type OrderCreatedEvent struct {
	OrderID   OrderID
	Total     Money
	CreatedAt time.Time
}

// OrderUpdatedEvent represents a domain event when an order is updated
type OrderUpdatedEvent struct {
	OrderID   OrderID
	Status    OrderStatus
	UpdatedAt time.Time
}