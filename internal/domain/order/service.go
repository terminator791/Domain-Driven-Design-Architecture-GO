package order

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
)

// DomainService represents complex business logic that doesn't belong to a single entity
type DomainService struct {
	orderRepo   Repository
	productRepo ProductRepository
}

func NewDomainService(orderRepo Repository, productRepo ProductRepository) *DomainService {
	return &DomainService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// ValidateOrderCreation validates business rules for order creation
func (s *DomainService) ValidateOrderCreation(ctx context.Context, order *Order) error {
	// Validate that all products exist and have sufficient stock
	for _, item := range order.Items() {
		product, err := s.productRepo.FindByID(ctx, item.ProductID())
		if err != nil {
			return common.ErrProductNotFound
		}

		// Check stock availability (business rule)
		if !product.HasSufficientStock(item.Quantity()) {
			return common.ErrInsufficientStock
		}

		// Validate that the unit price matches the product price
		if !item.UnitPrice().Equals(product.Price()) {
			return common.ErrInvalidPrice
		}
	}

	return nil
}

// CalculateOrderTotal calculates the total for an order with potential discounts
func (s *DomainService) CalculateOrderTotal(ctx context.Context, order *Order) (Money, error) {
	total := order.Total()

	// Apply business rules for discounts (example)
	// If order has more than 5 items, apply 5% discount
	if len(order.Items()) > 5 {
		discountAmount := total.Amount() * 0.05
		discount, err := NewMoney(discountAmount, total.Currency())
		if err != nil {
			return Money{}, err
		}
		total, err = total.Add(Money{amount: -discount.Amount(), currency: discount.Currency()})
		if err != nil {
			return Money{}, err
		}
	}

	return total, nil
}

// ProductRepository interface for products (used by domain service)
type ProductRepository interface {
	FindByID(ctx context.Context, id ProductID) (ProductInfo, error)
}

// ProductInfo represents product information needed by the domain service
type ProductInfo interface {
	ID() ProductID
	Name() string
	Price() Money
	StockLevel() int
	IsAvailable() bool
	HasSufficientStock(quantity Quantity) bool
}