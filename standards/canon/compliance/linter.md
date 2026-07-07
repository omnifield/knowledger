# Линтер — как enforce'ится

Правила ([golden-rules](golden-rules.md)) держит машина, не память. Три точки контакта.

## Механика

- **AST-линтер** (`@omnifield/compliance`) — читает исходник как AST, проверяет импорты/JSX/классы против правил слоя. Знает раскладку слоёв ([layers](../architecture/layers.md)) и зоны.
- **Vite-плагин** (`CompliancePlugin`) — гоняет линтер в dev-сборке. Structural-нарушения выводит как `[STRUCTURAL ERROR — CI will fail]`, **не блокируя HMR** (разработка не встаёт, но ошибка видна сразу).
- **CI-gate** (`compliance:check`) — отдельный job пайплайна; structural-нарушение = красный CI. Cosmetic (`warn`) — не валит.

## Severity в двух режимах

Один и тот же structural-violation: в dev — предупреждение с пометкой «CI упадёт», в CI — падение. Это даёт быструю обратную связь локально и жёсткий гейт на въезде в main. Cosmetic логируется с `file:line:column` + hint везде.

## Allowlist

Transitional-allowlist ([golden-rules](golden-rules.md)) читается линтером: path+kind суппрессится на время cleanup с `reason`+`ttl`. Не для постоянного обхода — для миграции зоны к эталону.

## Список portable-tier / disallowed-deps

Явные списки (portable-пакеты, запрещённые внешние deps per слой) — в конфиге линтера, не размазаны по коду. Меняются осознанно (owner-builders / architect), видны в одном месте.

## Принцип

Канон, не покрытый линтером, = канон на честном слове. Новое правило канона → правило линтера (canon-first). Правило меняется → правится линтер, все будущие проверки следуют новой версии автоматически.
