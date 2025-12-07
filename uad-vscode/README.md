# UAD Language Support for VS Code / Cursor

VS Code / Cursor extension for the **UAD (Unified Adversarial Dynamics) Programming Language**, providing syntax highlighting and language support.

## 功能特性 (Features)

- ✅ **語法高亮** (Syntax Highlighting)：完整的關鍵字、類型、字面量、操作符高亮
- ✅ **註解支援**：單行 (`//`) 和多行 (`/* */`) 註解
- ✅ **括號配對**：自動配對和自動關閉括號、大括號、方括號
- ✅ **字串支援**：雙引號字串，支援轉義序列
- ✅ **關鍵字分類**：不同類型的關鍵字使用不同的顏色高亮

## 支援的語法特性

### 核心關鍵字

- **控制流**：`fn`, `return`, `if`, `else`, `match`, `while`, `for`, `break`, `continue`
- **宣告**：`let`, `struct`, `enum`, `type`, `import`, `module`, `pub`
- **字面量**：`true`, `false`, `nil`
- **模式匹配**：`case`, `when`, `in`

### 領域特定關鍵字

- **核心領域**：`action`, `judge`, `agent`, `event`, `emit`
- **Model DSL**：`action_class`, `erh_profile`, `scenario`, `cognitive_siem`, `from`, `dataset`, `prime_threshold`, `fit_alpha`
- **Musical DSL**：`score`, `track`, `bars`, `motif`, `use`, `at`, `variation`, `as`, `transpose`
- **String Theory DSL**：`string`, `brane`, `on`, `coupling`, `resonance`, `modes`, `with`, `strength`, `dimensions`
- **Entanglement**：`entangle`

### 原始類型

- `Int`, `Float`, `Bool`, `String`, `Time`, `Duration`

### 字面量

- **整數**：十進制 (`42`)、十六進制 (`0xFF`)、二進制 (`0b1010`)
- **浮點數**：標準格式 (`3.14`)、科學計數法 (`1.5e-10`)
- **字串**：雙引號字串 (`"hello"`)，支援轉義序列 (`\n`, `\t`, `\r`, `\\`, `\"`, `\0`, `\u{XXXX}`)
- **Duration**：時間持續字面量 (`10s`, `5m`, `2h`, `3d`)
- **Time**：時間戳字面量 (`@2024-01-15T10:30:00Z`)

### 操作符

- **算術**：`+`, `-`, `*`, `/`, `%`
- **比較**：`==`, `!=`, `<`, `>`, `<=`, `>=`
- **邏輯**：`&&`, `||`, `!`
- **特殊**：`->`, `=>`, `.`, `::`, `..`

## 安裝方法

### 方法 1：從原始碼安裝（開發模式）

1. 複製此 extension 目錄到你的本地機器
2. 打開 VS Code 或 Cursor
3. 按 `F1` 或 `Cmd+Shift+P` (Mac) / `Ctrl+Shift+P` (Windows/Linux) 打開命令面板
4. 輸入 `Extensions: Install from VSIX...` 或直接在開發模式下使用：
   - 在 VS Code 中，點擊左側的 Extensions 圖標
   - 點擊右上角的 `...` 選單，選擇 `Install from VSIX...`
   - 選擇打包好的 `.vsix` 文件

### 方法 2：直接使用（開發模式）

如果你想在開發模式下直接使用：

1. 將 `uad-vscode` 目錄複製到 VS Code 的 extensions 目錄：
   - **macOS**: `~/.vscode/extensions/` 或 `~/.cursor/extensions/`
   - **Windows**: `%USERPROFILE%\.vscode\extensions\` 或 `%USERPROFILE%\.cursor\extensions\`
   - **Linux**: `~/.vscode/extensions/` 或 `~/.cursor/extensions/`

2. 重新啟動 VS Code 或 Cursor

### 方法 3：打包為 VSIX（推薦用於分發）

1. 安裝 `vsce`（VS Code Extension Manager）：
   ```bash
   npm install -g @vscode/vsce
   ```

2. 在 `uad-vscode` 目錄中執行：
   ```bash
   cd uad-vscode
   vsce package
   ```

3. 這會生成一個 `.vsix` 文件，可以在 VS Code 或 Cursor 中安裝

## 使用說明

安裝完成後，VS Code / Cursor 會自動識別 `.uad` 文件並應用語法高亮。

### 測試語法高亮

創建一個新的 `.uad` 文件，例如 `test.uad`：

```uad
// Hello World example
fn main() {
    println("Hello, UAD!");
    
    let x: Int = 42;
    let y: Float = 3.14;
    let name: String = "UAD Language";
    
    if x > 0 {
        print("Positive number");
    } else {
        print("Negative number");
    }
    
    struct Person {
        name: String,
        age: Int,
    }
    
    let person = Person {
        name: "Alice",
        age: 25,
    };
}
```

你應該能看到：
- 關鍵字（`fn`, `let`, `if`, `struct` 等）以不同顏色高亮
- 類型（`Int`, `Float`, `String`）以類型顏色顯示
- 字串以字串顏色顯示
- 數字以數字顏色顯示
- 註解以註解顏色顯示

## 檔案結構

```
uad-vscode/
├── package.json                 # Extension 配置檔案
├── language-configuration.json  # 語言配置（註解、括號等）
├── syntaxes/
│   └── uad.tmLanguage.json      # TextMate grammar（語法高亮規則）
└── README.md                    # 本文件
```

## 相容性

- **VS Code**：版本 1.74.0 或更高
- **Cursor**：完全相容（Cursor 基於 VS Code，支援相同的 extension 格式）

## 開發與貢獻

### 修改語法高亮規則

編輯 `syntaxes/uad.tmLanguage.json` 文件。這個文件使用 TextMate grammar 格式（JSON），定義了所有語法高亮規則。

### 修改語言配置

編輯 `language-configuration.json` 文件來調整註解、括號配對等行為。

### 測試變更

1. 在開發模式下，修改文件後重新載入 VS Code / Cursor 視窗（`Cmd+R` / `Ctrl+R`）
2. 打開一個 `.uad` 文件查看效果

## 相關資源

- **UAD 語言專案**：https://github.com/dennislee928/uad-lang
- **語言規格文件**：`docs/specs/CORE_LANGUAGE_SPEC.md`
- **範例程式碼**：`examples/` 目錄

## 授權

Apache-2.0 License

## 問題回報

如有問題或建議，請在 GitHub 專案中提交 Issue。

## 版本歷史

### 0.1.0 (2024-12)

- 初始版本
- 完整的語法高亮支援
- 支援所有 UAD 語言關鍵字和特性
- 支援核心語言和所有擴展 DSL（Musical、String Theory、Entanglement、Model DSL）

