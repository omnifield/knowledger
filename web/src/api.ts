// API-клиент knowledger-фронта. Single-origin: зовём backend через ДВЕРЬ-контракт
// `/api/knowledger/…` (дверь снимает `/api` -> нативный `knowledger:8040/knowledger/…`);
// локально это делает vite-proxy. Auth — token-stub `Bearer <handle>` (реальная identity позже).

export const API_BASE = "/api/knowledger";

// token-stub: любой непустой handle проходит auth-middleware бэка.
export const TOKEN = "web";

// --- канон-DTO (зеркало service-слоя бэка, internal/kb) ---------------------

export interface Workspace {
  id: string;
  key: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface NodeDTO {
  id: string;
  workspace_id: string;
  key: string;
  seq: number;
  parent_id: string | null;
  kind: string | null;
  title: string;
  body: string;
  ord: string;
  created_at: string;
  updated_at: string;
  // origin: "native" (узел базы), "proposal" (входящее cross-product предложение в inbox) или
  // "declined". proposed_by/source_ws — провенанс предложки.
  origin: string;
  proposed_by?: string;
  source_ws?: string;
}

// TreeNode — узел + рекурсивно поддерево (subtree-fetch отдаёт всё дерево одним запросом).
export interface TreeNode extends NodeDTO {
  children: TreeNode[];
}

export interface Ref {
  id: string;
  from_node: string;
  to_node: string;
  kind: string;
  created_at: string;
}

export interface Tag {
  id: string;
  workspace_id: string;
  name: string;
  created_at: string;
}

export interface Activity {
  id: string;
  node_id: string;
  actor: string;
  kind: string;
  data?: unknown;
  created_at: string;
}

// --- чистые хелперы (юнит-тестируемы без сети) -----------------------------

// apiUrl клеит путь к door-базе (нормализует ведущий слэш).
export function apiUrl(path: string): string {
  return `${API_BASE}${path.startsWith("/") ? path : `/${path}`}`;
}

// authHeaders — token-stub заголовок для всех запросов.
export function authHeaders(token: string = TOKEN): Record<string, string> {
  return { Authorization: `Bearer ${token}` };
}

// --- транспорт -------------------------------------------------------------

async function get<T>(path: string): Promise<T> {
  const res = await fetch(apiUrl(path), { headers: authHeaders() });
  if (!res.ok) {
    throw new Error(`GET ${path} -> ${res.status} ${res.statusText}`);
  }
  return (await res.json()) as T;
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(apiUrl(path), {
    method: "POST",
    headers: { ...authHeaders(), "Content-Type": "application/json" },
    body: JSON.stringify(body ?? {}),
  });
  if (!res.ok) {
    let detail = res.statusText;
    try {
      const j = (await res.json()) as { error?: string };
      if (j.error) detail = j.error;
    } catch {
      /* тело не JSON — оставляем statusText */
    }
    throw new Error(`POST ${path} -> ${res.status} ${detail}`);
  }
  return (await res.json()) as T;
}

export const listWorkspaces = () => get<Workspace[]>("/workspaces");
export const workspaceTree = (ws: string) => get<TreeNode[]>(`/workspaces/${ws}/tree`);
export const nodeRefs = (key: string) => get<Ref[]>(`/nodes/${key}/refs`);
export const nodeBacklinks = (key: string) => get<Ref[]>(`/nodes/${key}/backlinks`);
export const nodeTags = (key: string) => get<Tag[]>(`/nodes/${key}/tags`);
export const nodeActivity = (key: string) => get<Activity[]>(`/nodes/${key}/activity`);

// --- cross-product предложки (inbox + гейт) --------------------------------

export const listInbox = (ws: string) => get<NodeDTO[]>(`/workspaces/${ws}/inbox`);

// acceptProposal приземляет предложку в базу (origin->native); parent опц.
export const acceptProposal = (key: string, body: { parent_id?: string | null }) =>
  post<NodeDTO>(`/nodes/${key}/accept`, body);

// declineProposal отклоняет предложку (origin->declined); опц. коммент.
export const declineProposal = (key: string, body: { comment?: string }) =>
  post<NodeDTO>(`/nodes/${key}/decline`, body);
