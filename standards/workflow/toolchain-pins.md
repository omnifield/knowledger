# Toolchain pins — репо декларирует свой тулчейн (машина = cattle)

Паттерн для **всех** репо экосистемы. Происхождение: инцидент 2026-07-08 (новая машина без
uv/python → ручная установка + пути) и закрывшая его эскалация brainer
`briefs/escalation-toolchain-pins.md`. Никакая версия инструмента не живёт «в голове» или
«в истории терминала» — только пином в репо.

## Пины (полный набор)

| Что пинит | Где | Кто исполняет |
|---|---|---|
| Python | `.python-version` в корне uv-workspace | uv сам качает managed CPython. **Системный Python не ставим никогда.** |
| сам uv | root `pyproject.toml` → `[tool.uv] required-version = ">=X.Y,<X.Z"` | uv отказывается работать не той версией (self-enforcing). Без пина — дрейф резолвера/lock-формата |
| pnpm | `packageManager: "pnpm@x.y.z"` в root `package.json` | **сам pnpm ≥10** (`manage-package-manager-versions`, дефолт из коробки): любой pnpm 10.x скачивает и запускает запиненную версию |
| node | `engines.node` + root `.npmrc` → `engine-strict=true` | pnpm. ⚠️ Без `engine-strict` engines — **warning, не гейт** |
| Go | `go.mod` → директивы `go X.Y.Z` (точная) + `toolchain` | go ≥1.21 сам качает запиненный тулчейн; системный go — только базовый слой (bootstrap devopser). Канон сервиса — `../canon/languages/go.md` |

## ⚠️ Corepack — НЕ опора

Официально deprecated, выпиливается из будущих мажоров Node; требует ручного `corepack enable`
(противоречит «ставится само»). Пин `packageManager` исполняет сам pnpm ≥10 (см. таблицу).
`corepack enable` не выполняем нигде — ни в bootstrap'ах, ни в CI, ни в доках.

## Базовый слой машины (единственные prerequisites)

**git · node LTS · pnpm ≥10 · uv · Docker · claude CLI** — ставит devopser
`workstation/bootstrap.ps1` одной идемпотентной командой. Всё остальное самособирается
из пинов (`uv sync` / `pnpm install`). Правило: поставил на машину что-то руками →
это gap workstation-bootstrap'а, фиксится в devopser.

## CI читает те же пины

`pnpm/action-setup` **без** явной версии (читает `packageManager`), `astral-sh/setup-uv`
(читает `.python-version`). Версии инструментов в workflow-файлах не дублировать — один
источник правды, локалка и CI питаются из него.
