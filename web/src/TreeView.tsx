import { createSignal, For, Show } from "solid-js";
import type { TreeNode } from "./api";
import { kindLabel } from "./format";
import { useTree } from "./treeCtx";

// NodeRow — рекурсивная строка дерева базы знаний (read-first): раскрытие детей, key, заголовок,
// kind-бейдж. Клик по строке -> выбор узла (панель деталей справа).
export function NodeRow(props: { node: TreeNode; depth: number }) {
  const tree = useTree();
  // depth статичен per-инстанс (позиция в дереве не меняется) -> безопасно как начальное значение.
  // eslint-disable-next-line solid/reactivity
  const [open, setOpen] = createSignal(props.depth < 2);
  const hasChildren = () => props.node.children.length > 0;
  const selected = () => tree.selectedKey() === props.node.key;

  return (
    <li class="node">
      <div
        class="node-row"
        classList={{ selected: selected() }}
        style={{ "padding-left": `${props.depth * 1.1}rem` }}
      >
        <button
          class="twisty"
          classList={{ leaf: !hasChildren() }}
          disabled={!hasChildren()}
          onClick={() => setOpen(!open())}
        >
          {hasChildren() ? (open() ? "▾" : "▸") : "•"}
        </button>
        <button class="node-main" onClick={() => tree.select(props.node)}>
          <span class="key">{props.node.key}</span>
          <span class="title">{props.node.title}</span>
          <Show when={props.node.kind}>
            <span class="kind-badge">{kindLabel(props.node.kind)}</span>
          </Show>
        </button>
      </div>
      <Show when={hasChildren() && open()}>
        <ul class="tree">
          <For each={props.node.children}>{(c) => <NodeRow node={c} depth={props.depth + 1} />}</For>
        </ul>
      </Show>
    </li>
  );
}
