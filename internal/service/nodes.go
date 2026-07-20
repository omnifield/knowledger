package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// CreateNodeInput — вход создания узла. ParentID принимает UUID ИЛИ key (резолвится).
type CreateNodeInput struct {
	Title    string
	Body     string
	Kind     *string
	ParentID *string
	Actor    string
}

// UpdateNodeInput — частичный PATCH. Пары Set*/значение различают «поле не прислано» от
// «поле выставлено в null» (очистка kind/parent).
type UpdateNodeInput struct {
	Title     *string
	Body      *string
	SetKind   bool
	Kind      *string // при SetKind: nil => очистить kind
	SetParent bool
	ParentID  *string // при SetParent: nil => сделать корнем
	Actor     string
}

// NodeFilter — фильтр списка (накладывается in-memory). ParentSet=true включает фильтр по родителю.
type NodeFilter struct {
	ParentSet bool
	ParentID  *string // resolved UUID; nil => только корни
}

// CreateNode создаёт узел: аллокация стабильного key через транзакционный per-workspace счётчик,
// ord = дробный индекс порядка (v1: zero-padded seq — append в конец; реордер-между — позже),
// enforce FK (workspace/parent существуют и согласованы), activity=created.
func (s *Service) CreateNode(ctx context.Context, wsIDOrKey string, in CreateNodeInput) (kb.Node, error) {
	if in.Title == "" {
		return kb.Node{}, fmt.Errorf("%w: title required", ErrValidation)
	}
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		parent, err := s.resolveParentInWorkspace(ctx, q, ws.ID, in.ParentID)
		if err != nil {
			return err
		}
		seq, err := q.BumpWorkspaceNodeSeq(ctx, store.BumpWorkspaceNodeSeqParams{
			UpdatedAt: s.nowStr(), ID: ws.ID,
		})
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s-%d", ws.Key, seq)
		now := s.nowStr()
		node, err := q.CreateNode(ctx, store.CreateNodeParams{
			ID: s.newID(), WorkspaceID: ws.ID, Seq: seq, Key: key,
			ParentID: parent, Kind: ptrToNull(in.Kind), Title: in.Title, Body: in.Body,
			Ord: fmt.Sprintf("%012d", seq), CreatedAt: now, UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		if err := s.addActivity(ctx, q, node.ID, in.Actor, "created", map[string]any{"key": key}); err != nil {
			return err
		}
		out = mapNode(node)
		return nil
	})
	if err != nil {
		return kb.Node{}, err
	}
	return out, nil
}

// GetNode резолвит узел по UUID или key (dual-id).
func (s *Service) GetNode(ctx context.Context, idOrKey string) (kb.Node, error) {
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		out = mapNode(n)
		return nil
	})
	return out, err
}

// ListNodes возвращает узлы workspace с опц-фильтром parent (in-memory).
func (s *Service) ListNodes(ctx context.Context, wsIDOrKey string, f NodeFilter) ([]kb.Node, error) {
	out := []kb.Node{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		nodes, err := q.ListNodesByWorkspace(ctx, ws.ID)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			if f.ParentSet && !matchNull(n.ParentID, f.ParentID) {
				continue
			}
			out = append(out, mapNode(n))
		}
		return nil
	})
	return out, err
}

// ListChildren — прямые дети узла (dual-id родитель).
func (s *Service) ListChildren(ctx context.Context, idOrKey string) ([]kb.Node, error) {
	out := []kb.Node{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		parent, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		kids, err := q.ListChildren(ctx, sql.NullString{String: parent.ID, Valid: true})
		if err != nil {
			return err
		}
		for _, n := range kids {
			out = append(out, mapNode(n))
		}
		return nil
	})
	return out, err
}

// UpdateNode применяет частичный PATCH, enforce'я FK/цикл.
func (s *Service) UpdateNode(ctx context.Context, idOrKey string, in UpdateNodeInput) (kb.Node, error) {
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		p := store.UpdateNodeParams{
			Kind: n.Kind, Title: n.Title, Body: n.Body, Ord: n.Ord,
			ParentID: n.ParentID, UpdatedAt: s.nowStr(), ID: n.ID,
		}
		if in.Title != nil {
			if *in.Title == "" {
				return fmt.Errorf("%w: title cannot be empty", ErrValidation)
			}
			p.Title = *in.Title
		}
		if in.Body != nil {
			p.Body = *in.Body
		}
		if in.SetKind {
			p.Kind = ptrToNull(in.Kind)
		}
		if in.SetParent {
			parent, err := s.resolveParentForUpdate(ctx, q, n, in.ParentID)
			if err != nil {
				return err
			}
			p.ParentID = parent
		}
		updated, err := q.UpdateNode(ctx, p)
		if err != nil {
			return err
		}
		out = mapNode(updated)
		return nil
	})
	return out, err
}

// DeleteNode удаляет узел. Запрещаем удаление узла с детьми (app-level FK enforce). Зависимые
// строки (рёбра в обе стороны, тэг-членства, ревизии, activity) чистим в той же транзакции:
// БД-FK RESTRICT, а у узла всегда есть хотя бы activity=created — иначе DELETE падает на FK.
func (s *Service) DeleteNode(ctx context.Context, idOrKey string) error {
	return s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		cnt, err := q.CountChildren(ctx, sql.NullString{String: n.ID, Valid: true})
		if err != nil {
			return err
		}
		if cnt > 0 {
			return fmt.Errorf("%w: node %s has %d children; reparent or delete them first", ErrConflict, n.Key, cnt)
		}
		if err := q.DeleteRefsForNode(ctx, store.DeleteRefsForNodeParams{FromNode: n.ID, ToNode: n.ID}); err != nil {
			return err
		}
		if err := q.DeleteNodeTagsForNode(ctx, n.ID); err != nil {
			return err
		}
		if err := q.DeleteRevisionsForNode(ctx, n.ID); err != nil {
			return err
		}
		if err := q.DeleteActivityForNode(ctx, n.ID); err != nil {
			return err
		}
		return q.DeleteNode(ctx, n.ID)
	})
}

// --- FK/цикл enforce (app-level) -------------------------------------------

// resolveParentInWorkspace: nil-вход -> корень; иначе parent обязан существовать и лежать в том же
// workspace (дерево не пересекает границы workspace; кросс-ws — это ref, не parent_id).
func (s *Service) resolveParentInWorkspace(ctx context.Context, q *store.Queries, wsID string, idOrKey *string) (sql.NullString, error) {
	if idOrKey == nil {
		return sql.NullString{}, nil
	}
	p, err := s.resolveNode(ctx, q, *idOrKey)
	if err != nil {
		return sql.NullString{}, err
	}
	if p.WorkspaceID != wsID {
		return sql.NullString{}, fmt.Errorf("%w: parent %s is in another workspace", ErrValidation, p.Key)
	}
	return sql.NullString{String: p.ID, Valid: true}, nil
}

// resolveParentForUpdate — как выше, плюс запрет self-parent и циклов (parent не может быть
// самим узлом или его потомком).
func (s *Service) resolveParentForUpdate(ctx context.Context, q *store.Queries, node store.Node, idOrKey *string) (sql.NullString, error) {
	if idOrKey == nil {
		return sql.NullString{}, nil
	}
	p, err := s.resolveNode(ctx, q, *idOrKey)
	if err != nil {
		return sql.NullString{}, err
	}
	if p.WorkspaceID != node.WorkspaceID {
		return sql.NullString{}, fmt.Errorf("%w: parent %s is in another workspace", ErrValidation, p.Key)
	}
	if p.ID == node.ID {
		return sql.NullString{}, fmt.Errorf("%w: node cannot be its own parent", ErrValidation)
	}
	cur := p
	for cur.ParentID.Valid {
		if cur.ParentID.String == node.ID {
			return sql.NullString{}, fmt.Errorf("%w: reparent would create a cycle", ErrValidation)
		}
		cur, err = q.GetNodeByID(ctx, cur.ParentID.String)
		if err != nil {
			return sql.NullString{}, err
		}
	}
	return sql.NullString{String: p.ID, Valid: true}, nil
}

// addActivity пишет запись timeline. data сериализуется в JSON; пустой actor -> "system".
func (s *Service) addActivity(ctx context.Context, q *store.Queries, nodeID, actor, kind string, data any) error {
	raw := ""
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		raw = string(b)
	}
	if actor == "" {
		actor = "system"
	}
	_, err := q.CreateActivity(ctx, store.CreateActivityParams{
		ID: s.newID(), NodeID: nodeID, Actor: actor, Kind: kind, Data: raw, CreatedAt: s.nowStr(),
	})
	return err
}

// matchNull: узловое nullable-поле == искомому (nil == «пусто»/корень).
func matchNull(field sql.NullString, want *string) bool {
	if want == nil {
		return !field.Valid
	}
	return field.Valid && field.String == *want
}
