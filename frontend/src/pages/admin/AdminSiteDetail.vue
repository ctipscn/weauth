<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { adminSitesApi, adminUsersApi } from '@/lib/api'
import type { AdminUser } from '@/lib/api'
import { ArrowLeft, Globe, Copy, Check, Trash2, Users, Loader2, AlertTriangle, Calendar } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const siteId = computed(() => Number(route.params.id))

const deletingDomain = ref<{ id: number; domain: string } | null>(null)

const { data: site, isLoading } = useQuery({
  queryKey: ['admin-site', siteId],
  queryFn: () => adminSitesApi.get(siteId.value),
  retry: false,
})

const { data: users } = useQuery({
  queryKey: ['admin-users'],
  queryFn: adminUsersApi.list,
  retry: false,
})

const owner = computed<AdminUser | undefined>(() => {
  if (!site.value || !users.value) return undefined
  return users.value.find((u: AdminUser) => u.id === site.value!.user_id)
})

const deleteDomainMutation = useMutation({
  mutationFn: () => adminSitesApi.deleteDomain(siteId.value, deletingDomain.value!.id),
  onSuccess: () => {
    deletingDomain.value = null
    queryClient.invalidateQueries({ queryKey: ['admin-site', siteId] })
  },
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
    </div>

    <div v-else-if="!site" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
      <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
        <Globe class="w-6 h-6 text-slate-400" />
      </div>
      <h3 class="text-base font-medium text-slate-900">站点不存在</h3>
      <p class="text-sm text-slate-500 mt-1 mb-4">找不到该站点或已被删除</p>
      <button
        @click="router.push('/console/admin/sites')"
        class="inline-flex items-center h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
      >
        <ArrowLeft class="w-4 h-4 mr-1.5" />
        返回站点列表
      </button>
    </div>

    <template v-else>
      <!-- Back -->
      <button
        @click="router.push('/console/admin/sites')"
        class="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-700 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        返回站点列表
      </button>

      <!-- Header -->
      <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-100 bg-gradient-to-r from-violet-50 to-blue-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-white border border-violet-200 flex items-center justify-center">
                <Globe class="w-5 h-5 text-violet-600" />
              </div>
              <div>
                <h1 class="text-lg font-semibold text-slate-900">{{ site.name }}</h1>
                <div class="flex items-center gap-2 text-xs text-slate-500 mt-0.5">
                  <span>ID: {{ site.id }}</span>
                  <span>·</span>
                  <span>Sitekey: <code class="font-mono text-violet-600">{{ site.sitekey }}</code></span>
                </div>
              </div>
            </div>
            <span
              :class="cn(
                'px-2.5 py-1 rounded-lg text-xs font-medium',
                site.enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'
              )"
            >
              {{ site.enabled ? '已启用' : '已禁用' }}
            </span>
          </div>
        </div>

        <!-- Owner Info -->
        <div class="px-6 py-3 border-b border-slate-100 bg-slate-50/50">
          <div class="flex items-center gap-2 text-sm">
            <Users class="w-4 h-4 text-slate-400" />
            <span class="text-slate-600">所有者：</span>
            <template v-if="owner">
              <div class="flex items-center gap-1.5">
                <span class="font-medium text-slate-900">{{ owner.username }}</span>
                <span class="text-slate-400">({{ owner.email }})</span>
                <span
                  :class="cn(
                    'px-1.5 py-0.5 rounded text-xs font-medium',
                    owner.role === 'admin' ? 'bg-purple-50 text-purple-600' : 'bg-slate-100 text-slate-600'
                  )"
                >
                  {{ owner.role === 'admin' ? '管理员' : '用户' }}
                </span>
                <span
                  :class="cn(
                    'inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-xs font-medium',
                    owner.enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600'
                  )"
                >
                  {{ owner.enabled ? '启用' : '禁用' }}
                </span>
              </div>
            </template>
            <template v-else>
              <span class="text-slate-400">用户 #{{ site.user_id }}</span>
            </template>
          </div>
        </div>

        <!-- Info Grid -->
        <div class="divide-y divide-slate-100">
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">创建时间</span>
            <span class="text-sm text-slate-900">{{ formatDate(site.created_at) }}</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">更新时间</span>
            <span class="text-sm text-slate-900">{{ formatDate(site.updated_at) }}</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">基础难度</span>
            <span class="text-sm text-slate-900 font-medium">{{ site.base_difficulty }}</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">最大难度</span>
            <span class="text-sm text-slate-900 font-medium">{{ site.max_difficulty }}</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">挑战超时</span>
            <span class="text-sm text-slate-900">{{ site.challenge_timeout }}s</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">Token 超时</span>
            <span class="text-sm text-slate-900">{{ site.token_timeout }}s</span>
          </div>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">批量模式</span>
            <span
              :class="cn(
                'px-2 py-0.5 rounded text-xs font-medium',
                site.batch_mode_enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'
              )"
            >
              {{ site.batch_mode_enabled ? '已启用' : '已禁用' }}
            </span>
          </div>
          <template v-if="site.batch_mode_enabled">
            <div class="px-6 py-3 flex items-center justify-between">
              <span class="text-sm text-slate-600">批量难度</span>
              <span class="text-sm text-slate-900">{{ site.batch_difficulty }}</span>
            </div>
            <div class="px-6 py-3 flex items-center justify-between">
              <span class="text-sm text-slate-600">批次数范围</span>
              <span class="text-sm text-slate-900">{{ site.min_batch_count }} ~ {{ site.max_batch_count }}</span>
            </div>
          </template>
          <div class="px-6 py-3 flex items-center justify-between">
            <span class="text-sm text-slate-600">Secret</span>
            <code class="px-2 py-0.5 bg-slate-100 rounded text-xs font-mono text-slate-500">{{ site.secret === '[redacted]' ? '[已隐藏]' : site.secret }}</code>
          </div>
        </div>
      </div>

      <!-- Domains -->
      <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 flex items-center justify-between">
          <h2 class="font-medium text-slate-900 text-sm flex items-center gap-2">
            <Globe class="w-4 h-4 text-blue-600" />
            域名白名单
          </h2>
          <span class="text-xs text-slate-400">{{ site.domains?.length ?? 0 }} 个域名</span>
        </div>

        <div v-if="!site.domains?.length" class="px-5 py-6 text-center text-sm text-slate-400">
          暂无域名白名单
        </div>

        <div v-else class="divide-y divide-slate-100">
          <div
            v-for="domain in site.domains"
            :key="domain.id"
            class="px-5 py-3 flex items-center justify-between hover:bg-slate-50/50 transition-colors"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center flex-shrink-0">
                <Globe class="w-4 h-4 text-blue-600" />
              </div>
              <div class="min-w-0">
                <p class="text-sm font-medium text-slate-900 font-mono">{{ domain.domain }}</p>
                <p class="text-xs text-slate-500">ID: {{ domain.id }} · {{ formatDate(domain.created_at) }}</p>
              </div>
            </div>
            <button
              @click="deletingDomain = { id: domain.id, domain: domain.domain }"
              class="inline-flex items-center gap-1 h-7 px-2.5 text-xs text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors font-medium flex-shrink-0"
            >
              <Trash2 class="w-3.5 h-3.5" />
              删除
            </button>
          </div>
        </div>
      </div>

      <!-- Domain Delete Confirm Modal -->
      <div v-if="deletingDomain" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
        <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
          <div class="px-6 py-4 border-b border-red-100 bg-gradient-to-r from-red-50 to-orange-50">
            <h2 class="font-semibold text-red-700 flex items-center gap-2">
              <AlertTriangle class="w-5 h-5" />
              确认删除域名
            </h2>
          </div>
          <div class="p-6">
            <div class="flex items-center gap-3 p-3 rounded-lg bg-amber-50 border border-amber-100 mb-4">
              <div class="w-9 h-9 rounded-full bg-amber-100 flex items-center justify-center flex-shrink-0">
                <Globe class="w-5 h-5 text-amber-600" />
              </div>
              <div>
                <p class="text-sm font-medium text-amber-800">即将从白名单中移除</p>
                <p class="text-sm font-mono text-amber-700">{{ deletingDomain.domain }}</p>
              </div>
            </div>
            <p class="text-sm text-slate-600 mb-4">删除后该域名将无法用于此站点的验证请求。</p>
            <div class="flex justify-end gap-3">
              <button
                @click="deletingDomain = null"
                class="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
              >
                取消
              </button>
              <button
                @click="deleteDomainMutation.mutate()"
                :disabled="deleteDomainMutation.isPending.value"
                class="h-9 px-4 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-lg transition-colors text-sm font-medium disabled:opacity-50 flex items-center gap-1.5 shadow-sm"
              >
                <Loader2 v-if="deleteDomainMutation.isPending.value" class="h-4 w-4 animate-spin" />
                <template v-else>
                  <Trash2 class="h-4 w-4" />
                  确认删除
                </template>
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>