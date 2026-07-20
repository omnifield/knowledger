package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// CreateTagInput — вход создания per-workspace тэга.
type CreateTagInput struct {
	Name string
}

// CreateTag создаёт тэг в workspace (name уникален в пределах ws).
func (s *Service) CreateTag(ctx context.Context, wsIDOrKey string, in CreateTagInput) (kb.Tag, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return kb.Tag{}, fmt.Errorf("%w: tag name required", ErrValidation)
	}
	var out kb.Tag
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		tag, err := q.CreateTag(ctx, store.CreateTagParams{
			ID: s.newID(), WorkspaceID: ws.ID, Name: name, CreatedAt: s.nowStr(),
		})
		if err != nil {
			return err
		}
		out = mapTag(tag)
		return nil
	})
	return out, err
}

// ListTags — тэги workspace.
func (s *Service) ListTags(ctx context.Context, wsIDOrKey string) ([]kb.Tag, error) {
	out := []kb.Tag{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := s.resolveWorkspace(ctx, q, wsIDOrKey)
		if err != nil {
			return err
		}
		tags, err := q.ListTagsByWorkspace(ctx, ws.ID)
		if err != nil {
			return err
		}
		for _, t := range tags {
			out = append(out, mapTag(t))
		}
		return nil
	})
	return out, err
}

// ListNodeTags — тэги, навешенные на узел (dual-id).
func (s *Service) ListNodeTags(ctx context.Context, idOrKey string) ([]kb.Tag, error) {
	out := []kb.Tag{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		tags, err := q.ListNodeTags(ctx, n.ID)
		if err != nil {
			return err
		}
		for _, t := range tags {
			out = append(out, mapTag(t))
		}
		return nil
	})
	return out, err
}

// AddNodeTag навешивает тэг на узел (оба обязаны существовать; повторное навешивание -> конфликт).
func (s *Service) AddNodeTag(ctx context.Context, nodeIDOrKey, tagID string) error {
	return s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, nodeIDOrKey)
		if err != nil {
			return err
		}
		if _, err := q.GetTag(ctx, tagID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w: tag %q", ErrNotFound, tagID)
			}
			return err
		}
		existing, err := q.ListNodeTags(ctx, n.ID)
		if err != nil {
			return err
		}
		for _, t := range existing {
			if t.ID == tagID {
				return fmt.Errorf("%w: tag already on node %s", ErrConflict, n.Key)
			}
		}
		return q.AddNodeTag(ctx, store.AddNodeTagParams{NodeID: n.ID, TagID: tagID})
	})
}

// RemoveNodeTag снимает тэг с узла (идемпотентно).
func (s *Service) RemoveNodeTag(ctx context.Context, nodeIDOrKey, tagID string) error {
	return s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, nodeIDOrKey)
		if err != nil {
			return err
		}
		return q.RemoveNodeTag(ctx, store.RemoveNodeTagParams{NodeID: n.ID, TagID: tagID})
	})
}
