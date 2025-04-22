# Makefile for go-captcha-service project

# Variables
BINARY_NAME=go-captcha-service
VERSION?=0.0.1
BUILD_DIR=build
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/arm/v7 windows/amd64
DOCKER_IMAGE?=wenlng/go-captcha-service
GO=go
GOFLAGS=-ldflags="-w -s" -v -a -trimpath
COPY_BUILD_FILES=config.json ecosystem.config.js

# Default Target
.PHONY: all
all: build

# Install Dependencies
.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod download
	npm install -g pm2
	@if ! command -v protoc >/dev/null; then \
		echo "Installing protoc..."; \
		$(GO) install github.com/golang/protobuf/protoc-gen-go@latest; \
	fi
	@if ! command -v grpcurl >/dev/null; then \
		echo "Installing grpcurl..."; \
		$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest; \
	fi
	@if ! command -v yq >/dev/null; then \
		echo "Installing yq..."; \
		$(GO) install github.com/mikefarah/yq/v4@latest; \
	fi
	@if ! command -v modd >/dev/null; then \
	  	echo "Installing modd..."; \
		$(GO) install github.com/cortesi/modd/cmd/modd@latest; \
	fi

# Generate gRPC code
.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. proto/api.proto

.PHONY: start-dev
start-dev:
	modd -f modd.conf
	@echo "Starting modd successfully"

.PHONY: start
start:
	go run ./cmd/go-captcha-service/main.go -config config.dev.json -gocaptcha-config gocaptcha.dev.json
	@echo "Starting service successfully"

# Build the application
.PHONY: build
build: proto
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-captcha-service

# Cross-platform construction
.PHONY: build-multi
build-multi: proto
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output=$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then output=$$output.exe; fi; \
		echo "Building $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build $(GOFLAGS) -o $$output ./cmd/go-captcha-service || exit 1; \
	done

# Packaging Binaries
.PHONY: package
package: build-multi
	@mkdir -p $(BUILD_DIR)/packages
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		binary=$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then binary=$$binary.exe; fi; \
		package=$(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-$$os-$$arch.tar.gz; \
		echo "Packaging $$os/$$arch..."; \
		tar -czf $$package -C $(BUILD_DIR) $(BINARY_NAME)-$$os-$$arch config.json ecosystem.config.js; \
	done

# Run tests
.PHONY: test
test: proto
	$(GO) test -v ./...

# Coverage report
.PHONY: cover
cover: proto
	$(GO) test -cover -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean up
.PHONY: clean
clean:
	rm -rf $(BINARY_NAME) $(BUILD_DIR) coverage.out coverage.html testdata.etcd*
	rm -rf bin
	rm -f proto/*.pb.go
	rm -f docs/openapi.json

# Format code
fmt:
	$(GO) fmt ./...

# Local Docker build (binary)
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):latest .

# Multi-architecture Docker build and push (binary)
.PHONY: docker-build-multi
docker-build-multi:
	docker buildx build \
		--platform linux/amd64,linux/arm64,linux/arm/v7,windows/amd64 \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):amd64-$(VERSION) \
		-t $(DOCKER_IMAGE):arm64-$(VERSION) \
		-t $(DOCKER_IMAGE):armv7-$(VERSION) \
		--push .

# Multi-architecture Docker build and push (binary)
.PHONY: docker-proxy-build-multi
docker-proxy-build-multi:
	export http_proxy=http://127.0.0.1:7890
	export https_proxy=http://127.0.0.1:7890
	docker buildx build \
		--platform linux/amd64,linux/arm64,linux/arm/v7,windows/amd64 \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):amd64-$(VERSION) \
		-t $(DOCKER_IMAGE):arm64-$(VERSION) \
		-t $(DOCKER_IMAGE):armv7-$(VERSION) \
		--push .

# Run a local Docker container (binary)
.PHONY: docker-run
docker-run:
	docker run -d -p 8080:8080 -p 50051:50051 $(DOCKER_IMAGE):latest

# Run a local PM2 service (binary)
.PHONY: pm2-run
pm2-run: build
	pm2 start ecosystem.config.js

# PM2 deployment
.PHONY: pm2-deploy
pm2-deploy: build
	pm2 start ecosystem.config.js --env production

# Docker compose
compose:
	docker-compose up -d

# Help Information
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps             : Install dependencies"
	@echo "  proto            : Generate Protobuf code"
	@echo "  start            : Opening the development environment"
	@echo "  start-dev        : Opening the hot reload development environment"
	@echo "  build            : Build binary for current platform"
	@echo "  build-multi      : Build binaries for all platforms"
	@echo "  package          : Package binaries with config.json"
	@echo "  test             : Run tests"
	@echo "  cover            : Generate test coverage report"
	@echo "  clean            : Remove build artifacts"
	@echo "  docker-build     : Build Docker image locally (binary)"
	@echo "  docker-build-multi : Build and push multi-arch Docker image (binary)"
	@echo "  docker-run       : Run Docker container locally (binary)"
	@echo "  pm2-run          : Run with PM2 locally"
	@echo "  pm2-deploy  	  : Deploy with PM2 in production"
	@echo "  help             : Show this help message"