# .uad-core Language Specification (v0.1)

## 1. Lexical Structure

### 1.1 Identifiers

Identifiers follow the pattern: `[A-Za-z_][A-Za-z0-9_]*`

Examples: `x`, `myVar`, `_internal`, `calculate_alpha`

### 1.2 Keywords (Reserved Words)

Complete list of reserved keywords:

**Control Flow:**

- `fn`, `return`, `if`, `else`, `match`, `while`, `for`, `break`, `continue`

**Declarations:**

- `let`, `struct`, `enum`, `type`, `import`, `module`, `pub`

**Literals:**

- `true`, `false`, `nil`

**Pattern Matching:**

- `case`, `when`, `in`

**Domain-Specific:**

- `action`, `judge`, `agent`, `event`, `emit`

**Future Reserved:**

- `async`, `await`, `trait`, `impl`, `where`, `use`, `as`, `const`, `mut`

### 1.3 Literals

**Integer Literals:**

- Decimal: `[0-9]+`
- Hexadecimal: `0x[0-9A-Fa-f]+`
- Binary: `0b[01]+`
- Examples: `42`, `0xFF`, `0b1010`

**Float Literals:**

- Standard: `[0-9]+\.[0-9]+`
- Scientific: `[0-9]+(\.[0-9]+)?[eE][+-]?[0-9]+`
- Examples: `3.14`, `1.5e-10`, `2.0E+5`

**String Literals:**

- Double-quoted: `"..."`
- Escape sequences: `\n`, `\t`, `\r`, `\\`, `\"`, `\0`
- Unicode: `\u{XXXX}`
- Raw strings (future): `r"..."`
- Examples: `"Hello, world!"`, `"Line 1\nLine 2"`, `"Unicode: \u{03B1}"`

**Boolean Literals:**

- `true`, `false`

**Time Literals (Domain-Specific):**

- Duration: `10s`, `5m`, `2h`, `3d`, `90d`
- Timestamp: `@2024-01-15T10:30:00Z`

### 1.4 Operators

**Arithmetic:**

- `+`, `-`, `*`, `/`, `%` (modulo)

**Comparison:**

- `==`, `!=`, `<`, `>`, `<=`, `>=`

**Logical:**

- `&&` (and), `||` (or), `!` (not)

**Bitwise (future):**

- `&`, `|`, `^`, `<<`, `>>`

**Assignment:**

- `=`

**Other:**

- `->` (function return type)
- `=>` (match arm)
- `.` (field access)
- `::` (path separator)
- `..` (range, future)

### 1.5 Delimiters

- `(`, `)` - parentheses
- `{`, `}` - braces
- `[`, `]` - brackets
- `,` - comma
- `;` - semicolon
- `:` - colon

### 1.6 Comments

**Single-line:**

```uad
// This is a single-line comment
```

**Multi-line:**

```uad
/* This is a
   multi-line comment */
```

**Documentation comments (future):**

```uad
/// This function calculates alpha
fn compute_alpha() { ... }
```

---

## 2. Grammar (BNF)

### 2.1 Module Structure

```bnf
Module ::= Decl*

Decl ::= FnDecl
       | StructDecl
       | EnumDecl
       | TypeAlias
       | ImportDecl
```

### 2.2 Declarations

```bnf
FnDecl ::= "fn" Ident "(" ParamList? ")" ("->" Type)? Block

ParamList ::= Param ("," Param)*
Param ::= Ident ":" Type

StructDecl ::= "struct" Ident "{" FieldList? "}"
FieldList ::= Field ("," Field)* ","?
Field ::= Ident ":" Type

EnumDecl ::= "enum" Ident "{" VariantList? "}"
VariantList ::= Variant ("," Variant)* ","?
Variant ::= Ident ("(" TypeList ")")?

TypeAlias ::= "type" Ident "=" Type

ImportDecl ::= "import" StringLiteral
```

### 2.3 Types

```bnf
Type ::= PrimitiveType
       | Ident
       | ArrayType
       | MapType
       | FunctionType
       | "(" Type ")"

PrimitiveType ::= "Int" | "Float" | "Bool" | "String" | "Time" | "Duration"

ArrayType ::= "[" Type "]"

MapType ::= "Map" "[" Type "," Type "]"

FunctionType ::= "fn" "(" TypeList? ")" "->" Type

TypeList ::= Type ("," Type)*
```

### 2.4 Statements

```bnf
Stmt ::= LetStmt
       | ExprStmt
       | ReturnStmt
       | AssignStmt
       | WhileStmt
       | ForStmt
       | BreakStmt
       | ContinueStmt

LetStmt ::= "let" Ident (":" Type)? "=" Expr ";"

ExprStmt ::= Expr ";"

ReturnStmt ::= "return" Expr? ";"

AssignStmt ::= Expr "=" Expr ";"

WhileStmt ::= "while" Expr Block

ForStmt ::= "for" Ident "in" Expr Block

BreakStmt ::= "break" ";"

ContinueStmt ::= "continue" ";"

Block ::= "{" Stmt* (Expr)? "}"
```

### 2.5 Expressions

```bnf
Expr ::= Literal
       | Ident
       | BinaryExpr
       | UnaryExpr
       | CallExpr
       | IfExpr
       | MatchExpr
       | BlockExpr
       | StructLiteral
       | ArrayLiteral
       | MapLiteral
       | FieldAccess
       | IndexExpr
       | ParenExpr

Literal ::= IntLiteral | FloatLiteral | StringLiteral | BoolLiteral | TimeLiteral

BinaryExpr ::= Expr BinaryOp Expr

UnaryExpr ::= UnaryOp Expr

CallExpr ::= Expr "(" ArgList? ")"
ArgList ::= Expr ("," Expr)*

IfExpr ::= "if" Expr Block ("else" (IfExpr | Block))?

MatchExpr ::= "match" Expr "{" MatchArm* "}"
MatchArm ::= Pattern "=>" Expr ","?

Pattern ::= Literal
          | Ident
          | StructPattern
          | EnumPattern
          | "_"

StructPattern ::= Ident "{" FieldPatternList? "}"
FieldPatternList ::= FieldPattern ("," FieldPattern)*
FieldPattern ::= Ident (":" Pattern)?

EnumPattern ::= Ident ("(" PatternList ")")?
PatternList ::= Pattern ("," Pattern)*

BlockExpr ::= Block

StructLiteral ::= Ident "{" FieldInitList? "}"
FieldInitList ::= FieldInit ("," FieldInit)* ","?
FieldInit ::= Ident ":" Expr

ArrayLiteral ::= "[" ExprList? "]"
ExprList ::= Expr ("," Expr)* ","?

MapLiteral ::= "Map" "{" MapEntryList? "}"
MapEntryList ::= MapEntry ("," MapEntry)* ","?
MapEntry ::= Expr ":" Expr

FieldAccess ::= Expr "." Ident

IndexExpr ::= Expr "[" Expr "]"

ParenExpr ::= "(" Expr ")"
```

### 2.6 Operators

```bnf
BinaryOp ::= "+" | "-" | "*" | "/" | "%"
           | "==" | "!=" | "<" | ">" | "<=" | ">="
           | "&&" | "||"

UnaryOp ::= "-" | "!"
```

---

## 3. Operator Precedence

From highest to lowest precedence:

| Precedence | Operator(s)            | Associativity | Description                   |
| ---------- | ---------------------- | ------------- | ----------------------------- | ---- | ---------- |
| 10         | `.`, `[]`, `()`        | Left          | Field access, indexing, call  |
| 9          | `-`, `!` (unary)       | Right         | Unary minus, logical not      |
| 8          | `*`, `/`, `%`          | Left          | Multiplication, division, mod |
| 7          | `+`, `-`               | Left          | Addition, subtraction         |
| 6          | `<`, `>`, `<=`, `>=`   | Left          | Comparison                    |
| 5          | `==`, `!=`             | Left          | Equality                      |
| 4          | `&&`                   | Left          | Logical AND                   |
| 3          | `                      |               | `                             | Left | Logical OR |
| 2          | `=`                    | Right         | Assignment                    |
| 1          | `if`, `match`, `block` | N/A           | Control flow expressions      |

---

## 4. Type System

### 4.1 Primitive Types

| Type       | Description                     | Size     | Example                 |
| ---------- | ------------------------------- | -------- | ----------------------- |
| `Int`      | Signed 64-bit integer           | 8 bytes  | `42`, `-100`, `0xFF`    |
| `Float`    | IEEE 754 double-precision float | 8 bytes  | `3.14`, `1.5e-10`       |
| `Bool`     | Boolean value                   | 1 byte   | `true`, `false`         |
| `String`   | UTF-8 encoded string            | Variable | `"hello"`               |
| `Time`     | Absolute timestamp (UTC)        | 8 bytes  | `@2024-01-15T10:30:00Z` |
| `Duration` | Time duration                   | 8 bytes  | `10s`, `5m`, `90d`      |

### 4.2 Composite Types

**Struct (Product Type):**

```uad
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}
```

**Enum (Sum Type):**

```uad
enum JudgeKind {
  Human,
  Pipeline,
  Model(String),  // Variant with data
  Hybrid(String, Float),
}
```

**Array:**

```uad
let numbers: [Int] = [1, 2, 3, 4, 5];
let matrix: [[Float]] = [[1.0, 2.0], [3.0, 4.0]];
```

**Map:**

```uad
let scores: Map[String, Int] = Map {
  "alice": 95,
  "bob": 87,
};
```

**Function Type:**

```uad
let compute: fn(Float, Float) -> Float = fn(x, y) { x + y };
```

### 4.3 Type Inference Rules

**.uad-core uses bidirectional type checking:**

**1. Synthesis (⇒): Expression produces a type**

```
Γ ⊢ e ⇒ τ
```

**2. Checking (⇐): Expression checked against expected type**

```
Γ ⊢ e ⇐ τ
```

**Inference Rules:**

**[Literal-Int]**

```
─────────────────
Γ ⊢ n ⇒ Int
```

**[Literal-Float]**

```
─────────────────
Γ ⊢ f ⇒ Float
```

**[Literal-Bool]**

```
─────────────────────
Γ ⊢ true/false ⇒ Bool
```

**[Var]**

```
(x : τ) ∈ Γ
─────────────
Γ ⊢ x ⇒ τ
```

**[App]**

```
Γ ⊢ f ⇒ fn(τ₁) -> τ₂    Γ ⊢ e ⇐ τ₁
────────────────────────────────────
Γ ⊢ f(e) ⇒ τ₂
```

**[BinOp-Arith]**

```
op ∈ {+, -, *, /, %}
Γ ⊢ e₁ ⇒ Int    Γ ⊢ e₂ ⇒ Int
────────────────────────────────
Γ ⊢ e₁ op e₂ ⇒ Int
```

**[BinOp-Cmp]**

```
op ∈ {<, >, <=, >=, ==, !=}
Γ ⊢ e₁ ⇒ τ    Γ ⊢ e₂ ⇒ τ
────────────────────────────────
Γ ⊢ e₁ op e₂ ⇒ Bool
```

**[If]**

```
Γ ⊢ cond ⇐ Bool    Γ ⊢ then ⇒ τ    Γ ⊢ else ⇒ τ
────────────────────────────────────────────────
Γ ⊢ if cond then else ⇒ τ
```

**[Let]**

```
Γ ⊢ e ⇒ τ    Γ, x:τ ⊢ body ⇒ τ'
────────────────────────────────
Γ ⊢ let x = e; body ⇒ τ'
```

**[Struct-Literal]**

```
struct S { f₁:τ₁, ..., fₙ:τₙ } ∈ Γ
Γ ⊢ e₁ ⇐ τ₁    ...    Γ ⊢ eₙ ⇐ τₙ
─────────────────────────────────────
Γ ⊢ S { f₁:e₁, ..., fₙ:eₙ } ⇒ S
```

**[Field-Access]**

```
Γ ⊢ e ⇒ S    (f : τ) ∈ fields(S)
────────────────────────────────
Γ ⊢ e.f ⇒ τ
```

---

## 5. Semantic Rules

### 5.1 Variable Binding & Scoping

**Lexical Scoping:**

- .uad-core uses lexical (static) scoping
- Variables are visible from point of declaration until end of block
- Inner scopes can shadow outer scopes

**Example:**

```uad
let x = 10;           // Scope: entire module
{
  let x = 20;         // Shadows outer x
  print(x);           // Prints: 20
}
print(x);             // Prints: 10
```

**Rules:**

1. All variables must be declared before use
2. Variables are immutable by default (future: `mut` keyword for mutability)
3. Function parameters are in function scope
4. Block expressions create new scopes

### 5.2 Function Call Semantics

**Call-by-Value:**

- All arguments are evaluated left-to-right before function call
- Arguments are copied (value semantics for primitives)
- Structs and complex types passed by reference internally (implementation detail)

**Example:**

```uad
fn add(x: Int, y: Int) -> Int {
  return x + y;
}

let result = add(10, 20);  // result = 30
```

**Recursion:**

- Functions can call themselves (recursion is supported)
- Tail-call optimization is implementation-defined (future)

**Example:**

```uad
fn factorial(n: Int) -> Int {
  if n <= 1 {
    return 1;
  } else {
    return n * factorial(n - 1);
  }
}
```

### 5.3 Pattern Matching Semantics

**Match Expression:**

- Exhaustiveness checking: all cases must be covered
- Patterns are checked top-to-bottom
- First matching pattern wins

**Example:**

```uad
enum Result {
  Ok(Int),
  Err(String),
}

fn handle(r: Result) -> String {
  match r {
    Ok(value) => "Success: " + value,
    Err(msg) => "Error: " + msg,
  }
}
```

**Pattern Types:**

1. **Literal patterns:** `42`, `"hello"`, `true`
2. **Variable patterns:** `x`, `value`
3. **Wildcard pattern:** `_` (matches anything, discards value)
4. **Struct patterns:** `Point { x, y }`
5. **Enum patterns:** `Some(value)`, `None`

**Exhaustiveness:**

```uad
// Compile error: non-exhaustive match
match x {
  0 => "zero",
  1 => "one",
  // Missing: other integers
}

// Fixed:
match x {
  0 => "zero",
  1 => "one",
  _ => "other",  // Catch-all
}
```

### 5.4 Time & Duration Operations

**Duration Arithmetic:**

```uad
let d1: Duration = 10s;
let d2: Duration = 5m;
let total: Duration = d1 + d2;  // 5m10s

let doubled: Duration = d1 * 2;  // 20s
```

**Time Arithmetic:**

```uad
let now: Time = @2024-01-15T10:00:00Z;
let later: Time = now + 10s;     // Add duration to time
let duration: Duration = later - now;  // Subtract times
```

**Comparison:**

```uad
let d1 = 10s;
let d2 = 5m;
if d1 < d2 {
  print("d1 is shorter");
}
```

**Conversion (built-in functions):**

```uad
let seconds: Int = duration_to_seconds(10m);  // 600
let dur: Duration = seconds_to_duration(300);  // 5m
```

### 5.5 Block Expressions

**Blocks are expressions:**

- The last expression in a block is its value (if no semicolon)
- Empty blocks have type `()` (unit type, future)

**Example:**

```uad
let x = {
  let a = 10;
  let b = 20;
  a + b  // No semicolon: this is the block's value
};  // x = 30

let y = {
  let a = 10;
  a + 20;  // Semicolon: statement, not value
};  // y = () (unit)
```

### 5.6 Control Flow

**If as Expression:**

```uad
let x = if condition { 10 } else { 20 };
```

**While Loop:**

```uad
let mut i = 0;  // Future: mut keyword
while i < 10 {
  print(i);
  i = i + 1;
}
```

**For Loop (Future):**

```uad
for item in array {
  print(item);
}
```

---

## 6. Domain-Specific Types

### 6.1 Action Type

```uad
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
  timestamp: Time,
}
```

### 6.2 Judge Type

```uad
struct Judge {
  kind: JudgeKind,
  decision: Float,
  confidence: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model(String),
  Hybrid(String, String),
}
```

### 6.3 Agent Type (Future)

```uad
struct Agent {
  id: String,
  role: AgentRole,
  capability: Float,
}

enum AgentRole {
  Attacker,
  Defender,
  Analyst,
}
```

---

## 7. Built-in Functions

### 7.1 I/O Functions

```uad
fn print(s: String) -> ()
fn println(s: String) -> ()
fn read_line() -> String  // Future
```

### 7.2 Math Functions

```uad
fn abs(x: Float) -> Float
fn sqrt(x: Float) -> Float
fn log(x: Float) -> Float
fn exp(x: Float) -> Float
fn pow(base: Float, exp: Float) -> Float
fn sin(x: Float) -> Float
fn cos(x: Float) -> Float
```

### 7.3 String Functions

```uad
fn len(s: String) -> Int
fn substring(s: String, start: Int, end: Int) -> String
fn concat(s1: String, s2: String) -> String
```

### 7.4 Array Functions

```uad
fn array_len(arr: [T]) -> Int
fn array_push(arr: [T], item: T) -> [T]
fn array_get(arr: [T], index: Int) -> T
```

### 7.5 Time Functions

```uad
fn now() -> Time
fn duration_to_seconds(d: Duration) -> Int
fn seconds_to_duration(s: Int) -> Duration
```

---

## 8. Module System (Future)

```uad
// main.uad
import "stdlib/math.uad"
import "erh/core.uad"

fn main() {
  let x = math::sqrt(16.0);
  print(x);
}
```

---

## 9. Error Handling

**Current (v0.1):**

- Runtime panics for errors
- No exception handling

**Future:**

- Result type: `Result<T, E>`
- Option type: `Option<T>`
- Try operator: `?`

```uad
// Future syntax
fn divide(a: Int, b: Int) -> Result<Int, String> {
  if b == 0 {
    return Err("Division by zero");
  }
  return Ok(a / b);
}
```

---

## 10. Memory Model

**Value Semantics:**

- Primitives (Int, Float, Bool) are copied
- Strings are immutable and reference-counted (implementation detail)
- Structs are copied by default (shallow copy)
- Arrays and Maps are reference types (copy-on-write, implementation detail)

**Ownership (Future):**

- Move semantics for large structures
- Borrow checker (inspired by Rust, but simplified)

---

## 11. Compilation Phases

1. **Lexical Analysis:** Source → Tokens
2. **Parsing:** Tokens → AST
3. **Type Checking:** AST → Typed AST
4. **IR Generation:** Typed AST → .uad-IR
5. **Optimization:** IR → Optimized IR (future)
6. **Execution:** IR → VM execution

---

## 12. Future Features

- Generics: `fn identity<T>(x: T) -> T`
- Traits: `trait Show { fn show(self) -> String }`
- Async/Await: `async fn fetch() -> Response`
- Macros: `macro_rules! debug`
- FFI: Foreign Function Interface to C/Go
- SIMD: Vector operations

---

## 13. Standard Library Structure (Future)

```
stdlib/
  core.uad       - Core types and functions
  math.uad       - Mathematical functions
  string.uad     - String manipulation
  array.uad      - Array operations
  map.uad        - Map operations
  time.uad       - Time and duration
  io.uad         - Input/output
  erh/
    core.uad     - ERH core types
    analysis.uad - ERH analysis functions
  security/
    crypto.uad   - Cryptographic primitives
    siem.uad     - SIEM functions
```
