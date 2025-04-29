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

- `GET /`: Returns a "Hello World" message 