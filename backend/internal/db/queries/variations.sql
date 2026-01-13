-- name: CreateVariation :one
INSERT INTO recipe_variations (id, recipe_id, author_id, markdown_content, prep_time_minutes, cook_time_minutes, servings, difficulty, notes, is_published)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetVariationByID :one
SELECT * FROM recipe_variations WHERE id = $1;

-- name: GetVariationsByRecipe :many
SELECT * FROM recipe_variations 
WHERE recipe_id = $1 
ORDER BY created_at DESC;

-- name: GetPublishedVariationsByRecipe :many
SELECT * FROM recipe_variations 
WHERE recipe_id = $1 AND is_published = true
ORDER BY created_at DESC;

-- name: GetVariationByRecipeAndAuthor :one
SELECT * FROM recipe_variations 
WHERE recipe_id = $1 AND author_id = $2;

-- name: UpdateVariation :one
UPDATE recipe_variations
SET
    markdown_content = COALESCE(sqlc.narg('markdown_content'), markdown_content),
    prep_time_minutes = COALESCE(sqlc.narg('prep_time_minutes'), prep_time_minutes),
    cook_time_minutes = COALESCE(sqlc.narg('cook_time_minutes'), cook_time_minutes),
    servings = COALESCE(sqlc.narg('servings'), servings),
    difficulty = COALESCE(sqlc.narg('difficulty'), difficulty),
    notes = COALESCE(sqlc.narg('notes'), notes),
    is_published = COALESCE(sqlc.narg('is_published'), is_published),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteVariation :exec
DELETE FROM recipe_variations WHERE id = $1;

-- name: ListVariationsByAuthor :many
SELECT * FROM recipe_variations
WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetVariationsByRecipeWithAuthor :many
SELECT
	rv.*,
	u.email as author_email,
	u.role as author_role
FROM recipe_variations rv
JOIN users u ON rv.author_id = u.id
WHERE rv.recipe_id = $1
ORDER BY rv.created_at DESC;
