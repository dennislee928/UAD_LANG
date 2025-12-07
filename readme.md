# UAD Programming Language

<div align="center">

**Unified Adversarial Dynamics Language**  
*Domain-Specific Language for Adversarial Modeling, Ethical Risk, and Cognitive Security*

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/)

[Features](#features) • [Quick Start](#quick-start) • [Documentation](#documentation) • [Examples](#examples) • [Contributing](#contributing)

</div>

---

## 概述 (Overview)

UAD (.uad) 是一門專為 **對抗動態建模** (Adversarial Dynamics)、**道德風險量化** (Ethical Risk Hypothesis) 與 **認知安全系統** (Cognitive Security) 設計的專用語言。

與通用程式語言不同，UAD 將 **決策、風險和時間** 視為一等公民，提供原生的語義支持來描述複雜的對抗行為、量化道德後果，並模擬長期系統演化。

### 核心特性

- ✅ **強型別系統**：靜態型別檢查，確保模型正確性
- ✅ **時間語義**：樂理式 DSL（Score/Track/Bar/Motif）用於精確時間控制
- ✅ **對抗建模**：內建 Action/Judge/Agent 概念
- ✅ **弦理論語義**：描述複雜的耦合與共振關係
- ✅ **量子糾纏**：變數同步機制
- ✅ **雙模執行**：解釋器 + 虛擬機
- ✅ **Dev Container**：一鍵開發環境

完整白皮書請見 **[docs/WHITEPAPER.md](docs/WHITEPAPER.md)**

---

## 快速開始 (Quick Start)

### 前置需求

- Go 1.21 或更高版本
- Make
- Git

### 安裝與構建

```bash
# 克隆專案
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming

# 安裝依賴
go mod tidy

# 構建所有工具
make build

# 運行測試
make test

# 運行範例
make example
```

### 第一個 UAD 程式

```uad
// hello.uad
fn main() -> Int {
    println("Hello, UAD!");
    return 0;
}
```

---

## 文件 (Documentation)

### 📘 核心文件

- **[PARADIGM.md](docs/PARADIGM.md)** - UAD 語言範式與設計哲學
- **[SEMANTICS_OVERVIEW.md](docs/SEMANTICS_OVERVIEW.md)** - 語義與執行模型
- **[WHITEPAPER.md](docs/WHITEPAPER.md)** - 完整白皮書

### 📋 規格文件

- **[CORE_LANGUAGE_SPEC.md](docs/specs/CORE_LANGUAGE_SPEC.md)** - 核心語言規格
- **[MODEL_DSL_SPEC.md](docs/specs/MODEL_DSL_SPEC.md)** - 高階 DSL 規格
- **[IR_SPEC.md](docs/specs/IR_SPEC.md)** - 中間表示規格

### 🛠️ 開發文件

- **[CONTRIBUTING.md](CONTRIBUTING.md)** - 貢獻指南
- **[ROADMAP.md](docs/ROADMAP.md)** - 發展藍圖

---

## 語言特性 (Features)

### 時間語義 (M2.3)

```uad
score AttackSimulation {
    tempo: 120,
    track Attacker {
        bars 1..4 { motif reconnaissance; }
    }
}
```

### 弦理論語義 (M2.4)

```uad
string EthicalField {
    modes { integrity: Float }
}

coupling EthicalField.integrity {
    mode_pair (integrity, resilience) with strength 0.7;
}
```

### 量子糾纏 (M2.5)

```uad
let x: Int = 10;
let y: Int = 20;
entangle x, y;  // 共享相同的值
```

---

## 開發進度 (Status)

### ✅ 已完成

- [x] 語言核心 (Lexer, Parser, Type Checker, Interpreter, VM)
- [x] M2.3: 樂理 DSL (15 測試 ✅)
- [x] M2.4: 弦理論語義 (14 測試 ✅)
- [x] M2.5: 糾纏語義 (15 測試 ✅)
- [x] 測試覆蓋率: ~80% (121 測試)
- [x] CI/CD (GitHub Actions)
- [x] 文件系統完整

### 🚧 進行中

- [ ] M6: 實驗框架
- [ ] LSP 實作

---

## 貢獻 (Contributing)

我們歡迎所有形式的貢獻！請參閱 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 授權 (License)

Apache License 2.0 - 詳見 [LICENSE](LICENSE)

---

<div align="center">

**Made with ❤️ by the UAD Community**

</div>
