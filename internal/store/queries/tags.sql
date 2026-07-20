-- name: CreateTag :one
INSERT INTO tags (id, workspace_id, name, created_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTag :one
SELECT * FROM tags WHERE id = ?;

-- name: ListTagsByWorkspace :many
SELECT * FROM tags WHERE workspace_id = ? ORDER BY name;

-- name: AddNodeTag :exec
INSERT INTO node_tags (node_id, tag_id) VALUES (?, ?);

-- name: RemoveNodeTag :exec
DELETE FROM node_tags WHERE node_id = ? AND tag_id = ?;

-- All tag memberships of a node -- cleaned when the node is deleted.
-- name: DeleteNodeTagsForNode :exec
DELETE FROM node_tags WHERE node_id = ?;

-- Tags applied to a node (join across the m2m).
-- name: ListNodeTags :many
SELECT t.* FROM tags t
JOIN node_tags nt ON nt.tag_id = t.id
WHERE nt.node_id = ?
ORDER BY t.name;
