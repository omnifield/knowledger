# Canon digest — owner-frontend

Выжимка для owner-агента UI/HCA-зоны (framework-пакеты, apps). Полный канон — `../../canon/`,
карта — `../canon-map.md`. `⬛` — в контексте всегда, `📖` — читать первым действием.

---

## ⬛ Non-negotiables (всегда в силе)

**Правила импорта / слои** (`canon/architecture/import-rules.md`, `layers.md`)
- **No upward** (нижний слой не импортит верхний), **no horizontal** (`View.A` ⊥ `View.B`).
- Композиция сущностей — **только в Widget** (children/slots). Поведение вверх — `next()`. UI-события — тег-флоу. Данные из родителя — `useCtx().store.ctx`.
- View/Shape **stateless**: ноль состояния, ноль явных импортов.

**Kit-first** (`canon/components/kit-first.md`) + **модель компонентов** (`component-model.md`)
- Весь UI из kit. Потребитель — **props-only**: ноль raw-классов, ноль `<style>`/inline. Нужен вид, которого нет → расширяешь **kit** (owner-флоу), не обходишь классом.
- Композиция — **только в kit пресетом**. Карточка = сущность (данные по ключам; вид=пресет, полнота=слоты, точечно=props — НЕ ручной лайаут). Свои данные рисуются своими компонентами.

**Типы из zod** (`canon/principles/types-from-zod.md`)
- `z.infer<typeof schema>` — единственный источник домена. Ручной `interface`/`type` под домен запрещён. Соблазн «маленький тип для пропа» = сущность не выведена из своей Entity, СТОП.

**Модули + причина** (`canon/principles/modules-no-crutches.md`, `root-cause-not-symptom.md`)
- Cross-zone импорт реализации соседа — нарушение; связь через контракт/опубликованный API.
- Стоп-сигналы костыля (hardcoded / silent fallback / второй @source); 2 неудачных фикса → СТОП+корень. Gap соседнего пакета → фикс в **его** зоне, не workaround у себя.

**Ownership** (`canon/packages/ownership.md`)
- Только своя зона. Чужое → эскалация вверх к architect. Git: commit-only, push — architect.

---

## 📖 Read step-0

- `canon/architecture/ui-proxy-tag-flow.md` · `namespaces.md` — event-флоу, nested-неймспейсы.
- `canon/components/registration.md` · `tokens.md` — 2-слойная регистрация, frozen-токены.
- `canon/packages/anatomy.md` · `dependency-tiers.md` — структура зоны, portable-tier.
- `canon/principles/foundation-first.md` · `etalon-gate.md` — гейты старта и «готово».
- `canon/compliance/golden-rules.md` (severity: structural валит CI) · `linter.md`.
- Зона: свой `OWNERSHIP.md` / AI-anchor + `../shared-policy.md` (читают все первым).
