import { useState, useRef, useEffect } from 'react'
import { PowVerify } from '@/components/pow-verify'
import { Link } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { ArrowLeft, Check, Code, Copy, Clock, Zap, Timer, User, X, Info, Shield, RefreshCw } from 'lucide-react'
import Navbar from '@/components/Navbar'

// 预设难度测试
const difficultyPresets = [
  { difficulty: 3, label: '简单', desc: '~0.1s' },
  { difficulty: 4, label: '标准', desc: '~0.5s' },
  { difficulty: 5, label: '中等', desc: '~2s' },
  { difficulty: 6, label: '困难', desc: '~15s' },
]

export default function Demo() {
  const { token: authToken } = useAuth()
  const isLoggedIn = !!authToken

  const [sitekey, setSitekey] = useState('')
  const [action, setAction] = useState('demo')
  const [selectedTheme, setSelectedTheme] = useState<'light' | 'dark' | 'auto'>('light')
  const [selectedColor, setSelectedColor] = useState<'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'>('blue')
  const [customColorInput, setCustomColorInput] = useState('')
  const [autoStart, setAutoStart] = useState(false)
  const [token_, setToken_] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState<string | null>(null)
  const [solveTime, setSolveTime] = useState<number | null>(null)
  const [startTime, setStartTime] = useState<number | null>(null)

  // 难度测试状态
  const [testingDifficulty, setTestingDifficulty] = useState<number | null>(null)
  const [testResults, setTestResults] = useState<Record<number, number>>({})

  // 模拟表单提交状态
  const [formSubmitted, setFormSubmitted] = useState(false)
  const [loginToken, setLoginToken] = useState<string | null>(null)
  const [verifyResult, setVerifyResult] = useState<{ success: boolean; message: string } | null>(null)
  const [formSitekey, setFormSitekey] = useState('') // 用户输入的 sitekey
  const [secretKey, setSecretKey] = useState('') // 用户输入的 secret
  // const powVerifyRef = useRef<PowVerifyRef>(null) // PowVerify 组件的 ref（用于难度测试）
  const iframeRef = useRef<HTMLIFrameElement>(null) // iframe ref（用于模拟登录表单）
  const previewIframeRef = useRef<HTMLIFrameElement>(null) // iframe ref（用于预览面板）

  // 监听来自 iframe 的消息
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      // 开发环境放宽限制
      const allowedOrigins = [
        'http://localhost:5173',
        window.location.origin,
        'https://www.weauth.cn'
      ]
      
      if (!allowedOrigins.includes(event.origin)) {
        return
      }

      const { type, token, widgetId } = event.data

      // 模拟登录表单的 iframe
      if (widgetId === 'demo-login-widget') {
        if (type === 'weauth:success' && token) {
          setLoginToken(token)
        } else if (type === 'weauth:expired') {
          setLoginToken(null)
          setVerifyResult(null)
        }
      }
      
      // 预览面板的 iframe
      if (widgetId === 'preview-widget') {
        if (type === 'weauth:success' && token) {
          handleSuccess(token)
        } else if (type === 'weauth:error') {
          handleError(new Error(event.data.error || '验证失败'))
        }
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [])

  const handleSuccess = (newToken: string) => {
    const endTime = Date.now()
    if (startTime) {
      const time = endTime - startTime
      setSolveTime(time)
      if (testingDifficulty !== null) {
        setTestResults(prev => ({ ...prev, [testingDifficulty]: time }))
      }
    }
    setToken_(newToken)
    setError(null)
    setTestingDifficulty(null)
  }

  const handleError = (err: Error) => {
    setError(err.message)
    setToken_(null)
    setTestingDifficulty(null)
  }

  const handleCopy = async (text: string, field: string) => {
    await navigator.clipboard.writeText(text)
    setCopied(field)
    setTimeout(() => setCopied(null), 2000)
  }

  // 开始难度测试
  const startDifficultyTest = (difficulty: number) => {
    setTestingDifficulty(difficulty)
    setToken_(null)
    setError(null)
    setSolveTime(null)
    setStartTime(Date.now())
  }

  // 处理模拟登录表单提交
  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formSitekey) {
      setVerifyResult({ success: false, message: '请输入站点 Sitekey' })
      return
    }
    
    if (!loginToken) {
      setVerifyResult({ success: false, message: '请先完成人机验证' })
      return
    }
    
    if (!secretKey) {
      setVerifyResult({ success: false, message: '请输入站点 Secret 密钥' })
      return
    }
    
    setFormSubmitted(true)
    setVerifyResult(null)
    
    try {
      // 调用后端验证 API
      const response = await fetch('http://localhost:8080/api/pow/siteverify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          secret: secretKey,
          token: loginToken
        })
      })
      
      const result = await response.json()
      
      if (result.success) {
        setVerifyResult({ 
          success: true, 
          message: `✅ 验证通过！来源: ${result.hostname || 'N/A'}, 动作: ${result.action || 'N/A'}` 
        })
        // 验证成功后，Token 已被消耗，自动重置验证组件
        setTimeout(() => {
          setLoginToken(null)
          setVerifyResult(null)
          // 发送重置消息给 iframe
          if (iframeRef.current?.contentWindow) {
            iframeRef.current.contentWindow.postMessage(
              { type: 'weauth:reset', widgetId: 'demo-login-widget' },
              window.location.origin
            )
          }
        }, 2000) // 2秒后重置，让用户看到成功消息
      } else {
        const errorMsg = result['error-codes']?.[0] || 'unknown-error'
        const errorMessages: Record<string, string> = {
          'invalid-input-secret': 'Secret 密钥无效',
          'invalid-input-response': 'Token 无效或已过期',
          'timeout-or-duplicate': 'Token 已使用或已过期'
        }
        setVerifyResult({ 
          success: false, 
          message: '❌ ' + (errorMessages[errorMsg] || '验证失败')
        })
      }
    } catch (err) {
      setVerifyResult({ 
        success: false, 
        message: '❌ 网络错误，请检查后端服务是否运行' 
      })
    } finally {
      setFormSubmitted(false)
    }
  }

  // Token 过期时清除（已通过 iframe 消息处理，此处保留以备将来使用）
  // const handleLoginTokenExpire = () => {
  //   setLoginToken(null)
  //   setVerifyResult(null)
  // }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30">
      {/* Header */}
      <Navbar showMenu={true} />

      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6 pt-24">
        {/* Page Title */}
        <div className="mb-6">
          <Link
            to={isLoggedIn ? "/console/sites" : "/"}
            className="inline-flex items-center text-sm text-slate-500 hover:text-slate-700 transition-colors mb-3"
          >
            <ArrowLeft className="h-4 w-4 mr-1" />
            返回
          </Link>
          <h1 className="text-xl font-semibold text-slate-900">验证组件演示</h1>
          <p className="text-slate-500 text-sm mt-0.5">
            测试工作量证明验证组件
          </p>
        </div>

        {/* 模拟登录表单 - 类似 Turnstile 的使用场景 */}
        <div className="bg-white rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden mb-6">
          <div className="px-4 py-3 border-b border-slate-100">
            <h2 className="font-medium text-slate-900 flex items-center gap-2 text-sm">
              <User className="w-4 h-4 text-blue-500" />
              模拟登录表单
              <span className="text-xs text-slate-400 font-normal">（类似 Cloudflare Turnstile 用法）</span>
            </h2>
          </div>
          
          <form onSubmit={handleLoginSubmit} className="p-6">
            <div className="max-w-sm mx-auto space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">站点 Sitekey</label>
                <input
                  type="text"
                  value={formSitekey}
                  onChange={(e) => setFormSitekey(e.target.value)}
                  placeholder="输入站点 Sitekey（公开密钥）"
                  className="w-full h-10 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 text-sm font-mono focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                />
                <p className="text-xs text-slate-400 mt-1">
                  在「控制台 - 站点管理」复制 Sitekey
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">站点 Secret</label>
                <input
                  type="text"
                  value={secretKey}
                  onChange={(e) => setSecretKey(e.target.value)}
                  placeholder="输入站点 Secret 密钥（私密）"
                  className="w-full h-10 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 text-sm font-mono focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                />
                <p className="text-xs text-slate-400 mt-1">
                  在站点详情页点击「查看」获取 Secret
                </p>
              </div>

              <div className="p-3 rounded-lg bg-blue-50 border border-blue-200 text-blue-700 text-xs">
                <div className="flex items-start gap-2">
                  <Info className="w-4 h-4 mt-0.5 flex-shrink-0" />
                  <div>
                    <p className="font-medium mb-1">域名白名单提示</p>
                    <p>如果站点配置了域名白名单，请将 <code className="bg-blue-100 px-1 rounded">localhost:5173</code> 添加到白名单，或留空允许所有域名。</p>
                  </div>
                </div>
              </div>

              {/* WeAuth 验证组件 - iframe 方式 */}
              <div className="py-2">
                {formSitekey ? (
                  <iframe
                    ref={iframeRef}
                    key={formSitekey}
                    src={`/widget.html?sitekey=${encodeURIComponent(formSitekey)}&action=demo-login&theme=light&widgetId=demo-login-widget`}
                    width="300"
                    height="65"
                    style={{ border: 'none', display: 'block' }}
                    title="WeAuth Widget"
                  />
                ) : (
                  <div className="text-center py-8 text-slate-400 text-sm border-2 border-dashed border-slate-200 rounded-lg">
                    请先输入 Sitekey
                  </div>
                )}
              </div>

              <button
                type="submit"
                disabled={!loginToken || !formSitekey || !secretKey || formSubmitted}
                className="w-full h-10 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 disabled:bg-slate-100 disabled:border-slate-300 disabled:text-slate-400 disabled:cursor-not-allowed text-sm font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
              >
                {formSubmitted ? (
                  <>验证中...</>
                ) : (
                  <>
                    <Check className="w-4 h-4" />
                    验证 Token
                  </>
                )}
              </button>

              {/* 验证结果显示 */}
              {verifyResult && (
                <div className={`p-3 rounded-lg text-sm ${
                  verifyResult.success 
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' 
                      : 'bg-rose-50 text-rose-700 border border-rose-200'
                }`}>
                  <div className="flex items-center gap-2">
                    {verifyResult.success ? (
                      <Check className="w-4 h-4" />
                    ) : (
                      <X className="w-4 h-4" />
                    )}
                    <span className="font-medium">{verifyResult.message}</span>
                  </div>
                  {verifyResult.success && (
                    <p className="text-xs mt-1 text-emerald-600">
                      Token 已消耗，2 秒后自动重置
                    </p>
                  )}
                </div>
              )}

              {!loginToken && !verifyResult && (
                <p className="text-xs text-center text-slate-400">
                  填写 Sitekey 和 Secret，完成验证后点击「验证 Token」
                </p>
              )}
            </div>
          </form>
        </div>

        {/* Difficulty Tests */}
        <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden mb-6">
          <div className="px-4 py-3 border-b border-slate-100">
            <h2 className="font-medium text-slate-900 flex items-center gap-2 text-sm">
              <Zap className="w-4 h-4 text-amber-500" />
              难度测试
              <span className="text-xs text-slate-400 font-normal">（自动开始模式）</span>
            </h2>
          </div>
          
          <div className="p-4">
            <div className="flex flex-wrap gap-2">
              {difficultyPresets.map((preset) => (
                <button
                  key={preset.difficulty}
                  onClick={() => startDifficultyTest(preset.difficulty)}
                  disabled={testingDifficulty !== null}
                  className={`px-3 py-2 rounded-lg border text-sm transition-all ${
                    testingDifficulty === preset.difficulty
                      ? 'border-blue-500 bg-blue-50'
                      : testResults[preset.difficulty]
                      ? 'border-emerald-200 bg-emerald-50'
                      : 'border-slate-200 hover:border-slate-300 bg-white'
                  } disabled:opacity-50`}
                >
                  <span className="font-medium text-slate-900">{preset.label}</span>
                  <span className="text-slate-400 ml-1.5">{preset.desc}</span>
                  {testResults[preset.difficulty] && (
                    <span className="text-emerald-600 ml-2 font-medium">
                      {(testResults[preset.difficulty] / 1000).toFixed(2)}s
                    </span>
                  )}
                </button>
              ))}
            </div>

            {/* 难度测试组件 */}
            {testingDifficulty !== null && (
              <div className="mt-4 p-4 bg-slate-50 rounded-lg">
                <div className="flex items-center gap-2 mb-3">
                  <Timer className="w-4 h-4 text-blue-500 animate-pulse" />
                  <span className="text-sm text-slate-600">
                    测试难度 {testingDifficulty}...
                  </span>
                </div>
                  <PowVerify
                    key={`test-${testingDifficulty}-${customColorInput}`}
                    sitekey="test-site-key-12345678"
                    action={`difficulty-test-${testingDifficulty}`}
                    difficulty={testingDifficulty}
                    theme={selectedTheme}
                    themeColor={selectedColor}
                    customColor={customColorInput || undefined}
                    autoStart={true}
                    onSuccess={handleSuccess}
                    onError={handleError}
                  />
              </div>
            )}
          </div>
        </div>

        <div className="grid lg:grid-cols-2 gap-6">
          {/* Config Panel */}
          <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
          <div className="px-4 py-3 border-b border-slate-100">
            <h2 className="font-medium text-slate-900 flex items-center gap-2 text-sm">
              <Code className="w-4 h-4 text-indigo-500" />
              自定义配置
            </h2>
          </div>
            
            <div className="p-4 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">Sitekey</label>
                <input
                  type="text"
                  value={sitekey}
                  onChange={(e) => {
                    setSitekey(e.target.value)
                    setToken_(null)
                    setError(null)
                  }}
                  placeholder="输入您的 sitekey"
                  className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm font-mono focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">Action</label>
                <input
                  type="text"
                  value={action}
                  onChange={(e) => setAction(e.target.value)}
                  placeholder="login, register..."
                  className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">主题</label>
                <div className="flex gap-2">
                  {(['light', 'dark', 'auto'] as const).map((t) => (
                    <button
                      key={t}
                      onClick={() => setSelectedTheme(t)}
                      className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                        selectedTheme === t
                          ? 'bg-white border-2 border-indigo-600 text-indigo-600 hover:bg-indigo-50'
                          : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                      }`}
                    >
                      {t === 'light' ? '浅色' : t === 'dark' ? '深色' : '自动'}
                    </button>
                  ))}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">主题色（快捷选择）</label>
                <div className="grid grid-cols-3 gap-2">
                  {([
                    { value: 'blue', label: '蓝色', color: '#3B82F6' },
                    { value: 'violet', label: '紫色', color: '#8B5CF6' },
                    { value: 'emerald', label: '翠绿', color: '#10B981' },
                    { value: 'rose', label: '玫瑰', color: '#F43F5E' },
                    { value: 'orange', label: '橙色', color: '#F97316' },
                    { value: 'cyan', label: '青色', color: '#06B6D4' },
                  ] as const).map((c) => (
                    <button
                      key={c.value}
                      onClick={() => {
                        setCustomColorInput(c.color) // 填充到自定义颜色
                        setSelectedColor(c.value)
                      }}
                      className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors flex items-center gap-2 ${
                        customColorInput === c.color
                          ? 'bg-white border-2 border-slate-400 text-slate-900'
                          : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                      }`}
                      title={`点击填充 ${c.color}`}
                    >
                      <div className="w-3 h-3 rounded-full" style={{ backgroundColor: c.color }}></div>
                      {c.label}
                    </button>
                  ))}
                </div>
                <p className="text-xs text-slate-400 mt-1">
                  点击预设颜色会自动填充到下方输入框，可直接使用或微调
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">
                  自定义颜色
                </label>
                <div className="flex gap-2 items-center">
                  <input
                    type="text"
                    value={customColorInput}
                    onChange={(e) => setCustomColorInput(e.target.value)}
                    placeholder="#3B82F6 或 #38F"
                    className="flex-1 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm font-mono focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
                  />
                  {customColorInput && (
                    <div 
                      className="w-9 h-9 rounded-lg border-2 border-slate-300 flex-shrink-0"
                      style={{ backgroundColor: /^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$/.test(customColorInput) ? customColorInput : 'transparent' }}
                    />
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-1">
                  输入或修改 16 进制颜色（支持 #RRGGBB 或 #RGB 格式）
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">启动模式</label>
                <div className="flex gap-2">
                  <button
                    onClick={() => setAutoStart(false)}
                    className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                      !autoStart
                        ? 'bg-blue-600 text-white'
                        : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}
                  >
                    点击验证
                  </button>
                  <button
                    onClick={() => setAutoStart(true)}
                    className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                      autoStart
                        ? 'bg-blue-600 text-white'
                        : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}
                  >
                    自动验证
                  </button>
                </div>
                <p className="text-xs text-slate-400 mt-1">
                  {autoStart ? '组件加载后自动开始验证' : '用户点击后开始验证（推荐）'}
                </p>
              </div>
            </div>
          </div>

          {/* Preview Panel */}
          <div className="space-y-6">
            <div className={`rounded-2xl border overflow-hidden backdrop-blur-xl shadow-sm hover:shadow-md transition-shadow ${
              selectedTheme === 'dark' ? 'bg-slate-900/80 border-slate-700/60' : 'bg-white/80 border-slate-200/60'
            }`}>
              <div className={`px-4 py-3 border-b flex items-center justify-between ${
                selectedTheme === 'dark' ? 'border-slate-700' : 'border-slate-100'
              }`}>
                <h2 className={`font-medium text-sm ${
                  selectedTheme === 'dark' ? 'text-white' : 'text-slate-900'
                }`}>预览</h2>
                {sitekey && (
                  <button
                    onClick={() => {
                      // 重置预览 iframe
                      if (previewIframeRef.current?.contentWindow) {
                        previewIframeRef.current.contentWindow.postMessage(
                          { type: 'weauth:reset', widgetId: 'preview-widget' },
                          window.location.origin
                        )
                      }
                      // 重置本地状态
                      setToken_(null)
                      setError(null)
                      setSolveTime(null)
                    }}
                    className={`p-1.5 rounded transition-colors ${
                      selectedTheme === 'dark'
                        ? 'hover:bg-slate-800 text-slate-400 hover:text-slate-300'
                        : 'hover:bg-slate-100 text-slate-400 hover:text-slate-600'
                    }`}
                    title="重置验证组件"
                    aria-label="重置"
                  >
                    <RefreshCw className="w-4 h-4" />
                  </button>
                )}
              </div>
              <div className="p-8 flex justify-center">
                {sitekey ? (
                  <iframe
                    ref={previewIframeRef}
                    key={`preview-${sitekey}-${action}-${autoStart}-${customColorInput}-${selectedTheme}-${selectedColor}`}
                    src={`/widget.html?sitekey=${encodeURIComponent(sitekey)}&action=${encodeURIComponent(action)}&theme=${selectedTheme}&themeColor=${selectedColor}${customColorInput ? `&customColor=${encodeURIComponent(customColorInput)}` : ''}&size=normal&autoStart=${autoStart}&widgetId=preview-widget`}
                    width="600"
                    height="100"
                    style={{ border: 'none', display: 'block', maxWidth: '100%' }}
                    title="WeAuth Widget Preview"
                  />
                ) : (
                  <div className={`text-center py-8 ${
                    selectedTheme === 'dark' ? 'text-gray-500' : 'text-slate-400'
                  }`}>
                    <div className="w-10 h-10 mx-auto mb-2 opacity-30 rounded-lg bg-indigo-100 flex items-center justify-center">
                      <Shield className="w-6 h-6 text-indigo-400" />
                    </div>
                    <p className="text-sm">输入 sitekey 预览组件</p>
                  </div>
                )}
              </div>
            </div>

            {/* Result */}
            {(token_ || error) && !testingDifficulty && (
              <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
                <div className="px-4 py-3 border-b border-slate-100">
                  <h2 className="font-medium text-slate-900 text-sm">结果</h2>
                </div>
                <div className="p-4">
                  {token_ && (
                    <div className="space-y-3">
                      <div className="flex items-center gap-2 text-emerald-600">
                        <Check className="w-4 h-4" />
                        <span className="text-sm font-medium">验证成功</span>
                        {solveTime && (
                          <span className="text-slate-500 text-sm flex items-center gap-1">
                            <Clock className="w-3 h-3" />
                            {(solveTime / 1000).toFixed(2)}s
                          </span>
                        )}
                      </div>
                      <div>
                        <div className="flex items-center justify-between mb-1">
                          <label className="text-xs font-medium text-slate-500">Token</label>
                          <button
                            onClick={() => handleCopy(token_, 'token')}
                            className="text-xs text-slate-400 hover:text-slate-600 transition-colors flex items-center gap-1"
                          >
                            {copied === 'token' ? <Check className="w-3 h-3" /> : <Copy className="w-3 h-3" />}
                            {copied === 'token' ? '已复制' : '复制'}
                          </button>
                        </div>
                        <code className="block bg-slate-50 border border-slate-200 p-2 rounded text-xs text-slate-600 break-all font-mono">
                          {token_}
                        </code>
                      </div>
                    </div>
                  )}

                  {error && (
                    <div className="text-rose-600 text-sm">
                      <span className="font-medium">错误: </span>{error}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Integration Guide */}
        <div className="mt-6 bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
          <div className="px-4 py-3 border-b border-slate-100">
            <h2 className="font-medium text-slate-900 text-sm">集成方式</h2>
          </div>
          <div className="p-4 space-y-6">
            {/* Method 1: Script 标签自动渲染 */}
            <div>
              <h3 className="text-sm font-medium text-slate-700 mb-2">方式一：Script 标签自动渲染（推荐）</h3>
              <pre className="bg-slate-900 text-slate-300 p-3 rounded-lg text-xs overflow-x-auto">
                <code>{`<script src="https://www.weauth.cn/weauth.js" async defer></script>

<div class="weauth-widget"
     data-sitekey="${sitekey || 'your-sitekey'}"
     data-action="${action}"
     data-theme="${selectedTheme}"
     data-theme-color="${selectedColor}"
     ${customColorInput ? `data-custom-color="${customColorInput}"` : ''}
     data-auto-start="${autoStart}"
     data-callback="onVerify"></div>

<script>
  function onVerify(token) {
    console.log('验证成功:', token)
    // 提交表单或发送到后端
  }
</script>`}</code>
              </pre>
            </div>

            {/* Method 2: JavaScript API 手动渲染 */}
            <div>
              <h3 className="text-sm font-medium text-slate-700 mb-2">方式二：JavaScript API 手动渲染</h3>
              <pre className="bg-slate-900 text-slate-300 p-3 rounded-lg text-xs overflow-x-auto">
                <code>{`<div id="my-widget"></div>

<script src="https://www.weauth.cn/weauth.js" defer></script>
<script>
  window.onload = function() {
    const widgetId = weauth.render('#my-widget', {
      sitekey: '${sitekey || 'your-sitekey'}',
      action: '${action}',
      theme: '${selectedTheme}',
      themeColor: '${selectedColor}',
      ${customColorInput ? `customColor: '${customColorInput}',` : ''}
      autoStart: ${autoStart},
      callback: function(token) {
        console.log('Token:', token)
      }
    })
  }
</script>`}</code>
              </pre>
            </div>

            {/* Method 3: 纯 iframe */}
            <div>
              <h3 className="text-sm font-medium text-slate-700 mb-2">方式三：纯 iframe（高级）</h3>
              <pre className="bg-slate-900 text-slate-300 p-3 rounded-lg text-xs overflow-x-auto">
                <code>{`<iframe 
  src="https://www.weauth.cn/widget.html?sitekey=${sitekey || 'your-sitekey'}&action=${action}&theme=${selectedTheme}&themeColor=${selectedColor}${customColorInput ? `&customColor=${encodeURIComponent(customColorInput)}` : ''}&autoStart=${autoStart}"
  width="300" 
  height="65" 
  frameborder="0">
</iframe>

<script>
  window.addEventListener('message', function(event) {
    if (event.origin !== 'https://www.weauth.cn') return
    
    if (event.data.type === 'weauth:success') {
      console.log('Token:', event.data.token)
    }
  })
</script>`}</code>
              </pre>
            </div>

            {/* Method 4: 表单集成示例 */}
            <div>
              <h3 className="text-sm font-medium text-slate-700 mb-2">表单集成示例</h3>
              <pre className="bg-slate-900 text-slate-300 p-3 rounded-lg text-xs overflow-x-auto">
                <code>{`<form id="login-form" onsubmit="return handleSubmit(event)">
  <input type="text" name="username" required>
  <input type="password" name="password" required>
  
  <div class="weauth-widget"
       data-sitekey="${sitekey || 'your-sitekey'}"
       data-action="login"
       data-callback="onLoginVerify"></div>
  
  <button type="submit">登录</button>
</form>

<script>
  let verifyToken = null
  
  function onLoginVerify(token) {
    verifyToken = token
  }
  
  async function handleSubmit(event) {
    event.preventDefault()
    
    if (!verifyToken) {
      alert('请先完成人机验证')
      return false
    }
    
    const formData = new FormData(event.target)
    formData.append('weauth-token', verifyToken)
    
    // 提交到后端...
  }
</script>`}</code>
              </pre>
            </div>

            {/* Backend Verification */}
            <div>
              <h3 className="text-sm font-medium text-slate-700 mb-2">后端验证</h3>
              <pre className="bg-slate-900 text-slate-300 p-3 rounded-lg text-xs overflow-x-auto">
                <code>{`const response = await fetch('https://api.weauth.cn/api/pow/siteverify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    secret: 'your-secret-key',
    token: request.body['weauth-token']
  })
});

const { success, hostname, action } = await response.json();
if (success) {
  // 验证通过，继续处理登录/注册逻辑
}`}</code>
              </pre>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
