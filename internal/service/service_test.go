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

// Cross-product предложка: карантин в inbox (вне дерева) -> accept приземляет в роадмап; decline
// убирает из inbox; accept не-предложки -> валидация.
func TestProposalsFlow(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(t)

	if _, err := svc.CreateWorkspace(ctx, service.CreateWorkspaceInput{Key: "KNOW", Name: "KB"}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	root, _ := svc.CreateNode(ctx, "KNOW", service.CreateNodeInput{Title: "root"})

	prop, err := svc.CreateProposal(ctx, "KNOW", service.CreateProposalInput{
		Title: "article idea", SourceWs: "DEV", Actor: "devopser",
	})
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if prop.Origin != "proposal" || prop.ProposedBy != "devopser" || prop.SourceWs != "DEV" {
		t.Fatalf("proposal provenance wrong: %+v", prop)
	}

	// В inbox — есть; в роадмап-дереве — нет.
	inbox, _ := svc.ListInbox(ctx, "KNOW")
	if len(inbox) != 1 || inbox[0].ID != prop.ID {
		t.Fatalf("inbox = %+v, want the proposal", inbox)
	}
	tree, _ := svc.GetWorkspaceTree(ctx, "KNOW", -1)
	for _, n := range tree {
		if n.ID == prop.ID {
			t.Fatal("proposal leaked into the roadmap tree")
		}
	}

	// Accept -> в роадмап под root, origin native.
	accepted, err := svc.AcceptProposal(ctx, prop.Key, service.AcceptInput{ParentID: &root.Key, Actor: "know"})
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	if accepted.Origin != "native" || accepted.ParentID == nil || *accepted.ParentID != root.ID {
		t.Fatalf("accepted node wrong: %+v", accepted)
	}
	if inbox, _ := svc.ListInbox(ctx, "KNOW"); len(inbox) != 0 {
		t.Fatalf("inbox after accept = %+v, want empty", inbox)
	}

	// Accept не-предложки (обычный узел) -> валидация.
	if _, err := svc.AcceptProposal(ctx, root.Key, service.AcceptInput{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("accept non-proposal err = %v, want ErrValidation", err)
	}

	// Decline убирает из inbox.
	p2, _ := svc.CreateProposal(ctx, "KNOW", service.CreateProposalInput{Title: "reject me", Actor: "x"})
	if _, err := svc.DeclineProposal(ctx, p2.Key, service.DeclineInput{Comment: "nope", Actor: "know"}); err != nil {
		t.Fatalf("decline: %v", err)
	}
	if inbox, _ := svc.ListInbox(ctx, "KNOW"); len(inbox) != 0 {
		t.Fatalf("inbox after decline = %+v, want empty", inbox)
	}
}
