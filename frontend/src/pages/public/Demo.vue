<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { powApi } from '@/lib/api'
import { cn } from '@/lib/utils'
import {
  Shield, ShieldCheck, ShieldAlert,
  Sparkles, ArrowRight, Zap, Lock, Globe,
  Clock, Hash, Gauge, RefreshCw, Loader2,
} from 'lucide-vue-next'


const SITEKEY = 'trial-sitekey'
const SECRET = 'trial-secret'
const ORIGIN = 'http://localhost:5173'

const activeColor = ref<'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan'>('blue')

const colors = [
  { value: 'blue' as const,    label: '蓝色', class: 'bg-blue-500' },
  { value: 'violet' as const,  label: '紫色', class: 'bg-violet-500' },
  { value: 'emerald' as const, label: '翡翠', class: 'bg-emerald-500' },
  { value: 'rose' as const,    label: '玫瑰', class: 'bg-rose-500' },
  { value: 'orange' as const,  label: '橙色', class: 'bg-orange-500' },
  { value: 'cyan' as const,    label: '青色', class: 'bg-cyan-500' },
]

const themeColors: Record<
  'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'cyan', 
  { primary: string; hover: string; ring: string; text: string }
> = {
  blue:    { primary: '#0066ff', hover: '#0052cc', ring: 'rgba(0,102,255,0.25)', text: 'text-blue-700' },
  violet:  { primary: '#7c3aed', hover: '#6d28d9', ring: 'rgba(124,58,237,0.25)', text: 'text-violet-700' },
  emerald: { primary: '#059669', hover: '#047857', ring: 'rgba(5,150,105,0.25)', text: 'text-emerald-700' },
  rose:    { primary: '#e11d48', hover: '#be123c', ring: 'rgba(225,29,72,0.25)', text: 'text-rose-700' },
  orange:  { primary: '#ea580c', hover: '#c2410c', ring: 'rgba(234,88,12,0.25)', text: 'text-orange-700' },
  cyan:    { primary: '#0891b2', hover: '#0e7490', ring: 'rgba(8,145,178,0.25)', text: 'text-cyan-700' },
}

const theme = computed(() => themeColors[activeColor.value] || themeColors.blue)

type DemoStatus = 'idle' | 'requesting' | 'solving' | 'verifying' | 'completed' | 'error'
const status = ref<DemoStatus>('idle')
const errorMessage = ref<string | null>(null)
const tokenResult = ref<string | null>(null)
const siteverifyResult = ref<{ hostname: string; action: string; challenge_ts: string } | null>(null)

const iterations = ref(0)
const solveTimeMs = ref(0)
const solveSpeed = ref(0)
const cooldownRemaining = ref(0)

const hasResult = computed(() =>
  status.value === 'completed' || status.value === 'error'
)

const displayIterations = computed(() => {
  if (!hasResult.value) return '—'
  return iterations.value.toLocaleString()
})

const displayTime = computed(() => {
  if (!hasResult.value) return '—'
  return formattedTime(solveTimeMs.value)
})

const displaySpeed = computed(() => {
  if (!hasResult.value) return '—'
  return formattedSpeed(solveSpeed.value)
})

const workerRef = ref<Worker | null>(null)

function initWorker() {
  try {
    workerRef.value = new Worker(
      new URL('@/widget/pow-worker.ts', import.meta.url),
      { type: 'module' },
    )
  } catch {
    // Worker creation may fail in some environments
  }
}

initWorker()

onUnmounted(() => {
  workerRef.value?.terminate()
  if (cooldownTimer) clearInterval(cooldownTimer)
})

function solveWithWorker(
  challenge: string, difficulty: number
): Promise<{ nonce: string; iterations: number }> {
  return new Promise((resolve, reject) => {
    const w = workerRef.value
    if (!w) {
      reject(new Error('Worker 不可用'))
      return
    }

    const timeout = setTimeout(() => reject(new Error('计算超时')), 120000)

    w.onmessage = (e: MessageEvent) => {
      const data = e.data
      if (data.type === 'progress') {
        iterations.value = data.iterations ?? 0
      } else if (data.type === 'solved') {
        clearTimeout(timeout)
        resolve({ nonce: data.nonce, iterations: data.iterations })
      } else if (data.type === 'error') {
        clearTimeout(timeout)
        reject(new Error(data.message || '计算失败'))
      }
    }

    w.onerror = () => {
      clearTimeout(timeout)
      reject(new Error('Worker 内部错误'))
    }

    w.postMessage({ type: 'solve', challenge, difficulty })
  })
}

let cooldownTimer: ReturnType<typeof setInterval> | null = null

function startCooldown() {
  cooldownRemaining.value = 10
  if (cooldownTimer) clearInterval(cooldownTimer)
  cooldownTimer = setInterval(() => {
    cooldownRemaining.value--
    if (cooldownRemaining.value <= 0) {
      cooldownRemaining.value = 0
      if (cooldownTimer) clearInterval(cooldownTimer)
      cooldownTimer = null
    }
  }, 1000)
}

/* ── Verification flow ───────────────────────────────────────────── */
async function startVerification() {
  if (status.value === 'solving' || status.value === 'verifying' || status.value === 'requesting') return
  if (cooldownRemaining.value > 0) return

  /* reset */
  iterations.value = 0
  solveTimeMs.value = 0
  solveSpeed.value = 0
  errorMessage.value = null
  tokenResult.value = null
  siteverifyResult.value = null
  status.value = 'requesting'

  try {
    /* 1. get challenge */
    const challenge = await powApi.createChallenge(SITEKEY, ORIGIN, 'verify')
    if (!challenge) {
      throw new Error('后端返回挑战数据为空')
    }

    /* 2. extract challenge info (single or batch mode) */
    const isBatch = challenge.batch_mode && Array.isArray(challenge.challenges) && challenge.challenges.length > 0
    if (!challenge.challenge && !isBatch) {
      throw new Error('后端返回的挑战数据异常')
    }

    let challengeId = challenge.challenge_id

    /* 3. solve challenge(s) */
    status.value = 'solving'
    const startTime = performance.now()
    const diff = challenge.difficulty

    let allNonces: string[] = []
    let totalIters = 0

    if (isBatch) {
      /* 批量模式：依次求解全部 challenges */
      for (const c of challenge.challenges!) {
        const { nonce, iterations: iters } = await solveWithWorker(c, diff)
        allNonces.push(nonce)
        totalIters += iters
      }
    } else {
      /* 单次模式 */
      const { nonce, iterations: iters } = await solveWithWorker(challenge.challenge!, diff)
      allNonces = [nonce]
      totalIters = iters
    }

    const elapsed = performance.now() - startTime
    iterations.value = totalIters
    solveTimeMs.value = Math.round(elapsed)
    solveSpeed.value = elapsed > 0 ? Math.round(totalIters / (elapsed / 1000)) : 0

    /* 4. verify (single or batch) */
    status.value = 'verifying'
    const result = isBatch
      ? await powApi.verifyBatch(challengeId, allNonces, totalIters, Math.round(elapsed))
      : await powApi.verify(challengeId, allNonces[0]!, totalIters, Math.round(elapsed))

    const token = result.token

    /* 5. siteverify: 使用 secret + token 完成最终验证 */
    status.value = 'verifying'
    const svResult = await powApi.siteverify(SECRET, token)

    tokenResult.value = token
    siteverifyResult.value = svResult
    status.value = 'completed'
    startCooldown()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : '验证失败'
    status.value = 'error'
    startCooldown()
  }
}

const canVerify = computed(() => {
  if (status.value === 'solving' || status.value === 'verifying' || status.value === 'requesting') return false
  if (cooldownRemaining.value > 0) return false
  return true
})

const isCooldown = computed(() => cooldownRemaining.value > 0)

const buttonLabel = computed(() => {
  if (status.value === 'requesting') return '获取挑战中…'
  if (status.value === 'solving') return '计算中…'
  if (status.value === 'verifying') return '验证中…'
  if (cooldownRemaining.value > 0) return `请等待 ${cooldownRemaining.value}s`
  if (status.value === 'completed') return '再次体验'
  return '开始验证'
})

const statusIcon = computed(() => {
  switch (status.value) {
    case 'idle':       return Shield
    case 'requesting':
    case 'verifying':  return Loader2
    case 'solving':    return RefreshCw
    case 'completed':  return ShieldCheck
    case 'error':      return ShieldAlert
  }
})

const statusLabel = computed(() => {
  switch (status.value) {
    case 'idle':       return '就绪'
    case 'requesting': return '正在获取挑战…'
    case 'solving':    return '正在进行 PoW 计算…'
    case 'verifying':  return '正在提交验证…'
    case 'completed':  return '验证通过 ✓'
    case 'error':      return errorMessage.value || '验证失败'
  }
})

const statusTheme = computed(() => {
  switch (status.value) {
    case 'idle':       return { bg: 'bg-slate-100', text: 'text-slate-600', border: 'border-slate-200' }
    case 'requesting':
    case 'verifying':  return { bg: 'bg-blue-50',   text: 'text-blue-700',  border: 'border-blue-200' }
    case 'solving':    return { bg: 'bg-amber-50',  text: 'text-amber-700', border: 'border-amber-200' }
    case 'completed':  return { bg: 'bg-emerald-50',text: 'text-emerald-700',border: 'border-emerald-200' }
    case 'error':      return { bg: 'bg-red-50',    text: 'text-red-700',  border: 'border-red-200' }
  }
})

function formattedTime(ms: number): string {
  if (ms < 1000) return `${ms} ms`
  return `${(ms / 1000).toFixed(2)} s`
}

function formattedSpeed(hps: number): string {
  if (hps >= 1_000_000) return `${(hps / 1_000_000).toFixed(1)} M/s`
  if (hps >= 1_000) return `${(hps / 1_000).toFixed(1)} K/s`
  return `${hps} /s`
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30 pt-14">
    <div class="pt-24 pb-20">
      <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header -->
        <div class="text-center mb-10">
          <div class="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-200/50 rounded-full mb-6">
            <Sparkles class="w-4 h-4 text-blue-600" />
            <span class="text-sm font-medium text-blue-700">交互体验</span>
          </div>
          <h1 class="text-3xl sm:text-4xl font-bold text-slate-900 mb-4">实时体验 PoW 验证</h1>
          <p class="text-lg text-slate-600 max-w-2xl mx-auto">
            无需填写任何信息，点击下方按钮即可体验基于工作量证明的人机验证。<br />
            每次验证完成后需等待 <strong>10 秒</strong> 方可再次体验。
          </p>
        </div>

        <!-- Demo Card -->
        <div class="bg-white rounded-2xl shadow-lg border border-slate-200/60 p-6 sm:p-10">
          <!-- Color picker -->
          <div class="mb-8">
            <label class="block text-sm font-medium text-slate-700 mb-3">主题颜色</label>
            <div class="flex gap-3 flex-wrap">
              <button
                v-for="color in colors"
                :key="color.value"
                @click="activeColor = color.value"
                :class="cn(
                  'w-9 h-9 rounded-full border-2 transition-all duration-200',
                  activeColor === color.value
                    ? 'border-slate-900 scale-110 shadow-md'
                    : 'border-transparent hover:scale-110'
                )"
                :title="color.label"
              >
                <div :class="cn('w-full h-full rounded-full', color.class)" />
              </button>
            </div>
          </div>

          <!-- Status bar -->
          <div
            :class="cn(
              'flex items-center gap-3 px-4 py-3 rounded-xl border text-sm font-medium transition-all duration-300 mb-6',
              statusTheme.bg, statusTheme.text, statusTheme.border
            )"
          >
            <component
              :is="statusIcon"
              :class="cn(
                'w-5 h-5 shrink-0',
                (status === 'solving' || status === 'requesting' || status === 'verifying') && 'animate-spin'
              )"
            />
            <span>{{ statusLabel }}</span>
          </div>

          <!-- Real-time metrics (always visible) -->
          <div class="grid grid-cols-3 gap-4 mb-6">
            <!-- Iterations -->
            <div class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 text-center">
              <Hash class="w-4 h-4 mx-auto mb-1.5" :style="{ color: theme.primary }" />
              <div class="text-xs font-medium text-slate-500 mb-1">迭代次数</div>
              <div class="text-lg font-bold text-slate-900 font-mono">
                {{ displayIterations }}
              </div>
            </div>

            <!-- Time -->
            <div class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 text-center">
              <Clock class="w-4 h-4 mx-auto mb-1.5" :style="{ color: theme.primary }" />
              <div class="text-xs font-medium text-slate-500 mb-1">计算耗时</div>
              <div class="text-lg font-bold text-slate-900 font-mono">
                {{ displayTime }}
              </div>
            </div>

            <!-- Speed -->
            <div class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 text-center">
              <Gauge class="w-4 h-4 mx-auto mb-1.5" :style="{ color: theme.primary }" />
              <div class="text-xs font-medium text-slate-500 mb-1">计算速度</div>
              <div class="text-lg font-bold text-slate-900 font-mono">
                {{ displaySpeed }}
              </div>
            </div>
          </div>

          <!-- Verify button -->
          <div class="flex justify-center">
            <button
              @click="startVerification"
              :disabled="!canVerify"
              :style="{
                background: canVerify ? theme.primary : undefined,
              }"
              :class="cn(
                'relative inline-flex items-center gap-2.5 px-8 py-3 rounded-xl text-base font-semibold transition-all duration-200',
                canVerify
                  ? 'text-white cursor-pointer shadow-lg hover:shadow-xl active:scale-[0.97]'
                  : isCooldown
                    ? 'text-slate-400 bg-slate-200 cursor-not-allowed'
                    : 'text-white/60 bg-gray-400 cursor-not-allowed'
              )"
            >
              <component
                :is="statusIcon"
                :class="cn(
                  'w-5 h-5',
                  (status === 'solving' || status === 'requesting' || status === 'verifying') && 'animate-spin'
                )"
              />
              {{ buttonLabel }}
            </button>
          </div>

          <!-- Solver progress bar (only during solving) -->
          <div v-if="status === 'solving'" class="mt-4">
            <div class="h-1.5 bg-slate-200 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-300"
                :style="{ width: '70%', background: `linear-gradient(90deg, ${theme.primary}, ${theme.hover})` }"
              />
            </div>
            <p class="text-xs text-slate-400 text-center mt-1.5">正在计算中，请稍候…</p>
          </div>

          <!-- Result (success) -->
          <div
            v-if="tokenResult"
            class="mt-6 p-4 bg-emerald-50 border border-emerald-200 rounded-xl"
          >
            <div class="flex items-center gap-2 mb-2">
              <ShieldCheck class="w-5 h-5 text-emerald-600" />
              <span class="text-sm font-semibold text-emerald-700">验证成功 ✓</span>
            </div>
            <div class="text-xs text-emerald-600 font-medium mb-1.5">Token</div>
            <div class="text-xs text-emerald-600 break-all font-mono bg-white/60 rounded-lg p-2.5 border border-emerald-100/50 mb-2">
              {{ tokenResult }}
            </div>
            <div v-if="siteverifyResult" class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-emerald-600 mt-2 pt-2 border-t border-emerald-200">
              <span><strong>Hostname:</strong> {{ siteverifyResult.hostname }}</span>
              <span><strong>Action:</strong> {{ siteverifyResult.action }}</span>
              <span><strong>Challenge TS:</strong> {{ siteverifyResult.challenge_ts }}</span>
            </div>
          </div>
        </div>

        <!-- Feature grid -->
        <div class="grid sm:grid-cols-3 gap-4 mt-8">
          <div
            v-for="f in [
              { icon: Zap,  title: '零用户操作', desc: '用户无需输入验证码或点击图片' },
              { icon: Lock, title: '隐私安全',   desc: '不收集任何用户个人信息' },
              { icon: Globe,title: '跨平台',     desc: '支持所有现代浏览器和框架' },
            ]"
            :key="f.title"
            class="bg-white rounded-xl p-4 border border-slate-200/60"
          >
            <component :is="f.icon" class="w-5 h-5 text-blue-600 mb-2" />
            <h3 class="text-sm font-semibold text-slate-900 mb-1">{{ f.title }}</h3>
            <p class="text-xs text-slate-500">{{ f.desc }}</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>