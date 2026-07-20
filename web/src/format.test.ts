import { describe, expect, it } from "vitest";
import type { NodeDTO } from "./api";
import { kindLabel, proposalMeta, shortId } from "./format";

describe("kindLabel", () => {
  it("возвращает kind, если задан", () => {
    expect(kindLabel("article")).toBe("article");
  });
  it("падает на 'node' при null/пусто", () => {
    expect(kindLabel(null)).toBe("node");
    expect(kindLabel("")).toBe("node");
  });
});

describe("shortId", () => {
  it("усекает длинный uuid до 8", () => {
    expect(shortId("550e8400-e29b-41d4")).toBe("550e8400");
  });
  it("не трогает короткие", () => {
    expect(shortId("abc")).toBe("abc");
  });
});

describe("proposalMeta", () => {
  const base: NodeDTO = {
    id: "1", workspace_id: "w", key: "KNOW-1", seq: 1, parent_id: null, kind: null,
    title: "t", body: "", ord: "a0", created_at: "", updated_at: "", origin: "proposal",
  };
  it("клеит автора и source_ws", () => {
    expect(proposalMeta({ ...base, proposed_by: "devopser", source_ws: "DEV" })).toBe("от devopser · DEV");
  });
  it("без source_ws — только автор", () => {
    expect(proposalMeta({ ...base, proposed_by: "x" })).toBe("от x");
  });
});
