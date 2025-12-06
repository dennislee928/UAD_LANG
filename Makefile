.PHONY: all build test clean run-examples install fmt lint help

# Variables
BINARY_DIR := bin
UADC := $(BINARY_DIR)/uadc
UADVM := $(BINARY_DIR)/uadvm
UADREPL := $(BINARY_DIR)/uadrepl
UADI := $(BINARY_DIR)/uadi

GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Default target
all: build

# Build all binaries
build: $(UADC) $(UADVM) $(UADREPL) $(UADI)

# Build compiler
$(UADC): 
	@echo "Building uadc..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADC) ./cmd/uadc

# Build VM
$(UADVM):
	@echo "Building uadvm..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADVM) ./cmd/uadvm

# Build REPL
$(UADREPL):
	@echo "Building uadrepl..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADREPL) ./cmd/uadrepl

# Build Interpreter
$(UADI):
	@echo "Building uadi (interpreter)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADI) ./cmd/uadi

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.uadir

# Run example programs
run-examples: build
	@echo "Running examples..."
	@bash scripts/run_examples.sh

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	$(GO) install ./cmd/uadc
	$(GO) install ./cmd/uadvm
	$(GO) install ./cmd/uadrepl

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run ./...

# Run static analysis
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@bash scripts/dev_setup.sh

# Help
help:
	@echo "UAD Language Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all           - Build all binaries (default)"
	@echo "  build         - Build uadc, uadvm, and uadrepl"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Remove build artifacts"
	@echo "  run-examples  - Run example programs"
	@echo "  install       - Install binaries to GOPATH/bin"
	@echo "  fmt           - Format all Go code"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  vet           - Run go vet"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help message"

