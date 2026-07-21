# knowledger — REST API

KB-движок Omnifield: рекурсивное дерево type-less узлов (документы·секции·статьи) +
типизированные рёбра (`refs`) + теги + timeline (`activity`). Форма зеркалит tasker
(дерево + cross-product предложки), но узел knowledger — **контент** (`title`+`body`+`kind`),
а не задача: нет статусов/приоритетов/роллапа, зато есть граф ссылок и бэклинки.

## База и auth

- **Нативный префикс:** `http://knowledger:8040/knowledger/…` (сосед по docker-сети).
- **Через дверь хаба:** `http://gateway/api/knowledger/…` — дверь снимает `/api` →
  бьёт в `knowledger:8040/knowledger/…` (как tasker под `/tasker/`). Публичный вход.
- **UI:** `/knowledger/` (vite+solid) — только человеку; наполнение — через API.
- **Auth:** заголовок `Authorization: Bearer <handle>`. Token-stub — любой непустой
  handle проходит (реальная identity позже). Handle пишется как **actor** в activity —
  бери осмысленный (`Bearer devopser`), чтобы в истории было видно, кто наполнял.
  Пустой/без `Bearer ` → `401 {"error":"missing bearer token"}`.
- `healthz` — без auth.

## Модель

`workspace → node` (рекурсия через `parent_id`, любая глубина, single-tree — узел не
паренится в чужой workspace).

- **Workspace** — scope (продукт ИЛИ концерн). `key` короткий, `^[A-Z][A-Z0-9]*$`
  (UPPER, без `-`/пробелов; на входе нормализуется в UPPER). Уникален.
- **Node** — рекурсивный **type-less** узел. `key` = `<WS>-<n>` (напр. `KNOW-12`)
  генерится сам, **стабильный/immutable** (по нему ссылаются извне; переезд/переименование
  key не меняют — dual-id). `kind` — необязательный хинт (doc/section/…), НЕ enforced-тип.
  `body` — тело (текст/markdown). `ord` — дробный индекс порядка среди сиблингов.
- **Dual-id:** везде, где путь `{key}`, принимается И UUID, И стабильный key — резолвятся
  в один узел.
- **Ref** — направленное ребро `from → to`, одна forward-таблица на всё. Kinds:
  `link` (свободная ссылка), `transclude` (живой эмбед контента), `relates` (ненаправленная
  ассоциация), `blocks` (from блокирует to), `depends_on` (from зависит от to),
  `duplicate` (from дублирует to). Cross-workspace рёбра разрешены. **Бэклинки** —
  выводятся обратным запросом по `to_node`, отдельной таблицы нет.
- **Tag** — per-workspace метка, вешается на узел (membership node↔tag, НЕ ref).
- **Activity** — типизированный timeline узла (`kind`: commented/created/…; `data` — произвольный JSON).
- **Proposal** — cross-product предложка: пишется в inbox ЧУЖОГО workspace (`origin=proposal`,
  вне дерева); попадает в роадмап только через явный `accept` овнером (accept-gate).

## Эндпоинты

Все под `/knowledger/` (ниже — без префикса). Все, кроме `healthz`, требуют `Bearer`.

### Workspaces
| Метод · путь | Payload / примечание |
|---|---|
| `GET /workspaces` | список всех |
| `POST /workspaces` | `{"key":"KNOW","name":"Knowledger KB","description":"…"}` → `201` |
| `GET /workspaces/{ws}` | по key |
| `PATCH /workspaces/{ws}` | частичный: `{"name":"…"}` и/или `{"description":"…"}` |

### Nodes
| Метод · путь | Payload / примечание |
|---|---|
| `GET /workspaces/{ws}/nodes` | плоский список; `?parent=<key\|uuid>` — дети узла; `?parent=none` — только корни |
| `POST /workspaces/{ws}/nodes` | `{"title":"…","body":"…","kind":"section"?,"parent_id":"<key\|uuid>"?}` → `201`. `title` обязателен; без `parent_id` — корень |
| `GET /nodes/{key}` | узел по key ИЛИ UUID |
| `PATCH /nodes/{key}` | частичный: `title`,`body` (строки); `kind`,`parent_id` — строка ИЛИ `null` (очистка/reparent). Reparent валидирует цикл и cross-ws |
| `DELETE /nodes/{key}` | `204`. С детьми → `409` (сначала перевесить/удалить детей) |
| `GET /nodes/{key}/children` | прямые дети |
| `GET /nodes/{key}/tree` | узел + всё поддерево; `?depth=N` — лимит глубины |
| `GET /workspaces/{ws}/tree` | корни ws + их поддеревья; `?depth=N` |

### Refs / Backlinks
| Метод · путь | Payload / примечание |
|---|---|
| `GET /nodes/{key}/refs` | исходящие рёбра узла |
| `POST /nodes/{key}/refs` | `{"to_node":"<key\|uuid>","kind":"link"}` → `201`. `kind` ∈ link·transclude·relates·blocks·depends_on·duplicate; чужой ws в `to_node` ок |
| `GET /nodes/{key}/backlinks` | входящие рёбра (derived) |
| `DELETE /refs/{id}` | `204` (id ребра — из объекта Ref) |

### Tags
| Метод · путь | Payload / примечание |
|---|---|
| `GET /workspaces/{ws}/tags` | теги ws |
| `POST /workspaces/{ws}/tags` | `{"name":"canon"}` → `201` |
| `GET /nodes/{key}/tags` | теги узла |
| `POST /nodes/{key}/tags` | `{"tag_id":"<uuid>"}` навесить → `204` |
| `DELETE /nodes/{key}/tags/{tag_id}` | снять → `204` |

### Activity
| Метод · путь | Payload / примечание |
|---|---|
| `GET /nodes/{key}/activity` | timeline узла |
| `POST /nodes/{key}/activity` | `{"kind":"commented","data":{"text":"…"}}` → `201`. actor = handle из Bearer |

### Cross-product proposals (accept-gate)
| Метод · путь | Payload / примечание |
|---|---|
| `POST /workspaces/{ws}/proposals` | предложка в ЧУЖОЙ ws: `{"title":"…","body":"…","source_ws":"DEVOPSER"}` → `201` (`origin=proposal`, в inbox, НЕ в дереве). actor → `proposed_by` |
| `GET /workspaces/{ws}/inbox` | входящие предложки (`origin=proposal`), отдельно от дерева |
| `POST /nodes/{key}/accept` | приземлить в роадмап: тело опц. `{"parent_id":"<key\|uuid>"}` (без — корень). `origin→native`, key сохраняется. Не-предложка → `400` |
| `POST /nodes/{key}/decline` | отклонить: тело опц. `{"comment":"…"}`. `origin→declined`, остаётся в истории inbox |

### Health
`GET /knowledger/healthz` → `200 {"status":"ok","service":"knowledger"}` (без auth).

## Формы ответов

**Node** (и элемент дерева — плюс `children[]`):
```json
{
  "id":"uuid","workspace_id":"uuid","key":"KNOW-12","seq":12,
  "parent_id":"uuid|null","kind":"section|null","title":"…","body":"…",
  "ord":"a0","created_at":"…","updated_at":"…",
  "origin":"native|proposal|declined",
  "proposed_by":"devopser?","source_ws":"DEVOPSER?"
}
```
**Ref:** `{"id","from_node","to_node","kind","created_at"}` ·
**Tag:** `{"id","workspace_id","name","created_at"}` ·
**Activity:** `{"id","node_id","actor","kind","data?","created_at"}` ·
**Workspace:** `{"id","key","name","description","created_at","updated_at"}`.

## Ошибки

Тело всегда `{"error":"<msg>"}`. Статусы: `400` — валидация (пустой title, кривой ws-key,
неизвестный ref-kind, цикл reparent, accept не-предложки), `401` — нет/пустой Bearer,
`404` — узел/ws не найден, `409` — конфликт (дубль ws-key, delete узла с детьми),
`500` — внутренняя. Неизвестные поля в теле POST/PATCH отвергаются (`400`) — ловим опечатки.

## Быстрый старт (наполнение с нуля — KB сейчас пустая)

```bash
BASE=http://knowledger:8040/knowledger        # или http://gateway/api/knowledger через дверь
AUTH='Authorization: Bearer devopser'

# 1) создать workspace
curl -s -X POST -H "$AUTH" -H 'Content-Type: application/json' $BASE/workspaces \
  -d '{"key":"DEVOPSER","name":"Devopser KB","description":"База знаний девопсера"}'

# 2) корневой документ
curl -s -X POST -H "$AUTH" -H 'Content-Type: application/json' $BASE/workspaces/DEVOPSER/nodes \
  -d '{"title":"Дверь хаба","body":"# Контракт двери\n…","kind":"doc"}'   # → key DEVOPSER-1

# 3) секция-ребёнок
curl -s -X POST -H "$AUTH" -H 'Content-Type: application/json' $BASE/workspaces/DEVOPSER/nodes \
  -d '{"title":"Маршруты","body":"…","kind":"section","parent_id":"DEVOPSER-1"}'

# 4) ссылка между узлами
curl -s -X POST -H "$AUTH" -H 'Content-Type: application/json' $BASE/nodes/DEVOPSER-2/refs \
  -d '{"to_node":"DEVOPSER-1","kind":"relates"}'

# 5) прочитать всё дерево ws
curl -s -H "$AUTH" $BASE/workspaces/DEVOPSER/tree
```
