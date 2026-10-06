// Widget 配置
// 认证端点、widget 宿主地址等

/** Widget iframe 页面（widget.html）所在的基础地址 */
export const WIDGET_BASE_URL: string =
  import.meta.env.VITE_WIDGET_BASE_URL ||
  (import.meta.env.DEV
    ? 'http://localhost:5173'
    : 'https://www.weauth.cn')

/** 认证 API 基础地址（PoW 挑战、验证等端点） */
export const API_BASE_URL: string =
  import.meta.env.VITE_API_BASE_URL ||
  (import.meta.env.DEV
    ? 'http://localhost:8080'
    : 'https://api.weauth.cn')

/** PoW 认证端点 */
export const ENDPOINTS = {
  /** 创建 PoW 挑战 */
  challenge: '/api/pow/challenge',
  /** 验证 PoW 结果并签发 token */
  verify: '/api/pow/verify',
  /** 服务端校验 token（类似 reCAPTCHA siteverify） */
  siteverify: '/api/pow/siteverify',
} as const

/** Widget 默认尺寸 */
export const DEFAULT_WIDGET_WIDTH = '300px'
export const DEFAULT_WIDGET_HEIGHT = '100px'

// export const WIDGET_BASE_URL: string = 'http://localhost:5173'
// export const API_BASE_URL: string = 'http://localhost:8080'