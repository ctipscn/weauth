<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { adminUsersApi, adminSitesApi } from '@/lib/api'
import { Users, Globe, Shield, Activity, Loader2, ArrowRight, CheckCircle2, XCircle } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const { data: users, isLoading: usersLoading } = useQuery({
  queryKey: ['admin-users'],
  queryFn: adminUsersApi.list,
  retry: false,
})

const { data: sites, isLoading: sitesLoading } = useQuery({
  queryKey: ['admin-sites'],
  queryFn: adminSitesApi.list,
  retry: false,
})

const isLoading = computed(() => usersLoading.value || sitesLoading.value)
const totalUsers = computed(() => users.value?.length ?? 0)
const totalSites = computed(() => sites.value?.length ?? 0)
const adminUsers = computed(() => users.value?.filter(u => u.role === 'admin').length ?? 0)
const enabledSites = computed(() => sites.value?.filter(s => s.enabled).length ?? 0)
const totalDomains = computed(() => sites.value?.reduce((sum, s) => sum + (s.domains?.length ?? 0), 0) ?? 0)

const quickLinks = [
  { to: '/console/admin/users', label: '用户管理', desc: '查看和管理所有用户', icon: Users, color: 'blue' },
  { to: '/console/admin/sites', label: '站点管理', desc: '查看和管理所有站点', icon: Globe, color: 'violet' },
]

const statCards = computed(() => [
  { label: '总用户数', value: totalUsers.value, icon: Users, sub: `${adminUsers.value} 个管理员`, color: 'blue' },
  { label: '总站点数', value: totalSites.value, icon: Globe, sub: `${enabledSites.value} 个启用`, color: 'emerald' },
  { label: '总域名数', value: totalDomains.value, icon: Shield, sub: '所有站点域名白名单', color: 'amber' },
])
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-xl font-semibold text-slate-900">管理后台</h1>
      <p class="text-slate-500 text-sm mt-0.5">系统概览与管理操作</p>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
    </div>

    <template v-else>
      <!-- Stats -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div
          v-for="card in statCards"
          :key="card.label"
          class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-5 hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between mb-3">
            <span class="text-sm font-medium text-slate-500">{{ card.label }}</span>
            <div
              :class="cn(
                'w-9 h-9 rounded-lg flex items-center justify-center',
                card.color === 'blue' && 'bg-blue-50 text-blue-600',
                card.color === 'emerald' && 'bg-emerald-50 text-emerald-600',
                card.color === 'amber' && 'bg-amber-50 text-amber-600',
              )"
            >
              <component :is="card.icon" class="w-5 h-5" />
            </div>
          </div>
          <div class="text-2xl font-bold text-slate-900">{{ card.value }}</div>
          <div class="text-xs text-slate-400 mt-1">{{ card.sub }}</div>
        </div>
      </div>

      <!-- Quick Links -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <RouterLink
          v-for="link in quickLinks"
          :key="link.to"
          :to="link.to"
          class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-5 hover:shadow-md hover:border-blue-200 transition-all group"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center text-blue-600">
                <component :is="link.icon" class="w-5 h-5" />
              </div>
              <div>
                <p class="font-medium text-slate-900">{{ link.label }}</p>
                <p class="text-sm text-slate-500">{{ link.desc }}</p>
              </div>
            </div>
            <ArrowRight class="w-5 h-5 text-slate-300 group-hover:text-blue-500 transition-colors" />
          </div>
        </RouterLink>
      </div>

      <!-- Recent Users Preview -->
      <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 flex items-center justify-between">
          <h2 class="font-medium text-slate-900 text-sm">最近注册用户</h2>
          <RouterLink
            to="/console/admin/users"
            class="text-xs text-blue-600 hover:text-blue-700 font-medium"
          >
            查看全部 →
          </RouterLink>
        </div>
        <div class="divide-y divide-slate-100">
          <div
            v-for="user in users?.slice(0, 5)"
            :key="user.id"
            class="px-5 py-3 flex items-center justify-between hover:bg-slate-50 transition-colors"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div
                class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-xs font-medium text-white flex-shrink-0"
              >
                {{ user.username.charAt(0) }}
              </div>
              <div class="min-w-0">
                <p class="text-sm font-medium text-slate-900 truncate">{{ user.username }}</p>
                <p class="text-xs text-slate-500 truncate">{{ user.email }}</p>
              </div>
            </div>
            <div class="flex items-center gap-3 flex-shrink-0">
              <span
                :class="cn(
                  'px-2 py-0.5 rounded text-xs font-medium',
                  user.role === 'admin' ? 'bg-purple-50 text-purple-600' : 'bg-slate-100 text-slate-600'
                )"
              >
                {{ user.role === 'admin' ? '管理员' : '用户' }}
              </span>
              <span :class="cn(
                'flex items-center gap-1 text-xs',
                user.enabled ? 'text-emerald-600' : 'text-red-500'
              )">
                <CheckCircle2 v-if="user.enabled" class="w-3.5 h-3.5" />
                <XCircle v-else class="w-3.5 h-3.5" />
                {{ user.enabled ? '启用' : '禁用' }}
              </span>
            </div>
          </div>
          <div v-if="!users?.length" class="px-5 py-8 text-center text-sm text-slate-400">
            暂无用户数据
          </div>
        </div>
      </div>

      <!-- Recent Sites Preview -->
      <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 flex items-center justify-between">
          <h2 class="font-medium text-slate-900 text-sm">最近创建站点</h2>
          <RouterLink
            to="/console/admin/sites"
            class="text-xs text-blue-600 hover:text-blue-700 font-medium"
          >
            查看全部 →
          </RouterLink>
        </div>
        <div class="divide-y divide-slate-100">
          <div
            v-for="site in sites?.slice(0, 5)"
            :key="site.id"
            class="px-5 py-3 flex items-center justify-between hover:bg-slate-50 transition-colors"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-8 h-8 rounded-lg bg-violet-50 flex items-center justify-center flex-shrink-0">
                <Globe class="w-4 h-4 text-violet-600" />
              </div>
              <div class="min-w-0">
                <p class="text-sm font-medium text-slate-900 truncate">{{ site.name }}</p>
                <p class="text-xs text-slate-500">Sitekey: {{ site.sitekey }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2 flex-shrink-0">
              <span class="text-xs text-slate-500">用户 #{{ site.user_id }}</span>
              <span
                :class="cn(
                  'px-2 py-0.5 rounded text-xs font-medium',
                  site.enabled ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'
                )"
              >
                {{ site.enabled ? '启用' : '禁用' }}
              </span>
            </div>
          </div>
          <div v-if="!sites?.length" class="px-5 py-8 text-center text-sm text-slate-400">
            暂无站点数据
          </div>
        </div>
      </div>
    </template>
  </div>
</template>