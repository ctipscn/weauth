// Widget 配置
export const WIDGET_BASE_URL = import.meta.env.VITE_WIDGET_BASE_URL || 
  (import.meta.env.DEV 
    ? 'http://localhost:5173' 
    : 'https://www.weauth.cn')

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ||
  (import.meta.env.DEV
    ? 'http://localhost:8080/api'
    : 'https://api.weauth.cn/api')
