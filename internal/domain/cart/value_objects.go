package cart

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// CartID represents a unique identifier for a shopping cart (Value Object)
type CartID struct {
	value string
}

func NewCartID() CartID {
	return CartID{value: uuid.New().String()}
}

func NewCartIDFromString(id string) (CartID, error) {
	if id == "" {
		return CartID{}, errors.New("cart ID cannot be empty")
	}
	return CartID{value: id}, nil
}

func (id CartID) String() string {
	return id.value
}

func (id CartID) Equals(other CartID) bool {
	return id.value == other.value
}

// CartStatus represents the status of a shopping cart (Value Object)
type CartStatus string

const (
	CartStatusActive    CartStatus = "ACTIVE"
	CartStatusAbandoned CartStatus = "ABANDONED"
	CartStatusCheckedOut CartStatus = "CHECKED_OUT"
	CartStatusExpired   CartStatus = "EXPIRED"
)

func (s CartStatus) String() string {
	return string(s)
}

func (s CartStatus) IsValid() bool {
	switch s {
	case CartStatusActive, CartStatusAbandoned, CartStatusCheckedOut, CartStatusExpired:
		return true
	default:
		return false
	}
}

// CartItemAddedEvent represents a domain event when an item is added to cart
type CartItemAddedEvent struct {
	CartID    CartID
	ProductID order.ProductID
	Quantity  order.Quantity
	UnitPrice order.Money
	AddedAt   time.Time
}

// CartItemRemovedEvent represents a domain event when an item is removed from cart
type CartItemRemovedEvent struct {
	CartID    CartID
	ProductID order.ProductID
	RemovedAt time.Time
}

// CartCheckedOutEvent represents a domain event when cart is checked out
type CartCheckedOutEvent struct {
	CartID      CartID
	CustomerID  order.CustomerID
	Total       order.Money
	CheckedOutAt time.Time
}

// CartAbandonedEvent represents a domain event when cart is abandoned
type CartAbandonedEvent struct {
	CartID      CartID
	CustomerID  order.CustomerID
	AbandonedAt time.Time
}