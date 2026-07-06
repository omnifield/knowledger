# Неймспейсы

Глобальные реестры слоёв (`Views`, `Widgets`, `Shapes`, `Controllers`, `Features`, `Entities`) — **nested по структуре папок**.

## Правило

Папка = уровень неймспейса, имя файла = leaf.

```
widgets/forms/auth.tsx        → Widgets.Forms.Auth
views/viewer/loginForm.tsx    → Views.Viewer.LoginForm
views/hello.tsx               → Views.Hello          (корневой файл без папки)
```

**Не** flat (`Views.AuthLoginForm` — неправильно). Вложенность отражает дерево папок один-в-один.

## В пакете (домен)

Внутри доменного пакета неймспейс — по имени сущности/блока: `Learn.Lessons`, `Learn.Nav.Main`, `WebStudio.Inspector`. Вложенность безопасна для рендера (он вложенность поднимает), но **codegen `.Events`** привязан к leaf — вкладывать блок с собственным `__events` в под-неймспейс нельзя без дропа его событий (иначе `.Events`-кодген не найдёт). Прецедент: `Learn.Library.Info` (без своих событий) вложился штатно; nav-блоки — только после дропа собственного `__events`.

## Типы слотов

Живут в сгенерированном `CapsuleSlots` (`.capsule/@types/slots.d.ts`). Каждое property типизировано как `typeof import('@<layer>/...').default` — Ctrl+Click ведёт в источник. Требует `export default` в каждом файле слоя ([[layers]] конвенция export).
