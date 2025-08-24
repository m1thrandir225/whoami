-- name: CreatePasswordReset :one
INSERT INTO password_resets (
    user_id,
    token_hash,
    hotp_secret,
    counter,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets
WHERE token_hash = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkPasswordResetAsUsed :exec
UPDATE password_resets
SET used_at = NOW()
WHERE id = $1;

-- name: DeleteUnusedPasswordResets :exec
DELETE FROM password_resets
WHERE user_id = $1
AND used_at IS NULL;

-- name: GetUnusedPasswordResets :many
SELECT * FROM password_resets
WHERE user_id = $1
AND used_at IS NULL
ORDER BY created_at DESC;
