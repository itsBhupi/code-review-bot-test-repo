-- Migration rollback: Drop products table
-- Description: Removes the products table and all associated indexes
-- Date: 2025-11-10

-- Drop indexes first
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_sku;

-- Drop the table
DROP TABLE IF EXISTS products;

