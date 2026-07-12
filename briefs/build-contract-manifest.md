# Brief-заказ — knowledger: собрать `@omnifield/contract-manifest` (инк1, код)

| | |
|---|---|
| **Адресат** | knowledger-owner/architect |
| **От** | devopser-архитектор (координатор), 2026-07-12 |
| **Основание** | хаб-ядро (`devopser/briefs/hub-core-design.md` §5) строится на манифестах, БЕЗ полумер (решение user) → инк1 обязан стать реальным ПЕРВЫМ. Дизайн готов (`briefs/inc1-product-manifest-design.md`), нужен код |
| **Зона** | knowledger (contract-пакет). Файлы `omnifield.yaml` в продуктах — след. шаг (product-owner'ы), НЕ здесь |

## Задача — реализовать пакет из СВОЕГО же дизайна
`briefs/inc1-product-manifest-design.md` §2 — там нормативная Zod-схема. Собрать пакет:

1. **Zod-схема** `@omnifield/contract-manifest`: `ProductManifest` (+ `Route`, `Integration`,
   `ProductType`), `.strict()` + `superRefine` (reach обязателен для frontend/fullstack) — ровно §2.
2. **Эмит `omnifield.schema.json`** через `zod-to-json-schema` → кладётся в пакет (не-JS продукты
   и `$schema` редактора валидируют им; ingest-gate любого языка).
3. **Пакет публикуемый** (`@omnifield/contract-*`, semver, `apiVersion: omnifield.dev/v1`) —
   версионируемый шов, не через standards (§5 дизайна).
4. **Типы из Zod** (канон types-from-zod): один авторский источник → кросс-язычный артефакт.

## ⚠️ Модель поправлена (учесть при примерах)
Твой дизайн я уже выправил под поправку user (хаб ≠ brainer; brainer — продукт). В §6-примерах
доведи (это же liaison-поправка-3):
- **brainer на `/brainer`, НЕ на `/`** (корень `/` = лендинг ХАБА, не продукта) — я уже поправил в дизайне;
- **порты в примерах = РЕАЛЬНЫЕ внутренние порты продуктов** (что dev-сервер реально слушает;
  проверено живьём: weber `5173`, brainer `3500`/`8010`). ⚠️ Это НЕ «канон номеров» — порты
  внутренние, на разных контейнерах могут совпадать. Единственное требование: `manifest.reach.port`
  === реальный listening-порт продукта (иначе сгенерённый nginx-upstream `http://<repo>:<port>`
  промахнётся; ловит port-consistency-gate, `devopser/briefs/liaison-inc1-manifest-boundary.md`).
  `registry/ports.md` = человек-зеркало этих чисел, не центральная выдача (со временем — выведенное
  зеркало манифестов). Бери в примерах реальные числа = copy-correctness.

## DoD
- Пакет собирается; `omnifield.schema.json` эмитится из Zod (ноль дрейфа схема↔тип).
- Валидирует пример-манифест (weber-frontend + brainer-fullstack на канонических портах/маршрутах).
- `.strict()` ловит лишнее поле (structural guard тонкости); `superRefine` ловит frontend без reach.
- Пакет готов к `dlx`/публикации — devopser hub-core будет им валидировать скан манифестов.

## Границы
- Только contract-пакет. НЕ клади `omnifield.yaml` в продукты (это заказ product-owner'ам следом,
  когда пакет готов — его пришлю через git в их репо).
- Порты/маршруты примеров = зеркало `devopser/registry/ports.md` (контракт), не выдумывать.

## Связь
- `briefs/inc1-product-manifest-design.md` — твой дизайн (первоисточник схемы).
- `devopser/briefs/hub-core-design.md` — потребитель (реестр хаба сканит манифесты этим пакетом).
- `blueprints/workspace-platform-draft.md` — инк1 (+ баннер-поправка модели вверху).
