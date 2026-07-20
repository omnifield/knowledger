import { defineConfig } from "vitest/config";

// Юнит-тесты чистых хелперов (format.ts) — node-окружение, без DOM/сети.
export default defineConfig({
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
  },
});
