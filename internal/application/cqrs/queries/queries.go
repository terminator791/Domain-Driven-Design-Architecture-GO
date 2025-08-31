package queries

import (
	"time"

	"github.com/google/uuid"
)

// GetOrderByIDQuery represents a query to get an order by ID
type GetOrderByIDQuery struct {
	queryID string
	orderID string
}

func NewGetOrderByIDQuery(orderID string) *GetOrderByIDQuery {
	return &GetOrderByIDQuery{
		queryID: uuid.New().String(),
		orderID: orderID,
	}
}

func (q *GetOrderByIDQuery) QueryID() string {
	return q.queryID
}

func (q *GetOrderByIDQuery) QueryType() string {
	return "GetOrderByID"
}

func (q *GetOrderByIDQuery) GetFilter() interface{} {
	return struct {
		OrderID string `json:"order_id"`
	}{
		OrderID: q.orderID,
	}
}

func (q *GetOrderByIDQuery) OrderID() string {
	return q.orderID
}

// GetCustomerOrdersQuery represents a query to get all orders for a customer
type GetCustomerOrdersQuery struct {
	queryID    string
	customerID string
	page       int
	pageSize   int
	status     string
	fromDate   *time.Time
	toDate     *time.Time
}

func NewGetCustomerOrdersQuery(customerID string, page, pageSize int, status string) *GetCustomerOrdersQuery {
	return &GetCustomerOrdersQuery{
		queryID:    uuid.New().String(),
		customerID: customerID,
		page:       page,
		pageSize:   pageSize,
		status:     status,
	}
}

func (q *GetCustomerOrdersQuery) QueryID() string {
	return q.queryID
}

func (q *GetCustomerOrdersQuery) QueryType() string {
	return "GetCustomerOrders"
}

func (q *GetCustomerOrdersQuery) GetFilter() interface{} {
	return struct {
		CustomerID string     `json:"customer_id"`
		Page       int        `json:"page"`
		PageSize   int        `json:"page_size"`
		Status     string     `json:"status,omitempty"`
		FromDate   *time.Time `json:"from_date,omitempty"`
		ToDate     *time.Time `json:"to_date,omitempty"`
	}{
		CustomerID: q.customerID,
		Page:       q.page,
		PageSize:   q.pageSize,
		Status:     q.status,
		FromDate:   q.fromDate,
		ToDate:     q.toDate,
	}
}

func (q *GetCustomerOrdersQuery) CustomerID() string {
	return q.customerID
}

func (q *GetCustomerOrdersQuery) Page() int {
	return q.page
}

func (q *GetCustomerOrdersQuery) PageSize() int {
	return q.pageSize
}

func (q *GetCustomerOrdersQuery) Status() string {
	return q.status
}

func (q *GetCustomerOrdersQuery) SetDateRange(from, to time.Time) {
	q.fromDate = &from
	q.toDate = &to
}

// GetCartByCustomerQuery represents a query to get active cart for customer
type GetCartByCustomerQuery struct {
	queryID    string
	customerID string
}

func NewGetCartByCustomerQuery(customerID string) *GetCartByCustomerQuery {
	return &GetCartByCustomerQuery{
		queryID:    uuid.New().String(),
		customerID: customerID,
	}
}

func (q *GetCartByCustomerQuery) QueryID() string {
	return q.queryID
}

func (q *GetCartByCustomerQuery) QueryType() string {
	return "GetCartByCustomer"
}

func (q *GetCartByCustomerQuery) GetFilter() interface{} {
	return struct {
		CustomerID string `json:"customer_id"`
	}{
		CustomerID: q.customerID,
	}
}

func (q *GetCartByCustomerQuery) CustomerID() string {
	return q.customerID
}

// GetProductCatalogQuery represents a query to get products with filtering and pagination
type GetProductCatalogQuery struct {
	queryID     string
	page        int
	pageSize    int
	category    string
	searchTerm  string
	minPrice    float64
	maxPrice    float64
	sortBy      string
	sortOrder   string
	inStock     bool
}

func NewGetProductCatalogQuery(page, pageSize int) *GetProductCatalogQuery {
	return &GetProductCatalogQuery{
		queryID:   uuid.New().String(),
		page:      page,
		pageSize:  pageSize,
		sortBy:    "name",
		sortOrder: "ASC",
		inStock:   true,
	}
}

func (q *GetProductCatalogQuery) QueryID() string {
	return q.queryID
}

func (q *GetProductCatalogQuery) QueryType() string {
	return "GetProductCatalog"
}

func (q *GetProductCatalogQuery) GetFilter() interface{} {
	return struct {
		Page       int     `json:"page"`
		PageSize   int     `json:"page_size"`
		Category   string  `json:"category,omitempty"`
		SearchTerm string  `json:"search_term,omitempty"`
		MinPrice   float64 `json:"min_price,omitempty"`
		MaxPrice   float64 `json:"max_price,omitempty"`
		SortBy     string  `json:"sort_by"`
		SortOrder  string  `json:"sort_order"`
		InStock    bool    `json:"in_stock"`
	}{
		Page:       q.page,
		PageSize:   q.pageSize,
		Category:   q.category,
		SearchTerm: q.searchTerm,
		MinPrice:   q.minPrice,
		MaxPrice:   q.maxPrice,
		SortBy:     q.sortBy,
		SortOrder:  q.sortOrder,
		InStock:    q.inStock,
	}
}

func (q *GetProductCatalogQuery) Page() int {
	return q.page
}

func (q *GetProductCatalogQuery) PageSize() int {
	return q.pageSize
}

func (q *GetProductCatalogQuery) SetCategory(category string) {
	q.category = category
}

func (q *GetProductCatalogQuery) SetSearchTerm(searchTerm string) {
	q.searchTerm = searchTerm
}

func (q *GetProductCatalogQuery) SetPriceRange(min, max float64) {
	q.minPrice = min
	q.maxPrice = max
}

func (q *GetProductCatalogQuery) SetSort(sortBy, sortOrder string) {
	q.sortBy = sortBy
	q.sortOrder = sortOrder
}

func (q *GetProductCatalogQuery) SetInStockOnly(inStock bool) {
	q.inStock = inStock
}

// GetInventoryStatusQuery represents a query to get inventory status
type GetInventoryStatusQuery struct {
	queryID     string
	productID   string
	lowStockOnly bool
}

func NewGetInventoryStatusQuery(productID string, lowStockOnly bool) *GetInventoryStatusQuery {
	return &GetInventoryStatusQuery{
		queryID:      uuid.New().String(),
		productID:    productID,
		lowStockOnly: lowStockOnly,
	}
}

func (q *GetInventoryStatusQuery) QueryID() string {
	return q.queryID
}

func (q *GetInventoryStatusQuery) QueryType() string {
	return "GetInventoryStatus"
}

func (q *GetInventoryStatusQuery) GetFilter() interface{} {
	return struct {
		ProductID    string `json:"product_id,omitempty"`
		LowStockOnly bool   `json:"low_stock_only"`
	}{
		ProductID:    q.productID,
		LowStockOnly: q.lowStockOnly,
	}
}

func (q *GetInventoryStatusQuery) ProductID() string {
	return q.productID
}

func (q *GetInventoryStatusQuery) LowStockOnly() bool {
	return q.lowStockOnly
}

// GetPaymentHistoryQuery represents a query to get payment history
type GetPaymentHistoryQuery struct {
	queryID    string
	customerID string
	orderID    string
	status     string
	page       int
	pageSize   int
}

func NewGetPaymentHistoryQuery(customerID, orderID, status string, page, pageSize int) *GetPaymentHistoryQuery {
	return &GetPaymentHistoryQuery{
		queryID:    uuid.New().String(),
		customerID: customerID,
		orderID:    orderID,
		status:     status,
		page:       page,
		pageSize:   pageSize,
	}
}

func (q *GetPaymentHistoryQuery) QueryID() string {
	return q.queryID
}

func (q *GetPaymentHistoryQuery) QueryType() string {
	return "GetPaymentHistory"
}

func (q *GetPaymentHistoryQuery) GetFilter() interface{} {
	return struct {
		CustomerID string `json:"customer_id,omitempty"`
		OrderID    string `json:"order_id,omitempty"`
		Status     string `json:"status,omitempty"`
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
	}{
		CustomerID: q.customerID,
		OrderID:    q.orderID,
		Status:     q.status,
		Page:       q.page,
		PageSize:   q.pageSize,
	}
}

func (q *GetPaymentHistoryQuery) CustomerID() string {
	return q.customerID
}

func (q *GetPaymentHistoryQuery) OrderID() string {
	return q.orderID
}

func (q *GetPaymentHistoryQuery) Status() string {
	return q.status
}

func (q *GetPaymentHistoryQuery) Page() int {
	return q.page
}

func (q *GetPaymentHistoryQuery) PageSize() int {
	return q.pageSize
}