# Canon — HCA-канон Omnifield

Архитектурный канон экосистемы. Декомпозирован по темам — **плоский список запрещён** (общак, разрастётся). Новое правило кладётся в свою тему отдельным файлом; тема раздувается → дели на под-темы. Индекс держим актуальным.

Верхний закон — `../POLICY.md`. Здесь — детали.

**Визуальный индекс канона** (навигация по темам + port-status) — GitHub-борд:
🔗 https://github.com/orgs/omnifield/projects/2 (текст правил — в файлах ниже, борд — карта).

## Темы

### `principles/` — философия «почему»
Основания, из которых выводится всё остальное.
- `modules-no-crutches.md` — модули не монолит; независимость; совместимость дизайном.
- `root-cause-not-symptom.md` — лечим причину; стоп-сигналы костыля.
- `etalon-gate.md` — «готово» = код+тесты+трейсы+доки+раскладка.
- `types-from-zod.md` — типы только из zod (`z.infer`); ноль ручных `interface`/`type` под домен на любом уровне (app/package/contract).
- `foundation-first.md` — известные дыры базы закрываются ДО старта разработки, **даже не-блокеры**; сначала база без дыр, потом фичи.
- `ecosystem-self-sufficiency.md` — продукты самодостаточны, зависят только от нашей эко; сторонние API/SDK = уровень интеграций; контракты свои, вендорские типы не пересекают границу адаптера.
- `dogfood-product-flow.md` — наш флоу = флоу юзера; мы — юзер №0; «нам неудобно» = продукт не готов; наш способ работы = один из пресетов, не хардкод.

### `architecture/` — HCA «как устроено»
Слои, обёртки, механики, правила импорта.
- `layers.md` — Entity·View·Shape·Controller·Feature·Widget·Page: назначение, сигнатуры обёрток, что param vs global.
- `import-rules.md` — no upward / no horizontal; композиция только в Widget; цепочка через `next()`.
- `ui-proxy-tag-flow.md` — проксированный `Ui`, тег-флоу событий (децентрализованный N→1), именованные события opt-in.
- `namespaces.md` — nested по структуре папок (`widgets/forms/auth` → `Widgets.Forms.Auth`).

### `packages/` — модель пакета
- `anatomy.md` — **пакет = апп минус Page/Feature**; слои-папки `core/ entities/ views/ shapes/ widgets/ controllers/`; узкий barrel `web-core/wrappers` (без Page/Feature).
- `ownership.md` — owner-зоны, OWNERSHIP.md, releasability, границы.
- `dependency-tiers.md` — **portable-tier**: пакеты «в мир» с нулём ecosystem-deps (renderer/canvas/utils); их могут импортить наши, они наших — нет; CI enforce.

### `components/` — UI/kit модель
- `component-model.md` — композиция ТОЛЬКО в kit пресетом; карточка=сущность (данные по ключам, слоты вкл/выкл, мапер в пресете); список=пресеты форм.
- `kit-first.md` — весь UI из kit; потребители props-only, **ноль raw-классов**; Tailwind сканит только kit.
- `registration.md` — 2-слойная регистрация (manifest + `Ui`-namespace) для рендерера/стора studio.
- `tokens.md` — token-set (frozen); ноль raw-style; состояние в JS + токены.

### `languages/` — языковые каноны (что не закрывается общими принципами)
Брифы архитекторов на язык не опираются (функционал/архитектура/структура); язык-детали = owner + канон отсюда.
- `go.md` — тулчейн-пины (go.mod toolchain), раскладка сервиса (cmd/internal/migrations, свой go.mod = граница extract), стиль (errors %w, context-first, интерфейс у потребителя), stdlib-first транспорт, goose+sqlc (types-from-schema), slog/env-only, table-driven+`-race`.

### `compliance/` — enforcement
- `golden-rules.md` — правила + severity (`error` structural валит CI / `warn` cosmetic).
- `linter.md` — как enforce'ится (AST-линтер, Vite-плагин, CI-gate).

## Порт-статус

🔵 = вычитано из оракула `capsule` (CLAUDE.md, docs/01-architecture, docs/_meta), доведено до эталона.
Наполняется послойно при миграции соответствующего пакета — канон едет **вперёд** пакета (enforcement первым, см. [Canon-first](compliance/golden-rules.md#canon-first)).
