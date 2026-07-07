# Ярусы зависимостей — portable-tier

## Тезис

Часть пакетов — **portable-tier**: продукт «в мир» с **нулём ecosystem-deps**. Их могут импортить наши пакеты; они наши — **НЕТ**. Направление одностороннее, enforced CI.

Кандидаты: `renderer` (рендер UI по JSON-схеме, «обобщённый Widget», stateless), `canvas` (движок-адаптеры), `utils` (чистые хелперы). Признак — пакет полезен и вне Omnifield.

## Правило

- Portable-tier пакет **не импортит** `@omnifield/*` (кроме других portable-tier).
- Обратное — можно: любой наш пакет импортит portable-tier свободно.
- Никаких скрытых ecosystem-зависимостей (zod-shim, web-core, style-токены) внутри portable-пакета — они делают его непереносимым.

## Почему

Пакет, который мы отдаём наружу (или хотим отдать), с корнем в нашу экосистему — не переносим: потребитель тянет весь `@omnifield/*` хвост. Portable-tier = чистая граница: self-contained, zero-ecosystem, тестируется и публикуется без остального монорепо. Это частный, самый жёсткий случай [modules-no-crutches](../principles/modules-no-crutches.md).

## Enforcement

CI-gate (линтер / dep-cruise) режет `@omnifield/*`-импорт внутри portable-tier пакета. Нарушение — structural error, валит пайплайн ([golden-rules](../compliance/golden-rules.md)). Список portable-пакетов — явный, ведётся в конфиге линтера.
