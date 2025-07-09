-- name: CreateUserProfile :one
INSERT INTO user_profiles (
  user_id,
  first_name,
  last_name,
  phone,
  avatar_url,
  bio,
  timezone,
  locale
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
) RETURNING *;

-- name: UpdateUserProfile :exec
UPDATE user_profiles
SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5, bio = $6, timezone = $7, locale = $8
WHERE user_id = $1;

-- name: GetUserProfile :one
SELECT * FROM user_profiles WHERE user_id = $1;
