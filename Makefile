# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_NAME=$(WEDDING_SERVICE_HOSTNAME)
BUILD_DIR=./docker-volumes/main-service-files


# Default: clean, build, and run containers as daemon
all: stop-docker rm-executable build-go run-as-daemon

stop-docker:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo

rm-executable:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(APP_NAME)
	@echo

build-go:
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo

build-docker:
	@echo "Building Docker containers..."
	docker-compose up --build
	@echo

run-as-daemon:
	@echo "Starting Docker containers as daemons..."
	docker-compose up -d
	@echo

test:
	@echo "Running all tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test ./...
	@echo

test-coverage:
	@echo "Running tests with coverage for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "To view HTML coverage report, run: go tool cover -html=coverage.out"
	@echo

test-coverage-html: test-coverage
	@echo "Viewing test coverage as html for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go tool cover -html=coverage.out

bench:
	@echo "Running benchmarks for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -bench=.
	@echo
