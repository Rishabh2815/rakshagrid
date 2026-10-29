import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Production-ready config for Railway + Vercel
const backendURL = process.env.VITE_BACKEND_URL || "http://localhost:8080";

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "dist",
    sourcemap: false,
    minify: "esbuild",
  },
  server: {
    proxy: {
      "/v1": {
        target: backendURL,
        changeOrigin: true,
        ws: true,
      },
    },
  },
});
