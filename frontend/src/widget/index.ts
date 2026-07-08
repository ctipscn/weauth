import { WIDGET_BASE_URL } from './config'

// WeAuth 配置选项
export interface WeAuthOptions {
  sitekey: string
  action?: string
  difficulty?: number
  theme?: 'light' | 'dark' | 'auto'
  themeColor?: 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'
  customColor?: string
  size?: 'compact' | 'normal' | 'flexible'
  autoStart?: boolean
  callback?: (token: string) => void
  'error-callback'?: (error: Error) => void
  'expired-callback'?: () => void
}

// Widget 实例
interface WidgetInstance {
  iframe: HTMLIFrameElement
  container: HTMLElement
  options: WeAuthOptions
}

// 全局 API
export interface WeAuthAPI {
  render(container: string | HTMLElement, options: WeAuthOptions): string
  reset(widgetId: string): void
  execute(widgetId: string): void
  remove(widgetId: string): void
}

// Widget 实例存储
const widgetInstances = new Map<string, WidgetInstance>()
let widgetCounter = 0

// 生成唯一 widgetId
function generateWidgetId(): string {
  return `weauth-widget-${widgetCounter++}`
}

// 创建 iframe
function createIframe(
  container: HTMLElement,
  widgetId: string,
  options: WeAuthOptions
): HTMLIFrameElement {
  // 构建 iframe URL
  const params = new URLSearchParams({
    sitekey: options.sitekey,
    widgetId: widgetId,
  })

  if (options.action) params.set('action', options.action)
  if (options.difficulty) params.set('difficulty', options.difficulty.toString())
  if (options.theme) params.set('theme', options.theme)
  if (options.themeColor) params.set('themeColor', options.themeColor)
  if (options.customColor) params.set('customColor', options.customColor)
  if (options.size) params.set('size', options.size)
  if (options.autoStart !== undefined) params.set('autoStart', options.autoStart.toString())

  const iframeUrl = `${WIDGET_BASE_URL}/widget.html?${params.toString()}`

  // 创建 iframe
  const iframe = document.createElement('iframe')
  iframe.src = iframeUrl
  iframe.id = `weauth-iframe-${widgetId}`
  iframe.style.width = options.size === 'flexible' ? '100%' : '300px'
  iframe.style.height = '65px'
  iframe.style.border = 'none'
  iframe.style.overflow = 'hidden'
  iframe.style.display = 'block'
  iframe.setAttribute('scrolling', 'no')
  iframe.setAttribute('frameborder', '0')
  // 注意：sandbox 属性会限制 iframe 的能力，需要适当的权限
  iframe.setAttribute('sandbox', 'allow-same-origin allow-scripts allow-forms')

  container.appendChild(iframe)
  return iframe
}

// 渲染 Widget
function render(
  container: string | HTMLElement,
  options: WeAuthOptions
): string {
  // 验证必填参数
  if (!options.sitekey) {
    throw new Error('WeAuth: sitekey is required')
  }

  // 获取容器元素
  const containerElement = typeof container === 'string'
    ? document.querySelector<HTMLElement>(container)
    : container

  if (!containerElement) {
    throw new Error(`WeAuth: Container ${container} not found`)
  }

  // 生成 widgetId
  const widgetId = generateWidgetId()

  // 创建 iframe
  const iframe = createIframe(containerElement, widgetId, options)

  // 存储实例
  widgetInstances.set(widgetId, {
    iframe,
    container: containerElement,
    options,
  })

  return widgetId
}

// 重置 Widget
function reset(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget ${widgetId} not found`)
    return
  }

  // 检查 iframe 是否已加载
  if (!instance.iframe.contentWindow) {
    console.warn(`WeAuth: iframe contentWindow is not available for widget ${widgetId}`)
    return
  }

  // 发送重置消息到 iframe
  try {
    instance.iframe.contentWindow.postMessage(
      {
        type: 'weauth:reset',
        widgetId,
      },
      WIDGET_BASE_URL
    )
  } catch (error) {
    console.error(`WeAuth: Failed to send reset message to widget ${widgetId}:`, error)
  }
}

// 手动触发验证
function execute(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget ${widgetId} not found`)
    return
  }

  // 发送执行消息到 iframe
  instance.iframe.contentWindow?.postMessage(
    {
      type: 'weauth:execute',
      widgetId,
    },
    WIDGET_BASE_URL
  )
}

// 移除 Widget
function remove(widgetId: string): void {
  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    console.warn(`WeAuth: Widget ${widgetId} not found`)
    return
  }

  // 移除 iframe
  instance.iframe.remove()

  // 删除实例
  widgetInstances.delete(widgetId)
}

// 监听来自 iframe 的消息
window.addEventListener('message', (event) => {
  // 验证来源
  if (event.origin !== WIDGET_BASE_URL) {
    return
  }

  const { type, widgetId, token, error, height } = event.data

  const instance = widgetInstances.get(widgetId)
  if (!instance) {
    return
  }

  switch (type) {
    case 'weauth:success':
      instance.options.callback?.(token)
      break
    case 'weauth:error':
      instance.options['error-callback']?.(new Error(error))
      break
    case 'weauth:expired':
      instance.options['expired-callback']?.()
      break
    case 'weauth:resize':
      if (height) {
        instance.iframe.style.height = `${height}px`
      }
      break
  }
})

// 自动扫描和渲染
function autoRender() {
  const widgets = document.querySelectorAll('.weauth-widget')
  widgets.forEach((element) => {
    // 跳过已渲染的
    if (element.hasAttribute('data-widget-id')) {
      return
    }

    // 读取配置
    const options: WeAuthOptions = {
      sitekey: element.getAttribute('data-sitekey') || '',
      action: element.getAttribute('data-action') || undefined,
      difficulty: element.getAttribute('data-difficulty') 
        ? parseInt(element.getAttribute('data-difficulty')!) 
        : undefined,
      theme: (element.getAttribute('data-theme') as any) || undefined,
      themeColor: (element.getAttribute('data-theme-color') as any) || undefined,
      customColor: element.getAttribute('data-custom-color') || undefined,
      size: (element.getAttribute('data-size') as any) || undefined,
      autoStart: element.getAttribute('data-auto-start') === 'true',
    }

    // 从 data-callback 读取全局函数名
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

// WeAuth API
const weauth: WeAuthAPI = {
  render,
  reset,
  execute,
  remove,
}

// 挂载到 window
declare global {
  interface Window {
    weauth: WeAuthAPI
  }
}

if (typeof window !== 'undefined') {
  window.weauth = weauth
}

// DOMContentLoaded 或脚本加载后执行自动渲染
if (typeof document !== 'undefined') {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', autoRender)
  } else {
    // DOM 已加载完成，立即执行
    autoRender()
  }
}

export default weauth
