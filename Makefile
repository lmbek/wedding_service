# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Default: clean, build, and run containers as daemon
all: docker-stop rm-executable go-build-for-docker run-as-daemon

generate: go-generate

go-generate:
	cd src && go generate ./...

go-generate-cert:
	@echo "Running go generate..."
	cd src && go generate -tags self_sign_cert src/certificate/self_sign_cert/self_sign_cert.go
	@echo

go-generate-swagger:
	@echo "Running go generate..."
	cd src && go generate src/webserver/mux.go
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

build: go-build

go-build:
	@echo "Go building..."
	@echo "Building out/$(APP_NAME)"
	CGO_ENABLED=0 go build -o out/$(APP_NAME) .
	@echo

go-build-for-docker:
	@echo "Copying dependencies..."
	mkdir -p $(BUILD_DIR)/certificate/self_sign_cert/
	mkdir -p $(BUILD_DIR)/webserver/website/frontend/
	cp src/$(SELF_SIGNED_CERT_PATH) $(BUILD_DIR)/certificate/self_sign_cert/
	cp src/$(SELF_SIGNED_KEY_PATH) $(BUILD_DIR)/certificate/self_sign_cert/
	cp -r src/webserver/website/frontend/* $(BUILD_DIR)/webserver/website/frontend/
	cp .env $(BUILD_DIR)/.env
	@echo

	@echo "Building $(BUILD_DIR)/$(APP_NAME)"
	cd src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../$(BUILD_DIR)/$(APP_NAME) .
	@echo

go-build-for-github:
	@echo "Building $(APP_NAME)"
	CGO_ENABLED=0 go build -o $(APP_NAME) . #github don't need a build directory, we just test if it can build
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
	@echo "Running go"
	go run .
	@echo

test: go-test

go-test:
	@echo "Running all tests"
	go test -count=1 ./...
	@echo

test-v: go-test-v

go-test-v:
	@echo "Running all tests"
	go test -count=1 -v ./...
	@echo

test-cached: go-test-cached

go-test-cached:
	@echo "Running all tests"
	go test ./...
	@echo

test-coverage: go-test-coverage

go-test-coverage:
	@echo "Running tests with coverage"
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "To view HTML coverage report, run: go tool cover -html=coverage.out"
	@echo

test-coverage-html: go-test-coverage-html

go-test-coverage-html: test-coverage
	@echo "Viewing test coverage as html"
	go tool cover -html=coverage.out

bench: go-bench

go-bench:
	@echo "Running benchmarks"
	@rm -f cpu.prof mem.prof block.prof mutex.prof trace.out goroutine.prof
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
	@echo "Testing for race conditions"
	go test -count=1 -race ./...
	@echo "Race test completed."


race-v: go-race-v

go-race-v:
	@echo "Testing for race conditions"
	go test -count=1 -race -v ./...
	@echo "Race test completed."

fuzz: go-fuzz

go-fuzz:
	@echo "Running fuzz tests"
	go test -fuzz=Fuzz -fuzztime=10s
	@echo "Fuzz test completed."