export const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api'

/** 后端统一响应结构 */
export interface BaseResponse<T = unknown> {
  code: number
  data: T
  msg: string
}

/**
 * 生成随机请求 ID
 */
function generateRequestId(): string {
  const timestamp = Date.now().toString(36)
  const random = Math.random().toString(36).slice(2, 10)
  return `req-${timestamp}-${random}`
}

/**
 * 包装请求体为后端约定的 JSON 信封格式
 *
 * @see docs/restful-api.md 通用约定-统一请求格式
 */
export function envelope<T>(data: T): {
  meta: { request_id: string; timestamp: string }
  data: T
} {
  return {
    meta: {
      request_id: generateRequestId(),
      timestamp: new Date().toISOString(),
    },
    data,
  }
}

export async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('token')
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  const body = await response.json()

  // 兼容两种响应格式:
  // 1. 标准信封: { code: 0, msg: "success", data: { ... } }
  // 2. 直接数据: { challenge_id: "...", challenge: "...", ... } (PoW 等端点)
  if (body && typeof body === 'object' && 'code' in body && 'data' in body) {
    if (!response.ok) {
      throw new Error(body.msg || 'Request failed')
    }
    return body.data as T
  }

  if (!response.ok) {
    throw new Error('Request failed')
  }
  return body as T
}