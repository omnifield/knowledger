// Package service — доменное ядро knowledger: рекурсивное дерево узлов контента, dual-id
// резолв, стабильный per-workspace key, forward-рёбра с DERIVED backlinks, тэги и activity.
// FK/инварианты enforce'им ЗДЕСЬ (Jira-урок: не полагаться на логические FK БД). Возвращает
// доменные типы internal/kb; сырьё берёт из internal/store (sqlc). Зеркалит tasker.
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// Класс-ошибки домена (алиасы kb): httpapi маппит их в статус-коды (404/400/409).
var (
	// ErrNotFound — сущность не найдена.
	ErrNotFound = kb.ErrNotFound
	// ErrValidation — вход нарушил инвариант.
	ErrValidation = kb.ErrValidation
	// ErrConflict — операция конфликтует с текущим состоянием.
	ErrConflict = kb.ErrConflict
)

// Service — фасад над store. now/newID инъектируемы для детерминизма в тестах.
type Service struct {
	store *store.Store
	now   func() time.Time
	newID func() string
}

// New собирает сервис с дефолтным временем (UTC) и UUID-генератором.
func New(st *store.Store) *Service {
	return &Service{
		store: st,
		now:   func() time.Time { return time.Now().UTC() },
		newID: uuid.NewString,
	}
}

// nowStr — таймстемп в RFC3339Nano UTC (лексикографический порядок == хронологический).
func (s *Service) nowStr() string { return s.now().Format(time.RFC3339Nano) }

// --- helpers null<->ptr ----------------------------------------------------

func nullToPtr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	v := n.String
	return &v
}

func ptrToNull(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

// --- dual-id резолв --------------------------------------------------------

// resolveNode находит узел по UUID ИЛИ по стабильному key ("KNOW-12") — оба резолвят один узел.
func (s *Service) resolveNode(ctx context.Context, q *store.Queries, idOrKey string) (store.Node, error) {
	n, err := q.GetNodeByID(ctx, idOrKey)
	if err == nil {
		return n, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return store.Node{}, err
	}
	n, err = q.GetNodeByKey(ctx, idOrKey)
	if errors.Is(err, sql.ErrNoRows) {
		return store.Node{}, fmt.Errorf("%w: node %q", ErrNotFound, idOrKey)
	}
	if err != nil {
		return store.Node{}, err
	}
	return n, nil
}

// resolveWorkspace — по UUID или короткому key.
func (s *Service) resolveWorkspace(ctx context.Context, q *store.Queries, idOrKey string) (store.Workspace, error) {
	w, err := q.GetWorkspace(ctx, idOrKey)
	if err == nil {
		return w, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return store.Workspace{}, err
	}
	w, err = q.GetWorkspaceByKey(ctx, idOrKey)
	if errors.Is(err, sql.ErrNoRows) {
		return store.Workspace{}, fmt.Errorf("%w: workspace %q", ErrNotFound, idOrKey)
	}
	if err != nil {
		return store.Workspace{}, err
	}
	return w, nil
}

// --- мапперы store -> домен (kb) -------------------------------------------

func mapWorkspace(w store.Workspace) kb.Workspace {
	return kb.Workspace{
		ID: w.ID, Key: w.Key, Name: w.Name, Description: w.Description,
		CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
	}
}

func mapNode(n store.Node) kb.Node {
	return kb.Node{
		ID: n.ID, WorkspaceID: n.WorkspaceID, Key: n.Key, Seq: n.Seq,
		ParentID: nullToPtr(n.ParentID), Kind: nullToPtr(n.Kind),
		Title: n.Title, Body: n.Body, Ord: n.Ord,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}
}

func mapRef(r store.Ref) kb.Ref {
	return kb.Ref{
		ID: r.ID, FromNode: r.FromNode, ToNode: r.ToNode, Kind: r.Kind, CreatedAt: r.CreatedAt,
	}
}

func mapTag(t store.Tag) kb.Tag {
	return kb.Tag{ID: t.ID, WorkspaceID: t.WorkspaceID, Name: t.Name, CreatedAt: t.CreatedAt}
}

func mapActivity(a store.Activity) kb.Activity {
	var data json.RawMessage
	if a.Data != "" {
		data = json.RawMessage(a.Data)
	}
	return kb.Activity{
		ID: a.ID, NodeID: a.NodeID, Actor: a.Actor, Kind: a.Kind, Data: data, CreatedAt: a.CreatedAt,
	}
}
