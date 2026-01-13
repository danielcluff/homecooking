# Database Setup

This guide explains how to set up and use PostgreSQL with Docker for the HomeCooking backend.

## Prerequisites

- Docker and Docker Compose installed
- Make (optional, but recommended)

## Quick Start

### 1. Start PostgreSQL Container

From the project root directory:

```bash
docker-compose up -d
```

This will:
- Start a PostgreSQL 16 container
- Create the `homecooking` database
- Expose PostgreSQL on port 5432
- Automatically run migrations on first startup

### 2. Verify Database Connection

Check if the container is running:

```bash
docker ps | grep homecooking-db
```

Test the database connection:

```bash
docker exec homecooking-db psql -U postgres -c "\l"
```

### 3. Run the Backend Server

From the `backend` directory:

```bash
cd backend
make build  # Build the server
make run    # Run the server (loads .env automatically)
```

The server should connect to PostgreSQL and start on port 8080.

## Database Configuration

The backend uses these PostgreSQL connection settings (from `backend/.env`):

```
DATABASE_TYPE=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=homecooking
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
```

## Useful Commands

### Backend Makefile Commands

From the `backend` directory:

```bash
# Build the server
make build

# Run the server
make run

# Run migrations manually
make db-migrate

# Reset database (drops and recreates)
make db-reset

# Open PostgreSQL shell
make db-shell

# Run tests
make test
```

### Docker Commands

```bash
# Stop the database
docker-compose stop

# Start the database (if already created)
docker-compose start

# Stop and remove the database container
docker-compose down

# Stop and remove everything including the data volume
docker-compose down -v

# View database logs
docker logs homecooking-db

# Execute SQL commands
docker exec -i homecooking-db psql -U postgres -d homecooking -c "SELECT * FROM users;"
```

## Migrations

Migrations are located in `backend/internal/db/migrations/`:
- `001_init.up.sql` - Initial schema
- `002_add_recipe_variations.up.sql` - Recipe variations feature

### Running Migrations Manually

```bash
cd backend
make db-migrate
```

### Resetting the Database

To start fresh:

```bash
cd backend
make db-reset
make db-migrate
```

## Switching Between SQLite and PostgreSQL

### To use PostgreSQL (default):

Edit `backend/.env`:
```
DATABASE_TYPE=postgres
```

### To use SQLite:

Edit `backend/.env`:
```
DATABASE_TYPE=sqlite
DATABASE_PATH=./data.db
```

## Troubleshooting

### "Connection refused" error

Make sure the PostgreSQL container is running:
```bash
docker-compose up -d
```

### "Password authentication failed" error

The server needs environment variables loaded. Use `make run` instead of running the binary directly, or source the .env file:
```bash
export $(cat .env | grep -v '^#' | xargs)
./bin/server
```

### Migrations already applied

The docker-compose.yml mounts the migrations directory, so migrations run automatically on first container startup. If you see "already exists" notices, this is normal.

### Reset database completely

```bash
docker-compose down -v  # Remove container and volume
docker-compose up -d    # Recreate everything
cd backend && make db-migrate
```

## Connection String Format

If you need to connect with other tools:

```
postgresql://postgres:postgres@localhost:5432/homecooking
```

Or for external tools like pgAdmin, TablePlus, etc.:
- Host: `localhost`
- Port: `5432`
- Database: `homecooking`
- Username: `postgres`
- Password: `postgres`
