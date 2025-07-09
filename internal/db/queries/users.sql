-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash,
    role,
    privacy_settings
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;


-- name: UpdateUser :exec
UPDATE users
SET email = $2, username = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, password_changed_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: MarkEmailVerified :exec
UPDATE users
SET email_verified = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserPrivacySettings :exec
UPDATE users 
SET privacy_settings = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE users
SET active = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetUserWithProfile :one
SELECT u.id, u.email, u.username, u.email_verified, u.active, u.role, u.privacy_settings,  u.created_at, 
    u.last_login_at, u.updated_at, up.first_name, up.last_name, up.phone, up.bio
FROM users u 
LEFT JOIN user_profiles up ON u.id = up.user_id
WHERE u.id = $1;

