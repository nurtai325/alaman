-- name: GetUsers :many
SELECT * FROM users WHERE id BETWEEN $1 and $2;;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: InsertUser :exec
INSERT INTO users(name, phone, password) VALUES($1, $2, $3);
