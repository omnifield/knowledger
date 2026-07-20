-- name: CreateNode :one
INSERT INTO nodes (id, workspace_id, seq, key, parent_id, kind, title, body, ord, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetNodeByID :one
SELECT * FROM nodes WHERE id = ?;

-- name: GetNodeByKey :one
SELECT * FROM nodes WHERE key = ?;

-- Flat list of a workspace's nodes, ordered by the fractional index then seq for a
-- stable tie-break. parent/kind filters are applied in-memory by the service layer
-- (v0 scale; keeps SQL portable, no dynamic SQL).
-- name: ListNodesByWorkspace :many
SELECT * FROM nodes WHERE workspace_id = ? ORDER BY ord, seq;

-- Tree roots of a workspace (top level: no parent).
-- name: ListRootNodes :many
SELECT * FROM nodes WHERE workspace_id = ? AND parent_id IS NULL ORDER BY ord, seq;

-- name: ListChildren :many
SELECT * FROM nodes WHERE parent_id = ? ORDER BY ord, seq;

-- name: CountChildren :one
SELECT COUNT(*) FROM nodes WHERE parent_id = ?;

-- name: UpdateNode :one
UPDATE nodes
SET kind = ?, title = ?, body = ?, ord = ?, parent_id = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes WHERE id = ?;
