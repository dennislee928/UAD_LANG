# UAD_Programming – Cursor Plan

> Goal: 把 UAD 語言變成一個「雲端優先、可實驗、可發 paper」的語言專案：
>
> - 有乾淨的語法 / AST / VM 架構
> - 有明確的「樂理 + 弦理論 + 共振/糾纏」語義
> - 有標準化 repo 結構、CI/CD、測試與實驗管線

---

## 0. Meta – 使用方式

這份 `plan.md` 是給 Cursor / VS Code AI 使用的工作藍圖。

**原則：**

- 每次只處理「一個 Milestone」或其中的小段，避免一次修改太多檔案。
- 優先「整理現有東西」，再逐步加新功能。
- 對每個大變動，都要：
  - 保留原檔案的主要邏輯
  - 補上最基本的測試（即使是簡單 smoke test）
- 每次修改都要更新 `plan.md` 中的 Repo Snapshot，方便之後參考。
- 所有的 comments 都要用英文。
- 每做一個檔案修正就 commit 到 branch dev,每個 phase 結束後要做測試，測試成功後就 push 到 github 的 branch dev。然後建立 pull request 到 main branch
- 每次 pull request 都要寫清楚 commit message，描述修改的內容。
- 每次 pull request 都要寫清楚 description，描述修改的內容。
- 每次 pull request 都要寫清楚 labels，描述修改的內容。
- 每次 pull request 都要寫清楚 assignees，描述修改的內容。
- 每次 pull request 都要寫清楚 reviewers，描述修改的內容。
- 每次 pull request 都要寫清楚 milestone，描述修改的內容。
- 每次 pull request 都要寫清楚 project，描述修改的內容。
- 每次 pull request 都要寫清楚 dependencies，描述修改的內容。
- 每次 pull request 都要寫清楚 related issues，描述修改的內容。
- 每次 pull request 都要寫清楚 related pull requests，描述修改的內容。

---

## 1. High-Level Goals（高階目標）

1. 建立清楚的語言核心：
   - 一級概念：`score / bar / chord / agent / resonance / entangle`
   - AST 與 IR 的層次清楚（parser / semantic / runtime 分離）
2. 建立穩定的工具鏈：
   - 統一 build 指令（`make build`, `make test`, `make docs`）
   - devcontainer / Codespaces 支援
   - GitHub Actions CI（build + test）
3. 為「樂理 + 弦理論 + 共振/糾纏」語義留好鉤子：
   - 有 core types / modules，而不一定一次寫完全部實作
4. 為未來 paper / whitepaper 準備：
   - `docs/` / `spec/` 內有足夠清楚的語義文件
   - `examples/` 有代表性的 UAD 程式

---

## 2. Milestone 0 – Repo 快速盤點與基本清理

**目標**：讓專案狀態可被 Cursor/AI 安全操作，有最基本的一致性。

### 2.1 掃描現有結構

**Tasks for Cursor:**

- 列出目前 repo 的：
  - 頂層檔案：`README`, `main.*`, `src`, `parser` 等
  - 語言實作位置（parser / interpreter 現在在哪）
- 在 `plan.md` 下方另開一個 section：`## Repo Snapshot`，貼上簡短摘要：
  - 目前使用語言（Go / Rust / C++ / Python …）
  - 主要 entrypoint 檔案 / 資料夾
  - 是否已存在 tests / docs

### 2.2 最小清理（不改語義）

**Tasks for Cursor:**

- 若有明顯命名不一致（例如 `uad_lang` vs `UAD` 混用），在不影響 import 的前提下做輕微整理：
  - 只改註解 / README 敘述，不動 code 邏輯。
- 確認 `.gitignore` 是否存在，如沒有，建立一個基本版（忽略 build 產物、IDE 檔案等）。

---

## 3. Milestone 1 – 語言核心：Parser / AST / Runtime Skeleton

**目標**：把「語法 → AST → VM/Interpreter」的層級整理清楚，為後續樂理與弦論語義打底。

### 3.1 AST 統整與 Type 設計

**Concept target：**

- 必備 AST 節點：
  - `ScoreNode`
  - `BarNode`
  - `ChordNode`
  - `AgentDeclNode`
  - `ResonanceNode`
  - （之後預留）`EntangleNode`, `SyncNode`

**Tasks for Cursor:**

1. 找到目前 AST 定義檔（如 `ast.go`, `ast.rs`, `ast.ts` 等）。
2. 如果 AST 分散在多檔案：
   - 建議建立 `src/ast/` 或單一 `src/ast.<ext>`，集中定義核心節點。
3. 為上述 5–6 類核心 node 定義乾淨的型別/enum/struct：
   - 要有：
     - 位置信息（line/column）方便錯誤訊息
     - 子節點列表
     - semantic payload（例如 bar 的 duration、chord 的元素等）
4. 在 AST 型別上加精簡註解（英文）說明用途。

### 3.2 Parser 分層

**Concept target：**

- Parser 分為：
  - lexer/tokenizer
  - parser（產 AST）
  - error reporter

**Tasks for Cursor:**

1. 找出目前 parser 的實作檔。
2. 若 parser 與 runtime 混在同一檔案：
   - 重構成：
     - `src/parser/lexer.*`
     - `src/parser/parser.*`
3. 確保 parser interface 清楚：
   - 例如：`Parse(source: string) -> ScoreNode | Error`.
4. 新增最基本 parser 測試檔：
   - 測 `score {}`, `bar {}`, `chord [...]` 至少三種語法。
   - 確認 AST shape 正確。

### 3.3 Runtime / VM skeleton

**Concept target：**

- 即使只 prototype，也要有清楚分層：
  - `Evaluator` / `Interpreter` class / module
  - 之後才能掛上「時間格點 / 調度器」

**Tasks for Cursor:**

1. 找到解譯或執行邏輯檔（可能叫 `eval`, `vm`, `runtime` 等）。
2. 重構出：
   - `src/runtime/core.*`：執行 AST 的最小邏輯
   - interface 例如：`RunScore(score: ScoreNode, options: RuntimeConfig)`
3. 暫時將「時間」實作留簡單版本：
   - 單純照 AST 順序執行，但預留 `tick()` or `advanceBeat()` 介面。

---

## 4. Milestone 2 – 注入「樂理 + 弦理論」語義骨架

**目標**：先把概念接口放進去，不一定一次實作完所有細節。

### 4.1 樂理層：Score / Bar / Chord 語意

**Concept target：**

- `Score`：整體 program / simulation。
- `Bar`：時間區間（帶拍數/長度）。
- `Chord`：同一時間格點上要「一起」發生的動作。

**Tasks for Cursor:**

1. 在 `src/runtime/core.*` 中：
   - 引入明確的 `ScoreRuntime`, `BarRuntime`, `ChordRuntime` 結構或類別。
2. 執行順序調整為：
   - `ScoreRuntime.Run()`：
     - 依序走過每個 `BarNode`
     - 在每個 bar 中，依序處理 chord / events
3. 在 `BarRuntime` 中加時間欄位（如 `beats: int`），即使暫時沒用到，也保留。

### 4.2 弦理論層：Resonance（共振）骨架

**Concept target：**

- `ResonanceGraph`：說明哪些變數/agent 屬性之間有耦合。
- runtime 執行時，只要某個 node 更新，就通知 graph。

**Tasks for Cursor:**

1. 建立新檔案：`src/runtime/resonance.*`
   - 定義：
     - `ResonanceLink { source, target, mode }`
     - `ResonanceGraph { links: [] }`
2. 在 AST 的 `ResonanceNode` 對應：
   - 寫一個 semantic pass：從 AST 建出 `ResonanceGraph`
3. 在 runtime：
   - 對任何「狀態更新」操作（例如 `AgentState.x = ...`）包一個 helper：
     - `SetValueWithResonance(entity, field, value, graph)`
   - 目前可以先寫空的 / no-op，再逐版填。

### 4.3 Entanglement（量子糾纏）接口預留

**Concept target：**

- `EntanglementGroup`：一組變數 / field 共用同一底層值。

**Tasks for Cursor:**

1. 建立 `src/runtime/entanglement.*`
   - 定義：
     - `EntanglementGroup { id, members }`
     - helper：`BindEntangled(vars...)`
2. 在 AST 增加 `EntangleNode`（如果尚未有）：
   - 後續 semantic pass 用來建立 entanglement group。
3. 現階段實作可以先很簡單：
   - 給 entangled 變數一個 shared backing store（如一個 map）。

---

## 5. Milestone 3 – Repo 結構 / build / devcontainer

**目標**：讓專案「雲端優先」，打開 Codespaces / Cursor 就能 build & test。

### 5.1 建標準目錄

**Tasks for Cursor:**

- 若尚未有以下目錄，建立：
  - `src/`
  - `examples/`
  - `tests/`
  - `docs/`
  - `spec/`
- 把現有實作檔按功能遷移：
  - parser → `src/parser/`
  - runtime → `src/runtime/`
  - AST → `src/ast/`

### 5.2 Build script 統一

**Tasks for Cursor:**

1. 在 repo 根目錄新增 `Makefile`（或 `justfile`，擇一）：
   - `make build`：編譯主程式 / library
   - `make test`：跑全部測試
   - `make example`：編譯/執行一個範例（from `examples/`）
2. 在 README 裡新增一節：
   - 「Build & Run」
   - 告訴使用者 `make build && make test`。

### 5.3 devcontainer / Codespaces 設定

**Tasks for Cursor:**

1. 新增 `.devcontainer/devcontainer.json`：
   - 指向一個簡單 Docker image（例如基於官方語言 image）
   - 裝好必要 toolchain
2. 若需要，新增 `.devcontainer/Dockerfile`：
   - 安裝 build 工具（cmake / cargo / go / python …）
3. 在 README 補一句：
   - 「You can open this repo directly in GitHub Codespaces / VS Code Dev Containers.」

---

## 6. Milestone 4 – 測試 / CI / 實驗骨架

**目標**：讓語言每次修改都有最低限度的自動回歸保護。

### 6.1 測試最小集

**Tasks for Cursor:**

1. 在 `tests/` 內：
   - 新增 `parser_basic_test.*`：
     - 測 `score {}`, `bar {}`, `chord [ .. ]` → AST OK
   - 新增 `runtime_basic_test.*`：
     - 跑一個最簡單的 score，檢查執行無 panic / error。
2. 如果有 testing framework（go test / pytest / rust test），統一用一種。

### 6.2 GitHub Actions CI

**Tasks for Cursor:**

1. 新增 `.github/workflows/ci.yml`：
   - 觸發：`push` + `pull_request`
   - steps：
     - checkout
     - `make build`
     - `make test`
2. （可選）加上多平台 matrix（Linux / macOS）。

---

## 7. Milestone 5 – Docs / Spec / Concept Clarification

**目標**：把「UAD 在幹嘛」講清楚，方便未來寫 paper / whitepaper。

### 7.1 語言概念文件

**Tasks for Cursor:**

1. 在 `docs/` 新增：
   - `PARADIGM.md`：
     - 說明：
       - UAD 的一級概念：`score / bar / chord / agent / resonance / entangle`
       - 與 Go / Rust / Python 的基本差異（function vs chord, struct vs agent…）
     - 中英都可，至少標題與小節用英文。
   - `SEMANTICS_OVERVIEW.md`：
     - 描述：
       - Temporal harmony grid（時間格點）怎麼調度
       - Resonance graph & Entanglement group 的概念
2. 在 `spec/` 新增：
   - `LANGUAGE_SPEC_DRAFT.md`：
     - 收斂目前已有語法（不是 fully formal，但可當草稿）

### 7.2 README 精簡版本

**Tasks for Cursor:**

- 更新 `README.md`：
  - 加上：
    - 一段短 introductions：
      - UAD = music & string-theory inspired PL for time-structured, multi-agent systems
    - 簡單語法示例（score + bar + chord）
    - `Build`, `Run`, `Tests`, `Devcontainer` 用法

---

## 8. Milestone 6 – 實驗框架（ERH / Security / Psychohistory friendly）

**目標**：為未來 ERH / security / psychohistory 模擬留出「標準實驗入口」。

### 8.1 experiments/ 結構

**Tasks for Cursor:**

1. 建 `experiments/` 目錄：
   - e.g. `experiments/erh_demo.uad`
   - `experiments/configs/erh_demo.yaml`（描述跑幾次、哪些參數）
2. 寫一個簡單 CLI：
   - `cmd/uad-runner`（或類似）：
     - `uad-runner --config experiments/configs/erh_demo.yaml`
     - 讀 config → 讀 .uad → 執行 → 輸出 JSON 結果

### 8.2 （可選）CI 實驗 Workflow

**Tasks for Cursor:**

- 新增 `.github/workflows/experiments.yml`（可用 `workflow_dispatch` 觸發）：
  - build
  - 跑 `uad-runner` 幾個 demo config
  - 結果上傳為 artifact

---

## 9. Milestone 7 – 之後的進階（可暫時標 TODO）

**未來可以做，但現在先只在 plan 中標記：**

- 共振 / 糾纏 → 自動 refactor + pipeline routing
- WebAssembly backend（讓 UAD 在瀏覽器跑）
- VS Code / Cursor 語言 server（LSP）：
  - 語法高亮
  - AST/Resonance graph 可視化

---

## Repo Snapshot（由 Cursor 之後填）

> 由 Cursor 在第一次掃完 repo 後填寫，方便之後參考。

- 語言實作：`TODO by Cursor`（e.g. Go / Rust / C++ / Python）
- 主要 entrypoint：`TODO by Cursor`（e.g. `src/main.go`）
- Parser 檔案：`TODO by Cursor`
- Runtime 檔案：`TODO by Cursor`
- 已有測試？`TODO by Cursor`
- 已有 docs？`TODO by Cursor`
