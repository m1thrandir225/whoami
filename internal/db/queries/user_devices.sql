-- name: CreateUserDevice :one
INSERT INTO user_devices (
    user_id,
    device_id,
    device_name,
    device_type,
    user_agent,
    ip_address,
    trusted
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserDevicesByUserID :many
SELECT * FROM user_devices
WHERE user_id = $1
ORDER BY last_used_at DESC;

-- name: GetUserDeviceByID :one
SELECT * FROM user_devices
WHERE id = $1 AND user_id = $2;

-- name: UpdateUserDeviceLastUsed :one
UPDATE user_devices
SET last_used_at = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserDevice :one
UPDATE user_devices
SET device_name = $3,
    device_type = $4,
    user_agent = $5,
    trusted = $6
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteUserDevice :exec
DELETE FROM user_devices
WHERE id = $1 AND user_id = $2;

-- name: DeleteAllUserDevices :exec
DELETE FROM user_devices
WHERE user_id = $1;

-- name: GetUserDeviceByDeviceID :one
SELECT * FROM user_devices
WHERE user_id = $1 AND device_id = $2;

-- name: MarkDeviceAsTrusted :one
UPDATE user_devices
SET trusted = $3
WHERE id = $1 AND user_id = $2
RETURNING *;
