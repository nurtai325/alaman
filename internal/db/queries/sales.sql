-- name: GetSum :one
SELECT SUM(items_sum) FROM sales AS s
INNER JOIN leads l ON l.sale_id = s.id
WHERE s.payment_at >= $1;

-- name: GetLeadCount :one
SELECT COUNT(*) FROM leads
WHERE created_at >= $1;

-- name: GetSoldLeadCount :one
SELECT COUNT(*) FROM leads
WHERE created_at >= $1 AND sale_id IS NOT NULL;

-- name: GetSales :many
SELECT s.*, u.id AS user_id, u.name AS user_name FROM sales AS s
INNER JOIN leads l ON l.sale_id = s.id
INNER JOIN users u ON l.user_id = u.id
WHERE s.payment_at >= $1;

-- name: GetSaleItemsByTime :many
SELECT si.* FROM sale_items AS si
INNER JOIN sales s ON si.sale_id = s.id
WHERE s.payment_at >= $1;

-- name: GetNewLeadsCountByTime :one
SELECT COUNT(*) FROM leads
WHERE created_at >= $1;
