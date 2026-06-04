import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  base: "/PulseBlog/",
  plugins: [vue()],
  server: {
    port: 5173,
    strictPort: false
  }
});
