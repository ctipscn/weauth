<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { sitesApi } from '@/lib/api'
import type { Site } from '@/lib/api'
import { 
  Globe, Plus, Search, Settings, 
  Trash2, Copy, Check, AlertCircle, 
  Loader2, Shield, Zap, Server
} from 'lucide-vue-next'

import { cn } from '@/lib/utils'

const router = useRouter()
const queryClient = useQueryClient()

const showCreate = ref(false)
const newSiteName = ref('')
const newDomains = ref('')
const batchModeEnabled = ref(true)
const batchDifficulty = ref(4)
const createdSite = ref<(Site & { secret: string }) | null>(null)
const copiedField = ref<string | null>(null)
const searchQuery = ref('')

const { data: sites, isLoading } = useQuery({
  queryKey: ['sites'],
  queryFn: sitesApi.list,
  retry: false,
})

const filteredSites = computed(() => {
  if (!sites.value) return []
  if (!searchQuery.value.trim()) return sites.value
  const q = searchQuery.value.toLowerCase()
  return sites.value.filter(
    (s: any) =>
      s.name.toLowerCase().includes(q) ||
      s.sitekey.toLowerCase().includes(q)
  )
})

const createMutation = useMutation({
  mutationFn: () =>
    sitesApi.create(
      newSiteName.value,
      newDomains.value.split('\n').filter((d: string) => d.trim()),
      {
        batch_mode_enabled: batchModeEnabled.value,
        batch_difficulty: batchDifficulty.value,
      }
    ),
  onSuccess: (site: any) => {
    createdSite.value = site
    newSiteName.value = ''
    newDomains.value = ''
    batchModeEnabled.value = true
    batchDifficulty.value = 4
    showCreate.value = false
    queryClient.invalidateQueries({ queryKey: ['sites'] })
  },
})

const deleteMutation = useMutation({
  mutationFn: (id: number) => sitesApi.delete(id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['sites'] })
  },
})

const handleCopy = async (text: string, field: string) => {
  await navigator.clipboard.writeText(text)
  copiedField.value = field
  setTimeout(() => { copiedField.value = null }, 2000)
}

const handleCreate = (e: Event) => {
  e.preventDefault()
  createMutation.mutate()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">站点管理</h1>
        <p class="text-slate-500 text-sm mt-0.5">管理您的 PoW 验证站点和配置</p>
      </div>
      <button
        @click="showCreate = true"
        class="inline-flex items-center justify-center h-9 px-4 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 active:bg-blue-100 text-sm rounded-lg font-medium transition-all duration-150 shadow-sm"
      >
        <Plus class="h-4 w-4 mr-1.5" />
        创建站点
      </button>
    </div>

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索站点名称或 Sitekey..."
        class="w-full pl-10 pr-4 h-9 bg-white/80 border border-slate-200 rounded-lg text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
      />
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
    </div>

    <!-- Empty State -->
    <div v-else-if="!isLoading && sites?.length === 0" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
      <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
        <Globe class="w-6 h-6 text-slate-400" />
      </div>
      <h3 class="text-lg font-medium text-slate-900 mb-2">还没有创建站点</h3>
      <p class="text-slate-500 text-sm mb-6">创建您的第一个站点，开始使用 WeAuth 人机验证服务</p>
      <button
        @click="showCreate = true"
        class="inline-flex items-center h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg font-medium transition-colors"
      >
        <Plus class="h-4 w-4 mr-1.5" />
        创建站点
      </button>
    </div>

    <!-- No Results -->
    <div v-else-if="filteredSites.length === 0" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
      <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
        <Search class="w-6 h-6 text-slate-400" />
      </div>
      <h3 class="text-lg font-medium text-slate-900 mb-2">未找到匹配的站点</h3>
      <p class="text-slate-500 text-sm">尝试使用不同的搜索词</p>
    </div>

    <!-- Site List -->
    <div v-else class="space-y-3">
      <div
        v-for="site in filteredSites"
        :key="site.id"
        @click="router.push(`/console/sites/${site.id}`)"
        class="group bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md hover:border-slate-300/80 transition-all p-5 cursor-pointer"
      >
        <div class="flex items-start gap-4">
          <!-- Site Icon -->
          <div
            :class="cn(
              'w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0',
              site.enabled ? 'bg-violet-600 shadow-sm' : 'bg-slate-200'
            )"
          >
            <svg class="w-5 h-5" :class="site.enabled ? 'text-white' : 'text-slate-400'" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
            </svg>
          </div>

          <!-- Site Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-3 mb-1.5">
              <h3 class="font-semibold text-slate-900 truncate">{{ site.name }}</h3>
              <span
                :class="cn(
                  'px-2 py-0.5 text-xs rounded-full font-medium',
                  site.enabled
                    ? 'bg-emerald-100 text-emerald-700'
                    : 'bg-slate-100 text-slate-500'
                )"
              >
                {{ site.enabled ? '运行中' : '已禁用' }}
              </span>
            </div>

            <div class="flex items-center gap-4 text-xs text-slate-500">
              <div class="flex items-center gap-1.5">
                <Shield class="w-3.5 h-3.5" />
                <code class="font-mono bg-slate-100 px-1.5 py-0.5 rounded">{{ site.sitekey }}</code>
                <button
                  @click.stop="handleCopy(site.sitekey, `sitekey-${site.id}`)"
                  class="p-0.5 hover:text-slate-700 transition-colors"
                >
                  <Check v-if="copiedField === `sitekey-${site.id}`" class="w-3 h-3 text-emerald-500" />
                  <Copy v-else class="w-3 h-3" />
                </button>
              </div>
              <div class="flex items-center gap-1">
                <Zap class="w-3.5 h-3.5" />
                <span>难度 {{ site.batch_mode_enabled ? site.batch_difficulty : site.base_difficulty }}</span>
              </div>
              <div class="flex items-center gap-1">
                <Server class="w-3.5 h-3.5" />
                <span>{{ site.batch_mode_enabled ? '批量' : '动态' }}</span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click.stop
              @click="router.push(`/console/sites/${site.id}`)"
              class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
              title="管理"
            >
              <Settings class="w-4 h-4" />
            </button>
            <button
              @click.stop
              @click="deleteMutation.mutate(site.id)"
              class="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
              title="删除"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreate" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm p-4">
      <div class="bg-white rounded-xl border border-slate-200 shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-5 py-4 border-b border-slate-100">
          <h2 class="text-lg font-semibold text-slate-900">创建新站点</h2>
        </div>

        <form @submit="handleCreate" class="p-5 space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">站点名称</label>
            <input
              v-model="newSiteName"
              type="text"
              placeholder="例如: 我的博客"
              class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">
              域名白名单
              <span class="text-xs text-slate-400 font-normal">（每行一个，选填）</span>
            </label>
            <textarea
              v-model="newDomains"
              rows="3"
              placeholder="example.com&#10;*.example.com"
              class="w-full px-3 py-2 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all resize-none"
            />
          </div>

          <div class="flex items-center justify-between py-2 border-b border-slate-100">
            <div>
              <div class="text-sm font-medium text-slate-900">批量模式</div>
              <div class="text-xs text-slate-500">固定难度，根据风控调整 PoW 数量</div>
            </div>
            <button
              type="button"
              role="switch"
              :aria-checked="batchModeEnabled"
              @click="batchModeEnabled = !batchModeEnabled"
              :class="cn(
                'relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer',
                batchModeEnabled ? 'bg-emerald-600' : 'bg-gray-200'
              )"
            >
              <span
                :class="cn(
                  'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-lg ring-0 transition-transform duration-200 ease-in-out',
                  batchModeEnabled ? 'translate-x-5' : 'translate-x-0'
                )"
              />
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">
              {{ batchModeEnabled ? '批量难度' : '基础难度' }}
            </label>
            <input
              v-model.number="batchDifficulty"
              type="number"
              min="1"
              max="8"
              class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-700 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
            />
          </div>

          <div class="flex justify-end gap-3 pt-2">
            <button
              type="button"
              @click="showCreate = false"
              class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="createMutation.isPending.value || !newSiteName.trim()"
              class="h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ createMutation.isPending.value ? '创建中...' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Created Site Modal -->
    <div v-if="createdSite" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm p-4">
      <div class="bg-white rounded-xl border border-slate-200 shadow-xl w-full max-w-lg overflow-hidden">
        <div class="px-5 py-3 border-b border-violet-100 bg-gradient-to-r from-violet-50 to-purple-50">
          <div class="flex items-center gap-2 text-violet-600">
            <Check class="w-5 h-5" />
            <h2 class="font-semibold">站点创建成功</h2>
          </div>
        </div>
        <div class="p-5 space-y-4">
          <div class="flex items-center gap-2 p-2.5 rounded-lg bg-blue-50 border border-blue-200 text-blue-700 text-sm">
            <AlertCircle class="w-4 h-4 flex-shrink-0" />
            <span>站点创建成功！Secret 可在站点详情页查看</span>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-600 mb-1.5">
              Sitekey（前端使用）
            </label>
            <div class="flex items-center gap-2">
              <code class="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono">
                {{ createdSite.sitekey }}
              </code>
              <button
                @click="handleCopy(createdSite.sitekey, 'new-sitekey')"
                class="h-9 w-9 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
              >
                <Check v-if="copiedField === 'new-sitekey'" class="h-4 w-4 text-emerald-500" />
                <Copy v-else class="h-4 w-4" />
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-600 mb-1.5">
              Secret（后端使用，请勿泄露）
            </label>
            <div class="flex items-center gap-2">
              <code class="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                {{ createdSite.secret }}
              </code>
              <button
                @click="handleCopy(createdSite.secret, 'new-secret')"
                class="h-9 w-9 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
              >
                <Check v-if="copiedField === 'new-secret'" class="h-4 w-4 text-emerald-500" />
                <Copy v-else class="h-4 w-4" />
              </button>
            </div>
          </div>

          <div class="flex justify-end pt-2">
            <button
              @click="createdSite = null"
              class="h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
            >
              我已保存
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>