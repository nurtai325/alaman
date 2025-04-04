-- name: GetUsers :many
SELECT * FROM users 
WHERE deleted = FALSE
ORDER BY created_at DESC 
LIMIT $2
OFFSET $1;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 AND deleted = FALSE
LIMIT 1;

-- name: GetLogist :one
SELECT * FROM users 
WHERE role = 'логист' AND deleted = FALSE
LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users 
WHERE phone = $1 AND deleted = FALSE
LIMIT 1;

-- name: GetUsersCount :one
SELECT COUNT(*) 
FROM users;

-- name: InsertUser :one
INSERT INTO users(name, phone, password, role, jid)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, phone = $3, role = $4
WHERE id = $1 AND deleted = FALSE
RETURNING *;

-- name: DeleteUser :one
UPDATE users
SET deleted = TRUE
WHERE id = $1
RETURNING *;

-- name: ConnectUser :one
UPDATE users
SET jid = $2
WHERE id = $1 AND deleted = FALSE
RETURNING *;
