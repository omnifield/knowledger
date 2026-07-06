# Golden Rules — правила + severity

Правила канона, которые **enforced линтером** ([[linter]]). Формулировки слоёв — [[../architecture/import-rules]]; здесь — severity-модель.

## Правила

1. **No Upward Imports** — нижний слой не импортит верхний.
2. **No Horizontal Imports** — `View.A` ⊥ `View.B`, `Controller.A` ⊥ `Controller.B`; связь только через Widget/`next()`.
3. **Stateless View / Shape** — без состояния, без импортов кроме Solid и типов.
4. **Composition Only in Widgets** — склейка сущностей только на уровне Widget.

Плюс структурные границы: no app-package runtime-import, no disallowed external dep в слое, no cross-zone impl-import, portable-tier ⊥ ecosystem ([[../packages/dependency-tiers]]).

## Severity — два класса

| Класс | Правила | Эффект |
|---|---|---|
| **`error` (structural)** | `app-package-import` (runtime `@omnifield/*` в `apps/*/src`), `disallowed-import` (запрещённый внешний пакет в слое), portable-tier ecosystem-import | **Валят CI** (`compliance:check`). В dev — `[STRUCTURAL ERROR — CI will fail]` без блокировки HMR. |
| **`warn` (cosmetic)** | `native-jsx`, `native-js`, `raw-class`, `upward-import`, `horizontal-import`, `side-effect-fetch`, `unknown-alias`, `cross-zone-import` | Лог с `file:line:column` + hint, CI не валят. |

Override per-call: `check(path, code, { severity: { 'app-package-import': 'warn'|'off' } })` или через опции Vite-плагина.

## Transitional allowlist

Опциональный файл (`compliance-allowlist.json`, формат `{ path, kinds[], reason, ttl }`) — суппрессит structural-нарушения по path+kind на cleanup-окно. Инструмент миграции, не постоянный обход: каждая запись с `reason` + `ttl`. Эталонная зона = пустой allowlist.

## Canon-first

Линтер-правило заводится **ДО** app-кода, который оно защищает ([[../principles/root-cause-not-symptom]] — enforcement первым). Иначе канон держится на дисциплине памяти, а не на машине.
