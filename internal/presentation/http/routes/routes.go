package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/presentation/http/handlers"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/presentation/http/middleware"
)

func SetupRoutes(orderHandler *handlers.OrderHandler) *gin.Engine {
	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", orderHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}

		// Customer order routes
		customers := v1.Group("/customers")
		{
			customers.GET("/:customer_id/orders", orderHandler.GetOrdersByCustomer)
		}
	}

	return router
}