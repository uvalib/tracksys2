 /*global process */

 import { fileURLToPath, URL } from 'url'
 import { defineConfig } from 'vite'
 import vue from '@vitejs/plugin-vue'

 export default defineConfig({
    plugins: [vue()],
    resolve: {
       alias: {
          '@': fileURLToPath(new URL('./src', import.meta.url))
       }
    },
    server: { // this is used in dev mode only
       port: 8080,
       proxy: {
          '/api': {
             target: process.env.TRACKSYS2_SRV,  //export TRACKSYS2_SRV=http://localhost:8085
             changeOrigin: true
          },
          '/authenticate': {
             target: process.env.TRACKSYS2_SRV,
             changeOrigin: true
          },
          '/healthcheck': {
             target: process.env.TRACKSYS2_SRV,
             changeOrigin: true
          },
          '/version': {
             target: process.env.TRACKSYS2_SRV,
             changeOrigin: true
          },
       }
    },
    css: {
      preprocessorOptions: {
        scss: {
           additionalData: `
             @import "@/assets/styles/_mixins.scss";
          `
       },
      },
    },
 })


