package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/omnifield/knowledger/internal/service"
)

// --- workspaces ------------------------------------------------------------

func (a *API) handleListWorkspaces(w http.ResponseWriter, r *http.Request) {
	ws, err := a.svc.ListWorkspaces(r.Context())
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, ws)
}

func (a *API) handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ws, err := a.svc.CreateWorkspace(r.Context(), service.CreateWorkspaceInput{
		Key: body.Key, Name: body.Name, Description: body.Description, Actor: actorOf(r),
	})
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, ws)
}

func (a *API) handleGetWorkspace(w http.ResponseWriter, r *http.Request) {
	ws, err := a.svc.GetWorkspace(r.Context(), r.PathValue("ws"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, ws)
}

func (a *API) handleUpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	var raw map[string]json.RawMessage
	if err := decodeJSON(r, &raw); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var in service.UpdateWorkspaceInput
	if v, ok := raw["name"]; ok {
		s, err := asString(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "name must be a string")
			return
		}
		in.Name = &s
	}
	if v, ok := raw["description"]; ok {
		s, err := asString(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "description must be a string")
			return
		}
		in.Description = &s
	}
	ws, err := a.svc.UpdateWorkspace(r.Context(), r.PathValue("ws"), in)
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, ws)
}

// --- nodes -----------------------------------------------------------------

func (a *API) handleListNodes(w http.ResponseWriter, r *http.Request) {
	var f service.NodeFilter
	q := r.URL.Query()
	// ?parent=<idOrKey|none> — none/root возвращает только корни; отсутствие параметра — все.
	if q.Has("parent") {
		f.ParentSet = true
		if v := q.Get("parent"); v != "" && v != "none" && v != "root" && v != "null" {
			f.ParentID = &v
		}
	}
	// parent-фильтр по key/UUID: резолвим в UUID до in-memory сравнения (совпадает с n.ParentID).
	if f.ParentID != nil {
		n, err := a.svc.GetNode(r.Context(), *f.ParentID)
		if err != nil {
			a.writeServiceError(w, r, err)
			return
		}
		f.ParentID = &n.ID
	}
	nodes, err := a.svc.ListNodes(r.Context(), r.PathValue("ws"), f)
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (a *API) handleCreateNode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title    string  `json:"title"`
		Body     string  `json:"body"`
		Kind     *string `json:"kind"`
		ParentID *string `json:"parent_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	n, err := a.svc.CreateNode(r.Context(), r.PathValue("ws"), service.CreateNodeInput{
		Title: body.Title, Body: body.Body, Kind: body.Kind, ParentID: body.ParentID, Actor: actorOf(r),
	})
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (a *API) handleGetNode(w http.ResponseWriter, r *http.Request) {
	n, err := a.svc.GetNode(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (a *API) handleUpdateNode(w http.ResponseWriter, r *http.Request) {
	// Частичный PATCH: читаем в map, чтобы отличить «поле не прислано» от «null» (очистка).
	var raw map[string]json.RawMessage
	if err := decodeJSON(r, &raw); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in := service.UpdateNodeInput{Actor: actorOf(r)}
	if v, ok := raw["title"]; ok {
		s, err := asString(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "title must be a string")
			return
		}
		in.Title = &s
	}
	if v, ok := raw["body"]; ok {
		s, err := asString(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "body must be a string")
			return
		}
		in.Body = &s
	}
	if v, ok := raw["kind"]; ok {
		in.SetKind = true
		in.Kind = asNullableString(v)
	}
	if v, ok := raw["parent_id"]; ok {
		in.SetParent = true
		in.ParentID = asNullableString(v)
	}
	n, err := a.svc.UpdateNode(r.Context(), r.PathValue("key"), in)
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (a *API) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteNode(r.Context(), r.PathValue("key")); err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleListChildren(w http.ResponseWriter, r *http.Request) {
	kids, err := a.svc.ListChildren(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, kids)
}

// handleGetNodeTree — узел + всё поддерево (subtree-fetch). ?depth=N ограничивает глубину.
func (a *API) handleGetNodeTree(w http.ResponseWriter, r *http.Request) {
	tree, err := a.svc.GetSubtree(r.Context(), r.PathValue("key"), parseDepth(r))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tree)
}

// handleGetWorkspaceTree — корни workspace + их поддеревья. ?depth=N ограничивает глубину.
func (a *API) handleGetWorkspaceTree(w http.ResponseWriter, r *http.Request) {
	roots, err := a.svc.GetWorkspaceTree(r.Context(), r.PathValue("ws"), parseDepth(r))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, roots)
}

// parseDepth читает ?depth=N (>=0 — лимит глубины поддерева); отсутствие/кривое -> -1 (без лимита).
func parseDepth(r *http.Request) int {
	v := r.URL.Query().Get("depth")
	if v == "" {
		return -1
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return -1
	}
	return n
}

// --- refs / backlinks ------------------------------------------------------

func (a *API) handleListRefs(w http.ResponseWriter, r *http.Request) {
	refs, err := a.svc.ListRefs(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, refs)
}

func (a *API) handleCreateRef(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ToNode string `json:"to_node"`
		Kind   string `json:"kind"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ref, err := a.svc.AddRef(r.Context(), r.PathValue("key"), service.CreateRefInput{
		ToNode: body.ToNode, Kind: body.Kind, Actor: actorOf(r),
	})
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, ref)
}

func (a *API) handleListBacklinks(w http.ResponseWriter, r *http.Request) {
	refs, err := a.svc.ListBacklinks(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, refs)
}

func (a *API) handleDeleteRef(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteRef(r.Context(), r.PathValue("id")); err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- activity --------------------------------------------------------------

func (a *API) handleListActivity(w http.ResponseWriter, r *http.Request) {
	acts, err := a.svc.ListActivity(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, acts)
}

func (a *API) handleCreateActivity(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Kind string          `json:"kind"`
		Data json.RawMessage `json:"data"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	act, err := a.svc.AddActivity(r.Context(), r.PathValue("key"), service.AddActivityInput{
		Kind: body.Kind, Data: body.Data, Actor: actorOf(r),
	})
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, act)
}

// --- tags ------------------------------------------------------------------

func (a *API) handleListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := a.svc.ListTags(r.Context(), r.PathValue("ws"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}

func (a *API) handleCreateTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	tag, err := a.svc.CreateTag(r.Context(), r.PathValue("ws"), service.CreateTagInput{Name: body.Name})
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, tag)
}

func (a *API) handleListNodeTags(w http.ResponseWriter, r *http.Request) {
	tags, err := a.svc.ListNodeTags(r.Context(), r.PathValue("key"))
	if err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}

func (a *API) handleAttachTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TagID string `json:"tag_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := a.svc.AddNodeTag(r.Context(), r.PathValue("key"), body.TagID); err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleDetachTag(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.RemoveNodeTag(r.Context(), r.PathValue("key"), r.PathValue("tag_id")); err != nil {
		a.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- JSON-хелперы частичного PATCH -----------------------------------------

// asString извлекает Go-строку из JSON-строки.
func asString(raw json.RawMessage) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// asNullableString: JSON null -> nil (очистка поля); JSON-строка -> *string. Прочее -> nil.
func asNullableString(raw json.RawMessage) *string {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil
	}
	return &s
}
