-- name: CreateRevision :one
INSERT INTO revisions (id, node_id, snapshot, meta, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRevision :one
SELECT * FROM revisions WHERE id = ?;

-- Version history of a node, newest first.
-- name: ListRevisions :many
SELECT * FROM revisions WHERE node_id = ? ORDER BY created_at DESC, id DESC;

-- All revisions of a node -- cleaned when the node is deleted.
-- name: DeleteRevisionsForNode :exec
DELETE FROM revisions WHERE node_id = ?;
