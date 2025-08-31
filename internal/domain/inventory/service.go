package inventory

import (
	"context"
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// DomainService represents complex business logic for inventory operations
type DomainService struct {
	inventoryRepo    Repository
	reservationRepo  ReservationRepository
	eventPublisher   EventPublisher
}

func NewDomainService(inventoryRepo Repository, reservationRepo ReservationRepository, eventPublisher EventPublisher) *DomainService {
	return &DomainService{
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
		eventPublisher:  eventPublisher,
	}
}

// ReserveStockForOrder reserves stock for an order with business rules
func (s *DomainService) ReserveStockForOrder(ctx context.Context, orderItems []OrderItem, customerID order.CustomerID) ([]*StockReservation, error) {
	var reservations []*StockReservation

	// Attempt to reserve stock for each order item
	for _, item := range orderItems {
		inventoryItem, err := s.inventoryRepo.FindByProductID(ctx, item.ProductID)
		if err != nil {
			// Rollback any successful reservations
			s.rollbackReservations(ctx, reservations)
			return nil, common.ErrProductNotFound
		}

		// Business rule: Standard reservation duration is 30 minutes
		reservation, err := inventoryItem.ReserveStock(customerID, item.Quantity, 30*time.Minute)
		if err != nil {
			// Rollback any successful reservations
			s.rollbackReservations(ctx, reservations)
			return nil, err
		}

		// Save inventory item with updated stock levels
		err = s.inventoryRepo.Save(ctx, inventoryItem)
		if err != nil {
			// Rollback any successful reservations
			s.rollbackReservations(ctx, reservations)
			return nil, err
		}

		// Save the reservation
		err = s.reservationRepo.Save(ctx, reservation)
		if err != nil {
			// Rollback any successful reservations
			s.rollbackReservations(ctx, reservations)
			return nil, err
		}

		reservations = append(reservations, reservation)

		// Publish domain events
		for _, event := range inventoryItem.Events() {
			s.eventPublisher.Publish(ctx, event)
		}
		inventoryItem.ClearEvents()
	}

	return reservations, nil
}

// CommitReservationsForOrder commits all reservations for an order
func (s *DomainService) CommitReservationsForOrder(ctx context.Context, reservationIDs []ReservationID) error {
	for _, reservationID := range reservationIDs {
		reservation, err := s.reservationRepo.FindByID(ctx, reservationID)
		if err != nil {
			return err
		}

		inventoryItem, err := s.inventoryRepo.FindByProductID(ctx, reservation.ProductID())
		if err != nil {
			return err
		}

		err = inventoryItem.CommitReservation(reservationID)
		if err != nil {
			return err
		}

		// Save updated inventory item
		err = s.inventoryRepo.Save(ctx, inventoryItem)
		if err != nil {
			return err
		}

		// Save updated reservation
		err = s.reservationRepo.Save(ctx, reservation)
		if err != nil {
			return err
		}

		// Publish domain events
		for _, event := range inventoryItem.Events() {
			s.eventPublisher.Publish(ctx, event)
		}
		inventoryItem.ClearEvents()
	}

	return nil
}

// ReleaseReservationsForOrder releases all reservations for an order
func (s *DomainService) ReleaseReservationsForOrder(ctx context.Context, reservationIDs []ReservationID, reason string) error {
	for _, reservationID := range reservationIDs {
		reservation, err := s.reservationRepo.FindByID(ctx, reservationID)
		if err != nil {
			continue // Log error in real implementation
		}

		inventoryItem, err := s.inventoryRepo.FindByProductID(ctx, reservation.ProductID())
		if err != nil {
			continue // Log error in real implementation
		}

		err = inventoryItem.ReleaseReservation(reservationID, reason)
		if err != nil {
			continue // Log error in real implementation
		}

		// Save updated inventory item
		err = s.inventoryRepo.Save(ctx, inventoryItem)
		if err != nil {
			continue // Log error in real implementation
		}

		// Save updated reservation
		err = s.reservationRepo.Save(ctx, reservation)
		if err != nil {
			continue // Log error in real implementation
		}

		// Publish domain events
		for _, event := range inventoryItem.Events() {
			s.eventPublisher.Publish(ctx, event)
		}
		inventoryItem.ClearEvents()
	}

	return nil
}

// ProcessExpiredReservations processes all expired reservations
func (s *DomainService) ProcessExpiredReservations(ctx context.Context) error {
	itemsWithExpiredReservations, err := s.inventoryRepo.FindItemsWithExpiredReservations(ctx)
	if err != nil {
		return err
	}

	for _, item := range itemsWithExpiredReservations {
		err := item.ProcessExpiredReservations()
		if err != nil {
			continue // Log error in real implementation
		}

		// Save updated inventory item
		err = s.inventoryRepo.Save(ctx, item)
		if err != nil {
			continue // Log error in real implementation
		}

		// Publish domain events
		for _, event := range item.Events() {
			s.eventPublisher.Publish(ctx, event)
		}
		item.ClearEvents()
	}

	return nil
}

// CheckLowStockLevels checks for low stock and generates alerts
func (s *DomainService) CheckLowStockLevels(ctx context.Context) error {
	lowStockItems, err := s.inventoryRepo.FindLowStockItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range lowStockItems {
		// Generate low stock alert event
		alertLevel := item.GetAlertLevel()
		if alertLevel != AlertLevelNone {
			event := LowStockAlertEvent{
				ProductID:     item.ProductID(),
				CurrentStock:  item.StockLevel().Available(),
				MinimumStock:  item.MinStockLevel(),
				AlertLevel:    alertLevel,
				AlertedAt:     time.Now(),
			}
			
			err := s.eventPublisher.Publish(ctx, event)
			if err != nil {
				continue // Log error in real implementation
			}
		}
	}

	return nil
}

// ValidateStockAvailability validates if sufficient stock is available for an order
func (s *DomainService) ValidateStockAvailability(ctx context.Context, orderItems []OrderItem) error {
	for _, item := range orderItems {
		inventoryItem, err := s.inventoryRepo.FindByProductID(ctx, item.ProductID)
		if err != nil {
			return common.ErrProductNotFound
		}

		if !inventoryItem.StockLevel().CanReserve(item.Quantity) {
			return common.ErrInsufficientStock
		}
	}

	return nil
}

// GetStockLevel gets the current stock level for a product
func (s *DomainService) GetStockLevel(ctx context.Context, productID order.ProductID) (StockLevel, error) {
	inventoryItem, err := s.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		return StockLevel{}, err
	}

	return inventoryItem.StockLevel(), nil
}

// rollbackReservations is a helper method to rollback reservations in case of failure
func (s *DomainService) rollbackReservations(ctx context.Context, reservations []*StockReservation) {
	var reservationIDs []ReservationID
	for _, reservation := range reservations {
		reservationIDs = append(reservationIDs, reservation.ID())
	}
	
	// Best effort rollback - ignore errors
	s.ReleaseReservationsForOrder(ctx, reservationIDs, "rollback")
}

// OrderItem represents an item from an order for inventory operations
type OrderItem struct {
	ProductID order.ProductID
	Quantity  int
}