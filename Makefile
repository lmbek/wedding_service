# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Default: clean, build, and run containers as daemon
all: docker-stop rm-executable go-build run-as-daemon

generate: go-generate

go-generate:
	@echo "Running go generate..."
	go generate ./...
	@echo

down: docker-stop

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo

rm-executable:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(APP_NAME)
	@echo

go-build:
	@echo "Copying dependencies"
	mkdir -p $(BUILD_DIR)/certificate/
	cp $(LOCALHOST_CERT) $(BUILD_DIR)/certificate/
	cp $(LOCALHOST_KEY) $(BUILD_DIR)/certificate/
	@echo

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

run: go-run

go-run:
	@echo "Running go for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go run .
	@echo

test: go-test

go-test:
	@echo "Running all tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -count=1 ./...
	@echo

test-v: go-test-v

go-test-v:
	@echo "Running all tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -count=1 -v ./...
	@echo

test-cached: go-test-cached

go-test-cached:
	@echo "Running all tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test ./...
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
	@rm -f cpu.prof mem.prof block.prof mutex.prof trace.out goroutine.prof
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH)
	go test -bench=. -run=^$$ -benchtime=5s -benchmem \
		-cpuprofile=cpu.prof \
		-memprofile=mem.prof \
		-blockprofile=block.prof \
		-mutexprofile=mutex.prof \
		-trace trace.out
	@echo "Benchmark completed."

bench-analysis: go-bench analysis

analysis:
	@echo "Analyzing CPU profile..."
	@go tool pprof -top cpu.prof

	@echo "Analyzing memory profile..."
	@go tool pprof -top mem.prof

	@echo "Analyzing mutex profile..."
	@go tool pprof -top mutex.prof

	@echo "Analyzing block profile..."
	@go tool pprof -top block.prof

	@echo "To view detailed profiles:"
	@echo "  CPU:       go tool pprof -http=:8080 cpu.prof"
	@echo "  Mem:       go tool pprof -http=:8080 mem.prof"
	@echo "  Mutex:     go tool pprof -http=:8080 mutex.prof"
	@echo "  Block:     go tool pprof -http=:8080 block.prof"
	@echo "  Trace:     go tool trace trace.out"

race: go-race

go-race:
	@echo "Testing for race conditions for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH) \
	go test -count=1 -race ./...
	@echo "Race test completed."


race-v: go-race-v

go-race-v:
	@echo "Testing for race conditions for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH) \
	go test -count=1 -race -v ./...
	@echo "Race test completed."

fuzz: go-fuzz

go-fuzz:
	@echo "Running fuzz tests for $(TEST_GOOS)/$(TEST_GOARCH)..."
	GOOS=$(TEST_GOOS) GOARCH=$(TEST_GOARCH) \
	go test -fuzz=Fuzz -fuzztime=10s
	@echo "Fuzz test completed."