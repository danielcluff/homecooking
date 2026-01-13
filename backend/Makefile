# HomeCooking Backend Makefile

.PHONY: build run test test-coverage clean clean-coverage \
	test-repository test-service test-handler test-middleware \
	db-migrate db-reset db-shell

# Build
build:
	go build -o bin/server ./cmd/server

# Run
run:
	./run.sh

# Test targets
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean-coverage:
	rm -f coverage.out coverage.html

# Run tests by layer
test-repository:
	go test -v ./internal/repository/...

test-service:
	go test -v ./internal/services/...

test-handler:
	go test -v ./internal/handlers/...

test-middleware:
	go test -v ./internal/middleware/...

# Clean build artifacts
clean:
	rm -f bin/server

# Run tests with race detector
test-race:
	go test -race ./...

# Run tests with coverage for each package
test-detailed:
	@echo "=== Repository Tests ==="
	go test -v -coverprofile=repo_coverage.out ./internal/repository/...
	@echo ""
	@echo "=== Service Tests ==="
	go test -v -coverprofile=service_coverage.out ./internal/services/...
	@echo ""
	@echo "=== Handler Tests ==="
	go test -v -coverprofile=handler_coverage.out ./internal/handlers/...
	@echo ""
	@echo "=== Middleware Tests ==="
	go test -v -coverprofile=middleware_coverage.out ./internal/middleware/...

# Database commands
db-migrate:
	@echo "Running PostgreSQL migrations..."
	docker exec -i homecooking-db psql -U postgres -d homecooking < internal/db/migrations/001_init.up.sql
	docker exec -i homecooking-db psql -U postgres -d homecooking < internal/db/migrations/002_add_recipe_variations.up.sql
	@echo "Migrations complete!"

db-reset:
	@echo "Resetting database..."
	docker exec homecooking-db psql -U postgres -c "DROP DATABASE IF EXISTS homecooking;"
	docker exec homecooking-db psql -U postgres -c "CREATE DATABASE homecooking;"
	@echo "Database reset complete. Run 'make db-migrate' to apply migrations."

db-shell:
	@echo "Connecting to PostgreSQL database..."
	docker exec -it homecooking-db psql -U postgres -d homecooking
