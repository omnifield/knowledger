# Анатомия пакета

## Тезис

**Доменный пакет = АПП минус Page/Feature.** Авторится теми же HCA-обёртками, что и апп, и той же раскладкой слоёв-папок — но без корневого layout (Page) и без side-effect/API-слоя (Feature). Пакет — переиспользуемый домен, не приложение.

## Раскладка

```
packages/<domain>/src/
  core/            # cross-cutting координация: provider, context, coordination сверху-вниз
  entities/        # Entity — РЕАЛЬНАЯ (zod из shared-zod), single source of truth сущности
  views/           # stateless UI (JSX)
  shapes/          # presentation сущности (batch-template)
  widgets/         # композиция view/shape + controller
  controllers/     # поведение на FSM
  index.ts         # barrel
```

**НЕТ `pages/`, НЕТ `features/`.** Причина: Page = корневой layout приложения (не пакета); Feature = API/side-effects/навигация (это концерн аппа-потребителя, не переиспользуемого домена). Пакет, тянущий Feature, тянет рантайм и завязывается на конкретный апп — теряет переносимость.

## Обёртки — узкий barrel

Слои пакета берут обёртки из **`@omnifield/web-core/wrappers`** — barrel БЕЗ `Page`/`Feature` (enforce на уровне barrel, чтобы пакет физически не мог их использовать). Апп берёт полный barrel (с Page/Feature).

## Ключевые правила

- **Entity реальная** — `zod` из shared-zod, не plain config. Пакет владеет схемой домена.
- **`core/` — координация**, а не свалка. Cross-cutting (provider/context/coordination), не привязанное к одной сущности. Абстрактный переиспользуемый UI живёт в `views/`, не в «shared-складе с сущностями внутри».
- **Интерактив через тег-флоу** ([[../architecture/ui-proxy-tag-flow]]): пакетные View/Widget рисуются проксированным `Ui`, события всплывают в app-Feature по тегу. Пакет теги не навязывает — апп может перемапить (aliases). Именованная доменная семантика — через `useEmit` (opt-in).
- **Граница апп↔пакет — только структурная** (папка/пакет), не кодовая. Апп композирует домен, добавляет свой Page (layout) и Feature (API/навигацию).
- **Слоты — только под содержимое.** Слот включается, когда есть что в него положить; не «на будущее».
- **Импорты** — стандартные HCA-правила ([[../architecture/import-rules]]): no upward/horizontal; направление к `core/` как к точке координации.

## Чек-лист «пакет = эталон»

- [ ] Раскладка `core/ entities/ views/ shapes/ widgets/ controllers/`, без pages/features.
- [ ] Обёртки из `@omnifield/web-core/wrappers` (узкий barrel).
- [ ] Entity real (zod), сущность рисуется своими компонентами.
- [ ] Ноль cross-zone импортов реализации соседа.
- [ ] Проходит эталон-гейт ([[../principles/etalon-gate]]): тесты + трейсы + доки + раскладка.
