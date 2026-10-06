// 向认证端点发出请求的 interface 与函数
import { API_BASE_URL, ENDPOINTS } from './config'

/* ==================== 接口定义 ==================== */

/** 后端统一响应结构 */
export interface BaseResponse<T = unknown> {
  code: number
  data: T
  msg: string
}

/** PoW 挑战响应 */
export interface ChallengeResponse {
  challenge_id: string
  /** 单次模式下的挑战字符串 */
  challenge?: string
  /** 批量模式下的多个挑战字符串 */
  challenges?: string[]
  difficulty: number
  batch_mode: boolean
  batch_total: number
}

/** PoW 验证响应 */
export interface VerifyResponse {
  token: string
}

/** 站点验证响应（类似 reCAPTCHA siteverify） */
export interface SiteVerifyResponse {
  hostname: string
  action: string
  challenge_ts: string
}

/* ==================== 内部工具 ==================== */

/** 生成随机请求 ID */
function generateRequestId(): string {
  const timestamp = Date.now().toString(36)
  const random = Math.random().toString(36).slice(2, 10)
  return `req-${timestamp}-${random}`
}

/** 包装请求体为后端约定的 JSON 信封格式 */
function envelope<T>(data: T): {
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

/** 通用请求函数 */
async function rawRequest<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  const body = await response.json()

  // 兼容两种响应格式:
  // 1. 标准信封: { code: 0, msg: "success", data: { ... } }
  // 2. 直接数据: { challenge_id: "...", challenge: "...", ... }
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

/* ==================== 公开 API ==================== */

/** 创建 PoW 挑战 */
export function createChallenge(
  sitekey: string,
  origin: string,
  action?: string,
  difficulty?: number,
): Promise<ChallengeResponse> {
  return rawRequest(ENDPOINTS.challenge, {
    method: 'POST',
    body: JSON.stringify(envelope({ sitekey, origin, action, difficulty })),
  })
}

/** 验证 PoW 结果（单次模式） */
export function verifyChallenge(
  challengeId: string,
  nonce: string,
  iterations?: number,
  solveTimeMs?: number,
): Promise<VerifyResponse> {
  return rawRequest(ENDPOINTS.verify, {
    method: 'POST',
    body: JSON.stringify(
      envelope({
        challenge_id: challengeId,
        nonce,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: false,
      }),
    ),
  })
}

/** 验证 PoW 结果（批量模式） */
export function verifyBatchChallenge(
  challengeId: string,
  nonces: string[],
  iterations?: number,
  solveTimeMs?: number,
): Promise<VerifyResponse> {
  return rawRequest(ENDPOINTS.verify, {
    method: 'POST',
    body: JSON.stringify(
      envelope({
        challenge_id: challengeId,
        nonces,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: true,
      }),
    ),
  })
}

/** 服务端校验 token（类似 reCAPTCHA siteverify） */
export function siteVerify(
  secret: string,
  token: string,
): Promise<SiteVerifyResponse> {
  return rawRequest(ENDPOINTS.siteverify, {
    method: 'POST',
    body: JSON.stringify(envelope({ secret, token })),
  })
}