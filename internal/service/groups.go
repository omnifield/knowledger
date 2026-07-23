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

// Группы — организационный слой над workspace (сайдбар-папки). Адресуются по UUID (без dual-id:
// у них нет короткого key). Вложенность через parent_id; членство workspace — через group_id.

// CreateGroupInput — вход создания группы. ParentID nil => корневая группа.
type CreateGroupInput struct {
	Name     string
	ParentID *string
	Ord      string
}

// CreateGroup создаёт группу; при заданном ParentID проверяет, что родитель существует.
func (s *Service) CreateGroup(ctx context.Context, in CreateGroupInput) (kb.Group, error) {
	if strings.TrimSpace(in.Name) == "" {
		return kb.Group{}, fmt.Errorf("%w: name required", ErrValidation)
	}
	var out kb.Group
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		if in.ParentID != nil {
			if _, err := s.getGroup(ctx, q, *in.ParentID); err != nil {
				return err
			}
		}
		now := s.nowStr()
		g, err := q.CreateGroup(ctx, store.CreateGroupParams{
			ID: s.newID(), Name: strings.TrimSpace(in.Name),
			ParentID: ptrToNull(in.ParentID), Ord: in.Ord,
			CreatedAt: now, UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		out = mapGroup(g)
		return nil
	})
	return out, err
}

// ListGroups — все группы (плоско, с parent_id; дерево собирает потребитель).
func (s *Service) ListGroups(ctx context.Context) ([]kb.Group, error) {
	out := []kb.Group{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		gs, err := q.ListGroups(ctx)
		if err != nil {
			return err
		}
		for _, g := range gs {
			out = append(out, mapGroup(g))
		}
		return nil
	})
	return out, err
}

// UpdateGroupInput — частичный PATCH. ParentIDSet отличает «не прислан» от «null» (в корень).
type UpdateGroupInput struct {
	Name        *string
	ParentIDSet bool
	ParentID    *string
	Ord         *string
}

// UpdateGroup применяет частичный PATCH; reparent валидирует существование родителя и цикл.
func (s *Service) UpdateGroup(ctx context.Context, id string, in UpdateGroupInput) (kb.Group, error) {
	var out kb.Group
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		g, err := s.getGroup(ctx, q, id)
		if err != nil {
			return err
		}
		name := g.Name
		if in.Name != nil {
			if strings.TrimSpace(*in.Name) == "" {
				return fmt.Errorf("%w: name cannot be empty", ErrValidation)
			}
			name = strings.TrimSpace(*in.Name)
		}
		parent := g.ParentID
		if in.ParentIDSet {
			np := ptrToNull(in.ParentID)
			if np.Valid {
				if np.String == g.ID {
					return fmt.Errorf("%w: group cannot be its own parent", ErrValidation)
				}
				if _, err := s.getGroup(ctx, q, np.String); err != nil {
					return err
				}
				cyc, err := s.groupCycles(ctx, q, g.ID, np.String)
				if err != nil {
					return err
				}
				if cyc {
					return fmt.Errorf("%w: reparent would create a cycle", ErrValidation)
				}
			}
			parent = np
		}
		ord := g.Ord
		if in.Ord != nil {
			ord = *in.Ord
		}
		updated, err := q.UpdateGroup(ctx, store.UpdateGroupParams{
			Name: name, ParentID: parent, Ord: ord, UpdatedAt: s.nowStr(), ID: g.ID,
		})
		if err != nil {
			return err
		}
		out = mapGroup(updated)
		return nil
	})
	return out, err
}

// DeleteGroup удаляет пустую группу: с дочерними группами или членами-workspace -> ErrConflict
// (перевесить/переназначить сначала).
func (s *Service) DeleteGroup(ctx context.Context, id string) error {
	return s.store.Tx(ctx, func(q *store.Queries) error {
		g, err := s.getGroup(ctx, q, id)
		if err != nil {
			return err
		}
		self := sql.NullString{String: g.ID, Valid: true}
		kids, err := q.CountChildGroups(ctx, self)
		if err != nil {
			return err
		}
		if kids > 0 {
			return fmt.Errorf("%w: group has %d child groups; reparent or delete them first", ErrConflict, kids)
		}
		members, err := q.CountGroupWorkspaces(ctx, self)
		if err != nil {
			return err
		}
		if members > 0 {
			return fmt.Errorf("%w: group has %d workspaces; reassign them first", ErrConflict, members)
		}
		return q.DeleteGroup(ctx, g.ID)
	})
}

// getGroup — по UUID (у групп нет короткого key, только UUID).
func (s *Service) getGroup(ctx context.Context, q *store.Queries, id string) (store.Group, error) {
	g, err := q.GetGroup(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return store.Group{}, fmt.Errorf("%w: group %q", ErrNotFound, id)
	}
	return g, err
}

// groupCycles — true, если newParent лежит в поддереве group (reparent зациклил бы). Идём вверх
// от newParent к корню: встретили groupID -> цикл. Счётчик — страховка от битого parent-кольца.
func (s *Service) groupCycles(ctx context.Context, q *store.Queries, groupID, newParent string) (bool, error) {
	cur := newParent
	for i := 0; i < 10000; i++ {
		if cur == groupID {
			return true, nil
		}
		g, err := q.GetGroup(ctx, cur)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		if !g.ParentID.Valid {
			return false, nil
		}
		cur = g.ParentID.String
	}
	return true, nil
}
