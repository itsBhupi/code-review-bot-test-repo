-- Migration: Create products table
-- Description: Creates the products table with all necessary columns and indexes
-- Date: 2025-11-10

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    sku VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_created_at ON products(created_at DESC);
CREATE INDEX idx_products_name ON products(name);

-- Add comments for documentation
COMMENT ON TABLE products IS 'Stores product information for the e-commerce system';
COMMENT ON COLUMN products.id IS 'Primary key, auto-incremented';
COMMENT ON COLUMN products.name IS 'Product name (2-200 characters)';
COMMENT ON COLUMN products.description IS 'Product description (max 1000 characters)';
COMMENT ON COLUMN products.sku IS 'Stock Keeping Unit, unique identifier for product';
COMMENT ON COLUMN products.price IS 'Product price, must be greater than 0';
COMMENT ON COLUMN products.stock IS 'Current stock quantity, cannot be negative';
COMMENT ON COLUMN products.created_at IS 'Timestamp when product was created';
COMMENT ON COLUMN products.updated_at IS 'Timestamp when product was last updated';

