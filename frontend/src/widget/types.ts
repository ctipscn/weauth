// Widget 类型定义
// 从 pow-verify 提取的公共类型

/** PoW 验证状态 */
export type VerifyStatus = 'idle' | 'requesting' | 'solving' | 'verifying' | 'completed' | 'error'

/** Widget 验证配置（由 iframe URL 参数提供） */
export interface WidgetVerifyConfig {
  sitekey: string
  action?: string
  difficulty?: number
  autoStart?: boolean
  theme?: 'light' | 'dark' | 'auto'
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  customColor?: string
  size?: 'compact' | 'normal' | 'flexible'
}