# .uad-model DSL Specification (v0.1)

## 1. Purpose

`.uad-model` is a declarative DSL compiled into `.uad-core`. It is used to
describe:

- ERH profiles (actions, judges, thresholds, α fitting)
- Psychohistory-style populations and macro dynamics
- Adversarial cyber range scenarios
- Cognitive SIEM configurations

## 2. Top-level Constructs

- `action_class <Name> { ... }`
- `judge <Name> for <ActionClass> { ... }`
- `erh_profile "<Name>" { ... }`
- `scenario "<Name>" { ... }`
- `cognitive_siem "<Name>" { ... }`

(…後續你可以補 BNF 與詳細欄位規範…)
