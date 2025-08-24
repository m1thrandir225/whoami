-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id,
    action,
    resource_type,
    resource_id,
    ip_address,
    user_agent,
    details
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAuditLogsByUserID :many
SELECT * FROM audit_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetAuditLogsByAction :many
SELECT * FROM audit_logs
WHERE action = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetAuditLogsByResourceType :many
SELECT * FROM audit_logs
WHERE resource_type = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetAuditLogsByResourceID :many
SELECT * FROM audit_logs
WHERE resource_type = $1 AND resource_id = $2
ORDER BY created_at DESC
LIMIT $3;

-- name: GetAuditLogsByIP :many
SELECT * FROM audit_logs
WHERE ip_address = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetAuditLogsByDateRange :many
SELECT * FROM audit_logs
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC
LIMIT $3;

-- name: GetRecentAuditLogs :many
SELECT * FROM audit_logs
WHERE created_at > NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC
LIMIT $1;

-- name: DeleteOldAuditLogs :exec
DELETE FROM audit_logs
WHERE created_at < NOW() - INTERVAL '90 days';
