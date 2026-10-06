<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/lib/auth'
import { authApi } from '@/lib/api'

import { AlertCircle, ArrowLeft, Info } from 'lucide-vue-next'

const router = useRouter()
const auth = useAuthStore()

// 已登录用户直接跳转到控制台
onMounted(() => {
  if (auth.isAuthenticated) {
    router.push('/console/sites')
  }
})

const handleOIDCLogin = () => {
  // 跳转到后端 OIDC 授权入口
  window.location.href = authApi.getOIDCUrl()
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30 flex items-center justify-center p-4 pt-14">
    <div class="w-full max-w-md">
      <!-- Back Link -->
      <RouterLink
        to="/"
        class="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-700 transition-colors mb-8"
      >
        <ArrowLeft class="w-4 h-4" />
        返回首页
      </RouterLink>

      <!-- Card -->
      <div class="bg-white rounded-2xl shadow-lg border border-slate-200/60 p-6 sm:p-8">
        <!-- Logo -->
        <div class="text-center mb-6">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-blue-50 to-cyan-50 border border-blue-200/50 flex items-center justify-center mx-auto mb-4">
            <svg class="w-7 h-7 text-blue-600" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
            </svg>
          </div>
          <h1 class="text-xl font-bold text-slate-900">登录 WeAuth</h1>
          <p class="text-sm text-slate-500 mt-1">当前仅支持企业号登录</p>
        </div>

        <!-- Error Message (from OIDC callback) -->
        <div
          v-if="$route.query.error"
          class="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg mb-6"
        >
          <AlertCircle class="w-4 h-4 text-red-500 flex-shrink-0" />
          <p class="text-sm text-red-600">{{ $route.query.error }}</p>
        </div>

        <!-- OIDC Login Button -->
        <button
          type="button"
          @click="handleOIDCLogin"
          class="w-full flex items-center justify-center gap-3 py-3 px-4 rounded-xl bg-[#0066ff] hover:bg-[#0052cc] text-white font-medium text-sm transition-all shadow-md shadow-blue-500/25"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v6l4 2"/>
          </svg>
          使用 E时代通行证 登录
        </button>

        <!-- Unavailable Notice -->
        <div class="mt-6 flex items-start gap-2 p-3 bg-slate-50 border border-slate-200 rounded-lg">
          <Info class="w-4 h-4 text-slate-400 flex-shrink-0 mt-0.5" />
          <p class="text-xs text-slate-500 leading-relaxed">
            密码登录及邮箱注册暂未开放，请使用上方企业号登录。
          </p>
        </div>
      </div>
    </div>
  </div>
</template>