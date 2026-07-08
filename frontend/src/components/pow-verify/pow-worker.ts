// Web Worker for PoW computation

// Simple SHA-256 implementation for Web Worker
async function sha256(message: string): Promise<string> {
  const msgBuffer = new TextEncoder().encode(message)
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

// Check if hash starts with required number of zeros
function checkDifficulty(hash: string, difficulty: number): boolean {
  const prefix = '0'.repeat(difficulty)
  return hash.startsWith(prefix)
}

// Single PoW result type
interface SinglePoWResult {
  nonce: string
  hash: string
  iterations: number
}

// Main computation for single challenge
async function computeSinglePoW(
  challenge: string, 
  difficulty: number,
  onProgress?: (iterations: number) => void
): Promise<SinglePoWResult> {
  let nonce = 0
  let iterations = 0

  while (true) {
    const data = challenge + nonce.toString()
    const hash = await sha256(data)
    iterations++

    if (checkDifficulty(hash, difficulty)) {
      return {
        nonce: nonce.toString(),
        hash,
        iterations,
      }
    }

    // Report progress periodically
    if (onProgress && iterations % 10000 === 0) {
      onProgress(iterations)
    }

    nonce++

    // Yield to prevent blocking (every 1000 iterations)
    if (iterations % 1000 === 0) {
      await new Promise(resolve => setTimeout(resolve, 0))
    }
  }
}

// Legacy single mode computation (保持向后兼容)
async function computePoW(challenge: string, difficulty: number) {
  const startTime = performance.now()
  
  const result = await computeSinglePoW(challenge, difficulty, (iterations) => {
    const expectedIterations = Math.pow(16, difficulty)
    const progress = Math.min((iterations / expectedIterations) * 100, 99)
    self.postMessage({
      type: 'progress',
      iterations,
      progress,
    })
  })

  const solveTimeMs = Math.round(performance.now() - startTime)
  self.postMessage({
    type: 'result',
    nonce: result.nonce,
    hash: result.hash,
    iterations: result.iterations,
    solveTimeMs,
  })
}

// Batch mode computation (并行计算多个 PoW)
async function computeBatchPoW(challenges: string[], difficulty: number) {
  const startTime = performance.now()
  const totalChallenges = challenges.length
  let completedCount = 0
  let totalIterations = 0

  // 并行计算所有 challenges
  const results = await Promise.all(
    challenges.map(async (challenge) => {
      // 批量模式下不发送单个 challenge 的进度消息，避免进度条乱跳
      const result = await computeSinglePoW(challenge, difficulty)
      
      completedCount++
      totalIterations += result.iterations
      
      // 报告单个 challenge 完成
      self.postMessage({
        type: 'batch_progress',
        completedCount,
        totalChallenges,
        progress: (completedCount / totalChallenges) * 95, // 保留 5% 给提交阶段
      })
      
      return result
    })
  )

  const solveTimeMs = Math.round(performance.now() - startTime)
  
  // 返回批量结果
  self.postMessage({
    type: 'batch_result',
    nonces: results.map(r => r.nonce),
    hashes: results.map(r => r.hash),
    iterations: totalIterations,
    solveTimeMs,
    batchCount: totalChallenges,
  })
}

// Listen for messages
self.onmessage = async (e: MessageEvent) => {
  const { type, challenge, challenges, difficulty } = e.data

  if (type === 'start') {
    // 单个模式
    await computePoW(challenge, difficulty)
  } else if (type === 'start_batch') {
    // 批量模式
    await computeBatchPoW(challenges, difficulty)
  }
}

export {}
