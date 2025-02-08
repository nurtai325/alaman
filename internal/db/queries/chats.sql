-- name: GetChats :many
SELECT * FROM chats 
ORDER BY created_at DESC 
LIMIT $2 
OFFSET $1;

-- name: GetChat :one
SELECT * FROM chats 
WHERE id = $1 
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

-- name: DeleteChat :one
DELETE FROM chats
WHERE id = $1
RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages 
WHERE chat_id = $1
ORDER BY created_at DESC;

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
