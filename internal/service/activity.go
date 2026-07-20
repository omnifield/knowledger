package service

import (
	"context"
	"encoding/json"

	"github.com/omnifield/knowledger/internal/kb"
	"github.com/omnifield/knowledger/internal/store"
)

// AddActivityInput — вход добавления записи timeline (комментарий/событие).
type AddActivityInput struct {
	Kind  string
	Data  json.RawMessage
	Actor string
}

// AddActivity добавляет запись в timeline узла (dual-id).
func (s *Service) AddActivity(ctx context.Context, idOrKey string, in AddActivityInput) (kb.Activity, error) {
	var out kb.Activity
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		actor := in.Actor
		if actor == "" {
			actor = "system"
		}
		data := ""
		if len(in.Data) > 0 {
			data = string(in.Data)
		}
		act, err := q.CreateActivity(ctx, store.CreateActivityParams{
			ID: s.newID(), NodeID: n.ID, Actor: actor, Kind: in.Kind, Data: data, CreatedAt: s.nowStr(),
		})
		if err != nil {
			return err
		}
		out = mapActivity(act)
		return nil
	})
	return out, err
}

// ListActivity — timeline узла (dual-id), в хронологическом порядке.
func (s *Service) ListActivity(ctx context.Context, idOrKey string) ([]kb.Activity, error) {
	out := []kb.Activity{}
	err := s.store.Tx(ctx, func(q *store.Queries) error {
		n, err := s.resolveNode(ctx, q, idOrKey)
		if err != nil {
			return err
		}
		acts, err := q.ListActivity(ctx, n.ID)
		if err != nil {
			return err
		}
		for _, a := range acts {
			out = append(out, mapActivity(a))
		}
		return nil
	})
	return out, err
}
