# Makefile

# Variables
BINARY_NAME=url-shortener
WORKERS_BINARY_NAME=url-shortener-workers
BUILD_DIR=build
GOPATH=$(shell go env GOPATH)
COMPILEDAEMON=$(GOPATH)/bin/CompileDaemon

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

# Build the workers
build-workers:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(WORKERS_BINARY_NAME) cmd/workers/workers.go

# Run the application
run:
	@go run cmd/server/main.go

# Run the workers
run:
	@go run cmd/workers/workers.go

# Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Install CompileDaemon for hot reloading
install-daemon:
	@echo "Installing CompileDaemon..."
	@go install github.com/githubnemo/CompileDaemon@latest

# Run with hot reloading using CompileDaemon
dev: install-daemon
	@$(COMPILEDAEMON) --build="go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go" --command="./$(BUILD_DIR)/$(BINARY_NAME)" --color=true -pattern="(.+\.go|.+\.env)$$"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...
