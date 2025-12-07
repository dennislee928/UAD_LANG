# UAD Language Benchmarks

This directory contains benchmark implementations for comparing UAD with other languages.

## Directory Structure

```
benchmarks/
├── uad/          # UAD language implementations
├── go/           # Go language implementations
├── rust/         # Rust language implementations
└── results/      # Benchmark results (JSON files)
```

## Scenarios

### scenario1: Ransomware Kill-chain
Models a complete ransomware attack lifecycle.

### scenario2: SOC ERH Model
Security Operations Center with Error Rate Hypothesis modeling.

### scenario3: Psychohistory Macro Dynamics
Large-scale population dynamics simulation.

## Running Benchmarks

See [docs/benchmarks.md](../docs/benchmarks.md) for detailed instructions.

Quick start:
```bash
./run_bench.sh scenario1
```

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

