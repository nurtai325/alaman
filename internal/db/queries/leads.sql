-- name: GetNewLeads :many
SELECT * FROM leads AS l
WHERE user_id IS NULL
ORDER BY created_at DESC;

-- name: GetAssignedLeads :many
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL
ORDER BY created_at DESC;

-- name: GetInDeliveryLeads :many
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false
ORDER BY sold_at ASC;

-- name: GetCompletedLeads :many
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true
ORDER BY sold_at ASC;

-- name: InsertLead :one
INSERT INTO leads(phone)
VALUES ($1)
RETURNING *;

-- name: AssignLead :one
UPDATE leads
SET user_id = $2
WHERE id = $1
RETURNING *;

-- name: SetLeadInfo :one
UPDATE leads
SET name = $2, address = $3
WHERE id = $1
RETURNING *;

-- name: InsertSale :one
INSERT INTO sales(type, full_sum, delivery_cost, loan_cost, items_sum)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertSaleItem :one
INSERT INTO sale_items(product_id, sale_id, quantity)
VALUES($1, $2, $3)
RETURNING *;

-- name: SellLead :one
UPDAte leads
SET sale_id = $2, sold_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetFullLead :one
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE l.id = $1
LIMIT 1;

-- name: GetSaleItems :many
SELECT s.*, p.name AS product_name FROM sale_items AS s
INNER JOIN products p ON s.product_id = p.id
WHERE s.sale_id = $1;

-- name: CompleteLead :one
UPDAte leads
SET completed = true
WHERE id = $1
RETURNING *;
