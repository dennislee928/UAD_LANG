.PHONY: all build test clean run-examples install fmt lint help benchmark

# Variables
BINARY_DIR := bin
UADC := $(BINARY_DIR)/uadc
UADVM := $(BINARY_DIR)/uadvm
UADREPL := $(BINARY_DIR)/uadrepl
UADI := $(BINARY_DIR)/uadi
UADRUNNER := $(BINARY_DIR)/uad-runner

GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Default target
all: build

# Build all binaries
build: $(UADC) $(UADVM) $(UADREPL) $(UADI) $(UADRUNNER)

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

# Build Experiment Runner
$(UADRUNNER):
	@echo "Building uad-runner..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADRUNNER) ./cmd/uad-runner

# Build LSP Server
UADLSP := $(BINARY_DIR)/uad-lsp
$(UADLSP):
	@echo "Building uad-lsp..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADLSP) ./cmd/uad-lsp

build-lsp: $(UADLSP)

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

# Run a specific example
example: build
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make example FILE=examples/core/hello_world.uad"; \
		exit 1; \
	fi
	@echo "Running $(FILE)..."
	./bin/uadi -i $(FILE)

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	$(GO) install ./cmd/uadc
	$(GO) install ./cmd/uadvm
	$(GO) install ./cmd/uadrepl
	$(GO) install ./cmd/uadi
	$(GO) install ./cmd/uad-runner

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@echo "📚 Documentation available in docs/"
	@echo "  - docs/ARCHITECTURE.md - System architecture"
	@echo "  - docs/REPO_SNAPSHOT.md - Project status"
	@echo "  - docs/LANGUAGE_SPEC.md - Language specification"
	@echo "  - docs/specs/ - Formal specifications"
	@echo "  - docs/reports/ - Development reports"
	@echo ""
	@echo "To generate Go package documentation:"
	@echo "  godoc -http=:6060"

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

# Run experiments
run-experiments: build
	@echo "Running all experiments..."
	@for config in experiments/configs/*.yaml; do \
		echo ""; \
		echo "Running experiment: $$config"; \
		./bin/uad-runner -config $$config; \
	done

# Run specific experiment
experiment: build
	@if [ -z "$(CONFIG)" ]; then \
		echo "Usage: make experiment CONFIG=experiments/configs/erh_demo.yaml"; \
		exit 1; \
	fi
	@echo "Running experiment: $(CONFIG)"
	./bin/uad-runner -config $(CONFIG)

# Run benchmarks
benchmark: build
	@if [ -z "$(SCENARIO)" ]; then \
		echo "Running benchmark for scenario1..."; \
		./run_bench.sh scenario1; \
	else \
		echo "Running benchmark for $(SCENARIO)..."; \
		./run_bench.sh $(SCENARIO); \
	fi

# Help
help:
	@echo "UAD Language Makefile"
	@echo ""
	@echo "Build Targets:"
	@echo "  all           - Build all binaries (default)"
	@echo "  build         - Build uadc, uadvm, uadrepl, and uadi"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install binaries to GOPATH/bin"
	@echo ""
	@echo "Development Targets:"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format all Go code"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  vet           - Run go vet"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  dev-setup     - Setup development environment"
	@echo ""
	@echo "Execution Targets:"
	@echo "  run-examples    - Run all example programs"
	@echo "  example         - Run a specific example (make example FILE=path/to/file.uad)"
	@echo "  run-experiments - Run all experiments"
	@echo "  experiment      - Run specific experiment (make experiment CONFIG=path/to/config.yaml)"
	@echo "  benchmark       - Run benchmarks (make benchmark SCENARIO=scenario1)"
	@echo ""
	@echo "Documentation Targets:"
	@echo "  docs          - Show documentation locations"
	@echo "  help          - Show this help message"

