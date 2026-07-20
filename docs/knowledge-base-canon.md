# knowledge-base-canon — рыночный канон KB (для движка knowledger)

**Что это.** Дистиллят рыночного ресёрча OSS-продуктов баз знаний — модель данных,
API-конвенции, хранилище, связи, «самородки» и анти-паттерны. Аналог tasker'ского
`hub/patterns/task-manager-canon.md`: сначала изучаем рынок, потом строим Go-движок
по канону (и на этом же доке тестим функционал). Backend-first; одно крепкое ядро,
морды/интеграции — позже.

**Метод.** Fan-out веб-поиск по 5 углам → 24 источника → 114 утверждений →
adversarial-верификация (3 голоса, нужно 2/3 refute чтобы убить) → 23 подтверждено,
2 опровергнуто. Уровень доверия у каждого пункта помечен. Источники — в конце.
Дата: 2026-07-20.

**Однородность.** Внутренняя архитектура зеркалит tasker (Go-бэк, раскладка
`cmd/ · internal/{httpapi,service,store,config,webhook} · migrations · sqlc`, фронт
`vite+solid` потом). Рынок берём за **паттернами**, не за стеком.

---

## 1. Что knowledger уже наследует от tasker

tasker — родственный движок; его node-модель уже несёт половину канона. Наследуем
как есть (не переизобретаем):

- **Dual-id** — узел резолвится и по UUID, и по стабильному `key` (`WS-n`).
- **Рекурсивный type-less node** — `parent_id` self-FK = дерево; per-workspace scope.
- **Typed relations** (cross-workspace) отдельно от дерева; derived-rollup внутри дерева.
- **Activity timeline** с 1-го дня; **cross-product предложки → inbox → accept-gate**.
- REST под нативным префиксом (`/knowledger/…`), Bearer-actor, bind `0.0.0.0`.

Ресёрч ниже **подтверждает** эти выборы рынком и добавляет KB-специфику, которой у
tasker нет: контент/блоки, backlinks, transclusion, версии, полнотекст/семантика.

---

## 2. Кросс-каноничные бест-практис (подтверждено рынком)

### 2.1 Идентичность и адресация — **dual-id** [высокая увер., 3-0]
Адресуй контент **стабильным машинным ID** (UUID/int) и параллельно неси
**человеческий ключ** (name/slug), оба — в ответах API, чтобы ссылки переживали
переносы и переименования. Конвергенция трёх независимых продуктов:
- Logseq: `:block/uuid` (`:db.unique/identity`) + lowercase `:block/name` для страниц.
- BookStack: числовой `id` + URL-slug, оба в ответе.
- Outline: UUID `id` + 10-символьный `urlId`, lookup по любому.
- Directus: 4 неизменяемых PK-стратегии на коллекцию (auto-int/big-int/UUID/manual).

→ **Для knowledger:** ровно tasker'ский dual-id. Подтверждён как канон, не совпадение.

### 2.2 Контент — рекурсивный type-less node/block + дробный индекс [высокая, 3-0]
Документы, заголовки, абзацы, списки, страницы — всё **блоки** одного рекурсивного
типа, каждый адресуемый и переносимый. Иерархия — через `parent` + owning-page ref;
порядок сиблингов — **fractional index** (реордер без перенумерации).
- SiYuan: «everything is a block», перманентные ID → split/move/rename не рвут ссылки.
- Logseq: `:block/parent`, `:block/page`, `:block/order` = дробный индекс
  (`clj-fractional-indexing`, ключи вида `a0/a1/a0V`, лексикографически сортируемы).

→ **Для knowledger:** узел уже type-less. Добавить **колонку дробного порядка** для
сиблингов (не int-позиция). Гранулярность «блока» внутри статьи — как у tasker,
**отложить** (v1: узел = статья с телом; блочную модель — позже, см. §6).

### 2.3 Связи — forward-рёбра, backlinks **derived** [высокая, 3-0]
Links/backlinks/tags — cardinality-many ref-атрибуты на узле; backlinks **вычисляются
обратным запросом**, а не хранятся отдельной синхронизируемой таблицей.
- Logseq (DataScript/EAV): `:block/refs`, `:block/tags` (ref, many); backlinks =
  reverse-lookup `:block/_refs`. Структурно нет edge-таблицы-дубля.

→ **Для knowledger (SQL-перевод):** одна таблица forward-ребёр (`ref`: from→to, kind),
проиндексированная под reverse-lookup; backlinks и членство-в-теге = **запросы**, не
денормализованная копия. Держать в синхроне нечего.

### 2.4 Ссылки — **живые указатели на ID**, не текст-копии [высокая, 3-0]
Ref — указатель на стабильный ID, резолвится на рендере. Деградация ID = ссылка
превращается в мёртвый текст (Logseq #4491: cut-paste затёр `((uuid))` в plain text).
→ **Для knowledger:** хранить ссылки как ID-рёбра; **сохранять ID при copy/move**.

### 2.5 Transclusion / переиспользование — first-class [высокая для Docmost 3-0; SiYuan 2-1]
Один и тот же блок появляется в нескольких местах и остаётся в синхроне.
- Docmost **Synced Blocks** (shipped v0.90.0).
- SiYuan: embed с настраиваемым рендером (scoped к экспорту — не обобщать).
→ **Самородок:** transclusion = **reference-узел на канонический блок**, резолвится
(и опц. пере-оборачивается) на чтении. Не копия.

### 2.6 Версии и история — отдельный first-class ресурс [высокая, 3-0]
История — **не** инлайн-поле документа, а сиблинг-сущность по parent-ID.
- Directus: `Activity` (аудит-таймлайн) и `Revisions` (снапшоты + revert) — отдельные
  системные ресурсы, отдельно от CRUD.
- Payload: отдельная коллекция `_<slug>_versions`; каждая версия ссылается на parent,
  хранит **полную копию** под `version` + метаданные (autosave, timestamps); операции
  `findVersions/findVersionByID/restoreVersion`.
→ **Для knowledger:** `activity` уже есть (tasker). Добавить **`revision`** (node_id →
parent, snapshot полной копии, meta, created_at) + явные `versions`/`restore` в API.

### 2.7 Scope — spaces/collections-дерево, **не** фикс-глубина [высокая, 3-0; refuted строгая иерархия]
Гибкое `scope → вложенные узлы` произвольной глубины (Docmost: Spaces → nestable Pages,
drag-drop). Строгая фикс-глубина (якобы BookStack Shelves>Books>Chapters>Pages) —
**опровергнуто (1-2)**, это анти-паттерн.
→ **Для knowledger:** per-workspace scope + произвольная вложенность узлов (совпадает
с моделью). Не хардкодить уровни.

### 2.8 API-конвенции [высокая, 3-0]
- **Пагинация** offset/limit, эхом в конверте + `nextPath`-шорткат (Outline).
- **Поиск** — выделенный first-class эндпоинт, результаты **авто-ограничены scope
  токена** (Outline `documents.search`).
- **`depth`-параметр** — насколько глубоко разворачивать связанные объекты
  (Payload `?depth=2`, дефолт 2, `0`=только ID). Для графа/бэклинков — чистый способ
  управлять раскрытием связей в ответе.
- **Одна query-грамматика на все транспорты** — Payload делит один язык запросов между
  in-process Local API, REST и GraphQL. → Определить фильтр-грамматику в сервис-слое,
  отдавать через in-process + REST (+ опц. GraphQL), без per-transport дивергенции.

### 2.9 Schema-driven генерация (опц. самородок) [высокая, 3-0]
Directus: коллекция = таблица БД + метаданные; REST+GraphQL **авто-генерятся** из
схемы (database mirroring через общий AST), Relations — отдельный CRUD-ресурс (M2O/M2M/
M2A). → Для мостли-фикс node-модели: генерить CRUD/query-хендлеры из реестра
сущностей/связей, а не писать бойлерплейт на каждую сущность. Оценить, не переусложняет ли.

---

## 3. Рекомендованная форма Go-движка (backend-first)

### 3.1 Сущности / SQL-схема (Postgres + sqlc, goose-миграции — как tasker)
```
workspace(id, key, name, …)                    -- scope-тир (продукт|концерн), как tasker
node(id uuid pk, workspace_id, key,            -- dual-id (unique(workspace_id, seq))
     parent_id self-fk null,                   -- дерево (adjacency list)
     kind text null,                           -- type-less; kind опц. (article|section|…)
     title, body,                              -- v1: узел = статья с телом
     ord text,                                 -- ДРОБНЫЙ индекс порядка сиблингов
     created_at, updated_at)
ref(id, from_node, to_node, kind,              -- forward-рёбра: link|tag|transclude|relation
    created_at)                                -- backlinks = reverse-query по (to_node)
tag(id, workspace_id, name, …)                 -- или теги как node/ref; решить в KNOW-6
revision(id, node_id→node, snapshot jsonb,     -- версии = отдельный ресурс (Directus/Payload)
         meta, created_at)
activity(id, node_id, kind, data, actor, …)    -- типизированный таймлайн (как tasker)
```
Индексы: `node(workspace_id, parent_id)` для дерева; `ref(to_node, kind)` для backlinks;
GIN по tsvector для полнотекста (см. §3.3). Дерево — **adjacency list + recursive CTE**
(бенч: 12ms vs 18ms materialized-path на 50k/depth-6, PG16) [blog, средняя увер.].

### 3.2 API (`/knowledger/…`, зеркалит стиль tasker)
- `GET·POST /workspaces`, `GET /workspaces/{ws}`, `…/tree?depth=N`
- `GET·POST /workspaces/{ws}/nodes` (фильтры `?parent=`, `?tag=`), `GET·PATCH·DELETE /nodes/{key}` (dual-id)
- `GET /nodes/{key}/tree?depth=N`, `GET /nodes/{key}/children`
- `GET·POST /nodes/{key}/refs` (обогащённые: направление + backlinks + cross-ws), `DELETE /refs/{id}`
- `GET /nodes/{key}/backlinks` (derived reverse-query)
- `GET /search?q=…` — выделенный, auto-scoped токеном; `depth` для раскрытия связей
- `GET·POST /nodes/{key}/versions`, `POST /nodes/{key}/versions/{id}/restore`
- `GET·POST /nodes/{key}/activity`; `…/tags`
- **Наследовано от tasker:** `POST /workspaces/{ws}/proposals`, `GET …/inbox`,
  `POST /nodes/{key}/accept|decline` — cross-product наполнение KB через accept-gate.
- Конвенции: dual-id везде; offset/limit + конверт + `nextPath`; Bearer-actor; `/healthz` без auth.

### 3.3 Поиск [РЕКОМЕНДАЦИЯ, не верифицировано вживую — см. caveats]
Postgres-native гибрид: **BM25/tsvector + GIN** (лексика) ⊕ **pgvector + HNSW**
(семантика), слияние через **RRF** (reciprocal rank fusion). Старт — только полнотекст
(tsvector/GIN); вектор/эмбеддинги — за флагом, когда понадобится семантика. Не тащить
внешний индекс, пока Postgres хватает.

---

## 4. Steal-these (ранжировано)

1. **Dual-id** (уже в tasker; рынком переподтверждён как канон №1).
2. **Рекурсивный type-less node + дробный индекс** порядка сиблингов.
3. **Refs = forward-рёбра, backlinks derived** (одна edge-таблица, reverse-индекс).
4. **Transclusion = reference-узел**, резолвится на чтении (Docmost/SiYuan).
5. **Версии + activity — отдельные ресурсы** с `restore` (Directus/Payload).
6. **API: `depth`-раскрытие + конверт пагинации + выделенный scope-search**.
7. **Одна query-грамматика на все транспорты** (in-process + REST + опц. GraphQL).
8. **Schema-driven генерация хендлеров** из реестра сущностей (Directus) — опц., оценить ROI.

---

## 5. Анти-паттерны / грабли

- **Фикс-глубина иерархии** (Shelves>Books>Chapters>Pages) — жёстко, опровергнуто рынком;
  бери произвольную вложенность в scope.
- **Ссылки как текст-копии** → мёртвые ссылки при copy/move (Logseq #4491). Только ID-рёбра.
- **Денормализованная backlink-таблица** для синхрона → лучше derive обратным запросом.
- **Инлайн-версии в документе** → выноси в отдельный ресурс с restore.
- **RPC-per-method** (Outline `POST /api/:method`) — рабочий, но ломает REST-кэш/verb-семантику;
  для наполнения бери resource-REST (канон tasker), RPC терпим лишь для search/bulk.
- **Мид-миграция как источник** — «Logseq» = движущаяся цель (file-версия `:block/left`
  linked-list ↔ DB-версия дробный индекс); не цитировать как единое.

---

## 6. Что откладываем (backend-first, расширяемость заложена)

- **Блочная гранулярность внутри статьи** — v1 узел = статья с телом; блок-модель
  (как tasker отложил структурные блоки) — когда приземлится.
- **Морды** — `vite+solid` как у tasker, отдельным `web/`, после ядра.
- **CRDT/offline-sync** — для backend-first преждевременно; открытый вопрос (см. ниже).
- **Семантический/вектор-поиск** — за флагом поверх полнотекста.

---

## 7. Caveats (честно о покрытии)

- Покрытие неровное: сильные **первичные** источники по headless (Directus, Payload,
  Outline) и Logseq; тоньше по wiki-архетипу (BookStack/Docmost — по одному источнику);
  **не выжило** верифицированных утверждений по Wiki.js, XWiki, Confluence, AppFlowy,
  Anytype, AFFiNE, Trilium, Dendron, Foam, Strapi — канон выведен из подмножества.
- Часть модельных пунктов опирается на вторичку (DeepWiki для Logseq/SiYuan), но Logseq
  сверен с первичным кодом (`schema.cljs`, `clj-fractional-indexing`).
- SiYuan transclusion прошёл лишь 2-1 и scoped к экспорту — **не обобщать**.
- **Нет** верифицированных данных по: внутренностям полнотекст-vs-вектор поиска,
  CRDT/offline-sync, миграционному тулингу, permission-модели/мультитенантности.
  §3.3 и §6 (поиск, sync) — рекомендации, не факты.
- Опровергнуто и исключено: строгая 4-уровневая иерархия BookStack; 6-режимный
  `BlockRefMode` SiYuan.

### Открытые вопросы (в KNOW-6/7/8)
1. Полнотекст+семантика: tsvector/GIN vs внешний индекс vs pgvector — конкретная схема и индексы.
2. Offline-sync PKM (CRDT vs LWW vs OT) — нужен ли backend-first движку или отложить к блок-слою.
3. Точная SQL-форма EAV→Postgres: single node-таблица + ref-таблица + дробный порядок vs JSONB-refs; индексы под reverse-backlink.
4. Permission/мультитенантность wiki-архетипа (Wiki.js/XWiki/Confluence) — не выжило утверждений, добрать.

---

## 8. Источники (первичные — жирным)

- **Logseq** — issue [#4491](https://github.com/logseq/logseq/issues/4491), [#1389](https://github.com/logseq/logseq/issues/1389); [clj-fractional-indexing](https://github.com/logseq/clj-fractional-indexing); DeepWiki schema (вторичн.).
- **BookStack** — [demo API docs](https://demo.bookstackapp.com/api/docs).
- **Outline** — [developers](https://www.getoutline.com/developers), [openapi](https://github.com/outline/openapi).
- **Docmost** — [docs](https://docmost.com/docs/), [v0.90.0 release](https://github.com/docmost/docmost/releases/tag/v0.90.0).
- **Directus** — [collections/data-model](https://docs.directus.io/app/data-model/collections), [API reference](https://docs.directus.io/reference/introduction).
- **Payload** — [versions](https://payloadcms.com/docs/versions/overview), [concepts](https://payloadcms.com/docs/getting-started/concepts), [depth](https://payloadcms.com/docs/queries/depth).
- SiYuan — DeepWiki block model (вторичн.), openalternative.co.
- Поиск/схема (blog): ParadeDB hybrid-search, pgvector RRF (jkatz05), иерархия в Postgres (leonardqmarcq), recursive-CTE (cybertec).

> Канон-док. Правится по мере ресёрча/реализации. Финальный дом канона возможно —
> `hub/patterns/` (как у tasker); согласовать с devopser. На нём же тестируем функционал.
