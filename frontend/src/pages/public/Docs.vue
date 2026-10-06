<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/lib/auth'
import Navbar from '@/components/Navbar.vue'

import { ArrowLeft, Code, Book, Zap, Shield, Copy, Check, CheckCircle2 } from 'lucide-vue-next'

const auth = useAuthStore()
const isLoggedIn = !!auth.token

const copiedId = ref<string | null>(null)

const handleCopy = async (text: string, id: string) => {
  await navigator.clipboard.writeText(text)
  copiedId.value = id
  setTimeout(() => { copiedId.value = null }, 2000)
}

const copyCode = (id: string) => {
  const code = snippets[id]
  if (code) handleCopy(code, id)
}

const snippets: Record<string, string> = {
  'quick-1': '<script src="https://www.weauth.cn/weauth.js" async defer><\/script>',
  'quick-2': `<div class="weauth-widget"
     data-sitekey="your-sitekey-here"
     data-callback="onSuccess"></div>`,
  'quick-3': `function onSuccess(token) {
  console.log('验证成功:', token)
  // 提交表单或发送到后端
}`,
  'quick-4': `const response = await fetch('https://api.weauth.cn/api/pow/siteverify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    secret: 'your-secret-key',
    token: token
  })
})

const { success } = await response.json()`,

  'install-1': `<script src="https://www.weauth.cn/weauth.js" async defer><\/script>

<div class="weauth-widget"
     data-sitekey="your-sitekey"
     data-action="login"
     data-theme="light"
     data-callback="onVerify"></div>

<script>
  function onVerify(token) {
    document.getElementById('token-field').value = token
  }
<\/script>`,
  'install-2': `<div id="my-widget"></div>

<script src="https://www.weauth.cn/weauth.js" defer><\/script>
<script>
  window.onload = function() {
    const widgetId = weauth.render('#my-widget', {
      sitekey: 'your-sitekey',
      action: 'login',
      theme: 'light',
      callback: function(token) {
        console.log('Token:', token)
      }
    })
  }
<\/script>`,
  'install-3': `<iframe id="weauth-iframe"
        src="https://www.weauth.cn/widget.html?sitekey=your-key&widgetId=my-widget&theme=light"
        width="300" height="65" frameborder="0"></iframe>

<script>
  const iframe = document.getElementById('weauth-iframe')

  // 监听来自 iframe 的消息
  window.addEventListener('message', function(event) {
    if (event.origin !== 'https://www.weauth.cn') return

    if (event.data.type === 'weauth:success') {
      console.log('Token:', event.data.token)
    } else if (event.data.type === 'weauth:error') {
      console.error('验证失败:', event.data.error)
    } else if (event.data.type === 'weauth:expired') {
      console.log('Token 已过期')
    }
  })

  // 重置验证组件
  function resetWidget() {
    if (iframe.contentWindow) {
      iframe.contentWindow.postMessage({
        type: 'weauth:reset',
        widgetId: 'my-widget'
      }, 'https://www.weauth.cn')
    }
  }
<\/script>`,

  'api-js': `// 手动渲染
const widgetId = weauth.render(container, options)

// 重置组件
weauth.reset(widgetId)

// 手动触发验证（autoStart=false 时）
weauth.execute(widgetId)

// 移除组件
weauth.remove(widgetId)`,
  'api-iframe-reset': `// 获取 iframe 元素
const iframe = document.getElementById('weauth-iframe')

// 重置 iframe 中的验证组件
iframe.contentWindow.postMessage({
  type: 'weauth:reset',
  widgetId: 'your-widget-id'  // 必须与 URL 参数中的 widgetId 一致
}, 'https://www.weauth.cn')`,

  'example-login': `<form id="login-form" onsubmit="return handleSubmit(event)">
  <input type="text" name="username" required>
  <input type="password" name="password" required>

  <div class="weauth-widget"
       data-sitekey="your-sitekey"
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
<\/script>`,
  'example-react': `import { useEffect, useRef } from 'react'

function LoginForm() {
  const widgetRef = useRef<string>()

  useEffect(() => {
    widgetRef.current = window.weauth.render('#weauth-container', {
      sitekey: 'your-sitekey',
      callback: (token) => {
        console.log('Token:', token)
      }
    })

    return () => {
      if (widgetRef.current) {
        window.weauth.remove(widgetRef.current)
      }
    }
  }, [])

  return <div id="weauth-container"></div>
}`,

  'backend-node': `app.post('/api/login', async (req, res) => {
  const { 'weauth-token': token } = req.body

  const response = await fetch('https://api.weauth.cn/api/pow/siteverify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      secret: process.env.WEAUTH_SECRET,
      token: token
    })
  })

  const result = await response.json()

  if (!result.success) {
    return res.status(400).json({ error: '人机验证失败' })
  }

  // 继续登录逻辑...
})`,
  'backend-python': `import requests

@app.route('/api/login', methods=['POST'])
def login():
    token = request.json.get('weauth-token')

    response = requests.post('https://api.weauth.cn/api/pow/siteverify', json={
        'secret': os.environ['WEAUTH_SECRET'],
        'token': token
    })

    result = response.json()

    if not result['success']:
        return {'error': '人机验证失败'}, 400

    # 继续登录逻辑...`,
  'backend-go': `type VerifyRequest struct {
    Secret string \`json:"secret"\`
    Token  string \`json:"token"\`
}

type VerifyResponse struct {
    Success  bool   \`json:"success"\`
    Hostname string \`json:"hostname"\`
    Action   string \`json:"action"\`
}

func loginHandler(c *gin.Context) {
    var req LoginRequest
    c.BindJSON(&req)

    verifyReq := VerifyRequest{
        Secret: os.Getenv("WEAUTH_SECRET"),
        Token:  req.WeAuthToken,
    }

    resp, _ := http.Post(
        "https://api.weauth.cn/api/pow/siteverify",
        "application/json",
        bytes.NewBuffer(marshal(verifyReq)),
    )

    var result VerifyResponse
    json.NewDecoder(resp.Body).Decode(&result)

    if !result.Success {
        c.JSON(400, gin.H{"error": "人机验证失败"})
        return
    }

    // 继续登录逻辑...
}`,
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30">
    <Navbar />

    <main class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6 pt-24">
      <!-- Header -->
      <div class="mb-8">
        <RouterLink
          :to="isLoggedIn ? '/console/sites' : '/'"
          class="inline-flex items-center text-sm text-slate-500 hover:text-slate-700 transition-colors mb-4"
        >
          <ArrowLeft class="h-4 w-4 mr-1" />
          返回
        </RouterLink>
        <h1 class="text-3xl font-bold text-slate-900 mb-2">接入文档</h1>
        <p class="text-slate-600">
          WeAuth 是一个基于工作量证明（PoW）的人机验证服务，提供简单易用的 API 和组件
        </p>
      </div>

      <!-- 快速导航 -->
      <div class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Book class="w-5 h-5 text-blue-500" />
          快速导航
        </h2>
        <div class="grid sm:grid-cols-2 gap-3">
          <a href="#quick-start" class="text-blue-600 hover:text-blue-700 hover:underline">快速开始（3分钟集成）</a>
          <a href="#installation" class="text-blue-600 hover:text-blue-700 hover:underline">安装指南</a>
          <a href="#api-reference" class="text-blue-600 hover:text-blue-700 hover:underline">API 参考</a>
          <a href="#examples" class="text-blue-600 hover:text-blue-700 hover:underline">示例代码</a>
          <a href="#backend-verify" class="text-blue-600 hover:text-blue-700 hover:underline">后端验证</a>
          <a href="#faq" class="text-blue-600 hover:text-blue-700 hover:underline">常见问题</a>
        </div>
      </div>

      <!-- 快速开始 -->
      <section id="quick-start" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4 flex items-center gap-2">
          <Zap class="w-6 h-6 text-amber-500" />
          快速开始
        </h2>
        <p class="text-slate-600 mb-6">
          只需 3 分钟，即可在您的网站中集成 WeAuth 人机验证
        </p>

        <div class="space-y-6">
          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-3 flex items-center gap-2">
              <span class="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-bold">1</span>
              引入 Script
            </h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['quick-1'] }}</code></pre>
              <button
                @click="copyCode('quick-1')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'quick-1'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-3 flex items-center gap-2">
              <span class="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-bold">2</span>
              添加占位容器
            </h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['quick-2'] }}</code></pre>
              <button
                @click="copyCode('quick-2')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'quick-2'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-3 flex items-center gap-2">
              <span class="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-bold">3</span>
              处理回调
            </h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['quick-3'] }}</code></pre>
              <button
                @click="copyCode('quick-3')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'quick-3'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-3 flex items-center gap-2">
              <span class="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-bold">4</span>
              后端验证
            </h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['quick-4'] }}</code></pre>
              <button
                @click="copyCode('quick-4')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'quick-4'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p class="text-blue-800 text-sm flex items-start gap-2">
              <CheckCircle2 class="w-5 h-5 flex-shrink-0 mt-0.5" />
              <span>
                <strong>提示：</strong>在
                <RouterLink to="/console/sites" class="underline hover:text-blue-900">控制台</RouterLink>
                创建站点以获取 sitekey 和 secret
              </span>
            </p>
          </div>
        </div>
      </section>

      <!-- 安装指南 -->
      <section id="installation" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4 flex items-center gap-2">
          <Code class="w-6 h-6 text-indigo-500" />
          安装指南
        </h2>

        <div class="space-y-8">
          <!-- 方式一：自动渲染 -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">方式一：自动渲染（推荐）</h3>
            <p class="text-slate-600 mb-4">
              脚本加载完成后，会自动扫描页面中所有 <code class="bg-slate-100 px-1.5 py-0.5 rounded text-sm">.weauth-widget</code> 元素并渲染验证组件
            </p>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['install-1'] }}</code></pre>
              <button
                @click="copyCode('install-1')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'install-1'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- 方式二：手动渲染 -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">方式二：手动渲染</h3>
            <p class="text-slate-600 mb-4">
              使用 JavaScript API 手动控制渲染时机和位置
            </p>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['install-2'] }}</code></pre>
              <button
                @click="copyCode('install-2')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'install-2'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- 方式三：纯 iframe -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">方式三：纯 iframe（高级）</h3>
            <p class="text-slate-600 mb-4">
              直接使用 iframe 嵌入，适合需要完全控制的场景
            </p>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['install-3'] }}</code></pre>
              <button
                @click="copyCode('install-3')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'install-3'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- API 参考 -->
      <section id="api-reference" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4 flex items-center gap-2">
          <Shield class="w-6 h-6 text-emerald-500" />
          API 参考
        </h2>

        <div class="space-y-6">
          <!-- 配置选项 -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">配置选项</h3>
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-slate-200">
                    <th class="text-left py-3 px-4 font-semibold text-slate-700">参数</th>
                    <th class="text-left py-3 px-4 font-semibold text-slate-700">类型</th>
                    <th class="text-left py-3 px-4 font-semibold text-slate-700">默认值</th>
                    <th class="text-left py-3 px-4 font-semibold text-slate-700">说明</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">sitekey</code></td>
                    <td class="py-3 px-4 text-slate-600">string</td>
                    <td class="py-3 px-4 text-red-600">必填</td>
                    <td class="py-3 px-4 text-slate-600">站点公钥</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">action</code></td>
                    <td class="py-3 px-4 text-slate-600">string</td>
                    <td class="py-3 px-4 text-slate-500">-</td>
                    <td class="py-3 px-4 text-slate-600">操作标识（如 login, register）</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">theme</code></td>
                    <td class="py-3 px-4 text-slate-600">'light' | 'dark' | 'auto'</td>
                    <td class="py-3 px-4 text-slate-500">'light'</td>
                    <td class="py-3 px-4 text-slate-600">主题模式</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">themeColor</code></td>
                    <td class="py-3 px-4 text-slate-600">string</td>
                    <td class="py-3 px-4 text-slate-500">'blue'</td>
                    <td class="py-3 px-4 text-slate-600">预设主题色（blue, violet, emerald, rose, orange, cyan）</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">customColor</code></td>
                    <td class="py-3 px-4 text-slate-600">string</td>
                    <td class="py-3 px-4 text-slate-500">-</td>
                    <td class="py-3 px-4 text-slate-600">自定义颜色（#RRGGBB 格式）</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">size</code></td>
                    <td class="py-3 px-4 text-slate-600">'compact' | 'normal' | 'flexible'</td>
                    <td class="py-3 px-4 text-slate-500">'normal'</td>
                    <td class="py-3 px-4 text-slate-600">尺寸</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">autoStart</code></td>
                    <td class="py-3 px-4 text-slate-600">boolean</td>
                    <td class="py-3 px-4 text-slate-500">false</td>
                    <td class="py-3 px-4 text-slate-600">是否自动开始验证</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">callback</code></td>
                    <td class="py-3 px-4 text-slate-600">function</td>
                    <td class="py-3 px-4 text-slate-500">-</td>
                    <td class="py-3 px-4 text-slate-600">验证成功回调</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">error-callback</code></td>
                    <td class="py-3 px-4 text-slate-600">function</td>
                    <td class="py-3 px-4 text-slate-500">-</td>
                    <td class="py-3 px-4 text-slate-600">验证失败回调</td>
                  </tr>
                  <tr>
                    <td class="py-3 px-4"><code class="text-blue-600">expired-callback</code></td>
                    <td class="py-3 px-4 text-slate-600">function</td>
                    <td class="py-3 px-4 text-slate-500">-</td>
                    <td class="py-3 px-4 text-slate-600">Token 过期回调</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- JavaScript API -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">JavaScript API</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['api-js'] }}</code></pre>
              <button
                @click="copyCode('api-js')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'api-js'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
            <div class="mt-4 p-4 bg-slate-50 border border-slate-200 rounded-lg">
              <p class="text-sm text-slate-700 mb-2">
                <strong>注意：</strong>如果使用纯 iframe 方式，需要通过 postMessage 来重置：
              </p>
              <div class="relative">
                <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['api-iframe-reset'] }}</code></pre>
                <button
                  @click="copyCode('api-iframe-reset')"
                  class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                  aria-label="复制代码"
                >
                  <Check v-if="copiedId === 'api-iframe-reset'" class="w-4 h-4" />
                  <Copy v-else class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 示例代码 -->
      <section id="examples" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4 flex items-center gap-2">
          <Code class="w-6 h-6 text-blue-500" />
          示例代码
        </h2>

        <div class="space-y-8">
          <!-- 登录表单 -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">登录表单集成</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['example-login'] }}</code></pre>
              <button
                @click="copyCode('example-login')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'example-login'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- React 集成 -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">React 集成</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['example-react'] }}</code></pre>
              <button
                @click="copyCode('example-react')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'example-react'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 后端验证 -->
      <section id="backend-verify" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4 flex items-center gap-2">
          <Shield class="w-6 h-6 text-violet-500" />
          后端验证
        </h2>

        <div class="space-y-8">
          <!-- API 端点 -->
          <div class="bg-slate-50 border border-slate-200 rounded-lg p-4">
            <p class="text-sm text-slate-600 mb-2"><strong>API 端点：</strong></p>
            <code class="text-blue-600">POST https://api.weauth.cn/api/pow/siteverify</code>
          </div>

          <!-- Node.js -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">Node.js / Express</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['backend-node'] }}</code></pre>
              <button
                @click="copyCode('backend-node')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'backend-node'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Python -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">Python / Flask</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['backend-python'] }}</code></pre>
              <button
                @click="copyCode('backend-python')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'backend-python'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Go -->
          <div>
            <h3 class="text-xl font-semibold text-slate-900 mb-3">Go</h3>
            <div class="relative">
              <pre class="bg-slate-900 text-slate-300 p-4 rounded-lg text-sm overflow-x-auto"><code>{{ snippets['backend-go'] }}</code></pre>
              <button
                @click="copyCode('backend-go')"
                class="absolute top-3 right-3 p-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 transition-colors"
                aria-label="复制代码"
              >
                <Check v-if="copiedId === 'backend-go'" class="w-4 h-4" />
                <Copy v-else class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 常见问题 -->
      <section id="faq" class="bg-white rounded-2xl border border-slate-200/60 shadow-sm p-6 mb-6">
        <h2 class="text-2xl font-bold text-slate-900 mb-4">常见问题</h2>

        <div class="space-y-6">
          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: 域名白名单如何配置？</h3>
            <p class="text-slate-600">A: 在控制台的站点设置中添加允许的域名，支持子域名通配符。留空则允许所有域名（不推荐）。</p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: Token 有效期是多久？</h3>
            <p class="text-slate-600">A: 默认 60 秒，验证后立即失效（一次性）。Token 会自动缓存在浏览器中，35 秒内不会重复验证。</p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: 如何处理 Token 过期？</h3>
            <p class="text-slate-600">A: 组件会自动重置并提示用户重新验证，您也可以监听 <code class="bg-slate-100 px-1.5 py-0.5 rounded">expired-callback</code> 回调处理。</p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: 如何重置 iframe 中的验证组件？</h3>
            <p class="text-slate-600 mb-2">A: 使用 <code class="bg-slate-100 px-1.5 py-0.5 rounded">postMessage</code> 向 iframe 发送重置消息：</p>
            <pre class="bg-slate-100 p-3 rounded text-sm overflow-x-auto"><code>iframe.contentWindow.postMessage({
  type: 'weauth:reset',
  widgetId: 'your-widget-id'
}, 'https://www.weauth.cn')</code></pre>
            <p class="text-slate-600 mt-2">注意：<code class="bg-slate-100 px-1.5 py-0.5 rounded">widgetId</code> 必须与 iframe URL 参数中的 <code class="bg-slate-100 px-1.5 py-0.5 rounded">widgetId</code> 一致。</p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: iframe 被 CSP 阻止怎么办？</h3>
            <p class="text-slate-600">A: 在网站的 Content-Security-Policy 中添加 <code class="bg-slate-100 px-1.5 py-0.5 rounded">frame-src https://www.weauth.cn</code></p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: 支持离线使用吗？</h3>
            <p class="text-slate-600">A: 不支持，需要联网验证。这是工作量证明的特性，确保验证的有效性。</p>
          </div>

          <div>
            <h3 class="text-lg font-semibold text-slate-900 mb-2">Q: 如何自定义样式？</h3>
            <p class="text-slate-600">A: 可以通过 <code class="bg-slate-100 px-1.5 py-0.5 rounded">theme</code>、<code class="bg-slate-100 px-1.5 py-0.5 rounded">themeColor</code> 和 <code class="bg-slate-100 px-1.5 py-0.5 rounded">customColor</code> 参数调整主题和颜色。由于使用 iframe 隔离，无法直接修改内部样式。</p>
          </div>
        </div>
      </section>

      <!-- 底部 CTA -->
      <div class="bg-gradient-to-r from-blue-500 to-cyan-500 rounded-2xl p-8 text-center text-white">
        <h3 class="text-2xl font-bold mb-3">准备开始了吗？</h3>
        <p class="text-blue-50 mb-6">立即创建站点，获取您的 sitekey 和 secret</p>
        <RouterLink
          :to="isLoggedIn ? '/console/sites' : '/login'"
          class="inline-flex items-center gap-2 bg-white text-blue-600 px-6 py-3 rounded-lg font-semibold hover:bg-blue-50 transition-colors"
        >
          {{ isLoggedIn ? '前往控制台' : '登录开始使用' }}
        </RouterLink>
      </div>
    </main>
  </div>
</template>