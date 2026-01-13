-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1 LIMIT 1;

-- name: GetRecipeBySlug :one
SELECT * FROM recipes
WHERE slug = $1 LIMIT 1;

-- name: CreateRecipe :one
INSERT INTO recipes (id, title, slug, markdown_content, author_id, category_id, description, prep_time_minutes, cook_time_minutes, servings, difficulty, featured_image_path, is_published)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: UpdateRecipe :one
UPDATE recipes
SET
    title = COALESCE(sqlc.narg('title'), title),
    markdown_content = COALESCE(sqlc.narg('markdown_content'), markdown_content),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    description = COALESCE(sqlc.narg('description'), description),
    prep_time_minutes = COALESCE(sqlc.narg('prep_time_minutes'), prep_time_minutes),
    cook_time_minutes = COALESCE(sqlc.narg('cook_time_minutes'), cook_time_minutes),
    servings = COALESCE(sqlc.narg('servings'), servings),
    difficulty = COALESCE(sqlc.narg('difficulty'), difficulty),
    featured_image_path = COALESCE(sqlc.narg('featured_image_path'), featured_image_path),
    is_published = COALESCE(sqlc.narg('is_published'), is_published),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateRecipePublishedStatus :one
UPDATE recipes
SET
    is_published = $2,
    published_at = CASE WHEN $2 = true THEN NOW() ELSE published_at END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipes WHERE id = $1;

-- name: ListRecipes :many
SELECT * FROM recipes
WHERE is_published = true
ORDER BY published_at DESC
LIMIT $1 OFFSET $2;

-- name: ListRecipesByCategory :many
SELECT * FROM recipes
WHERE is_published = true
  AND category_id = $1
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchRecipes :many
SELECT * FROM recipes
WHERE is_published = true
  AND (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRecipesByAuthor :many
SELECT * FROM recipes
WHERE author_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRecipeFeaturedImage :one
UPDATE recipes
SET
    featured_image_path = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
