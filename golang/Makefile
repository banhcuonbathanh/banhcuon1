# Variables
SERVER_CMD = cmd/server/main.go
CLIENT_CMD = cmd/client/main.go
TEST_USER_DIR = cmd/client
PID_FILE = .server.pid

# Phony targets
.PHONY: all run-server stop-server run-client test test-create test-get clean

# Default target
all: run-server run-client test

# Run the server
run-server:
	@echo "Starting the gRPC server..."
	@go run $(SERVER_CMD) & echo $$! > $(PID_FILE)

# Stop the server
stop-server:
	@echo "Stopping the gRPC server..."
	@-kill `cat $(PID_FILE)` 2>/dev/null || true
	@rm -f $(PID_FILE)

# Run the client
run-client:
	@echo "Running the gRPC client..."
	@go run $(CLIENT_CMD)

# Run all tests
test:
	@echo "Running all tests..."
	@cd $(TEST_USER_DIR) && go test -v ./...

# Run specific test for CreateUser
test-create:
	@echo "Running CreateUser test..."
	@cd $(TEST_USER_DIR) && go test -v -run TestCreateUser

# Run specific test for GetUser
test-get:
	@echo "Running GetUser test..."
	@cd $(TEST_USER_DIR) && go test -v -run TestGetUser

# Clean build artifacts (if any)
clean:
	@echo "Cleaning build artifacts..."
	@go clean
	@rm -f $(TEST_USER_DIR)/*.test
	@rm -f $(PID_FILE)

# Help target
help:
	@echo "Available targets:"
	@echo "  run-server  : Start the gRPC server"
	@echo "  stop-server : Stop the gRPC server"
	@echo "  run-client  : Run the gRPC client"
	@echo "  test        : Run all tests"
	@echo "  test-create : Run CreateUser test"
	@echo "  test-get    : Run GetUser test"
	@echo "  clean       : Clean build artifacts"
	@echo "  help        : Show this help message"