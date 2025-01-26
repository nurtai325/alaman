-- name: GetProducts :many
SELECT * FROM products 
ORDER BY created_at DESC 
LIMIT $2 
OFFSET $1;

-- name: GetProduct :one
SELECT * FROM products 
WHERE id = $1 
LIMIT 1;

-- name: GetProductByName :one
SELECT * FROM products
WHERE name = $1 
LIMIT 1;

-- name: GetProductsCount :one
SELECT COUNT(*) 
FROM products;

-- name: InsertProduct :one
INSERT INTO products(name, in_stock, price, stock_price)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, stock_price = $4
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

-- name: AddStockProduct :one
UPDATE products
SET in_stock = $2
WHERE id = $1
RETURNING in_stock;
