# Omnifield Commons

**Нейтральная общая земля экосистемы Omnifield.** От неё зависят ВСЕ остальные репо, но она не зависит ни от кого. Держит две вещи с разным ритмом:

- **`contracts/`** — версионируемые мосты между репо (FE↔BE API-схемы, домен↔апп, *↔devops). Публикуются как `@omnifield/contract-*`. Меняются часто (на каждую cross-repo фичу).
- **`standards/`** — каноны и дисциплина (POLICY, HCA-канон, DoR/DoD, гигиена доков, agent-шаблоны, протокол интеграций). НЕ публикуются — наследуются через copy/submodule. Меняются редко (процесс).

Плюс `docs/` — общий глоссарий и cross-cutting ADR.

## Зачем отдельный репо

`backend/` не должен импортить из `framework/`, а FE — из `backend/`. Иначе репо связываются в комок. `commons/` = нейтральная земля, от которой зависят обе стороны шва. Изменение контракта = version-bump = **видимый брейк**, а не тихий бриж «наугад».

## Топология экосистемы (5 репо)

| Репо | Что | Тулчейн |
|---|---|---|
| **framework/** | все FE-пакеты (core·kit·builders·CLI·boosters·workspace·домены·canvas·desktop) + `testbed/` (безликие mechanics-апы) | nx/pnpm, `@omnifield/*` |
| **apps/** | продуктовые аппы с лицом (prod-условия, published `@omnifield/*`) | nx/pnpm |
| **backend/** | Python-сервисы (auth·lang·learn·community·llm·voice·image) | uv/FastAPI |
| **devops/** | docker·gateway(nginx)·WireGuard·CI·deploy·self-host bootstrap | shell/yaml |
| **commons/** ← ты здесь | contracts (мосты) + standards (каноны) + shared docs | mixed |

У каждого репо — свои архитектор+овнеры, свой скоуп/дока. **Дисциплина — из `commons/standards`. Общение — через `commons/contracts`.**

## Структура

```
commons/
├── standards/                  # НЕ публикуется — каноны/дисциплина
│   ├── POLICY.md               # верховные правила (читать ПЕРВЫМ)
│   ├── canon/                  # HCA-канон (декомпозирован по темам)
│   │   ├── principles/         #   философия «почему»
│   │   ├── architecture/       #   HCA «как устроено»
│   │   ├── packages/           #   модель пакета
│   │   ├── components/         #   UI/kit модель
│   │   └── compliance/         #   enforcement (линтер)
│   ├── workflow/               # роадмап-шаблон · DoR/DoD · docs-гигиена · git-scope
│   ├── integration-protocol.md # шов standards↔contracts, онбординг скоупов
│   └── agents/                 # шаблоны agent-конфигов (architect · owner · layer)
├── contracts/                  # публикуется @omnifield/contract-*
└── docs/                       # глоссарий · cross-cutting ADR
```

## Потребление

- **Контракты:** `pnpm add @omnifield/contract-<seam>` из реестра — репо НЕ нужно клонировать.
- **Стандарты:** git submodule / copy — каждый репо наследует канон и agent-шаблоны.

## Происхождение

v2-миграция из оракула `capsule` (ADR 077). 🔵 порт = зрелое, вычитанное из оракула. 🟢 ново = дисциплина, которой в v1 не было.
