-- name: GetLeadWhs :many
SELECT * FROM lead_whs 
ORDER BY created_at DESC 
LIMIT $2 
OFFSET $1;

-- name: GetLeadWh :one
SELECT * FROM lead_whs 
WHERE id = $1 
LIMIT 1;

-- name: GetLeadWhsCount :one
SELECT COUNT(*) 
FROM lead_whs;

-- name: InsertLeadWh :one
INSERT INTO lead_whs(name, phone)
VALUES($1, $2)
RETURNING *;

-- name: UpdateLeadWh :one
UPDATE lead_whs 
SET name = $2, phone = $3
WHERE id = $1
RETURNING *;

-- name: ConnectLeadWh :one
UPDATE lead_whs 
SET jid = $2
WHERE id = $1
RETURNING *;

-- name: DeleteLeadWh :one
DELETE FROM lead_whs 
WHERE id = $1
RETURNING *;
