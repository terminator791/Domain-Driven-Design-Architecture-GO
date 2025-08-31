package inventory

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// ReservationID represents a unique identifier for a stock reservation (Value Object)
type ReservationID struct {
	value string
}

func NewReservationID() ReservationID {
	return ReservationID{value: uuid.New().String()}
}

func NewReservationIDFromString(id string) (ReservationID, error) {
	if id == "" {
		return ReservationID{}, errors.New("reservation ID cannot be empty")
	}
	return ReservationID{value: id}, nil
}

func (id ReservationID) String() string {
	return id.value
}

func (id ReservationID) Equals(other ReservationID) bool {
	return id.value == other.value
}

// StockLevel represents stock quantity (Value Object)
type StockLevel struct {
	available int
	reserved  int
	total     int
}

func NewStockLevel(available, reserved int) StockLevel {
	return StockLevel{
		available: available,
		reserved:  reserved,
		total:     available + reserved,
	}
}

func (s StockLevel) Available() int {
	return s.available
}

func (s StockLevel) Reserved() int {
	return s.reserved
}

func (s StockLevel) Total() int {
	return s.total
}

func (s StockLevel) CanReserve(quantity int) bool {
	return s.available >= quantity
}

// ReservationStatus represents the status of a stock reservation (Value Object)
type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "ACTIVE"
	ReservationStatusExpired  ReservationStatus = "EXPIRED"
	ReservationStatusCommitted ReservationStatus = "COMMITTED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
)

func (s ReservationStatus) String() string {
	return string(s)
}

func (s ReservationStatus) IsValid() bool {
	switch s {
	case ReservationStatusActive, ReservationStatusExpired, ReservationStatusCommitted, ReservationStatusCancelled:
		return true
	default:
		return false
	}
}

// AlertLevel represents different inventory alert levels (Value Object)
type AlertLevel string

const (
	AlertLevelNone     AlertLevel = "NONE"
	AlertLevelLow      AlertLevel = "LOW"
	AlertLevelCritical AlertLevel = "CRITICAL"
	AlertLevelOutOfStock AlertLevel = "OUT_OF_STOCK"
)

func (a AlertLevel) String() string {
	return string(a)
}

// StockReservedEvent represents a domain event when stock is reserved
type StockReservedEvent struct {
	ProductID     order.ProductID
	ReservationID ReservationID
	Quantity      int
	ReservedBy    order.CustomerID
	ReservedAt    time.Time
	ExpiresAt     time.Time
}

// StockCommittedEvent represents a domain event when reserved stock is committed
type StockCommittedEvent struct {
	ProductID     order.ProductID
	ReservationID ReservationID
	Quantity      int
	CommittedAt   time.Time
}

// StockReleasedEvent represents a domain event when reserved stock is released
type StockReleasedEvent struct {
	ProductID     order.ProductID
	ReservationID ReservationID
	Quantity      int
	ReleasedAt    time.Time
	Reason        string
}

// StockReplenishedEvent represents a domain event when stock is replenished
type StockReplenishedEvent struct {
	ProductID      order.ProductID
	QuantityAdded  int
	NewStockLevel  int
	ReplenishedAt  time.Time
	ReplenishedBy  string
}

// LowStockAlertEvent represents a domain event when stock is low
type LowStockAlertEvent struct {
	ProductID     order.ProductID
	CurrentStock  int
	MinimumStock  int
	AlertLevel    AlertLevel
	AlertedAt     time.Time
}