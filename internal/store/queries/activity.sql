-- name: CreateActivity :one
INSERT INTO activity (id, node_id, actor, kind, data, created_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- Typed timeline of a node, oldest first (chronological read).
-- name: ListActivity :many
SELECT * FROM activity WHERE node_id = ? ORDER BY created_at, id;

-- All timeline entries of a node -- cleaned when the node is deleted.
-- name: DeleteActivityForNode :exec
DELETE FROM activity WHERE node_id = ?;
