package inventory

import (
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// StockReservation represents a reservation of stock for a specific customer (Entity)
type StockReservation struct {
	id         ReservationID
	productID  order.ProductID
	customerID order.CustomerID
	quantity   int
	status     ReservationStatus
	createdAt  time.Time
	expiresAt  time.Time
	updatedAt  time.Time
}

func NewStockReservation(productID order.ProductID, customerID order.CustomerID, quantity int, duration time.Duration) *StockReservation {
	now := time.Now()
	return &StockReservation{
		id:         NewReservationID(),
		productID:  productID,
		customerID: customerID,
		quantity:   quantity,
		status:     ReservationStatusActive,
		createdAt:  now,
		expiresAt:  now.Add(duration),
		updatedAt:  now,
	}
}

func (r *StockReservation) ID() ReservationID {
	return r.id
}

func (r *StockReservation) ProductID() order.ProductID {
	return r.productID
}

func (r *StockReservation) CustomerID() order.CustomerID {
	return r.customerID
}

func (r *StockReservation) Quantity() int {
	return r.quantity
}

func (r *StockReservation) Status() ReservationStatus {
	return r.status
}

func (r *StockReservation) CreatedAt() time.Time {
	return r.createdAt
}

func (r *StockReservation) ExpiresAt() time.Time {
	return r.expiresAt
}

func (r *StockReservation) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *StockReservation) IsExpired() bool {
	return time.Now().After(r.expiresAt)
}

func (r *StockReservation) Commit() error {
	if r.status != ReservationStatusActive {
		return common.ErrInvalidInput
	}
	
	if r.IsExpired() {
		return common.ErrReservationExpired
	}

	r.status = ReservationStatusCommitted
	r.updatedAt = time.Now()
	return nil
}

func (r *StockReservation) Cancel() error {
	if r.status != ReservationStatusActive {
		return common.ErrInvalidInput
	}

	r.status = ReservationStatusCancelled
	r.updatedAt = time.Now()
	return nil
}

func (r *StockReservation) Expire() error {
	if r.status != ReservationStatusActive {
		return common.ErrInvalidInput
	}

	r.status = ReservationStatusExpired
	r.updatedAt = time.Now()
	return nil
}

func (r *StockReservation) Extend(duration time.Duration) error {
	if r.status != ReservationStatusActive {
		return common.ErrInvalidInput
	}

	r.expiresAt = time.Now().Add(duration)
	r.updatedAt = time.Now()
	return nil
}

// InventoryItem represents the inventory state for a specific product (Aggregate Root Entity)
type InventoryItem struct {
	productID     order.ProductID
	stockLevel    StockLevel
	reservations  []StockReservation
	minStockLevel int
	maxStockLevel int
	reorderPoint  int
	lastUpdated   time.Time
	events        []interface{} // Domain events
}

func NewInventoryItem(productID order.ProductID, initialStock, minStock, maxStock, reorderPoint int) *InventoryItem {
	return &InventoryItem{
		productID:     productID,
		stockLevel:    NewStockLevel(initialStock, 0),
		reservations:  make([]StockReservation, 0),
		minStockLevel: minStock,
		maxStockLevel: maxStock,
		reorderPoint:  reorderPoint,
		lastUpdated:   time.Now(),
		events:        make([]interface{}, 0),
	}
}

func (i *InventoryItem) ProductID() order.ProductID {
	return i.productID
}

func (i *InventoryItem) StockLevel() StockLevel {
	return i.stockLevel
}

func (i *InventoryItem) Reservations() []StockReservation {
	return i.reservations
}

func (i *InventoryItem) MinStockLevel() int {
	return i.minStockLevel
}

func (i *InventoryItem) MaxStockLevel() int {
	return i.maxStockLevel
}

func (i *InventoryItem) ReorderPoint() int {
	return i.reorderPoint
}

func (i *InventoryItem) LastUpdated() time.Time {
	return i.lastUpdated
}

func (i *InventoryItem) Events() []interface{} {
	return i.events
}

func (i *InventoryItem) ClearEvents() {
	i.events = make([]interface{}, 0)
}

// ReserveStock reserves stock for a customer
func (i *InventoryItem) ReserveStock(customerID order.CustomerID, quantity int, duration time.Duration) (*StockReservation, error) {
	if !i.stockLevel.CanReserve(quantity) {
		return nil, common.ErrInsufficientStock
	}

	reservation := NewStockReservation(i.productID, customerID, quantity, duration)
	i.reservations = append(i.reservations, *reservation)
	
	// Update stock levels
	i.stockLevel = NewStockLevel(i.stockLevel.Available()-quantity, i.stockLevel.Reserved()+quantity)
	i.lastUpdated = time.Now()

	// Add domain event
	event := StockReservedEvent{
		ProductID:     i.productID,
		ReservationID: reservation.ID(),
		Quantity:      quantity,
		ReservedBy:    customerID,
		ReservedAt:    time.Now(),
		ExpiresAt:     reservation.ExpiresAt(),
	}
	i.events = append(i.events, event)

	return reservation, nil
}

// CommitReservation commits a reservation (converts reserved stock to sold)
func (i *InventoryItem) CommitReservation(reservationID ReservationID) error {
	for idx, reservation := range i.reservations {
		if reservation.ID().Equals(reservationID) {
			err := reservation.Commit()
			if err != nil {
				return err
			}

			// Update the reservation in slice
			i.reservations[idx] = reservation
			
			// Update stock levels (remove from reserved and total)
			i.stockLevel = NewStockLevel(i.stockLevel.Available(), i.stockLevel.Reserved()-reservation.Quantity())
			i.lastUpdated = time.Now()

			// Add domain event
			event := StockCommittedEvent{
				ProductID:     i.productID,
				ReservationID: reservationID,
				Quantity:      reservation.Quantity(),
				CommittedAt:   time.Now(),
			}
			i.events = append(i.events, event)

			// Check for low stock alerts
			i.checkStockLevels()

			return nil
		}
	}

	return common.ErrReservationNotFound
}

// ReleaseReservation releases a reservation (returns reserved stock to available)
func (i *InventoryItem) ReleaseReservation(reservationID ReservationID, reason string) error {
	for idx, reservation := range i.reservations {
		if reservation.ID().Equals(reservationID) {
			if reservation.Status() != ReservationStatusActive {
				return common.ErrInvalidInput
			}

			err := reservation.Cancel()
			if err != nil {
				return err
			}

			// Update the reservation in slice
			i.reservations[idx] = reservation
			
			// Return stock to available
			i.stockLevel = NewStockLevel(i.stockLevel.Available()+reservation.Quantity(), i.stockLevel.Reserved()-reservation.Quantity())
			i.lastUpdated = time.Now()

			// Add domain event
			event := StockReleasedEvent{
				ProductID:     i.productID,
				ReservationID: reservationID,
				Quantity:      reservation.Quantity(),
				ReleasedAt:    time.Now(),
				Reason:        reason,
			}
			i.events = append(i.events, event)

			return nil
		}
	}

	return common.ErrReservationNotFound
}

// ReplenishStock adds stock to inventory
func (i *InventoryItem) ReplenishStock(quantity int, replenishedBy string) error {
	if quantity <= 0 {
		return common.ErrInvalidInput
	}

	// Add to available stock
	i.stockLevel = NewStockLevel(i.stockLevel.Available()+quantity, i.stockLevel.Reserved())
	i.lastUpdated = time.Now()

	// Add domain event
	event := StockReplenishedEvent{
		ProductID:      i.productID,
		QuantityAdded:  quantity,
		NewStockLevel:  i.stockLevel.Total(),
		ReplenishedAt:  time.Now(),
		ReplenishedBy:  replenishedBy,
	}
	i.events = append(i.events, event)

	return nil
}

// ProcessExpiredReservations processes and expires old reservations
func (i *InventoryItem) ProcessExpiredReservations() error {
	for idx, reservation := range i.reservations {
		if reservation.IsExpired() && reservation.Status() == ReservationStatusActive {
			err := reservation.Expire()
			if err != nil {
				continue
			}

			// Update the reservation in slice
			i.reservations[idx] = reservation
			
			// Return stock to available
			i.stockLevel = NewStockLevel(i.stockLevel.Available()+reservation.Quantity(), i.stockLevel.Reserved()-reservation.Quantity())

			// Add domain event
			event := StockReleasedEvent{
				ProductID:     i.productID,
				ReservationID: reservation.ID(),
				Quantity:      reservation.Quantity(),
				ReleasedAt:    time.Now(),
				Reason:        "expired",
			}
			i.events = append(i.events, event)
		}
	}

	i.lastUpdated = time.Now()
	return nil
}

// GetAlertLevel determines the current alert level based on stock
func (i *InventoryItem) GetAlertLevel() AlertLevel {
	available := i.stockLevel.Available()
	
	if available == 0 {
		return AlertLevelOutOfStock
	} else if available <= i.minStockLevel {
		return AlertLevelCritical
	} else if available <= i.reorderPoint {
		return AlertLevelLow
	}
	
	return AlertLevelNone
}

// IsLowStock checks if stock is below reorder point
func (i *InventoryItem) IsLowStock() bool {
	return i.stockLevel.Available() <= i.reorderPoint
}

// checkStockLevels checks stock levels and generates alerts if needed
func (i *InventoryItem) checkStockLevels() {
	alertLevel := i.GetAlertLevel()
	
	if alertLevel != AlertLevelNone {
		event := LowStockAlertEvent{
			ProductID:     i.productID,
			CurrentStock:  i.stockLevel.Available(),
			MinimumStock:  i.minStockLevel,
			AlertLevel:    alertLevel,
			AlertedAt:     time.Now(),
		}
		i.events = append(i.events, event)
	}
}