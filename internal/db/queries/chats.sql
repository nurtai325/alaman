-- name: GetChats :many
SELECT ch.*, u.name AS user_name, u.phone AS user_phone, l.phone AS lead_phone FROM chats AS ch
INNER JOIN users u ON ch.user_id = u.id
INNER JOIN leads l ON ch.lead_id = l.id
ORDER BY updated_at DESC 
LIMIT $2 
OFFSET $1;

-- name: GetChat :one
SELECT ch.*, u.name AS user_name, u.phone AS user_phone, l.phone AS lead_phone FROM chats AS ch
INNER JOIN users u ON ch.user_id = u.id
INNER JOIN leads l ON ch.lead_id = l.id
WHERE ch.id = $1 
LIMIT 1;

-- name: GetChatByLeadId :one
SELECT * FROM chats 
WHERE lead_id = $1 
LIMIT 1;

-- name: GetChatsCount :one
SELECT COUNT(*) 
FROM chats;

-- name: InsertChat :one
INSERT INTO chats(lead_id, user_id)
VALUES($1, $2)
RETURNING *;

-- name: UpdateChat :one
UPDATE chats
SET updated_at = $2
WHERE id = $1
RETURNING *;

-- name: DeleteChat :one
DELETE FROM chats
WHERE id = $1
RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages 
WHERE chat_id = $1
ORDER BY created_at ASC;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 
LIMIT 1;

-- name: GetMessagesCount :one
SELECT COUNT(*) 
FROM messages;

-- name: InsertMessage :one
INSERT INTO messages(text, path, type, is_sent, audio_length, chat_id)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteMessage :one
DELETE FROM messages
WHERE id = $1
RETURNING *;
