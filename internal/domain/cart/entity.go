package cart

import (
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

// CartItem represents an item within a shopping cart (Entity)
type CartItem struct {
	productID order.ProductID
	quantity  order.Quantity
	unitPrice order.Money
	addedAt   time.Time
}

func NewCartItem(productID order.ProductID, quantity order.Quantity, unitPrice order.Money) CartItem {
	return CartItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
		addedAt:   time.Now(),
	}
}

func (item CartItem) ProductID() order.ProductID {
	return item.productID
}

func (item CartItem) Quantity() order.Quantity {
	return item.quantity
}

func (item CartItem) UnitPrice() order.Money {
	return item.unitPrice
}

func (item CartItem) AddedAt() time.Time {
	return item.addedAt
}

func (item CartItem) TotalPrice() order.Money {
	return item.unitPrice.Multiply(item.quantity.Value())
}

func (item *CartItem) UpdateQuantity(quantity order.Quantity) {
	item.quantity = quantity
}

// Cart represents a shopping cart aggregate root (Aggregate Root Entity)
type Cart struct {
	id         CartID
	customerID order.CustomerID
	items      []CartItem
	status     CartStatus
	total      order.Money
	createdAt  time.Time
	updatedAt  time.Time
	expiresAt  time.Time
	events     []interface{} // Domain events
}

func NewCart(customerID order.CustomerID) *Cart {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Cart expires in 24 hours
	
	return &Cart{
		id:         NewCartID(),
		customerID: customerID,
		items:      make([]CartItem, 0),
		status:     CartStatusActive,
		total:      order.Money{}, // Will be calculated when items are added
		createdAt:  now,
		updatedAt:  now,
		expiresAt:  expiresAt,
		events:     make([]interface{}, 0),
	}
}

func (c *Cart) ID() CartID {
	return c.id
}

func (c *Cart) CustomerID() order.CustomerID {
	return c.customerID
}

func (c *Cart) Items() []CartItem {
	return c.items
}

func (c *Cart) Status() CartStatus {
	return c.status
}

func (c *Cart) Total() order.Money {
	return c.total
}

func (c *Cart) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Cart) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Cart) ExpiresAt() time.Time {
	return c.expiresAt
}

func (c *Cart) Events() []interface{} {
	return c.events
}

func (c *Cart) ClearEvents() {
	c.events = make([]interface{}, 0)
}

// AddItem adds an item to the cart or updates quantity if item already exists
func (c *Cart) AddItem(productID order.ProductID, quantity order.Quantity, unitPrice order.Money) error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	// Check if cart is expired
	if time.Now().After(c.expiresAt) {
		c.status = CartStatusExpired
		return common.ErrInvalidInput
	}

	// Check if item already exists, if so, update quantity
	for i, item := range c.items {
		if item.productID.Equals(productID) {
			newQuantity := item.quantity.Add(quantity)
			c.items[i].UpdateQuantity(newQuantity)
			c.updatedAt = time.Now()
			
			// Add domain event
			event := CartItemAddedEvent{
				CartID:    c.id,
				ProductID: productID,
				Quantity:  quantity,
				UnitPrice: unitPrice,
				AddedAt:   time.Now(),
			}
			c.events = append(c.events, event)
			
			return c.recalculateTotal()
		}
	}

	// Add new item
	item := NewCartItem(productID, quantity, unitPrice)
	c.items = append(c.items, item)
	c.updatedAt = time.Now()

	// Add domain event
	event := CartItemAddedEvent{
		CartID:    c.id,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		AddedAt:   time.Now(),
	}
	c.events = append(c.events, event)

	return c.recalculateTotal()
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID order.ProductID) error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	for i, item := range c.items {
		if item.productID.Equals(productID) {
			// Remove item from slice
			c.items = append(c.items[:i], c.items[i+1:]...)
			c.updatedAt = time.Now()

			// Add domain event
			event := CartItemRemovedEvent{
				CartID:    c.id,
				ProductID: productID,
				RemovedAt: time.Now(),
			}
			c.events = append(c.events, event)

			return c.recalculateTotal()
		}
	}

	return common.ErrProductNotFound
}

// UpdateItemQuantity updates the quantity of a specific item
func (c *Cart) UpdateItemQuantity(productID order.ProductID, quantity order.Quantity) error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	for i, item := range c.items {
		if item.productID.Equals(productID) {
			c.items[i].UpdateQuantity(quantity)
			c.updatedAt = time.Now()
			
			return c.recalculateTotal()
		}
	}

	return common.ErrProductNotFound
}

// Clear removes all items from the cart
func (c *Cart) Clear() error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	c.items = make([]CartItem, 0)
	c.updatedAt = time.Now()
	return c.recalculateTotal()
}

// CheckOut transitions the cart to checked out status
func (c *Cart) CheckOut() error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	if len(c.items) == 0 {
		return common.ErrInvalidInput
	}

	c.status = CartStatusCheckedOut
	c.updatedAt = time.Now()

	// Add domain event
	event := CartCheckedOutEvent{
		CartID:       c.id,
		CustomerID:   c.customerID,
		Total:        c.total,
		CheckedOutAt: time.Now(),
	}
	c.events = append(c.events, event)

	return nil
}

// Abandon marks the cart as abandoned
func (c *Cart) Abandon() error {
	if c.status != CartStatusActive {
		return common.ErrInvalidInput
	}

	c.status = CartStatusAbandoned
	c.updatedAt = time.Now()

	// Add domain event
	event := CartAbandonedEvent{
		CartID:      c.id,
		CustomerID:  c.customerID,
		AbandonedAt: time.Now(),
	}
	c.events = append(c.events, event)

	return nil
}

// ExtendExpiry extends the cart expiry time
func (c *Cart) ExtendExpiry(duration time.Duration) {
	c.expiresAt = c.expiresAt.Add(duration)
	c.updatedAt = time.Now()
}

// IsExpired checks if the cart is expired
func (c *Cart) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

// ItemCount returns the total number of items in the cart
func (c *Cart) ItemCount() int {
	totalItems := 0
	for _, item := range c.items {
		totalItems += item.quantity.Value()
	}
	return totalItems
}

// recalculateTotal recalculates the total amount of the cart
func (c *Cart) recalculateTotal() error {
	if len(c.items) == 0 {
		c.total = order.Money{}
		return nil
	}

	firstItem := c.items[0]
	total := firstItem.TotalPrice()

	for i := 1; i < len(c.items); i++ {
		itemTotal := c.items[i].TotalPrice()
		var err error
		total, err = total.Add(itemTotal)
		if err != nil {
			return err
		}
	}

	c.total = total
	return nil
}