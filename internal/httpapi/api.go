// Package httpapi — REST-дверь knowledger под нативным префиксом /knowledger/ (дверь omnifield-hub
// rewrite'ит /api/knowledger -> knowledger:8040/knowledger/). Auth — token-stub (Bearer <handle>).
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/omnifield/knowledger/internal/service"
)

// API — HTTP-слой поверх сервис-ядра.
type API struct {
	svc *service.Service
	log *slog.Logger
}

// New собирает API.
func New(svc *service.Service, log *slog.Logger) *API {
	if log == nil {
		log = slog.Default()
	}
	return &API{svc: svc, log: log}
}

type ctxKey int

const actorKey ctxKey = iota

// Handler возвращает роутер со всеми маршрутами /knowledger/* (Go 1.22 method+path patterns).
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()

	// health — без auth (probe девбокса дергает без токена).
	mux.HandleFunc("GET /knowledger/healthz", a.handleHealthz)

	// workspaces
	mux.Handle("GET /knowledger/workspaces", a.auth(a.handleListWorkspaces))
	mux.Handle("POST /knowledger/workspaces", a.auth(a.handleCreateWorkspace))
	mux.Handle("GET /knowledger/workspaces/{ws}", a.auth(a.handleGetWorkspace))
	mux.Handle("PATCH /knowledger/workspaces/{ws}", a.auth(a.handleUpdateWorkspace))
	mux.Handle("GET /knowledger/workspaces/{ws}/nodes", a.auth(a.handleListNodes))
	mux.Handle("POST /knowledger/workspaces/{ws}/nodes", a.auth(a.handleCreateNode))
	mux.Handle("GET /knowledger/workspaces/{ws}/tree", a.auth(a.handleGetWorkspaceTree))
	mux.Handle("GET /knowledger/workspaces/{ws}/tags", a.auth(a.handleListTags))
	mux.Handle("POST /knowledger/workspaces/{ws}/tags", a.auth(a.handleCreateTag))
	// cross-product proposals: write into a foreign workspace's inbox; not the roadmap.
	mux.Handle("POST /knowledger/workspaces/{ws}/proposals", a.auth(a.handleCreateProposal))
	mux.Handle("GET /knowledger/workspaces/{ws}/inbox", a.auth(a.handleListInbox))

	// nodes (dual-id: {key} = UUID или стабильный key)
	mux.Handle("GET /knowledger/nodes/{key}", a.auth(a.handleGetNode))
	mux.Handle("PATCH /knowledger/nodes/{key}", a.auth(a.handleUpdateNode))
	mux.Handle("DELETE /knowledger/nodes/{key}", a.auth(a.handleDeleteNode))
	mux.Handle("GET /knowledger/nodes/{key}/children", a.auth(a.handleListChildren))
	// proposal gate: promote into the roadmap (accept) or reject (decline).
	mux.Handle("POST /knowledger/nodes/{key}/accept", a.auth(a.handleAcceptProposal))
	mux.Handle("POST /knowledger/nodes/{key}/decline", a.auth(a.handleDeclineProposal))
	mux.Handle("GET /knowledger/nodes/{key}/tree", a.auth(a.handleGetNodeTree))
	mux.Handle("GET /knowledger/nodes/{key}/refs", a.auth(a.handleListRefs))
	mux.Handle("POST /knowledger/nodes/{key}/refs", a.auth(a.handleCreateRef))
	mux.Handle("GET /knowledger/nodes/{key}/backlinks", a.auth(a.handleListBacklinks))
	mux.Handle("GET /knowledger/nodes/{key}/activity", a.auth(a.handleListActivity))
	mux.Handle("POST /knowledger/nodes/{key}/activity", a.auth(a.handleCreateActivity))
	mux.Handle("GET /knowledger/nodes/{key}/tags", a.auth(a.handleListNodeTags))
	mux.Handle("POST /knowledger/nodes/{key}/tags", a.auth(a.handleAttachTag))
	mux.Handle("DELETE /knowledger/nodes/{key}/tags/{tag_id}", a.auth(a.handleDetachTag))

	// refs
	mux.Handle("DELETE /knowledger/refs/{id}", a.auth(a.handleDeleteRef))

	return mux
}

func (a *API) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "knowledger"})
}

// auth — token-stub middleware: требует `Authorization: Bearer <handle>`, кладёт handle как actor
// в контекст. Пустой/кривой заголовок -> 401. Полноценная identity — позже (мост).
func (a *API) auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		const p = "Bearer "
		if !strings.HasPrefix(h, p) {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		handle := strings.TrimSpace(strings.TrimPrefix(h, p))
		if handle == "" {
			writeError(w, http.StatusUnauthorized, "empty bearer handle")
			return
		}
		ctx := context.WithValue(r.Context(), actorKey, handle)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func actorOf(r *http.Request) string {
	if v, ok := r.Context().Value(actorKey).(string); ok {
		return v
	}
	return "system"
}

// --- JSON I/O + маппинг ошибок ---------------------------------------------

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeServiceError маппит доменные класс-ошибки в HTTP-статусы.
func (a *API) writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		a.log.ErrorContext(r.Context(), "internal error", slog.String("err", err.Error()),
			slog.String("path", r.URL.Path))
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

// decodeJSON читает тело в v (запрет неизвестных полей -> ловим опечатки клиента).
func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
