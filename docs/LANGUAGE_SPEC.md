# .uad-core Language Specification (v0.1)

## 1. Lexical Structure

- Identifiers: `[A-Za-z_][A-Za-z0-9_]*`
- Keywords: `fn`, `struct`, `enum`, `let`, `if`, `else`, `true`, `false`,
  `match`, `return`, `import`, `module`
- Literals:
  - Integer: `[0-9]+`
  - Float: `[0-9]+\.[0-9]+`
  - String: double-quoted, supports `\"` and `\\`
  - Bool: `true` / `false`

## 2. Types

Primitive types:

- `Int`, `Float`, `Bool`, `String`, `Time`, `Duration`

Composite types:

- `struct` definitions
- `enum` (sum types)
- arrays: `[T]`
- maps: `Map[K, V]`

(…後續你可以再補型別推導、語義細節與 BNF…)
