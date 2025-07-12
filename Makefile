.PHONY: build test clean generate docker-build docker-build-all

# Binary name
FULCRUM_AGENT_BINARY=fcfmagent

# Docker settings
DOCKER_IMAGE=cfm-fulcrum
DOCKER_TAG=latest
DOCKER_REGISTRY?=

# Build settings
BUILD_DIR=bin
FULCRUM_AGENT_PATH=./cmd/agent/main.go

# Environment variables
export CGO_ENABLED=0

# Install development tools
install-tools: install-mockery

install-mockery:
	go install github.com/vektra/mockery/v2@latest

# Create all generated code (including mocks)
generate:
	@echo "Generating code..."
	@which mockery > /dev/null || (echo "Installing mockery..." && go install github.com/vektra/mockery/v2@latest)
	go generate ./...
	@echo "Tidying dependencies..."
	go mod tidy

build: build-fulcrum-agent

# Build the server
build-fulcrum-agent:
	go build -o $(BUILD_DIR)/$(FULCRUM_AGENT_BINARY) $(FULCRUM_AGENT_PATH)

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker build with registry prefix
docker-build-registry:
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker build for multiple platforms
docker-build-multiarch:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run tests with code generation
test: generate
	go test -v ./...

# Run tests without generating code
test-only:
	go test -v ./...

# Clean build artifacts and generated code
clean:
	rm -rf $(BUILD_DIR)
	rm -rf mocks/
	go clean

# Clean only generated code
clean-generated:
	rm -rf mocks/

# Build for multiple platforms
build-all:
	# Fulcrum agent binaries
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(FULCRUM_AGENT_BINARY)-linux-amd64 $(FULCRUM_AGENT_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(FULCRUM_AGENT_BINARY)-darwin-amd64 $(FULCRUM_AGENT_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(FULCRUM_AGENT_BINARY)-darwin-arm64 $(FULCRUM_AGENT_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(FULCRUM_AGENT_BINARY)-windows-amd64.exe $(FULCRUM_AGENT_PATH)