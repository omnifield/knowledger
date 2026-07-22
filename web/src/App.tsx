import { createEffect, createResource, createSignal, For, onCleanup, Show } from "solid-js";
import type { NodeDTO, TreeNode } from "./api";
import { acceptProposal, declineProposal, getNode, listInbox, listWorkspaces, workspaceTree } from "./api";
import { proposalMeta } from "./format";
import { NodeDetail } from "./NodeDetail";
import { NodeRow } from "./TreeView";
import { TreeProvider } from "./treeCtx";

// --- hash-роутинг: URL = источник правды навигации -------------------------
// Токен в хэше — наш стабильный key: `#/FUND-2` (узел) или `#/FUND` (раздел).
// ws узла выводится из его key (`<WS>-<n>`; ws-ключи без дефисов, так что split
// по последнему `-` однозначен). Deep-link, back/forward и reload работают даром.
interface Route {
  ws?: string;
  key?: string;
}
function parseHash(): Route {
  const t = decodeURIComponent(location.hash.replace(/^#\/?/, "")).trim();
  if (/^[A-Z][A-Z0-9]*-\d+$/.test(t)) return { ws: t.replace(/-\d+$/, ""), key: t };
  if (/^[A-Z][A-Z0-9]*$/.test(t)) return { ws: t };
  return {};
}
function goTo(token: string, replace = false): void {
  const h = `#/${token}`;
  if (replace) location.replace(h);
  else location.hash = h;
}

// App — read-first обзор базы знаний: пространства (workspaces) -> дерево узлов + входящие
// cross-product предложки (accept/decline) -> панель деталей выбранного узла. Навигация — по URL.
export function App() {
  const [workspaces] = createResource(listWorkspaces);
  const [route, setRoute] = createSignal<Route>(parseHash());
  const onHash = () => setRoute(parseHash());
  window.addEventListener("hashchange", onHash);
  onCleanup(() => window.removeEventListener("hashchange", onHash));

  const wsKey = () => route().ws;
  const nodeKey = () => route().key;
  const [actionErr, setActionErr] = createSignal<string | undefined>();

  // Пустой URL -> подставляем первое пространство (replace, чтобы не копить историю).
  createEffect(() => {
    const list = workspaces();
    if (list && list.length > 0 && !route().ws) {
      goTo(list[0].key, true);
      setRoute(parseHash());
    }
  });

  const [tree, { refetch: refetchTree }] = createResource(wsKey, workspaceTree);
  const [inbox, { refetch: refetchInbox }] = createResource(wsKey, listInbox);
  // Выбранный узел — фетчится по key из URL (работает и для deep-link без клика по дереву).
  const [sel] = createResource(nodeKey, getNode);
  // Выбранное пространство целиком — для описания раздела в правой панели, пока узел не выбран.
  const activeWs = () => workspaces()?.find((w) => w.key === wsKey());

  // Смена пространства сбрасывает ошибку действия.
  createEffect(() => {
    wsKey();
    setActionErr(undefined);
  });

  // accept: приземляем предложку корнем; узел уходит из inbox, входит в дерево. Возврат к разделу.
  const onAccept = async (p: NodeDTO) => {
    setActionErr(undefined);
    try {
      await acceptProposal(p.key, {});
      const ws = wsKey();
      if (ws) goTo(ws);
      void refetchInbox();
      void refetchTree();
    } catch (e) {
      setActionErr(String(e));
    }
  };

  // decline: отклоняем с опц. комментом; уходит из inbox. Возврат к разделу.
  const onDecline = async (p: NodeDTO) => {
    setActionErr(undefined);
    const comment = window.prompt(`Отклонить ${p.key}? Причина (опц.):`);
    if (comment === null) {
      return; // отмена диалога
    }
    try {
      await declineProposal(p.key, { comment });
      const ws = wsKey();
      if (ws) goTo(ws);
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
                        classList={{ active: ws.key === wsKey() }}
                        onClick={() => goTo(ws.key)}
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

        <TreeProvider value={{ selectedKey: nodeKey, select: (n: NodeDTO) => goTo(n.key) }}>
          <section class="tree-pane">
            <Show when={wsKey()} fallback={<p class="muted">выберите пространство слева</p>}>
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
                          <button class="proposal-main" onClick={() => goTo(p.key)}>
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

          {/* Правая панель: узел из URL (с состояниями загрузки/ошибки), иначе — описание раздела. */}
          <Show
            when={nodeKey()}
            fallback={
              <Show
                when={activeWs()}
                fallback={
                  <aside class="detail empty">
                    <p class="muted">выберите узел — тело, теги, ссылки, история</p>
                  </aside>
                }
              >
                {(ws) => (
                  <aside class="detail">
                    <div class="detail-head">
                      <span class="key">{ws().key}</span>
                    </div>
                    <h2>{ws().name}</h2>
                    <Show when={ws().description} fallback={<p class="muted">нет описания раздела</p>}>
                      <p class="body">{ws().description}</p>
                    </Show>
                  </aside>
                )}
              </Show>
            }
          >
            <Show when={!sel.loading} fallback={<aside class="detail empty"><p class="muted">загрузка узла…</p></aside>}>
              <Show
                when={!sel.error}
                fallback={
                  <aside class="detail empty">
                    <p class="error">узел {nodeKey()} не найден</p>
                  </aside>
                }
              >
                <Show when={sel()}>{(node) => <NodeDetail node={node()} />}</Show>
              </Show>
            </Show>
          </Show>
        </TreeProvider>
      </main>
    </div>
  );
}
