package store_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/omnifield/knowledger/internal/store"
)

func openTest(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// Миграции применились -> таблицы есть, а per-workspace счётчик монотонен и транзакционен.
func TestMigrationsAndSeqCounter(t *testing.T) {
	ctx := context.Background()
	st := openTest(t)

	ws, err := st.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		ID: "ws1", Key: "KNOW", Name: "n", Description: "", CreatedAt: "t0", UpdatedAt: "t0",
	})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	for want := int64(1); want <= 3; want++ {
		got, err := st.BumpWorkspaceNodeSeq(ctx, store.BumpWorkspaceNodeSeqParams{UpdatedAt: "t1", ID: ws.ID})
		if err != nil {
			t.Fatalf("bump seq: %v", err)
		}
		if got != want {
			t.Fatalf("seq = %d, want %d (counter must be monotonic)", got, want)
		}
	}
}

// dual-id: узел резолвится и по UUID, и по стабильному key.
func TestNodeDualIDResolve(t *testing.T) {
	ctx := context.Background()
	st := openTest(t)
	mustWorkspace(t, st, "ws1", "KNOW")

	n, err := st.CreateNode(ctx, store.CreateNodeParams{
		ID: "node1", WorkspaceID: "ws1", Seq: 1, Key: "KNOW-1",
		ParentID: sql.NullString{}, Kind: sql.NullString{},
		Title: "root", Body: "", Ord: "a0", CreatedAt: "t0", UpdatedAt: "t0",
	})
	if err != nil {
		t.Fatalf("create node: %v", err)
	}

	byID, err := st.GetNodeByID(ctx, n.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	byKey, err := st.GetNodeByKey(ctx, "KNOW-1")
	if err != nil {
		t.Fatalf("get by key: %v", err)
	}
	if byID.ID != byKey.ID {
		t.Fatalf("dual-id resolved different nodes: %s vs %s", byID.ID, byKey.ID)
	}
}

// Backlinks — DERIVED обратным запросом: ребро A->B видно и из ListRefsFrom(A), и из ListBacklinks(B).
func TestRefBacklinksDerived(t *testing.T) {
	ctx := context.Background()
	st := openTest(t)
	mustWorkspace(t, st, "ws1", "KNOW")
	mustNode(t, st, "a", "ws1", 1, "KNOW-1")
	mustNode(t, st, "b", "ws1", 2, "KNOW-2")

	if _, err := st.CreateRef(ctx, store.CreateRefParams{
		ID: "ref1", FromNode: "a", ToNode: "b", Kind: "link", CreatedAt: "t0",
	}); err != nil {
		t.Fatalf("create ref: %v", err)
	}

	from, err := st.ListRefsFrom(ctx, "a")
	if err != nil {
		t.Fatalf("list refs from: %v", err)
	}
	if len(from) != 1 || from[0].ToNode != "b" {
		t.Fatalf("ListRefsFrom(a) = %+v, want one edge to b", from)
	}

	back, err := st.ListBacklinks(ctx, "b")
	if err != nil {
		t.Fatalf("list backlinks: %v", err)
	}
	if len(back) != 1 || back[0].FromNode != "a" {
		t.Fatalf("ListBacklinks(b) = %+v, want one edge from a", back)
	}
}

// Транзакция откатывается при ошибке fn (частичная запись не сохраняется).
func TestTxRollback(t *testing.T) {
	ctx := context.Background()
	st := openTest(t)

	sentinel := context.Canceled
	err := st.Tx(ctx, func(q *store.Queries) error {
		if _, err := q.CreateWorkspace(ctx, store.CreateWorkspaceParams{ID: "ws1", Key: "KNOW", Name: "n", CreatedAt: "t0", UpdatedAt: "t0"}); err != nil {
			return err
		}
		return sentinel // форсим rollback
	})
	if err != sentinel {
		t.Fatalf("Tx err = %v, want sentinel", err)
	}
	if _, err := st.GetWorkspace(ctx, "ws1"); err == nil {
		t.Fatal("workspace persisted despite rollback")
	}
}

func mustWorkspace(t *testing.T, st *store.Store, id, key string) {
	t.Helper()
	if _, err := st.CreateWorkspace(context.Background(), store.CreateWorkspaceParams{
		ID: id, Key: key, Name: "n", CreatedAt: "t0", UpdatedAt: "t0",
	}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
}

func mustNode(t *testing.T, st *store.Store, id, wsID string, seq int64, key string) {
	t.Helper()
	if _, err := st.CreateNode(context.Background(), store.CreateNodeParams{
		ID: id, WorkspaceID: wsID, Seq: seq, Key: key,
		Title: "n", Ord: "a0", CreatedAt: "t0", UpdatedAt: "t0",
	}); err != nil {
		t.Fatalf("create node: %v", err)
	}
}
