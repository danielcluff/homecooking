-- name: CreateRecipeImage :one
INSERT INTO recipe_images (recipe_id, file_path, webp_path, thumbnail_path, caption, order_index)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRecipeImageByID :one
SELECT * FROM recipe_images
WHERE id = $1 LIMIT 1;

-- name: UpdateRecipeImage :one
UPDATE recipe_images
SET
    caption = COALESCE(sqlc.narg('caption'), caption),
    order_index = COALESCE(sqlc.narg('order_index'), order_index)
WHERE id = $1
RETURNING *;

-- name: DeleteRecipeImage :exec
DELETE FROM recipe_images WHERE id = $1;

-- name: GetRecipeImages :many
SELECT * FROM recipe_images
WHERE recipe_id = $1
ORDER BY order_index;

-- name: GetRecipeWithImages :one
SELECT 
    r.*,
    COALESCE(
        json_agg(
            json_build_object(
                'id', ri.id,
                'file_path', ri.file_path,
                'webp_path', ri.webp_path,
                'thumbnail_path', ri.thumbnail_path,
                'caption', ri.caption,
                'order_index', ri.order_index
            ) ORDER BY ri.order_index
        ) FILTER (WHERE ri.id IS NOT NULL), 
        '[]'::json
    ) as body_images
FROM recipes r
LEFT JOIN recipe_images ri ON r.id = ri.recipe_id
WHERE r.id = $1
GROUP BY r.id;
