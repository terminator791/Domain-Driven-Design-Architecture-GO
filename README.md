# Domain-Driven Design (DDD) Order Management System in Go

This project demonstrates a complete implementation of **Domain-Driven Design (DDD)** principles using Go, built as an Order Management System API. It showcases the core DDD concepts while maintaining clean architecture and separation of concerns.

## 🏗️ Architecture Overview

The project follows a layered architecture aligned with DDD principles:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  (HTTP Handlers, Routes, Middleware)                       │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                        │
│  (Use Cases, Commands, Handlers, DTOs)                     │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                           │
│  (Entities, Value Objects, Domain Services, Repositories)  │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                      │
│  (Database, Repositories, Configuration)                   │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 DDD Concepts Implemented

### 1. **Ubiquitous Language**
- **Order**: A customer's request for one or more products
- **OrderItem**: Individual products within an order with quantities
- **Customer**: The person placing orders
- **Product**: Items available for purchase
- **OrderStatus**: Current state of an order (PENDING, CONFIRMED, SHIPPED, etc.)

### 2. **Bounded Context**
- **Order Management Context**: Handles order creation, updates, and tracking
- Clear boundaries between Order, Product, and Customer domains

### 3. **Domain Layer Components**

#### **Entities** (`internal/domain/*/entity.go`)
- **Order**: Aggregate root containing order items and business logic
- **OrderItem**: Entity representing individual items in an order
- **Product**: Entity representing products in the catalog
- **Customer**: Entity representing customers

#### **Value Objects** (`internal/domain/order/value_objects.go`)
- **OrderID, ProductID, CustomerID**: Strongly-typed identifiers
- **Money**: Represents monetary values with currency
- **Quantity**: Represents item quantities with validation
- **OrderStatus**: Enumerated order states

#### **Aggregates**
- **Order Aggregate**: Order (root) + OrderItems
- Ensures consistency boundaries and invariants

#### **Domain Services** (`internal/domain/order/service.go`)
- **OrderDomainService**: Complex business logic that doesn't belong to entities
- Validates order creation rules
- Calculates discounts and totals

#### **Domain Events**
- **OrderCreatedEvent**: Published when an order is created
- **OrderUpdatedEvent**: Published when order status changes

#### **Repository Interfaces** (`internal/domain/*/repository.go`)
- Abstract data access without coupling to infrastructure
- Follows Repository pattern

### 4. **Application Layer** (`internal/application/order/`)

#### **Commands and Handlers**
- **CreateOrderRequest/Handler**: Creates new orders
- **UpdateOrderStatusRequest/Handler**: Updates order status
- **GetOrderHandler**: Retrieves order information

#### **Application Services**
- **OrderApplicationService**: Orchestrates use cases
- Coordinates between domain services and repositories

#### **DTOs (Data Transfer Objects)**
- **OrderResponse**: API response models
- **CreateOrderResponse**: Command response models

### 5. **Infrastructure Layer** (`internal/infrastructure/`)

#### **Database**
- PostgreSQL with proper migrations
- Repository implementations using sqlx

#### **Configuration**
- Viper-based configuration management
- Environment variable support

## 🚀 Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15 Alpine
- **Database Access**: sqlx
- **Configuration**: Viper
- **Containerization**: Docker & Docker Compose
- **Architecture**: Domain-Driven Design (DDD)

## 📁 Project Structure

```
cmd/
  api/
    main.go                 # Application entry point
internal/
  domain/                   # Domain layer (business logic)
    order/
      entity.go             # Order aggregate and entities
      value_objects.go      # Value objects (Money, Quantity, etc.)
      repository.go         # Repository interfaces
      service.go            # Domain services
    product/
      entity.go             # Product entity
    customer/
      entity.go             # Customer entity
    common/
      errors.go             # Domain errors
  application/              # Application layer (use cases)
    order/
      commands.go           # Command objects
      handlers.go           # Command handlers
      service.go            # Application services
      dtos.go               # Data transfer objects
  infrastructure/           # Infrastructure layer (external concerns)
    config/
      config.go             # Configuration management
    database/
      postgres.go           # Database connection and migrations
    repositories/
      order_repository.go   # Order repository implementation
      product_repository.go # Product repository implementation
      customer_repository.go # Customer repository implementation
  presentation/             # Presentation layer (HTTP API)
    http/
      handlers/
        order_handler.go    # HTTP handlers
      middleware/
        error_handler.go    # HTTP middleware
      routes/
        routes.go           # Route definitions
pkg/
  logger/
    logger.go               # Logging utilities
docker-compose.yml          # Docker composition
Dockerfile                  # Docker image definition
config.yaml                 # Application configuration
```

## 🏃‍♂️ Running the Application

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Using Docker Compose (Recommended)
```bash
# Start the complete stack (API + PostgreSQL)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop the stack
docker-compose down
```

### Local Development
```bash
# Install dependencies
go mod download

# Start PostgreSQL (using Docker)
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=orderdb -p 5432:5432 -d postgres:15-alpine

# Run the application
go run cmd/api/main.go
```

## 📚 API Documentation

The API provides endpoints for order management:

### Health Check
```http
GET /health
```

### Order Operations
```http
# Create a new order
POST /api/v1/orders
Content-Type: application/json

{
  "customer_id": "customer-uuid",
  "items": [
    {
      "product_id": "product-uuid",
      "quantity": 2,
      "unit_price": 19.99,
      "currency": "USD"
    }
  ]
}

# Get order by ID
GET /api/v1/orders/{order_id}

# Update order status
PUT /api/v1/orders/{order_id}/status
Content-Type: application/json

{
  "status": "CONFIRMED"
}

# Get orders by customer
GET /api/v1/customers/{customer_id}/orders
```

### Example Usage
```bash
# Health check
curl http://localhost:8080/health

# Create an order (requires existing customer and product IDs)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "123e4567-e89b-12d3-a456-426614174000",
    "items": [
      {
        "product_id": "123e4567-e89b-12d3-a456-426614174001",
        "quantity": 2,
        "unit_price": 29.99,
        "currency": "USD"
      }
    ]
  }'
```

## 🧪 Database Schema

The application automatically creates the following tables:

- **customers**: Customer information
- **products**: Product catalog
- **orders**: Order headers
- **order_items**: Order line items

## 🔑 Key DDD Benefits Demonstrated

1. **Clear Business Logic**: Domain entities encapsulate business rules
2. **Testability**: Business logic is independent of infrastructure
3. **Maintainability**: Clear separation of concerns
4. **Flexibility**: Easy to modify business rules without affecting other layers
5. **Domain Events**: Loose coupling through event-driven architecture
6. **Rich Domain Model**: Behavior-rich entities rather than anemic models

## 🛠️ Advanced DDD Patterns Implemented

1. **Aggregate Pattern**: Order as aggregate root with OrderItems
2. **Repository Pattern**: Abstract data access
3. **Domain Services**: Complex business logic
4. **Value Objects**: Immutable, validated data types
5. **Domain Events**: Eventual consistency and decoupling
6. **Application Services**: Use case orchestration
7. **Anti-Corruption Layer**: Clean interfaces between layers

## 🔮 Future Enhancements

- Event sourcing for order state changes
- CQRS (Command Query Responsibility Segregation)
- Saga pattern for distributed transactions
- Integration events for microservices communication
- Domain event store and replay capabilities

## 📖 Learning Resources

This implementation demonstrates practical DDD concepts. For deeper understanding:

- "Domain-Driven Design" by Eric Evans
- "Implementing Domain-Driven Design" by Vaughn Vernon
- Clean Architecture principles
- Hexagonal Architecture patterns

## 🤝 Contributing

This project serves as a learning resource for DDD in Go. Contributions that enhance the demonstration of DDD principles are welcome!