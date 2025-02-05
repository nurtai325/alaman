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
INSERT INTO products(name, in_stock, price, sale_count, stock_price)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, stock_price = $4, sale_count = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

-- name: AddStockProduct :one
UPDATE products
SET in_stock = in_stock + $2
WHERE id = $1
RETURNING in_stock;

-- name: RemoveStockProduct :one
UPDATE products
SET in_stock = in_stock - $2
WHERE id = $1
RETURNING in_stock;

-- name: GetProductChanges :many
SELECT * FROM product_changes
WHERE created_at > $1 AND created_at < $2
ORDER BY created_at ASC;

-- name: InsertProductChange :one
INSERT INTO product_changes(quantity, is_income, product_id)
VALUES($1 ,$2 , $3)
RETURNING *;
