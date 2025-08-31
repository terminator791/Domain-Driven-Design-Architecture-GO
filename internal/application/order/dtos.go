package order

import "time"

// OrderResponse represents the response data for an order
type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Items      []OrderItemResponse `json:"items"`
	Status     string              `json:"status"`
	Total      MoneyResponse       `json:"total"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// OrderItemResponse represents an order item in the response
type OrderItemResponse struct {
	ProductID string        `json:"product_id"`
	Quantity  int           `json:"quantity"`
	UnitPrice MoneyResponse `json:"unit_price"`
	Total     MoneyResponse `json:"total"`
}

// MoneyResponse represents money in the response
type MoneyResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// OrderListResponse represents a list of orders
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}

// CreateOrderResponse represents the response after creating an order
type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
}