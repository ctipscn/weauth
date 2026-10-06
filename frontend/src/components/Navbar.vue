<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/lib/auth'
import { cn } from '@/lib/utils'
import { LogOut, ChevronDown, Settings } from 'lucide-vue-next'

const route = useRoute()
const auth = useAuthStore()

const userMenuOpen = ref(false)



const toggleUserMenu = () => {
  userMenuOpen.value = !userMenuOpen.value
}

const handleLogout = async () => {
  userMenuOpen.value = false
  await auth.logout()
  window.location.href = '/'
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.user-menu-container')) {
    userMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})
</script>

<style scoped>
.btn-outline {
  @apply px-4 py-2 bg-white border-2 border-blue-600 text-blue-600 text-sm font-medium rounded-lg transition-all;
}
.btn-outline:hover {
  @apply bg-blue-50 border-blue-700 text-blue-700;
}
</style>


<template>
  <nav class="fixed top-0 w-full z-50 bg-white/90 backdrop-blur-md border-b border-slate-200/60 drop-shadow-md">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <div class="flex items-center gap-8">
          <RouterLink to="/" class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-blue-50 border border-blue-200 flex items-center justify-center">
              <svg class="w-4 h-4 text-blue-600" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
              </svg>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span class="text-slate-900 font-semibold hidden sm:block tracking-tight">微验</span>
              <span class="text-slate-500 text-sm hidden sm:block">WeAuth</span>
            </div>
          </RouterLink>

          <!-- Desktop Public Nav (always visible) -->
          <nav class="hidden md:flex items-center gap-8">
            <a href="/#features" class="text-sm text-slate-600 hover:text-slate-900 transition-colors">特性</a>
            <a href="/#pricing" class="text-sm text-slate-600 hover:text-slate-900 transition-colors">定价</a>
            <RouterLink to="/demo" class="text-sm text-slate-600 hover:text-slate-900 transition-colors">演示</RouterLink>
            <RouterLink to="/docs" class="text-sm text-slate-600 hover:text-slate-900 transition-colors">文档</RouterLink>
            <a
              href="https://cloud.emoera.com/"
              target="_blank"
              rel="noopener noreferrer"
              class="text-sm text-slate-600 hover:text-slate-900 transition-colors flex items-center gap-1"
            >
              E时代云服务
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
          </nav>
        </div>

        <!-- Right Side -->
        <div class="flex items-center gap-3">
          <!-- 控制台按钮-->
          <RouterLink
            v-if="!route.path.startsWith('/console')"
            to="/console/sites"
            class="btn-outline"
          >
            控制台
          </RouterLink>

          <!-- Logged in: user menu -->
          <template v-if="auth.isAuthenticated">
            <div class="relative user-menu-container">
              <button
                @click="toggleUserMenu"
                class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-slate-100 transition-colors"
              >
                <div
                  v-if="auth.user?.avatar"
                  class="w-7 h-7 rounded-full overflow-hidden"
                >
                  <img :src="auth.user.avatar" :alt="auth.user.username" class="w-full h-full object-cover" />
                </div>
                <div
                  v-else
                  class="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-xs font-medium text-white shadow-sm"
                >
                  {{ auth.user?.username?.charAt(0) || 'U' }}
                </div>
                <span class="text-sm text-slate-700 hidden sm:block max-w-[100px] truncate">
                  {{ auth.user?.username || '用户' }}
                </span>
                <ChevronDown :class="cn('w-4 h-4 text-slate-400 transition-transform hidden sm:block', userMenuOpen && 'rotate-180')" />
              </button>

              <!-- Dropdown -->
              <div
                v-if="userMenuOpen"
                class="absolute right-0 mt-1 w-52 bg-white rounded-lg border border-slate-200 shadow-lg overflow-hidden z-50"
              >
                <div class="px-3 py-2.5 border-b border-slate-100 bg-slate-50">
                  <p class="text-sm font-medium text-slate-900 truncate">{{ auth.user?.username }}</p>
                  <p class="text-xs text-slate-500 truncate">{{ auth.user?.email }}</p>
                  <p class="text-xs text-slate-400 mt-0.5">ID: {{ auth.user?.id }}</p>
                </div>
                <div class="py-1">
                  <RouterLink
                    to="/console/profile"
                    @click="userMenuOpen = false"
                    class="flex items-center gap-2 px-3 py-2 text-sm text-slate-600 hover:text-slate-900 hover:bg-slate-50 transition-colors"
                  >
                    <Settings class="w-4 h-4" />
                    个人中心
                  </RouterLink>
                  <button
                    @click="handleLogout"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-500 hover:text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <LogOut class="w-4 h-4" />
                    退出登录
                  </button>
                </div>
              </div>
            </div>
          </template>

          <!-- Logged out: login button only -->
          <template v-else>
            <RouterLink to="/login" class="btn-outline">登录</RouterLink>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>