import { request, envelope } from './client'
import type { ChallengeResponse, SiteVerifyResponse } from './types'

export const powApi = {
  createChallenge: (sitekey: string, origin: string, action?: string, difficulty?: number) =>
    request<ChallengeResponse>('/pow/challenge', {
      method: 'POST',
      body: JSON.stringify(envelope({ sitekey, origin, action, difficulty })),
    }),

  verify: (challengeId: string, nonce: string, iterations?: number, solveTimeMs?: number) =>
    request<{ token: string }>('/pow/verify', {
      method: 'POST',
      body: JSON.stringify(envelope({
        challenge_id: challengeId,
        nonce,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: false,
      })),
    }),

  // 批量模式验证
  verifyBatch: (challengeId: string, nonces: string[], iterations?: number, solveTimeMs?: number) =>
    request<{ token: string }>('/pow/verify', {
      method: 'POST',
      body: JSON.stringify(envelope({
        challenge_id: challengeId,
        nonces,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: true,
      })),
    }),

  // 站点验证：使用 secret + token 验证 PoW 结果（类似 reCAPTCHA siteverify）
  siteverify: (secret: string, token: string) =>
    request<SiteVerifyResponse>('/pow/siteverify', {
      method: 'POST',
      body: JSON.stringify(envelope({ secret, token })),
    }),
}