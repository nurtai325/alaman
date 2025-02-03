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
