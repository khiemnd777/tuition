import { defineConfig } from "vite";

export default defineConfig({
  base: "./",
  build: {
    outDir: "dist",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          spreadsheet: ["xlsx"],
          archive: ["jszip"],
          qrcode: ["qrcode"],
          sanitizer: ["dompurify"],
        },
      },
    },
  },
  server: {
    port: 5277,
    strictPort: true,
  },
  preview: {
    port: 5278,
    strictPort: true,
  },
  test: {
    environment: "node",
  },
});
