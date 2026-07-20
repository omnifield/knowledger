package service_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/omnifield/knowledger/internal/service"
	"github.com/omnifield/knowledger/internal/store"
)

func newSvc(t *testing.T) *service.Service {
	t.Helper()
	st, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "svc.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return service.New(st)
}

// Полный happy-path ядра: workspace -> дерево -> ref -> backlinks (derived) + инварианты.
func TestServiceCoreFlow(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(t)

	ws, err := svc.CreateWorkspace(ctx, service.CreateWorkspaceInput{Key: "know", Name: "KB"})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	if ws.Key != "KNOW" {
		t.Fatalf("ws key = %q, want KNOW (uppercased)", ws.Key)
	}

	root, err := svc.CreateNode(ctx, "KNOW", service.CreateNodeInput{Title: "Go"})
	if err != nil {
		t.Fatalf("create root: %v", err)
	}
	if root.Key != "KNOW-1" {
		t.Fatalf("root key = %q, want KNOW-1", root.Key)
	}

	child, err := svc.CreateNode(ctx, "KNOW", service.CreateNodeInput{Title: "goroutines", ParentID: &root.Key})
	if err != nil {
		t.Fatalf("create child: %v", err)
	}
	if child.ParentID == nil || *child.ParentID != root.ID {
		t.Fatalf("child parent = %v, want %s", child.ParentID, root.ID)
	}

	// dual-id: по key и по uuid резолвится один узел.
	byKey, err := svc.GetNode(ctx, "KNOW-1")
	if err != nil {
		t.Fatalf("get by key: %v", err)
	}
	byID, err := svc.GetNode(ctx, root.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if byKey.ID != byID.ID {
		t.Fatalf("dual-id mismatch: %s vs %s", byKey.ID, byID.ID)
	}

	// ref + DERIVED backlinks: child -link-> root видно в backlinks(root).
	if _, err := svc.AddRef(ctx, child.Key, service.CreateRefInput{ToNode: root.Key, Kind: "link"}); err != nil {
		t.Fatalf("add ref: %v", err)
	}
	back, err := svc.ListBacklinks(ctx, root.Key)
	if err != nil {
		t.Fatalf("backlinks: %v", err)
	}
	if len(back) != 1 || back[0].FromNode != child.ID {
		t.Fatalf("backlinks(root) = %+v, want one from child", back)
	}

	// bad ref kind -> validation.
	if _, err := svc.AddRef(ctx, child.Key, service.CreateRefInput{ToNode: root.Key, Kind: "bogus"}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("bad ref kind err = %v, want ErrValidation", err)
	}
}

// Удаление узла с детьми запрещено (ErrConflict); дубль ws-key — тоже.
func TestServiceInvariants(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(t)

	if _, err := svc.CreateWorkspace(ctx, service.CreateWorkspaceInput{Key: "KNOW", Name: "KB"}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	if _, err := svc.CreateWorkspace(ctx, service.CreateWorkspaceInput{Key: "KNOW", Name: "dup"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("dup ws err = %v, want ErrConflict", err)
	}

	root, _ := svc.CreateNode(ctx, "KNOW", service.CreateNodeInput{Title: "root"})
	if _, err := svc.CreateNode(ctx, "KNOW", service.CreateNodeInput{Title: "child", ParentID: &root.Key}); err != nil {
		t.Fatalf("create child: %v", err)
	}
	if err := svc.DeleteNode(ctx, root.Key); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("delete parent err = %v, want ErrConflict", err)
	}

	// unknown node -> not found.
	if _, err := svc.GetNode(ctx, "KNOW-999"); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("get missing err = %v, want ErrNotFound", err)
	}
}
