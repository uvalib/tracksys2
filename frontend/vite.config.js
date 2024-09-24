/*global process */

import { fileURLToPath, URL } from 'url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const hash = Math.floor(Math.random() * 90000) + 10000

export default defineConfig({
   define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
   },
   plugins: [vue()],
   resolve: {
      alias: {
         '@': fileURLToPath(new URL('./src', import.meta.url))
      }
   },
   build: {
      rollupOptions: {
         output: {
            entryFileNames: `[name]` + hash + `.js`,
            chunkFileNames: `[name]` + hash + `.js`,
            assetFileNames: `[name]` + hash + `.[ext]`
         }
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
         '/cleanup': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
         '/config': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
         '/healthcheck': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
         '/pdf': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
         '/upload_search_image': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
         '/version': {
            target: process.env.TRACKSYS2_SRV,
            changeOrigin: true
         },
      }
   },
})


