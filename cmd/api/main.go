package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/application/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/config"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/database"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/repositories"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/presentation/http/handlers"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/presentation/http/routes"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Errorf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	logger.Info("Database migrations completed successfully")

	// Initialize repositories
	orderRepo := repositories.NewPostgresOrderRepository(db)
	productRepo := repositories.NewPostgresProductRepository(db)
	customerRepo := repositories.NewPostgresCustomerRepository(db)
	eventPub := repositories.NewSimpleEventPublisher()

	// Initialize application service
	orderService := order.NewApplicationService(orderRepo, productRepo, customerRepo, eventPub)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup routes
	router := routes.SetupRoutes(orderHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	logger.Infof("Order Management API is running on %s", server.Addr)
	logger.Info("API Documentation:")
	logger.Info("  Health Check: GET /health")
	logger.Info("  Create Order: POST /api/v1/orders")
	logger.Info("  Get Order: GET /api/v1/orders/:id")
	logger.Info("  Update Order Status: PUT /api/v1/orders/:id/status")
	logger.Info("  Get Customer Orders: GET /api/v1/customers/:customer_id/orders")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}