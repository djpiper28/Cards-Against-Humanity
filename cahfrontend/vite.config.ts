import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
// import devtools from 'solid-devtools/vite'

export default defineConfig({
  plugins: [
    solidPlugin(),
    // devtools({
    //   autoname: true,
    // }),
  ],
  server: {
    port: 3000,
    hmr: {
      port: 3001,
    },
  },
  build: {
    target: "esnext",
  },
});
