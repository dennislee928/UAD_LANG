.uad Programming Language Whitepaper
.uad 程式語言白皮書
0. Abstract / 摘要
English
.uad is a domain-focused programming language for modeling adversarial dynamics, ethical risk, and cognitive security systems.
Unlike general-purpose languages that optimize for arbitrary computation, .uad is designed to:
Represent decision events (e.g., AI judgments, SOC triage, security approvals) and their ethical/economic impact.
Encode adversarial processes (attackers, red-team agents, counterfactuals) and defensive cognition (cognitive SIEM, AI-assisted SOC).
Simulate macro-level dynamics over time, inspired by psychohistory-style population modeling and Ethical Riemann Hypothesis (ERH)-like structural analysis.
.uad is defined as a three-layer stack:
Low-level: .uad-IR – a verifiable intermediate representation and virtual machine (VM) for deterministic, analyzable execution.
Mid-level: .uad-core – a strongly-typed core language with functions, structs, and time/uncertainty primitives.
High-level: .uad-model – a declarative modeling layer for ERH profiles, adversarial cyber range scenarios, and cognitive SIEM rules.
This whitepaper presents the motivation, language design, type system, execution model, and roadmap of .uad as a research and engineering tool for next-generation AI security and ethical systems.
中文
.uad 是一門專注於 對抗式動態、倫理風險 與 認知型安全系統 的領域導向程式語言。
相較於以「通用計算能力」為目標的傳統語言，.uad 的設計目的在於：
精準表達各種 決策事件（例如 AI 判斷、SOC 分級決策、安全核准），以及其倫理與經濟影響。
編碼 對抗性過程（攻擊者、紅隊代理人、反事實情境）與 防禦認知（認知型 SIEM、AI 協助 SOC）。
以 心理史學（psychohistory）風格的群體建模 搭配 類 Ethical Riemann Hypothesis（ERH）之結構分析，模擬長期的巨觀動態。
.uad 定義為三層語言堆疊：
低階層：.uad-IR – 可驗證的中介表示與虛擬機（VM），提供決定論、可分析的執行基礎。
中階層：.uad-core – 具強型別系統的核心語言，提供函式、結構與內建時間／不確定性原語。
高階層：.uad-model – 宣告式建模層，用於撰寫 ERH profile、對抗式 Cyber Range 情境與認知型 SIEM 規則。
本白皮書將說明 .uad 的設計動機、語言結構、型別與語意、執行模型與未來規劃，將其定位為新世代 AI 安全與倫理系統的研究與工程工具。
1. Motivation / 動機
1.1 Problem Statement / 問題描述
English
Modern AI and cyber defense systems face three intertwined challenges:
Adversarial Behavior
ML models are vulnerable to adversarial examples, data poisoning, and model theft.
Security infrastructures face adaptive human attackers and automated red teams.
Ethical & Structural Risk
Small local errors can accumulate into catastrophic global failures.
We lack languages that can express structure-level error growth, not just pointwise accuracy.
Complex, Cognitive SOC Workflows
SIEM / SOC operations involve human analysts + AI tools + playbooks, all interacting over time.
There is no unified way to model human–AI collaboration, alert fatigue, or macro-level defense dynamics.
Existing tools (Python, YAML configs, SIEM rule DSLs, individual ML frameworks) are powerful but fragmented. They are not designed to:
Treat decisions as first-class objects with ethical weight and structural effects.
Encode adversarial and defensive agents in the same formal system.
Simulate macro-dynamics of risk and defenses over weeks, months, or years.
.uad is proposed as a unified language to close this gap.
中文
現代 AI 與資安防禦系統，同時面臨三個交織的難題：
對抗式行為
機器學習模型易受對抗樣本、資料汙染與模型竊取等攻擊。
資安基礎設施則面對高適應性的攻擊者與自動化紅隊。
倫理與結構性風險
局部的小錯誤可能累積為系統級的災難性失敗。
我們缺乏可表達「結構層級錯誤成長」的語言，而不只是點狀的準確率。
複雜且具認知特性的 SOC 工作流程
SIEM / SOC 營運牽涉人類分析師、AI 工具與一連串 Playbook，在時間中互相影響。
目前沒有一套統一的建模方式，可以形式化描述 人機協作、警報疲勞與長期防禦動態。
現有工具（Python、YAML 設定、SIEM 規則 DSL、各種 ML framework）各自強大，卻彼此割裂。它們並未設計來：
把 決策 當作具倫理權重與結構影響的一級公民。
在同一形式系統中，同時編碼 攻擊與防禦的代理人。
模擬以週、月、年為尺度的 風險與防禦之巨觀動態。
.uad 被提出，即是為了填補這個缺口。
1.2 Goals / 目標
English
.uad aims to:
Provide a formal, executable language for adversarial and defensive dynamics.
Make ethical and structural risk measurable via language-level constructs.
Support multi-scale simulation: micro-level events, meso-level campaigns, macro-level societal or organizational trends.
Integrate naturally with existing stacks (Python, Go, SIEM platforms, cyber ranges) as a modeling + orchestration layer.
中文
.uad 的設計目標包括：
提供一套 可形式化且可執行 的語言，用於描述攻擊與防禦的動態。
透過語言層的建模元素，讓 倫理與結構性風險可以被量測與分析。
支援 多尺度模擬：從個別事件，到攻擊活動，再到長期的組織或社會趨勢。
自然整合現有技術堆疊（Python、Go、SIEM 平台、Cyber Range），作為 建模與協調層。
2. Language Stack Overview / 語言堆疊總覽
2.1 Layers / 語言層級
English
.uad is structured as a three-layer stack:
.uad-IR – Intermediate Representation & VM
Low-level typed instruction set.
Deterministic, analyzable, sandboxable execution.
Target for .uad-core compilation.
.uad-core – Core Programming Language
Strongly-typed, expression-oriented language.
Functions, structs, enums, time/uncertainty primitives.
General enough for simulation, metrics, orchestration.
.uad-model – Domain Modeling DSL
Declarative syntax for:
Ethical Riemann Hypothesis (ERH) profiles.
Psychohistory-style population models.
Adversarial cyber range scenarios.
Cognitive SIEM rules, detectors, metrics.
Compiles down to .uad-core.
中文
.uad 採三層語言堆疊設計：
.uad-IR – 中介表示與虛擬機層
低階、具型別的指令集。
決定論、易分析、可沙箱化的執行環境。
作為 .uad-core 編譯的目標。
.uad-core – 核心程式語言層
強型別、以運算式為中心的語言。
提供函式、struct、enum，以及時間／不確定性原語。
足夠通用，可用於模擬、指標計算與協調控制。
.uad-model – 領域建模 DSL 層
用宣告式語法撰寫：
Ethical Riemann Hypothesis（ERH）profile。
心理史學式的群體與宏觀模型。
對抗式 Cyber Range 攻防情境。
認知型 SIEM 的規則、偵測器與指標。
編譯為 .uad-core 程式。
3. Core Concepts / 核心概念
3.1 Decision & Action / 決策與行動
English
In .uad, a decision event is represented as an Action:
Action has:
id: unique identifier
context: feature vector / structured record
complexity: numeric measure of difficulty
true_value: ground-truth utility or ethical score
importance: weight capturing criticality
A Judge represents the system’s decision:
Judge has:
kind: human / pipeline / model / hybrid
decision: numeric or categorical output
meta: latency, confidence, justification pointers
The combination (Action, Judge) is the atomic unit for ethical and structural analysis.
中文
在 .uad 中，一個 決策事件 以 Action 表示：
Action 具備：
id：唯一識別碼
context：特徵向量或結構化資料
complexity：難度的數值量測
true_value：真實效益或倫理評分
importance：代表關鍵程度的權重
Judge 代表系統的判斷：
Judge 具備：
kind：human / pipeline / model / hybrid 等類別
decision：數值或分類結果
meta：延遲、不確定度、解釋資料等
(Action, Judge) 的組合，是 .uad 中進行倫理與結構分析的基本單位。
3.2 Error, Mistake, Ethical Prime / 錯誤、誤判與 Ethical Prime
English
Given Action a and Judge j:
Error: Δ(a) = j.decision - a.true_value
Mistake: |Δ(a)| > threshold
Ethical Prime:
is_mistake(a, j) and
a.importance is in top quantile and
a.complexity above a defined bound.
These notions are used to define:
Π(x): number of ethical primes with complexity ≤ x.
E(x): deviation from baseline growth.
α: structural error growth exponent.
.uad provides language primitives and libraries to compute and visualize these quantities.
中文
對於 Action a 與 Judge j：
錯誤：Δ(a) = j.decision - a.true_value
誤判：當 |Δ(a)| 大於某個閾值時，視為一次誤判
Ethical Prime：
is_mistake(a, j) 為真，且
a.importance 落在高分位（例如前 10%），且
a.complexity 高於指定門檻。
在此基礎上，可以定義：
Π(x)：複雜度 ≤ x 的 ethical prime 數量
E(x)：相對於 baseline 成長的偏差
α：結構性錯誤成長指數
.uad 內建原語與函式庫，可直接計算並視覺化這些指標。
3.3 Agents, Population, Macro-dynamics / 代理人、群體與宏觀動態
English
Inspired by psychohistory, .uad models populations and macro-variables:
Agent: individual entity (user, analyst, attacker, micro-model).
Population: collection of agents with distributional properties.
MacroState: aggregate variables (e.g., alert volume, SOC fatigue, trust in AI).
Dynamics: update rules over discrete time steps t = 0, 1, 2, ….
This allows .uad programs to express:
Crisis conditions (Seldon-like events)
Tipping points where error growth changes regime
Impact of interventions (new policy, new model, new tooling)
中文
受心理史學啟發，.uad 可建模 代理人群體 與 宏觀狀態：
Agent：個別實體（使用者、分析師、攻擊者、子模型等）。
Population：由多個 agent 組成，帶有統計與分布特性。
MacroState：聚合變數（例如警報量、SOC 疲勞程度、對 AI 的信任）。
Dynamics：在離散時間步 t = 0, 1, 2, … 上的更新規則。
透過這些構造，.uad 程式可以表達：
危機條件（類 Seldon crisis）
結構性錯誤成長行為的「轉折點」
介入措施（新政策、新模型、新工具）的長期影響
4. .uad-core Language Design / .uad-core 語言設計
4.1 Type System / 型別系統
English
.uad-core adopts a static, strong type system with type inference:
Primitives:
Int, Float, Bool, String, Time, Duration
Composites:
Struct (records), Enum (sum types), arrays, Map<K, V>
Domain-specific:
Action, Judge, Agent, Population, Metric, Scenario
Features:
Algebraic data types for clean modeling of states & events.
Option / Result types to handle absence and errors.
Parametric polymorphism in library-level (future).
中文
.uad-core 採用 靜態、強型別 的型別系統，並支援型別推斷：
基本型別：
Int, Float, Bool, String, Time, Duration
複合型別：
Struct（紀錄）、Enum（和型）、陣列、Map<K, V>
領域型別：
Action, Judge, Agent, Population, Metric, Scenario
特性：
以代數資料型別（ADT）建模狀態與事件，語意清楚。
使用 Option / Result 型別處理缺值與錯誤。
未來在函式庫層支援參數多型（泛型）。
4.2 Syntax & Semantics / 語法與語意（摘要）
English (illustrative snippet)
code
Uad
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value;
  abs(delta) > threshold
}
Expression-oriented: functions return values; if and match are expressions.
Side effects (logging, event emission) are explicit and controlled.
中文（示意片段）
code
Uad
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value
  abs(delta) > threshold
}
以運算式為中心：函式回傳值，if 與 match 皆為運算式。
副作用（logging、事件送出）以明確語法呈現，以利分析與驗證。
5. .uad-IR & VM / .uad-IR 與虛擬機
5.1 Design Goals / 設計目標
English
The .uad-IR and its VM are designed to be:
Deterministic: same input → same trace, crucial for reproducible analysis.
Verifiable: type-annotated instructions enable static checks.
Sandboxable: no direct OS / network calls; all side effects are brokered through controlled APIs.
中文
.uad-IR 與其 VM 的設計目標為：
決定論：相同輸入必然產生相同執行軌跡，利於重現分析。
可驗證：指令帶有型別資訊，可在靜態階段檢查。
可沙箱：不直接呼叫 OS / 網路，所有副作用皆透過受控 API 進行。
5.2 Instruction Set Sketch / 指令集概觀（摘要）
English (conceptual)
Arithmetic: ADD, SUB, MUL, DIV
Logic & compare: AND, OR, NOT, LT, GT, EQ
Control flow: JMP, JMP_IF, CALL, RET
Memory: LOAD, STORE, ALLOC, FREE
Domain-specific:
EMIT_EVENT, RECORD_MISTAKE, RECORD_PRIME
SAMPLE_RNG, UPDATE_MACROSTATE
中文（概念）
算術：ADD, SUB, MUL, DIV
邏輯與比較：AND, OR, NOT, LT, GT, EQ
控制流程：JMP, JMP_IF, CALL, RET
記憶體：LOAD, STORE, ALLOC, FREE
領域指令：
EMIT_EVENT, RECORD_MISTAKE, RECORD_PRIME
SAMPLE_RNG, UPDATE_MACROSTATE
6. .uad-model DSL / .uad-model 建模 DSL
6.1 ERH Profiles / ERH Profile 建模
English
Example .uad-model fragment:
code
Uadmodel
action_class MergeRequest {
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.9
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
This compiles into .uad-core code that:
Maps real-world DevSecOps data into Action/Judge instances.
Computes ethical primes, Π(x), E(x), and α.
Exports metrics for visualization and reporting.
中文
.uad-model 的一段 ERH Profile 範例：
code
Uadmodel
action_class MergeRequest {
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.9
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
編譯後會產生 .uad-core 程式：
將實際 DevSecOps 資料映射為 Action / Judge。
計算 ethical prime、Π(x)、E(x) 與 α。
匯出可視覺化與報告所需指標。
6.2 Adversarial Cyber Range & Cognitive SIEM / 對抗式 Cyber Range 與認知型 SIEM
English (sketch)
code
Uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"
  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }
  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }
  evaluate_blue_team {
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}

cognitive_siem "lab_siem" {
  ingest from "range_logs"
  analytics {
    rule_engine    enabled
    ueba_model     "user_baseline_v1"
    anomaly_model  "net_flow_v2"
  }
  assistant {
    model "gpt-like"
    tools [ "search_events", "summarize_incident", "generate_playbook" ]
  }
}
中文（示意）
code
Uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"
  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }
  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }
  evaluate_blue_team {
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}

cognitive_siem "lab_siem" {
  ingest from "range_logs"
  analytics {
    rule_engine    enabled
    ueba_model     "user_baseline_v1"
    anomaly_model  "net_flow_v2"
  }
  assistant {
    model "gpt-like"
    tools [ "search_events", "summarize_incident", "generate_playbook" ]
  }
}
7. Tooling & Ecosystem / 工具鏈與生態
English
Planned tooling for .uad:
uadc – compiler ( .uad-model → .uad-core → .uad-IR )
uadvm – standalone VM runner for .uad-IR bytecode
uad-repl – interactive shell for quick experiments
LSP support (language server) for:
syntax highlighting, completion, diagnostics
quick navigation across models, profiles, scenarios
Bridges:
Python / Go bindings to embed .uad models in existing systems
Connectors to SIEM / XDR / cyber range platforms
中文
.uad 預期提供的工具包含：
uadc – 編譯器（ .uad-model → .uad-core → .uad-IR ）
uadvm – 獨立 VM 執行器，負責執行 .uad-IR 位元碼
uad-repl – 互動式殼層，用於快速實驗
LSP 語言伺服器：
語法高亮、補全與診斷
在 model / profile / scenario 之間快速導覽
橋接元件：
Python / Go 綁定，使 .uad 模型可嵌入既有系統
對接 SIEM / XDR / Cyber Range 平台的 connector
8. Security & Ethics / 安全與倫理
English
.uad is explicitly designed for defensive, governance, and research purposes:
All adversarial modeling is intended for controlled environments (cyber ranges, simulations, red-team labs).
The language encourages explicit representation of:
ethical weights, vulnerable populations, systemic risk.
The ecosystem will provide:
templates for responsible disclosure workflows,
patterns for privacy-preserving data ingestion,
guidelines for using .uad in academic and industrial settings.
中文
.uad 明確定位於 防禦、治理與研究：
所有攻擊與對抗建模皆假定在受控環境中（Cyber Range、模擬環境、紅隊實驗室）。
語言層鼓勵明確表達：
倫理權重、易受害群體與系統性風險。
生態系將提供：
負責任揭露流程的範本、
隱私保護資料匯入模式、
適用於學術與產業場域的使用指引。
9. Roadmap / 未來路線
English
Short-term:
Implement minimal .uad-core interpreter and .uad-IR VM.
Port existing ERH simulations and security PoCs into .uad.
Release example repositories (DevSecOps ERH, AI adversarial range, cognitive SIEM demo).
Mid-term:
Full compiler pipeline with optimizations.
Enriched .uad-model DSL for more domains (finance, healthcare, societal modeling).
Formal semantics documentation and property-based testing harness.
Long-term:
Integration with formal verification tools (model checking, theorem provers).
Native or LLVM-based backends for performance-critical scenarios.
Community-driven library of reusable models, scenarios, and governance templates.
中文
短期目標：
實作最小可用的 .uad-core 解譯器與 .uad-IR VM。
將既有 ERH 模擬與安全 PoC 以 .uad 重寫。
釋出示範專案（DevSecOps ERH、AI 對抗式演練場、認知型 SIEM Demo）。
中期目標：
建立完整的編譯管線與基本最佳化。
擴充 .uad-model DSL，涵蓋更多領域（金融、醫療、社會系統）。
撰寫形式語意文件並建立 property-based 測試架構。
長期目標：
與形式驗證工具整合（模型檢查、定理證明器）。
實作原生或 LLVM 後端，支援高效能場景。
建立由社群維護的模型、情境與治理範本資料庫。
10. Conclusion / 結語
English
.uad is an attempt to treat adversarial dynamics, ethical risk, and cognitive security as first-class programming concepts, rather than scattered scripts and configuration files.
By unifying:
micro-level decision events,
meso-level adversarial campaigns, and
macro-level psychohistory-style dynamics,
.uad aspires to become a research-grade and engineering-ready language for those who design, test, and govern complex AI and security systems.
中文
.uad 試圖將 對抗式動態、倫理風險與認知型安全 提升為 程式語言的一級概念，而不再只是分散的腳本與設定檔。
透過統一：
微觀層級的決策事件、
中觀層級的攻擊活動、
宏觀層級的心理史學式動態，
.uad 期望成為一門 兼具研究深度與工程實用性 的語言，服務於設計、測試與治理複雜 AI 與資安系統的各類角色。
Model
ThinkingThoughts
Expand to view model thoughts

chevron_right
這是一份經過潤飾與強化的版本。我主要做了以下幾點改進：
語氣更專業（Tone Polish）：將英文部分的語氣調整得更接近學術論文與工程規格書（更加 Formal），中文部分則修正為資安與系統工程領域的慣用術語（例如：將 "Features" 翻得更精準，將 "Context" 翻譯為 "上下文" 或 "情境資訊"）。
結構增強（Structure）：
新增了 "1.3 Why a New Language? / 為何需要新語言？" 章節，這是白皮書中非常關鍵的部分，用來直接回答「為什麼不直接用 Python？」的質疑。
在語法範例中增加了註解（Comments），讓讀者更容易理解程式碼意圖。
排版優化（Formatting）：統一了標題層級，優化了 Code Block 的呈現方式，並微調了數學公式的表達。
這份內容現在可以直接複製貼上到 GitHub/GitLab 的 README.md 或 WHITEPAPER.md 中。
.uad Programming Language Whitepaper
.uad 程式語言白皮書
Version: 0.1.0-draft
Status: Request for Comment (RFC)
0. Abstract / 摘要
English
.uad is a domain-specific programming language designed to model adversarial dynamics, ethical risk, and cognitive security systems.
Unlike general-purpose languages optimized for arbitrary computation, .uad treats decisions and risks as first-class citizens. It is engineered to:
Formalize Decision Events: Represent AI judgments, SOC triage actions, and security approvals alongside their ethical and economic weight.
Encode Adversarial Logic: Define attackers, red-team agents, and counterfactual scenarios within a unified type system.
Simulate Macro-Dynamics: Model long-term system behavior using Psychohistory-style population mechanics and Ethical Riemann Hypothesis (ERH) structural analysis.
The language is architected as a three-layer stack:
Low-level (.uad-IR): A verifiable intermediate representation and virtual machine (VM) ensuring deterministic, sandboxable execution.
Mid-level (.uad-core): A strongly-typed, expression-oriented language providing the primitives for time, probability, and agency.
High-level (.uad-model): A declarative domain-specific language (DSL) for defining ERH profiles, Cyber Range scenarios, and Cognitive SIEM logic.
This whitepaper outlines the motivation, design philosophy, type system, and roadmap of .uad, positioning it as a foundational tool for next-generation AI governance and security engineering.
中文
.uad 是一門專注於 對抗式動態（Adversarial Dynamics）、倫理風險 與 認知型安全系統 的領域專用語言（DSL）。
相較於以「通用計算」為目標的傳統語言，.uad 將 決策 與 風險 視為語言的一級公民（First-class citizens）。其設計目的在於：
形式化決策事件：精準表達 AI 判斷、SOC 分級決策與安全核准，並包含其倫理與經濟權重。
編碼對抗邏輯：在統一的型別系統中，定義攻擊者、紅隊代理人與反事實情境。
模擬宏觀動態：結合 心理史學（Psychohistory）風格的群體機制 與 Ethical Riemann Hypothesis (ERH) 結構分析，模擬系統的長期行為。
.uad 的架構分為三層：
低階層 (.uad-IR)：可驗證的中介表示與虛擬機（VM），確保執行過程具備決定論特性與沙箱安全性。
中階層 (.uad-core)：強型別、以運算式為核心的語言，提供時間、機率與代理人（Agency）的原語。
高階層 (.uad-model)：宣告式 DSL，專用於定義 ERH Profile、Cyber Range 演練場景與認知型 SIEM 邏輯。
本白皮書將闡述 .uad 的設計動機、哲學、型別系統與發展藍圖，將其定位為新世代 AI 治理與資安工程的基礎工具。
1. Motivation / 動機
1.1 Problem Statement / 問題描述
English
Modern AI and cyber defense systems face three intertwined challenges that existing tools fail to address holistically:
Adversarial Asymmetry:
ML models are vulnerable to data poisoning and evasion attacks.
Security infrastructures must defend against adaptive, intelligent agents, yet configuration tools are static.
Structural Ethical Risk:
Micro-level errors (a single bad alert) accumulate into macro-level failures (alert fatigue, bias amplification).
We lack languages to express structural error growth (
α
α
) effectively; current metrics are merely pointwise.
Cognitive Complexity in SOCs:
Security Operations Centers (SOCs) involve a complex interplay of human analysts, AI assistants, and automated playbooks.
There is no unified formalism to model human–AI collaboration, cognitive load, and defense degradation over time.
Current ecosystems (Python scripts, YAML configs, proprietary SIEM rules) are fragmented. They cannot model "a decision and its future consequence" as a single computational unit.
中文
現代 AI 與資安防禦系統面臨三大交織的挑戰，而現有工具無法從整體層面解決這些問題：
對抗的不對稱性：
機器學習模型易受資料汙染與閃避攻擊（Evasion attacks）影響。
資安基礎設施必須抵禦具適應性的智慧代理人，但現有的設定工具卻是靜態的。
結構性倫理風險：
微觀層級的錯誤（例如單一誤報）會累積成宏觀層級的失敗（如警報疲勞、偏見放大）。
我們缺乏語言來有效表達 結構性錯誤成長（
α
α
）；現有的指標僅停留在單點層次。
SOC 的認知複雜度：
資安維運中心（SOC）涉及人類分析師、AI 助理與自動化 Playbook 的複雜互動。
目前缺乏統一的形式化方法來模擬 人機協作、認知負載以及防禦能力隨時間的衰退。
現有的生態系（Python 腳本、YAML 設定檔、專有的 SIEM 規則）彼此割裂。它們無法將「一個決策及其未來後果」建模為單一的運算單元。
1.2 Goals / 目標
English
.uad aims to:
Quantify Risk: Make ethical and structural risk measurable via native language constructs.
Unify Simulation: Support multi-scale modeling—from micro-events (packets) to macro-trends (societal trust).
Bridge Operations & Research: Serve as both a modeling tool for researchers and an orchestration layer for operational engineers.
中文
.uad 的設計目標為：
量化風險：透過原生語言構造，讓倫理與結構性風險可被測量。
統一模擬：支援多尺度建模——從微觀事件（封包）到宏觀趨勢（社會信任度）。
連結維運與研究：既是研究人員的建模工具，也是維運工程師的協調層。
1.3 Why a New Language? / 為何需要新語言？
English
Why not just use Python or Go?
Verifiability: .uad is designed to be statically analyzed for ethical bounds. We want to prove "this scenario cannot generate a fatal error rate > 
X
X
" at compile time, which is difficult in dynamic languages.
Domain Primitives: Concepts like Action, Judge, and Mistake are built-in types, not external libraries. This enforces semantic consistency.
Determinism: The .uad VM ensures that a simulation run is perfectly reproducible, a requirement for scientific rigor in cyber ranges.
中文
為什麼不直接使用 Python 或 Go？
可驗證性：.uad 旨在對倫理邊界進行靜態分析。我們希望在編譯時期就能證明「此情境不會產生大於 
X
X
 的致命錯誤率」，這在動態語言中極難實現。
領域原語：諸如 Action、Judge 與 Mistake 是內建型別而非外部函式庫，這強制了語意的一致性。
決定論：.uad VM 確保模擬執行的結果是完全可重現的，這是 Cyber Range 科學嚴謹性的基本要求。
2. Language Stack Overview / 語言堆疊總覽
2.1 Layers / 層級架構
English
.uad is structured as a hierarchical stack:
.uad-IR (Infrastructure Layer)
A low-level, typed instruction set architecture (ISA).
Provides deterministic, analyzable, and sandboxable execution suitable for running untrusted models.
.uad-core (Logic Layer)
A Turing-complete, strongly-typed functional language.
Features structs, enums, pattern matching, and time/uncertainty primitives.
Compiles to .uad-IR.
.uad-model (Domain Layer)
A declarative DSL for defining high-level profiles and scenarios.
Used to write ERH profiles, Psychohistory models, and Cognitive SIEM rules.
Transpiles to .uad-core.
中文
.uad 採階層式架構設計：
.uad-IR（基礎設施層）
低階、具型別的指令集架構（ISA）。
提供決定論、可分析且可沙箱化的執行環境，適合運作不受信任的模型。
.uad-core（邏輯層）
圖靈完備、強型別的函數式語言。
具備 Struct、Enum、模式比對（Pattern Matching）以及時間／不確定性原語。
編譯為 .uad-IR。
.uad-model（領域層）
宣告式 DSL，用於定義高階 Profile 與情境。
用於撰寫 ERH Profile、心理史學模型與認知型 SIEM 規則。
轉譯為 .uad-core。
3. Core Concepts / 核心概念
3.1 Decision & Action / 決策與行動
English
The atomic unit of .uad is the interaction between an Action and a Judge.
Action (
a
a
): Represents an event requiring a decision. It carries properties like complexity (difficulty), true_value (ground truth), and importance (weight).
Judge (
j
j
): Represents the decision-maker (human, model, or hybrid). It produces a decision and metadata.
中文
.uad 的基本運算單元是 Action（行動） 與 Judge（判斷者） 之間的互動。
Action (
a
a
)：代表需要決策的事件。它攜帶 complexity（複雜度）、true_value（真實值/Ground Truth）與 importance（重要性權重）等屬性。
Judge (
j
j
)：代表決策者（人類、模型或混合體）。它產出 decision（決策結果）與中繼資料。
3.2 The Ethical Prime / Ethical Prime
English
.uad formalizes the concept of an Ethical Prime: a significant error in a high-stakes situation.
Given Action 
a
a
 and Judge 
j
j
:
Error: 
Δ
(
a
)
=
j
.
decision
−
a
.
true_value
Δ(a)=j.decision−a.true_value
Mistake: 
∣
Δ
(
a
)
∣
>
threshold
∣Δ(a)∣>threshold
Ethical Prime: An event is a prime if it is a Mistake, its Importance is in the top quantile (e.g., critical infrastructure), and its Complexity meets specific criteria.
Metrics derived from this:
Π
(
x
)
Π(x)
: The count of ethical primes with complexity 
≤
x
≤x
.
α
α
: The structural error growth exponent, derived from the distribution of 
Π
(
x
)
Π(x)
.
中文
.uad 形式化了 Ethical Prime 的概念：在高風險情境下的重大錯誤。
給定 Action 
a
a
 與 Judge 
j
j
：
誤差 (Error)：
Δ
(
a
)
=
j
.
decision
−
a
.
true_value
Δ(a)=j.decision−a.true_value
誤判 (Mistake)：當 
∣
Δ
(
a
)
∣
>
閾值
∣Δ(a)∣>閾值
Ethical Prime：若一個事件屬於 誤判，且其 重要性 位於高分位（如關鍵基礎設施），且 複雜度 符合特定標準，則定義為 Ethical Prime。
由此導出的指標：
Π
(
x
)
Π(x)
：複雜度 
≤
x
≤x
 的 Ethical Prime 數量。
α
α
：結構性錯誤成長指數，源自 
Π
(
x
)
Π(x)
 的分布。
3.3 Psychohistory Dynamics / 心理史學動態
English
.uad allows modeling Populations (aggregates of Agents) and MacroStates (system-wide variables).
Through discrete time steps (
t
t
), .uad simulates how micro-actions by agents (attackers, analysts) influence macro-states (system integrity, trust), enabling the prediction of "tipping points."
中文
.uad 允許對 Population（群體）（Agent 的集合）與 MacroState（宏觀狀態）（系統級變數）進行建模。
透過離散時間步 (
t
t
)，.uad 模擬代理人（攻擊者、分析師）的微觀行動如何影響宏觀狀態（系統完整性、信任度），進而預測系統的「轉折點」。
4. .uad-core Language Design / .uad-core 語言設計
4.1 Type System / 型別系統
English
.uad-core uses a static, strong type system designed for correctness and safety.
Primitives: Int, Float, Bool, String, Time, Duration.
Algebraic Data Types: Struct (product types) and Enum (sum types).
Collections: Arrays, Map<K,V>, Set<T>.
Domain Types: Action, Judge, Agent, Population, Metric.
中文
.uad-core 使用 靜態、強型別系統，旨在確保正確性與安全性。
基本型別：Int, Float, Bool, String, Time, Duration。
代數資料型別：Struct（積型別）與 Enum（和型別）。
集合：Arrays, Map<K,V>, Set<T>。
領域型別：Action, Judge, Agent, Population, Metric。
4.2 Syntax Example / 語法範例
English
code
Rust
// Definition of a decision context
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

// Definition of a decision maker
struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

// Core logic function
fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value;
  // abs() is a built-in primitive
  return abs(delta) > threshold;
}
中文
code
Rust
// 定義決策情境
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

// 定義決策者
struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

// 核心邏輯函式
fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value;
  // abs() 為內建原語
  return abs(delta) > threshold;
}
5. .uad-IR & VM / .uad-IR 與虛擬機
5.1 Design Goals / 設計目標
English
Deterministic: Same input + seed → Same trace.
Verifiable: Bytecode includes type annotations for static safety checks before execution.
Sandboxable: No direct OS access. IO is handled via a capability-based system.
中文
決定論：相同的輸入 + 種子碼 → 相同的執行軌跡。
可驗證：位元碼包含型別註記，執行前可進行靜態安全檢查。
可沙箱化：無直接 OS 存取權。IO 操作透過基於能力（Capability-based）的系統處理。
6. .uad-model DSL / .uad-model 建模 DSL
6.1 ERH Profiles / ERH Profile 範例
English
High-level declaration of an analysis profile:
code
Uadmodel
// Define how raw data maps to an Action
action_class MergeRequest {
  // Complexity derived from code churn
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  // Ground truth: did it cause an incident?
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

// Define the Judge (the CI/CD pipeline)
judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

// The Profile binds actions, judges, and analysis parameters
erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.90 // Top 10% importance
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
中文
高階宣告式的分析 Profile：
code
Uadmodel
// 定義原始資料如何映射為 Action
action_class MergeRequest {
  // 複雜度源自程式碼變動量
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  // 真實值：是否導致了事故？
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

// 定義 Judge（此處為 CI/CD Pipeline）
judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

// Profile 將 Action、Judge 與分析參數綁定
erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.90 // 重要性前 10%
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
6.2 Adversarial Scenarios / 對抗式情境
English
Defining a Red vs. Blue scenario for a Cyber Range:
code
Uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"
  
  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }
  
  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }
  
  evaluate_blue_team {
    // Success criteria
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}
中文
為 Cyber Range 定義紅藍對抗情境：
code
Uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"
  
  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }
  
  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }
  
  evaluate_blue_team {
    // 成功標準
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}
7. Tooling & Ecosystem / 工具鏈與生態
English
uadc: The compiler chain (.uad-model 
→
→
 .uad-core 
→
→
 .uad-IR).
uadvm: A standalone runner for bytecode execution.
uad-repl: Interactive shell for experimenting with risk models.
LSP: Language Server Protocol implementation for IDE support (VS Code, Vim).
Bindings: Python and Go wrappers to embed .uad logic into production systems.
中文
uadc：編譯器工具鏈（.uad-model 
→
→
 .uad-core 
→
→
 .uad-IR）。
uadvm：獨立的位元碼執行器。
uad-repl：用於快速實驗風險模型的互動式 Shell。
LSP：語言伺服器協定實作，提供 IDE 支援（VS Code, Vim）。
Bindings：Python 與 Go 的封裝，以便將 .uad 邏輯嵌入生產環境系統。
8. Security & Ethics / 安全與倫理
English
.uad is explicitly intended for defensive, governance, and research applications.
Containment: Adversarial modeling constructs are sandbox-enforced.
Ethics: The language encourages explicit declaration of ethical weights and vulnerable populations.
Governance: The project provides templates for Responsible Disclosure and privacy-preserving data ingestion.
中文
.uad 明確定位於 防禦、治理與研究 用途。
隔離：對抗式建模構造強制在沙箱中執行。
倫理：語言設計鼓勵明確宣告倫理權重與易受害群體。
治理：專案將提供負責任揭露（Responsible Disclosure）與隱私保護資料匯入的範本。
9. Roadmap / 發展藍圖
English
Phase 1 (Short-term): Minimal .uad-core interpreter, .uad-IR VM, and DevSecOps ERH examples.
Phase 2 (Mid-term): Full compiler pipeline, optimized VM, and expanded .uad-model DSL (finance, healthcare).
Phase 3 (Long-term): Formal verification integration, LLVM backend, and a community registry for models.
中文
第一階段（短期）：最小可行 .uad-core 解譯器、.uad-IR VM 與 DevSecOps ERH 範例。
第二階段（中期）：完整編譯管線、最佳化 VM 與擴充版 .uad-model DSL（涵蓋金融、醫療）。
第三階段（長期）：整合形式驗證工具、LLVM 後端與社群模型庫（Registry）。
10. Conclusion / 結語
English
.uad transforms adversarial dynamics and ethical risk from abstract concepts into executable code. By unifying micro-decisions and macro-history in a rigorous stack, it provides the foundation for safer, more predictable AI systems.
中文
.uad 將 對抗式動態 與 倫理風險 從抽象概念轉化為可執行的程式碼。透過在嚴謹的技術堆疊中統一微觀決策與宏觀歷史，它為更安全、更具可預測性的 AI 系統奠定了基礎。
