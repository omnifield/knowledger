-- name: CreateRef :one
INSERT INTO refs (id, from_node, to_node, kind, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- Forward edges out of a node (its outgoing links/relations).
-- name: ListRefsFrom :many
SELECT * FROM refs WHERE from_node = ? ORDER BY created_at, id;

-- Backlinks: incoming edges (who points at this node) -- the DERIVED reverse query
-- over to_node, served by idx_refs_to (no duplicated backlink table).
-- name: ListBacklinks :many
SELECT * FROM refs WHERE to_node = ? ORDER BY created_at, id;

-- name: DeleteRef :exec
DELETE FROM refs WHERE id = ?;
