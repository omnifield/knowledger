import { createResource, For, Show } from "solid-js";
import type { NodeDTO } from "./api";
import { getNode, nodeActivity, nodeBacklinks, nodeRefs, nodeTags } from "./api";
import { kindLabel, shortId } from "./format";

// RefTarget — резолвит UUID ребра в узел (dual-id) и рендерит кликабельную ссылку `#/<KEY>`
// с key+title. Пока грузится — усечённый UUID как заглушка. Так блок «Ссылки» согласован
// с URL-навигацией (key, не сырой айдишник) и по ссылкам можно переходить.
function RefTarget(props: { id: string }) {
  const [node] = createResource(() => props.id, getNode);
  return (
    <Show when={node()} fallback={<span class="mono">{shortId(props.id)}</span>}>
      {(n) => (
        <a class="ref-node" href={`#/${n().key}`}>
          <span class="key">{n().key}</span> {n().title}
        </a>
      )}
    </Show>
  );
}

// NodeDetail — панель выбранного узла: тело, теги, исходящие ссылки, обратные ссылки (derived) и
// история. Каждый список — свой resource, перезагружается при смене выбранного узла (по key).
export function NodeDetail(props: { node: NodeDTO }) {
  const key = () => props.node.key;
  const [refs] = createResource(key, nodeRefs);
  const [backlinks] = createResource(key, nodeBacklinks);
  const [tags] = createResource(key, nodeTags);
  const [activity] = createResource(key, nodeActivity);

  return (
    <aside class="detail">
      <div class="detail-head">
        <span class="key">{props.node.key}</span>
        <Show when={props.node.kind}>
          <span class="kind-badge">{kindLabel(props.node.kind)}</span>
        </Show>
        <Show when={props.node.origin !== "native"}>
          <span class="origin-badge">{props.node.origin}</span>
        </Show>
      </div>
      <h2>{props.node.title}</h2>
      <Show when={props.node.body} fallback={<p class="muted">нет тела</p>}>
        <p class="body">{props.node.body}</p>
      </Show>

      <section>
        <h3>Теги</h3>
        <Show when={(tags()?.length ?? 0) > 0} fallback={<p class="muted">нет</p>}>
          <div class="tag-row">
            <For each={tags()}>{(t) => <span class="tag">{t.name}</span>}</For>
          </div>
        </Show>
      </section>

      <section>
        <h3>Ссылки</h3>
        <Show when={(refs()?.length ?? 0) > 0} fallback={<p class="muted">нет</p>}>
          <ul class="ref-list">
            <For each={refs()}>
              {(r) => (
                <li>
                  <span class="ref-kind">{r.kind}</span> → <RefTarget id={r.to_node} />
                </li>
              )}
            </For>
          </ul>
        </Show>
      </section>

      <section>
        <h3>
          Обратные ссылки <span class="derived">derived</span>
        </h3>
        <Show when={(backlinks()?.length ?? 0) > 0} fallback={<p class="muted">нет</p>}>
          <ul class="ref-list">
            <For each={backlinks()}>
              {(r) => (
                <li>
                  <RefTarget id={r.from_node} /> <span class="ref-kind">{r.kind}</span> →
                </li>
              )}
            </For>
          </ul>
        </Show>
      </section>

      <section>
        <h3>История</h3>
        <Show when={(activity()?.length ?? 0) > 0} fallback={<p class="muted">нет</p>}>
          <ul class="activity">
            <For each={activity()}>
              {(a) => (
                <li>
                  <span class="act-kind">{a.kind}</span> <span class="muted">{a.actor}</span>
                </li>
              )}
            </For>
          </ul>
        </Show>
      </section>
    </aside>
  );
}
