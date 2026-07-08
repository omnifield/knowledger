# Canon-map — какой роли какой срез канона

Стоппер: канон = 18 файлов в 5 темах, агент читает не всё → ошибки. Фикс — **срез per-роль**:
не «прочитай весь канон», а «железное — в промпте, глубина — по короткому must-read списку».

## Архетипы ролей

| Архетип | Кто (пример по репо) | Домен |
|---|---|---|
| **architect** | main-сессия любого репо | триаж / арх-решения / координация / git |
| **owner-frontend** | owner-web-*, owner HCA-зоны | UI/HCA (framework, apps) |
| **owner-backend** | owner-backend (brainer), owner devopser-стеков | python / rust сервисы |
| **owner-contract** | owner commons/standards | правила экосистемы, контракты |
| **layer-HCA** | view/shape/widget/page/controller/feature/entity/ui-component | один HCA-артефакт |

## Матрица

`⬛` = inline в промпт (non-negotiable, агент физически видит) · `📖` = must-read step-0 (ссылка) · `·` = skip

| canon | architect | owner-front | owner-back | owner-contract | layer-HCA |
|---|:--:|:--:|:--:|:--:|:--:|
| `principles/modules-no-crutches` | ⬛ | ⬛ | ⬛ | ⬛ | ⬛ |
| `principles/root-cause-not-symptom` | ⬛ | ⬛ | ⬛ | ⬛ | ⬛ |
| `principles/foundation-first` | ⬛ | 📖 | 📖 | ⬛ | · |
| `principles/etalon-gate` | ⬛ | 📖 | 📖 | 📖 | 📖 |
| `principles/types-from-zod` | 📖 | ⬛ | 📖¹ | 📖 | ⬛ |
| `architecture/layers` | 📖 | ⬛ | · | · | ⬛ |
| `architecture/import-rules` | 📖 | ⬛ | · | · | ⬛ |
| `architecture/ui-proxy-tag-flow` | · | 📖 | · | · | 📖 |
| `architecture/namespaces` | · | 📖 | · | · | 📖 |
| `packages/anatomy` | 📖 | 📖 | 📖 | · | · |
| `packages/ownership` | ⬛ | ⬛ | ⬛ | 📖 | · |
| `packages/dependency-tiers` | 📖 | 📖 | 📖 | · | · |
| `components/component-model` | 📖 | ⬛ | · | · | 📖 |
| `components/kit-first` | · | ⬛ | · | · | ⬛ |
| `components/registration` | · | 📖 | · | · | 📖 |
| `components/tokens` | · | 📖 | · | · | 📖 |
| `compliance/golden-rules` | 📖 | ⬛ | 📖 | 📖 | ⬛ |
| `compliance/linter` | 📖 | 📖 | · | · | · |

¹ `types-from-zod` для backend = типы из схемы (pydantic/SQLAlchemy как source of truth), тот же принцип «ноль ручных типов под домен», другой инструмент.

## Гибрид-инъекция (как срез попадает в агента)

1. **Digest на архетип** — `canon-digests/<архетип>.md`: `⬛`-правила **выжимкой** (1 экран) + список `📖` must-read.
2. **Репо наследует** — `.claude/agents/<agent>.md` каждого репо **инлайнит** блок своего архетипа из digest'а (через include при сборке промпта или копией с пометкой источника) + `📖`-ссылки как «step-0 read».
3. **Обновление канона** → правим canon-файл → обновляем digest → агенты наследуют новое. Digest **не дублирует** канон дословно — это выжимка «нельзя», указывающая на источник.

**Почему гибрид:** железное всегда в контексте (нельзя не увидеть — решает стоппер), детали не раздувают промпт (по ссылке on-demand). Только-ссылка = агент сливается (текущая боль); только-инлайн всего = промпт-жир и дубль при апдейте.

## Правило ведения

- Новый canon-файл → добавь строку в матрицу + реши, кому `⬛`/`📖`.
- Новый архетип → колонка + свой digest.
- Digest меняется только вслед за canon-файлом, не самостоятельно.
