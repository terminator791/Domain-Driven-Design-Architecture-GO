package cart

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// DomainService represents complex business logic for cart operations
type DomainService struct {
	cartRepo    Repository
	productRepo ProductRepository
}

func NewDomainService(cartRepo Repository, productRepo ProductRepository) *DomainService {
	return &DomainService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetOrCreateActiveCart gets an active cart for a customer or creates a new one
func (s *DomainService) GetOrCreateActiveCart(ctx context.Context, customerID order.CustomerID) (*Cart, error) {
	// Try to find existing active cart
	existingCart, err := s.cartRepo.FindActiveByCustomerID(ctx, customerID)
	if err == nil && existingCart != nil {
		// Check if cart is expired
		if existingCart.IsExpired() {
			err := existingCart.Abandon()
			if err != nil {
				return nil, err
			}
			// Save the abandoned cart
			err = s.cartRepo.Save(ctx, existingCart)
			if err != nil {
				return nil, err
			}
			// Create new cart
			return NewCart(customerID), nil
		}
		return existingCart, nil
	}

	// Create new cart if none exists or error occurred
	return NewCart(customerID), nil
}

// ValidateAddItem validates business rules for adding an item to cart
func (s *DomainService) ValidateAddItem(ctx context.Context, cart *Cart, productID order.ProductID, quantity order.Quantity) error {
	// Validate that product exists
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return common.ErrProductNotFound
	}

	// Check if product is available
	if !product.IsAvailable() {
		return common.ErrProductNotAvailable
	}

	// Check stock availability
	if !product.HasSufficientStock(quantity) {
		return common.ErrInsufficientStock
	}

	// Business rule: Maximum 10 items per product in cart
	existingQuantity := 0
	for _, item := range cart.Items() {
		if item.ProductID().Equals(productID) {
			existingQuantity = item.Quantity().Value()
			break
		}
	}
	
	if existingQuantity+quantity.Value() > 10 {
		return common.ErrMaxQuantityExceeded
	}

	// Business rule: Maximum 50 total items in cart
	if cart.ItemCount()+quantity.Value() > 50 {
		return common.ErrCartTotalItemsExceeded
	}

	return nil
}

// ProcessAbandonedCarts processes expired carts and marks them as abandoned
func (s *DomainService) ProcessAbandonedCarts(ctx context.Context) error {
	expiredCarts, err := s.cartRepo.FindExpiredCarts(ctx)
	if err != nil {
		return err
	}

	for _, cart := range expiredCarts {
		if cart.Status() == CartStatusActive {
			err := cart.Abandon()
			if err != nil {
				continue // Log error in real implementation
			}
			
			err = s.cartRepo.Save(ctx, cart)
			if err != nil {
				continue // Log error in real implementation
			}
		}
	}

	return nil
}

// CalculateCartWithDiscounts applies business rules for cart-level discounts
func (s *DomainService) CalculateCartWithDiscounts(ctx context.Context, cart *Cart) (order.Money, error) {
	total := cart.Total()

	// Business rule: Free shipping for orders over $100
	if total.Amount() >= 100.0 {
		// This would typically be handled by a shipping service
		// For now, we'll just return the total as-is
	}

	// Business rule: 5% discount for carts with more than 5 different products
	if len(cart.Items()) > 5 {
		discountAmount := total.Amount() * 0.05
		discount, err := order.NewMoney(discountAmount, total.Currency())
		if err != nil {
			return order.Money{}, err
		}
		total, err = total.Add(order.Money{})
		if err != nil {
			return order.Money{}, err
		}
		// Subtract discount (add negative amount)
		discountedTotal, err := order.NewMoney(total.Amount()-discount.Amount(), total.Currency())
		if err != nil {
			return order.Money{}, err
		}
		total = discountedTotal
	}

	return total, nil
}

// ConvertCartToOrder converts a cart to an order (used during checkout)
func (s *DomainService) ConvertCartToOrder(ctx context.Context, cart *Cart) (*order.Order, error) {
	if cart.Status() != CartStatusActive {
		return nil, common.ErrInvalidInput
	}

	if len(cart.Items()) == 0 {
		return nil, common.ErrInvalidInput
	}

	// Create new order
	newOrder := order.NewOrder(cart.CustomerID())

	// Add items from cart to order
	for _, cartItem := range cart.Items() {
		err := newOrder.AddItem(cartItem.ProductID(), cartItem.Quantity(), cartItem.UnitPrice())
		if err != nil {
			return nil, err
		}
	}

	return newOrder, nil
}

// ProductRepository interface for products (used by cart domain service)
type ProductRepository interface {
	FindByID(ctx context.Context, id order.ProductID) (ProductInfo, error)
}

// ProductInfo represents product information needed by the cart domain service
type ProductInfo interface {
	ID() order.ProductID
	Name() string
	Price() order.Money
	StockLevel() int
	IsAvailable() bool
	HasSufficientStock(quantity order.Quantity) bool
}