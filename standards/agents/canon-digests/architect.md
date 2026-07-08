# Canon digest — architect

Выжимка для architect (main-сессия). Полный канон — `../../canon/`, карта — `../canon-map.md`.
Architect — **верх эскалации**: знает канон широко (границы + куда что), детали — по ссылкам.

---

## ⬛ Non-negotiables (всегда в силе)

**Модули, не монолит** (`canon/principles/modules-no-crutches.md`)
- Экосистема независимых модулей. Cross-зонная связь — контракт в нейтральной земле (`commons`),
  не прямой импорт реализации. Нужен адаптер/шим между своими → граница неверна, пересматриваем (ADR).
- Изменение контракта = version-bump = видимый брейк, не тихий бридж.

**Лечим причину, не симптом** (`canon/principles/root-cause-not-symptom.md`)
- Костыль во фреймворке = compound debt. Канон мешает → это сигнал пересмотра канона (**ADR**), не обход.
- Диагностика фактами, не гипотезами. 2 неудачных фикса → СТОП + корень.

**Foundation-first** (`canon/principles/foundation-first.md`)
- Известные дыры базы гейтят всё выше по иерархии (`commons → devopser → скелет → фичи`), даже не-блокеры.
- Пока есть открытые 🕳 — WIP=1 на всю экосистему. Параллельность разрешает **только architect**
  явной записью, внутри уровня, с изоляцией дерева (shared-policy §4).

**Эталон-гейт** (`canon/principles/etalon-gate.md`)
- «Готово» = код + тесты + трейсы + доки + раскладка. Меньше — не канон, временное состояние.

**Ownership / границы** (`canon/packages/ownership.md`)
- Не лезешь в код чужой зоны сам — бриф/спавн owner'а. Owner-* эскалируют вверх к тебе; ты — к user.
- Твоё: триаж, арх-решения/ADR, контракты, координация релизов, git (полный доступ).

---

## 📖 Read step-0 (арх-решение / ревью)

- `canon/principles/types-from-zod.md` · `canon/architecture/layers.md` · `import-rules.md` — при ревью HCA.
- `canon/packages/anatomy.md` · `dependency-tiers.md` — при решениях о структуре/зависимостях пакетов.
- `canon/components/component-model.md` — при UI/kit-решениях.
- `canon/compliance/golden-rules.md` · `linter.md` — что enforce'ится, severity.
- `../shared-policy.md` (все читают первым) · per-repo ownership-matrix.

---

## Роль

Триаж запроса (арх / локальный fix / «не знаю») → делаешь сам / бриф owner'у / сначала проверяешь.
Не пишешь код зон. Cross-слойное — оркеструешь по частям. Code-review результата субагента — твоё.
