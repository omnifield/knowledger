# ROADMAP — сводная (cross-repo картина)

Один источник правды «где всё сейчас» (`standards/workflow/roadmap-template.md`).
Обновляется при каждом сдвиге. Вопрос «что сейчас происходит и кто что делает» —
ответ ЗДЕСЬ, не в переписке.

## Таблица

| Feature | Stage | Maturity | Owner-repo | Blocked-by | DoR | DoD | Link |
|---|---|---|---|---|---|---|---|
| brainer interface-MVP (фронт+бэк против контракта) | App | done | brainer | — | ✅ | ✅ | brainer PR#1 |
| brainer repo-skeleton (pnpm+nx+uv+CI, эталон продукт-репо) | DevOps | ready | brainer | — | ✅ | ⬜ | brainer `briefs/repo-skeleton.md` |
| commons foundation-first + toolchain-pins (правила) | Contract | done | commons | — | ✅ | ✅ | `standards/canon/principles/foundation-first.md` |
| devopser workstation-bootstrap v1 | DevOps | done | devopser | — | ✅ | ✅ | devopser `briefs/workstation-bootstrap.md` |
| devopser workstation амендмент №2 (corepack out, pnpm в слой) | DevOps | ready | devopser | — | ✅ | ⬜ | тот же бриф, §Амендмент №2 |
| devopser infra-migration (стеки из оракула + registry) | DevOps | ready | devopser | — | ✅ | ⬜ | devopser `briefs/infra-migration.md` |

## Очередь запусков (кто, что, после кого)

Правило чтения: бери верхний item со своей ролью, у которого `Blocked-by` пуст.

| # | Кто запускается | Репо / scope | Задача | После кого |
|---|---|---|---|---|
| 1 | architect brainer (`-Scope main`) | brainer | закоммитить оба брифа → выполнить repo-skeleton по шагам | ни после кого — можно сейчас |
| 2 | owner-workstation (`-Scope workstation`) | devopser | переделать bootstrap по амендменту №2 | ни после кого — параллельно с №1 |
| 3 | architect devopser (`-Scope main`) | devopser | infra-migration (стеки observability/gateway/storage + registry) | после №2 (чтобы workstation-канон не менялся под ногами) |

User в этой схеме: запускает сессии по очереди выше + push/merge там, где просит architect.
