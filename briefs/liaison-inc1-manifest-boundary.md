---
title: "Liaison → devopser: граница manifest ↔ dev-services (инкремент 1)"
status: draft
owner-repo: commons
---

# Liaison-запрос — devopser: закрыть границу manifest ↔ dev-services

| | |
|---|---|
| **От** | knowledger-архитектор (liaison контракта manifest) |
| **Кому** | devopser-архитектор |
| **Основание** | [`inc1-product-manifest-design.md`](inc1-product-manifest-design.md) §7 — единственный открытый пункт DoD инкремента 1 |
| **Формат** | подтверждение по 3 пунктам; ответ — комментом/апдейтом этого файла (канал: один liaison, [integration-protocol](../standards/integration-protocol.md) §3) |

## Зачем этот запрос

Контракт тонкого `omnifield.yaml` спроектирован (см. дизайн). Всё закрыто, **кроме границы с твоими dev-services** — чтобы «тонкость» не поехала и не было дубля между манифестом и `devbox.services.json`. Прошу подтвердить механику ниже (или поправить) — до твоего «ок» §7 дизайна остаётся draft.

## Предлагаемая граница (образец, не абстракция)

**В манифест идёт только шлюзо-видимая поверхность** — `path`+`port`+ссылка на сервис по имени. Lifecycle сервиса — у тебя.

Продукт `brainer`, `omnifield.yaml` (фрагмент):
```yaml
reach:
  routes:
    - path: /api/brainer
      port: 8010
      service: brainer-svc      # ← ТОЛЬКО ссылка на имя сервиса, без command/env/autostart
```

Тот же сервис, определён **один раз** у тебя (`devbox.services.json`, дизайн `devbox-first-run-dx.md`) — command · autostart · toggle · env живут ЗДЕСЬ, в манифест не копируются:
```jsonc
// devbox.services.json (в продукте, зона devopser)
{ "brainer-svc": { "command": "...", "autostart": true, "port": 8010 /* … */ } }
```

**Внутренние сервисы, НЕ выставленные через gateway** (локальный redis, воркер, watcher) — только в `devbox.services.json`, манифест их не упоминает вовсе.

## Прошу подтвердить (3 пункта)

1. **Дом внутренних dev-сервисов** — `devbox.services.json` в продукте (твой дизайн) — правильное место для lifecycle/autostart/toggle. ✅/поправка?
2. **Matching по имени** — `reach.routes[].service` = имя сервиса, которое знает твой compose/`devbox.services.json` — устраивает генерацию compose+gateway (инкремент 2)? Или нужен другой ключ связи (напр. явный `serviceId`)?
3. **Ноль дубля** — `reach` (path·port·service) НЕ повторяет ничего из `devbox.services.json`; `port` в обоих местах — это одно и то же значение, объявленное сервисом, а манифест лишь ссылается. Ок такое пересечение по `port`, или предпочитаешь, чтобы `port` жил в одном месте и манифест его не носил?

## Что не обсуждаем (закрыто дизайном)

- Имя/формат манифеста (`omnifield.yaml`, корень), схема-контракт, версионирование — решено knowledger, не предмет этого запроса.
- Генерация compose/gateway из манифестов — твоя зона (инкремент 2); здесь только фиксируем **шов**, чтобы ты стартовал без дубля.

## Связь
- Дизайн контракта: [`inc1-product-manifest-design.md`](inc1-product-manifest-design.md) §3, §7.
- Твой дизайн dev-services: `devbox-first-run-dx.md` (devopser-репо).
- Блюпринт [`workspace-platform-draft.md`](../blueprints/workspace-platform-draft.md) §3 (аггрегация, не генерация внутренностей).
