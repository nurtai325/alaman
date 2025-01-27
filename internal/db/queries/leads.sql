-- name: GetNewLeads :many
SELECT * FROM leads AS l
WHERE user_id = NULL
ORDER BY created_at DESC;

-- name: GetAssignedLeads :many
SELECT * FROM leads
WHERE user_id != NULL AND sale_id = NULL
ORDER BY created_at DESC;

-- name: GetInDeliveryLeads :many
SELECT * FROM leads
WHERE user_id != NULL AND sale_id != NULL AND completed = false
ORDER BY created_at DESC;

-- name: GetCompletedLeads :many
SELECT * FROM leads
WHERE completed = true
ORDER BY created_at DESC;
