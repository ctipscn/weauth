<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { adminUsersApi } from '@/lib/api'
import type { AdminUser } from '@/lib/api'
import { Users, Search, Trash2, AlertTriangle, Loader2, CheckCircle2, XCircle, Shield, Calendar, ArrowUpDown } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const queryClient = useQueryClient()

const searchQuery = ref('')
const deletingUser = ref<AdminUser | null>(null)
const cascadeSites = ref(false)

const { data: users, isLoading } = useQuery({
  queryKey: ['admin-users'],
  queryFn: adminUsersApi.list,
  retry: false,
})

const filteredUsers = computed(() => {
  if (!users.value) return []
  if (!searchQuery.value.trim()) return users.value
  const q = searchQuery.value.toLowerCase()
  return users.value.filter(
    (u: AdminUser) =>
      u.username.toLowerCase().includes(q) ||
      u.email.toLowerCase().includes(q) ||
      String(u.id).includes(q)
  )
})

const deleteMutation = useMutation({
  mutationFn: () => adminUsersApi.delete(deletingUser.value!.id, cascadeSites.value || undefined),
  onSuccess: () => {
    deletingUser.value = null
    cascadeSites.value = false
    queryClient.invalidateQueries({ queryKey: ['admin-users'] })
  },
})

const confirmDelete = (user: AdminUser) => {
  deletingUser.value = user
  cascadeSites.value = false
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
        <h1 class="text-xl font-semibold text-slate-900">用户管理</h1>
        <p class="text-slate-500 text-sm mt-0.5">查看和管理系统中所有用户</p>
      </div>
      <div class="text-sm text-slate-400">
        共 {{ users?.length ?? 0 }} 个用户
      </div>
    </div>

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索用户名、邮箱或用户 ID..."
        class="w-full pl-10 pr-4 h-9 bg-white/80 border border-slate-200 rounded-lg text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
      />
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
    </div>

    <!-- Empty -->
    <div v-else-if="!isLoading && users?.length === 0" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
      <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
        <Users class="w-6 h-6 text-slate-400" />
      </div>
      <h3 class="text-base font-medium text-slate-900">暂无用户</h3>
      <p class="text-sm text-slate-500 mt-1">系统中还没有注册用户</p>
    </div>

    <!-- Users Table -->
    <div v-else class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-100 bg-slate-50/50">
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">ID</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">用户</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">邮箱</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">角色</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">状态</th>
              <th class="text-left px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">注册时间</th>
              <th class="text-right px-4 py-3 font-medium text-slate-500 text-xs uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="user in filteredUsers"
              :key="user.id"
              class="hover:bg-slate-50/50 transition-colors"
            >
              <td class="px-4 py-3">
                <code class="text-blue-600 font-bold font-mono text-sm">{{ user.id }}</code>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div
                    class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-xs font-medium text-white flex-shrink-0"
                  >
                    {{ user.username.charAt(0) }}
                  </div>
                  <span class="font-medium text-slate-900">{{ user.username }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-slate-600">{{ user.email }}</td>
              <td class="px-4 py-3">
                <span
                  :class="cn(
                    'inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium',
                    user.role === 'admin'
                      ? 'bg-purple-50 text-purple-600'
                      : 'bg-slate-100 text-slate-600'
                  )"
                >
                  <Shield v-if="user.role === 'admin'" class="w-3 h-3" />
                  {{ user.role === 'admin' ? '管理员' : '用户' }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span
                  :class="cn(
                    'inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium',
                    user.enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600'
                  )"
                >
                  <CheckCircle2 v-if="user.enabled" class="w-3 h-3" />
                  <XCircle v-else class="w-3 h-3" />
                  {{ user.enabled ? '启用' : '禁用' }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-500 text-xs">
                <div class="flex items-center gap-1">
                  <Calendar class="w-3 h-3" />
                  {{ formatDate(user.created_at) }}
                </div>
              </td>
              <td class="px-4 py-3 text-right">
                <button
                  @click="confirmDelete(user)"
                  class="inline-flex items-center gap-1 h-7 px-2.5 text-xs text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors font-medium"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="deletingUser" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-red-100 bg-gradient-to-r from-red-50 to-orange-50">
          <h2 class="font-semibold text-red-700 flex items-center gap-2">
            <AlertTriangle class="w-5 h-5" />
            确认删除用户
          </h2>
        </div>
        <div class="p-6 space-y-4">
          <div class="flex items-center gap-3 p-3 rounded-lg bg-amber-50 border border-amber-100">
            <div class="w-9 h-9 rounded-full bg-amber-100 flex items-center justify-center flex-shrink-0">
              <AlertTriangle class="w-5 h-5 text-amber-600" />
            </div>
            <div>
              <p class="text-sm font-medium text-amber-800">即将删除用户</p>
              <p class="text-sm text-amber-700">
                <strong>{{ deletingUser.username }}</strong> (ID: {{ deletingUser.id }})
              </p>
            </div>
          </div>

          <label class="flex items-start gap-3 p-3 rounded-lg border border-slate-200 hover:border-orange-200 cursor-pointer transition-colors">
            <input
              v-model="cascadeSites"
              type="checkbox"
              class="mt-0.5 h-4 w-4 rounded border-slate-300 text-orange-500 focus:ring-orange-500"
            />
            <div>
              <p class="text-sm font-medium text-slate-900">级联删除站点</p>
              <p class="text-xs text-slate-500 mt-0.5">
                同时删除该用户下的所有站点及关联数据（域名、IP 策略、挑战记录等）。
                不勾选则仅软删除用户，站点数据保留。
              </p>
            </div>
          </label>

          <div class="flex items-center gap-2 p-2.5 rounded-lg bg-red-50 border border-red-100 text-red-600 text-xs">
            <AlertTriangle class="w-4 h-4 flex-shrink-0" />
            <span>此操作不可恢复，请谨慎确认。</span>
          </div>

          <div class="flex justify-end gap-3 pt-2">
            <button
              @click="deletingUser = null; cascadeSites = false"
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