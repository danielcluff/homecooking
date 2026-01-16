# Technical Details and Specifications

## Overview
HomeCooking is a modern web application designed for recipe management and sharing. It consists of a Go backend API and an Astro frontend, supporting multi-user access, AI integration, and various database backends.

## Architecture
The application follows a layered architecture with clear separation of concerns:

- **Backend (Go)**:
  - **Entry Point**: `backend/cmd/server/main.go` - initializes the server, repositories, services, handlers, and routes using `http.ServeMux`.
  - **Data Layer**: Uses `sqlc` for type-safe SQL queries. Migrations in `backend/internal/db/migrations/`.
  - **Business Logic**: Services in `backend/internal/services/` handle domain logic (recipe management, authentication).
  - **Data Access**: Repositories in `backend/internal/repository/` abstract database operations.
  - **API Layer**: Handlers in `backend/internal/handlers/` process HTTP requests and return JSON responses.
  - **Middleware**: Authentication, logging, and CORS in `backend/internal/middleware/`.
  - **Models**: Domain models in `backend/internal/models/` (users, recipes, variations, etc.).
  - **AI & Storage**: Optional AI service for recipe extraction/enhancement using OpenAI API. Local filesystem storage for images.
  - **Configuration**: Environment-based config from `.env` files.

- **Frontend (Astro)**:
  - **Pages**: File-based routing in `frontend/src/pages/` (e.g., index.astro, admin panels, recipe views).
  - **Components**: Reusable UI in `frontend/src/components/`.
  - **Layouts**: Page layouts in `frontend/src/layouts/`.
  - **Utilities**: Helpers in `frontend/src/lib/`.
  - **Styling**: Tailwind CSS v4 for responsive design.

- **Database**:
  - Supports PostgreSQL (recommended), SQLite (development), MySQL.
  - Schema: users, recipes, categories, tags, recipe_groups, recipe_variations, share_codes, user_invites.
  - Migrations for schema evolution.

## Technologies
- **Backend**: Go 1.25.2, golang-jwt/jwt/v5, google/uuid, lib/pq, mattn/go-sqlite3, stretchr/testify.
- **Frontend**: Astro 5.x, Tailwind CSS v4, marked.
- **Development**: Makefile, Docker Compose, Git.

## Key Features Specifications
- **Recipe Management**: Markdown-formatted recipes with metadata (prep time, cook time, servings, difficulty).
- **Categorization & Tagging**: Organize recipes into categories and apply custom tags.
- **Recipe Groups**: Group related recipes for meal planning.
- **Recipe Variations**: Authenticated users can create personal variations.
- **Image Support**: Upload featured/inline images with optimization; served statically.
- **AI Integration**: Optional recipe extraction from images and text enhancement via OpenAI.
- **Sharing & Access**: Share via codes, invites; role-based permissions (admin, editor, user).
- **Authentication**: JWT with refresh tokens, bcrypt hashing.
- **Deployment**: Backend port 8080, frontend 4321; Docker for database.

## File Structure
- `backend/`: Go backend code.
- `frontend/`: Astro frontend.
- `docs/`: Documentation including phase completions.
- Root: README.md, docker-compose.yml, setup scripts.

## Dependencies
Listed in go.mod and package.json; includes standard Go libraries and Astro ecosystem packages.

## Testing
Backend testing with testify; coverage reports available.

This specification covers the core technical aspects of the HomeCooking codebase.