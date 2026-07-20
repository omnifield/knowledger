package service

// Cross-product proposals — карантинная полоса на workspace с гейтом приёмки.
//
// Архитектор продукта A кладёт предложку в workspace продукта B. Предложка — обычный узел
// (стабильный key, activity, тело) но с origin='proposal' и без parent, поэтому ИСКЛЮЧЕНА из
// роадмапа/дерева B (roots/list фильтруют origin='native') и видна только в inbox B. В роадмап
// входит лишь через явный accept (origin->'native' + parent, выбранный архитектором B). decline
// -> origin='declined' (терминал, вне inbox и дерева; у KB нет статусов).
//
// Agent-agnostic (как git-flow): движок энфорсит только СТРУКТУРНЫЙ гейт — ничто чужое не
// приземляется в роадмап без accept, ждёт в inbox — и пишет proposed_by/actor. КТО вправе
// принимать — концерн identity (мост позже), не здесь.

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// Значения node.origin. 'native' — обычный узел роадмапа/дерева; 'proposal' — входящая
// cross-product предложка в inbox; 'declined' — отклонённая (история, вне inbox и дерева).
const (
	originNative   = "native"
	originProposal = "proposal"
	originDeclined = "declined"
)

// CreateProposalInput — cross-product предложка в чужой workspace. SourceWs — key своего
// workspace (бэклинк); Actor пишется как proposed_by.
type CreateProposalInput struct {
	Title    string
	Body     string
	SourceWs string
	Actor    string
}

// CreateProposal создаёт карантинный узел (origin='proposal', без parent) в целевом workspace.
// Аллоцирует стабильный key как любой узел (чтобы предлагающий мог сослаться), но вне роадмапа.
func (s *Service) CreateProposal(ctx context.Context, wsIDOrKey string, in CreateProposalInput) (kb.Node, error) {
	if in.Title == "" {
		return kb.Node{}, fmt.Errorf("%w: title required", ErrValidation)
	}
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		seq, err := q.BumpWorkspaceNodeSeq(ctx, store.BumpWorkspaceNodeSeqParams{UpdatedAt: s.nowStr(), ID: ws.ID})
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s-%d", ws.Key, seq)
		now := s.nowStr()
		node, err := q.CreateNode(ctx, store.CreateNodeParams{
			ID: s.newID(), WorkspaceID: ws.ID, Seq: seq, Key: key,
			ParentID: sql.NullString{}, Kind: sql.NullString{}, Title: in.Title, Body: in.Body,
			Ord: fmt.Sprintf("%012d", seq), Origin: originProposal, ProposedBy: in.Actor, SourceWs: in.SourceWs,
			CreatedAt: now, UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		if err := s.addActivity(ctx, q, node.ID, in.Actor, "proposed",
			map[string]any{"key": key, "source_ws": in.SourceWs}); err != nil {
			return err
		}
		out = mapNode(node)
		return nil
	})
	return out, err
}

// ListInbox — входящие предложки workspace (origin='proposal').
func (s *Service) ListInbox(ctx context.Context, wsIDOrKey string) ([]kb.Node, error) {
	out := []kb.Node{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		nodes, err := q.ListInboxNodes(ctx, ws.ID)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			out = append(out, mapNode(n))
		}
		return nil
	})
	return out, err
}

// AcceptInput — куда в принимающем роадмапе ложится предложка. ParentID (UUID или key, тот же
// workspace; nil => корень роадмапа).
type AcceptInput struct {
	ParentID *string
	Actor    string
}

// AcceptProposal приземляет предложку в роадмап: origin -> 'native' + выбранный parent. Стабильный
// key сохраняется. Не-предложка -> ошибка валидации.
func (s *Service) AcceptProposal(ctx context.Context, idOrKey string, in AcceptInput) (kb.Node, error) {
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		if n.Origin != originProposal {
			return fmt.Errorf("%w: node %s is not a proposal", ErrValidation, n.Key)
		}
		parent, err := s.resolveParentInWorkspace(ctx, q, n.WorkspaceID, in.ParentID)
		if err != nil {
			return err
		}
		updated, err := q.AcceptProposalNode(ctx, store.AcceptProposalNodeParams{
			ParentID: parent, UpdatedAt: s.nowStr(), ID: n.ID,
		})
		if err != nil {
			return err
		}
		if err := s.addActivity(ctx, q, updated.ID, in.Actor, "proposal_accepted",
			map[string]any{"from": n.ProposedBy, "source_ws": n.SourceWs}); err != nil {
			return err
		}
		out = mapNode(updated)
		return nil
	})
	return out, err
}

// DeclineInput — опциональная причина отклонения.
type DeclineInput struct {
	Comment string
	Actor   string
}

// DeclineProposal отклоняет предложку: origin -> 'declined' (вне inbox и дерева, как история).
// Не-предложка -> ошибка валидации.
func (s *Service) DeclineProposal(ctx context.Context, idOrKey string, in DeclineInput) (kb.Node, error) {
	var out kb.Node
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		if n.Origin != originProposal {
			return fmt.Errorf("%w: node %s is not a proposal", ErrValidation, n.Key)
		}
		updated, err := q.DeclineProposalNode(ctx, store.DeclineProposalNodeParams{
			UpdatedAt: s.nowStr(), ID: n.ID,
		})
		if err != nil {
			return err
		}
		if err := s.addActivity(ctx, q, updated.ID, in.Actor, "proposal_declined",
			map[string]any{"comment": in.Comment}); err != nil {
			return err
		}
		out = mapNode(updated)
		return nil
	})
	return out, err
}
