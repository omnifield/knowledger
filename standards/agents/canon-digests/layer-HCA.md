# Canon digest — layer-HCA

Выжимка для layer-агента (один HCA-артефакт: view/shape/widget/page/controller/feature/entity/
ui-component). Без git, один артефакт за вызов. Полный канон — `../../canon/`, карта — `../canon-map.md`.

---

## ⬛ Non-negotiables (всегда в силе)

**Правила импорта** (`canon/architecture/import-rules.md`)
- **No upward**, **no horizontal**. View не импортит View, Controller не знает о Controller.
- Композиция — только в Widget. Поведение вверх — `next()`. Данные — `useCtx().store.ctx`.
- **Stateless View/Shape**: ноль состояния, ноль явных импортов (JSX-рантайм — авто).

**Kit-first** (`canon/components/kit-first.md`)
- Весь UI из kit, **props-only**: ноль `class`/`className`/`<style>` в артефакте. Нет нужного вида → это задача kit (эскалируй), не raw-класс.

**Модель компонентов** (`canon/components/component-model.md`)
- Композиция — пресетом в kit. Карточка=сущность (данные по ключам, слоты вкл/выкл), НЕ ручной div-лайаут. Свои данные — своими компонентами.

**Типы из zod** (`canon/principles/types-from-zod.md`)
- `z.infer<typeof schema>` — единственный источник. Ноль ручных `interface`/`type` под домен/props.

**Причина, не симптом** (`canon/principles/root-cause-not-symptom.md`)
- Не костыль ради «сделать». Артефакт не строится по канону → **эскалируй вверх** (owner/architect), не обходи.

---

## 📖 Read step-0 (по типу артефакта)

- `canon/architecture/layers.md` — сигнатура обёртки твоего слоя, param vs global.
- `canon/architecture/namespaces.md` — путь папки = неймспейс (`widgets/forms/auth` → `Widgets.Forms.Auth`).
- `canon/architecture/ui-proxy-tag-flow.md` — event-теги (controller/feature/widget).
- `canon/components/registration.md` · `tokens.md` — регистрация, токены (view/shape/widget/ui-component).
- `canon/principles/etalon-gate.md` — «готово».
- `canon/compliance/golden-rules.md` — твои enforce'имые «нельзя».

---

## Правила листа

Один артефакт за вызов. Без git, без Bash. Максимум **один** уточняющий вопрос.
Задача шире одного артефакта / cross-слойная → это не твоё, эскалируй owner/architect.
