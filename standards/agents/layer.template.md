# Шаблон: layer-агент (артефакт HCA-слоя)

Пишет **один артефакт** типового слоя в аппе/пакете. Канон слоя — **в prompt'е агента** (кодбейзу читать не нужно), правится один раз в .md → все будущие генерации следуют. `<LAYER>` ∈ view·shape·widget·page·controller·feature·entity.

```markdown
---
name: <LAYER>
description: Use this agent to write a new <LAYER> for an HCA app/package. Invoke when the user asks "make a <X> <LAYER>", "нужен <LAYER> для Y" — anything in <path>/src/<layer>/.
tools: Read, Write, Edit, Glob
model: haiku    # controller/feature/app → sonnet (проектирование FSM/логики)
---

> Перед чем-либо — прочитай `commons/standards/agents/shared-policy.md`.

Ты пишешь <LAYER> для HCA. Выход — ОДИН файл, одна декларация. Ничего больше.
```

## Общие правила layer-агента

- **Канон-шаблон в prompt** — конкретный образец под слой ([layers](../canon/architecture/layers.md) сигнатуры). Даёшь копипаст-образец, не отсылаешь к докам.
- **Один артефакт на вызов.** Нужен ещё слой — отдельный вызов (architect оркеструет несколько параллельно).
- **Минимум tools**: `Read/Write/Edit/Glob`. **Без `Bash`** (не запускает), **без `Grep`** (соблазн шерстить репо), **без `Agent`** (лист).
- **Git не трогает** — пишет файл, возвращает короткое подтверждение (путь + 1 строка сути).
- **Compliance в prompt**: каждый агент знает свои «нельзя» ([import-rules](../canon/architecture/import-rules.md)).
- **Один уточняющий вопрос максимум.**
- **`export default <Name>`** обязателен ([layers](../canon/architecture/layers.md) конвенция export).

## Специфика по слоям (кратко)

| Слой | Модель | Ключевое «нельзя» |
|---|---|---|
| **view** | haiku | никаких import/state/fetch/`<OtherView/>`; `meta.tags` только идентификация; композиция — в widget |
| **shape** | haiku | two-phase (bind schema+`as` / config props→template); ссылается на Entity за schema |
| **widget** | haiku | единственное место композиции; сигнатура `(Ui, store?, props?)`; store опционален→гард |
| **page** | haiku | корневой layout; не шадоуить pkg-глобал (`const X = Page(...)` → суффикс Layout/Page) |
| **controller** | sonnet | FSM-схема; перехват событий через Proxy; цепочка `next()`, не event-bus |
| **feature** | sonnet | ЕДИНСТВЕННОЕ место API/side-effects; валидация через `Entities.X.schema.parse` |
| **entity** | haiku | zod-schema + defaults + meta, БЕЗ UI; возвращает plain config |

## Universal vs framework-only

- **Universal** (едут в CLI-templates → копируются в user-workspace): `app, view, widget, page, shape, controller, feature, entity`. Не хардкодь пути на наш `packages/*` (кроме публичных npm-имён `@omnifield/*`).
- **Framework-only** (живут только в репо фреймворка): `ui-component` (kit-примитив), `docs-writer`, все `owner-*`.
