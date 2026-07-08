# Go — канон сервиса

Верхний закон — `../../POLICY.md`. Родственные: [types-from-zod](../principles/types-from-zod.md)
(здесь тот же закон в форме «типы из схемы, не руками»), [etalon-gate](../principles/etalon-gate.md),
[ecosystem-self-sufficiency](../principles/ecosystem-self-sufficiency.md) (вендор за границей
адаптера), `../../workflow/toolchain-pins.md` (пины), `../../workflow/testing.md`.

Первый потребитель — `chater` (brainer, extract-candidate). Канон для **всех** наших
Go-сервисов; брифы архитекторов на язык не опираются — язык-детали закрывает этот файл + owner.

## Тулчейн (пины — по паттерну toolchain-pins)

- **Версия Go пинится в `go.mod`**: директива `go X.Y.Z` (точная, не мажор) + `toolchain`.
  Go ≥1.21 сам качает запиненный тулчейн — системный Go любой свежести, ручных установок
  версий нет. Базовый слой машины получает `go` через devopser bootstrap.
- **Форматирование не обсуждается**: `gofmt` (через `gofumpt` — строже, без конфига).
  CI гейтит diff.
- **Линт**: `golangci-lint`, конфиг `.golangci.yml` в корне сервиса. Baseline enable:
  `errcheck, govet, staticcheck, unused, ineffassign, misspell, gocritic, revive, gofumpt`.
  Версия golangci-lint пинится в CI и в `tools`-комментарии сервиса.

## Раскладка сервиса

```
<service>/
  cmd/<service>/main.go   # только wiring: config → deps → serve. НОЛЬ логики
  internal/               # ВСЯ реализация (компилятор сам фенсит внешний импорт)
    <domain>/             # пакеты по доменам (rooms, messages, participants)
    httpapi/              # транспорт: роуты/хендлеры (тонкие, без логики)
    store/                # доступ к данным (sqlc-генерация + обвязка)
  migrations/             # SQL-миграции (goose), нумерованные
  api/                    # JSON-схемы/контракт наружу (если сервис их владелец)
  go.mod                  # свой модуль = сам себе граница (extract = git mv)
```

- **Никакого `pkg/`** — публичного Go-API у сервисов нет; наружу только HTTP/wire-контракт.
- Свой `go.mod` на сервис (не общий workspace-модуль) — граница extract'а.

## Стиль (минимум правил, остальное несёт линт)

1. **Ошибки — значения**: возврат `error` последним; оборачивание `fmt.Errorf("...: %w", err)`;
   sentinel/typed-ошибки в пакете-владельце. `panic` — только в `main` на fatal-wiring.
   Глотать ошибку (`_ =`) без комментария-причины — нарушение (root-cause канон).
2. **`context.Context` — первый аргумент** всего, что ходит в IO; не хранить в структурах.
3. **Интерфейс объявляет ПОТРЕБИТЕЛЬ** (accept interfaces, return structs), там где
   потребляется, минимальный по методам. Это Go-форма нашего «контракт через interface,
   не через реализацию соседа».
4. **Ноль глобального мутабельного состояния**; зависимости — явным конструктором
   `New<X>(deps)`.
5. Конкурентность: горутина не запускается без ответа «кто её останавливает» (context/
   errgroup); каналы — деталь реализации пакета, наружу не торчат.

## Транспорт

- HTTP: **stdlib `net/http`** (ServeMux ≥1.22 умеет method+path). Фреймворк не берём;
  упёрлись в реальный gap — сначала эскалация (prefer-existing-libs, но stdlib-first).
- WebSocket: `coder/websocket` (maintained, context-first). SSE — руками поверх stdlib
  (это 20 строк, библиотека не нужна).
- Graceful shutdown обязателен (`signal.NotifyContext` + `srv.Shutdown`).

## Данные

- **SQLite → Postgres drop-in** (паттерн `backend/learn`): SQL пишем совместимым, драйвер
  за build-tag/конфигом не прячем — выбор БД = конфиг.
- **Миграции: `goose`** (SQL-файлы в `migrations/`, идемпотентный `up` на старте dev).
- **Запросы: `sqlc`** — типы и методы генерятся из SQL-схемы. Ручной маппинг
  row→struct — нарушение (types-from-schema).

## Контракты (wire)

Типы wire-формата (события kernel-контракта, HTTP-DTO чужих контрактов) **генерятся из
JSON-схемы владельца контракта**, руками не пишутся. Свой контракт сервис описывает
JSON-схемой в `api/` — из неё генерятся и его Go-типы, и TS-типы потребителей. Один
источник — ноль дрейфа.

## Логи · конфиг · наблюдаемость

- Логи: **stdlib `log/slog`**, structured, JSON-handler в проде; уровни через конфиг.
- Конфиг: **env-only** (12-factor, как python-сервисы); маленькая `Config`-структура с
  явным парсом в `cmd/`; никаких конфиг-файлов-форматов без решения архитектора.
- Трейсы эталон-гейта: замеры горячих путей (spawn/fanout/query) логгером — observability
  вшита, см. etalon-gate.

## Тесты

- **stdlib `testing`, table-driven**; `testify` допустим для assert/require, моки руками
  (интерфейсы и так минимальные — см. Стиль §3).
- **`go test -race` в CI всегда** (fanout-домены без race-детектора не живут).
- Интеграционные тесты store — на реальном SQLite (файл во временной директории), не моках.

## CI-гейт сервиса

`gofumpt -l` (пусто) → `go vet` → `golangci-lint run` → `go test -race ./...` → `go build ./...`.
Всё читает пины из репо, версии в workflow не дублируются (toolchain-pins).
