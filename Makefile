# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_NAME=$(WEDDING_SERVICE_HOSTNAME)
BUILD_DIR=./docker-volumes/main-service-files


# Default: clean, build, and run containers as daemon
all: docker-stop rm-executable go-build run-as-daemon

go-generate:
	@echo "Running go generate..."
	go generate ./...
	@echo

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo

rm-executable:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(APP_NAME)
	@echo

go-build:
	@echo "Building $(APP_NAME) for $(DOCKER_GOOS)/$(DOCKER_GOARCH)..."
	GOOS=$(DOCKER_GOOS) GOARCH=$(DOCKER_GOARCH) CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo

docker-build:
	@echo "Building Docker containers..."
	docker-compose up --build
	@echo

run-as-daemon:
	@echo "Starting Docker containers as daemons..."
	docker-compose up -d
	@echo

test: go-test

go-test:
	@echo "Running all tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -v ./...
	@echo

test-coverage: go-test-coverage

go-test-coverage:
	@echo "Running tests with coverage for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "To view HTML coverage report, run: go tool cover -html=coverage.out"
	@echo

test-coverage-html: go-test-coverage-html

go-test-coverage-html: test-coverage
	@echo "Viewing test coverage as html for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go tool cover -html=coverage.out

bench: go-bench

go-bench:
	@echo "Running benchmarks for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -bench=.
	@echo
