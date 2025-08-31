package main

import (
	"context"
	"fmt"
	"log"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/customer"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/product"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/config"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/database"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/infrastructure/repositories"
)

func main() {
	fmt.Println("Creating sample data...")

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	customerRepo := repositories.NewPostgresCustomerRepository(db)
	productRepo := repositories.NewPostgresProductRepository(db)

	ctx := context.Background()

	// Create sample customers
	customers := []*customer.Customer{
		customer.NewCustomer("John Doe", "john.doe@example.com", "+1234567890"),
		customer.NewCustomer("Jane Smith", "jane.smith@example.com", "+1234567891"),
		customer.NewCustomer("Bob Johnson", "bob.johnson@example.com", "+1234567892"),
	}

	fmt.Println("Creating customers...")
	for _, c := range customers {
		err := customerRepo.Save(ctx, c)
		if err != nil {
			log.Printf("Failed to save customer %s: %v", c.Name(), err)
		} else {
			fmt.Printf("Created customer: %s (ID: %s)\n", c.Name(), c.ID().String())
		}
	}

	// Create sample products
	products := []*product.Product{}
	
	// Product 1: Laptop
	laptopPrice, _ := order.NewMoney(999.99, "USD")
	laptop := product.NewProduct("MacBook Pro", "Apple MacBook Pro 16-inch", laptopPrice, 10)
	products = append(products, laptop)

	// Product 2: Mouse
	mousePrice, _ := order.NewMoney(79.99, "USD")
	mouse := product.NewProduct("Magic Mouse", "Apple Magic Mouse", mousePrice, 50)
	products = append(products, mouse)

	// Product 3: Keyboard
	keyboardPrice, _ := order.NewMoney(199.99, "USD")
	keyboard := product.NewProduct("Magic Keyboard", "Apple Magic Keyboard", keyboardPrice, 30)
	products = append(products, keyboard)

	// Product 4: Monitor
	monitorPrice, _ := order.NewMoney(599.99, "USD")
	monitor := product.NewProduct("Studio Display", "Apple Studio Display", monitorPrice, 5)
	products = append(products, monitor)

	fmt.Println("Creating products...")
	for _, p := range products {
		err := productRepo.Save(ctx, p)
		if err != nil {
			log.Printf("Failed to save product %s: %v", p.Name(), err)
		} else {
			fmt.Printf("Created product: %s (ID: %s, Price: %s)\n", 
				p.Name(), p.ID().String(), p.Price().String())
		}
	}

	fmt.Println("\nSample data created successfully!")
	fmt.Println("\nYou can now test the API with the following IDs:")
	
	fmt.Println("\nCustomers:")
	for _, c := range customers {
		fmt.Printf("- %s: %s\n", c.Name(), c.ID().String())
	}
	
	fmt.Println("\nProducts:")
	for _, p := range products {
		fmt.Printf("- %s: %s\n", p.Name(), p.ID().String())
	}

	fmt.Println("\nExample API calls:")
	fmt.Printf(`
# Create an order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "%s",
    "items": [
      {
        "product_id": "%s",
        "quantity": 1,
        "unit_price": %.2f,
        "currency": "USD"
      },
      {
        "product_id": "%s",
        "quantity": 2,
        "unit_price": %.2f,
        "currency": "USD"
      }
    ]
  }'
`, customers[0].ID().String(), products[0].ID().String(), products[0].Price().Amount(), 
   products[1].ID().String(), products[1].Price().Amount())
}