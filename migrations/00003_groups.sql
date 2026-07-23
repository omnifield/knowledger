-- Группы — организационный слой НАД workspace: вложенные «папки» (Продукты · Концепты ·
-- Каталог) произвольной глубины для сайдбара. Группа НЕ трогает дерево узлов/контент —
-- только раскладку зон. Одна группа на workspace (папки), без группы → корень сайдбара.
-- SQLite<->Postgres совместимо: TEXT-колонки, self-FK parent_id (вложенность), ADD COLUMN
-- с REFERENCES (nullable, дефолт NULL). Общая меха (MECH-4): паттерн лифтится в brainer/chatter.

-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL,
  parent_id  TEXT REFERENCES groups(id),     -- nullable self-FK → вложенность любой глубины
  ord        TEXT NOT NULL DEFAULT '',       -- дробный индекс порядка среди сиблингов
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_groups_parent ON groups(parent_id);
-- +goose StatementEnd
-- +goose StatementBegin
-- Членство workspace в группе (одна группа = папки). NULL → без группы (корень сайдбара).
ALTER TABLE workspaces ADD COLUMN group_id TEXT REFERENCES groups(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workspaces DROP COLUMN group_id;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_groups_parent;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
