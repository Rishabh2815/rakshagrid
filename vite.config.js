import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Dev server proxies /v1 to the Go backend so the PWA can be developed
// against a local `docker compose up` without CORS configuration.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/v1": {
        target: "http://localhost:8080",
        changeOrigin: true,
        ws: true,
      },
    },
  },
});
