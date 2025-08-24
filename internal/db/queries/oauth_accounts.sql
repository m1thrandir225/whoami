-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (
    user_id,
    provider,
    provider_user_id,
    email,
    name,
    avatar_url,
    access_token,
    refresh_token,
    token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetOAuthAccountByID :one
SELECT * FROM oauth_accounts
WHERE id = $1 AND user_id = $2;

-- name: GetOAuthAccountByProvider :one
SELECT * FROM oauth_accounts
WHERE provider = $1 AND provider_user_id = $2;

-- name: GetOAuthAccountsByUserID :many
SELECT * FROM oauth_accounts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetOAuthAccountByEmail :one
SELECT * FROM oauth_accounts
WHERE email = $1 AND provider = $2;

-- name: UpdateOAuthAccount :one
UPDATE oauth_accounts
SET email = $3,
    name = $4,
    avatar_url = $5,
    access_token = $6,
    refresh_token = $7,
    token_expires_at = $8,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteOAuthAccount :exec
DELETE FROM oauth_accounts
WHERE id = $1 AND user_id = $2;

-- name: DeleteOAuthAccountByProvider :exec
DELETE FROM oauth_accounts
WHERE user_id = $1 AND provider = $2;

-- name: UpdateOAuthTokens :one
UPDATE oauth_accounts
SET access_token = $3,
    refresh_token = $4,
    token_expires_at = $5,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;
