.PHONY: fmt lint test test-race build clean deps check install-tools help

# Packages to build/test (exclude specs directory)
PACKAGES := $(shell go list ./... | grep -v '/specs/')

# Go files to format (exclude specs directory)
GO_FILES := $(shell find . -name '*.go' -not -path './specs/*' -not -path './vendor/*')

# Default target
.DEFAULT_GOAL := help

# Help target
help:
	@echo "Available targets:"
	@echo "  make fmt           - Format code using gofmt (and goimports if installed)"
	@echo "  make lint          - Run linters (golangci-lint or go vet)"
	@echo "  make test          - Run tests with coverage"
	@echo "  make test-race     - Run tests with race detector (requires CGO)"
	@echo "  make coverage      - Generate HTML coverage report"
	@echo "  make build         - Build all packages"
	@echo "  make clean         - Clean build artifacts and coverage files"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make install-tools - Install development tools (goimports, golangci-lint)"
	@echo "  make check         - Run fmt + lint + test"
	@echo "  make help          - Show this help message"

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		goimports -w -local github.com/Zillaforge/cloud-sdk $(GO_FILES); \
	else \
		echo "⚠️  goimports not found. Run 'make install-tools' to install it."; \
	fi

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./internal/... ./modules/... .; \
	else \
		echo "⚠️  golangci-lint not found. Run 'make install-tools' to install it."; \
		echo "Running basic go vet instead..."; \
		go vet $(PACKAGES); \
	fi

# Run tests (without race detector for Alpine compatibility)
test:
	@echo "Running tests..."
	@go test -cover ./...

# Run tests with race detector (requires CGO)
test-race:
	@echo "Running tests with race detector..."
	@if [ "$$CGO_ENABLED" = "1" ]; then \
		go test -v -race -coverprofile=coverage.out $(PACKAGES); \
	else \
		echo "⚠️  Race detector requires CGO. Set CGO_ENABLED=1 to enable."; \
		echo "Running tests without race detector..."; \
		go test -v -coverprofile=coverage.out $(PACKAGES); \
	fi

# Run tests with coverage report
coverage: test
	@go test -v -coverprofile=coverage.out $(PACKAGES)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build
build:
	@echo "Building..."
	@go build $(PACKAGES)

# Clean
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean $(PACKAGES)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed successfully"

# Run all checks
check: fmt lint test
	@echo "All checks passed!"
