package order

// CreateOrderRequest represents the command to create a new order
type CreateOrderRequest struct {
	CustomerID string                      `json:"customer_id" binding:"required"`
	Items      []CreateOrderItemRequest    `json:"items" binding:"required,min=1"`
}

// CreateOrderItemRequest represents an item in the create order request
type CreateOrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
	Currency  string  `json:"currency" binding:"required"`
}

// UpdateOrderStatusRequest represents the command to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// AddOrderItemRequest represents the command to add an item to an order
type AddOrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
	Currency  string  `json:"currency" binding:"required"`
}

// RemoveOrderItemRequest represents the command to remove an item from an order
type RemoveOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}