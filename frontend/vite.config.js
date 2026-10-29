import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Vercel services handles /api routing via vercel.json
// No need for local proxy in production
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "dist",
    sourcemap: false,
    minify: "esbuild",
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        ws: true,
        // For local dev: map /api requests to backend server
      },
    },
  },
});
