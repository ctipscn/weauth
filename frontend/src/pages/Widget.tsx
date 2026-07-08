import { useState, useEffect, useRef } from 'react'
import { PowVerify, type PowVerifyRef } from '@/components/pow-verify'
import type { PowVerifyProps } from '@/components/pow-verify/types'

export default function Widget() {
  const [config, setConfig] = useState<Omit<PowVerifyProps, 'onSuccess' | 'onError' | 'onExpire'> | null>(null)
  const powVerifyRef = useRef<PowVerifyRef>(null)
  const widgetIdRef = useRef<string>('')

  useEffect(() => {
    // 从 URL 参数读取配置
    const params = new URLSearchParams(window.location.search)
    
    const sitekey = params.get('sitekey')
    if (!sitekey) {
      console.error('WeAuth: sitekey is required')
      return
    }

    widgetIdRef.current = params.get('widgetId') || 'weauth-widget-0'

    const widgetConfig: Omit<PowVerifyProps, 'onSuccess' | 'onError' | 'onExpire'> = {
      sitekey,
      action: params.get('action') || undefined,
      difficulty: params.get('difficulty') ? parseInt(params.get('difficulty')!) : undefined,
      theme: (params.get('theme') as 'light' | 'dark' | 'auto') || 'light',
      themeColor: (params.get('themeColor') as any) || 'blue',
      customColor: params.get('customColor') || undefined,
      size: (params.get('size') as 'compact' | 'normal' | 'flexible') || 'normal',
      autoStart: params.get('autoStart') === 'true',
    }

    setConfig(widgetConfig)
  }, [])

  // 监听来自父页面的控制消息
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      // 验证消息格式
      if (!event.data || typeof event.data !== 'object') {
        return
      }

      const { type, widgetId } = event.data

      // 确保消息是发给当前 widget 的
      if (widgetId && widgetId !== widgetIdRef.current) {
        return
      }

      // 只处理 weauth 相关的消息类型
      if (!type || !type.startsWith('weauth:')) {
        return
      }

      // 验证消息来源（允许所有来源，因为 iframe 可能被嵌入到任何网站）
      // 但可以通过检查消息格式来确保是合法的 weauth 消息
      // 注意：在生产环境中，如果需要更严格的安全控制，可以在这里添加 origin 检查

      switch (type) {
        case 'weauth:reset':
          if (powVerifyRef.current) {
            powVerifyRef.current.reset()
          } else {
            console.warn('WeAuth: powVerifyRef is not ready, reset will be ignored')
          }
          break
        case 'weauth:execute':
          // 预留：手动触发验证（如果组件支持的话）
          break
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [])

  // 发送消息给父页面
  const sendToParent = (message: any) => {
    window.parent.postMessage({
      ...message,
      widgetId: widgetIdRef.current,
    }, '*') // 生产环境应指定具体的 targetOrigin
  }

  const handleSuccess = (token: string) => {
    sendToParent({
      type: 'weauth:success',
      token,
    })
  }

  const handleError = (error: Error) => {
    sendToParent({
      type: 'weauth:error',
      error: error.message,
    })
  }

  const handleExpire = () => {
    sendToParent({
      type: 'weauth:expired',
    })
  }

  // 通知父页面高度变化（可选功能）
  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const height = entry.contentRect.height
        sendToParent({
          type: 'weauth:resize',
          height: Math.ceil(height),
        })
      }
    })

    const widgetRoot = document.getElementById('widget-root')
    if (widgetRoot) {
      resizeObserver.observe(widgetRoot)
    }

    return () => resizeObserver.disconnect()
  }, [config])

  if (!config) {
    return (
      <div style={{ 
        padding: '20px', 
        textAlign: 'center', 
        color: '#666',
        fontSize: '14px' 
      }}>
        加载中...
      </div>
    )
  }

  return (
    <div id="widget-root" style={{ margin: 0, padding: 0 }}>
      <PowVerify
        ref={powVerifyRef}
        {...config}
        onSuccess={handleSuccess}
        onError={handleError}
        onExpire={handleExpire}
      />
    </div>
  )
}
