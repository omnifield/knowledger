-- Cross-product proposals: узел может прийти в ЧУЖОЙ workspace как карантинная предложка
-- (origin='proposal'), исключённая из роадмапа/дерева, пока явный accept не приземлит её
-- (origin->'native' + реальный parent). decline -> origin='declined' (история, вне inbox и
-- дерева) -- у KB нет статусов, поэтому терминальное состояние несёт сам origin. Agent-agnostic:
-- движок энфорсит только СТРУКТУРНЫЙ гейт (без accept в роадмап ничто не входит) + пишет, кто
-- предложил; КТО вправе принимать -- концерн identity, не здесь.
-- SQLite<->Postgres совместимо: TEXT-колонки, ADD COLUMN с дефолтами (бэкфилл existing -> 'native').

-- +goose Up
-- +goose StatementBegin
ALTER TABLE nodes ADD COLUMN origin TEXT NOT NULL DEFAULT 'native';
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE nodes ADD COLUMN proposed_by TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE nodes ADD COLUMN source_ws TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd
-- +goose StatementBegin
-- Inbox-lookup: предложки workspace. Он же питает роадмап-фильтр (origin='native').
CREATE INDEX idx_nodes_origin ON nodes(workspace_id, origin);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_nodes_origin;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE nodes DROP COLUMN source_ws;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE nodes DROP COLUMN proposed_by;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE nodes DROP COLUMN origin;
-- +goose StatementEnd
