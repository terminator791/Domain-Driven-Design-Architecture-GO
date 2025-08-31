package cart

import (
	"testing"
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

func TestCartCreation(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	if cart.ID().String() == "" {
		t.Error("Cart ID should not be empty")
	}

	if !cart.CustomerID().Equals(customerID) {
		t.Error("Cart customer ID should match provided customer ID")
	}

	if cart.Status() != CartStatusActive {
		t.Error("New cart should have ACTIVE status")
	}

	if len(cart.Items()) != 0 {
		t.Error("New cart should have no items")
	}

	if cart.ItemCount() != 0 {
		t.Error("New cart should have zero item count")
	}

	if cart.IsExpired() {
		t.Error("New cart should not be expired")
	}
}

func TestAddItemToCart(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(2)
	unitPrice, _ := order.NewMoney(10.99, "USD")

	err := cart.AddItem(productID, quantity, unitPrice)
	if err != nil {
		t.Errorf("Adding item to cart should not error: %v", err)
	}

	if len(cart.Items()) != 1 {
		t.Error("Cart should have 1 item after adding")
	}

	if cart.ItemCount() != 2 {
		t.Error("Cart should have 2 total items")
	}

	// Test adding same product again (should update quantity)
	err = cart.AddItem(productID, quantity, unitPrice)
	if err != nil {
		t.Errorf("Adding same item again should not error: %v", err)
	}

	if len(cart.Items()) != 1 {
		t.Error("Cart should still have 1 unique item")
	}

	if cart.ItemCount() != 4 {
		t.Error("Cart should have 4 total items after adding same product again")
	}

	// Verify total is calculated correctly
	expectedTotal, _ := order.NewMoney(43.96, "USD") // 4 * 10.99
	if !cart.Total().Equals(expectedTotal) {
		t.Errorf("Cart total should be %v, got %v", expectedTotal, cart.Total())
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(2)
	unitPrice, _ := order.NewMoney(10.99, "USD")

	// Add item first
	cart.AddItem(productID, quantity, unitPrice)

	// Remove item
	err := cart.RemoveItem(productID)
	if err != nil {
		t.Errorf("Removing item should not error: %v", err)
	}

	if len(cart.Items()) != 0 {
		t.Error("Cart should have no items after removal")
	}

	if cart.ItemCount() != 0 {
		t.Error("Cart should have zero total items after removal")
	}
}

func TestCartCheckout(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(1)
	unitPrice, _ := order.NewMoney(25.50, "USD")

	// Add item to cart
	cart.AddItem(productID, quantity, unitPrice)

	// Test checkout
	err := cart.CheckOut()
	if err != nil {
		t.Errorf("Checkout should not error: %v", err)
	}

	if cart.Status() != CartStatusCheckedOut {
		t.Error("Cart status should be CHECKED_OUT after checkout")
	}

	// Test that we can't add items after checkout
	err = cart.AddItem(productID, quantity, unitPrice)
	if err == nil {
		t.Error("Adding items to checked out cart should error")
	}
}

func TestCartExpiry(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	// Set expiry to past time to simulate expired cart
	pastTime := time.Now().Add(-1 * time.Hour)
	cart.expiresAt = pastTime

	if !cart.IsExpired() {
		t.Error("Cart with past expiry should be expired")
	}

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(1)
	unitPrice, _ := order.NewMoney(10.00, "USD")

	// Try to add item to expired cart
	err := cart.AddItem(productID, quantity, unitPrice)
	if err == nil {
		t.Error("Adding item to expired cart should error")
	}

	if cart.Status() != CartStatusExpired {
		t.Error("Cart status should be EXPIRED after trying to add to expired cart")
	}
}

func TestCartAbandonment(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(1)
	unitPrice, _ := order.NewMoney(15.00, "USD")

	// Add item to cart
	cart.AddItem(productID, quantity, unitPrice)

	// Abandon cart
	err := cart.Abandon()
	if err != nil {
		t.Errorf("Abandoning cart should not error: %v", err)
	}

	if cart.Status() != CartStatusAbandoned {
		t.Error("Cart status should be ABANDONED after abandonment")
	}

	// Test that we can't add items after abandonment
	err = cart.AddItem(productID, quantity, unitPrice)
	if err == nil {
		t.Error("Adding items to abandoned cart should error")
	}
}

func TestCartClear(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	// Add multiple items
	productID1 := order.NewProductID()
	productID2 := order.NewProductID()
	quantity, _ := order.NewQuantity(2)
	unitPrice, _ := order.NewMoney(10.00, "USD")

	cart.AddItem(productID1, quantity, unitPrice)
	cart.AddItem(productID2, quantity, unitPrice)

	// Clear cart
	err := cart.Clear()
	if err != nil {
		t.Errorf("Clearing cart should not error: %v", err)
	}

	if len(cart.Items()) != 0 {
		t.Error("Cart should have no items after clearing")
	}

	if cart.ItemCount() != 0 {
		t.Error("Cart should have zero total items after clearing")
	}

	expectedTotal, _ := order.NewMoney(0, "")
	if cart.Total().Amount() != expectedTotal.Amount() {
		t.Error("Cart total should be zero after clearing")
	}
}

func TestCartExtendExpiry(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	originalExpiry := cart.ExpiresAt()
	
	// Extend expiry by 1 hour
	cart.ExtendExpiry(1 * time.Hour)
	
	if !cart.ExpiresAt().After(originalExpiry) {
		t.Error("Cart expiry should be extended")
	}
}

func TestCartDomainEvents(t *testing.T) {
	customerID := order.NewCustomerID()
	cart := NewCart(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(1)
	unitPrice, _ := order.NewMoney(10.00, "USD")

	// Clear any initial events
	cart.ClearEvents()

	// Add item (should generate event)
	cart.AddItem(productID, quantity, unitPrice)
	
	if len(cart.Events()) != 1 {
		t.Error("Adding item should generate 1 domain event")
	}

	// Remove item (should generate another event)
	cart.RemoveItem(productID)
	
	if len(cart.Events()) != 2 {
		t.Error("Cart should have 2 domain events after add and remove")
	}

	// Clear events
	cart.ClearEvents()
	
	if len(cart.Events()) != 0 {
		t.Error("Cart should have no events after clearing")
	}
}