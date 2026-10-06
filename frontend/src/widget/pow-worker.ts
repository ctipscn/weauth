interface SolveMessage {
  type: 'solve'
  challenge: string
  difficulty: number
}

interface StopMessage {
  type: 'stop'
}

type WorkerMessage = SolveMessage | StopMessage

interface SolvedResult {
  type: 'solved'
  nonce: string
  iterations: number
}

interface ProgressResult {
  type: 'progress'
  progress: number
  iterations: number
}

type WorkerResult = SolvedResult | ProgressResult

let running = true

self.onmessage = (e: MessageEvent<WorkerMessage>) => {
  const msg = e.data

  if (msg.type === 'stop') {
    running = false
    return
  }

  if (msg.type === 'solve') {
    running = true
    const { challenge, difficulty } = msg
    // 启动异步求解，不阻塞后续 stop 消息
    solve(challenge, difficulty).then((result) => {
      if (result && running) {
        self.postMessage({ type: 'solved', ...result } satisfies SolvedResult)
      }
    })
  }
}

async function solve(
  challenge: string,
  difficulty: number,
): Promise<{ nonce: string; iterations: number } | null> {
  const target = '0'.repeat(difficulty)
  let nonce = 0
  const batchSize = 100
  const reportInterval = 5000
  const encoder = new TextEncoder()

  while (running) {
    for (let i = 0; i < batchSize; i++) {
      const data = encoder.encode(challenge + nonce.toString(16))
      const hashHex = await sha256Hex(data)

      if (hashHex.startsWith(target)) {
        return { nonce: nonce.toString(16), iterations: nonce + 1 }
      }
      nonce++
    }

    // 报告进度（含实时迭代次数）
    if (nonce % reportInterval === 0) {
      self.postMessage({
        type: 'progress',
        progress: Math.min((nonce / 100000) * 100, 99),
        iterations: nonce,
      } satisfies ProgressResult)
    }
  }

  return null
}

/**
 * 使用 Web Crypto API 计算 SHA-256 并返回十六进制字符串。
 * 在 Web Worker 中 crypto.subtle 完全可用。
 */
async function sha256Hex(data: Uint8Array): Promise<string> {
  const hashBuffer = await crypto.subtle.digest('SHA-256', data as BufferSource)
  const hashArray = new Uint8Array(hashBuffer)
  // 将每个字节转为两位十六进制字符串
  let hex = ''
  for (let i = 0; i < hashArray.length; i++) {
    hex += hashArray[i]!.toString(16).padStart(2, '0')
  }
  return hex
}

export {}