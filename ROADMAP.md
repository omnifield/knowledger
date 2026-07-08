# ROADMAP — сводная (cross-repo картина)

Один источник правды «где всё сейчас» (`standards/workflow/roadmap-template.md`).
Обновляется при каждом сдвиге. Вопрос «что сейчас происходит и кто что делает» —
ответ ЗДЕСЬ, не в переписке.

## Таблица

| Feature | Stage | Maturity | Owner-repo | Blocked-by | DoR | DoD | Link |
|---|---|---|---|---|---|---|---|
| brainer interface-MVP (фронт+бэк против контракта) | App | done | brainer | — | ✅ | ✅ | brainer PR#1 |
| brainer repo-skeleton (pnpm+nx+uv+CI, эталон продукт-репо) | DevOps | done | brainer | мерж PR#2 (CI зелёный) | ✅ | ✅ | brainer PR#2 |
| commons foundation-first + toolchain-pins (правила) | Contract | done | commons | — | ✅ | ✅ | `standards/canon/principles/foundation-first.md` |
| devopser workstation-bootstrap (v1 + амендмент №2 реализованы) | DevOps | doing | devopser | e2e на чистой машине (дыры 1–2 эскалации) | ✅ | ⬜ | devopser `workstation/escalation-bootstrap-gaps.md` |
| devopser infra-migration (стеки из оракула + registry) | DevOps | ready | devopser | — | ✅ | ⬜ | devopser `briefs/infra-migration.md` |
| 🕳 brainer: разгрести shared tree (устаревшая ветка + смешанный WIP двух owner'ов) | DevOps | done | brainer | — | ✅ | ✅ | main чист/синхронен, WIP спасён в `fix/`-ветки (обе на origin) |
| 🕳 brainer fix: config.py реестр от `parents[4]` (path-костыль; тесты зависят от имени папки чекаута) | Backend | ready | brainer | очередь WIP=1 (№3) | ✅ | ⬜ | brainer `briefs/fix-backend-repo-registry.md` |
| 🕳 brainer fix: vite port 5173 → 3500 (контракт DEPLOY.md + devopser `registry/ports.md`) | App | ready | brainer | очередь WIP=1 (№4) | ✅ | ⬜ | brainer `briefs/fix-frontend-port-3500.md` |
| 🕳 process: параллельные owner'ы на одном working tree без изоляции | Contract | done | commons | — | ✅ | ✅ | shared-policy §4/4.1 + foundation-first WIP=1 |
| testing-стандарт (6 принципов; боли v1 зафиксированы) | Contract | proposed | commons | обсуждение user+architect (развилки: win-CI, playwright) | ⬜ | ⬜ | `standards/workflow/testing.md` 🚧 DRAFT |

## Очередь запусков (кто, что, после кого)

Правило чтения: бери верхний item со своей ролью, у которого `Blocked-by` пуст.
Иерархия базы (foundation-first): дыры devopser гейтят ВСЁ выше — brainer стоит,
пока workstation DoD не зелёный.

**Очередь СТРОГО последовательная (WIP=1, foundation-first): следующий стартует только
когда предыдущий закрыт. Никаких параллельных сессий ни на одном дереве.**

| # | Кто запускается | Репо / scope | Задача |
|---|---|---|---|
| 1 | architect brainer (`-Scope main`) | brainer | разгрести дерево: спасти WIP обоих owner'ов (config.py+tests, vite.config.ts) в отдельные ветки/стэши, переключить на актуальный main, удалить смерженную ветку |
| 2 | **user + owner-workstation** | devopser | e2e-прогон bootstrap на чистой среде (Sandbox / VM / свежая тачка) → закрыть дыры 1–2 + «≥10» в README/ps1 |
| 3 | owner-backend (`-Scope backend`) | brainer | 🕳 config.py по брифу (один, на чистом main) |
| 4 | owner-frontend (`-Scope frontend`) | brainer | 🕳 vite port → 3500 по брифу |
| 5 | architect devopser (`-Scope main`) | devopser | infra-migration (стеки + registry) |
| 6 | **фичевая разработка** | все | только когда таблица выше без открытых 🕳; параллельность — по явной записи architect'а |

User в этой схеме: запускает сессии по очереди выше + push/merge там, где просит architect.
