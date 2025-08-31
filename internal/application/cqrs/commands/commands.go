package commands

import (
	"github.com/google/uuid"
)

// CreateOrderCommand represents a command to create an order
type CreateOrderCommand struct {
	commandID   string
	customerID  string
	items       []CreateOrderItemCommand
}

type CreateOrderItemCommand struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Currency  string  `json:"currency"`
}

func NewCreateOrderCommand(customerID string, items []CreateOrderItemCommand) *CreateOrderCommand {
	return &CreateOrderCommand{
		commandID:  uuid.New().String(),
		customerID: customerID,
		items:      items,
	}
}

func (c *CreateOrderCommand) CommandID() string {
	return c.commandID
}

func (c *CreateOrderCommand) CommandType() string {
	return "CreateOrder"
}

func (c *CreateOrderCommand) AggregateID() string {
	return c.customerID
}

func (c *CreateOrderCommand) GetPayload() interface{} {
	return struct {
		CustomerID string                     `json:"customer_id"`
		Items      []CreateOrderItemCommand   `json:"items"`
	}{
		CustomerID: c.customerID,
		Items:      c.items,
	}
}

func (c *CreateOrderCommand) CustomerID() string {
	return c.customerID
}

func (c *CreateOrderCommand) Items() []CreateOrderItemCommand {
	return c.items
}

// AddToCartCommand represents a command to add item to cart
type AddToCartCommand struct {
	commandID  string
	customerID string
	productID  string
	quantity   int
	unitPrice  float64
	currency   string
}

func NewAddToCartCommand(customerID, productID string, quantity int, unitPrice float64, currency string) *AddToCartCommand {
	return &AddToCartCommand{
		commandID:  uuid.New().String(),
		customerID: customerID,
		productID:  productID,
		quantity:   quantity,
		unitPrice:  unitPrice,
		currency:   currency,
	}
}

func (c *AddToCartCommand) CommandID() string {
	return c.commandID
}

func (c *AddToCartCommand) CommandType() string {
	return "AddToCart"
}

func (c *AddToCartCommand) AggregateID() string {
	return c.customerID
}

func (c *AddToCartCommand) GetPayload() interface{} {
	return struct {
		CustomerID string  `json:"customer_id"`
		ProductID  string  `json:"product_id"`
		Quantity   int     `json:"quantity"`
		UnitPrice  float64 `json:"unit_price"`
		Currency   string  `json:"currency"`
	}{
		CustomerID: c.customerID,
		ProductID:  c.productID,
		Quantity:   c.quantity,
		UnitPrice:  c.unitPrice,
		Currency:   c.currency,
	}
}

func (c *AddToCartCommand) CustomerID() string {
	return c.customerID
}

func (c *AddToCartCommand) ProductID() string {
	return c.productID
}

func (c *AddToCartCommand) Quantity() int {
	return c.quantity
}

func (c *AddToCartCommand) UnitPrice() float64 {
	return c.unitPrice
}

func (c *AddToCartCommand) Currency() string {
	return c.currency
}

// ProcessPaymentCommand represents a command to process payment
type ProcessPaymentCommand struct {
	commandID       string
	orderID         string
	customerID      string
	paymentMethodID string
	amount          float64
	currency        string
}

func NewProcessPaymentCommand(orderID, customerID, paymentMethodID string, amount float64, currency string) *ProcessPaymentCommand {
	return &ProcessPaymentCommand{
		commandID:       uuid.New().String(),
		orderID:         orderID,
		customerID:      customerID,
		paymentMethodID: paymentMethodID,
		amount:          amount,
		currency:        currency,
	}
}

func (c *ProcessPaymentCommand) CommandID() string {
	return c.commandID
}

func (c *ProcessPaymentCommand) CommandType() string {
	return "ProcessPayment"
}

func (c *ProcessPaymentCommand) AggregateID() string {
	return c.orderID
}

func (c *ProcessPaymentCommand) GetPayload() interface{} {
	return struct {
		OrderID         string  `json:"order_id"`
		CustomerID      string  `json:"customer_id"`
		PaymentMethodID string  `json:"payment_method_id"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
	}{
		OrderID:         c.orderID,
		CustomerID:      c.customerID,
		PaymentMethodID: c.paymentMethodID,
		Amount:          c.amount,
		Currency:        c.currency,
	}
}

func (c *ProcessPaymentCommand) OrderID() string {
	return c.orderID
}

func (c *ProcessPaymentCommand) CustomerID() string {
	return c.customerID
}

func (c *ProcessPaymentCommand) PaymentMethodID() string {
	return c.paymentMethodID
}

func (c *ProcessPaymentCommand) Amount() float64 {
	return c.amount
}

func (c *ProcessPaymentCommand) Currency() string {
	return c.currency
}

// UpdateInventoryCommand represents a command to update inventory
type UpdateInventoryCommand struct {
	commandID    string
	productID    string
	quantity     int
	operation    string // "ADD" or "REMOVE"
	reason       string
	updatedBy    string
}

func NewUpdateInventoryCommand(productID string, quantity int, operation, reason, updatedBy string) *UpdateInventoryCommand {
	return &UpdateInventoryCommand{
		commandID: uuid.New().String(),
		productID: productID,
		quantity:  quantity,
		operation: operation,
		reason:    reason,
		updatedBy: updatedBy,
	}
}

func (c *UpdateInventoryCommand) CommandID() string {
	return c.commandID
}

func (c *UpdateInventoryCommand) CommandType() string {
	return "UpdateInventory"
}

func (c *UpdateInventoryCommand) AggregateID() string {
	return c.productID
}

func (c *UpdateInventoryCommand) GetPayload() interface{} {
	return struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Operation string `json:"operation"`
		Reason    string `json:"reason"`
		UpdatedBy string `json:"updated_by"`
	}{
		ProductID: c.productID,
		Quantity:  c.quantity,
		Operation: c.operation,
		Reason:    c.reason,
		UpdatedBy: c.updatedBy,
	}
}

func (c *UpdateInventoryCommand) ProductID() string {
	return c.productID
}

func (c *UpdateInventoryCommand) Quantity() int {
	return c.quantity
}

func (c *UpdateInventoryCommand) Operation() string {
	return c.operation
}

func (c *UpdateInventoryCommand) Reason() string {
	return c.reason
}

func (c *UpdateInventoryCommand) UpdatedBy() string {
	return c.updatedBy
}