export interface PowVerifyProps {
  /** 站点密钥 */
  sitekey: string
  /** 操作类型（如 login, register） */
  action?: string
  /** 覆盖难度（仅测试站点有效） */
  difficulty?: number
  /** 验证成功回调 */
  onSuccess: (token: string) => void
  /** 验证失败回调 */
  onError?: (error: Error) => void
  /** Token 过期回调 */
  onExpire?: () => void
  /** 
   * 主题设置
   * - 'light': 浅色主题
   * - 'dark': 深色主题  
   * - 'auto': 自动检测系统主题
   */
  theme?: 'light' | 'dark' | 'auto'
  /**
   * 主题色
   * - 'blue': 蓝色（默认）
   * - 'violet': 紫色
   * - 'emerald': 翠绿色
   * - 'rose': 玫瑰红
   * - 'orange': 橙色
   * - 'cyan': 青色
   */
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  /**
   * 自定义主题色（16进制颜色值）
   * - 例如: '#3B82F6', '#8B5CF6', '#10B981'
   * - 优先级高于 themeColor
   * - 仅支持 #RRGGBB 或 #RGB 格式
   */
  customColor?: string
  /**
   * 显示模式
   * - 'always': 始终显示
   * - 'execute': 仅在执行时显示
   * - 'interaction-only': 仅交互时显示
   */
  appearance?: 'always' | 'execute' | 'interaction-only'
  /**
   * 尺寸
   * - 'compact': 紧凑模式
   * - 'normal': 标准模式
   * - 'flexible': 自适应宽度
   */
  size?: 'compact' | 'normal' | 'flexible'
  /** 自定义类名 */
  className?: string
  /**
   * 是否自动开始验证
   * - true: 组件加载后自动开始验证
   * - false: 需要用户点击才开始验证（默认，类似 Turnstile）
   */
  autoStart?: boolean
}

export type VerifyStatus = 'idle' | 'loading' | 'computing' | 'verifying' | 'success' | 'error'

export interface ChallengeData {
  challenge_id: string
  challenge?: string       // 单个模式
  challenges?: string[]    // 批量模式
  difficulty: number
  batch_mode: boolean
  batch_total: number
}

export interface WorkerMessage {
  type: 'start' | 'start_batch' | 'progress' | 'result' | 'batch_progress' | 'batch_result'
  challenge?: string
  challenges?: string[]
  difficulty?: number
  nonce?: string
  nonces?: string[]
  hash?: string
  hashes?: string[]
  iterations?: number
  progress?: number
  currentChallenge?: number
  totalChallenges?: number
  completedCount?: number
  batchCount?: number
  solveTimeMs?: number
}

export interface WorkerResult {
  nonce: string
  hash: string
  iterations: number
  solveTimeMs: number
}

export interface BatchWorkerResult {
  nonces: string[]
  hashes: string[]
  iterations: number
  solveTimeMs: number
  batchCount: number
}

// Script 方式的配置选项
export interface WeAuthWidgetOptions {
  sitekey: string
  action?: string
  theme?: 'light' | 'dark' | 'auto'
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  customColor?: string
  appearance?: 'always' | 'execute' | 'interaction-only'
  size?: 'compact' | 'normal' | 'flexible'
  callback?: (token: string) => void
  'error-callback'?: (error: Error) => void
  'expired-callback'?: () => void
}
