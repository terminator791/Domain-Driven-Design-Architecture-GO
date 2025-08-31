# Enterprise E-commerce Domain-Driven Design (DDD) System in Go

This project demonstrates a comprehensive implementation of **Domain-Driven Design (DDD)** principles using Go, built as an enterprise-level e-commerce system. It showcases advanced DDD concepts, CQRS, event sourcing principles, and complex business domains while maintaining clean architecture and separation of concerns.

## 🏗️ Enterprise Architecture Overview

The project follows a sophisticated layered architecture with multiple bounded contexts:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Presentation Layer                               │
│   (HTTP Handlers, Routes, Middleware, WebAPI Controllers)                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                          Application Layer                                 │
│   (CQRS Bus, Commands, Queries, Handlers, Read Models, DTOs)              │
├─────────────────────────────────────────────────────────────────────────────┤
│                   Multiple Domain Bounded Contexts                         │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   Order     │   Cart      │ Inventory   │  Payment    │ Promotion   │  │
│  │   Domain    │   Domain    │   Domain    │   Domain    │   Domain    │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘  │
│  (Entities, Value Objects, Domain Services, Events, Aggregates)            │
├─────────────────────────────────────────────────────────────────────────────┤
│                        Infrastructure Layer                                │
│  (Repositories, Database, Event Store, External APIs, Messaging)           │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🎯 Enterprise DDD Concepts Implemented

### 1. **Multiple Bounded Contexts**
- **Order Management**: Order lifecycle, status tracking, fulfillment
- **Shopping Cart**: Session-based cart management, item persistence, expiry
- **Inventory Management**: Stock levels, reservations, low-stock alerts
- **Payment Processing**: Multi-method payments, authorization/capture, refunds
- **Promotion Engine**: Coupons, discounts, promotional campaigns

### 2. **Advanced Ubiquitous Language**

#### Order Context
- **Order**: Customer's purchase request with multiple items
- **OrderItem**: Individual products with quantities and pricing
- **OrderStatus**: PENDING → CONFIRMED → SHIPPED → DELIVERED

#### Cart Context  
- **Cart**: Session-based shopping container with expiry
- **CartItem**: Products added to cart with real-time pricing
- **CartStatus**: ACTIVE → ABANDONED → CHECKED_OUT → EXPIRED

#### Inventory Context
- **InventoryItem**: Stock tracking with available/reserved quantities  
- **StockReservation**: Temporary stock allocation with expiry
- **AlertLevel**: NONE → LOW → CRITICAL → OUT_OF_STOCK

#### Payment Context
- **Payment**: Financial transaction with authorization flow
- **PaymentMethod**: Customer's stored payment instruments
- **PaymentStatus**: PENDING → AUTHORIZED → CAPTURED → REFUNDED

#### Promotion Context
- **Promotion**: Marketing campaigns with business rules
- **Coupon**: Individual discount codes with usage tracking
- **DiscountValue**: Percentage or fixed amount discounts

### 3. **Advanced Domain Patterns**

#### **Complex Aggregates**
- **Order Aggregate**: Order + OrderItems + Domain Events
- **Cart Aggregate**: Cart + CartItems + Business Rules  
- **Payment Aggregate**: Payment + Refunds + Status Transitions
- **Promotion Aggregate**: Promotion + Coupons + Usage Tracking
- **Inventory Aggregate**: InventoryItem + Reservations + Alerts

#### **Entities** (`internal/domain/*/entity.go`)
- **Order**: Aggregate root containing order items and business logic
- **OrderItem**: Entity representing individual items in an order
- **Product**: Entity representing products in the catalog
- **Customer**: Entity representing customers

#### **Rich Value Objects**
- **Money**: Currency-aware monetary calculations
- **StockLevel**: Available/Reserved/Total quantity tracking  
- **DiscountValue**: Percentage and fixed amount discounts
- **ReservationID, PaymentID, CartID**: Strong typing for identifiers
- **PaymentStatus, CartStatus, AlertLevel**: Enumerated states with business rules

#### **Domain Services** 
- **OrderDomainService**: Complex order validation and pricing
- **CartDomainService**: Cart lifecycle and business rules
- **InventoryDomainService**: Stock reservation and replenishment
- **PromotionEngine**: Discount calculation and coupon validation

#### **Domain Events & Event Sourcing**
- **OrderCreatedEvent, PaymentCapturedEvent**: Business significant events
- **CartAbandonedEvent, StockReservedEvent**: Cross-domain communication
- **PromotionAppliedEvent, LowStockAlertEvent**: Business process triggers
- Event-driven architecture for loose coupling

#### **Repository Patterns**
- Abstract persistence interfaces without infrastructure coupling
- Separate read/write models for CQRS optimization
- Event sourcing capabilities for audit trail

### 4. **CQRS (Command Query Responsibility Segregation)**

#### **Command Side - Write Models** (`internal/application/cqrs/commands/`)
- **CreateOrderCommand**: Order creation with validation
- **AddToCartCommand**: Cart item management  
- **ProcessPaymentCommand**: Payment processing workflow
- **UpdateInventoryCommand**: Stock level adjustments

#### **Query Side - Read Models** (`internal/application/readmodels/`)
- **OrderReadModel**: Optimized order views with customer data
- **ProductCatalogReadModel**: Searchable product catalog with filters
- **InventoryReadModel**: Real-time stock levels and alerts
- **PaymentReadModel**: Payment history and transaction details
- **SalesReportReadModel**: Analytics and business intelligence

#### **Command/Query Buses** (`internal/application/cqrs/`)
- **CommandBus**: Routes commands to appropriate handlers
- **QueryBus**: Routes queries to optimized read handlers  
- **EventBus**: Publishes domain events to event handlers

### 5. **Enterprise Business Features**

#### **Shopping Cart Management**
- Session-based cart with automatic expiry (24 hours)
- Real-time pricing updates and stock validation
- Cart abandonment recovery and analytics
- Business rules: max items per product (10), max total items (50)

#### **Advanced Inventory Management**  
- **Stock Reservations**: Temporary allocation during checkout (30 min TTL)
- **Multi-level Alerts**: None → Low → Critical → Out of Stock
- **Automated Replenishment**: Reorder point triggers and notifications
- **Real-time Tracking**: Available, Reserved, and Total stock levels

#### **Enterprise Payment Processing**
- **Multi-method Support**: Credit/Debit cards, PayPal, Bank transfers, Digital wallets, Crypto
- **Authorization/Capture Flow**: Two-phase payment processing
- **Partial Refunds**: Business rule-based refund management
- **Payment Method Storage**: Customer payment instrument management

#### **Promotion Engine & Discounts**
- **Flexible Discount Types**: Percentage, Fixed amount, Buy-X-Get-Y, Free shipping
- **Advanced Conditions**: Minimum order amount, eligible products/categories, customer tiers
- **Usage Limits**: Total uses and per-customer restrictions
- **Coupon Management**: Individual codes with expiry and customer targeting

### 6. **Advanced DDD Patterns Implemented**

#### **Event Sourcing Foundations**
- Domain events capture all business-significant state changes
- Event-driven communication between bounded contexts
- Audit trail for compliance and business analytics

#### **CQRS Separation**
- Commands modify state through domain aggregates
- Queries use optimized read models for performance
- Separate scaling and optimization strategies

#### **Saga Pattern Ready**
- Event-driven workflows for complex business processes
- Compensation actions for distributed transaction management
- Order processing workflow: Cart → Inventory → Payment → Fulfillment

### 7. **Infrastructure & Persistence Layer**

#### **Database**
- PostgreSQL with optimized schemas for read/write separation
- Repository implementations using sqlx for performance
- Database migrations with proper indexing strategies

#### **Configuration** 
- Viper-based configuration management
- Environment-specific settings and secrets management
- Docker Compose for development environment

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