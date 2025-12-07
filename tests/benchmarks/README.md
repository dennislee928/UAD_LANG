# UAD Language Benchmarks

This directory contains benchmark implementations for comparing UAD with other languages.

## Directory Structure

```
tests/benchmarks/
├── uad/          # UAD language implementations
├── go/           # Go language implementations
├── rust/         # Rust language implementations
├── results/      # Benchmark results (JSON files)
└── run_bench.sh  # Benchmark execution script
```

## Prerequisites

Before running benchmarks, ensure you have:

1. **hyperfine** - Statistical benchmarking tool
   ```bash
   # Install via cargo
   cargo install hyperfine
   
   # Or via Homebrew (macOS)
   brew install hyperfine
   
   # Or via package manager (Linux)
   # Ubuntu/Debian: sudo apt install hyperfine
   # Fedora: sudo dnf install hyperfine
   ```

2. **Go** - https://go.dev/dl/

3. **Rust** - https://rustup.rs/

4. **UAD Tools** - Build from project root:
   ```bash
   cd ../..
   make build
   ```

## Scenarios

### scenario1: Ransomware Kill-chain
Models a complete ransomware attack lifecycle.

### scenario2: SOC ERH Model
Security Operations Center with Error Rate Hypothesis modeling.

### scenario3: Psychohistory Macro Dynamics
Large-scale population dynamics simulation.

## Installation

Before running benchmarks, install all dependencies. See [INSTALL.md](INSTALL.md) for detailed installation instructions.

Quick check:
```bash
./run_bench.sh
```

The script will verify all dependencies and provide installation commands for any missing tools.

## Running Benchmarks

See [../../docs/benchmarks.md](../../docs/benchmarks.md) for detailed methodology.

Quick start:
```bash
cd tests/benchmarks
./run_bench.sh scenario1
```

Results will be saved in `results/` directory as JSON files.

## Adding New Scenarios

1. Create implementation in each language directory:
   - `uad/scenario_name.uad`
   - `go/scenario_name.go`
   - `rust/scenario_name.rs`

2. Ensure all implementations:
   - Solve the same problem
   - Use appropriate language idioms
   - Are optimized for performance
   - Include error handling

3. Update `run_bench.sh` if needed for scenario-specific setup

4. Document the scenario in `docs/benchmarks.md`

