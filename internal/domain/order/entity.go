package order

import (
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
)

// OrderItem represents an item within an order (Entity)
type OrderItem struct {
	productID ProductID
	quantity  Quantity
	unitPrice Money
}

func NewOrderItem(productID ProductID, quantity Quantity, unitPrice Money) OrderItem {
	return OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
	}
}

func (item OrderItem) ProductID() ProductID {
	return item.productID
}

func (item OrderItem) Quantity() Quantity {
	return item.quantity
}

func (item OrderItem) UnitPrice() Money {
	return item.unitPrice
}

func (item OrderItem) TotalPrice() Money {
	return item.unitPrice.Multiply(item.quantity.Value())
}

// Order represents an order aggregate root (Aggregate Root Entity)
type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     OrderStatus
	total      Money
	createdAt  time.Time
	updatedAt  time.Time
	events     []interface{} // Domain events
}

func NewOrder(customerID CustomerID) *Order {
	now := time.Now()
	return &Order{
		id:         NewOrderID(),
		customerID: customerID,
		items:      make([]OrderItem, 0),
		status:     OrderStatusPending,
		total:      Money{}, // Will be calculated when items are added
		createdAt:  now,
		updatedAt:  now,
		events:     make([]interface{}, 0),
	}
}

func (o *Order) ID() OrderID {
	return o.id
}

func (o *Order) CustomerID() CustomerID {
	return o.customerID
}

func (o *Order) Items() []OrderItem {
	return o.items
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Total() Money {
	return o.total
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Order) Events() []interface{} {
	return o.events
}

func (o *Order) ClearEvents() {
	o.events = make([]interface{}, 0)
}

// AddItem adds an item to the order and recalculates total
func (o *Order) AddItem(productID ProductID, quantity Quantity, unitPrice Money) error {
	if o.status != OrderStatusPending {
		return common.ErrInvalidInput
	}

	// Check if item already exists, if so, update quantity
	for i, item := range o.items {
		if item.productID.Equals(productID) {
			newQuantity := item.quantity.Add(quantity)
			o.items[i] = NewOrderItem(productID, newQuantity, unitPrice)
			return o.recalculateTotal()
		}
	}

	// Add new item
	item := NewOrderItem(productID, quantity, unitPrice)
	o.items = append(o.items, item)
	o.updatedAt = time.Now()

	return o.recalculateTotal()
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productID ProductID) error {
	if o.status != OrderStatusPending {
		return common.ErrInvalidInput
	}

	for i, item := range o.items {
		if item.productID.Equals(productID) {
			// Remove item by replacing with last element and truncating
			o.items[i] = o.items[len(o.items)-1]
			o.items = o.items[:len(o.items)-1]
			o.updatedAt = time.Now()
			return o.recalculateTotal()
		}
	}

	return common.ErrEntityNotFound
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if !status.IsValid() {
		return common.ErrInvalidInput
	}

	// Validate status transition
	if !o.canTransitionTo(status) {
		return common.ErrInvalidInput
	}

	o.status = status
	o.updatedAt = time.Now()

	// Add domain event
	event := OrderUpdatedEvent{
		OrderID:   o.id,
		Status:    status,
		UpdatedAt: o.updatedAt,
	}
	o.events = append(o.events, event)

	return nil
}

// Confirm confirms the order and creates domain event
func (o *Order) Confirm() error {
	if len(o.items) == 0 {
		return common.ErrInvalidInput
	}

	err := o.UpdateStatus(OrderStatusConfirmed)
	if err != nil {
		return err
	}

	// Add order created event
	event := OrderCreatedEvent{
		OrderID:   o.id,
		Total:     o.total,
		CreatedAt: o.createdAt,
	}
	o.events = append(o.events, event)

	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	return o.UpdateStatus(OrderStatusCancelled)
}

// recalculateTotal recalculates the total amount of the order
func (o *Order) recalculateTotal() error {
	if len(o.items) == 0 {
		o.total = Money{}
		return nil
	}

	firstItem := o.items[0]
	total := firstItem.TotalPrice()

	for i := 1; i < len(o.items); i++ {
		itemTotal := o.items[i].TotalPrice()
		var err error
		total, err = total.Add(itemTotal)
		if err != nil {
			return err
		}
	}

	o.total = total
	return nil
}

// canTransitionTo checks if the order can transition to the given status
func (o *Order) canTransitionTo(newStatus OrderStatus) bool {
	switch o.status {
	case OrderStatusPending:
		return newStatus == OrderStatusConfirmed || newStatus == OrderStatusCancelled
	case OrderStatusConfirmed:
		return newStatus == OrderStatusShipped || newStatus == OrderStatusCancelled
	case OrderStatusShipped:
		return newStatus == OrderStatusDelivered
	case OrderStatusDelivered:
		return false // Final state
	case OrderStatusCancelled:
		return false // Final state
	default:
		return false
	}
}