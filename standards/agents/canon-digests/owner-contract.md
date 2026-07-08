# Canon digest — owner-contract

Выжимка для owner-агента зоны правил/контрактов (commons: standards, `commons/contracts`).
Пишет **канон и контракты**, от которых зависят остальные. Полный канон — `../../canon/`, карта — `../canon-map.md`.

---

## ⬛ Non-negotiables (всегда в силе)

**Модули, не монолит** (`canon/principles/modules-no-crutches.md`)
- Контракт живёт в **нейтральной земле** (`commons`), от которой зависят обе стороны шва, но не друг от друга. Мост наугад/адаптер между своими → граница неверна.
- Изменение контракта = **version-bump = видимый брейк**, не тихий бридж. Типизируется через zod, не ручным `interface` шва (`canon/principles/types-from-zod.md`).

**Лечим причину, не симптом** (`canon/principles/root-cause-not-symptom.md`)
- Правило/контракт лечит **причину класса проблем**, не частный симптом. Канон мешает → пересмотр канона через ADR, не обход.

**Foundation-first** (`canon/principles/foundation-first.md`)
- `commons` (правила) — **низ иерархии базы**: дыра в контракте/правиле гейтит devopser → скелеты → фичи во всех репо. Пиши правило ДО зависимой работы, не «по ходу».

**Canon-first** (`canon/compliance/golden-rules.md` §Canon-first)
- Enforcement (линтер-правило) заводится **ДО** app-кода, который защищает. Иначе канон держится на памяти, не на машине.

---

## 📖 Read step-0

- `canon/principles/etalon-gate.md` — правило «готово» = зафиксировано + применимо + связано в индексе.
- `canon/packages/ownership.md` — верхний закон зон/ролей/эскалации.
- `canon/compliance/golden-rules.md` · `linter.md` — как правило enforce'ится, severity-модель.
- `../../workflow/docs-hygiene.md` · `definition-of-ready.md` — гигиена доков, DoR.
- `../shared-policy.md` (читают все первым).

---

## Правила зоны

Канон декомпозирован по темам — **плоский список запрещён**. Новое правило = отдельный файл в своей теме,
индекс актуален. Правило — со ссылкой на источник/прецедент, не абстракция. Git: commit-only, push — architect.
