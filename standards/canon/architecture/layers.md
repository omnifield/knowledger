# Слои HCA

**HCA (Hyper-Controlled Architecture)**, философия *«UI is a Shadow»* — интерфейс это немая проекция логики; вся власть в Controller/Feature, общение с UI идёт через Proxy и meta-теги.

Слои снизу вверх. Запрещены **upward** и **horizontal** импорты ([import-rules](import-rules.md)). Любая композиция между сущностями — только в Widget.

| Слой | Папка | Что это | Обёртка |
|---|---|---|---|
| **Entity** | `entities/` | Domain data layer: zod schema + defaults + meta. **БЕЗ UI**. Single source of truth для сущностей (User, Product, Order). Возвращает plain config (не component). | `Entity(({ zod }) => ({ schema, defaults? }))` |
| **View** | `views/` | Stateless UI в виде JSX. Только Solid JSX + `data-meta`. Не знает про FSM, API, router. | `View((Ui, props?) => JSX)` |
| **Shape** | `shapes/` | **Presentation**: как нарисовать сущность через batch-template (`Ui.DataTable`, `Ui.List`). Ссылается на Entity за schema/defaults. Two-phase: фаза-1 bind, фаза-2 config. | `Shape((ui, { zod }) => ({ schema, as }), (ui, props) => ({ ...config }))` |
| **Controller** | `controllers/` | Поведение на FSM-схеме. Через Proxy перехватывает `onClick`/`onInput` у потомков. | `Controller((services) => ({ initial, states }))` |
| **Feature** | `features/` | Domain logic / side effects. **Только тут разрешены API.** Валидирует через `Entities.X.schema.parse(...)`. | `Feature((services) => ({ initial, states }))` |
| **Widget** | `widgets/` | Композиция View/Shape + Controller/Feature. **Единственное место, где можно «склеивать».** | `Widget((Ui, store?, props?) => JSX)` |
| **Page** | `pages/` | Корневой layout, оборачивает Widget. | `Page((Ui, props?) => JSX)` |

Имена `Page/Widget/View/Shape/Controller/Feature/Entity` — **глобальные**, инжектятся авто-импортом. В коде их **не импортируют**.

## Entity vs Shape — критическое различие

| Concern | Entity | Shape |
|---|---|---|
| Что описывает | Сущность (User row) | Презентацию (таблицу users) |
| Содержит UI template | ❌ | ✅ (`as: ui.DataTable` / `ui.List`) |
| Содержит columns / itemAs | ❌ | ✅ |
| Reusable across presentations | ✅ | ❌ (specific к layout) |
| Возвращает | plain config object | component-функцию |

Правило: **сущность → Entity. Как нарисовать → Shape.** Shape ссылается на Entity (`schema: Entities.Users.schema`).

## Что param vs global

- **`Ui`** (UI-kit примитивы) — приходит **первым параметром** обёртки. Per-instance проксированная копия под текущий `ControllerContext` (event-binding, meta-registration). Разные ref у разных instance'ов. См. [ui-proxy-tag-flow](ui-proxy-tag-flow.md).
- **`Views/Widgets/Shapes/Controllers/Features/Entities`** — **глобалы** (один stable ref на приложение). Доступны прямо из тела фабрики без объявления в args.
- **`useCtx`** — глобал; в Views/Widgets даёт доступ к `ControllerContext` (reactive store).
- **`services`** (Controller/Feature) — параметр, per-instance (`api`/`store`/`state`, деструктуризация `{ router, utils, <pkg>Api }`).
- **`store`** — 2-й аргумент Widget; **не** `useCtx` в user-виджете; store опционален → гард.
- **`props`** — опциональный аргумент View/Widget/Page (template-pattern).

`zod`/`utils` инжектятся **только** в logic/data-слои (Entity/Controller/Feature) объектом (tool-injection); в UI-слои (View/Shape/Widget/Page) — НЕ инжектятся.

## Конвенция export

Каждый файл слоя **обязан** заканчиваться `export default <Name>;`. Без него TS не увидит default → сломается типизация slot-кодгена и навигация Ctrl+Click.
