-- Канон-схема knowledger: база знаний как рекурсивный type-less узел + forward-рёбра
-- (backlinks derived) + тэги + версии + activity. SQL — совместимый SQLite↔Postgres
-- (канон go.md): только TEXT/INTEGER, ноль AUTOINCREMENT/SERIAL (id/seq генерит app-слой),
-- ноль engine-специфики. FK ОБЪЯВЛЕНЫ (документируют инвариант + PRAGMA foreign_keys для
-- SQLite), но enforce'им ещё и сами в сервис-слое (Jira-урок). Таймстемпы — TEXT RFC3339.
-- Дизайн — docs/knowledge-base-canon.md. Proposal/inbox-колонки придут отдельной миграцией
-- (KNOW-9), как у tasker 00002.

-- +goose Up
-- +goose StatementBegin
CREATE TABLE workspaces (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,          -- короткий человекочитаемый префикс, напр. "KNOW"
  name        TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  node_seq    INTEGER NOT NULL DEFAULT 0,    -- per-workspace счётчик node.seq (транзакционный, не убывает)
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Node — рекурсивный type-less узел контента (документ/раздел/статья — всё узлы). Дерево через
-- parent_id (primary-структура), НЕ через ref. key = "<WS>-<seq>" СТАБИЛЬНЫЙ/immutable (dual-id).
-- kind — опциональная подсказка (не enforced-тип). ord — ДРОБНЫЙ индекс порядка сиблингов
-- (реордер без перенумерации), TEXT для лексикографической сортировки.
CREATE TABLE nodes (
  id           TEXT PRIMARY KEY,
  workspace_id TEXT NOT NULL REFERENCES workspaces(id),
  seq          INTEGER NOT NULL,             -- per-workspace порядковый (unique с workspace_id)
  key          TEXT NOT NULL,                -- "<WS>-<seq>", immutable
  parent_id    TEXT REFERENCES nodes(id),    -- nullable self-FK → рекурсия/дерево
  kind         TEXT,                         -- nullable: type-less (article|section|tag|...)
  title        TEXT NOT NULL,
  body         TEXT NOT NULL DEFAULT '',
  ord          TEXT NOT NULL DEFAULT '',     -- дробный индекс порядка сиблингов
  created_at   TEXT NOT NULL,
  updated_at   TEXT NOT NULL,
  UNIQUE (workspace_id, seq)
);
-- +goose StatementEnd
-- +goose StatementBegin
-- key глобально уникален (WS-префиксы уникальны) → dual-id резолв по key однозначен без ws-контекста.
CREATE UNIQUE INDEX idx_nodes_key ON nodes(key);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_nodes_ws ON nodes(workspace_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_nodes_parent ON nodes(parent_id);
-- +goose StatementEnd

-- +goose StatementBegin
-- Ref — forward-ребро узел→узел. Одна таблица несёт content-links (link|transclude) и typed-
-- relations (relates|blocks|depends_on|duplicate); cross-workspace допустимо. BACKLINKS DERIVED:
-- обратный запрос по to_node (idx_refs_to), НЕ отдельная синхронизируемая таблица.
CREATE TABLE refs (
  id         TEXT PRIMARY KEY,
  from_node  TEXT NOT NULL REFERENCES nodes(id),
  to_node    TEXT NOT NULL REFERENCES nodes(id),
  kind       TEXT NOT NULL,                  -- link|transclude|relates|blocks|depends_on|duplicate
  created_at TEXT NOT NULL,
  UNIQUE (from_node, to_node, kind)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_refs_from ON refs(from_node);
-- +goose StatementEnd
-- +goose StatementBegin
-- Обратный индекс — backlinks/«кто ссылается на узел» одним индекс-сканом.
CREATE INDEX idx_refs_to ON refs(to_node);
-- +goose StatementEnd

-- +goose StatementBegin
-- Tag — per-workspace метка (отдельная сущность, m2m как tasker labels — не ref).
CREATE TABLE tags (
  id           TEXT PRIMARY KEY,
  workspace_id TEXT NOT NULL REFERENCES workspaces(id),
  name         TEXT NOT NULL,
  created_at   TEXT NOT NULL,
  UNIQUE (workspace_id, name)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE node_tags (
  node_id TEXT NOT NULL REFERENCES nodes(id),
  tag_id  TEXT NOT NULL REFERENCES tags(id),
  PRIMARY KEY (node_id, tag_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Revision — неизменяемый снимок узла как ОТДЕЛЬНЫЙ ресурс (не инлайн-состояние). snapshot —
-- полная копия узла на момент версии (JSON-blob); restore переприменяет её.
CREATE TABLE revisions (
  id         TEXT PRIMARY KEY,
  node_id    TEXT NOT NULL REFERENCES nodes(id),
  snapshot   TEXT NOT NULL DEFAULT '',       -- JSON-копия узла
  meta       TEXT NOT NULL DEFAULT '',       -- JSON-метаданные (autosave/причина/...)
  created_at TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_revisions_node ON revisions(node_id);
-- +goose StatementEnd

-- +goose StatementBegin
-- Activity — typed timeline с 1-го дня. data = JSON-blob (payload события). comments — kind='commented'.
CREATE TABLE activity (
  id         TEXT PRIMARY KEY,
  node_id    TEXT NOT NULL REFERENCES nodes(id),
  actor      TEXT NOT NULL,
  kind       TEXT NOT NULL,                  -- created|updated|ref_added|tagged|commented|...
  data       TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_activity_node ON activity(node_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activity;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS revisions;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS node_tags;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS refs;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS nodes;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS workspaces;
-- +goose StatementEnd
