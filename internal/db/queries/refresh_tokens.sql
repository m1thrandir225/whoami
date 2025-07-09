-- name: CreateRefreshToken :one

INSERT INTO refresh_tokens (
  user_id,
  token_hash,
  device_info,
  expires_at
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens 
WHERE token_hash = $1 AND expires_at > NOW() and revoked_at IS NULL;

-- name: UpdateRefreshTokenLastUsed :exec
UPDATE refresh_tokens
SET last_used_at = NOW()
WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1;

-- name: RevokeAllUserRefreshTokens :exec 
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: CleanupExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

-- name: GetActiveRefreshTokensByUser :many
SELECT * FROM refresh_tokens
WHERE user_id = $1 AND expires_at > NOW() AND revoked_at IS NULL 
ORDER BY created_at DESC;
