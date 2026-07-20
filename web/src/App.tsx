import { createEffect, createResource, createSignal, For, Show } from "solid-js";
import type { NodeDTO, TreeNode } from "./api";
import { acceptProposal, declineProposal, listInbox, listWorkspaces, workspaceTree } from "./api";
import { proposalMeta } from "./format";
import { NodeDetail } from "./NodeDetail";
import { NodeRow } from "./TreeView";
import { TreeProvider } from "./treeCtx";

// App — read-first обзор базы знаний: пространства (workspaces) -> дерево узлов + входящие
// cross-product предложки (accept/decline) -> панель деталей выбранного узла.
export function App() {
  const [workspaces] = createResource(listWorkspaces);
  const [selected, setSelected] = createSignal<string | undefined>();
  const [sel, setSel] = createSignal<NodeDTO | undefined>();
  const [actionErr, setActionErr] = createSignal<string | undefined>();

  // Автовыбор первого пространства, как только список приехал.
  createEffect(() => {
    const list = workspaces();
    if (list && list.length > 0 && selected() === undefined) {
      setSelected(list[0].key);
    }
  });

  const [tree, { refetch: refetchTree }] = createResource(selected, workspaceTree);
  const [inbox, { refetch: refetchInbox }] = createResource(selected, listInbox);

  // Смена пространства сбрасывает выбранный узел (он из другого дерева) и ошибку действия.
  createEffect(() => {
    selected();
    setSel(undefined);
    setActionErr(undefined);
  });

  // accept: приземляем предложку в базу корнем; узел уходит из inbox, входит в дерево.
  const onAccept = async (p: NodeDTO) => {
    setActionErr(undefined);
    try {
      await acceptProposal(p.key, {});
      setSel(undefined);
      void refetchInbox();
      void refetchTree();
    } catch (e) {
      setActionErr(String(e));
    }
  };

  // decline: отклоняем с опц. комментом (origin->declined); уходит из inbox.
  const onDecline = async (p: NodeDTO) => {
    setActionErr(undefined);
    const comment = window.prompt(`Отклонить ${p.key}? Причина (опц.):`);
    if (comment === null) {
      return; // отмена диалога
    }
    try {
      await declineProposal(p.key, { comment });
      setSel(undefined);
      void refetchInbox();
    } catch (e) {
      setActionErr(String(e));
    }
  };

  return (
    <div class="app">
      <header class="topbar">
        <h1>knowledger</h1>
        <span class="tagline">база знаний — read-first</span>
      </header>

      <main class="layout">
        <nav class="sidebar">
          <h2>Пространства</h2>
          <Show when={!workspaces.loading} fallback={<p class="muted">загрузка…</p>}>
            <Show when={!workspaces.error} fallback={<p class="error">{String(workspaces.error)}</p>}>
              <ul class="ws-list">
                <For each={workspaces()} fallback={<p class="muted">пусто — создайте через API</p>}>
                  {(ws) => (
                    <li>
                      <button
                        class="ws-item"
                        classList={{ active: ws.key === selected() }}
                        onClick={() => setSelected(ws.key)}
                      >
                        <span class="ws-key">{ws.key}</span>
                        <span class="ws-name">{ws.name}</span>
                      </button>
                    </li>
                  )}
                </For>
              </ul>
            </Show>
          </Show>
        </nav>

        <TreeProvider value={{ selectedKey: () => sel()?.key, select: (n: NodeDTO) => setSel(n) }}>
          <section class="tree-pane">
            <Show when={selected()} fallback={<p class="muted">выберите пространство слева</p>}>
              {/* Входящие — cross-product предложки (origin=proposal): вне базы, с гейтом accept/decline. */}
              <Show when={(inbox()?.length ?? 0) > 0}>
                <div class="inbox">
                  <div class="inbox-head">
                    <h3>
                      Входящие <span class="inbox-count">{inbox()?.length}</span>
                    </h3>
                    <span class="muted">предложки других продуктов — вне базы, пока не приняты</span>
                  </div>
                  <Show when={actionErr()}>
                    <p class="error">{actionErr()}</p>
                  </Show>
                  <ul class="inbox-list">
                    <For each={inbox()}>
                      {(p) => (
                        <li class="proposal">
                          <button class="proposal-main" onClick={() => setSel(p)}>
                            <span class="key">{p.key}</span>
                            <span class="title">{p.title}</span>
                            <span class="proposal-from muted">{proposalMeta(p)}</span>
                          </button>
                          <span class="proposal-actions">
                            <button class="btn accept" onClick={() => onAccept(p)}>
                              Принять
                            </button>
                            <button class="btn decline" onClick={() => onDecline(p)}>
                              Отклонить
                            </button>
                          </span>
                        </li>
                      )}
                    </For>
                  </ul>
                </div>
              </Show>

              <Show when={!tree.loading} fallback={<p class="muted">загрузка дерева…</p>}>
                <Show when={!tree.error} fallback={<p class="error">{String(tree.error)}</p>}>
                  <ul class="tree">
                    <For each={tree()} fallback={<p class="muted">нет узлов — создайте через API</p>}>
                      {(root: TreeNode) => <NodeRow node={root} depth={0} />}
                    </For>
                  </ul>
                </Show>
              </Show>
            </Show>
          </section>

          <Show
            when={sel()}
            fallback={
              <aside class="detail empty">
                <p class="muted">выберите узел — тело, теги, ссылки, история</p>
              </aside>
            }
          >
            {(node) => <NodeDetail node={node()} />}
          </Show>
        </TreeProvider>
      </main>
    </div>
  );
}
