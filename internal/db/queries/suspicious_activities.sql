-- name: CreateSuspiciousActivity :one
INSERT INTO suspicious_activities (
    user_id,
    activity_type,
    ip_address,
    user_agent,
    description,
    metadata,
    severity
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSuspiciousActivitiesByUserID :many
SELECT * FROM suspicious_activities
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetSuspiciousActivitiesByIP :many
SELECT * FROM suspicious_activities
WHERE ip_address = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetUnresolvedSuspiciousActivities :many
SELECT * FROM suspicious_activities
WHERE resolved = false
ORDER BY created_at DESC
LIMIT $1;

-- name: ResolveSuspiciousActivity :exec
UPDATE suspicious_activities
SET resolved = true
WHERE id = $1;

-- name: GetSuspiciousActivityCountByUser :one
SELECT COUNT(*) FROM suspicious_activities
WHERE user_id = $1
AND created_at > NOW() - INTERVAL '24 hours';

-- name: GetSuspiciousActivityCountByIP :one
SELECT COUNT(*) FROM suspicious_activities
WHERE ip_address = $1
AND created_at > NOW() - INTERVAL '24 hours';
