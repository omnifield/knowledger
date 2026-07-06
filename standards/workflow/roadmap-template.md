# Roadmap — живой статус (шаблон)

Один источник правды «где всё сейчас». Копируется в каждый репо как `ROADMAP.md` (свои фичи) + сводная версия в `commons` (cross-repo картина). **Никакой дата-археологии** — статус читается из таблицы, а не «по датам догадайся».

## Как вести

- Одна строка = одна фича/item.
- Обновляй `Stage` / `Maturity` / `Blocked-by` при каждом сдвиге — в паузах, перед рестартом сессии.
- `Blocked-by` пуст → можно двигать. Заполнен → ждёт апстрим (см. пайплайн в `definition-of-ready.md`).
- Закрыл item → `Maturity: done/etalon`, не удаляй строку (история стадий видна).

## Таблица

| Feature | Stage | Maturity | Owner-repo | Blocked-by | DoR | DoD | Link |
|---|---|---|---|---|---|---|---|
| _Learn lessons (пример)_ | Domain | doing | framework | — | ✅ | ⬜ | ADR-069 / PR#… |
| _Moderator MVP (пример)_ | Contract | proposed | commons | — | ⬜ | ⬜ | ADR-077 |

**Stage:** `Contract` · `Backend` · `Domain` · `App` · `DevOps`
**Maturity:** `proposed` · `ready` · `doing` · `done` · `etalon` · `parked`
**DoR/DoD:** ✅ met · ⬜ not yet (детали чек-листов — `definition-of-ready.md`)

## Сводная (только в commons)

Дополнительно — раздел «Cross-repo фичи в полёте»: фичи, чьи стадии живут в разных репо, с текущей стадией и кто следующий разблокируется. Чтобы cross-repo цепочка не терялась между репо-роадмапами.
