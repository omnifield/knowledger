package service

import (
	"context"
	"database/sql"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// GetSubtree возвращает узел и всё его поддерево (один запрос вместо N обходов). maxDepth < 0 —
// без ограничения; maxDepth = 0 — только сам узел (дети не раскрываются). Рёбра-refs сюда не
// входят (дерево — parent/child; связи — refs/backlinks).
func (s *Service) GetSubtree(ctx context.Context, idOrKey string, maxDepth int) (kb.TreeNode, error) {
	var out kb.TreeNode
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		out, err = s.buildSubtree(ctx, q, n, maxDepth)
		return err
	})
	return out, err
}

// GetWorkspaceTree возвращает корни workspace + их поддеревья. maxDepth — как в GetSubtree.
func (s *Service) GetWorkspaceTree(ctx context.Context, wsIDOrKey string, maxDepth int) ([]kb.TreeNode, error) {
	out := []kb.TreeNode{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		roots, err := q.ListRootNodes(ctx, ws.ID)
		if err != nil {
			return err
		}
		for _, r := range roots {
			t, err := s.buildSubtree(ctx, q, r, maxDepth)
			if err != nil {
				return err
			}
			out = append(out, t)
		}
		return nil
	})
	return out, err
}

// buildSubtree рекурсивно собирает узел и его детей. depth — оставшаяся глубина спуска
// (< 0 — без лимита; = 0 — стоп, дети не раскрываются).
func (s *Service) buildSubtree(ctx context.Context, q *store.Queries, n store.Node, depth int) (kb.TreeNode, error) {
	node := kb.TreeNode{Node: mapNode(n), Children: []kb.TreeNode{}}
	if depth == 0 {
		return node, nil
	}
	kids, err := q.ListChildren(ctx, sql.NullString{String: n.ID, Valid: true})
	if err != nil {
		return kb.TreeNode{}, err
	}
	for _, k := range kids {
		child, err := s.buildSubtree(ctx, q, k, depth-1)
		if err != nil {
			return kb.TreeNode{}, err
		}
		node.Children = append(node.Children, child)
	}
	return node, nil
}
