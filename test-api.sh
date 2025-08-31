#!/bin/bash

echo "Quick API Test Script"
echo "===================="

# Check if the API is running
echo "Checking if API is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ API is running!"
    
    echo ""
    echo "Testing health endpoint:"
    curl -s http://localhost:8080/health | jq '.' || curl -s http://localhost:8080/health

    echo -e "\n\nExample API calls you can try:"
    echo ""
    echo "1. Create an order (replace with actual customer and product IDs from seed data):"
    echo 'curl -X POST http://localhost:8080/api/v1/orders \'
    echo '  -H "Content-Type: application/json" \'
    echo '  -d '\''{'
    echo '    "customer_id": "CUSTOMER_ID_FROM_SEED",''
    echo '    "items": ['
    echo '      {'
    echo '        "product_id": "PRODUCT_ID_FROM_SEED",'
    echo '        "quantity": 1,'
    echo '        "unit_price": 999.99,'
    echo '        "currency": "USD"'
    echo '      }'
    echo '    ]'
    echo '  }'\'''

    echo ""
    echo "2. Get order by ID:"
    echo "curl http://localhost:8080/api/v1/orders/ORDER_ID"

    echo ""
    echo "3. Update order status:"
    echo 'curl -X PUT http://localhost:8080/api/v1/orders/ORDER_ID/status \'
    echo '  -H "Content-Type: application/json" \'
    echo '  -d '\''{"status": "CONFIRMED"}'\'''

else
    echo "❌ API is not running. Start it with:"
    echo "  make docker-up && make seed"
    echo "  or"
    echo "  make run (in another terminal, ensure PostgreSQL is running)"
fi