import { defineConfig, loadEnv } from 'vite'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    build: {
      lib: {
        entry: path.resolve(__dirname, 'src/widget/index.ts'),
        name: 'weauth',
        fileName: 'weauth',
        formats: ['umd'], // UMD 格式，兼容性最好
      },
      outDir: 'dist',
      emptyOutDir: false, // 不清理 dist，因为主应用也在那里
      rollupOptions: {
        output: {
          assetFileNames: 'weauth.[ext]',
          entryFileNames: 'weauth.js', // 指定输出文件名为 weauth.js
        },
        // 确保外部依赖不会被打包（weauth.js 应该是独立的）
        external: [],
      },
    },
    define: {
      'process.env.NODE_ENV': JSON.stringify('production'),
      // 替换环境变量 - 使用字符串替换方式
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
