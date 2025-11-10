# Go Gin GORM PostgreSQL Server

A basic "Hello World" server built with Go, Gin framework, GORM, and PostgreSQL.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Make sure PostgreSQL is running and accessible

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Create a PostgreSQL database:
```sql
CREATE DATABASE test_db;
```

3. Update database connection details in `main.go` if needed:
```go
dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable"
```

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Hello Endpoints
- `GET /api/hello` - Returns a "Hello World" message
- `GET /api/hello/:name` - Returns a personalized greeting

### User Endpoints
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get a specific user by ID
- `POST /api/users` - Create a new user

### Product Endpoints
- `GET /api/products` - List all products (supports pagination with `?limit=50&offset=0`)
- `GET /api/products/:id` - Get a specific product by ID
- `GET /api/products/sku/:sku` - Get a product by SKU
- `POST /api/products` - Create a new product
- `PUT /api/products/:id` - Update a product (full update)
- `PATCH /api/products/:id` - Partially update a product
- `DELETE /api/products/:id` - Delete a product

### System Endpoints
- `GET /api/health` - Health check endpoint
- `GET /api/version` - Get API version information

## Product API Examples

### Create a Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop for developers",
    "sku": "LAP-001",
    "price": 1299.99,
    "stock": 50
  }'
```

### Get All Products
```bash
curl http://localhost:8080/api/products?limit=10&offset=0
```

### Get Product by ID
```bash
curl http://localhost:8080/api/products/1
```

### Get Product by SKU
```bash
curl http://localhost:8080/api/products/sku/LAP-001
```

### Update a Product
```bash
curl -X PATCH http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 1199.99,
    "stock": 45
  }'
```

### Delete a Product
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

## Database Setup

The application requires PostgreSQL. Make sure to run the migrations in the `migrations/` directory:

```bash
# Create products table
psql -U postgres -d test_db -f migrations/001_create_products_table.sql
```

To roll back the migration:
```bash
psql -U postgres -d test_db -f migrations/001_create_products_table_down.sql
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```
 