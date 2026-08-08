import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, 'VITE_')

  const proxyPath  = env.VITE_API_BASE_URL || '/api'
  const proxyTarget = env.VITE_API_TARGET || 'http://127.0.0.1:8989'
  const outDir = "dist"
  const rollupOptions = {
    output: {
      entryFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      chunkFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      assetFileNames: 'assets/087AC4D233B64EB0[name].[hash].[ext]'
    }
  }
  return {
    plugins: [vue()],
    resolve: {
      alias: { '@': path.resolve(__dirname, 'src') }
    },
    server: {
      port: 5173,
      proxy: {
        [proxyPath]: {
          target: proxyTarget,
          changeOrigin: true,
          rewrite: (p) => p.replace(new RegExp('^' + proxyPath), '')
        },
        '/uploads': {
          target: proxyTarget,
          changeOrigin: true,
        },
      }
    },
    build: {
      minify: 'esbuild', // 使用 esbuild 压缩（Vite 内置，无需额外依赖）
      manifest: false, // 是否产出manifest.json
      sourcemap: false, // 是否产出sourcemap.json
      outDir: outDir, // 产出目录
      rollupOptions
    }
  }
})
