-- name: CreateEmailVerification :one
INSERT INTO email_verifications (
    user_id,
    token_hash,
    expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetEmailVerificationByToken :one
SELECT * FROM email_verifications
WHERE token_hash = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkEmailVerificationAsUsed :exec
UPDATE email_verifications
SET used_at = NOW()
WHERE id = $1;

-- name: DeleteUnverifiedTokens :exec
DELETE FROM email_verifications
WHERE user_id = $1
AND used_at IS NULL;

-- name: GetUnverifiedVerifications :many
SELECT * FROM email_verifications
WHERE user_id = $1
AND used_at IS NULL
ORDER BY created_at DESC;
