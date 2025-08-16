#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to make API calls with consistent formatting
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${BLUE}${description}${NC}"
    echo "Endpoint: ${method} ${endpoint}"
    
    if [ -n "$data" ]; then
        echo "Data: ${data}"
    fi
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s "${endpoint}")
    elif [ "$method" = "POST" ]; then
        response=$(curl -s -X POST -H "Content-Type: application/json" -d "${data}" "${endpoint}")
    elif [ "$method" = "PUT" ]; then
        response=$(curl -s -X PUT -H "Content-Type: application/json" -d "${data}" "${endpoint}")
    elif [ "$method" = "DELETE" ]; then
        response=$(curl -s -X DELETE "${endpoint}")
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Success${NC}"
        echo "Response: ${response}"
    else
        echo -e "${RED}❌ Failed${NC}"
        echo "Response: ${response}"
    fi
}

echo -e "${YELLOW}🚀 Expense API Demo - Full CRUD Operations${NC}"
echo "=========================================="
echo ""

# Check API health
echo -e "${BLUE}Checking API health...${NC}"
health_response=$(curl -s http://localhost:8080/health)
if echo "$health_response" | grep -q "healthy"; then
    echo -e "${GREEN}✅ API is running${NC}"
else
    echo -e "${RED}❌ API is not responding${NC}"
    exit 1
fi

echo ""

# 1. Get all categories
api_call "GET" "http://localhost:8080/api/categories" "" "1. Getting all categories"

# 2. Create a new category
api_call "POST" "http://localhost:8080/api/categories" '{"name": "Entertainment", "type": "expense"}' "2. Creating a new category"

# 3. Get the newly created category (assuming it got ID 10)
api_call "GET" "http://localhost:8080/api/categories/10" "" "3. Getting the new category"

# 4. Update the category
api_call "PUT" "http://localhost:8080/api/categories/10" '{"name": "Movies & Entertainment"}' "4. Updating the category"

# 5. Create a transaction using the new category
api_call "POST" "http://localhost:8080/api/transactions" '{"transaction_id": "TXN789012345", "amount": 25.50, "type": "expense", "category_id": 10, "description": "Movie tickets"}' "5. Creating a transaction"

# 6. Get all transactions
api_call "GET" "http://localhost:8080/api/transactions" "" "6. Getting all transactions"

# 7. Get the newly created transaction (assuming it got ID 2)
api_call "GET" "http://localhost:8080/api/transactions/2" "" "7. Getting specific transaction"

# 8. Update the transaction
api_call "PUT" "http://localhost:8080/api/transactions/2" '{"transaction_id": "TXN789012345-UPD", "amount": 30.00, "description": "Movie tickets and popcorn"}' "8. Updating transaction"

# 9. Test date range filtering
api_call "GET" "http://localhost:8080/api/transactions/date-range?start_date=2025-08-01&end_date=2025-08-31" "" "9. Getting transactions by date range (August 2025)"

# 10. Get aggregated data
api_call "GET" "http://localhost:8080/api/transactions/aggregate" "" "10. Getting aggregated data"

# 11. Get summary
api_call "GET" "http://localhost:8080/api/transactions/summary" "" "11. Getting summary overview"

# 12. Try to delete category with transactions (should fail)
api_call "DELETE" "http://localhost:8080/api/categories/10" "" "12. Trying to delete category with transactions (should fail)"

# 13. Delete the transaction first
api_call "DELETE" "http://localhost:8080/api/transactions/2" "" "13. Deleting transaction"

# 14. Now delete the category (should succeed)
api_call "DELETE" "http://localhost:8080/api/categories/10" "" "14. Deleting category (should succeed)"

echo ""
echo -e "${YELLOW}Final State:${NC}"
echo "============="

# 15. Final transactions list
api_call "GET" "http://localhost:8080/api/transactions" "" "15. Final transactions list"

# 16. Final categories list
api_call "GET" "http://localhost:8080/api/categories" "" "16. Final categories list"

echo ""
echo -e "${GREEN}🎉 Demo completed successfully!${NC}"
echo ""
echo "The API now supports full CRUD operations for both transactions and categories."
echo "Check the API.md file for complete documentation of all endpoints." 