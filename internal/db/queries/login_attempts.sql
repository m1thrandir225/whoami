-- name: CreateLoginAttempt :one
INSERT INTO login_attempts (
    user_id,
    email,
    ip_address,
    user_agent,
    success,
    failure_reason
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetLoginAttemptsByUserID :many
SELECT * FROM login_attempts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetLoginAttemptsByEmail :many
SELECT * FROM login_attempts
WHERE email = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetLoginAttemptsByIP :many
SELECT * FROM login_attempts
WHERE ip_address = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetFailedLoginAttemptsByUserID :many
SELECT * FROM login_attempts
WHERE user_id = $1 AND success = false
ORDER BY created_at DESC
LIMIT $2;

-- name: GetFailedLoginAttemptsByEmail :many
SELECT * FROM login_attempts
WHERE email = $1 AND success = false
ORDER BY created_at DESC
LIMIT $2;

-- name: GetFailedLoginAttemptsByIP :many
SELECT * FROM login_attempts
WHERE ip_address = $1 AND success = false
ORDER BY created_at DESC
LIMIT $2;

-- name: GetRecentFailedAttemptsByUserID :many
SELECT * FROM login_attempts
WHERE user_id = $1
AND success = false
AND created_at > NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC;

-- name: GetRecentFailedAttemptsByEmail :many
SELECT * FROM login_attempts
WHERE email = $1
AND success = false
AND created_at > NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC;

-- name: GetRecentFailedAttemptsByIP :many
SELECT * FROM login_attempts
WHERE ip_address = $1
AND success = false
AND created_at > NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC;

-- name: DeleteOldLoginAttempts :exec
DELETE FROM login_attempts
WHERE created_at < NOW() - INTERVAL '90 days';
