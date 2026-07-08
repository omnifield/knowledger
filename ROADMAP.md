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
| 🕳 brainer fix: config.py реестр от `parents[4]` (path-костыль; тесты зависят от имени папки чекаута) | Backend | ready | brainer | — (дыро-фикс не гейтится) | ✅ | ⬜ | brainer PR#2 §заметки |
| 🕳 brainer fix: vite port 5173 → 3500 (контракт DEPLOY.md + devopser `registry/ports.md`) | App | ready | brainer | — (дыро-фикс не гейтится) | ✅ | ⬜ | brainer PR#2 §заметки |

## Очередь запусков (кто, что, после кого)

Правило чтения: бери верхний item со своей ролью, у которого `Blocked-by` пуст.
Иерархия базы (foundation-first): дыры devopser гейтят ВСЁ выше — brainer стоит,
пока workstation DoD не зелёный.

| # | Кто запускается | Репо / scope | Задача | После кого |
|---|---|---|---|---|
| 1 | user (или architect brainer) | brainer | смержить PR#2 (CI зелёный, отклонения одобрены) | сейчас |
| 2 | **user + owner-workstation** | devopser | e2e-прогон bootstrap на чистой машине (Windows Sandbox / VM / следующая свежая тачка) → закрыть дыры 1–2 + «≥10» в README/ps1 | сейчас — это ГЛАВНЫЙ гейт |
| 3 | owner-backend (`-Scope backend`) | brainer | 🕳 config.py: реестр репо конфигом/env, не `parents[4]` | после мержа №1 (дыро-фикс, e2e не ждёт) |
| 4 | owner-frontend (`-Scope frontend`) | brainer | 🕳 vite port → 3500 по контракту | после мержа №1 (дыро-фикс, e2e не ждёт) |
| 5 | architect devopser (`-Scope main`) | devopser | infra-migration (стеки + registry) | параллельно №2 (тот же уровень базы) |
| 6 | **фичевая разработка (любая)** | все | — | ТОЛЬКО после №2–4 (все дыры закрыты) |

User в этой схеме: запускает сессии по очереди выше + push/merge там, где просит architect.
