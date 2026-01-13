-- name: CreateRecipeGroup :one
INSERT INTO recipe_groups (id, name, slug, description, icon)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRecipeGroupByID :one
SELECT * FROM recipe_groups
WHERE id = $1 LIMIT 1;

-- name: GetRecipeGroupBySlug :one
SELECT * FROM recipe_groups
WHERE slug = $1 LIMIT 1;

-- name: ListRecipeGroups :many
SELECT * FROM recipe_groups
ORDER BY created_at DESC;

-- name: UpdateRecipeGroup :one
UPDATE recipe_groups
SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    icon = COALESCE(sqlc.narg('icon'), icon)
WHERE id = $1
RETURNING *;

-- name: DeleteRecipeGroup :exec
DELETE FROM recipe_groups WHERE id = $1;

-- name: AddRecipeToGroup :exec
INSERT INTO recipe_groupings (group_id, recipe_id, order_index)
VALUES ($1, $2, $3)
ON CONFLICT (group_id, recipe_id) DO UPDATE SET order_index = $3;

-- name: RemoveRecipeFromGroup :exec
DELETE FROM recipe_groupings
WHERE group_id = $1 AND recipe_id = $2;

-- name: GetRecipesInGroup :many
SELECT r.* 
FROM recipes r
JOIN recipe_groupings rg ON r.id = rg.recipe_id
WHERE rg.group_id = $1
  AND r.is_published = true
ORDER BY rg.order_index;

-- name: GetRecipeGroupWithRecipes :one
SELECT 
    g.*,
    COALESCE(
        json_agg(
            json_build_object(
                'id', r.id,
                'title', r.title,
                'slug', r.slug,
                'featured_image_path', r.featured_image_path,
                'description', r.description,
                'order_index', rg.order_index
            ) ORDER BY rg.order_index
        ) FILTER (WHERE r.id IS NOT NULL), 
        '[]'::json
    ) as recipes
FROM recipe_groups g
LEFT JOIN recipe_groupings rg ON g.id = rg.group_id
LEFT JOIN recipes r ON rg.recipe_id = r.id
WHERE g.id = $1
GROUP BY g.id;

-- name: GetGroupsForRecipe :many
SELECT g.* 
FROM recipe_groups g
JOIN recipe_groupings rg ON g.id = rg.group_id
WHERE rg.recipe_id = $1
ORDER BY g.name;
