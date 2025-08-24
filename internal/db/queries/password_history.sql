-- name: CreatePasswordHistory :one
INSERT INTO password_history (
    user_id,
    password_hash,
    created_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPasswordHistoryByUserID :many
SELECT * FROM password_history
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: DeleteOldPasswordHistory :exec
DELETE FROM password_history
WHERE user_id = $1
AND created_at < NOW() - INTERVAL '1 year';

-- name: CheckPasswordInHistory :one
SELECT COUNT(*) FROM password_history
WHERE user_id = $1
AND password_hash = $2;
