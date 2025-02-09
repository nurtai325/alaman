-- name: GetReports :many
SELECT * FROM reports
ORDER BY created_at DESC;

-- name: GetReport :one
SELECT * FROM reports
WHERE id = $1
LIMIT 1;

-- name: InsertReport :one
INSERT INTO reports(name, path, start_at, end_at)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateReport :one
UPDATE reports
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteReport :one
DELETE FROM reports
WHERE id = $1
RETURNING *;

-- name: GetReportByProduct :one
SELECT 
	COALESCE(SUM(sl.sale_count), 0) AS sale_count_sum, 
	COALESCE(SUM(sl.quantity), 0) AS sold, 
	COALESCE(SUM(sl.price), 0) AS sold_sum
FROM sale_items AS sl
WHERE sl.product_id = $1 
AND sl.created_at > $2 
AND sl.created_at < $3;

-- name: GetProductIncoming :one
SELECT SUM(quantity)
FROM product_changes
WHERE product_id = $1 AND created_at > $2 AND created_at < $3 AND is_income = TRUE;

-- name: GetProductOutcoming :one
SELECT SUM(quantity)
FROM product_changes
WHERE product_id = $1 AND created_at > $2 AND created_at < $3 AND is_income = FALSE;
