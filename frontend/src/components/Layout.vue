<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute, RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '@/lib/auth'
import { cn } from '@/lib/utils'
import { Globe, User, Settings, Menu, X, BarChart3, ExternalLink, LogOut, LayoutDashboard, Users, ChevronDown, Shield } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const sidebarOpen = ref(false)

const sidebarItems = computed(() => {
  const items: { path: string; label: string; icon: any }[] = [
    { path: '/console/sites', label: '站点管理', icon: Globe },
  ]

  if (auth.user?.role === 'admin') {
    items.push(
      { path: '/console/admin/dashboard', label: '管理后台', icon: LayoutDashboard },
      { path: '/console/admin/users', label: '用户管理', icon: Users },
      { path: '/console/admin/sites', label: '所有站点', icon: Globe },
    )
  }

  items.push(
    { path: '/console/profile', label: '个人中心', icon: User },
    { path: '/console/privacy', label: '隐私设置', icon: Shield },
  )

  return items
})

const handleLogout = async () => {
  await auth.logout()
  window.location.href = '/'
}

const isActive = (path: string) => {
  if (path === '/console/admin/sites') {
    return route.path.startsWith('/console/admin/sites')
  }
  if (path === '/console/admin/dashboard' || path === '/console/admin/users') {
    return route.path === path
  }
  if (path === '/console/sites') {
    return route.path === '/console/sites' || (route.path.startsWith('/console/sites') && !route.path.startsWith('/console/admin/sites'))
  }
  return route.path === path
}
</script>

<template>
  <!-- Unauthenticated state -->
  <div v-if="!auth.isAuthenticated && !auth.isLoading" class="min-h-screen flex items-center justify-center pt-14">
    <div class="text-center">
      <p class="text-slate-500 text-sm mb-4">请先登录</p>
      <RouterLink to="/login" replace class="inline-flex items-center h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors">
        前往登录
      </RouterLink>
    </div>
  </div>

  <!-- Loading state -->
  <div v-else-if="auth.isLoading" class="min-h-screen flex items-center justify-center pt-14">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
  </div>

  <!-- Main layout -->
  <div v-else class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30 pt-14">
    <!-- Mobile sidebar backdrop -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-40 bg-black/30 backdrop-blur-sm md:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Sidebar toggle button (mobile) -->
    <button
      class="fixed bottom-6 right-6 z-50 md:hidden w-12 h-12 bg-blue-600 text-white rounded-full shadow-lg flex items-center justify-center hover:bg-blue-700 transition-colors"
      @click="sidebarOpen = !sidebarOpen"
    >
      <Menu v-if="!sidebarOpen" class="w-5 h-5" />
      <X v-else class="w-5 h-5" />
    </button>

    <!-- Centered layout container -->
    <div class="max-w-7xl mx-auto flex min-h-[calc(100vh-3.5rem)]">
      <!-- Sidebar -->
      <aside
        :class="cn(
          'w-56 bg-white/95 backdrop-blur-md border-r border-slate-200/60 flex-shrink-0',
          'transition-transform duration-200',
          sidebarOpen
            ? 'fixed top-14 left-0 z-40 h-[calc(100vh-3.5rem)] translate-x-0'
            : 'hidden',
          'md:sticky md:top-14 md:block md:h-[calc(100vh-3.5rem)] md:translate-x-0'
        )"
      >
        <div class="flex flex-col h-full">
          <!-- Navigation items -->
          <nav class="flex-1 overflow-y-auto p-3 pt-5 space-y-1">
            <RouterLink
              v-for="item in sidebarItems"
              :key="item.path"
              :to="item.path"
              @click="sidebarOpen = false"
              :class="cn(
                'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                isActive(item.path)
                  ? 'bg-[#0066ff]/10 text-[#0066ff]'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
              )"
            >
              <component :is="item.icon" class="w-4.5 h-4.5 flex-shrink-0" />
              {{ item.label }}
            </RouterLink>
          </nav>

          <!-- Bottom section -->
          <div class="p-3 border-t border-slate-100">
            <!-- Demo link -->
            <RouterLink
              to="/demo"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-slate-500 hover:text-slate-700 hover:bg-slate-100 transition-colors"
            >
              <ExternalLink class="w-4.5 h-4.5 flex-shrink-0" />
              查看演示
            </RouterLink>

            <!-- User info & logout -->
            <div class="mt-1 px-3 py-2.5 flex items-center gap-2 border-t border-slate-100 pt-3">
              <div
                class="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-xs font-medium text-white flex-shrink-0"
              >
                {{ auth.user?.username?.charAt(0) || 'U' }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-xs font-medium text-slate-900 truncate">{{ auth.user?.username || '用户' }}</p>
              </div>
              <button
                @click="handleLogout"
                class="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-md transition-colors"
                title="退出登录"
              >
                <LogOut class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </aside>

      <!-- Main content -->
      <main class="flex-1 min-h-[calc(100vh-3.5rem)] min-w-0">
        <div class="p-4 sm:p-6 lg:p-8">
          <RouterView />
        </div>
      </main>
    </div>
  </div>
</template>