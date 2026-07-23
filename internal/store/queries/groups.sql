-- name: CreateGroup :one
INSERT INTO groups (id, name, parent_id, ord, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetGroup :one
SELECT * FROM groups WHERE id = ?;

-- name: ListGroups :many
SELECT * FROM groups ORDER BY ord, created_at;

-- name: UpdateGroup :one
UPDATE groups
SET name = ?, parent_id = ?, ord = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = ?;

-- name: CountChildGroups :one
SELECT COUNT(*) FROM groups WHERE parent_id = ?;

-- name: CountGroupWorkspaces :one
SELECT COUNT(*) FROM workspaces WHERE group_id = ?;
