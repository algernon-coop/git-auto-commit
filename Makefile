.PHONY: build test lint clean install help

# Build variables
BINARY_NAME=git-auto-commit
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	go build ${LDFLAGS} -o ${BINARY_NAME} .

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

lint: ## Run linter
	golangci-lint run

clean: ## Clean build artifacts
	rm -f ${BINARY_NAME}
	rm -f coverage.out

install: build ## Install the binary to $GOPATH/bin
	go install ${LDFLAGS}

run: build ## Build and run the binary
	./${BINARY_NAME}

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
