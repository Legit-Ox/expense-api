# Expense API

A RESTful API for managing personal expenses and income built with Go, Fiber, and GORM. Supports both PostgreSQL and SQLite databases.

## Features

- **Transaction Management**: Create, read, update, and delete expense/income transactions
- **Category Management**: Organize transactions by categories with full CRUD operations
- **Aggregation**: Get summaries by category and totals
- **Date Range Filtering**: Filter transactions by specific date ranges
- **Flexible Database**: Support for PostgreSQL (production) and SQLite (development)
- **CORS Enabled**: Ready for web and mobile applications
- **Docker Support**: Easy deployment and development setup
- **Data Validation**: Comprehensive validation for data integrity
- **Soft Deletes**: Categories support soft deletion with referential integrity checks

## Enhanced API Capabilities

The API now provides full CRUD (Create, Read, Update, Delete) operations for both transactions and categories:

### Transaction Operations
- **Create**: Add new expense or income transactions
- **Read**: Retrieve all transactions, specific transactions, or filter by type/date range
- **Update**: Modify existing transaction details
- **Delete**: Remove transactions from the system
- **Aggregate**: Get financial summaries and category breakdowns

### Category Operations
- **Create**: Add new expense or income categories
- **Read**: List all categories or get specific category details
- **Update**: Modify category names and types
- **Delete**: Remove categories (only if no transactions reference them)

### Advanced Features
- **Date Range Filtering**: Query transactions within specific time periods
- **Type Filtering**: Filter transactions by expense or income type
- **Referential Integrity**: Prevent deletion of categories with associated transactions
- **Data Validation**: Comprehensive input validation and error handling

## Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) - Fast HTTP framework for Go
- **ORM**: [GORM](https://gorm.io/) - Go ORM library
- **Database**: PostgreSQL (production) / SQLite (development)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing with testify assertions

## API Endpoints

### Transactions
- `POST /api/transactions` - Create a new transaction
- `GET /api/transactions` - List all transactions (with optional type filter)
- `GET /api/transactions/:id` - Get a specific transaction
- `PUT /api/transactions/:id` - Update a transaction
- `DELETE /api/transactions/:id` - Delete a transaction
- `GET /api/transactions/aggregate` - Get aggregated data by category
- `GET /api/transactions/date-range` - Get transactions within a date range

### Categories
- `POST /api/categories` - Create a new category
- `GET /api/categories` - List all categories
- `GET /api/categories/:id` - Get a specific category
- `PUT /api/categories/:id` - Update a category
- `DELETE /api/categories/:id` - Delete a category (only if no transactions exist)

### Health Check
- `GET /health` - API health status

## Data Models

### Transaction
```json
{
  "id": 1,
  "transaction_id": "TXN123456789",
  "amount": 50.0,
  "type": "expense",
  "category_id": 1,
  "description": "Lunch",
  "date": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z"
}
```

### Category
```json
{
  "id": 1,
  "name": "Food",
  "type": "expense"
}
```

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose (for local development)
- PostgreSQL (for production)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd expense-api
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your database configuration
   ```

3. **Using Docker Compose (recommended)**
   ```bash
   docker-compose up -d
   ```
   This will start both PostgreSQL and the API service.

4. **Using Go directly**
   ```bash
   # Download dependencies
   go mod tidy
   
   # Set environment variables
   export DB_URL="sqlite://./expense.db"
   export PORT=8080
   
   # Run the application
   go run main.go
   ```

### API Usage Examples

#### Create a Transaction
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_id": "TXN123456789",
    "amount": 50.0,
    "type": "expense",
    "category_id": 1,
    "description": "Lunch",
    "date": "2024-01-15T12:00:00Z"
  }'
```

#### Get All Transactions
```bash
curl http://localhost:8080/api/transactions
```

#### Filter Transactions by Type
```bash
curl "http://localhost:8080/api/transactions?type=expense"
```

#### Get Aggregated Data
```bash
curl http://localhost:8080/api/transactions/aggregate
```

#### Create a Category
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Entertainment",
    "type": "expense"
  }'
```

## Testing

Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Database Configuration

### Environment Variables
- `DB_URL`: Database connection string
  - PostgreSQL: `postgres://user:password@host:port/dbname?sslmode=disable`
  - SQLite: `sqlite://./expense.db`
- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development/production)

### Default Categories
The API automatically creates these default categories:

**Expenses:**
- Food
- Transport
- Bills
- Shopping

**Income:**
- Salary
- Freelance
- Investments

## Deployment

### Render (Recommended for Free Tier)

1. Connect your GitHub repository to Render
2. Create a new Web Service
3. Use the provided `render.yaml` configuration
4. Render will automatically deploy your API

### Railway

1. Install Railway CLI: `npm i -g @railway/cli`
2. Login: `railway login`
3. Deploy: `railway up`

### Manual Docker Deployment

```bash
# Build the image
docker build -t expense-api .

# Run the container
docker run -p 8080:8080 \
  -e DB_URL="your-database-url" \
  -e PORT=8080 \
  expense-api
```

## Development

### Project Structure
```
expense-api/
├── models/          # Data models and structs
├── database/        # Database connection and migrations
├── handlers/        # HTTP request handlers
├── main.go         # Application entry point
├── Dockerfile      # Container configuration
├── docker-compose.yml # Local development setup
├── render.yaml     # Render deployment config
├── railway.json    # Railway deployment config
└── README.md       # This file
```

### Adding New Endpoints

1. Create handler functions in the appropriate handler file
2. Add routes in `main.go`
3. Update tests if needed

### Database Migrations

The application uses GORM's auto-migration feature. To add new fields:

1. Update the model structs in `models/models.go`
2. Restart the application - GORM will automatically migrate the schema

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For issues and questions:
- Create an issue in the GitHub repository
- Check the API documentation at `/health` endpoint
- Review the test files for usage examples 