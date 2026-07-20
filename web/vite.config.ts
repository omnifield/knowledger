import { defineConfig } from "vite";
import solid from "vite-plugin-solid";

// base "/knowledger/" — единый origin через дверь хаба (/api снимается дверью). server.host=true —
// bind 0.0.0.0 (G1: сосед по docker-сети/браузер снаружи достучится). DEV-proxy: локально (без
// двери) фронт бьёт в бэк :8040 напрямую, зеркаля дверь — снимаем /api, нативный префикс
// /knowledger/. В проде единый origin даёт дверь (:8080), proxy не участвует. Пресет
// @omnifield/vite-preset подключим, когда приземлится omnifield.yaml (сейчас — прямой конфиг).
export default defineConfig({
  base: "/knowledger/",
  plugins: [solid()],
  server: {
    host: true,
    proxy: {
      "/api/knowledger": {
        target: "http://localhost:8040",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
    },
  },
});
