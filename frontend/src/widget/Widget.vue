<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { cn } from '@/lib/utils'
import { useVerify } from './verify'
import type { WidgetVerifyConfig } from './types'
import { parseWidgetParams, isWeAuthMessage, MESSAGE_TYPES } from './utils'
import {
  Shield,
  ShieldCheck,
  ShieldAlert,
  Loader2,
  RefreshCw,
  Zap,
} from 'lucide-vue-next'

const config = ref<WidgetVerifyConfig | null>(null)
const widgetIdRef = ref('')

// 验证组合式函数
const { status, token, errorMessage, progress, startVerification, reset } = useVerify()

// 清理函数
let cleanupMessageHandler: (() => void) | null = null
let cleanupResizeObserver: (() => void) | null = null

// 主题色映射
interface ThemeColorSet {
  primary: string
  hover: string
  bg: string
  ring: string
}

const themeColors: Record<string, ThemeColorSet> = {
  blue: { primary: '#0066ff', hover: '#0052cc', bg: '#f0f5ff', ring: '#0066ff33' },
  violet: { primary: '#7c3aed', hover: '#6d28d9', bg: '#f5f3ff', ring: '#7c3aed33' },
  emerald: { primary: '#059669', hover: '#047857', bg: '#ecfdf5', ring: '#05966933' },
  rose: { primary: '#e11d48', hover: '#be123c', bg: '#fff1f2', ring: '#e11d4833' },
  orange: { primary: '#ea580c', hover: '#c2410c', bg: '#fff7ed', ring: '#ea580c33' },
  cyan: { primary: '#0891b2', hover: '#0e7490', bg: '#ecfeff', ring: '#0891b233' },
}

// 当前主题色
const colors = computed((): ThemeColorSet => {
  const tc = config.value?.themeColor ?? 'blue'
  return themeColors[tc]!
})

// 状态 UI 配置
interface StatusUiConfig {
  icon: typeof Shield | typeof ShieldCheck | typeof ShieldAlert | typeof Loader2 | typeof RefreshCw | typeof Zap
  text: string
  subtext: string
  color: string
  bgColor: string
  iconBgColor: string
}

const statusConfig = computed((): StatusUiConfig => {
  switch (status.value) {
    case 'idle':
      return {
        icon: Shield,
        text: '点击开始验证',
        subtext: '完成一次简单计算即可通过',
        color: 'text-slate-700',
        bgColor: 'bg-slate-50',
        iconBgColor: 'bg-slate-100',
      }
    case 'requesting':
      return {
        icon: Loader2,
        text: '获取挑战中...',
        subtext: '正在与服务器建立安全连接',
        color: 'text-blue-700',
        bgColor: 'bg-blue-50',
        iconBgColor: 'bg-blue-100',
      }
    case 'solving':
      return {
        icon: Zap,
        text: '计算中...',
        subtext: '正在执行工作量证明计算',
        color: 'text-amber-700',
        bgColor: 'bg-amber-50',
        iconBgColor: 'bg-amber-100',
      }
    case 'verifying':
      return {
        icon: Loader2,
        text: '验证中...',
        subtext: '正在提交验证结果',
        color: 'text-blue-700',
        bgColor: 'bg-blue-50',
        iconBgColor: 'bg-blue-100',
      }
    case 'completed':
      return {
        icon: ShieldCheck,
        text: '验证通过',
        subtext: '身份验证成功',
        color: 'text-emerald-700',
        bgColor: 'bg-emerald-50',
        iconBgColor: 'bg-emerald-100',
      }
    case 'error':
      return {
        icon: ShieldAlert,
        text: errorMessage.value || '验证失败',
        subtext: '点击重新尝试',
        color: 'text-red-700',
        bgColor: 'bg-red-50',
        iconBgColor: 'bg-red-100',
      }
    default:
      return {
        icon: Shield,
        text: '点击开始验证',
        subtext: '完成一次简单计算即可通过',
        color: 'text-slate-700',
        bgColor: 'bg-slate-50',
        iconBgColor: 'bg-slate-100',
      }
  }
})

// 是否处于等待/计算中的状态
const isBusy = computed(() =>
  status.value === 'requesting' || status.value === 'solving' || status.value === 'verifying'
)

// 是否紧凑模式
const isCompact = computed(() => config.value?.size === 'compact')

// 是否可点击
const isClickable = computed(() =>
  status.value === 'idle' || status.value === 'error'
)

/* ==================== 通信 ==================== */

function sendToParent(message: Record<string, unknown>): void {
  window.parent.postMessage(
    { ...message, widgetId: widgetIdRef.value },
    '*',
  )
}

async function handleClick(): Promise<void> {
  if (isBusy.value) return
  if (status.value === 'completed') return

  const cfg = config.value
  if (!cfg) return

  const t = await startVerification(cfg.sitekey, window.location.origin, cfg.action, cfg.difficulty)

  if (t) {
    sendToParent({ type: MESSAGE_TYPES.success, token: t })
  } else if (errorMessage.value) {
    sendToParent({ type: MESSAGE_TYPES.error, error: errorMessage.value })
  }
}

/* ==================== 生命周期 ==================== */

onMounted(() => {
  const params = parseWidgetParams()

  if (!params.sitekey) {
    console.error('WeAuth: sitekey is required')
    return
  }

  widgetIdRef.value = params.widgetId

  config.value = {
    sitekey: params.sitekey,
    action: params.action,
    difficulty: params.difficulty,
    theme: params.theme,
    themeColor: params.themeColor,
    customColor: params.customColor,
    size: params.size,
    autoStart: params.autoStart,
  }

  const handleMessage = (event: MessageEvent) => {
    if (!event.data || typeof event.data !== 'object') return
    const data = event.data as Record<string, unknown>
    if (!isWeAuthMessage(data)) return
    const { type, widgetId } = data
    if (widgetId && widgetId !== widgetIdRef.value) return
    if (type === MESSAGE_TYPES.reset) {
      reset()
    } else if (type === MESSAGE_TYPES.execute) {
      handleClick()
    }
  }

  window.addEventListener('message', handleMessage)
  cleanupMessageHandler = () => window.removeEventListener('message', handleMessage)

  const resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      const height = entry.contentRect.height
      window.parent.postMessage(
        { type: MESSAGE_TYPES.resize, widgetId: widgetIdRef.value, height: Math.ceil(height) },
        '*',
      )
    }
  })

  const widgetRoot = document.getElementById('widget-root')
  if (widgetRoot) {
    resizeObserver.observe(widgetRoot)
  }
  cleanupResizeObserver = () => resizeObserver.disconnect()

  if (params.autoStart) {
    handleClick()
  }
})

onUnmounted(() => {
  cleanupMessageHandler?.()
  cleanupResizeObserver?.()
})
</script>

<template>
  <div v-if="!config" class="flex items-center justify-center py-5 text-slate-400 text-sm">
    <Loader2 class="w-4 h-4 mr-2 animate-spin" />
    加载中...
  </div>
  <div v-else>
    <!-- ============ 紧凑模式 (size=compact) ============ -->
    <template v-if="isCompact">
      <div
        :class="cn(
          'weauth-block compact relative flex flex-col rounded-lg border transition-all duration-300 select-none overflow-hidden',
          isClickable ? 'cursor-pointer hover:shadow-sm active:scale-[0.99]' : 'cursor-default',
          status === 'idle' && 'border-slate-200 bg-white hover:border-slate-300',
          status === 'completed' && 'border-emerald-200 bg-white',
          status === 'error' && 'border-red-200 bg-white hover:border-red-300',
          isBusy && 'border-blue-200 bg-white',
        )"
        @click="handleClick"
        role="button"
        :tabindex="isClickable ? 0 : -1"
        @keydown.enter="handleClick"
        @keydown.space.prevent="handleClick"
      >
        <!-- 两行布局：图标跨行居中，右侧上下分隔 -->
        <div class="flex flex-row items-center gap-2 px-3 py-2">
          <!-- 图标 -->
          <div
            :class="cn(
              'flex items-center justify-center w-6 h-6 rounded-full transition-all duration-300 shrink-0',
              statusConfig.iconBgColor,
              isBusy && 'animate-pulse',
            )"
          >
            <Component
              :is="statusConfig.icon"
              :class="cn(
                'w-3.5 h-3.5',
                statusConfig.color,
                (status === 'requesting' || status === 'verifying') && 'animate-spin',
              )"
            />
          </div>

          <!-- 右侧两行文字 -->
          <div class="flex flex-col min-w-0 flex-1">
            <div class="flex items-center gap-1.5">
              <span
                :class="cn(
                  'text-xs font-medium leading-tight truncate',
                  status === 'completed' ? 'text-emerald-700' :
                  status === 'error' ? 'text-red-700' :
                  'text-slate-700'
                )"
              >
                {{ statusConfig.text }}
              </span>
              <!-- solving 进度文字 -->
              <span
                v-if="status === 'solving'"
                class="text-[10px] font-medium text-amber-500 shrink-0"
              >{{ progress }}%</span>
            </div>

            <span
              :class="cn(
                'text-[10px] leading-tight truncate',
                status === 'completed' ? 'text-emerald-500' :
                status === 'error' ? 'text-red-500' :
                'text-slate-400'
              )"
            >
              <Shield class="inline w-2.5 h-2.5 -mt-0.5 mr-0.5" :style="{ color: colors.primary }" />
              WeAuth
            </span>
          </div>
        </div>

        <!-- solving 进度条 -->
        <div
          v-if="status === 'solving'"
          class="absolute bottom-0 left-0 right-0 h-0.5 bg-slate-100 overflow-hidden"
        >
          <div
            class="h-full rounded-full transition-all duration-300 ease-out"
            :style="{
              width: `${progress}%`,
              background: colors.primary,
            }"
          />
        </div>
      </div>
    </template>

    <!-- ============ 普通 / 自适应模式 (size=normal|flexible) ============ -->
    <template v-else>
      <div
        :class="cn(
          'weauth-block relative flex flex-col rounded-xl border transition-all duration-300 select-none overflow-hidden',
          isClickable ? 'cursor-pointer hover:shadow-md active:scale-[0.99]' : 'cursor-default',
          status === 'idle' && 'border-slate-200 bg-white hover:border-slate-300 hover:shadow-sm',
          status === 'completed' && 'border-emerald-200 bg-white',
          status === 'error' && 'border-red-200 bg-white hover:border-red-300 hover:shadow-sm',
          isBusy && 'border-blue-200 bg-white',
        )"
        @click="handleClick"
        role="button"
        :tabindex="isClickable ? 0 : -1"
        @keydown.enter="handleClick"
        @keydown.space.prevent="handleClick"
      >
        <!-- 主交互行：图标 + 状态文字 -->
        <div
          :class="cn(
            'flex flex-row items-center gap-3 px-4 py-3 transition-colors duration-200',
            isClickable && 'hover:bg-slate-50/50',
          )"
        >
          <!-- 图标 -->
          <div
            :class="cn(
              'flex items-center justify-center w-9 h-9 rounded-full transition-all duration-300 shrink-0',
              statusConfig.iconBgColor,
              isBusy && 'animate-pulse',
            )"
          >
            <Component
              :is="statusConfig.icon"
              :class="cn(
                'w-[18px] h-[18px]',
                statusConfig.color,
                (status === 'requesting' || status === 'verifying') && 'animate-spin',
              )"
            />
          </div>

          <!-- 文字区 -->
          <div class="flex flex-col min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span
                :class="cn(
                  'text-sm font-medium leading-tight truncate',
                  status === 'completed' ? 'text-emerald-700' :
                  status === 'error' ? 'text-red-700' :
                  'text-slate-700'
                )"
              >
                {{ statusConfig.text }}
              </span>

              <!-- solving 时右侧显示进度 -->
              <span
                v-if="status === 'solving'"
                class="text-xs font-medium text-amber-500 shrink-0"
              >{{ progress }}%</span>
            </div>

            <span
              :class="cn(
                'text-[11px] leading-tight truncate',
                status === 'completed' ? 'text-emerald-500' :
                status === 'error' ? 'text-red-500' :
                'text-slate-400'
              )"
            >
              {{ statusConfig.subtext }}
            </span>
          </div>

          <!-- solving 进度条（嵌套在行右侧） -->
          <div
            v-if="status === 'solving'"
            class="absolute bottom-0 left-0 right-0 h-1 bg-slate-100 overflow-hidden"
          >
            <div
              class="h-full rounded-full transition-all duration-300 ease-out"
              :style="{
                width: `${progress}%`,
                background: colors.primary,
              }"
            />
          </div>
        </div>

        <!-- 分隔线 -->
        <div
          :class="cn(
            'mx-4 h-px',
            status === 'completed' ? 'bg-emerald-100' :
            status === 'error' ? 'bg-red-100' :
            'bg-slate-100'
          )"
        />

        <!-- 底部：WeAuth 品牌信息 -->
        <div
          :class="cn(
            'flex items-center justify-center gap-1.5 px-4 py-1.5 text-[10px]',
            status === 'completed' ? 'bg-emerald-50/50 text-emerald-500' :
            status === 'error' ? 'bg-red-50/50 text-red-500' :
            'bg-slate-50/50 text-slate-400'
          )"
        >
          <Shield class="w-2.5 h-2.5 shrink-0" :style="{ color: colors.primary }" />
          <span>由 <strong>WeAuth</strong> 提供安全验证</span>
          <template v-if="status === 'idle'">
            <span class="mx-0.5">·</span>
            <span>工作量证明 · 保护隐私</span>
          </template>
          <template v-if="status === 'completed'">
            <span class="mx-0.5">·</span>
            <span>身份已确认</span>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.weauth-block {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>