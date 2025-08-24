-- name: CreateAccountLockout :one
INSERT INTO account_lockouts (
    user_id,
    ip_address,
    lockout_type,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAccountLockoutByUserID :one
SELECT * FROM account_lockouts
WHERE user_id = $1
AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAccountLockoutByIP :one
SELECT * FROM account_lockouts
WHERE ip_address = $1
AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAccountLockoutByUserAndIP :one
SELECT * FROM account_lockouts
WHERE user_id = $1
AND ip_address = $2
AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteExpiredLockouts :exec
DELETE FROM account_lockouts
WHERE expires_at <= NOW();

-- name: DeleteAccountLockoutByID :exec
DELETE FROM account_lockouts
WHERE id = $1;
