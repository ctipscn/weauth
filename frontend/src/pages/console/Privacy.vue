<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { userApi } from '@/lib/api'
import { useAuthStore } from '@/lib/auth'
import { cn } from '@/lib/utils'
import {
  Loader2, AlertCircle,
  Check, Key, Shield, Trash2,
} from 'lucide-vue-next'

const router = useRouter()
const auth = useAuthStore()
const queryClient = useQueryClient()

const isOidc = computed(() => auth.isOidc)

// --- Password form state ---
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)

// --- Deactivation state ---
const confirmDeactivate = ref(false)
const deactivateMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)

const { data: profile } = useQuery({
  queryKey: ['profile'],
  queryFn: userApi.getProfile,
})

// --- 修改密码 (非 OIDC: passwd / bot) ---
const changePasswordMutation = useMutation({
  mutationFn: (data: { old_password: string; new_password: string }) =>
    userApi.updatePassword(data),
  onSuccess: () => {
    message.value = { type: 'success', text: '密码修改成功' }
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    queryClient.invalidateQueries({ queryKey: ['profile'] })
  },
  onError: (error: any) => {
    message.value = { type: 'error', text: error instanceof Error ? error.message : '修改失败' }
  },
})

const handleChangePassword = (e: Event) => {
  e.preventDefault()
  message.value = null

  if (oldPassword.value.length < 6) {
    message.value = { type: 'error', text: '当前密码至少需要 6 个字符' }
    return
  }
  if (newPassword.value.length < 6) {
    message.value = { type: 'error', text: '新密码至少需要 6 个字符' }
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    message.value = { type: 'error', text: '两次输入的新密码不一致' }
    return
  }

  changePasswordMutation.mutate({
    old_password: oldPassword.value,
    new_password: newPassword.value,
  })
}

// --- 设置密码 (OIDC) ---
const setPasswordMutation = useMutation({
  mutationFn: (password: string) => userApi.setPassword(password),
  onSuccess: () => {
    message.value = { type: 'success', text: '密码设置成功' }
    newPassword.value = ''
    confirmPassword.value = ''
    queryClient.invalidateQueries({ queryKey: ['profile'] })
  },
  onError: (error: any) => {
    message.value = { type: 'error', text: error instanceof Error ? error.message : '设置失败' }
  },
})

const handleSetPassword = (e: Event) => {
  e.preventDefault()
  message.value = null

  if (newPassword.value.length < 6) {
    message.value = { type: 'error', text: '密码至少需要 6 个字符' }
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    message.value = { type: 'error', text: '两次输入的密码不一致' }
    return
  }

  setPasswordMutation.mutate(newPassword.value)
}

// --- Deactivation ---
const deactivateMutation = useMutation({
  mutationFn: userApi.deactivate,
  onSuccess: () => {
    deactivateMessage.value = { type: 'success', text: '账号已注销' }
    setTimeout(() => {
      auth.logout()
      router.push('/')
    }, 1500)
  },
  onError: (error: any) => {
    deactivateMessage.value = { type: 'error', text: error instanceof Error ? error.message : '注销失败' }
    confirmDeactivate.value = false
  },
})

const handleDeactivate = () => {
  deactivateMessage.value = null
  deactivateMutation.mutate()
}

const cancelDeactivate = () => {
  confirmDeactivate.value = false
  deactivateMessage.value = null
}
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-xl font-semibold text-slate-900">隐私设置</h1>
      <p class="text-slate-500 text-sm mt-0.5">管理您的密码和账户安全</p>
    </div>

    <!-- ─── OIDC: 设置密码 ─── -->
    <div v-if="isOidc" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div class="px-5 py-3 border-b border-slate-100">
        <h2 class="font-medium text-slate-900 text-sm flex items-center gap-2">
          <Shield class="w-4 h-4 text-blue-600" />
          设置密码
        </h2>
        <p class="text-xs text-slate-500 mt-0.5">
          设置密码后，您也可以通过账号密码方式登录
        </p>
      </div>

      <form @submit="handleSetPassword" class="p-5 space-y-4">
        <div
          v-if="message"
          :class="cn(
            'flex items-center gap-2 p-2.5 rounded-lg text-sm',
            message.type === 'success'
              ? 'bg-emerald-50 border border-emerald-100 text-emerald-600'
              : 'bg-red-50 border border-red-100 text-red-600'
          )"
        >
          <Check v-if="message.type === 'success'" class="w-4 h-4 shrink-0" />
          <AlertCircle v-else class="w-4 h-4 shrink-0" />
          {{ message.text }}
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">密码</label>
          <input
            v-model="newPassword"
            type="password"
            placeholder="至少 6 个字符"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="6"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">确认密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="再次输入"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="6"
          />
        </div>

        <button
          type="submit"
          :disabled="setPasswordMutation.isPending.value"
          class="w-full h-9 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-lg text-sm font-medium transition-all duration-150 disabled:opacity-50 flex items-center justify-center shadow-sm"
        >
          <Loader2 v-if="setPasswordMutation.isPending.value" class="h-4 w-4 animate-spin" />
          <template v-else>
            <Check class="h-4 w-4 mr-1.5" />
            设置密码
          </template>
        </button>
      </form>
    </div>

    <!-- ─── 非 OIDC (passwd / bot): 修改密码 ─── -->
    <div v-else class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div class="px-5 py-3 border-b border-slate-100">
        <h2 class="font-medium text-slate-900 text-sm flex items-center gap-2">
          <Key class="w-4 h-4 text-blue-600" />
          修改密码
        </h2>
        <p class="text-xs text-slate-500 mt-0.5">输入当前密码并设置新密码</p>
      </div>

      <form @submit="handleChangePassword" class="p-5 space-y-4">
        <div
          v-if="message"
          :class="cn(
            'flex items-center gap-2 p-2.5 rounded-lg text-sm',
            message.type === 'success'
              ? 'bg-emerald-50 border border-emerald-100 text-emerald-600'
              : 'bg-red-50 border border-red-100 text-red-600'
          )"
        >
          <Check v-if="message.type === 'success'" class="w-4 h-4 shrink-0" />
          <AlertCircle v-else class="w-4 h-4 shrink-0" />
          {{ message.text }}
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">当前密码</label>
          <input
            v-model="oldPassword"
            type="password"
            placeholder="输入当前密码"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="6"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">新密码</label>
          <input
            v-model="newPassword"
            type="password"
            placeholder="至少 6 个字符"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="6"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">确认新密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="再次输入新密码"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="6"
          />
        </div>

        <button
          type="submit"
          :disabled="changePasswordMutation.isPending.value"
          class="w-full h-9 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-lg text-sm font-medium transition-all duration-150 disabled:opacity-50 flex items-center justify-center shadow-sm"
        >
          <Loader2 v-if="changePasswordMutation.isPending.value" class="h-4 w-4 animate-spin" />
          <template v-else>
            <Check class="h-4 w-4 mr-1.5" />
            更新密码
          </template>
        </button>
      </form>
    </div>

    <!-- ─── 注销账号 ─── -->
    <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-red-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div class="px-5 py-3 border-b border-red-100">
        <h2 class="font-medium text-red-600 text-sm flex items-center gap-2">
          <Trash2 class="w-4 h-4" />
          注销账号
        </h2>
        <p class="text-xs text-slate-500 mt-0.5">
          注销后账号将被禁用，所有关联数据将被标记为不可用，此操作不可撤销
        </p>
      </div>

      <div class="p-5 space-y-4">
        <div
          v-if="deactivateMessage"
          :class="cn(
            'flex items-center gap-2 p-2.5 rounded-lg text-sm',
            deactivateMessage.type === 'success'
              ? 'bg-emerald-50 border border-emerald-100 text-emerald-600'
              : 'bg-red-50 border border-red-100 text-red-600'
          )"
        >
          <Check v-if="deactivateMessage.type === 'success'" class="w-4 h-4 shrink-0" />
          <AlertCircle v-else class="w-4 h-4 shrink-0" />
          {{ deactivateMessage.text }}
        </div>

        <template v-if="!confirmDeactivate">
          <p class="text-sm text-slate-600">
            如需注销账号，请点击下方按钮。系统将立即禁用您的账号，您将无法再登录或使用任何 WeAuth 服务。
          </p>
          <button
            type="button"
            class="w-full h-9 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-lg text-sm font-medium transition-all duration-150 flex items-center justify-center shadow-sm"
            @click="confirmDeactivate = true"
          >
            <Trash2 class="h-4 w-4 mr-1.5" />
            注销账号
          </button>
        </template>

        <template v-else>
          <div class="bg-red-50 border border-red-200 rounded-lg p-3 text-sm text-red-700">
            <p class="font-medium mb-1">确认注销？</p>
            <p>此操作将立即禁用您的账号，关联的站点和域名将无法使用，且<strong>不可撤销</strong>。请确认您的决定。</p>
          </div>
          <div class="flex gap-3">
            <button
              type="button"
              class="flex-1 h-9 bg-slate-100 hover:bg-slate-200 active:bg-slate-300 text-slate-700 rounded-lg text-sm font-medium transition-all duration-150"
              @click="cancelDeactivate"
              :disabled="deactivateMutation.isPending.value"
            >
              取消
            </button>
            <button
              type="button"
              class="flex-1 h-9 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-lg text-sm font-medium transition-all duration-150 disabled:opacity-50 flex items-center justify-center shadow-sm"
              @click="handleDeactivate"
              :disabled="deactivateMutation.isPending.value"
            >
              <Loader2 v-if="deactivateMutation.isPending.value" class="h-4 w-4 animate-spin" />
              <template v-else>
                <AlertCircle class="h-4 w-4 mr-1.5" />
                确认注销
              </template>
            </button>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>