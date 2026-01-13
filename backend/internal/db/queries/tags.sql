-- name: CreateTag :one
INSERT INTO tags (id, name, slug, color)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = $1 LIMIT 1;

-- name: GetTagBySlug :one
SELECT * FROM tags
WHERE slug = $1 LIMIT 1;

-- name: ListTags :many
SELECT * FROM tags
ORDER BY name ASC;

-- name: UpdateTag :one
UPDATE tags
SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    color = COALESCE(sqlc.narg('color'), color)
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: AddTagToRecipe :exec
INSERT INTO recipe_tags (recipe_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromRecipe :exec
DELETE FROM recipe_tags
WHERE recipe_id = $1 AND tag_id = $2;

-- name: GetRecipeTags :many
SELECT t.* FROM tags t
JOIN recipe_tags rt ON t.id = rt.tag_id
WHERE rt.recipe_id = $1;
