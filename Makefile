# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Default: clean, build, and run containers as daemon
all: docker-stop rm-executable go-build-for-docker run-as-daemon

generate: go-generate

go-generate:
	cd services/wedding_service && go generate -tags self_sign_cert ./...

go-generate-cert:
	@echo "Running go generate..."
	cd services/wedding_service && go generate -tags self_sign_cert webserver/certificate/self_sign_cert/self_sign_cert.go
	@echo

go-generate-swagger:
	@echo "Running go generate..."
	cd services/wedding_service && go generate webserver/mux.go
	@echo

down: docker-stop

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down --remove-orphans
	@echo

rm-executable:
	@echo "Cleaning build artifacts..."
	rm -f bin/$(APP_NAME)
	@echo

build: go-build

go-build:
	@echo "Go building..."
	@echo "Building bin/$(APP_NAME) from services/wedding_service"
	mkdir -p bin
	cd services/wedding_service && CGO_ENABLED=0 go build -o ../../bin/$(APP_NAME) .
	@echo

go-build-for-docker:
	@echo "Building Docker images for services..."
	docker-compose build gateway_service wedding_service
	@echo



go-build-for-github:
	@echo "Building $(APP_NAME) from services/wedding_service"
	cd services/wedding_service && CGO_ENABLED=0 go build -o ../../$(APP_NAME) .
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
	@echo "Running wedding_service"
	cd services/wedding_service && go run .
	@echo

test: go-test

go-test:
	@echo "Running all tests for wedding_service module"
	cd services/wedding_service && go test -count=1 ./...
	@echo

test-v: go-test-v

go-test-v:
	@echo "Running all tests (verbose) for wedding_service module"
	cd services/wedding_service && go test -count=1 -v ./...
	@echo

test-cached: go-test-cached

go-test-cached:
	@echo "Running cached tests for wedding_service module"
	cd services/wedding_service && go test ./...
	@echo

test-coverage: go-test-coverage

go-test-coverage:
	@echo "Running tests with coverage (excluding swagger) for wedding_service module"
	cd services/wedding_service && go test -coverprofile=coverage-raw.out ./...
	cd services/wedding_service && grep -v "/swagger/" coverage-raw.out > coverage.out || true
	cd services/wedding_service && go tool cover -func=coverage.out
	cp services/wedding_service/coverage.out ./coverage.out
	@echo "To view HTML coverage report, run:"
	@echo "  (cd services/wedding_service && go tool cover -html=coverage.out)"
	@echo

test-coverage-html: go-test-coverage-html

go-test-coverage-html: test-coverage
	@echo "Viewing test coverage as html"
	cd services/wedding_service && go tool cover -html=coverage.out

bench: go-bench

go-bench:
	@echo "Running benchmarks for wedding_service module"
	@rm -f cpu.prof mem.prof block.prof mutex.prof trace.out goroutine.prof
	cd services/wedding_service && go test -bench=. -run=^$$ -benchtime=5s -benchmem \
		-cpuprofile=cpu.prof \
		-memprofile=mem.prof \
		-blockprofile=block.prof \
		-mutexprofile=mutex.prof \
		-trace trace.out
	@echo "Benchmark completed."

bench-analysis: go-bench analysis

analysis:
	@echo "Analyzing CPU profile..."
	@cd services/wedding_service && go tool pprof -top cpu.prof
	@echo "Analyzing memory profile..."
	@cd services/wedding_service && go tool pprof -top mem.prof
	@echo "Analyzing mutex profile..."
	@cd services/wedding_service && go tool pprof -top mutex.prof
	@echo "Analyzing block profile..."
	@cd services/wedding_service && go tool pprof -top block.prof
	@echo "To view detailed profiles:"
	@echo "  CPU:       (cd services/wedding_service && go tool pprof -http=:9090 cpu.prof)"
	@echo "  Mem:       (cd services/wedding_service && go tool pprof -http=:9090 mem.prof)"
	@echo "  Mutex:     (cd services/wedding_service && go tool pprof -http=:9090 mutex.prof)"
	@echo "  Block:     (cd services/wedding_service && go tool pprof -http=:9090 block.prof)"
	@echo "  Trace:     (cd services/wedding_service && go tool trace trace.out)"

race: go-race

go-race:
	@echo "Testing for race conditions (wedding_service module)"
	cd services/wedding_service && go test -count=1 -race ./...
	@echo "Race test completed."


race-v: go-race-v

go-race-v:
	@echo "Testing for race conditions (verbose, wedding_service module)"
	cd services/wedding_service && go test -count=1 -race -v ./...
	@echo "Race test completed."

fuzz: go-fuzz

go-fuzz:
	@echo "Running fuzz tests (wedding_service module)"
	cd services/wedding_service && go test -fuzz=Fuzz -fuzztime=10s
	@echo "Fuzz test completed."
