<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { adminSitesApi } from '@/lib/api'
import type { AdminSite } from '@/lib/api'
import { Globe, Search, ExternalLink, Trash2, Loader2, AlertTriangle, Users } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const router = useRouter()
const queryClient = useQueryClient()

const searchQuery = ref('')
const deletingSite = ref<AdminSite | null>(null)

const { data: sites, isLoading } = useQuery({
  queryKey: ['admin-sites'],
  queryFn: adminSitesApi.list,
  retry: false,
})

const filteredSites = computed(() => {
  if (!sites.value) return []
  if (!searchQuery.value.trim()) return sites.value
  const q = searchQuery.value.toLowerCase()
  return sites.value.filter(
    (s: AdminSite) =>
      s.name.toLowerCase().includes(q) ||
      s.sitekey.toLowerCase().includes(q) ||
      String(s.user_id).includes(q)
  )
})

const deleteMutation = useMutation({
  mutationFn: () => adminSitesApi.delete(deletingSite.value!.id),
  onSuccess: () => {
    deletingSite.value = null
    queryClient.invalidateQueries({ queryKey: ['admin-sites'] })
  },
})

const confirmDelete = (site: AdminSite) => {
  deletingSite.value = site
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">站点管理</h1>
        <p class="text-slate-500 text-sm mt-0.5">查看和管理系统中所有验证站点</p>
      </div>
      <div class="text-sm text-slate-400">
        共 {{ sites?.length ?? 0 }} 个站点
      </div>
    </div>

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索站点名称、Sitekey 或用户 ID..."
        class="w-full pl-10 pr-4 h-9 bg-white/80 border border-slate-200 rounded-lg text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
      />
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
    </div>

    <!-- Empty -->
    <div v-else-if="!isLoading && sites?.length === 0" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
      <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
        <Globe class="w-6 h-6 text-slate-400" />
      </div>
      <h3 class="text-base font-medium text-slate-900">暂无站点</h3>
      <p class="text-sm text-slate-500 mt-1">系统中还没有创建验证站点</p>
    </div>

    <!-- Sites Table -->
    <div v-else class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-100 bg-slate-50/50">
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">ID</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">站点</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">Sitekey</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">所有者</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">难度</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">域名数</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">状态</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">创建时间</th>
              <th class="text-right px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="site in filteredSites"
              :key="site.id"
              class="hover:bg-slate-50/50 transition-colors"
            >
              <td class="px-4 py-3">
                <code class="text-violet-600 font-mono text-sm">{{ site.id }}</code>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg bg-violet-50 flex items-center justify-center flex-shrink-0">
                    <Globe class="w-4 h-4 text-violet-600" />
                  </div>
                  <div class="min-w-0">
                    <p class="font-medium text-slate-900 truncate max-w-[180px]">{{ site.name }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <code class="text-xs text-slate-500 font-mono">{{ site.sitekey }}</code>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-1.5">
                  <Users class="w-3.5 h-3.5 text-slate-400" />
                  <span class="text-slate-600">#{{ site.user_id }}</span>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="text-slate-600">{{ site.base_difficulty }}→{{ site.max_difficulty }}</span>
              </td>
              <td class="px-4 py-3">
                <span class="text-slate-600">{{ site.domains?.length ?? 0 }}</span>
              </td>
              <td class="px-4 py-3">
                <span
                  :class="cn(
                    'px-2 py-0.5 rounded text-xs font-medium',
                    site.enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'
                  )"
                >
                  {{ site.enabled ? '启用' : '禁用' }}
                </span>
              </td>
              <td class="px-4 py-3 text-xs text-slate-500 whitespace-nowrap">
                {{ formatDate(site.created_at) }}
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="router.push('/console/admin/sites/' + site.id)"
                    class="inline-flex items-center gap-1 h-7 px-2.5 text-xs text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-md transition-colors font-medium"
                  >
                    <ExternalLink class="w-3.5 h-3.5" />
                    详情
                  </button>
                  <button
                    @click="confirmDelete(site)"
                    class="inline-flex items-center gap-1 h-7 px-2.5 text-xs text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors font-medium"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="deletingSite" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-red-100 bg-gradient-to-r from-red-50 to-orange-50">
          <h2 class="font-semibold text-red-700 flex items-center gap-2">
            <AlertTriangle class="w-5 h-5" />
            确认删除站点
          </h2>
        </div>
        <div class="p-6">
          <div class="flex items-center gap-3 p-3 rounded-lg bg-amber-50 border border-amber-100 mb-4">
            <div class="w-9 h-9 rounded-full bg-amber-100 flex items-center justify-center flex-shrink-0">
              <AlertTriangle class="w-5 h-5 text-amber-600" />
            </div>
            <div>
              <p class="text-sm font-medium text-amber-800">即将删除站点</p>
              <p class="text-sm text-amber-700">
                <strong>{{ deletingSite.name }}</strong> (ID: {{ deletingSite.id }})
                <span class="block text-xs text-amber-500 mt-0.5">所有者: 用户 #{{ deletingSite.user_id }}</span>
              </p>
            </div>
          </div>

          <p class="text-sm text-slate-600 mb-4">
            此操作不可恢复。站点及其所有配置（域名白名单、IP 策略、挑战记录等）将永久删除。
          </p>

          <div class="flex justify-end gap-3">
            <button
              @click="deletingSite = null"
              class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
            >
              取消
            </button>
            <button
              @click="deleteMutation.mutate()"
              :disabled="deleteMutation.isPending.value"
              class="h-9 px-4 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-lg transition-colors text-sm font-medium disabled:opacity-50 flex items-center gap-1.5 shadow-sm"
            >
              <Loader2 v-if="deleteMutation.isPending.value" class="h-4 w-4 animate-spin" />
              <template v-else>
                <Trash2 class="h-4 w-4" />
                确认删除
              </template>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>