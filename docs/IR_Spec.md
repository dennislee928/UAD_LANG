# .uad-IR Specification (v0.1)

## 1. Goals

- Deterministic execution for reproducible simulations
- Typed instructions for static verification
- Simple stack-based VM

## 2. Value Kinds

- `i64`, `f64`, `bool`, `str`, `addr` (heap reference)

## 3. Instruction Set (draft)

- Arithmetic: `ADD`, `SUB`, `MUL`, `DIV`
- Comparison: `LT`, `GT`, `EQ`, `NEQ`
- Logic: `AND`, `OR`, `NOT`
- Control flow: `JMP`, `JMP_IF`, `CALL`, `RET`
- Memory: `LOAD`, `STORE`, `ALLOC`, `FREE`
- Domain: `EMIT_EVENT`, `RECORD_MISTAKE`, `RECORD_PRIME`, `SAMPLE_RNG`,
  `UPDATE_MACROSTATE`

(…之後你可以補上具體編碼格式與樣例…)
