import type { NodeDTO } from "./api";

// kindLabel — отображаемый вид node.kind (type-less: пусто -> "node").
export function kindLabel(kind: string | null): string {
  return kind && kind.trim() !== "" ? kind : "node";
}

// shortId — усечённый UUID для компактного показа ссылок.
export function shortId(id: string): string {
  return id.length > 8 ? id.slice(0, 8) : id;
}

// proposalMeta — провенанс входящей предложки ("от <кто> · <ws>").
export function proposalMeta(p: NodeDTO): string {
  const from = p.proposed_by && p.proposed_by !== "" ? p.proposed_by : "?";
  const src = p.source_ws && p.source_ws !== "" ? ` · ${p.source_ws}` : "";
  return `от ${from}${src}`;
}
