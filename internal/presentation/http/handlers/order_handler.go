package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/application/order"
)

type OrderHandler struct {
	orderService *order.ApplicationService
}

func NewOrderHandler(orderService *order.ApplicationService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.Error(gin.Error{Err: fmt.Errorf("order ID is required")})
		return
	}

	response, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOrdersByCustomer handles GET /customers/:customer_id/orders
func (h *OrderHandler) GetOrdersByCustomer(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		c.Error(gin.Error{Err: fmt.Errorf("customer ID is required")})
		return
	}

	response, err := h.orderService.GetOrdersByCustomer(c.Request.Context(), customerID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrderStatus handles PUT /orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.Error(gin.Error{Err: fmt.Errorf("order ID is required")})
		return
	}

	var req order.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := h.orderService.UpdateOrderStatus(c.Request.Context(), orderID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// HealthCheck handles GET /health
func (h *OrderHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "order-management-api",
		"version": "1.0.0",
	})
}