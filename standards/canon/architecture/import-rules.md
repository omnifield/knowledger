# Правила импорта

Enforced линтером ([golden-rules](../compliance/golden-rules.md)). Слой защищён независимо от того, откуда берётся ref.

## Четыре правила

1. **No Upward Imports.** Нижний слой не импортирует верхний. Entity не знает про Widget; View не знает про Controller.
2. **No Horizontal Imports.** `View.A` не импортирует `View.B`. `Controller.A` не знает о `Controller.B`. Только композиция в Widget или цепочка через `next()` к родительской Feature.
3. **Stateless View / Shape.** Никакого состояния, ноль явных импортов (JSX-рантайм — авто, не ручной import). Ноль ручных `type`/`interface` на уровне аппа — типы только из zod-Entity (`z.infer`), глобальных деклараций или vendor-типов ([types-from-zod](../principles/types-from-zod.md)).
4. **Composition Only in Widgets.** Одна View не может «жёстко» использовать другую — только через children/slots на уровне Widget.

## Почему

Слои — направленный граф. Upward/horizontal рёбра превращают его в клубок: нельзя переиспользовать узел без его соседей, нельзя тестировать изолированно, нельзя двигать. Композиция вынесена в **одну точку** (Widget) — там и только там разрешено «склеивать».

## Как связывать без импорта

- **Композиция** сущностей → **Widget** (children/slots).
- **Поведение вверх по цепочке** → `next()`: Controller делегирует необработанное родительскому Controller/Feature прямым вызовом (не через event-bus), естественный `await`-возврат.
- **События от UI** → тег-флоу ([ui-proxy-tag-flow](ui-proxy-tag-flow.md)): децентрализованный поиск N→1, инпут всплывает к хендлеру по тегу.
- **Именованная семантика между зонами** → `useEmit`/типизированные события (onLogin/onImport), opt-in.
- **Данные из родителя** → `useCtx().store.ctx` (reactive).

## Границы зон

Cross-zone импорт реализации соседа — нарушение ([modules-no-crutches](../principles/modules-no-crutches.md)). Между пакетами/репо — только контракт (`commons/contracts`) или опубликованный API, не прямой импорт внутренностей.
