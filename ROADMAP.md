# ROADMAP — живой борд

Единый источник правды «где всё сейчас» — **GitHub Project**, не этот файл.
Вопрос «что происходит и кто что делает» — ответ на борде, не в переписке и не в markdown.

## 🔗 https://github.com/orgs/omnifield/projects/1

Борд поверх всех репо (`brainer` · `commons` · `writer` · `devopser`). Поля = колонки старого roadmap:

| Поле | Значения |
|---|---|
| `Stage` | `Proposed` · `Ready` · `Doing` · `Blocked` · `Done` |
| `Layer` | `Contract` · `DevOps` · `Backend` · `App` · `Frontend` |
| `Repo` | авто из issue |
| `Queue` | номер в очереди запусков (WIP=1) |
| `Blocked-by` | ссылка/текст блокера |
| `DoR` / `DoD` | ✅ met · ⬜ not yet |

## Как читать (WIP=1, foundation-first)

- Вид **Board** — канбан по `Stage`; вид **Table** + сортировка по `Queue` — очередь запусков.
- Бери верхний item своей роли с пустым `Blocked-by` и минимальным `Queue`.
- Очередь **строго последовательная**: следующий стартует только когда предыдущий закрыт.
- Дыры (foundation-first) гейтят всё выше: пока workstation DoD не зелёный — brainer стоит.

## Задача = issue

- Карточка открывается в issue репо-владельца: сверху 🧭 **TL;DR** (тейки для user, без деталей),
  ниже — **бриф исполнителю** (факты / fix / verify).
- Вопросы — **комментами** в issue (канал общения, история сохраняется).
- User: запускает сессии по очереди борда + push/merge там, где просит architect.

---

_Markdown-таблица ROADMAP ретайрнута 2026-07-08 — заменена бордом (стоппер «roadmap неинформативен»).
История стадий живёт в самих issue._
