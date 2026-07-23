package kb

import "encoding/json"

// Workspace — a scope tier: a product OR a concern. Key is short and stable.
// GroupID is an optional sidebar group membership (organizational layer only —
// it never touches the node tree or content); nil means ungrouped (sidebar root).
type Workspace struct {
	ID          string  `json:"id"`
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	GroupID     *string `json:"group_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Group — a nestable sidebar folder over workspaces (organizational layer only,
// not knowledge). ParentID (self-ref) gives arbitrary nesting; nil is a root
// group. Ord is a fractional sibling index. Reusable mech (MECH-4).
type Group struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ParentID  *string `json:"parent_id"`
	Ord       string  `json:"ord"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// Node — a recursive, type-less content node. Documents, sections and articles
// are all nodes; Kind is an optional hint, not an enforced type. Key is stable
// and immutable so external references survive moves and renames (the dual-id
// rule). Ord is a fractional index for sibling ordering — reorder a sibling
// without renumbering the rest.
type Node struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"workspace_id"`
	Key         string  `json:"key"`
	Seq         int64   `json:"seq"`
	ParentID    *string `json:"parent_id"`
	Kind        *string `json:"kind"`
	Title       string  `json:"title"`
	Body        string  `json:"body"`
	Ord         string  `json:"ord"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	// Origin: "native" (roadmap/tree node), "proposal" (cross-product suggestion in the inbox,
	// not yet accepted) or "declined". ProposedBy/SourceWs carry a proposal's provenance.
	Origin     string `json:"origin"`
	ProposedBy string `json:"proposed_by,omitempty"`
	SourceWs   string `json:"source_ws,omitempty"`
}

// Ref — a forward reference edge from one node to another. A single edge table
// carries both content links (link/tag/transclude) and typed relations
// (relates/blocks/depends_on/duplicate); backlinks are DERIVED by a reverse
// query over ToNode, never kept in a duplicated table to be synchronised.
type Ref struct {
	ID        string `json:"id"`
	FromNode  string `json:"from_node"`
	ToNode    string `json:"to_node"`
	Kind      string `json:"kind"`
	CreatedAt string `json:"created_at"`
}

// Tag — a per-workspace label applied to nodes through a tag Ref.
type Tag struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
}

// Revision — an immutable point-in-time snapshot of a node, kept as a separate
// first-class resource rather than inline node state. Snapshot is the full node
// copy at that revision; a restore re-applies it.
type Revision struct {
	ID        string          `json:"id"`
	NodeID    string          `json:"node_id"`
	Snapshot  json.RawMessage `json:"snapshot"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"created_at"`
}

// Activity — a typed timeline entry for a node (kind = commented, created, …).
// Data is an arbitrary JSON payload.
type Activity struct {
	ID        string          `json:"id"`
	NodeID    string          `json:"node_id"`
	Actor     string          `json:"actor"`
	Kind      string          `json:"kind"`
	Data      json.RawMessage `json:"data,omitempty"`
	CreatedAt string          `json:"created_at"`
}

// TreeNode — a node plus its subtree, so a whole tree ships in one request.
// Cross-node relations are not included here (the tree is parent/child only).
type TreeNode struct {
	Node
	Children []TreeNode `json:"children"`
}
