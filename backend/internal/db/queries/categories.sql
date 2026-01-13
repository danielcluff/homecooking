-- name: CreateCategory :one
INSERT INTO categories (id, name, slug, icon, description, order_index)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY order_index ASC, name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    icon = COALESCE(sqlc.narg('icon'), icon),
    description = COALESCE(sqlc.narg('description'), description),
    order_index = COALESCE(sqlc.narg('order_index'), order_index)
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
