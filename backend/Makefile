# HomeCooking Backend Makefile

.PHONY: build run test test-coverage clean clean-coverage \
	test-repository test-service test-handler test-middleware

# Build
build:
	go build -o bin/server ./cmd/server

# Run
run:
	go run ./cmd/server

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
