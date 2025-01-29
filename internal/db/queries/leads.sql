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
SELECT * FROM leads
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false
ORDER BY created_at DESC;

-- name: GetCompletedLeads :many
SELECT * FROM leads
WHERE completed = true
ORDER BY created_at DESC;

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
SET sale_id = $2
WHERE id = $1
RETURNING *;

-- name: GetFullLead :one
SELECT l.* FROM leads AS l
WHERE id = $1
LIMIT 1;
