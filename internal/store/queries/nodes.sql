-- name: CreateNode :one
INSERT INTO nodes (id, workspace_id, seq, key, parent_id, kind, title, body, ord, origin, proposed_by, source_ws, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetNodeByID :one
SELECT * FROM nodes WHERE id = ?;

-- name: GetNodeByKey :one
SELECT * FROM nodes WHERE key = ?;

-- Flat list of a workspace's roadmap nodes (origin='native' -- proposals/declined live outside
-- the roadmap). Ordered by the fractional index then seq. parent filter applied in-memory.
-- name: ListNodesByWorkspace :many
SELECT * FROM nodes WHERE workspace_id = ? AND origin = 'native' ORDER BY ord, seq;

-- Tree roots of a workspace (top level: no parent). Excludes un-accepted proposals (they are
-- parent-less like roots but must not surface in the roadmap until accepted).
-- name: ListRootNodes :many
SELECT * FROM nodes WHERE workspace_id = ? AND parent_id IS NULL AND origin = 'native' ORDER BY ord, seq;

-- Inbox: pending cross-product proposals of a workspace, awaiting accept/decline.
-- name: ListInboxNodes :many
SELECT * FROM nodes WHERE workspace_id = ? AND origin = 'proposal' ORDER BY seq;

-- name: ListChildren :many
SELECT * FROM nodes WHERE parent_id = ? ORDER BY ord, seq;

-- name: CountChildren :one
SELECT COUNT(*) FROM nodes WHERE parent_id = ?;

-- name: UpdateNode :one
UPDATE nodes
SET kind = ?, title = ?, body = ?, ord = ?, parent_id = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- Accept a proposal into the roadmap: flip origin to native and set the parent the receiving
-- architect chose. Stable key is preserved (inbound references stay valid).
-- name: AcceptProposalNode :one
UPDATE nodes SET parent_id = ?, origin = 'native', updated_at = ? WHERE id = ?
RETURNING *;

-- Decline a proposal: origin -> 'declined' (terminal; out of inbox and roadmap, kept as history).
-- name: DeclineProposalNode :one
UPDATE nodes SET origin = 'declined', updated_at = ? WHERE id = ?
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes WHERE id = ?;
