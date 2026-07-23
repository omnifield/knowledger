package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// wsKeyRe — key workspace: короткий префикс node.key ("KNOW" -> "KNOW-1"). Без '-' (иначе ломает
// разбор node.key) и без пробелов; нормализуем в UPPER.
var wsKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)

// CreateWorkspaceInput — вход создания workspace.
type CreateWorkspaceInput struct {
	Key         string
	Name        string
	Description string
	Actor       string
}

// CreateWorkspace создаёт workspace. key уникален, нормализован в UPPER.
func (s *Service) CreateWorkspace(ctx context.Context, in CreateWorkspaceInput) (kb.Workspace, error) {
	key := strings.ToUpper(strings.TrimSpace(in.Key))
	if !wsKeyRe.MatchString(key) {
		return kb.Workspace{}, fmt.Errorf("%w: key must match [A-Z][A-Z0-9]* (no '-', no spaces)", ErrValidation)
	}
	if strings.TrimSpace(in.Name) == "" {
		return kb.Workspace{}, fmt.Errorf("%w: name required", ErrValidation)
	}
	var out kb.Workspace
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		if _, err := q.GetWorkspaceByKey(ctx, key); err == nil {
			return fmt.Errorf("%w: workspace key %q already exists", ErrConflict, key)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		now := s.nowStr()
		ws, err := q.CreateWorkspace(ctx, store.CreateWorkspaceParams{
			ID: s.newID(), Key: key, Name: in.Name, Description: in.Description,
			CreatedAt: now, UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		out = mapWorkspace(ws)
		return nil
	})
	if err != nil {
		return kb.Workspace{}, err
	}
	return out, nil
}

// GetWorkspace — dual-id (UUID или key).
func (s *Service) GetWorkspace(ctx context.Context, idOrKey string) (kb.Workspace, error) {
	var out kb.Workspace
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		w, err := s.resolveWorkspace(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		out = mapWorkspace(w)
		return nil
	})
	return out, err
}

// ListWorkspaces — все workspace.
func (s *Service) ListWorkspaces(ctx context.Context) ([]kb.Workspace, error) {
	out := []kb.Workspace{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		ws, err := q.ListWorkspaces(ctx)
		if err != nil {
			return err
		}
		for _, w := range ws {
			out = append(out, mapWorkspace(w))
		}
		return nil
	})
	return out, err
}

// UpdateWorkspaceInput — частичный PATCH имени/описания/группы. GroupIDSet отличает «поле не
// прислано» от «null» (снять группу -> в корень сайдбара).
type UpdateWorkspaceInput struct {
	Name        *string
	Description *string
	GroupIDSet  bool
	GroupID     *string
}

// UpdateWorkspace применяет частичный PATCH.
func (s *Service) UpdateWorkspace(ctx context.Context, idOrKey string, in UpdateWorkspaceInput) (kb.Workspace, error) {
	var out kb.Workspace
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		w, err := s.resolveWorkspace(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		name, desc := w.Name, w.Description
		if in.Name != nil {
			if strings.TrimSpace(*in.Name) == "" {
				return fmt.Errorf("%w: name cannot be empty", ErrValidation)
			}
			name = *in.Name
		}
		if in.Description != nil {
			desc = *in.Description
		}
		updated, err := q.UpdateWorkspace(ctx, store.UpdateWorkspaceParams{
			Name: name, Description: desc, UpdatedAt: s.nowStr(), ID: w.ID,
		})
		if err != nil {
			return err
		}
		// Назначение/снятие группы — отдельным апдейтом в той же Tx (валидируем существование группы).
		if in.GroupIDSet {
			g := ptrToNull(in.GroupID)
			if g.Valid {
				if _, err := s.getGroup(ctx, q, g.String); err != nil {
					return err
				}
			}
			updated, err = q.SetWorkspaceGroup(ctx, store.SetWorkspaceGroupParams{
				GroupID: g, UpdatedAt: s.nowStr(), ID: w.ID,
			})
			if err != nil {
				return err
			}
		}
		out = mapWorkspace(updated)
		return nil
	})
	return out, err
}
