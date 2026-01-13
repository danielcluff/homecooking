-- name: GetSetting :one
SELECT * FROM app_settings
WHERE key = $1 LIMIT 1;

-- name: UpsertSetting :one
INSERT INTO app_settings (key, value, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = NOW()
RETURNING *;

-- name: DeleteSetting :exec
DELETE FROM app_settings WHERE key = $1;

-- name: ListSettings :many
SELECT * FROM app_settings
ORDER BY key;
