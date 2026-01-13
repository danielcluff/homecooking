# Phase 2: Core Features - In Progress

## Progress Summary

### Completed ✓

1. **Recipe Repository**
   - Full CRUD operations
   - Search functionality
   - Update featured image
   - Update publish status
   - Proper SQL null handling

2. **Recipe Service**
   - Business logic for recipe operations
   - Slug generation from title
   - Authorization checks (owner only can edit/delete)
   - Validation of required fields

3. **Recipe Handlers**
   - GET /api/v1/recipes (list with pagination)
   - GET /api/v1/recipes/{id} (get by ID)
   - GET /api/v1/recipes/slug/{slug} (get by slug)
   - GET /api/v1/recipes/search (search recipes)
   - POST /api/v1/recipes (create - authenticated)
   - PUT /api/v1/recipes/{id} (update - authenticated, owner only)
   - POST /api/v1/recipes/{id}/publish (publish/unpublish - authenticated, owner only)
   - DELETE /api/v1/recipes/{id} (delete - authenticated, owner only)

4. **Category Repository & Service**
   - Full CRUD operations for categories
   - List all categories
   - Get by ID and slug
   - Slug generation from name

5. **Category Handlers**
   - GET /api/v1/categories (list)
   - GET /api/v1/categories/{id} (get)
   - POST /api/v1/categories (create - authenticated)
   - PUT /api/v1/categories/{id} (update - authenticated)
   - DELETE /api/v1/categories/{id} (delete - authenticated)

6. **Tag Repository & Service**
   - Full CRUD operations for tags
   - Add/remove tags from recipes
   - Get tags for a specific recipe
   - Default color handling (#6366f1)

7. **Tag Handlers**
   - GET /api/v1/tags (list)
   - GET /api/v1/tags/{id} (get)
   - GET /api/v1/recipes/{recipeId}/tags (get recipe tags)
   - POST /api/v1/tags (create - authenticated)
   - PUT /api/v1/tags/{id} (update - authenticated)
   - DELETE /api/v1/tags/{id} (delete - authenticated)
   - POST /api/v1/recipes/{recipeId}/tags/{tagId} (add tag to recipe - authenticated)
   - DELETE /api/v1/recipes/{recipeId}/tags/{tagId} (remove tag from recipe - authenticated)

8. **Type Fixes**
   - Changed int32 in models to match sqlc
   - Proper null handling helpers
   - Fixed all build errors

9. **Frontend Admin Layout**
   - Admin navigation bar with routes
   - Responsive design
   - Dashboard page with stats cards
   - Quick action links

### In Progress 🚧

- Recipe management UI
- Category management UI
- Tag management UI
- Markdown editor component
- File upload system

### Next Steps

1. Create recipe management pages (list, new, edit)
2. Create category management pages
3. Create tag management pages
4. Add markdown editor component
5. Implement file upload system
6. Recipe groups backend and frontend
7. Recipe groups UI

## Updated API Endpoints (Added in Phase 2)

```
Recipes:
- GET /api/v1/recipes/search?q={query}&limit={limit}&offset={offset}
- POST /api/v1/recipes/{id}/publish

Categories:
- GET /api/v1/categories
- GET /api/v1/categories/{id}
- POST /api/v1/categories
- PUT /api/v1/categories/{id}
- DELETE /api/v1/categories/{id}

Tags:
- GET /api/v1/tags
- GET /api/v1/tags/{id}
- GET /api/v1/recipes/{recipeId}/tags
- POST /api/v1/tags
- PUT /api/v1/tags/{id}
- DELETE /api/v1/tags/{id}
- POST /api/v1/recipes/{recipeId}/tags/{tagId}
- DELETE /api/v1/recipes/{recipeId}/tags/{tagId}
```

## Backend Status

- ✅ Recipe CRUD fully functional
- ✅ Category CRUD fully functional
- ✅ Tag CRUD fully functional
- ✅ Authentication working
- ✅ Authorization checks in place
- 🚧 Recipe groups (next)
- ❌ File upload (Phase 2)
- ❌ Image optimization (Phase 3)

## Frontend Status

- ✅ Basic pages (home, recipes, login, register)
- ✅ Tailwind CSS v4 integration
- ✅ Admin layout created
- ✅ Admin dashboard page
- 🚧 Recipe form (next)
- 🚧 Category/tag UI (next)
- ❌ Markdown editor (Phase 2)
- ❌ Image upload (Phase 3)

### Recent Work (Dec 27, 2025)

11. **Frontend Pages Created**
    - Recipe list page (with status badges)
    - Recipe new form page (with markdown editor)
    - Category list page  
    - Category new form page
    - Tag list page
    - Tag new form page
    - Login page with API integration
    - API client utility (src/lib/api.ts)

12. **Type Fixes**
    - Changed int32 in models to match sqlc
    - Proper null handling helpers
    - Fixed all build errors

### Current Issue

- Frontend build failing with JSON parse error
- Need to fix script handling in Astro pages

### Next Steps

1. Fix frontend build errors
2. Add authentication protection to admin pages
3. Create recipe edit page
4. Create recipe detail page with markdown
5. Add file upload system
6. Recipe groups backend and UI
