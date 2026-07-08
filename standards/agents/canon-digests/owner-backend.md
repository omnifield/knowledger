# Canon digest — owner-backend

Выжимка для owner-агента backend-сервисов (python / rust). Полный канон — `../../canon/`,
карта — `../canon-map.md`. Это **не замена** канона: `⬛` — железные «нельзя» в контексте,
`📖` — прочитать первым действием.

---

## ⬛ Non-negotiables (всегда в силе)

**Модули, не монолит** (`canon/principles/modules-no-crutches.md`)
- Сервис самодостаточен: не пускает корни в соседа. `backend/` ⊥ `framework/` ⊥ `apps/`.
- Связь между зонами — только через опубликованный пакет / контракт (`commons/contracts`,
  версионируемый), НЕ прямой импорт реализации соседа.
- Нужен адаптер/шим, чтобы состыковать своих → граница спроектирована неверно. СТОП, не пиши адаптер.

**Лечим причину, не симптом** (`canon/principles/root-cause-not-symptom.md`)
- Стоп-сигналы костыля — подход в корне неверный, СТОП: **hardcoded путь**, **silent fallback**
  («не нашлось — тихо возьмём другое»), **try-catch-swallow**, куча нестабильных шагов.
- Gap пакета/сервиса → фикс в **зоне-владельце**, не workaround в потребителе.
- «Не знаю, работает ли» → проверь инструментом (запусти тест / прочитай source), не гипотезы.
- **2 неудачных фикса подряд → СТОП + диагностика корня**, не третья попытка вслепую.

**Границы зоны / ownership** (`canon/packages/ownership.md`)
- Правишь **только свою зону**. Чужая зона / cross-зонный концерн → эскалация **вверх** к architect,
  не лезешь сам. Эскалация строго вверх, никогда вбок.
- Git: **commit-only** (conventional `fix(backend): …`), push/merge — architect после ревью.

**Типы из схемы** (`canon/principles/types-from-zod.md`, backend-проекция)
- Source of truth типов — схема (pydantic / SQLAlchemy), не ручные дубли под домен.

---

## 📖 Read step-0 (перед нетривиальной правкой)

- `canon/principles/foundation-first.md` — дыра базы гейтит всё выше; WIP=1 пока база не зелёная.
- `canon/principles/etalon-gate.md` — что значит «готово» (код + тесты + трейсы + доки + раскладка).
- `canon/packages/anatomy.md` — структура зоны/пакета.
- `canon/packages/dependency-tiers.md` — portable-tier, что кому можно импортить.
- `canon/compliance/golden-rules.md` — enforce'имые правила + severity.
- Зона: свой `OWNERSHIP.md` / AI-anchor (per-repo) + `../shared-policy.md` (читают все первым).

---

## Не для owner-backend (skip)

`architecture/*` (HCA-слои, ui-proxy, namespaces) и `components/*` (kit/tokens/registration) —
это UI/HCA-домен. Если задача вдруг упирается в них — она не backend, эскалируй architect'у.
