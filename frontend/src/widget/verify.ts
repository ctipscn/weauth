// PoW 验证组合式函数
// 封装验证状态机：idle → requesting → solving → verifying → completed (或 error)
// 处理 Worker 的创建/销毁生命周期、单次求解与批量求解

import { onMounted, onUnmounted, ref, type Ref } from 'vue'
import { createChallenge, verifyChallenge, verifyBatchChallenge } from './request'
import { createPowWorker, solveWithWorker, destroyWorker } from './worker'
import type { VerifyStatus } from './types'
import type { ChallengeResponse } from './request'

/** 不可同时运行的状态集合 */
const RUNNING_STATES = new Set<VerifyStatus>(['requesting', 'solving', 'verifying'])

export function useVerify() {
  // --- 响应式状态 ---
  const status: Ref<VerifyStatus> = ref('idle')
  const token: Ref<string | null> = ref(null)
  const errorMessage: Ref<string | null> = ref(null)
  const progress = ref(0)

  // --- Worker 实例（非响应式） ---
  let worker: Worker | null = null

  // --- 生命周期 ---
  onMounted(() => {
    try {
      worker = createPowWorker()
    } catch (err) {
      console.warn('[widget] Failed to create PoW worker:', err)
    }
  })

  onUnmounted(() => {
    destroyWorker(worker)
    worker = null
  })

  // --- 内部方法 ---

  function isRunning(s: VerifyStatus): boolean {
    return RUNNING_STATES.has(s)
  }

  /** 求解单个 challenge */
  async function solveSingle(challenge: {
    challenge_id: string
    challenge: string
    difficulty: number
  }): Promise<string | null> {
    if (!worker) {
      errorMessage.value = 'Worker 不可用'
      status.value = 'error'
      return null
    }

    status.value = 'solving'
    const startTime = Date.now()

    try {
      const solveResult = await solveWithWorker(
        worker,
        challenge.challenge,
        challenge.difficulty,
        (p) => { progress.value = p },
      )

      const solveTimeMs = Date.now() - startTime

      status.value = 'verifying'
      const result = await verifyChallenge(
        challenge.challenge_id,
        solveResult.nonce,
        solveResult.iterations,
        solveTimeMs,
      )

      token.value = result.token
      status.value = 'completed'
      progress.value = 100
      return result.token
    } catch (err) {
      errorMessage.value = err instanceof Error ? err.message : '验证失败'
      status.value = 'error'
      return null
    }
  }

  /** 批量求解多个 challenge 并一次性验证 */
  async function solveBatch(challenge: ChallengeResponse): Promise<string | null> {
    if (!worker || !challenge.challenges || challenge.challenges.length === 0) {
      errorMessage.value = '批量求解失败：数据不足'
      status.value = 'error'
      return null
    }

    status.value = 'solving'
    const startTime = Date.now()

    const nonces: string[] = []
    let totalIterations = 0
    const { challenges } = challenge

    try {
      for (let i = 0; i < challenges.length; i++) {
        const ch = challenges[i]!
        progress.value = Math.round((i / challenges.length) * 100)

        const result = await solveWithWorker(
          worker, ch, challenge.difficulty,
          (p) => {
            progress.value = Math.round(((i + p / 100) / challenges.length) * 100)
          },
        )

        nonces.push(result.nonce)
        totalIterations += result.iterations
      }

      const solveTimeMs = Date.now() - startTime

      status.value = 'verifying'
      const verifyResult = await verifyBatchChallenge(
        challenge.challenge_id,
        nonces,
        totalIterations,
        solveTimeMs,
      )

      token.value = verifyResult.token
      status.value = 'completed'
      progress.value = 100
      return verifyResult.token
    } catch (err) {
      errorMessage.value = err instanceof Error ? err.message : '批量验证失败'
      status.value = 'error'
      return null
    }
  }

  // --- 公开方法 ---

  /**
   * 开始验证流程。如果已有流程正在运行则跳过。
   * @returns 成功时返回 token，失败时返回 null
   */
  async function startVerification(
    sitekey: string,
    origin: string,
    action?: string,
    difficulty?: number,
  ): Promise<string | null> {
    if (isRunning(status.value)) {
      console.warn('[widget] Verification already in progress')
      return null
    }

    status.value = 'requesting'
    errorMessage.value = null
    token.value = null
    progress.value = 0

    try {
      const challenge = await createChallenge(sitekey, origin, action, difficulty)

      if (!challenge) {
        throw new Error('服务器未返回 challenge 数据')
      }

      if (challenge.batch_mode && challenge.challenges && challenge.challenges.length > 0) {
        return await solveBatch(challenge)
      }

      if (challenge.challenge) {
        return await solveSingle({
          challenge_id: challenge.challenge_id,
          challenge: challenge.challenge,
          difficulty: challenge.difficulty ?? difficulty ?? 4,
        })
      }

      throw new Error('无效的 challenge 响应')
    } catch (err) {
      errorMessage.value = err instanceof Error ? err.message : '验证失败'
      status.value = 'error'
      return null
    }
  }

  /** 重置所有状态到 idle */
  function reset(): void {
    status.value = 'idle'
    token.value = null
    errorMessage.value = null
    progress.value = 0
  }

  return {
    status,
    token,
    errorMessage,
    progress,
    startVerification,
    reset,
  }
}

/** isRunning 判断（供外部使用） */
export { RUNNING_STATES }