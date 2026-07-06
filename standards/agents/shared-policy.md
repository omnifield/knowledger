# Agent shared-policy — читают ВСЕ агенты первым

Cross-cutting правила для любого агента (architect / owner / layer) любого репо. Верхний закон — `../POLICY.md`; детали канона — `../canon/`; дисциплина — `../workflow/`. Здесь — то, что иначе дублировалось бы в каждом agent-конфиге.

## 1. Boundaries

- У каждой зоны (пакет/сервис) — **owner**. Не лезь в чужую зону своими руками.
- Задача требует правок в чужой зоне:
  - **Тривиально** (typo, missing export, stale ref, bump до next minor) — запроси fix через `Agent(subagent_type='owner-<zone>')`, опиши конкретно.
  - **Нетривиально** (новый API, refactor, breaking, дизайн) — НЕ ЛЕЗЬ. Эскалируй: «для X нужно Y в зоне Z; делегировать owner-Z». Решает architect репо.
- Cross-repo концерн → **контракт** в `commons/contracts`, не переброс брифами ([[../canon/packages/ownership]]).

## 2. Доки (часть DoD)

- После изменения — синхронизируй: AI-anchor (`docs/_meta/<zone>.md`) + user-guide (где применимо) + per-zone README.
- Нет доки — создай; протухла — почини; старое помечай `superseded` ([[../workflow/docs-hygiene]]).
- Канон **не дублируется** в prompt: owner ссылается на AI-anchor как single source of truth, обновляет anchor — не копирует контент.

## 3. Тесты + трейсы (DoD)

- **Definition of done** = код + тесты + трейсы + доки + раскладка ([[../canon/principles/etalon-gate]]).
- Pure-логика — unit (vitest, node); DOM/Solid — jsdom. Баг → сначала характеризационный тест (repro), потом fix.
- Не тестируемо изолированно (только реальный браузер) — задокументируй почему в тесте-плейсхолдере; закрытие через browser-eyeball product owner'а.
- Трейсы = инструментирование perf-логгерами ([[../canon/principles/etalon-gate]]), не «допишем позже».

## 4. Commit-каденс + git-роли

- Работа **по этапам**: этап → проверка → коммит; не коммитим непроверенное ([[../workflow/commit-cadence]]).
- **Pre-commit гейт** — test+lint+build зелёные ПЕРЕД коммитом.
- **Git-роли** (per-repo механика — hooks/marker, но канон общий):
  - **Architect (main-сессия)** — полный git (commit/push/merge).
  - **Owner-субагент** — **commit-only** (под git-gate); push/merge делает architect после ревью.
  - **Layer-агент** — git **не трогает** вообще (пишет артефакт, возвращает управление).
- Хук заблокировал операцию — **НЕ обходи** (`bash -c`, `&&`, `--no-verify`). STOP + эскалация.

## 5. Cross-zone контекст (owner)

- Owner знает свою **release-группу** и **consumer'ов** своей зоны. Изменение публичного API без согласования с consumer'ом — нет.
- Сомнение — спроси architect'а, не делай sweeping refactor «на удачу» ([[../canon/principles/root-cause-not-symptom]]).

## 6. Стиль

- Комментарии — `// Why:` для неочевидного, не дублируй код.
- Один уточняющий вопрос максимум — лучше написать и переписать, чем спросить пять раз.
