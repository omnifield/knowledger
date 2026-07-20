package service

import (
	"context"
	"fmt"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// CreateRefInput — вход создания forward-ребра из узла (from) в ToNode.
type CreateRefInput struct {
	ToNode string // UUID или key целевого узла (cross-workspace допустим)
	Kind   string // kb.Ref* : link|transclude|relates|blocks|depends_on|duplicate
	Actor  string
}

// AddRef создаёт forward-ребро from -> to заданного вида. Оба узла обязаны существовать; kind
// валиден; дубль (from,to,kind) -> конфликт. Backlinks для to появляются автоматически (DERIVED).
func (s *Service) AddRef(ctx context.Context, fromIDOrKey string, in CreateRefInput) (kb.Ref, error) {
	if !kb.ValidRefKind(in.Kind) {
		return kb.Ref{}, fmt.Errorf("%w: unknown ref kind %q", ErrValidation, in.Kind)
	}
	var out kb.Ref
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		from, err := s.resolveNode(ctx, q, fromIDOrKey)
		if err != nil {
			return err
		}
		to, err := s.resolveNode(ctx, q, in.ToNode)
		if err != nil {
			return err
		}
		existing, err := q.ListRefsFrom(ctx, from.ID)
		if err != nil {
			return err
		}
		for _, r := range existing {
			if r.ToNode == to.ID && r.Kind == in.Kind {
				return fmt.Errorf("%w: ref %s -> %s (%s) already exists", ErrConflict, from.Key, to.Key, in.Kind)
			}
		}
		ref, err := q.CreateRef(ctx, store.CreateRefParams{
			ID: s.newID(), FromNode: from.ID, ToNode: to.ID, Kind: in.Kind, CreatedAt: s.nowStr(),
		})
		if err != nil {
			return err
		}
		out = mapRef(ref)
		return nil
	})
	return out, err
}

// ListRefs — исходящие рёбра узла (forward).
func (s *Service) ListRefs(ctx context.Context, idOrKey string) ([]kb.Ref, error) {
	out := []kb.Ref{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		refs, err := q.ListRefsFrom(ctx, n.ID)
		if err != nil {
			return err
		}
		for _, r := range refs {
			out = append(out, mapRef(r))
		}
		return nil
	})
	return out, err
}

// ListBacklinks — DERIVED: входящие рёбра (кто ссылается на узел), обратным запросом по to_node.
func (s *Service) ListBacklinks(ctx context.Context, idOrKey string) ([]kb.Ref, error) {
	out := []kb.Ref{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		refs, err := q.ListBacklinks(ctx, n.ID)
		if err != nil {
			return err
		}
		for _, r := range refs {
			out = append(out, mapRef(r))
		}
		return nil
	})
	return out, err
}

// DeleteRef удаляет ребро по id (идемпотентно: отсутствующее -> no-op).
func (s *Service) DeleteRef(ctx context.Context, id string) error {
	return s.store.Tx(ctx, func(q *store.Queries) error {
		return q.DeleteRef(ctx, id)
	})
}
