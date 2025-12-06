````markdown
# .uad Language – Implementation Plan (for Cursor)
Version: 0.1
Owner: Dennis  
Status: Draft – Initial bring-up

---

## 0. Guiding Principles

- 先做「**最小可跑**」而不是「最完整」：
  - 目標一：`.uad-core` 有基本語法 + AST interpreter，可以印字串、做加減乘除。
  - 目標二：把「一個簡單的 ERH 計算」搬進 `.uad-core` 跑起來。
- 所有東西都往三層架構收斂：
  1. `.uad-core` 語言（parser + type checker + AST → IR）
  2. `.uad-IR` + VM
  3. `.uad-model` DSL（轉成 `.uad-core`）

建議開發順序：  
**Phase 1: .uad-core + AST interpreter → Phase 2: IR + VM → Phase 3: .uad-model → Phase 4: ERH / Security integration**

---

## 1. Repo Bootstrap

**目標：** 在 Cursor 裡有可編譯的 Go skeleton，`uadc` / `uadvm` 可以 build（但暫時不做實事）。

### Tasks

- [ ] 建立專案目錄 `uad-lang/`
- [ ] 新增 `go.mod`：

  ```bash
  go mod init github.com/dennislee928/uad-lang
````

* [ ] 建立基本目錄（可用 Cursor multi-file 新增）：

  ```text
  uad-lang/
    cmd/uadc/main.go
    cmd/uadvm/main.go
    internal/common/{errors.go,position.go}
    internal/lexer/{tokens.go,lexer.go}
    internal/ast/core_nodes.go
    README.md
    docs/{WHITEPAPER.md,LANGUAGE_SPEC.md,IR_SPEC.md,MODEL_LANG_SPEC.md}
    examples/core/hello_world.uad
  ```

* [ ] 寫最小版 `Makefile`：

  * `make build` 建 `bin/uadc` 與 `bin/uadvm`
  * `make test` 跑 `go test ./...`

---

## 2. Phase 1 – .uad-core 最小語言 + AST Interpreter

**目標：**
不碰 IR / VM，先用 AST interpreter 讓 `.uad-core` 可以：

* 定義 `fn main()`
* 印字串 `print("Hello, .uad!")`
* 做簡單算術與 `let` 綁定
* `if/else` 分支

### 2.1 Lexer 子集

**功能範圍（v0.1）**

* Token 類型：

  * `Ident`, `Int`, `Float`, `String`, `Bool`
  * 關鍵字：`fn`, `let`, `if`, `else`, `true`, `false`, `return`
  * 運算子：`+ - * / == != < > =`
  * 括號：`() { } , ;`
* 暫不支援：

  * `struct`, `enum`, `match`, module system

**Tasks**

* [ ] 在 `internal/lexer/tokens.go` 定義 `TokenType` 與 `Token`
* [ ] 在 `internal/lexer/lexer.go` 實作：

  * [ ] 跳過空白與註解
  * [ ] `lexIdentOrKeyword`
  * [ ] `lexNumber`（先只支援 Int，之後加 Float）
  * [ ] `lexString`
  * [ ] 基本運算子與分隔符 `+ - * / == != < > = ( ) { } , ;`
* [ ] 寫最小單元測試 `internal/lexer/lexer_test.go`

  * 給一小段程式碼，assert token 序列

---

### 2.2 AST 定義（.uad-core subset）

**v0.1 AST 範圍**

Expressions:

* `Ident`
* `Literal`（Int / Float / String / Bool）
* `BinaryExpr`（`+ - * / == != < >`）
* `CallExpr`
* `IfExpr`

Statements:

* `LetStmt`
* `ExprStmt`
* `ReturnStmt`

Declarations:

* `FnDecl`
* `Module`（暫時用單一檔案 module）

**Tasks**

* [ ] 在 `internal/ast/core_nodes.go`：

  * [ ] 定義 `Node`, `Expr`, `Stmt` 介面
  * [ ] 實作 `Ident`, `Literal`, `BinaryExpr`, `CallExpr`, `IfExpr`
  * [ ] 實作 `LetStmt`, `ExprStmt`, `ReturnStmt`
  * [ ] 實作 `FnDecl`, `Module`

---

### 2.3 Parser（.uad-core subset）

**語法目標（v0.1）**

```uad
fn main() {
  let x = 1 + 2 * 3;
  if x > 3 {
    print("big");
  } else {
    print("small");
  }
}
```

**Parsing Strategy**

* 手刻 recursive descent parser（類似 Crafting Interpreters / Writing an Interpreter in Go 的作法）
* 運算子優先順序：

  * 最高：`* /`
  * 其次：`+ -`
  * 再來：比較 `== != < >`

**Tasks**

* [ ] 新增 `internal/parser/core_parser.go`
* [ ] 實作：

  * [ ] `ParseModule(tokens []Token) (*ast.Module, error)`
  * [ ] `parseFnDecl`
  * [ ] `parseStmt`（支援 `let` / `return` / expression statement）
  * [ ] `parseExpr`：

    * [ ] `parsePrimary`（Literal / Ident / Call）
    * [ ] `parseBinary`（用 precedence climbing 或 Pratt parser）
    * [ ] `parseIfExpr`
* [ ] 寫簡單 parser 測試 `core_parser_test.go`：

  * [ ] 解析 `hello_world.uad` 不出錯
  * [ ] 解析一個帶 `if` 的例子

---

### 2.4 AST Interpreter（.uad-core）

**目標：**

* 執行 `Module`，找到 `fn main()`, 呼叫它
* 支援：

  * 整數運算
  * `let` 綁定在 function local scope
  * `return` 終止函式
  * 內建函式 `print`

**Tasks**

* [ ] 新增 `internal/vm/interpreter.go`（暫時用 AST interpreter，以後 VM 接手）

  * [ ] 定義 `Env`（symbol → value）
  * [ ] 定義 `Value` 型別（Int / String / Bool）
  * [ ] `EvalExpr(expr, env)`：

    * [ ] `Literal` → 對應 `Value`
    * [ ] `Ident` → 從 `Env` lookup
    * [ ] `BinaryExpr` → 根據運算子計算
    * [ ] `CallExpr`：

      * 若是 `print` → 走 stdout
      * 若是 user 函式 → 建立新 `Env`，綁定參數，eval body
    * [ ] `IfExpr` → 分支 eval
  * [ ] `ExecStmt(stmt, env)`：

    * [ ] `LetStmt` → `env.Set()`
    * [ ] `ExprStmt` → `EvalExpr`
    * [ ] `ReturnStmt` → 用 Go error or special `Return` struct 傳遞 result
* [ ] 在 `cmd/uadvm/main.go` 暫時直接：

  * 讀 `.uad` 原始碼 → lexer → parser → AST → interpreter 執行
  * 等 IR 完成之後再改成讀 `.uadir`

**驗收條件**

* [ ] `examples/core/hello_world.uad`：

  ```uad
  fn main() {
    print("Hello, .uad core!");
  }
  ```

  終端機輸出 `"Hello, .uad core!"`

* [ ] `examples/core/erh_basic.uad`（之後補）可以計算一組 numbers 的平均或簡易錯誤率。

---

## 3. Phase 2 – .uad-IR + VM

**目標：**
把 interpreter 後端改成 IR + VM，建立真正的 `.uad-IR` 層。

### 3.1 IR 定義

**Tasks**

* [ ] 在 `internal/ir/ir.go` 定義：

  * [ ] `OpCode` enum（`OpAdd, OpSub, OpMul, OpDiv, OpEq, OpJmp, OpJmpIf, OpCall, OpRet, ...`）
  * [ ] `Instr` struct（`Op`, `Operands []int64`）
  * [ ] `Function`, `Module` 型別
* [ ] 在 `internal/ir/builder.go` 定義 AST → IR lowering：

  * [ ] `BuildModule(ast.Module) (*ir.Module, error)`
  * [ ] `buildFnDecl`, `buildExpr`, `buildStmt` 等 helper

### 3.2 VM Loop

**Tasks**

* [ ] 在 `internal/vm/vm.go`：

  * [ ] 定義 `VM` struct：stack, call stack, instruction pointer
  * [ ] `func (v *VM) Run(mod *ir.Module) error`

    * [ ] 依序執行指令：

      * `OpConst`, `OpLoad`, `OpStore`, `OpAdd`, …, `OpJmp`, `OpJmpIf`, `OpCall`, `OpRet`
  * [ ] 在 VM runtime 中實作 `print` builtin（透過特殊 function id / opcode）
* [ ] 調整 `cmd/uadc`：

  * 由 `.uad` → AST → type-check → IR → `out.uadir`（暫時可用簡單 JSON 格式）
* [ ] 調整 `cmd/uadvm`：

  * 讀 `.uadir` → `ir.Module` → VM.Run()

**驗收條件**

* [ ] 原本用 AST interpreter 跑的 `hello_world.uad` 現在改用：

  * `uadc` 編譯成 `out.uadir`
  * `uadvm out.uadir` 執行 → 仍能印 `"Hello, .uad core!"`

---

## 4. Phase 3 – .uad-model DSL & Desugaring

**目標：**
讓你可以用 `.uad-model` 描述 ERH profile / scenario，compiler 負責產生 `.uad-core`。

### 4.1 Model AST & Parser

**Tasks**

* [ ] 在 `internal/ast/model_nodes.go` 定義：

  * `ModelModule`, `ActionClassDecl`, `JudgeDecl`, `ErhProfileDecl`, `ScenarioDecl`, `CognitiveSiemDecl` 等
* [ ] 在 `internal/parser/model_parser.go`：

  * 寫一個簡化 parser，先支援：

    * `action_class` + `judge` + `erh_profile`
  * 不急著支援 scenario / cognitive_siem，全都留 TODO

### 4.2 Desugar: Model → Core

**Tasks**

* [ ] 在 `internal/model/desugar.go`：

  * [ ] `func DesugarModel(m *ast.ModelModule) (*ast.Module, error)`
  * 實作策略（v0.1）：

    * `action_class` → 對應到一個 `struct Action_XXX` 定義（之後可直接手寫成 core AST）
    * `judge` → 轉成 `fn judge_XXX(a: Action_XXX) -> Float`
    * `erh_profile` → 轉成一組函式：

      * `fn run_profile_X(dataset: Dataset) -> ProfileResult`
* [ ] 在 `cmd/uadc/compile()`：

  * 若副檔名是 `.uadmodel`：

    * lexer + model parser → Model AST → `model.Desugar` → Core AST → Type-check → IR

**驗收條件**

* [ ] `examples/model/devsecops_erh.uadmodel` 能被 `uadc` 編譯成 IR，不爆炸。
* [ ] 暫時不要求算出真正 α，只要 pipeline 通了就算 Phase 3 達標。

---

## 5. Phase 4 – ERH / Security Integration

**目標：**
把你現有的 ERH / ERH_ON_SECURITY_POC 逐步搬進 `.uad` 生態，變成可跑的實驗。

### 5.1 ERH Library

**Tasks**

* [ ] 在 `runtime/stdlib/erh.uad`：

  * [ ] 定義 `struct Action`, `struct Judge`, `fn is_mistake(...)`, `fn is_prime(...)`
  * [ ] 寫一個簡單 `fn compute_alpha(data: [ActionJudgePair]) -> Float`
* [ ] 在 `.uad-model`：

  * [ ] 定義一個 `erh_profile` 對應到 GitLab MR / Security Logs 的示範資料
* [ ] 寫一個 `examples/core/erh_devsecops.uad`：

  * 用 hard-coded data 跑 `compute_alpha`

### 5.2 Security / Cyber Range / SIEM（之後）

* [ ] 先寫 spec / TODO，不急著實作：

  * `scenario`, `cognitive_siem` 的語義與執行模型
* [ ] 等 `.uad-core` 跑穩，再接你的 AI Adversarial Cyber Range & Cognitive SIEM 設計。

---

## 6. 開發習慣建議（for Cursor flow）

* 每一個 Phase 開一個 branch：

  * `feature/core-interpreter`
  * `feature/ir-vm`
  * `feature/model-dsl`
* 每一個大檔案（lexer/parser/vm）：

  * 先寫 TODO skeleton，再逐步用 Cursor 填空。
* 優先讓：

  * `examples/core/hello_world.uad`
  * `examples/core/erh_basic.uad`

  這兩個 example 先跑起來，再往複雜的 ERH / scenario 前進。

---

## 7. Checkpoint Summary

可以當成你在 Cursor 裡的「階段里程碑」 checklist：

* [ ] Phase 1：

  * [ ] Lexer subset OK
  * [ ] Core AST + parser OK
  * [ ] AST interpreter 跑 `hello_world.uad`
* [ ] Phase 2：

  * [ ] IR spec + builder OK
  * [ ] VM loop OK
  * [ ] `uadc + uadvm` pipeline 跑 `hello_world.uad`
* [ ] Phase 3：

  * [ ] Model AST + parser OK
  * [ ] Desugar Model → Core OK
  * [ ] `devsecops_erh.uadmodel` 可被編譯成 IR
* [ ] Phase 4：

  * [ ] ERH stdlib 初版
  * [ ] ERH DevSecOps example 可跑出基本統計（即使 α 先用簡化版）

---

```
::contentReference[oaicite:1]{index=1}
```
