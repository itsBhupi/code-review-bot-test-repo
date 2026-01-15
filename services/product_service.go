package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("product not found")
	// ErrInvalidProductData is returned when product data is invalid
	ErrInvalidProductData = errors.New("invalid product data")
	// ErrProductAlreadyExists is returned when trying to create a duplicate product
	ErrProductAlreadyExists = errors.New("product with this SKU already exists")
)

// Product represents a product in the system
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	Name        string  `json:"name" binding:"required,min=2,max=200"`
	Description string  `json:"description" binding:"max=1000"`
	SKU         string  `json:"sku" binding:"required,min=3,max=50"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

// UpdateProductInput represents the input for updating a product
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

// ProductService handles product business logic
type ProductService struct {
	db *sql.DB
}

// NewProductService creates a new instance of ProductService
func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product in the database
func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*Product, error) {
	// Validate input
	if err := s.validateProductInput(input); err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	defer tx.Rollback()

	// Check if SKU already exists
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM products WHERE sku = $1", input.SKU).Scan(&count)
	if err != nil {
		log.Error().Err(err).Str("sku", input.SKU).Msg("Failed to check SKU existence")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	if count > 0 {
		return nil, ErrProductAlreadyExists
	}

	// Insert product
	now := time.Now()
	var product Product
	err = tx.QueryRowContext(ctx,
		`INSERT INTO products (name, description, sku, price, stock, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) 
		 RETURNING id, name, description, sku, price, stock, created_at, updated_at`,
		input.Name,
		input.Description,
		input.SKU,
		input.Price,
		input.Stock,
		now,
		now,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to insert product")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	log.Info().Int("productID", product.ID).Str("sku", product.SKU).Msg("Product created successfully")
	return &product, nil
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id int) (*Product, error) {
	if id <= 0 {
		return nil, ErrInvalidProductData
	}

	var product Product
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, description, sku, price, stock, created_at, updated_at 
		 FROM products WHERE id = $1`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		log.Error().Err(err).Int("productID", id).Msg("Failed to get product")
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*Product, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, ErrInvalidProductData
	}

	var product Product
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, description, sku, price, stock, created_at, updated_at 
		 FROM products WHERE sku = $1`,
		sku,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		log.Error().Err(err).Str("sku", sku).Msg("Failed to get product by SKU")
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// ListProducts retrieves all products with optional pagination
func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]Product, error) {
	// Set default and max limits
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, description, sku, price, stock, created_at, updated_at 
		 FROM products 
		 ORDER BY created_at DESC 
		 LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query products")
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.SKU,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan product row")
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating product rows")
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id int, input UpdateProductInput) (*Product, error) {
	if id <= 0 {
		return nil, ErrInvalidProductData
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	defer tx.Rollback()

	// Check if product exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		log.Error().Err(err).Int("productID", id).Msg("Failed to check product existence")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	if !exists {
		return nil, ErrProductNotFound
	}

	// Build dynamic update query
	query := "UPDATE products SET updated_at = $1"
	args := []interface{}{time.Now()}
	paramCount := 1

	if input.Name != nil {
		paramCount++
		query += fmt.Sprintf(", name = $%d", paramCount)
		args = append(args, *input.Name)
	}

	if input.Description != nil {
		paramCount++
		query += fmt.Sprintf(", description = $%d", paramCount)
		args = append(args, *input.Description)
	}

	if input.Price != nil {
		if *input.Price <= 0 {
			return nil, ErrInvalidProductData
		}
		paramCount++
		query += fmt.Sprintf(", price = $%d", paramCount)
		args = append(args, *input.Price)
	}

	if input.Stock != nil {
		if *input.Stock < 0 {
			return nil, ErrInvalidProductData
		}
		paramCount++
		query += fmt.Sprintf(", stock = $%d", paramCount)
		args = append(args, *input.Stock)
	}

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	args = append(args, id)

	// Execute update
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Int("productID", id).Msg("Failed to update product")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Retrieve updated product
	product, err := s.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Info().Int("productID", id).Msg("Product updated successfully")
	return product, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidProductData
	}

	result, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Error().Err(err).Int("productID", id).Msg("Failed to delete product")
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows affected")
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	log.Info().Int("productID", id).Msg("Product deleted successfully")
	return nil
}

// validateProductInput validates product input data
func (s *ProductService) validateProductInput(input CreateProductInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidProductData)
	}

	if len(input.Name) < 2 || len(input.Name) > 200 {
		return fmt.Errorf("%w: name must be between 2 and 200 characters", ErrInvalidProductData)
	}

	if strings.TrimSpace(input.SKU) == "" {
		return fmt.Errorf("%w: SKU is required", ErrInvalidProductData)
	}

	if len(input.SKU) < 3 || len(input.SKU) > 50 {
		return fmt.Errorf("%w: SKU must be between 3 and 50 characters", ErrInvalidProductData)
	}

	if input.Price <= 0 {
		return fmt.Errorf("%w: price must be greater than 0", ErrInvalidProductData)
	}

	if input.Stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrInvalidProductData)
	}

	if len(input.Description) > 1000 {
		return fmt.Errorf("%w: description cannot exceed 1000 characters", ErrInvalidProductData)
	}

	return nil
}

