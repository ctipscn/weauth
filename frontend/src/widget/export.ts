// 导出暴露给外部 HTML 的方法
// 宿主页面通过 window.weauth 使用这些 API
import { WIDGET_BASE_URL } from './config'
import { generateWidgetId, buildIframeUrl, isWeAuthMessage, MESSAGE_TYPES } from './utils'

/* ==================== 类型定义 ==================== */

/** WeAuth widget 选项 */
export interface WeAuthOptions {
  sitekey: string
  action?: string
  difficulty?: number
  theme?: 'light' | 'dark' | 'auto'
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  customColor?: string
  size?: 'compact' | 'normal' | 'flexible'
  autoStart?: boolean
  /** 验证成功回调 */
  callback?: (token: string) => void
  /** 验证错误回调 */
  'error-callback'?: (error: Error) => void
  /** token 过期回调 */
  'expired-callback'?: () => void
}

/** WeAuth 全局 API 接口 */
export interface WeAuthAPI {
  /**
   * 在指定容器中渲染并返回 widgetId
   * @param container CSS 选择器或 HTMLElement
   * @param options  配置选项
   */
  render(container: string | HTMLElement, options: WeAuthOptions): string
  /** 重置 widget（重新开始验证流程） */
  reset(widgetId: string): void
  /** 手动触发验证 */
  execute(widgetId: string): void
  /** 移除 widget（从 DOM 中卸载） */
  remove(widgetId: string): void
}

/* ==================== 内部状态 ==================== */

interface WidgetInstance {
  iframe: HTMLIFrameElement
  container: HTMLElement
  options: WeAuthOptions
}

const widgetInstances = new Map<string, WidgetInstance>()

/* ==================== 核心方法 ==================== */

/**
 * 在主页面中嵌入 WeAuth widget iframe。
 * 返回唯一的 widgetId。
 */
function render(
  container: string | HTMLElement,
  options: WeAuthOptions,
): string {
  if (!options.sitekey) {
    throw new Error('WeAuth: sitekey is required')
  }

  const containerElement =
    typeof container === 'string'
      ? document.querySelector<HTMLElement>(container)
      : container

  if (!containerElement) {
    throw new Error(`WeAuth: Container "${container}" not found`)
  }

  const widgetId = generateWidgetId()
  const iframeUrl = buildIframeUrl(options, widgetId)

  const iframe = document.createElement('iframe')
  iframe.src = iframeUrl
  iframe.id = `weauth-iframe-${widgetId}`
  iframe.style.width = options.size === 'flexible' ? '100%' : options.size === 'compact' ? '200px' : '300px'
  iframe.style.height = '100px'  // 初始高度，组件会通过 postMessage 动态调整
  iframe.style.border = 'none'
  iframe.style.overflow = 'hidden'
  iframe.style.display = 'block'
  iframe.setAttribute('scrolling', 'no')
  iframe.setAttribute('frameborder', '0')
  iframe.setAttribute('sandbox', 'allow-same-origin allow-scripts allow-forms')

  containerElement.appendChild(iframe)

  widgetInstances.set(widgetId, { iframe, container: containerElement, options })

  return widgetId
}

/**
 * 向 iframe 发送重置消息。
 */
function reset(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget "${widgetId}" not found`)
    return
  }

  if (!instance.iframe.contentWindow) {
    console.warn(`WeAuth: iframe contentWindow unavailable for widget "${widgetId}"`)
    return
  }

  try {
    instance.iframe.contentWindow.postMessage(
      { type: MESSAGE_TYPES.reset, widgetId },
      WIDGET_BASE_URL,
    )
  } catch (error) {
    console.error(`WeAuth: Failed to send reset to widget "${widgetId}":`, error)
  }
}

/**
 * 向 iframe 发送手动触发验证消息。
 */
function execute(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget "${widgetId}" not found`)
    return
  }

  instance.iframe.contentWindow?.postMessage(
    { type: MESSAGE_TYPES.execute, widgetId },
    WIDGET_BASE_URL,
  )
}

/**
 * 从 DOM 移除 widget iframe 并清理状态。
 */
function remove(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget "${widgetId}" not found`)
    return
  }

  instance.iframe.remove()
  widgetInstances.delete(widgetId)
}

/* ==================== 自动渲染 ==================== */

/**
 * 扫描页面上所有带有 `.weauth-widget` class 的元素，
 * 根据其 data-* 属性自动渲染 widget。
 */
function autoRender(): void {
  const widgets = document.querySelectorAll('.weauth-widget')
  widgets.forEach((element) => {
    if (element.hasAttribute('data-widget-id')) {
      return // 已渲染，跳过
    }

    const options: WeAuthOptions = {
      sitekey: element.getAttribute('data-sitekey') || '',
      action: element.getAttribute('data-action') || undefined,
      difficulty: element.getAttribute('data-difficulty')
        ? parseInt(element.getAttribute('data-difficulty')!, 10)
        : undefined,
      theme: (element.getAttribute('data-theme') as WeAuthOptions['theme']) || undefined,
      themeColor: (element.getAttribute('data-theme-color') as WeAuthOptions['themeColor']) || undefined,
      customColor: element.getAttribute('data-custom-color') || undefined,
      size: (element.getAttribute('data-size') as WeAuthOptions['size']) || undefined,
      autoStart: element.getAttribute('data-auto-start') === 'true',
    }

    const callbackName = element.getAttribute('data-callback')
    if (callbackName && typeof (window as any)[callbackName] === 'function') {
      options.callback = (window as any)[callbackName]
    }

    const errorCallbackName = element.getAttribute('data-error-callback')
    if (errorCallbackName && typeof (window as any)[errorCallbackName] === 'function') {
      options['error-callback'] = (window as any)[errorCallbackName]
    }

    const expiredCallbackName = element.getAttribute('data-expired-callback')
    if (expiredCallbackName && typeof (window as any)[expiredCallbackName] === 'function') {
      options['expired-callback'] = (window as any)[expiredCallbackName]
    }

    try {
      const widgetId = render(element as HTMLElement, options)
      element.setAttribute('data-widget-id', widgetId)
    } catch (err) {
      console.error('WeAuth: Failed to render widget', err)
    }
  })
}

/* ==================== postMessage 监听 ==================== */

/** 监听来自 iframe 的 postMessage 并分发给对应回调 */
window.addEventListener('message', (event: MessageEvent) => {
  const data = event.data
  if (!isWeAuthMessage(data)) {
    return
  }

  const { type, widgetId, token, error, height } = data as Record<string, unknown>
  const instance = widgetInstances.get(widgetId as string)
  if (!instance) {
    return
  }

  switch (type) {
    case MESSAGE_TYPES.success:
      instance.options.callback?.(token as string)
      break

    case MESSAGE_TYPES.error:
      instance.options['error-callback']?.(new Error(error as string))
      break

    case MESSAGE_TYPES.expired:
      instance.options['expired-callback']?.()
      break

    case MESSAGE_TYPES.resize:
      if (typeof height === 'number') {
        instance.iframe.style.height = `${height}px`
      }
      break
  }
})

/* ==================== 全局挂载 & 自动渲染 ==================== */

const weauth: WeAuthAPI = { render, reset, execute, remove }

// 暴露到全局
if (typeof window !== 'undefined') {
  ;(window as any).weauth = weauth
}

// 如果页面已加载完成则立即扫描，否则等 DOMContentLoaded
if (typeof document !== 'undefined') {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', autoRender)
  } else {
    autoRender()
  }
}

export { weauth as default, autoRender }
export type { WidgetInstance }