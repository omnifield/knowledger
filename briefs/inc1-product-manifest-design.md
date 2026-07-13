---
title: "Инкремент 1 — product-manifest контракт (дизайн)"
status: draft
owner-repo: commons
---

# Дизайн — Инкремент 1: контракт тонкого product-manifest

| | |
|---|---|
| **Автор** | knowledger-архитектор |
| **Заказ** | [`inc1-product-manifest.md`](inc1-product-manifest.md) |
| **Основание** | блюпринт [`workspace-platform-draft.md`](../blueprints/workspace-platform-draft.md) §3 (два слоя) + §9 (решения user) |
| **Статус** | draft — на ревью user + согласование границы dev-services с devopser (см. §7) |

---

## 0. TL;DR (тейки для user)

- **Файл:** `omnifield.yaml` в **корне репо каждого продукта** (JS и не-JS одинаково). Отдельный чистый файл, не `project.json`.
- **Источник правды схемы:** Zod в `@omnifield/contract-manifest` (канон [types-from-zod](../standards/canon/principles/types-from-zod.md)). Из неё эмитится `omnifield.schema.json` — им валидируют не-JS продукты и подсвечивает редактор. Один авторский источник → один кросс-язычный артефакт.
- **Тонкость держится структурно:** схема `.strict()` — **любое лишнее поле = ошибка валидации**, а не «договорились не дописывать». Соблазн утечь расширенное внутрь манифеста ловит валидатор, не память.
- **Три обязательных поля** (`apiVersion`·`name`·`type`), остальное — по необходимости читателя. Поля нет, если его **никто не читает** (матрица §4).
- **Граница с dev-services:** манифест несёт только **шлюзо-видимую поверхность** (route+port). Внутренние dev-сервисы (autostart/toggle) остаются в продукте (`devbox.services.json`, devopser). Пересечения нет — манифест ссылается на сервис по имени, не переопределяет его lifecycle. **Требует финального ок devopser (§7).**

---

## 1. Форма и место (Q1, Q2)

- **Имя файла:** `omnifield.yaml`. **Место:** корень репо продукта.
- **Формат:** YAML — человекопишущийся, комментируемый, язык-нейтральный. Продукт **не тянет npm-пакет**, чтобы нести манифест: это просто файл (канон [containers-only](../standards/canon/principles/containers-only-execution.md) — «файлы на машине», движняк валидации — на въезде внутри контейнера).
- **Универсальность не-JS (Q2 — подтверждено):** Go-chater, Python-backend несут ровно тот же `omnifield.yaml` в корне. Валидация — **ingest-gate** (devopser-агрегация / brainer-чтение) прогоняет файл против `omnifield.schema.json`; в любом языке есть JSON-Schema-валидатор. Продукт не обязан быть JS.
- **Редакторская подсказка (опциональна):** первой строкой
  ```yaml
  # yaml-language-server: $schema=<workspace-origin>/omnifield.schema.json
  ```
  `$schema` указывает на схему, **отданную своим origin воркспейса** (или вендоренную из `@omnifield/contract-manifest`), НЕ на публичный CDN — self-host-first / air-gapped ([integration-protocol](../standards/integration-protocol.md): внешние ресурсы не хардкодим). Строка — только сахар для IDE; нормативная валидация — на ingest-gate, не в редакторе.

---

## 2. Схема-контракт (нормативная — Zod)

Источник правды. Живёт в `@omnifield/contract-manifest`, публикуется как `@omnifield/contract-*` (README commons). `.strict()` — структурный guard тонкости.

```ts
import { z } from 'zod'

export const ProductType = z.enum(['frontend', 'backend', 'fullstack', 'service'])

/** Один шлюзо-видимый маршрут. НЕ описывает lifecycle сервиса — только как достучаться. */
export const Route = z.object({
  path: z.string().startsWith('/'),            // location-префикс в gateway (nginx)
  port: z.number().int().positive(),           // порт КОНТЕЙНЕРА, отдающий этот путь
  service: z.string().optional(),              // имя compose-сервиса, если != name; ТОНКАЯ ссылка, не определение
}).strict()

export const Integration = z.object({
  scopes: z.array(z.string()).default([]),     // роль-скоупы участия (role-tier, блюпринт §7)
  spawnEligible: z.boolean().default(false),   // хаб вправе спавнить агентов в контейнер продукта
  deps: z.array(z.string()).default([]),       // имена ДРУГИХ ПРОДУКТОВ (не npm-пакетов)
}).strict()

export const ProductManifest = z.object({
  apiVersion: z.literal('omnifield.dev/v1'),   // мажор контракта (см. §5)
  name: z.string().regex(/^[a-z][a-z0-9-]*$/), // уникальный id; = базовое имя compose-сервиса
  type: ProductType,
  title: z.string().optional(),                // человекочитаемая метка карты хаба; default = name
  description: z.string().optional(),          // подпись карточки хаба
  reach: z.object({ routes: z.array(Route).min(1) }).strict().optional(),
  integration: Integration.default({ scopes: [], spawnEligible: false, deps: [] }),
})
  .strict()                                    // ← лишний ключ = ошибка. Утечку расширенного ловит валидатор.
  .superRefine((m, ctx) => {
    // reach обязателен для того, что вообще ходит через дверь
    if ((m.type === 'frontend' || m.type === 'fullstack') && !m.reach) {
      ctx.addIssue({ code: 'custom', path: ['reach'],
        message: `type '${m.type}' обязан объявить reach.routes` })
    }
  })

export type ProductManifest = z.infer<typeof ProductManifest>
```

**Кросс-язычный артефакт:** из этой Zod-схемы `zod-to-json-schema` эмитит `omnifield.schema.json`, который кладётся в пакет. Не-JS продукты/валидаторы и `$schema` редактора используют **его**. Ровно паттерн [go.md](../standards/canon/languages/go.md) §Контракты (из одной схемы владельца — типы всех языков), с Zod как человекопишущимся истоком (канон types-from-zod), а JSON-Schema — как выведенным артефактом. Ноль дрейфа.

### Обязательные vs опциональные

| Поле | Обяз.? | Default (тончайшее поведение) |
|---|---|---|
| `apiVersion` | ✅ | — |
| `name` | ✅ | — |
| `type` | ✅ | — |
| `reach.routes[]` | ✅ для `frontend`/`fullstack`; опц. для `backend`/`service` | — (headless-`service` без двери валиден) |
| `title` / `description` | ⬜ | `title` = `name`; описания нет |
| `integration.scopes` | ⬜ | `[]` |
| `integration.spawnEligible` | ⬜ | `false` |
| `integration.deps` | ⬜ | `[]` |

Отсутствие блока `integration` целиком = все дефолты = тончайший манифест (3 обязательных поля + `reach`).

---

## 3. Граница «манифест ↔ внутри продукта» (DoD)

| Концерн | Где | Почему |
|---|---|---|
| identity (`name`·`type`·`title`) | **манифест** | это и есть визитка |
| шлюзо-видимая поверхность (`reach.routes`) | **манифест** | интерфейс продукт↔хаб/gateway |
| индекс интеграции (`scopes`·`spawnEligible`·`deps`) | **манифест** | ровно то, что хабу надо для карты/спавна |
| Dockerfile / compose-сервис (build·command·env·volumes) | **внутри продукта** | самодостаточная runnable-конфигурация (блюпринт §3) |
| dev-сервисы lifecycle (autostart·toggle·внутренние сервисы) | **внутри продукта** (`devbox.services.json`, devopser) | расширенная настройка живёт внутри; §7 |
| runtime-конфиг · секреты · feature-flags | **внутри продукта** | хаб их не видит и не хранит |
| зависимости пакетов | **внутри продукта** (`package.json`/`go.mod`/`pyproject`) | `deps` манифеста = продукты, НЕ пакеты |

Правило разбора: если что-то нужно, **чтобы продукт поднялся сам** `docker compose up` — это внутри продукта. Если нужно, **чтобы хаб/gateway его интегрировали** — это манифест. `docker compose up` одного продукта НЕ должен читать манифест (канон самодостаточности).

---

## 4. Кто что читает (Q3) — матрица держит тонкость

Поле легитимно, **только если у него есть читатель**. Предложенное поле без читателя = отклоняется (это и есть механизм против scope-creep, дополняющий `.strict()`).

| Поле | devopser (compose+gateway) | хаб (реестр/карта) + brainer-продукт (спавн) |
|---|---|---|
| `apiVersion` | ✅ валидация | ✅ валидация |
| `name` | ✅ ключ сервиса | ✅ id узла карты |
| `type` | ✅ выбор шаблона/группировки | ✅ иконка/группа карты |
| `title` / `description` | — | ✅ карточка |
| `reach.routes[].path` | ✅ nginx `location` | ✅ open-in ссылка |
| `reach.routes[].port` | ✅ upstream | ~ deep-link |
| `reach.routes[].service` | ✅ цель upstream | — |
| `integration.scopes` | — | ✅ проводка ролей (role-tier §7) |
| `integration.spawnEligible` | — | ✅ кнопка спавна |
| `integration.deps` | ✅ порядок старта (`depends_on`) | ✅ рёбра графа карты |

- **devopser** читает `identity` + `reach` + `deps` → генерит compose и проводит gateway (инкремент 2). НЕ трогает `scopes`/`spawnEligible`.
- **brainer** читает `identity` + `reach` + весь `integration` → живая карта, маршрутизация, спавн (инкремент 4).
- Оба читают **только**; никто не пишет манифест — его пишет продукт. Аггрегация, не генерация внутренностей (блюпринт §3).

---

## 5. Версионирование контракта (Q4)

- **Схема публикуется как `@omnifield/contract-manifest`** (semver), README commons: contracts = версионируемые мосты, `@omnifield/contract-*`.
- **Экземпляр манифеста пинит мажор** через `apiVersion: omnifield.dev/v1`. Ломающая смена поля = major-bump пакета + новый `apiVersion` = **видимый брейк** и точка синхронизации ([integration-protocol](../standards/integration-protocol.md) §2), не тихий дрейф.
- Аддитивные опциональные поля = minor. Ingest-gate потребителя пинит нужный мажор через диапазон зависимости.
- **Не** наследуется через standards (это версионируемый шов между репо, а не дисциплина) — живёт в `contracts/`, не в `standards/`.

---

## 6. Примеры манифестов (DoD)

### weber — frontend

```yaml
# yaml-language-server: $schema=<workspace-origin>/omnifield.schema.json
apiVersion: omnifield.dev/v1
name: weber
type: frontend
title: Weber
description: Веб-редактор карточек знаний
reach:
  routes:
    - path: /weber
      port: 5173
integration:
  deps: [brainer]        # ходит в brainer API; spawnEligible/scopes — дефолты (тончайше)
```

### brainer — fullstack

```yaml
apiVersion: omnifield.dev/v1
name: brainer
type: fullstack
title: Brainer
description: Пульт агентов — карта и спавн Claude-сессий (ПРОДУКТ, втягивается в хаб)
reach:
  routes:
    - path: /brainer          # продукт за СВОИМ префиксом; корень / = лендинг ХАБА, не продукта
      port: 3500              # реальный внутренний listening-порт (nginx upstream http://brainer-web:3500)
      service: brainer-web
    - path: /api/brainer
      port: 8010
      service: brainer-svc
integration:
  scopes: [workspace, control-plane]
  spawnEligible: true
```

---

## 7. Граница с dev-services (Q5) — согласовать с devopser

**Тезис:** манифест несёт **только шлюзо-видимую поверхность**; lifecycle dev-сервисов — внутри продукта. Дублирования нет.

- **В манифест идёт** только то, что должно пройти через дверь: `reach.routes[]` = `path`+`port` (+ `service`, если имя compose-сервиса отличается от `name`). Это **ссылка**, не определение.
- **Внутри продукта остаётся** (`devbox.services.json`, дизайн devopser [`devbox-first-run-dx.md`](https://github.com/omnifield/devopser)): команда запуска, autostart, toggle, env, а также **внутренние сервисы, которые НЕ выставляются через gateway** (локальный redis, воркер, watcher) — их манифест не упоминает вовсе.
- **Точка стыка (единственная):** сервис, который И dev-service, И шлюзо-видим. Правило анти-дубля: сервис **определяется один раз** в compose/`devbox.services.json` продукта; манифест лишь **ссылается на него по имени** (`reach.routes[].service`) + отдаёт `path`+`port`. Lifecycle/command/env в манифест не копируются. **Имя в `reach.routes[].service` обязано совпадать** с именем сервиса, которое знает compose devopser.

**🔶 Открытый пункт для devopser (liaison):** подтвердить, что (а) `devbox.services.json` — правильный дом внутренних dev-сервисов, (б) matching по имени `service` ↔ compose-сервис устраивает генерацию инкремента 2, (в) ничего из `reach` не дублирует `devbox.services.json`. До ок devopser этот раздел — draft.

---

## 8. Поглощение старого (тонкий индекс, не хранилище)

- Хардкод `BRAINER_REPOS` (brainer `config.py`) и devopser `registry/products.md` **выводятся** сканом всех `omnifield.yaml` воркспейса — как **тонкий индекс интеграции**, НЕ копия всех настроек продукта (блюпринт §3). Реестр перестаёт быть рукописным источником; источник = per-product манифест.

---

## 9. DoD-чеклист

- [x] Контракт-схема: поля · типы · required/opt (§2 Zod + таблица).
- [x] Пример манифеста для 2 продуктов: weber (frontend), brainer (fullstack) (§6).
- [x] Явная граница «в манифесте / внутри продукта» (§3).
- [ ] Согласовано с devopser-дизайном dev-services без дублирования — **ждёт ок devopser** (§7).

## Связь

- Блюпринт [`workspace-platform-draft.md`](../blueprints/workspace-platform-draft.md) §3, §6, §9.
- Канон: [types-from-zod](../standards/canon/principles/types-from-zod.md), [ecosystem-self-sufficiency](../standards/canon/principles/ecosystem-self-sufficiency.md), [containers-only-execution](../standards/canon/principles/containers-only-execution.md), [go.md](../standards/canon/languages/go.md), [integration-protocol](../standards/integration-protocol.md).
- Заказ: [`inc1-product-manifest.md`](inc1-product-manifest.md).
