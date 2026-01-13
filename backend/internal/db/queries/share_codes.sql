-- name: CreateShareCode :one
INSERT INTO share_codes (recipe_id, code, expires_at, max_uses)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetShareCodeByCode :one
SELECT sc.*, r.title as recipe_title, r.slug as recipe_slug
FROM share_codes sc
JOIN recipes r ON sc.recipe_id = r.id
WHERE sc.code = $1
  AND (sc.expires_at IS NULL OR sc.expires_at > NOW())
  AND (sc.max_uses IS NULL OR sc.use_count < sc.max_uses)
LIMIT 1;

-- name: IncrementShareCodeUse :exec
UPDATE share_codes
SET use_count = use_count + 1
WHERE id = $1;

-- name: DeleteShareCode :exec
DELETE FROM share_codes WHERE id = $1;

-- name: GetShareCodesForRecipe :many
SELECT * FROM share_codes
WHERE recipe_id = $1
ORDER BY created_at DESC;
