import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath } from 'url'
import path from 'path'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    build: {
      lib: {
        entry: path.resolve(__dirname, 'src/widget/index.ts'),
        name: 'weauth',
        fileName: 'weauth',
        formats: ['umd'],
      },
      outDir: 'dist',
      emptyOutDir: false,
      rollupOptions: {
        output: {
          assetFileNames: 'weauth.[ext]',
          entryFileNames: 'weauth.js',
        },
        external: [],
      },
    },
    define: {
      'process.env.NODE_ENV': JSON.stringify('production'),
      'import.meta.env.VITE_WIDGET_BASE_URL': JSON.stringify(
        env.VITE_WIDGET_BASE_URL || 'https://www.weauth.cn'
      ),
      'import.meta.env.VITE_API_BASE_URL': JSON.stringify(
        env.VITE_API_BASE_URL || 'https://api.weauth.cn/api'
      ),
      'import.meta.env.DEV': JSON.stringify(false),
      'import.meta.env.PROD': JSON.stringify(true),
      'import.meta.env.MODE': JSON.stringify('production'),
    },
  }
})