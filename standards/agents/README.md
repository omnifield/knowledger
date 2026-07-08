# Agents — флоу делегирования (шаблоны)

Репо-агностик каноны agent-системы. Каждый репо **наследует** эти шаблоны в свой `.claude/agents/` и допиливает per-repo механику (hooks/git-gate, scope-identity, пины модели). Все агенты первым читают [shared-policy](shared-policy.md).

## Три яруса

```
        product owner (user)
                │  запрос
                ▼
     ┌──────────────────────┐
     │  ARCHITECT (main)     │  триаж · арх-решения · координация · git
     └──────────┬───────────┘
        брифы / Agent-спавн
                ▼
     ┌──────────────────────┐
     │  OWNER-<zone>         │  код зоны · тесты · доки · commit-only
     └──────────┬───────────┘
        (в аппе) делегирует
                ▼
     ┌──────────────────────┐
     │  LAYER-агент          │  один артефакт слоя · без git
     └──────────────────────┘
        эскалация — строго ВВЕРХ
```

| Ярус | Роль | Модель | Git | Шаблон |
|---|---|---|---|---|
| **Architect** | триаж/архитектура/координация | opus | полный | [architect.template](architect.template.md) |
| **Owner** | full lifecycle зоны | sonnet | commit-only (под gate) | [owner.template](owner.template.md) |
| **Layer** | один HCA-артефакт | haiku (ctrl/feat→sonnet) | нет | [layer.template](layer.template.md) |

## Как вызываются

- **Architect → owner**: бриф-файл (`docs/_meta/briefs/*.md`), owner-сессию запускает product owner (scope-identity) ЛИБО architect спавнит субагента `Agent(subagent_type='owner-<zone>')`.
- **Architect → layer** (или owner-app → layer): `Agent(subagent_type='<layer>', prompt='...')`; несколько параллельно в одном сообщении, когда задачи независимы.
- **Owner → owner** (сосед по release-группе, тривиальный fix): `Agent(subagent_type='owner-<X>')` по политике boundaries ([shared-policy](shared-policy.md) §1).

## Что НЕ делегируется субагентам

| Задача | Кто |
|---|---|
| Арх-решение, ADR, контракт | architect |
| Cross-слойный / cross-пакетный refactor | architect оркеструет по частям |
| Code-review результата субагента | architect |
| Чистая косметика (1 файл, rename) | architect напрямую |

## Реестр зон (per-repo)

Каждый репо ведёт свою **ownership-matrix** (зона → owner → release-группа) в своей доке (`docs/_meta/agents.md` или аналог). Здесь — только форма ролей, не список зон конкретного репо.

## Принципы дизайна

1. **Канон-срез per-роль** — гибрид-инъекция: железные «нельзя» (`⬛`) инлайнятся в промпт, глубина (`📖`) — ссылкой. Карта роль→срез: [canon-map](canon-map.md); готовые выжимки: [canon-digests/](canon-digests/). Агент видит не 18 файлов, а свой digest.
2. **Минимум tools**: layer = `Read/Write/Edit/Glob`; owner = `+Bash`; никогда `Grep`/`Agent` у листьев.
3. **Один артефакт на layer-вызов**; owner — пакетная правка в своей зоне.
4. **Compliance-правила в каждом prompt** — агент знает свои «нельзя».
5. **Один уточняющий вопрос максимум.**
6. **Эволюция шаблона** — правится .md агента, не код: «обновился один файл — все будущие генерации правильные».

## Связанное

- [shared-policy](shared-policy.md) — читают все агенты первым.
- [git-scope](../workflow/git-scope.md) · [commit-cadence](../workflow/commit-cadence.md) — scope и каденс коммитов.
- [ownership](../canon/packages/ownership.md) — зоны/роли/эскалация (верхний закон).
