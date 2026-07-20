import js from "@eslint/js";
import solid from "eslint-plugin-solid/configs/typescript";
import tseslint from "typescript-eslint";

// Flat-config (eslint 9): базовый js + typescript-eslint + solid-правила для TSX. no-undef
// выключаем на TS — типы уже ловят необъявленное, а браузерные globals (window/fetch/document)
// иначе ложно краснеют.
export default tseslint.config(
  { ignores: ["dist/", "node_modules/"] },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ["src/**/*.{ts,tsx}"],
    ...solid,
    rules: {
      ...solid.rules,
      "no-undef": "off",
    },
  },
);
