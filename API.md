# Expense API Documentation

## Base URL
- Local: `http://localhost:8080`
- Production: `https://your-domain.com`

## Authentication
No authentication required for this API.

## Endpoints

### Health Check

#### GET /health
Check if the API is running.

**Response:**
```json
{
  "status": "healthy",
  "message": "Expense API is running"
}
```

### Transactions

#### POST /api/transactions
Create a new transaction.

**Request Body:**
```json
{
  "transaction_id": "TXN123456789",
  "amount": 50.0,
  "type": "expense",
  "category_id": 1,
  "description": "Lunch",
  "date": "2024-01-15T12:00:00Z"
}
```

**Fields:**
- `transaction_id` (string, optional): Bank transaction reference number
- `amount` (float, required): Transaction amount
- `type` (string, required): Either "expense" or "income"
- `category_id` (integer, required): ID of the category
- `description` (string, required): Transaction description
- `date` (string, optional): ISO 8601 date string (defaults to current time)

**Response (201 Created):**
```json
{
  "id": 1,
  "transaction_id": "TXN123456789",
  "amount": 50.0,
  "type": "expense",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Food",
    "type": "expense"
  },
  "description": "Lunch",
  "date": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data or category not found
- `500 Internal Server Error`: Database error

#### POST /api/transactions/bulk
Create multiple transactions in a single request.

**Request Body:**
```json
{
  "transactions": [
    {
      "transaction_id": "TXN123456789",
      "amount": 50.0,
      "type": "expense",
      "category_id": 1,
      "description": "Lunch",
      "date": "2024-01-15T12:00:00Z"
    },
    {
      "transaction_id": "TXN987654321",
      "amount": 25.0,
      "type": "expense",
      "category_id": 2,
      "description": "Coffee"
    },
    {
      "transaction_id": "TXN555666777",
      "amount": 1000.0,
      "type": "income",
      "category_id": 3,
      "description": "Freelance payment"
    }
  ]
}
```

**Fields:**
- `transactions` (array, required): Array of transaction objects (max 100)
  - Each transaction object has the same fields as the single transaction POST endpoint

**Response (201 Created - All Success):**
```json
{
  "success": [
    {
      "id": 1,
      "amount": 50.0,
      "type": "expense",
      "category_id": 1,
      "category": "Food",
      "description": "Lunch",
      "date": "2024-01-15T12:00:00Z",
      "created_at": "2024-01-15T12:00:00Z"
    }
  ],
  "failed": [],
  "total_count": 3,
  "success_count": 3,
  "failed_count": 0
}
```

**Response (207 Multi-Status - Partial Success):**
```json
{
  "success": [
    {
      "id": 1,
      "amount": 50.0,
      "type": "expense",
      "category_id": 1,
      "category": "Food",
      "description": "Lunch",
      "date": "2024-01-15T12:00:00Z",
      "created_at": "2024-01-15T12:00:00Z"
    }
  ],
  "failed": [
    {
      "index": 1,
      "transaction": {
        "amount": 25.0,
        "type": "expense",
        "category_id": 999,
        "description": "Coffee"
      },
      "error": "Category not found"
    }
  ],
  "total_count": 2,
  "success_count": 1,
  "failed_count": 1
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body, no transactions provided, or all transactions failed
- `207 Multi-Status`: Some transactions succeeded, some failed
- `500 Internal Server Error`: Database error

#### GET /api/transactions
Get all transactions with optional filtering.

**Query Parameters:**
- `type` (optional): Filter by "expense" or "income"

**Examples:**
- `GET /api/transactions` - Get all transactions
- `GET /api/transactions?type=expense` - Get only expenses
- `GET /api/transactions?type=income` - Get only income

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "transaction_id": "TXN123456789",
    "amount": 50.0,
    "type": "expense",
    "category_id": 1,
    "category": "Food",
    "description": "Lunch",
    "date": "2024-01-15T12:00:00Z",
    "created_at": "2024-01-15T12:00:00Z"
  }
]
```

#### GET /api/transactions/aggregate
Get aggregated transaction data by category.

**Response (200 OK):**
```json
{
  "categories": {
    "Food": 150.0,
    "Transport": 75.0,
    "Salary": 5000.0
  },
  "total_income": 5000.0,
  "total_expenses": 225.0,
  "net_amount": 4775.0
}
```

#### GET /api/transactions/:id
Get a specific transaction by ID.

**Response (200 OK):**
```json
{
  "id": 1,
  "transaction_id": "TXN123456789",
  "amount": 50.0,
  "type": "expense",
  "category_id": 1,
  "category": "Food",
  "description": "Lunch",
  "date": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Transaction not found

#### PUT /api/transactions/:id
Update an existing transaction.

**Request Body:**
```json
{
  "amount": 55.0,
  "description": "Lunch with dessert"
}
```

**Fields (all optional):**
- `transaction_id` (string): Bank transaction reference number
- `amount` (float): Transaction amount
- `type` (string): Either "expense" or "income"
- `category_id` (integer): ID of the category
- `description` (string): Transaction description
- `date` (string): ISO 8601 date string

**Response (200 OK):**
```json
{
  "id": 1,
  "amount": 55.0,
  "type": "expense",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Food",
    "type": "expense"
  },
  "description": "Lunch with dessert",
  "date": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data or category not found
- `404 Not Found`: Transaction not found
- `500 Internal Server Error`: Database error

#### PATCH /api/transactions/:id/category
Update the category of an existing transaction.

**Request Body:**
```json
{
  "category_id": 2
}
```

**Fields:**
- `category_id` (integer, required): ID of the new category

**Response (200 OK):**
```json
{
  "id": 1,
  "transaction_id": "TXN123456789",
  "amount": 50.0,
  "type": "expense",
  "category_id": 2,
  "category": "Transport",
  "description": "Lunch",
  "date": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data, category not found, or category type mismatch
- `404 Not Found`: Transaction not found
- `500 Internal Server Error`: Database error

#### DELETE /api/transactions/:id
Delete a transaction.

**Response (200 OK):**
```json
{
  "message": "Transaction deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Transaction not found
- `500 Internal Server Error`: Database error

#### DELETE /api/transactions/bulk
Delete multiple transactions in a single request.

**Request Body:**
```json
{
  "transaction_ids": [1, 2, 3, 4, 5]
}
```

**Fields:**
- `transaction_ids` (array, required): Array of transaction IDs to delete (max 1000)

**Response (200 OK - All Success):**
```json
{
  "deleted": [1, 2, 3, 4, 5],
  "failed": [],
  "total_count": 5,
  "deleted_count": 5,
  "failed_count": 0
}
```

**Response (207 Multi-Status - Partial Success):**
```json
{
  "deleted": [1, 2, 4],
  "failed": [
    {
      "transaction_id": 3,
      "error": "Transaction not found"
    },
    {
      "transaction_id": 5,
      "error": "Transaction not found"
    }
  ],
  "total_count": 5,
  "deleted_count": 3,
  "failed_count": 2
}
```

**Response (400 Bad Request - All Failed):**
```json
{
  "deleted": [],
  "failed": [
    {
      "transaction_id": 1,
      "error": "Transaction not found"
    },
    {
      "transaction_id": 2,
      "error": "Transaction not found"
    }
  ],
  "total_count": 2,
  "deleted_count": 0,
  "failed_count": 2
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body, no transaction IDs provided, too many IDs (>1000), or all deletions failed
- `207 Multi-Status`: Some deletions succeeded, some failed
- `500 Internal Server Error`: Database error

#### GET /api/transactions/date-range
Get transactions within a specific date range.

**Query Parameters:**
- `start_date` (required): Start date in YYYY-MM-DD format
- `end_date` (required): End date in YYYY-MM-DD format
- `type` (optional): Filter by "expense" or "income"

**Examples:**
- `GET /api/transactions/date-range?start_date=2024-01-01&end_date=2024-01-31`
- `GET /api/transactions/date-range?start_date=2024-01-01&end_date=2024-01-31&type=expense`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "transaction_id": "TXN123456789",
    "amount": 50.0,
    "type": "expense",
    "category_id": 1,
    "category": "Food",
    "description": "Lunch",
    "date": "2024-01-15T12:00:00Z",
    "created_at": "2024-01-15T12:00:00Z"
  }
]
```

**Error Responses:**
- `400 Bad Request`: Missing or invalid date parameters
- `500 Internal Server Error`: Database error

### Categories

#### POST /api/categories
Create a new category.

**Request Body:**
```json
{
  "name": "Entertainment",
  "type": "expense"
}
```

**Fields:**
- `name` (string, required): Category name (must be unique)
- `type` (string, required): Either "expense" or "income"

**Response (201 Created):**
```json
{
  "id": 8,
  "name": "Entertainment",
  "type": "expense",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data or duplicate category name
- `500 Internal Server Error`: Database error

#### GET /api/categories
Get all categories.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Food",
    "type": "expense"
  },
  {
    "id": 2,
    "name": "Transport",
    "type": "expense"
  },
  {
    "id": 5,
    "name": "Salary",
    "type": "income"
  }
]
```

#### GET /api/categories/:id
Get a specific category by ID.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Food",
  "type": "expense"
}
```

**Error Responses:**
- `404 Not Found`: Category not found

#### PUT /api/categories/:id
Update an existing category.

**Request Body:**
```json
{
  "name": "Food & Dining",
  "type": "expense"
}
```

**Fields (all optional):**
- `name` (string): Category name (must be unique)
- `type` (string): Either "expense" or "income"

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Food & Dining",
  "type": "expense",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z",
  "deleted_at": null
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data or duplicate category name
- `404 Not Found`: Category not found
- `500 Internal Server Error`: Database error

#### DELETE /api/categories/:id
Delete a category.

**Response (200 OK):**
```json
{
  "message": "Category deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Category has associated transactions
- `404 Not Found`: Category not found
- `500 Internal Server Error`: Database error

## Error Handling

All error responses follow this format:
```json
{
  "error": "Error message description"
}
```

## HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Data Validation

### Transaction Type
Must be exactly "expense" or "income" (case-sensitive).

### Category Type
Must be exactly "expense" or "income" (case-sensitive).

### Amount
Must be a positive number.

### Category ID
Must reference an existing category, and the category type must match the transaction type.

### Category Name
Must be unique across all categories.

## CORS

The API supports CORS for cross-origin requests:
- All origins allowed (`*`)
- Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Headers: Origin, Content-Type, Accept, Authorization

## Rate Limiting

Currently no rate limiting is implemented.

## Pagination

Currently no pagination is implemented. All endpoints return all results.

## Examples

### Complete Workflow

1. **Check API health:**
   ```bash
   curl http://localhost:8080/health
   ```

2. **View default categories:**
   ```bash
   curl http://localhost:8080/api/categories
   ```

3. **Create a transaction:**
   ```bash
   curl -X POST http://localhost:8080/api/transactions \
     -H "Content-Type: application/json" \
     -d '{
       "transaction_id": "TXN789012345",
       "amount": 25.50,
       "type": "expense",
       "category_id": 1,
       "description": "Coffee and pastry"
     }'
   ```

4. **View all transactions:**
   ```bash
   curl http://localhost:8080/api/transactions
   ```

5. **Get aggregated data:**
   ```bash
   curl http://localhost:8080/api/transactions/aggregate
   ```

6. **Get a specific transaction:**
   ```bash
   curl http://localhost:8080/api/transactions/1
   ```

7. **Update a transaction:**
   ```bash
   curl -X PUT http://localhost:8080/api/transactions/1 \
     -H "Content-Type: application/json" \
     -d '{
       "amount": 55.0,
       "description": "Lunch with dessert"
     }'
   ```

8. **Change transaction category:**
   ```bash
   curl -X PATCH http://localhost:8080/api/transactions/1/category \
     -H "Content-Type: application/json" \
     -d '{
       "category_id": 2
     }'
   ```

9. **Delete a transaction:**
   ```bash
   curl -X DELETE http://localhost:8080/api/transactions/1
   ```

10. **Delete multiple transactions (bulk delete):**
    ```bash
    curl -X DELETE http://localhost:8080/api/transactions/bulk \
      -H "Content-Type: application/json" \
      -d '{
        "transaction_ids": [1, 2, 3, 4, 5]
      }'
    ```

11. **Get transactions by date range:**
    ```bash
    curl "http://localhost:8080/api/transactions/date-range?start_date=2024-01-01&end_date=2024-01-31"
    ```

12. **Get a specific category:**
    ```bash
    curl http://localhost:8080/api/categories/1
    ```

13. **Update a category:**
    ```bash
    curl -X PUT http://localhost:8080/api/categories/1 \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Food & Dining"
      }'
    ```

14. **Delete a category (only if no transactions exist):**
    ```bash
    curl -X DELETE http://localhost:8080/api/categories/1
    ```

### Using with JavaScript/Fetch

```javascript
// Create a transaction
const response = await fetch('http://localhost:8080/api/transactions', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    transaction_id: 'TXN123456789',
    amount: 50.0,
    type: 'expense',
    category_id: 1,
    description: 'Lunch'
  })
});

const transaction = await response.json();
console.log('Created transaction:', transaction);
``` 