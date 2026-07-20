import { createContext, useContext } from "solid-js";
import type { NodeDTO } from "./api";

// TreeCtx — выбор узла для панели деталей (общий для дерева и inbox).
export interface TreeCtx {
  selectedKey: () => string | undefined;
  select: (n: NodeDTO) => void;
}

const Ctx = createContext<TreeCtx>();

// TreeProvider — провайдер контекста выбора.
export const TreeProvider = Ctx.Provider;

// useTree — доступ к контексту выбора (бросает вне провайдера).
export function useTree(): TreeCtx {
  const c = useContext(Ctx);
  if (!c) {
    throw new Error("useTree вне TreeProvider");
  }
  return c;
}
