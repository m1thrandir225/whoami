-- name: CreateDataExport :one
INSERT INTO data_exports (
    user_id,
    export_type,
    status,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetDataExportByID :one
SELECT * FROM data_exports
WHERE id = $1 AND user_id = $2;

-- name: GetDataExportsByUserID :many
SELECT * FROM data_exports
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateDataExportStatus :one
UPDATE data_exports
SET status = $2,
    completed_at = CASE WHEN $2 = 'completed' THEN NOW() ELSE completed_at END
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: UpdateDataExportFile :one
UPDATE data_exports
SET file_path = $2,
    file_size = $3
WHERE id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteDataExport :exec
DELETE FROM data_exports
WHERE id = $1 AND user_id = $2;

-- name: DeleteExpiredDataExports :exec
DELETE FROM data_exports
WHERE expires_at < NOW();

-- name: GetPendingDataExports :many
SELECT * FROM data_exports
WHERE status = 'pending'
ORDER BY created_at ASC;
