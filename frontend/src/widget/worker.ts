// Web Worker 管理
// 创建 PoW 求解 Worker、发送求解任务、销毁

/** 创建 PoW 求解 Worker */
export function createPowWorker(): Worker {
  return new Worker(
    new URL('./pow-worker.ts', import.meta.url),
    { type: 'module' },
  )
}

/** 求解进度回调 */
export type SolveProgressCallback = (progress: number, iterations: number) => void

/**
 * 在 Worker 中求解指定 challenge 的 PoW。
 * 返回 { nonce, iterations }，或在超时/失败时 reject。
 */
export function solveWithWorker(
  worker: Worker,
  challenge: string,
  difficulty: number,
  onProgress?: SolveProgressCallback,
): Promise<{ nonce: string; iterations: number }> {
  return new Promise((resolve, reject) => {
    const timeout = setTimeout(() => reject(new Error('计算超时')), 120_000)

    worker.onmessage = (e: MessageEvent) => {
      const data = e.data
      if (data.type === 'progress') {
        onProgress?.(data.progress, data.iterations)
      } else if (data.type === 'solved') {
        clearTimeout(timeout)
        resolve({ nonce: data.nonce, iterations: data.iterations })
      } else if (data.type === 'error') {
        clearTimeout(timeout)
        reject(new Error(data.message || '计算失败'))
      }
    }

    worker.onerror = () => {
      clearTimeout(timeout)
      reject(new Error('Worker error'))
    }

    worker.postMessage({ type: 'solve', challenge, difficulty })
  })
}

/** 销毁 Worker */
export function destroyWorker(worker: Worker | null): void {
  worker?.terminate()
}