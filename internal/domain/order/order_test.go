package order_test

import (
	"testing"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

func TestOrderCreation(t *testing.T) {
	// Create customer ID
	customerID := order.NewCustomerID()

	// Create new order
	o := order.NewOrder(customerID)

	// Verify order was created correctly
	if o.ID().String() == "" {
		t.Error("Order ID should not be empty")
	}

	if !o.CustomerID().Equals(customerID) {
		t.Error("Customer ID should match")
	}

	if o.Status() != order.OrderStatusPending {
		t.Error("New order should have PENDING status")
	}

	if len(o.Items()) != 0 {
		t.Error("New order should have no items")
	}
}

func TestAddItemToOrder(t *testing.T) {
	// Create order
	customerID := order.NewCustomerID()
	o := order.NewOrder(customerID)

	// Create item details
	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(2)
	unitPrice, _ := order.NewMoney(10.50, "USD")

	// Add item to order
	err := o.AddItem(productID, quantity, unitPrice)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	// Verify item was added
	if len(o.Items()) != 1 {
		t.Error("Order should have 1 item")
	}

	item := o.Items()[0]
	if !item.ProductID().Equals(productID) {
		t.Error("Product ID should match")
	}

	if !item.Quantity().Equals(quantity) {
		t.Error("Quantity should match")
	}

	if !item.UnitPrice().Equals(unitPrice) {
		t.Error("Unit price should match")
	}

	// Verify total was calculated
	expectedTotal, _ := order.NewMoney(21.00, "USD") // 2 * 10.50
	if !o.Total().Equals(expectedTotal) {
		t.Errorf("Expected total %s, got %s", expectedTotal.String(), o.Total().String())
	}
}

func TestOrderStatusTransition(t *testing.T) {
	// Create order with items
	customerID := order.NewCustomerID()
	o := order.NewOrder(customerID)

	productID := order.NewProductID()
	quantity, _ := order.NewQuantity(1)
	unitPrice, _ := order.NewMoney(100.00, "USD")

	o.AddItem(productID, quantity, unitPrice)

	// Confirm order
	err := o.Confirm()
	if err != nil {
		t.Fatalf("Failed to confirm order: %v", err)
	}

	if o.Status() != order.OrderStatusConfirmed {
		t.Error("Order should be CONFIRMED")
	}

	// Test valid transition
	err = o.UpdateStatus(order.OrderStatusShipped)
	if err != nil {
		t.Fatalf("Failed to update status to SHIPPED: %v", err)
	}

	// Test invalid transition
	err = o.UpdateStatus(order.OrderStatusPending)
	if err == nil {
		t.Error("Should not allow transition from SHIPPED to PENDING")
	}
}

func TestMoneyValueObject(t *testing.T) {
	// Test valid money creation
	money, err := order.NewMoney(10.50, "USD")
	if err != nil {
		t.Fatalf("Failed to create money: %v", err)
	}

	if money.Amount() != 10.50 {
		t.Error("Amount should be 10.50")
	}

	if money.Currency() != "USD" {
		t.Error("Currency should be USD")
	}

	// Test invalid money creation
	_, err = order.NewMoney(-10.50, "USD")
	if err == nil {
		t.Error("Should not allow negative amounts")
	}

	_, err = order.NewMoney(10.50, "")
	if err == nil {
		t.Error("Should not allow empty currency")
	}

	// Test money addition
	money2, _ := order.NewMoney(5.25, "USD")
	total, err := money.Add(money2)
	if err != nil {
		t.Fatalf("Failed to add money: %v", err)
	}

	if total.Amount() != 15.75 {
		t.Error("Total should be 15.75")
	}

	// Test money multiplication
	doubled := money.Multiply(2)
	if doubled.Amount() != 21.00 {
		t.Error("Doubled amount should be 21.00")
	}
}