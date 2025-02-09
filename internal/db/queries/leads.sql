-- name: GetNewLeads :many
SELECT * FROM leads AS l
WHERE user_id IS NULL
ORDER BY created_at DESC;

-- name: GetLeadByPhone :one
SELECT * FROM leads 
WHERE phone = $1 
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAssignedLeads :many
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL
ORDER BY created_at DESC;

-- name: GetAssignedLeadsByUser :many
SELECT l.*, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL AND user_id = $1
ORDER BY created_at DESC;

-- name: GetInDeliveryLeads :many
SELECT l.*, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false
ORDER BY sold_at DESC;

-- name: GetInDeliveryLeadsByUser :many
SELECT l.*, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false AND user_id = $1
ORDER BY sold_at DESC;

-- name: GetCompletedLeads :many
SELECT l.*, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true
ORDER BY sold_at DESC;

-- name: GetCompletedLeadsByUser :many
SELECT l.*, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true AND user_id = $1
ORDER BY sold_at DESC;

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
INSERT INTO sales(type, full_sum, delivery_cost, loan_cost, items_sum, delivery_type, payment_at)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: InsertSaleItem :one
INSERT INTO sale_items(price, product_name, sale_id, quantity, product_id, sale_count)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SellLead :one
UPDAte leads
SET sale_id = $2, sold_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetFullLead :one
SELECT l.*, u.name AS user_name, u.phone AS user_phone, s.full_sum, s.delivery_cost, s.loan_cost, s.delivery_type, s.payment_at, s.type AS sale_type
FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE l.id = $1
LIMIT 1;

-- name: GetSaleItems :many
SELECT * FROM sale_items AS s
WHERE s.sale_id = $1;

-- name: CompleteLead :one
UPDATE leads
SET completed = true
WHERE id = $1
RETURNING *;
