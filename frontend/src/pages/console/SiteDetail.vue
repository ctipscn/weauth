<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { sitesApi } from '@/lib/api'
import { cn } from '@/lib/utils'
import { ArrowLeft, Copy, Check, Plus, Trash2, RefreshCw, Globe, Shield, TrendingUp, Activity, Save, RotateCcw, Layers, Zap, Info, Eye, AlertTriangle, XCircle, BarChart3, CheckCircle2, Loader2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const siteId = computed(() => Number(route.params.id))

const copiedField = ref<string | null>(null)
const newDomain = ref('')
const showSecretConfirm = ref(false)
const newSecret = ref<string | null>(null)
const showViewSecret = ref(false)
const viewingSecret = ref<string | null>(null)
const showDeleteConfirm = ref(false)

// Form state
const basicForm = ref<{ name: string; enabled: boolean } | null>(null)
const modeForm = ref<{
  batch_mode_enabled: boolean
  batch_difficulty: number
  min_batch_count: number
  max_batch_count: number
  base_difficulty: number
  max_difficulty: number
} | null>(null)
const policyForm = ref<{
  enabled: boolean
  time_window_seconds: number
  request_threshold: number
  batch_increment: number
  difficulty_increment: number
  failure_weight: number
} | null>(null)

const { data: site, isLoading } = useQuery({
  queryKey: ['site', siteId],
  queryFn: () => sitesApi.get(siteId.value),
  retry: false,
})

const { data: ipPolicy } = useQuery({
  queryKey: ['site-policy', siteId],
  queryFn: () => sitesApi.getIpPolicy(siteId.value),
  retry: false,
})

const { data: ipRecords } = useQuery({
  queryKey: ['site-ip-records', siteId],
  queryFn: () => sitesApi.listIpRecords(siteId.value),
  retry: false,
})

const { data: stats } = useQuery({
  queryKey: ['site-stats', siteId],
  queryFn: () => sitesApi.getStats(siteId.value),
  retry: false,
})

// Initialize forms when data loads
watch(site, (val) => {
  if (val && !basicForm.value) {
    basicForm.value = { name: val.name, enabled: val.enabled }
    modeForm.value = {
      batch_mode_enabled: val.batch_mode_enabled,
      batch_difficulty: val.batch_difficulty,
      min_batch_count: val.min_batch_count,
      max_batch_count: val.max_batch_count,
      base_difficulty: val.base_difficulty,
      max_difficulty: val.max_difficulty,
    }
  }
})

watch(ipPolicy, (val) => {
  if (val && !policyForm.value) {
    policyForm.value = {
      enabled: val.enabled,
      time_window_seconds: val.time_window_seconds,
      request_threshold: val.request_threshold,
      batch_increment: val.batch_increment || 1,
      difficulty_increment: val.difficulty_increment,
      failure_weight: val.failure_weight,
    }
  }
})

// Dirty checks
const isBasicDirty = computed(() => {
  if (!site.value || !basicForm.value) return false
  return basicForm.value.name !== site.value.name || basicForm.value.enabled !== site.value.enabled
})

const isModeDirty = computed(() => {
  if (!site.value || !modeForm.value) return false
  const m = modeForm.value
  const s = site.value
  return (
    m.batch_mode_enabled !== s.batch_mode_enabled ||
    m.batch_difficulty !== s.batch_difficulty ||
    m.min_batch_count !== s.min_batch_count ||
    m.max_batch_count !== s.max_batch_count ||
    m.base_difficulty !== s.base_difficulty ||
    m.max_difficulty !== s.max_difficulty
  )
})

const isPolicyDirty = computed(() => {
  if (!ipPolicy.value || !policyForm.value) return false
  const p = policyForm.value
  const ip = ipPolicy.value
  return (
    p.enabled !== ip.enabled ||
    p.time_window_seconds !== ip.time_window_seconds ||
    p.request_threshold !== ip.request_threshold ||
    p.batch_increment !== (ip.batch_increment || 1) ||
    p.difficulty_increment !== ip.difficulty_increment ||
    p.failure_weight !== ip.failure_weight
  )
})

// Reset forms
const resetBasicForm = () => {
  if (site.value) {
    basicForm.value = { name: site.value.name, enabled: site.value.enabled }
  }
}

const resetModeForm = () => {
  if (site.value) {
    const s = site.value
    modeForm.value = {
      batch_mode_enabled: s.batch_mode_enabled,
      batch_difficulty: s.batch_difficulty,
      min_batch_count: s.min_batch_count,
      max_batch_count: s.max_batch_count,
      base_difficulty: s.base_difficulty,
      max_difficulty: s.max_difficulty,
    }
  }
}

const resetPolicyForm = () => {
  if (ipPolicy.value) {
    const ip = ipPolicy.value
    policyForm.value = {
      enabled: ip.enabled,
      time_window_seconds: ip.time_window_seconds,
      request_threshold: ip.request_threshold,
      batch_increment: ip.batch_increment || 1,
      difficulty_increment: ip.difficulty_increment,
      failure_weight: ip.failure_weight,
    }
  }
}

// Mutations
const updateMutation = useMutation({
  mutationFn: (data: any) => sitesApi.update(siteId.value, data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['site', siteId] })
  },
})

const updatePolicyMutation = useMutation({
  mutationFn: (data: any) => sitesApi.updateIpPolicy(siteId.value, data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['site-policy', siteId] })
  },
})

const addDomainMutation = useMutation({
  mutationFn: (domain: string) => sitesApi.addDomain(siteId.value, domain),
  onSuccess: () => {
    newDomain.value = ''
    queryClient.invalidateQueries({ queryKey: ['site', siteId] })
  },
})

const deleteDomainMutation = useMutation({
  mutationFn: (domainId: number) => sitesApi.deleteDomain(siteId.value, domainId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['site', siteId] })
  },
})

const regenerateSecretMutation = useMutation({
  mutationFn: () => sitesApi.regenerateSecret(siteId.value),
  onSuccess: (data) => {
    newSecret.value = data.secret
    showSecretConfirm.value = false
  },
})

const viewSecretMutation = useMutation({
  mutationFn: () => sitesApi.getSecret(siteId.value),
  onSuccess: (data) => {
    viewingSecret.value = data.secret
    showViewSecret.value = false
  },
})

const deleteIpRecordMutation = useMutation({
  mutationFn: (ip: string) => sitesApi.deleteIpRecord(siteId.value, ip),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['site-ip-records', siteId] })
  },
})

const deleteSiteMutation = useMutation({
  mutationFn: () => sitesApi.delete(siteId.value),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['sites'] })
    router.push('/console/sites')
  },
})

const handleSaveBasic = async () => {
  if (!basicForm.value) return
  await updateMutation.mutateAsync(basicForm.value)
}

const handleSaveMode = async () => {
  if (!modeForm.value) return
  await updateMutation.mutateAsync(modeForm.value)
}

const handleSavePolicy = async () => {
  if (!policyForm.value) return
  await updatePolicyMutation.mutateAsync(policyForm.value)
}

const handleCopy = async (text: string, field: string) => {
  await navigator.clipboard.writeText(text)
  copiedField.value = field
  setTimeout(() => { copiedField.value = null }, 2000)
}

const successRate = computed(() => {
  if (!stats.value) return 0
  return stats.value.total_challenges > 0
    ? Math.round((stats.value.solved_challenges / stats.value.total_challenges) * 100)
    : 0
})

const toggleChecked = (val: any) => !val
</script>

<template>
  <div v-if="isLoading" class="flex items-center justify-center h-64">
    <Loader2 class="animate-spin h-6 w-6 text-violet-600" />
  </div>

  <div v-else-if="!site" class="text-center py-12">
    <p class="text-slate-500">站点不存在</p>
  </div>

  <div v-else class="max-w-4xl mx-auto space-y-8">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <button @click="router.push('/console/sites')" class="p-2 hover:bg-slate-100 rounded-lg transition-colors">
        <ArrowLeft class="h-5 w-5 text-slate-600" />
      </button>
      <div class="flex-1">
        <div class="flex items-center gap-3">
          <h1 class="text-xl font-semibold text-slate-900">{{ site.name }}</h1>
          <span
            :class="cn(
              'px-2 py-0.5 text-xs rounded-full font-medium',
              site.enabled ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-500'
            )"
          >
            {{ site.enabled ? '运行中' : '已禁用' }}
          </span>
        </div>
        <div class="flex items-center gap-2 mt-1">
          <code class="text-xs text-slate-500 bg-slate-100 px-2 py-0.5 rounded font-mono">
            {{ site.sitekey }}
          </code>
          <button @click="handleCopy(site.sitekey, 'sitekey')" class="text-slate-400 hover:text-slate-600 transition-colors">
            <Check v-if="copiedField === 'sitekey'" class="h-3.5 w-3.5 text-emerald-500" />
            <Copy v-else class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Stats -->
    <div v-if="stats" class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-slate-200 p-4">
        <div class="flex items-center gap-2 mb-2">
          <div class="p-1.5 rounded-lg bg-blue-50">
            <Activity class="h-4 w-4 text-blue-500" />
          </div>
          <span class="text-xs text-slate-500">24h 请求</span>
        </div>
        <div class="text-2xl font-semibold text-slate-900">{{ stats.last_24h_challenges }}</div>
      </div>
      <div class="bg-white rounded-xl border border-slate-200 p-4">
        <div class="flex items-center gap-2 mb-2">
          <div class="p-1.5 rounded-lg bg-emerald-50">
            <CheckCircle2 class="h-4 w-4 text-emerald-500" />
          </div>
          <span class="text-xs text-slate-500">24h 成功</span>
        </div>
        <div class="text-2xl font-semibold text-slate-900">{{ stats.last_24h_solved }}</div>
      </div>
      <div class="bg-white rounded-xl border border-slate-200 p-4">
        <div class="flex items-center gap-2 mb-2">
          <div class="p-1.5 rounded-lg bg-purple-50">
            <BarChart3 class="h-4 w-4 text-purple-500" />
          </div>
          <span class="text-xs text-slate-500">总请求</span>
        </div>
        <div class="text-2xl font-semibold text-slate-900">{{ stats.total_challenges }}</div>
      </div>
      <div class="bg-white rounded-xl border border-slate-200 p-4">
        <div class="flex items-center gap-2 mb-2">
          <div class="p-1.5 rounded-lg bg-amber-50">
            <TrendingUp class="h-4 w-4 text-amber-500" />
          </div>
          <span class="text-xs text-slate-500">成功率</span>
        </div>
        <div class="text-2xl font-semibold text-slate-900">{{ successRate }}%</div>
      </div>
    </div>

    <!-- SVG Bar Chart -->
    <div v-if="stats" class="bg-white rounded-xl border border-slate-200 p-6">
      <h3 class="text-sm font-semibold text-slate-900 mb-4 flex items-center gap-2">
        <BarChart3 class="w-4 h-4 text-blue-500" />
        验证统计概览
      </h3>
      <div class="flex items-end justify-center gap-8 h-40">
        <div class="flex flex-col items-center gap-2">
          <div class="text-xs font-medium text-slate-600">{{ stats.total_challenges }}</div>
          <div class="w-12 bg-blue-500 rounded-t-md transition-all" :style="{ height: Math.min(120, (stats.total_challenges / Math.max(stats.total_challenges, 1)) * 120) + 'px' }"></div>
          <span class="text-xs text-slate-500">总挑战</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <div class="text-xs font-medium text-slate-600">{{ stats.solved_challenges }}</div>
          <div class="w-12 bg-emerald-500 rounded-t-md transition-all" :style="{ height: Math.min(120, (stats.solved_challenges / Math.max(stats.total_challenges, 1)) * 120) + 'px' }"></div>
          <span class="text-xs text-slate-500">已解决</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <div class="text-xs font-medium text-slate-600">{{ stats.last_24h_challenges }}</div>
          <div class="w-12 bg-purple-500 rounded-t-md transition-all" :style="{ height: Math.min(120, (stats.last_24h_challenges / Math.max(stats.total_challenges, 1)) * 120) + 'px' }"></div>
          <span class="text-xs text-slate-500">24h</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <div class="text-xs font-medium text-slate-600">{{ successRate }}%</div>
          <div class="w-12 bg-amber-500 rounded-t-md transition-all" :style="{ height: Math.min(120, (successRate / 100) * 120) + 'px' }"></div>
          <span class="text-xs text-slate-500">成功率</span>
        </div>
      </div>
    </div>

    <!-- 基本设置 -->
    <div v-if="basicForm" class="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h2 class="font-semibold text-slate-900">基本设置</h2>
        </div>
        <div v-if="isBasicDirty" class="flex items-center gap-2">
          <button @click="resetBasicForm" class="h-8 px-3 text-sm text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors flex items-center gap-1.5">
            <RotateCcw class="w-3.5 h-3.5" />
            取消
          </button>
          <button @click="handleSaveBasic" :disabled="updateMutation.isPending.value" class="h-8 px-3 text-sm bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            <Save class="w-3.5 h-3.5" />
            {{ updateMutation.isPending.value ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
      <div class="px-6">
        <div class="flex items-center justify-between py-4 border-b border-slate-100">
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate-900">站点名称</div>
            <div class="text-xs text-slate-500 mt-0.5">用于标识此站点</div>
          </div>
          <div class="flex-shrink-0">
            <input
              v-model="basicForm.name"
              type="text"
              class="w-48 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
            />
          </div>
        </div>
        <div class="flex items-center justify-between py-4 border-b border-slate-100">
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate-900">启用状态</div>
            <div class="text-xs text-slate-500 mt-0.5">关闭后将拒绝所有验证请求</div>
          </div>
          <div class="flex-shrink-0">
            <button
              type="button"
              role="switch"
              :aria-checked="basicForm.enabled"
              @click="basicForm.enabled = !basicForm.enabled"
              :class="cn(
                'relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer',
                basicForm.enabled ? 'bg-emerald-600' : 'bg-gray-200'
              )"
            >
              <span
                :class="cn(
                  'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-lg ring-0 transition-transform duration-200 ease-in-out',
                  basicForm.enabled ? 'translate-x-5' : 'translate-x-0'
                )"
              />
            </button>
          </div>
        </div>
        <div class="flex items-center justify-between py-4">
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate-900">Secret 密钥</div>
            <div class="text-xs text-slate-500 mt-0.5">用于后端验证，请妥善保管</div>
          </div>
          <div class="flex-shrink-0 flex items-center gap-2">
            <button @click="showViewSecret = true" class="h-8 px-3 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors flex items-center gap-1.5">
              <Eye class="w-3.5 h-3.5" />
              查看
            </button>
            <button @click="showSecretConfirm = true" class="h-8 px-3 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors flex items-center gap-1.5">
              <RefreshCw class="w-3.5 h-3.5" />
              重新生成
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 验证模式 -->
    <div v-if="modeForm" class="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h2 class="font-semibold text-slate-900">验证模式</h2>
          <p class="text-xs text-slate-500 mt-1">选择 PoW 验证的计算策略</p>
        </div>
        <div v-if="isModeDirty" class="flex items-center gap-2">
          <button @click="resetModeForm" class="h-8 px-3 text-sm text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors flex items-center gap-1.5">
            <RotateCcw class="w-3.5 h-3.5" />
            取消
          </button>
          <button @click="handleSaveMode" :disabled="updateMutation.isPending.value" class="h-8 px-3 text-sm bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            <Save class="w-3.5 h-3.5" />
            {{ updateMutation.isPending.value ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
      <div class="px-6">
        <div class="flex items-center justify-between py-4 border-b border-slate-100">
          <div class="flex-1 pr-4">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-slate-900">批量模式</span>
              <span class="px-1.5 py-0.5 text-xs rounded bg-emerald-100 text-emerald-700 font-medium">推荐</span>
            </div>
            <div class="text-xs text-slate-500 mt-0.5">固定难度，根据风控动态调整计算数量</div>
          </div>
          <div class="flex-shrink-0">
            <button
              type="button"
              role="switch"
              :aria-checked="modeForm.batch_mode_enabled"
              @click="modeForm.batch_mode_enabled = !modeForm.batch_mode_enabled"
              :class="cn(
                'relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer',
                modeForm.batch_mode_enabled ? 'bg-emerald-600' : 'bg-gray-200'
              )"
            >
              <span
                :class="cn(
                  'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-lg ring-0 transition-transform duration-200 ease-in-out',
                  modeForm.batch_mode_enabled ? 'translate-x-5' : 'translate-x-0'
                )"
              />
            </button>
          </div>
        </div>

        <div class="py-4 border-b border-slate-100">
          <div :class="cn('flex items-start gap-3 p-3 rounded-lg', modeForm.batch_mode_enabled ? 'bg-purple-50' : 'bg-slate-50')">
            <div :class="cn('p-1.5 rounded-lg', modeForm.batch_mode_enabled ? 'bg-purple-100' : 'bg-slate-200')">
              <Layers v-if="modeForm.batch_mode_enabled" class="h-4 w-4 text-purple-600" />
              <Zap v-else class="h-4 w-4 text-amber-600" />
            </div>
            <div class="flex-1 text-sm">
              <div class="font-medium text-slate-900">
                {{ modeForm.batch_mode_enabled ? '批量模式' : '动态难度模式' }}
              </div>
              <p class="text-slate-600 mt-0.5">
                {{ modeForm.batch_mode_enabled
                  ? `客户端计算 ${modeForm.min_batch_count}-${modeForm.max_batch_count} 个难度 ${modeForm.batch_difficulty} 的 PoW，风险越高数量越多`
                  : `客户端计算 1 个 PoW，难度范围 ${modeForm.base_difficulty}-${modeForm.max_difficulty}，风险越高难度越大`
                }}
              </p>
            </div>
          </div>
        </div>

        <template v-if="modeForm.batch_mode_enabled">
          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">批量难度</div>
              <div class="text-xs text-slate-500 mt-0.5">每个 PoW 的计算难度</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm.batch_difficulty"
                @input="modeForm.batch_difficulty = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="8"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">最小数量</div>
              <div class="text-xs text-slate-500 mt-0.5">低风险时的 PoW 数量</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm.min_batch_count"
                @input="modeForm.min_batch_count = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="10"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">最大数量</div>
              <div class="text-xs text-slate-500 mt-0.5">高风险时的 PoW 数量</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm.max_batch_count"
                @input="modeForm.max_batch_count = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="20"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
        </template>
        <template v-else>
          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">基础难度</div>
              <div class="text-xs text-slate-500 mt-0.5">低风险时的难度等级</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm.base_difficulty"
                @input="modeForm.base_difficulty = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="10"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">最大难度</div>
              <div class="text-xs text-slate-500 mt-0.5">高风险时的难度等级</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm.max_difficulty"
                @input="modeForm.max_difficulty = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="12"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 风控策略 -->
    <div v-if="policyForm" class="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h2 class="font-semibold text-slate-900">风控策略</h2>
          <p class="text-xs text-slate-500 mt-1">根据 IP 请求频率动态调整验证强度</p>
        </div>
        <div v-if="isPolicyDirty" class="flex items-center gap-2">
          <button @click="resetPolicyForm" class="h-8 px-3 text-sm text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors flex items-center gap-1.5">
            <RotateCcw class="w-3.5 h-3.5" />
            取消
          </button>
          <button @click="handleSavePolicy" :disabled="updatePolicyMutation.isPending.value" class="h-8 px-3 text-sm bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            <Save class="w-3.5 h-3.5" />
            {{ updatePolicyMutation.isPending.value ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
      <div class="px-6">
        <div class="flex items-center justify-between py-4 border-b border-slate-100">
          <div class="flex-1 pr-4">
            <div class="text-sm font-medium text-slate-900">启用风控</div>
            <div class="text-xs text-slate-500 mt-0.5">关闭后所有请求使用基础配置</div>
          </div>
          <div class="flex-shrink-0">
            <button
              type="button"
              role="switch"
              :aria-checked="policyForm.enabled"
              @click="policyForm.enabled = !policyForm.enabled"
              :class="cn(
                'relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer',
                policyForm.enabled ? 'bg-emerald-600' : 'bg-gray-200'
              )"
            >
              <span
                :class="cn(
                  'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-lg ring-0 transition-transform duration-200 ease-in-out',
                  policyForm.enabled ? 'translate-x-5' : 'translate-x-0'
                )"
              />
            </button>
          </div>
        </div>

        <template v-if="policyForm.enabled">
          <div class="py-4 border-b border-slate-100">
            <div class="flex items-start gap-2 p-3 rounded-lg bg-amber-50 text-amber-800">
              <Info class="h-4 w-4 mt-0.5 flex-shrink-0" />
              <p class="text-sm">
                {{ modeForm?.batch_mode_enabled
                  ? `同一 IP 在 ${policyForm.time_window_seconds} 秒内每 ${policyForm.request_threshold} 次请求，增加 ${policyForm.batch_increment} 个计算任务`
                  : `同一 IP 在 ${policyForm.time_window_seconds} 秒内每 ${policyForm.request_threshold} 次请求，增加 ${policyForm.difficulty_increment} 级难度`
                }}
              </p>
            </div>
          </div>

          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">时间窗口</div>
              <div class="text-xs text-slate-500 mt-0.5">统计请求的时间范围（秒）</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="policyForm.time_window_seconds"
                @input="policyForm.time_window_seconds = Number(($event.target as HTMLInputElement).value)"
                type="number"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">请求阈值</div>
              <div class="text-xs text-slate-500 mt-0.5">触发增量的请求次数</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="policyForm.request_threshold"
                @input="policyForm.request_threshold = Number(($event.target as HTMLInputElement).value)"
                type="number"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4 border-b border-slate-100">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">{{ modeForm?.batch_mode_enabled ? '批次增量' : '难度增量' }}</div>
              <div class="text-xs text-slate-500 mt-0.5">{{ modeForm?.batch_mode_enabled ? '每次触发增加的 PoW 数量' : '每次触发增加的难度等级' }}</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="modeForm?.batch_mode_enabled ? policyForm.batch_increment : policyForm.difficulty_increment"
                @input="modeForm?.batch_mode_enabled ? (policyForm.batch_increment = Number(($event.target as HTMLInputElement).value)) : (policyForm.difficulty_increment = Number(($event.target as HTMLInputElement).value))"
                type="number"
                min="1"
                max="5"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
          <div class="flex items-center justify-between py-4">
            <div class="flex-1 pr-4">
              <div class="text-sm font-medium text-slate-900">失败权重</div>
              <div class="text-xs text-slate-500 mt-0.5">失败请求的计数权重</div>
            </div>
            <div class="flex-shrink-0">
              <input
                :value="policyForm.failure_weight"
                @input="policyForm.failure_weight = Number(($event.target as HTMLInputElement).value)"
                type="number"
                min="1"
                max="10"
                step="0.1"
                class="w-24 h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
              />
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 域名白名单 -->
    <section class="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-100">
        <h2 class="font-semibold text-slate-900">域名白名单</h2>
        <p class="text-xs text-slate-500 mt-1">留空允许所有域名，支持通配符 *.example.com</p>
      </div>
      <div class="px-6 py-4">
        <div class="flex gap-2 mb-4">
          <input
            v-model="newDomain"
            type="text"
            placeholder="输入域名，按回车添加"
            @keydown.enter="newDomain && addDomainMutation.mutate(newDomain)"
            class="flex-1 h-9 px-3 rounded-lg border border-slate-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
          />
          <button
            @click="newDomain && addDomainMutation.mutate(newDomain)"
            :disabled="!newDomain || addDomainMutation.isPending.value"
            class="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
          >
            <Plus class="w-4 h-4" />
            添加
          </button>
        </div>

        <div v-if="site.domains && site.domains.length > 0" class="space-y-2">
          <div
            v-for="domain in site.domains"
            :key="domain.id"
            class="flex items-center justify-between px-3 py-2.5 bg-slate-50 rounded-lg group hover:bg-slate-100 transition-colors"
          >
            <div class="flex items-center gap-2">
              <Globe class="h-4 w-4 text-slate-400" />
              <code class="text-sm text-slate-700">{{ domain.domain }}</code>
            </div>
            <button
              @click="deleteDomainMutation.mutate(domain.id)"
              class="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded opacity-0 group-hover:opacity-100 transition-all"
            >
              <Trash2 class="h-4 w-4" />
            </button>
          </div>
        </div>
        <div v-else class="text-center py-8 text-slate-400 text-sm">
          未配置域名限制，允许所有来源
        </div>
      </div>
    </section>

    <!-- IP 记录 -->
    <section class="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-100">
        <h2 class="font-semibold text-slate-900">IP 访问记录</h2>
        <p class="text-xs text-slate-500 mt-1">监控各 IP 的请求状态和风控级别</p>
      </div>
      <div class="overflow-x-auto">
        <table v-if="ipRecords && ipRecords.length > 0" class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-100">
              <th class="text-left py-3 px-6 font-medium text-slate-500">IP 地址</th>
              <th class="text-left py-3 px-4 font-medium text-slate-500">请求数</th>
              <th class="text-left py-3 px-4 font-medium text-slate-500">失败数</th>
              <th class="text-left py-3 px-4 font-medium text-slate-500">当前难度</th>
              <th class="text-left py-3 px-4 font-medium text-slate-500">最后请求</th>
              <th class="text-right py-3 px-6 font-medium text-slate-500">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="record in ipRecords"
              :key="record.id"
              class="border-b border-slate-50 hover:bg-slate-50/50 transition-colors"
            >
              <td class="py-3 px-6 font-mono text-slate-700">{{ record.client_ip }}</td>
              <td class="py-3 px-4 text-slate-600">{{ record.request_count }}</td>
              <td class="py-3 px-4">
                <span :class="record.failure_count > 0 ? 'text-rose-600 font-medium' : 'text-slate-600'">
                  {{ record.failure_count }}
                </span>
              </td>
              <td class="py-3 px-4">
                <span class="px-2 py-1 bg-blue-50 text-blue-700 rounded text-xs font-medium">
                  {{ record.current_difficulty }}
                </span>
              </td>
              <td class="py-3 px-4 text-slate-500 text-xs">
                {{ new Date(record.last_request_at).toLocaleString('zh-CN') }}
              </td>
              <td class="py-3 px-6 text-right">
                <button
                  @click="deleteIpRecordMutation.mutate(record.client_ip)"
                  class="text-xs text-rose-600 hover:text-rose-700 font-medium"
                >
                  重置
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="text-center py-12 text-slate-400 text-sm">
          暂无访问记录
        </div>
      </div>
    </section>

    <!-- Danger Zone -->
    <section class="bg-white rounded-xl border border-red-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-red-100 bg-red-50">
        <div class="flex items-center gap-2">
          <AlertTriangle class="w-5 h-5 text-red-500" />
          <h2 class="font-semibold text-red-700">危险区域</h2>
        </div>
        <p class="text-xs text-red-500 mt-1">删除操作不可恢复，请谨慎操作</p>
      </div>
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium text-slate-900">删除站点</div>
            <p class="text-xs text-slate-500 mt-0.5">永久删除此站点及其所有配置</p>
          </div>
          <button
            @click="showDeleteConfirm = true"
            class="h-9 px-4 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-1.5"
          >
            <Trash2 class="w-4 h-4" />
            删除站点
          </button>
        </div>
      </div>
    </section>

    <!-- Secret Confirm Modal -->
    <div v-if="showSecretConfirm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-100">
          <h2 class="font-semibold text-slate-900 flex items-center gap-2">
            <XCircle class="w-5 h-5 text-red-500" />
            确认重新生成 Secret
          </h2>
        </div>
        <div class="p-6">
          <p class="text-slate-600 text-sm mb-6">
            重新生成后，旧的 Secret 将立即失效。所有使用旧 Secret 的服务需要更新配置。
          </p>
          <div class="flex justify-end gap-3">
            <button @click="showSecretConfirm = false" class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium">
              取消
            </button>
            <button
              @click="regenerateSecretMutation.mutate()"
              :disabled="regenerateSecretMutation.isPending.value"
              class="h-9 px-4 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm font-medium disabled:opacity-50"
            >
              {{ regenerateSecretMutation.isPending.value ? '处理中...' : '确认重新生成' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- View Secret Confirm Modal -->
    <div v-if="showViewSecret" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-100">
          <h2 class="font-semibold text-slate-900 flex items-center gap-2">
            <Shield class="w-5 h-5 text-amber-500" />
            确认查看 Secret
          </h2>
        </div>
        <div class="p-6">
          <p class="text-slate-600 text-sm mb-6">
            Secret 是敏感信息，请确保周围环境安全。建议在私密场所查看。
          </p>
          <div class="flex justify-end gap-3">
            <button @click="showViewSecret = false" class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium">
              取消
            </button>
            <button
              @click="viewSecretMutation.mutate()"
              :disabled="viewSecretMutation.isPending.value"
              class="h-9 px-4 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ viewSecretMutation.isPending.value ? '处理中...' : '确认查看' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Viewing Secret Modal -->
    <div v-if="viewingSecret" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-amber-100 bg-amber-50">
          <h2 class="font-semibold text-amber-700 flex items-center gap-2">
            <Shield class="w-5 h-5" />
            站点 Secret
          </h2>
        </div>
        <div class="p-6 space-y-4">
          <div class="flex items-center gap-2 p-3 rounded-lg bg-amber-50 border border-amber-200 text-amber-700 text-sm">
            <AlertTriangle class="w-4 h-4 flex-shrink-0" />
            <span>请妥善保管此密钥，不要泄露给他人</span>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">Secret 密钥</label>
            <div class="flex items-center gap-2">
              <code class="flex-1 px-3 py-2.5 bg-slate-100 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                {{ viewingSecret }}
              </code>
              <button
                @click="handleCopy(viewingSecret, 'viewing-secret')"
                class="h-10 w-10 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
              >
                <Check v-if="copiedField === 'viewing-secret'" class="h-5 w-5 text-emerald-500" />
                <Copy v-else class="h-5 w-5" />
              </button>
            </div>
          </div>
          <div class="flex justify-end pt-2">
            <button @click="viewingSecret = null" class="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white rounded-lg transition-colors text-sm font-medium">
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- New Secret Modal -->
    <div v-if="newSecret" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-violet-100 bg-gradient-to-r from-violet-50 to-purple-50">
          <h2 class="font-semibold text-violet-700 flex items-center gap-2">
            <CheckCircle2 class="w-5 h-5" />
            Secret 已重新生成
          </h2>
        </div>
        <div class="p-6 space-y-4">
          <div class="flex items-center gap-2 p-3 rounded-lg bg-blue-50 border border-blue-200 text-blue-700 text-sm">
            <Shield class="w-4 h-4 flex-shrink-0" />
            <span>新的 Secret 已生成，可随时在此页面查看</span>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">新 Secret</label>
            <div class="flex items-center gap-2">
              <code class="flex-1 px-3 py-2.5 bg-slate-100 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                {{ newSecret }}
              </code>
              <button
                @click="handleCopy(newSecret, 'new-secret')"
                class="h-10 w-10 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
              >
                <Check v-if="copiedField === 'new-secret'" class="h-5 w-5 text-emerald-500" />
                <Copy v-else class="h-5 w-5" />
              </button>
            </div>
          </div>
          <div class="flex justify-end pt-2">
            <button @click="newSecret = null" class="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white rounded-lg transition-colors text-sm font-medium">
              我已保存
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-red-100 bg-red-50">
          <h2 class="font-semibold text-red-700 flex items-center gap-2">
            <AlertTriangle class="w-5 h-5" />
            确认删除站点
          </h2>
        </div>
        <div class="p-6">
          <p class="text-slate-600 text-sm mb-6">
            此操作不可恢复。站点 <strong class="text-slate-900">{{ site.name }}</strong> 及其所有配置和记录将永久删除。
          </p>
          <div class="flex justify-end gap-3">
            <button @click="showDeleteConfirm = false" class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium">
              取消
            </button>
            <button
              @click="deleteSiteMutation.mutate()"
              :disabled="deleteSiteMutation.isPending.value"
              class="h-9 px-4 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm font-medium disabled:opacity-50"
            >
              {{ deleteSiteMutation.isPending.value ? '删除中...' : '确认删除' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>