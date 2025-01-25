-- name: GetUsers :many
SELECT * FROM users 
WHERE active = true
ORDER BY created_at DESC 
LIMIT $2 
OFFSET $1;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 
LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users 
WHERE phone = $1 
LIMIT 1;

-- name: GetUsersCount :one
SELECT COUNT(*) 
FROM users;

-- name: InsertUser :one
INSERT INTO users(name, phone, password, role)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, phone = $3, role = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;
