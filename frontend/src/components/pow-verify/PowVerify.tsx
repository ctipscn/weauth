import { useState, useEffect, useRef, useCallback, useImperativeHandle, forwardRef } from 'react'
import { cn } from '@/lib/utils'
import { powApi } from '@/lib/api'
import type { PowVerifyProps, VerifyStatus, ChallengeData, WorkerResult, BatchWorkerResult } from './types'
import { Check, X, RefreshCw, Loader2 } from 'lucide-react'

// Token 缓存 key 前缀
const TOKEN_CACHE_PREFIX = 'weauth_token_'
// 默认缓存时间 35 秒（服务端 token 有效期 60 秒）
const DEFAULT_CACHE_TTL = 35 * 1000

// 验证颜色格式（仅允许 #RRGGBB 或 #RGB 格式，防止 CSS 注入）
function isValidHexColor(color: string): boolean {
  return /^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$/.test(color)
}

// 将短格式转换为长格式 (#RGB -> #RRGGBB)
function normalizeHexColor(color: string): string {
  if (color.length === 4) {
    return '#' + color[1] + color[1] + color[2] + color[2] + color[3] + color[3]
  }
  return color
}

// 根据颜色亮度判断是否为深色（保留供未来使用）
// function isColorDark(hexColor: string): boolean {
//   const hex = normalizeHexColor(hexColor).substring(1)
//   const r = parseInt(hex.substring(0, 2), 16)
//   const g = parseInt(hex.substring(2, 4), 16)
//   const b = parseInt(hex.substring(4, 6), 16)
//   // 使用相对亮度公式
//   const brightness = (r * 299 + g * 587 + b * 114) / 1000
//   return brightness < 128
// }

// 调整颜色亮度（用于生成深色主题的淡色变体）
function adjustColorBrightness(hexColor: string, percent: number): string {
  const hex = normalizeHexColor(hexColor).substring(1)
  const r = Math.min(255, Math.max(0, parseInt(hex.substring(0, 2), 16) + percent))
  const g = Math.min(255, Math.max(0, parseInt(hex.substring(2, 4), 16) + percent))
  const b = Math.min(255, Math.max(0, parseInt(hex.substring(4, 6), 16) + percent))
  
  return '#' + [r, g, b].map(x => {
    const hex = x.toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }).join('')
}

// 获取缓存的 token 和过期时间
function getCachedToken(sitekey: string, action?: string): { token: string; expiresAt: number } | null {
  const key = TOKEN_CACHE_PREFIX + sitekey + (action ? '_' + action : '')
  const cached = sessionStorage.getItem(key)
  if (!cached) return null
  
  try {
    const { token, expiresAt } = JSON.parse(cached)
    if (Date.now() < expiresAt) {
      return { token, expiresAt }
    }
    // 已过期，清除
    sessionStorage.removeItem(key)
  } catch {
    sessionStorage.removeItem(key)
  }
  return null
}

// 缓存 token（默认 35 秒有效）
function cacheToken(sitekey: string, action: string | undefined, token: string, ttlMs = DEFAULT_CACHE_TTL) {
  const key = TOKEN_CACHE_PREFIX + sitekey + (action ? '_' + action : '')
  const data = {
    token,
    expiresAt: Date.now() + ttlMs
  }
  sessionStorage.setItem(key, JSON.stringify(data))
}

// 清除缓存的 token
function clearCachedToken(sitekey: string, action?: string) {
  const key = TOKEN_CACHE_PREFIX + sitekey + (action ? '_' + action : '')
  sessionStorage.removeItem(key)
}

// PowVerify 组件暴露的方法
export interface PowVerifyRef {
  reset: () => void
}

export const PowVerify = forwardRef<PowVerifyRef, PowVerifyProps>(({
  sitekey,
  action,
  difficulty: overrideDifficulty,
  onSuccess,
  onError,
  onExpire, // Token 过期时回调
  theme = 'auto',
  themeColor = 'blue', // 默认蓝色主题
  customColor, // 自定义颜色
  size = 'normal',
  autoStart = false, // 是否自动开始验证
}, ref) => {
  const [status, setStatus] = useState<VerifyStatus>('idle')
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [batchInfo, setBatchInfo] = useState<{ current: number; total: number } | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const workerRef = useRef<Worker | null>(null)
  const challengeRef = useRef<ChallengeData | null>(null)
  const expiryTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // 清理定时器
  const clearTimers = useCallback(() => {
    if (expiryTimerRef.current) {
      clearTimeout(expiryTimerRef.current)
      expiryTimerRef.current = null
    }
  }, [])

  // 设置过期定时器
  const setupExpiryTimer = useCallback((expiresAt: number) => {
    clearTimers()
    
    const now = Date.now()
    const timeUntilExpiry = expiresAt - now
    
    if (timeUntilExpiry <= 0) {
      // 已过期
      setToken(null)
      setStatus('idle')
      clearCachedToken(sitekey, action)
      onExpire?.()
      return
    }

    // 设置过期后自动恢复待验证状态
    expiryTimerRef.current = setTimeout(() => {
      setToken(null)
      setStatus('idle')
      clearCachedToken(sitekey, action)
      onExpire?.()
    }, timeUntilExpiry)
  }, [sitekey, action, clearTimers, onExpire])

  // 重置组件状态（类似 Turnstile 的 reset() 方法）
  const reset = useCallback(() => {
    // 清除 worker
    if (workerRef.current) {
      workerRef.current.terminate()
      workerRef.current = null
    }
    
    // 清除定时器
    clearTimers()
    
    // 清除缓存的 token
    clearCachedToken(sitekey, action)
    
    // 重置状态
    setToken(null)
    setStatus('idle')
    setProgress(0)
    setError(null)
    setBatchInfo(null)
    challengeRef.current = null
    
    // 通知父组件
    onExpire?.()
  }, [sitekey, action, clearTimers, onExpire])

  // 暴露 reset 方法给父组件
  useImperativeHandle(ref, () => ({
    reset
  }), [reset])

  // 检查缓存的 token
  useEffect(() => {
    const cached = getCachedToken(sitekey, action)
    if (cached) {
      setToken(cached.token)
      setStatus('success')
      onSuccess(cached.token)
      setupExpiryTimer(cached.expiresAt)
    }
    
    return () => clearTimers()
  }, [sitekey, action, onSuccess, setupExpiryTimer, clearTimers])

  // Cleanup worker on unmount
  useEffect(() => {
    return () => {
      if (workerRef.current) {
        workerRef.current.terminate()
      }
    }
  }, [])

  // Start verification process
  const startVerification = useCallback(async () => {
    // 如果已经有有效 token，直接返回
    const cached = getCachedToken(sitekey, action)
    if (cached) {
      setToken(cached.token)
      setStatus('success')
      onSuccess(cached.token)
      setupExpiryTimer(cached.expiresAt)
      return
    }

    try {
      setStatus('loading')
      setProgress(0)
      setError(null)
      setBatchInfo(null)

      // Get challenge from server
      const origin = window.location.origin
      const challenge = await powApi.createChallenge(sitekey, origin, action, overrideDifficulty)
      challengeRef.current = challenge
      
      // 如果是批量模式，显示批次信息
      if (challenge.batch_mode && challenge.batch_total > 1) {
        setBatchInfo({ current: 0, total: challenge.batch_total })
      }
      
      setStatus('computing')

      // Create Web Worker
      const worker = new Worker(
        new URL('./pow-worker.ts', import.meta.url),
        { type: 'module' }
      )
      workerRef.current = worker

      // Handle worker messages
      worker.onmessage = async (e) => {
        const data = e.data

        if (data.type === 'progress') {
          // 单个模式进度（批量模式下忽略，避免进度条乱跳）
          if (!batchInfo) {
            setProgress(Math.min(data.progress, 95))
          }
        } else if (data.type === 'batch_progress') {
          // 批量模式进度
          setBatchInfo({ current: data.completedCount, total: data.totalChallenges })
          setProgress(Math.min(data.progress, 95))
        } else if (data.type === 'result') {
          // 单个模式结果
          const result: WorkerResult = {
            nonce: data.nonce,
            hash: data.hash,
            iterations: data.iterations,
            solveTimeMs: data.solveTimeMs,
          }
          
          setProgress(98)
          setStatus('verifying')

          try {
            // Verify with server
            const verifyResult = await powApi.verify(
              challengeRef.current!.challenge_id,
              result.nonce,
              result.iterations,
              result.solveTimeMs
            )

            setProgress(100)
            setStatus('success')
            setToken(verifyResult.token)
            // 缓存 token 并设置过期定时器
            cacheToken(sitekey, action, verifyResult.token)
            const expiresAt = Date.now() + DEFAULT_CACHE_TTL
            setupExpiryTimer(expiresAt)
            onSuccess(verifyResult.token)
          } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Verification failed'
            setError(errorMessage)
            setStatus('error')
            onError?.(err instanceof Error ? err : new Error(errorMessage))
          }

          worker.terminate()
          workerRef.current = null
        } else if (data.type === 'batch_result') {
          // 批量模式结果
          const result: BatchWorkerResult = {
            nonces: data.nonces,
            hashes: data.hashes,
            iterations: data.iterations,
            solveTimeMs: data.solveTimeMs,
            batchCount: data.batchCount,
          }
          
          setProgress(98)
          setStatus('verifying')

          try {
            // Batch verify with server
            const verifyResult = await powApi.verifyBatch(
              challengeRef.current!.challenge_id,
              result.nonces,
              result.iterations,
              result.solveTimeMs
            )

            setProgress(100)
            setStatus('success')
            setToken(verifyResult.token)
            // 缓存 token 并设置过期定时器
            cacheToken(sitekey, action, verifyResult.token)
            const expiresAt = Date.now() + DEFAULT_CACHE_TTL
            setupExpiryTimer(expiresAt)
            onSuccess(verifyResult.token)
          } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Verification failed'
            setError(errorMessage)
            setStatus('error')
            onError?.(err instanceof Error ? err : new Error(errorMessage))
          }

          worker.terminate()
          workerRef.current = null
        }
      }

      worker.onerror = (err) => {
        setError('Worker error: ' + err.message)
        setStatus('error')
        onError?.(new Error(err.message))
        worker.terminate()
        workerRef.current = null
      }

      // Start computation based on mode
      if (challenge.batch_mode && challenge.challenges) {
        // 批量模式
        worker.postMessage({
          type: 'start_batch',
          challenges: challenge.challenges,
          difficulty: challenge.difficulty,
        })
      } else {
        // 单个模式
        worker.postMessage({
          type: 'start',
          challenge: challenge.challenge,
          difficulty: challenge.difficulty,
        })
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to start verification'
      setError(errorMessage)
      setStatus('error')
      onError?.(err instanceof Error ? err : new Error(errorMessage))
    }
  }, [sitekey, action, overrideDifficulty, onSuccess, onError, setupExpiryTimer])

  // Auto-start on mount if enabled
  useEffect(() => {
    if (autoStart && status === 'idle' && !token) {
      startVerification()
    }
  }, [autoStart, status, token, startVerification])

  // Retry handler
  const handleRetry = () => {
    if (workerRef.current) {
      workerRef.current.terminate()
      workerRef.current = null
    }
    startVerification()
  }

  // 点击开始验证
  const handleClick = () => {
    if (status === 'idle') {
      startVerification()
    } else if (status === 'error') {
      handleRetry()
    }
  }

  // Determine theme
  const resolvedTheme = theme === 'auto' 
    ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
    : theme
  const isDark = resolvedTheme === 'dark'

  // Size variants
  const sizeClasses = {
    compact: 'max-w-[260px]',
    normal: 'max-w-[300px]',
    flexible: 'w-full',
  }

  // 是否正在处理中
  const isProcessing = status === 'loading' || status === 'computing' || status === 'verifying'

  // 预设主题色映射
  const themeColorMap = {
    blue: '#3B82F6',
    violet: '#8B5CF6',
    emerald: '#10B981',
    rose: '#F43F5E',
    orange: '#F97316',
    cyan: '#06B6D4',
  }

  // 验证和处理自定义颜色
  const validatedCustomColor = customColor && isValidHexColor(customColor) 
    ? normalizeHexColor(customColor) 
    : null

  // 确定最终使用的颜色
  const finalColor = validatedCustomColor || themeColorMap[themeColor]

  // 生成颜色样式
  const colorStyles = {
    success: finalColor,
    loading: isDark ? `${finalColor}33` : `${finalColor}1A`, // 20% 和 10% 透明度
    loadingIcon: isDark ? adjustColorBrightness(finalColor, 40) : finalColor,
    progress: finalColor,
  }

  return (
    <div
      className={cn(
        'rounded-lg border transition-all',
        sizeClasses[size],
        isDark 
          ? 'bg-slate-900 border-slate-700 text-white' 
          : 'bg-white border-slate-200 text-slate-900 shadow-sm'
      )}
      data-theme={resolvedTheme}
      data-status={status}
    >
      {/* 主内容区 */}
      <div 
        className={cn(
          'flex items-center gap-3 p-3',
          (status === 'idle' || status === 'error') && 'cursor-pointer',
          status === 'idle' && (isDark ? 'hover:bg-slate-800/50' : 'hover:bg-slate-50')
        )}
        onClick={handleClick}
      >
        {/* Checkbox 样式的图标区域 */}
        <div className="relative flex-shrink-0">
          {status === 'idle' ? (
            // 未验证 - 空复选框
            <div className={cn(
              'w-6 h-6 rounded border-2 transition-colors',
              isDark 
                ? 'border-slate-600 hover:border-slate-500' 
                : 'border-slate-300 hover:border-slate-400'
            )} />
          ) : status === 'success' ? (
            // 验证成功 - 主题色勾选
            <div 
              className="w-6 h-6 rounded flex items-center justify-center"
              style={{ backgroundColor: colorStyles.success }}
            >
              <Check className="w-4 h-4 text-white" strokeWidth={3} />
            </div>
          ) : status === 'error' ? (
            // 验证失败 - 红色叉
            <div className="w-6 h-6 rounded bg-red-500 flex items-center justify-center">
              <X className="w-4 h-4 text-white" strokeWidth={3} />
            </div>
          ) : (
            // 验证中 - 旋转加载
            <div 
              className="w-6 h-6 rounded flex items-center justify-center"
              style={{ backgroundColor: colorStyles.loading }}
            >
              <Loader2 
                className="w-4 h-4 animate-spin"
                style={{ color: colorStyles.loadingIcon }}
              />
            </div>
          )}
        </div>

        {/* 文字区域 */}
        <div className="flex-1 min-w-0">
          <div className={cn(
            'text-sm transition-colors',
            isDark ? 'text-slate-200' : 'text-slate-700'
          )}>
            {status === 'idle' && '点击进行人机验证'}
            {status === 'loading' && '初始化...'}
            {status === 'computing' && (
              batchInfo && batchInfo.total > 1 
                ? `验证中 (${batchInfo.current}/${batchInfo.total})`
                : '验证中...'
            )}
            {status === 'verifying' && '提交验证...'}
            {status === 'success' && '验证成功'}
            {status === 'error' && (
              <span className="text-red-500" title={error || undefined}>
                {error || '验证失败'}，点击重试
              </span>
            )}
          </div>

          {/* 进度条 */}
          {isProcessing && (
            <div className={cn(
              'mt-1.5 h-1 rounded-full overflow-hidden',
              isDark ? 'bg-slate-800' : 'bg-slate-100'
            )}>
              <div
                className="h-full rounded-full transition-all duration-300 ease-out"
                style={{ backgroundColor: colorStyles.progress, width: `${progress}%` }}
              />
            </div>
          )}
        </div>

        {/* 重试按钮 */}
        {status === 'error' && (
          <button
            onClick={(e) => {
              e.stopPropagation()
              handleRetry()
            }}
            className={cn(
              'p-1.5 rounded transition-colors flex-shrink-0',
              isDark 
                ? 'hover:bg-slate-800 text-slate-400 hover:text-slate-300' 
                : 'hover:bg-slate-100 text-slate-400 hover:text-slate-600'
            )}
            aria-label="重试"
          >
            <RefreshCw className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* 品牌区 */}
      <div className={cn(
        'px-3 py-2 border-t text-xs flex items-center justify-between',
        isDark ? 'border-slate-800 text-slate-500' : 'border-slate-100 text-slate-400'
      )}>
        <span className="flex items-center gap-1.5">
          <svg className="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
          </svg>
          微验 WeAuth
        </span>
        <a
          href="https://www.weauth.cn/privacy"
          target="_blank"
          rel="noopener noreferrer"
          className={cn(
            'hover:underline transition-colors',
            isDark ? 'hover:text-slate-400' : 'hover:text-slate-500'
          )}
        >
          隐私
        </a>
      </div>

      {/* 隐藏的 token input，方便表单提交 */}
      {token && (
        <input type="hidden" name="weauth-token" value={token} />
      )}
    </div>
  )
})

// 导出供 script 方式使用的渲染函数
export function renderPowVerify(
  container: HTMLElement | string,
  options: Omit<PowVerifyProps, 'onSuccess' | 'onError'> & {
    callback?: (token: string) => void
    'error-callback'?: (error: Error) => void
  }
) {
  const element = typeof container === 'string' 
    ? document.querySelector(container) 
    : container
    
  if (!element) {
    console.error('WeAuth: Container not found')
    return
  }

  // 这里会在实际使用时通过 React 渲染
  console.log('WeAuth widget initialized', options)
}
