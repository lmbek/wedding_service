# Load .env variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

GOOS=linux
GOARCH=amd64

APP_NAME=$(WEDDING_SERVICE_HOSTNAME)
BUILD_DIR=./docker-volumes/main-service-files

# Use the command Make to fully rerun everything fast
.PHONY: all

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
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo

build-docker:
	@echo "Building Docker containers..."
	docker-compose up --build
	@echo

run-as-daemon:
	@echo "Starting Docker containers as daemons..."
	docker-compose up -d
	@echo


