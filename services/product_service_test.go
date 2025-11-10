package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ProductService) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	service := NewProductService(db)
	return db, mock, service
}

func TestNewProductService(t *testing.T) {
	db, _, _ := setupTestDB(t)
	defer db.Close()

	service := NewProductService(db)
	assert.NotNil(t, service)
	assert.NotNil(t, service.db)
}

func TestCreateProduct_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	input := CreateProductInput{
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-001",
		Price:       99.99,
		Stock:       10,
	}

	now := time.Now()
	expectedProduct := &Product{
		ID:          1,
		Name:        input.Name,
		Description: input.Description,
		SKU:         input.SKU,
		Price:       input.Price,
		Stock:       input.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect SKU check query
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products WHERE sku = \\$1").
		WithArgs(input.SKU).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Expect insert query
	mock.ExpectQuery("INSERT INTO products").
		WithArgs(input.Name, input.Description, input.SKU, input.Price, input.Stock, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "sku", "price", "stock", "created_at", "updated_at"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.SKU,
				expectedProduct.Price, expectedProduct.Stock, expectedProduct.CreatedAt, expectedProduct.UpdatedAt))

	// Expect commit
	mock.ExpectCommit()

	product, err := service.CreateProduct(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	assert.Equal(t, expectedProduct.SKU, product.SKU)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_DuplicateSKU(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	input := CreateProductInput{
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-001",
		Price:       99.99,
		Stock:       10,
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect SKU check query - returns 1 indicating SKU exists
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products WHERE sku = \\$1").
		WithArgs(input.SKU).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Expect rollback
	mock.ExpectRollback()

	product, err := service.CreateProduct(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, ErrProductAlreadyExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	db, _, service := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name  string
		input CreateProductInput
	}{
		{
			name: "Empty name",
			input: CreateProductInput{
				Name:  "",
				SKU:   "TEST-001",
				Price: 99.99,
				Stock: 10,
			},
		},
		{
			name: "Name too short",
			input: CreateProductInput{
				Name:  "A",
				SKU:   "TEST-001",
				Price: 99.99,
				Stock: 10,
			},
		},
		{
			name: "Empty SKU",
			input: CreateProductInput{
				Name:  "Test Product",
				SKU:   "",
				Price: 99.99,
				Stock: 10,
			},
		},
		{
			name: "Invalid price",
			input: CreateProductInput{
				Name:  "Test Product",
				SKU:   "TEST-001",
				Price: -10,
				Stock: 10,
			},
		},
		{
			name: "Negative stock",
			input: CreateProductInput{
				Name:  "Test Product",
				SKU:   "TEST-001",
				Price: 99.99,
				Stock: -5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.CreateProduct(context.Background(), tt.input)
			assert.Error(t, err)
			assert.Nil(t, product)
			assert.ErrorIs(t, err, ErrInvalidProductData)
		})
	}
}

func TestGetProductByID_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	expectedProduct := &Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-001",
		Price:       99.99,
		Stock:       10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery("SELECT id, name, description, sku, price, stock, created_at, updated_at FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "sku", "price", "stock", "created_at", "updated_at"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.SKU,
				expectedProduct.Price, expectedProduct.Stock, expectedProduct.CreatedAt, expectedProduct.UpdatedAt))

	product, err := service.GetProductByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductByID_NotFound(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, description, sku, price, stock, created_at, updated_at FROM products WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	product, err := service.GetProductByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductByID_InvalidID(t *testing.T) {
	db, _, service := setupTestDB(t)
	defer db.Close()

	product, err := service.GetProductByID(context.Background(), -1)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, ErrInvalidProductData)
}

func TestGetProductBySKU_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	expectedProduct := &Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-001",
		Price:       99.99,
		Stock:       10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery("SELECT id, name, description, sku, price, stock, created_at, updated_at FROM products WHERE sku = \\$1").
		WithArgs("TEST-001").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "sku", "price", "stock", "created_at", "updated_at"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.SKU,
				expectedProduct.Price, expectedProduct.Stock, expectedProduct.CreatedAt, expectedProduct.UpdatedAt))

	product, err := service.GetProductBySKU(context.Background(), "TEST-001")

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.SKU, product.SKU)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListProducts_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "sku", "price", "stock", "created_at", "updated_at"}).
		AddRow(1, "Product 1", "Description 1", "SKU-001", 99.99, 10, now, now).
		AddRow(2, "Product 2", "Description 2", "SKU-002", 149.99, 5, now, now)

	mock.ExpectQuery("SELECT id, name, description, sku, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
		WithArgs(50, 0).
		WillReturnRows(rows)

	products, err := service.ListProducts(context.Background(), 50, 0)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 2", products[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	newName := "Updated Product"
	newPrice := 149.99
	input := UpdateProductInput{
		Name:  &newName,
		Price: &newPrice,
	}

	now := time.Now()
	expectedProduct := &Product{
		ID:          1,
		Name:        newName,
		Description: "Original Description",
		SKU:         "TEST-001",
		Price:       newPrice,
		Stock:       10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect existence check
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM products WHERE id = \\$1\\)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Expect update query
	mock.ExpectExec("UPDATE products SET updated_at = \\$1, name = \\$2, price = \\$3 WHERE id = \\$4").
		WithArgs(sqlmock.AnyArg(), newName, newPrice, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect commit
	mock.ExpectCommit()

	// Expect get query after update
	mock.ExpectQuery("SELECT id, name, description, sku, price, stock, created_at, updated_at FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "sku", "price", "stock", "created_at", "updated_at"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.SKU,
				expectedProduct.Price, expectedProduct.Stock, expectedProduct.CreatedAt, expectedProduct.UpdatedAt))

	product, err := service.UpdateProduct(context.Background(), 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, newName, product.Name)
	assert.Equal(t, newPrice, product.Price)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_NotFound(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	newName := "Updated Product"
	input := UpdateProductInput{
		Name: &newName,
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect existence check - returns false
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM products WHERE id = \\$1\\)").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Expect rollback
	mock.ExpectRollback()

	product, err := service.UpdateProduct(context.Background(), 999, input)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct_Success(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct_NotFound(t *testing.T) {
	db, mock, service := setupTestDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := service.DeleteProduct(context.Background(), 999)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrProductNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct_InvalidID(t *testing.T) {
	db, _, service := setupTestDB(t)
	defer db.Close()

	err := service.DeleteProduct(context.Background(), -1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidProductData)
}

