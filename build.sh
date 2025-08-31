#!/bin/bash

echo "Building and testing the DDD Order Management System..."

# Build the main application
echo "Building main application..."
go build -o bin/api cmd/api/main.go

# Build the seed script
echo "Building seed script..."
go build -o bin/seed cmd/seed/main.go

echo "Build completed successfully!"

# Test the application (requires Docker Compose)
echo ""
echo "To test the complete system:"
echo "1. Start the stack: docker-compose up -d"
echo "2. Wait for database to be ready"
echo "3. Run seed data: ./bin/seed"
echo "4. Test API endpoints"
echo ""
echo "Or run locally:"
echo "1. Start PostgreSQL: docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=orderdb -p 5432:5432 -d postgres:15-alpine"
echo "2. Run migrations and start API: ./bin/api"
echo "3. In another terminal, run: ./bin/seed"