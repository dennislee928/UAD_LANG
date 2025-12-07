# UAD 語言範式 (UAD Language Paradigm)

## 目錄

1. [核心哲學](#核心哲學)
2. [語言特性](#語言特性)
3. [對抗動態學 (Adversarial Dynamics)](#對抗動態學)
4. [道德風險假說 (Ethical Risk Hypothesis, ERH)](#道德風險假說)
5. [認知安全 (Cognitive Security)](#認知安全)
6. [時間語義 (Temporal Semantics)](#時間語義)
7. [共振與糾纏 (Resonance & Entanglement)](#共振與糾纏)
8. [與其他語言的差異](#與其他語言的差異)

---

## 核心哲學

UAD (Unified Adversarial Dynamics) 語言是專為**對抗性系統建模**和**道德風險量化**而設計的特定領域語言 (DSL)。它的核心理念是：

> **「在充滿對抗的環境中，決策的道德後果不是靜態的，而是動態演化的。」**

UAD 語言提供了一套完整的語義框架，用於描述、分析和預測在對抗性情境下的行為模式、道德風險以及認知安全問題。

### 設計原則

1. **表達力優先 (Expressiveness First)**

   - 提供直觀的語法來描述複雜的對抗行為
   - 支持多層次的抽象（從基本行動到複雜策略）

2. **型別安全 (Type Safety)**

   - 強型別系統確保模型正確性
   - 編譯時檢查避免語義錯誤

3. **時間敏感 (Time-Aware)**

   - 內建時間語義，支持事件調度和時序約束
   - 樂理式的 DSL（Score/Track/Bar）用於精確時間控制

4. **多範式融合 (Multi-Paradigm)**
   - 函數式：純函數、不可變值、模式匹配
   - 物件導向：結構體、方法、封裝
   - 宣告式：高階 DSL 用於描述模型結構

---

## 語言特性

### 1. 核心概念三元組：Action, Judge, Agent

UAD 的世界觀建立在三個基本概念上：

```uad
// Action: 表示一個具體的行為或決策
struct Action {
    id: String,
    complexity: Float,    // 行動複雜度
    importance: Float,    // 道德重要性
    timestamp: Duration,
}

// Judge: 評估 Action 的道德後果
struct Judge {
    erh_score: Float,     // 道德風險假說評分
    confidence: Float,    // 評估置信度
    reasoning: String,    // 推理過程
}

// Agent: 執行 Action 並產生 Judge
struct Agent {
    name: String,
    strategy: String,     // 決策策略
    history: [Action],    // 行動歷史
}
```

### 2. 型別系統

- **基本型別**：`Int`, `Float`, `Bool`, `String`, `Duration`, `Nil`
- **複合型別**：結構體 (struct)、列舉 (enum)、陣列 `[T]`、映射 `Map[K, V]`
- **函數型別**：`fn(T1, T2) -> T3`
- **特殊型別**：`Time`, `TimeRange`, `Score`, `String` (弦場)

### 3. 模式匹配與決策

```uad
fn evaluate_action(action: Action) -> Judge {
    match action.complexity {
        low if action.importance < 0.5 => Judge {
            erh_score: 0.1,
            confidence: 0.9,
            reasoning: "Low risk, low importance",
        },
        high if action.importance > 0.8 => Judge {
            erh_score: 0.8,
            confidence: 0.7,
            reasoning: "High risk due to complexity and importance",
        },
        _ => Judge {
            erh_score: 0.5,
            confidence: 0.5,
            reasoning: "Moderate risk",
        },
    }
}
```

---

## 對抗動態學

對抗動態學是 UAD 語言的核心應用領域。它描述在多智能體環境中，行為者之間的策略互動和演化過程。

### 核心方程

在 UAD 中，對抗性互動可以用以下概念建模：

1. **行動空間 (Action Space)**：所有可能的行動
2. **策略空間 (Strategy Space)**：決策規則的集合
3. **效用函數 (Utility Function)**：評估行動結果的函數
4. **納許均衡 (Nash Equilibrium)**：無人有動機單方面改變策略的狀態

```uad
// 定義策略
fn strategy_optimal(state: GameState, agent: Agent) -> Action {
    let candidates = generate_actions(state);
    let best = candidates
        .map(|a| (a, evaluate_utility(a, agent)))
        .max_by(|_, u| u);
    return best.0;
}

// 模擬多回合對抗
fn simulate_adversarial_game(agents: [Agent], rounds: Int) -> GameHistory {
    let mut history = GameHistory::new();
    for round in 0..rounds {
        for agent in agents {
            let action = agent.strategy(history.current_state(), agent);
            history.record(action);
            history.update_state(action);
        }
    }
    return history;
}
```

---

## 道德風險假說

**道德風險假說 (Ethical Risk Hypothesis, ERH)** 是 UAD 語言的理論基礎，認為：

> **在對抗性環境中，行為者往往會選擇「道德風險較高但短期效用較大」的策略，從而導致整體系統的道德退化。**

### ERH 評分機制

UAD 提供內建的 ERH 評分系統：

```uad
fn compute_erh_score(action: Action, context: Context) -> Float {
    let complexity_penalty = action.complexity * 0.3;
    let importance_weight = action.importance * 0.4;
    let history_factor = context.agent.history
        .filter(|a| a.erh_score > 0.5)
        .count() as Float * 0.1;

    return (complexity_penalty + importance_weight + history_factor)
        .clamp(0.0, 1.0);
}
```

### ERH Profile

```uad
erh_profile DefaultProfile {
    prime_threshold: 0.6,
    fit_alpha: 0.8,
    dataset: "historical_actions.csv",
}
```

---

## 認知安全

認知安全是 UAD 關注的另一個核心議題。它涉及：

1. **認知偏誤檢測**
2. **資訊操縱識別**
3. **決策過程透明化**

```uad
struct CognitiveRisk {
    bias_type: String,
    severity: Float,
    mitigation: String,
}

fn detect_cognitive_bias(action: Action, agent: Agent) -> [CognitiveRisk] {
    let mut risks = [];

    // 確認偏誤 (Confirmation Bias)
    if action.aligns_with(agent.prior_beliefs) && action.importance < 0.3 {
        risks.push(CognitiveRisk {
            bias_type: "confirmation_bias",
            severity: 0.7,
            mitigation: "Seek disconfirming evidence",
        });
    }

    // 過度自信 (Overconfidence)
    if agent.confidence > 0.9 && agent.experience < 10 {
        risks.push(CognitiveRisk {
            bias_type: "overconfidence",
            severity: 0.8,
            mitigation: "Request peer review",
        });
    }

    return risks;
}
```

---

## 時間語義

UAD 語言具有先進的時間語義，受音樂理論啟發。

### 樂理式時間結構

```uad
score AttackSimulation {
    tempo: 120,  // BPM
    time_signature: 4/4,

    track Attacker {
        bars 1..4 {
            motif reconnaissance;
        }
        bars 5..8 {
            use motif attack at bars 5..6;
            use motif retreat at bars 7..8;
        }
    }

    track Defender {
        bars 1..8 {
            motif monitor;
        }
    }
}

motif reconnaissance {
    duration: 4 beats,
    actions: [
        Action { type: "scan", timestamp: 0 },
        Action { type: "probe", timestamp: 2 },
    ]
}
```

### 時間調度

- **Tick**：最小時間單位
- **Beat**：節拍（與 BPM 相關）
- **Bar**：小節（多個 beat 組成）
- **Event Queue**：事件調度隊列

---

## 共振與糾纏

### 弦理論語義 (M2.4)

UAD 引入「弦場」概念來建模複雜的耦合關係：

```uad
string EthicalField {
    modes {
        integrity: Float,
        transparency: Float,
        accountability: Float,
    }
}

string SystemStability {
    modes {
        resilience: Float,
        adaptability: Float,
    }
}

// 定義耦合關係
coupling EthicalField.integrity {
    mode_pair (integrity, resilience) with strength 0.7;
}

// 共振規則
resonance when freq(EthicalField.integrity) ≈ freq(SystemStability.resilience) {
    mark "High coherence state - system is stable and ethical";
}
```

### 量子糾纏語義 (M2.5)

變數可以「糾纏」，共享相同的值：

```uad
let x: Int = 10;
let y: Int = 20;
let z: Int = 30;

// 糾纏三個變數
entangle x, y, z;

// 修改任一變數，其他變數同步更新
x = 42;
// 現在 y == 42, z == 42
```

---

## 與其他語言的差異

### vs. Rust

- **UAD**：專注於對抗動態和道德評估
- **Rust**：系統程式設計，記憶體安全

### vs. Python

- **UAD**：強型別、編譯時檢查、時間語義
- **Python**：動態型別、腳本式、通用目的

### vs. Erlang/Elixir

- **UAD**：內建對抗行為建模、ERH 評分
- **Erlang**：並發、容錯、Actor 模型

### 獨特之處

1. **專為對抗場景設計**：內建 Action/Judge/Agent 概念
2. **時間即一等公民**：樂理式時間結構、事件調度
3. **道德評估原生支持**：ERH 評分、認知偏誤檢測
4. **多範式融合**：函數式 + 宣告式 + 物件導向
5. **弦場與糾纏**：創新的耦合和同步語義

---

## 結語

UAD 語言不僅僅是一個程式語言，它是一套**思考框架**，用於分析和理解在充滿對抗的複雜系統中，行為、道德和時間如何交織在一起。

透過 UAD，我們可以：

✅ **建模對抗行為**：描述攻防策略、博弈過程  
✅ **量化道德風險**：使用 ERH 評分系統  
✅ **確保認知安全**：檢測偏誤、透明化決策  
✅ **精確控制時間**：樂理式時間結構  
✅ **表達複雜耦合**：弦場與糾纏語義

**UAD = 對抗 + 道德 + 時間 + 型別安全**

---

## 參考資料

- [LANGUAGE_SPEC.md](../docs/LANGUAGE_SPEC.md)
- [MODEL_LANG_SPEC.md](../docs/MODEL_LANG_SPEC.md)
- [SEMANTICS_OVERVIEW.md](./SEMANTICS_OVERVIEW.md)
- [WHITEPAPER.md](../docs/WHITEPAPER.md)

---

_本文件描述 UAD 語言的核心範式與設計哲學。_  
_最後更新：2025_
