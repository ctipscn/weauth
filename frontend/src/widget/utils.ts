// 工具函数：从 URL 解析参数、生成 widgetId、构建 iframe URL 等
import { WIDGET_BASE_URL } from './config'

/* ==================== URL 参数解析（iframe 内部使用） ==================== */

/** iframe URL 中携带的 widget 配置 */
export interface WidgetUrlParams {
  sitekey: string
  widgetId: string
  action?: string
  difficulty?: number
  theme?: 'light' | 'dark' | 'auto'
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  customColor?: string
  size?: 'compact' | 'normal' | 'flexible'
  autoStart?: boolean
}

/**
 * 从 iframe 当前 URL 的查询参数解析 widget 配置。
 * 在 widget.html 内部调用。
 */
export function parseWidgetParams(search: string = window.location.search): WidgetUrlParams {
  const params = new URLSearchParams(search)

  const sitekey = params.get('sitekey') || ''
  const widgetId = params.get('widgetId') || 'weauth-widget-0'

  const difficulty = params.get('difficulty')
  const theme = params.get('theme') as WidgetUrlParams['theme']
  const themeColor = params.get('themeColor') as WidgetUrlParams['themeColor']
  const size = params.get('size') as WidgetUrlParams['size']

  return {
    sitekey,
    widgetId,
    action: params.get('action') || undefined,
    difficulty: difficulty ? parseInt(difficulty, 10) : undefined,
    theme: theme || 'light',
    themeColor: themeColor || 'blue',
    customColor: params.get('customColor') || undefined,
    size: size || 'normal',
    autoStart: params.get('autoStart') === 'true',
  }
}

/* ==================== iframe 构建（宿主页面使用） ==================== */

/** 宿主页面传入的 iframe 构建参数（export.ts 中的 WeAuthOptions 子集） */
export interface IframeOptions {
  sitekey: string
  action?: string
  difficulty?: number
  theme?: string
  themeColor?: string
  customColor?: string
  size?: string
  autoStart?: boolean
}

/** 生成唯一 widgetId（模块级计数器保证递增） */
let widgetCounter = 0

export function generateWidgetId(): string {
  return `weauth-widget-${widgetCounter++}`
}

/** 构建 widget iframe 的 URL */
export function buildIframeUrl(options: IframeOptions, widgetId: string): string {
  const params = new URLSearchParams({
    sitekey: options.sitekey,
    widgetId,
  })

  if (options.action) params.set('action', options.action)
  if (options.difficulty !== undefined) params.set('difficulty', options.difficulty.toString())
  if (options.theme) params.set('theme', options.theme)
  if (options.themeColor) params.set('themeColor', options.themeColor)
  if (options.customColor) params.set('customColor', options.customColor)
  if (options.size) params.set('size', options.size)
  if (options.autoStart !== undefined) params.set('autoStart', options.autoStart.toString())

  return `${WIDGET_BASE_URL}/widget.html?${params.toString()}`
}

/* ==================== postMessage 消息辅助 ==================== */

/** iframe 与宿主之间约定的消息类型 */
export const MESSAGE_TYPES = {
  /** 宿主 → iframe：重置验证 */
  reset: 'weauth:reset',
  /** 宿主 → iframe：手动触发验证 */
  execute: 'weauth:execute',
  /** iframe → 宿主：验证成功 */
  success: 'weauth:success',
  /** iframe → 宿主：验证失败 */
  error: 'weauth:error',
  /** iframe → 宿主：token 过期 */
  expired: 'weauth:expired',
  /** iframe → 宿主：尺寸变化 */
  resize: 'weauth:resize',
} as const

/** 判断消息是否为 WeAuth 控制消息 */
export function isWeAuthMessage(data: unknown): data is Record<string, unknown> {
  const d = data as Record<string, unknown>
  return (
    typeof data === 'object' &&
    data !== null &&
    typeof d.type === 'string' &&
    d.type.startsWith('weauth:')
  )
}