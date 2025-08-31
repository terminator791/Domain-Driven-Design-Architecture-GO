package readmodels

import (
	"time"
)

// OrderReadModel represents the read-optimized view of an order
type OrderReadModel struct {
	ID                string                   `json:"id" db:"id"`
	CustomerID        string                   `json:"customer_id" db:"customer_id"`
	CustomerName      string                   `json:"customer_name" db:"customer_name"`
	CustomerEmail     string                   `json:"customer_email" db:"customer_email"`
	Status            string                   `json:"status" db:"status"`
	TotalAmount       float64                  `json:"total_amount" db:"total_amount"`
	TotalCurrency     string                   `json:"total_currency" db:"total_currency"`
	ItemCount         int                      `json:"item_count" db:"item_count"`
	CreatedAt         time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at" db:"updated_at"`
	Items             []OrderItemReadModel     `json:"items"`
	PaymentStatus     string                   `json:"payment_status" db:"payment_status"`
	ShippingStatus    string                   `json:"shipping_status" db:"shipping_status"`
	EstimatedDelivery *time.Time               `json:"estimated_delivery,omitempty" db:"estimated_delivery"`
}

// OrderItemReadModel represents the read-optimized view of an order item
type OrderItemReadModel struct {
	ProductID       string  `json:"product_id" db:"product_id"`
	ProductName     string  `json:"product_name" db:"product_name"`
	ProductImage    string  `json:"product_image" db:"product_image"`
	Quantity        int     `json:"quantity" db:"quantity"`
	UnitPrice       float64 `json:"unit_price" db:"unit_price"`
	TotalPrice      float64 `json:"total_price" db:"total_price"`
	Currency        string  `json:"currency" db:"currency"`
}

// CartReadModel represents the read-optimized view of a shopping cart
type CartReadModel struct {
	ID           string                 `json:"id" db:"id"`
	CustomerID   string                 `json:"customer_id" db:"customer_id"`
	Status       string                 `json:"status" db:"status"`
	TotalAmount  float64                `json:"total_amount" db:"total_amount"`
	TotalCurrency string                `json:"total_currency" db:"total_currency"`
	ItemCount    int                    `json:"item_count" db:"item_count"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	ExpiresAt    time.Time              `json:"expires_at" db:"expires_at"`
	Items        []CartItemReadModel    `json:"items"`
}

// CartItemReadModel represents the read-optimized view of a cart item
type CartItemReadModel struct {
	ProductID    string  `json:"product_id" db:"product_id"`
	ProductName  string  `json:"product_name" db:"product_name"`
	ProductImage string  `json:"product_image" db:"product_image"`
	Quantity     int     `json:"quantity" db:"quantity"`
	UnitPrice    float64 `json:"unit_price" db:"unit_price"`
	TotalPrice   float64 `json:"total_price" db:"total_price"`
	Currency     string  `json:"currency" db:"currency"`
	AddedAt      time.Time `json:"added_at" db:"added_at"`
	InStock      bool    `json:"in_stock" db:"in_stock"`
	StockLevel   int     `json:"stock_level" db:"stock_level"`
}

// ProductCatalogReadModel represents the read-optimized view of products for catalog
type ProductCatalogReadModel struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Price         float64   `json:"price" db:"price"`
	Currency      string    `json:"currency" db:"currency"`
	Category      string    `json:"category" db:"category"`
	SubCategory   string    `json:"sub_category" db:"sub_category"`
	Brand         string    `json:"brand" db:"brand"`
	SKU           string    `json:"sku" db:"sku"`
	ImageURL      string    `json:"image_url" db:"image_url"`
	Images        []string  `json:"images"`
	StockLevel    int       `json:"stock_level" db:"stock_level"`
	IsAvailable   bool      `json:"is_available" db:"is_available"`
	Rating        float64   `json:"rating" db:"rating"`
	ReviewCount   int       `json:"review_count" db:"review_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Tags          []string  `json:"tags"`
	Specifications map[string]string `json:"specifications"`
}

// InventoryReadModel represents the read-optimized view of inventory status
type InventoryReadModel struct {
	ProductID       string    `json:"product_id" db:"product_id"`
	ProductName     string    `json:"product_name" db:"product_name"`
	ProductSKU      string    `json:"product_sku" db:"product_sku"`
	AvailableStock  int       `json:"available_stock" db:"available_stock"`
	ReservedStock   int       `json:"reserved_stock" db:"reserved_stock"`
	TotalStock      int       `json:"total_stock" db:"total_stock"`
	MinStockLevel   int       `json:"min_stock_level" db:"min_stock_level"`
	ReorderPoint    int       `json:"reorder_point" db:"reorder_point"`
	AlertLevel      string    `json:"alert_level" db:"alert_level"`
	LastUpdated     time.Time `json:"last_updated" db:"last_updated"`
	PendingOrders   int       `json:"pending_orders" db:"pending_orders"`
	ReservationCount int      `json:"reservation_count" db:"reservation_count"`
}

// PaymentReadModel represents the read-optimized view of payments
type PaymentReadModel struct {
	ID                string               `json:"id" db:"id"`
	OrderID           string               `json:"order_id" db:"order_id"`
	CustomerID        string               `json:"customer_id" db:"customer_id"`
	PaymentMethodID   string               `json:"payment_method_id" db:"payment_method_id"`
	PaymentMethodType string               `json:"payment_method_type" db:"payment_method_type"`
	PaymentMethodName string               `json:"payment_method_name" db:"payment_method_name"`
	Amount            float64              `json:"amount" db:"amount"`
	Currency          string               `json:"currency" db:"currency"`
	Status            string               `json:"status" db:"status"`
	AuthorizedAmount  float64              `json:"authorized_amount" db:"authorized_amount"`
	CapturedAmount    float64              `json:"captured_amount" db:"captured_amount"`
	RefundedAmount    float64              `json:"refunded_amount" db:"refunded_amount"`
	TransactionID     string               `json:"transaction_id" db:"transaction_id"`
	CreatedAt         time.Time            `json:"created_at" db:"created_at"`
	AuthorizedAt      *time.Time           `json:"authorized_at,omitempty" db:"authorized_at"`
	CapturedAt        *time.Time           `json:"captured_at,omitempty" db:"captured_at"`
	Refunds           []RefundReadModel    `json:"refunds"`
}

// RefundReadModel represents the read-optimized view of refunds
type RefundReadModel struct {
	ID          string    `json:"id" db:"id"`
	PaymentID   string    `json:"payment_id" db:"payment_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Currency    string    `json:"currency" db:"currency"`
	Reason      string    `json:"reason" db:"reason"`
	RefundID    string    `json:"refund_id" db:"refund_id"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// CustomerOrderSummaryReadModel represents summary of customer orders
type CustomerOrderSummaryReadModel struct {
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	CustomerName      string    `json:"customer_name" db:"customer_name"`
	CustomerEmail     string    `json:"customer_email" db:"customer_email"`
	TotalOrders       int       `json:"total_orders" db:"total_orders"`
	TotalSpent        float64   `json:"total_spent" db:"total_spent"`
	Currency          string    `json:"currency" db:"currency"`
	AverageOrderValue float64   `json:"average_order_value" db:"average_order_value"`
	LastOrderDate     *time.Time `json:"last_order_date,omitempty" db:"last_order_date"`
	FirstOrderDate    *time.Time `json:"first_order_date,omitempty" db:"first_order_date"`
	MostPurchasedCategory string `json:"most_purchased_category" db:"most_purchased_category"`
}

// SalesReportReadModel represents sales analytics data
type SalesReportReadModel struct {
	Period             string  `json:"period" db:"period"`
	TotalOrders        int     `json:"total_orders" db:"total_orders"`
	TotalRevenue       float64 `json:"total_revenue" db:"total_revenue"`
	Currency           string  `json:"currency" db:"currency"`
	AverageOrderValue  float64 `json:"average_order_value" db:"average_order_value"`
	NewCustomers       int     `json:"new_customers" db:"new_customers"`
	ReturningCustomers int     `json:"returning_customers" db:"returning_customers"`
	TopSellingProducts []TopProductReadModel `json:"top_selling_products"`
	TopCategories      []CategorySalesReadModel `json:"top_categories"`
}

// TopProductReadModel represents top selling products
type TopProductReadModel struct {
	ProductID    string  `json:"product_id" db:"product_id"`
	ProductName  string  `json:"product_name" db:"product_name"`
	QuantitySold int     `json:"quantity_sold" db:"quantity_sold"`
	Revenue      float64 `json:"revenue" db:"revenue"`
	Currency     string  `json:"currency" db:"currency"`
}

// CategorySalesReadModel represents category sales data
type CategorySalesReadModel struct {
	Category     string  `json:"category" db:"category"`
	TotalOrders  int     `json:"total_orders" db:"total_orders"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
	Currency     string  `json:"currency" db:"currency"`
}

// PaginationResult represents pagination metadata
type PaginationResult struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// PaginatedResult wraps results with pagination
type PaginatedResult struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationResult `json:"pagination"`
}