<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { userApi } from '@/lib/api'
import { useAuthStore } from '@/lib/auth'
import { cn } from '@/lib/utils'
import {
  User,
  Calendar,
  Loader2,
  Pencil,
  Check,
  AlertCircle,
} from 'lucide-vue-next'

const auth = useAuthStore()
const queryClient = useQueryClient()

const { data: profile, isLoading } = useQuery({
  queryKey: ['profile'],
  queryFn: userApi.getProfile,
})

// Edit form state
const editing = ref(false)
const editUsername = ref('')
const editAvatar = ref('')
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)

const hasChanges = computed(() => {
  if (!profile.value) return false
  return editUsername.value !== profile.value.username ||
         editAvatar.value !== (profile.value.avatar || '')
})

const startEditing = () => {
  if (!profile.value) return
  editUsername.value = profile.value.username
  editAvatar.value = profile.value.avatar || ''
  editing.value = true
  message.value = null
}

const cancelEditing = () => {
  editing.value = false
  message.value = null
}

const updateProfileMutation = useMutation({
  mutationFn: (data: { username?: string; avatar?: string }) =>
    userApi.updateProfile(data),
  onSuccess: (updatedProfile) => {
    message.value = { type: 'success', text: '资料更新成功' }
    editing.value = false
    // Refresh profile data
    queryClient.setQueryData(['profile'], updatedProfile)
    queryClient.invalidateQueries({ queryKey: ['profile'] })
    // Sync auth store
    if (auth.user) {
      auth.user = {
        ...auth.user,
        username: updatedProfile.username,
        avatar: updatedProfile.avatar,
      }
    }
  },
  onError: (error: any) => {
    message.value = { type: 'error', text: error instanceof Error ? error.message : '更新失败' }
  },
})

const handleSave = () => {
  message.value = null

  const username = editUsername.value.trim()
  if (!username) {
    message.value = { type: 'error', text: '用户名不能为空' }
    return
  }
  if (username.length < 2) {
    message.value = { type: 'error', text: '用户名至少需要 2 个字符' }
    return
  }

  const data: { username?: string; avatar?: string } = {}
  if (username !== profile.value?.username) data.username = username
  if (editAvatar.value !== (profile.value?.avatar || '')) {
    data.avatar = editAvatar.value || ''
  }

  if (Object.keys(data).length === 0) {
    editing.value = false
    return
  }

  updateProfileMutation.mutate(data)
}
</script>

<template>
  <div v-if="isLoading" class="flex items-center justify-center h-64">
    <Loader2 class="animate-spin h-6 w-6 text-blue-600" />
  </div>

  <div v-else class="max-w-2xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">个人中心</h1>
        <p class="text-slate-500 text-sm mt-0.5">管理您的账户信息</p>
      </div>
      <button
        v-if="!editing"
        type="button"
        class="flex items-center gap-1.5 h-8 px-3 bg-white border border-slate-200 hover:bg-slate-50 active:bg-slate-100 text-slate-600 rounded-lg text-sm font-medium transition-all duration-150 shadow-sm"
        @click="startEditing"
      >
        <Pencil class="w-3.5 h-3.5" />
        编辑资料
      </button>
    </div>

    <!-- Profile Card -->
    <div class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <!-- Avatar Section -->
      <div class="p-5 border-b border-slate-100 bg-slate-50">
        <div class="flex items-center gap-4">
          <img
            v-if="profile?.avatar"
            :src="profile.avatar"
            :alt="profile.username"
            class="h-14 w-14 rounded-xl object-cover"
          />
          <div
            v-else
            class="h-14 w-14 rounded-xl bg-blue-600 flex items-center justify-center text-xl font-semibold text-white shadow-sm"
          >
            {{ profile?.username?.charAt(0) || 'U' }}
          </div>
          <div class="min-w-0">
            <p class="text-sm text-slate-500 truncate">
              用户名：<span class="text-slate-900 font-medium">{{ profile?.username }}</span>
            </p>
            <p class="text-sm text-slate-500 truncate">
              邮箱：<span class="text-slate-900">{{ profile?.email }}</span>
            </p>
          </div>
        </div>
      </div>

      <!-- Info Grid -->
      <div class="divide-y divide-slate-100">
        <div class="flex items-center justify-between px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
              <User class="w-4 h-4 text-slate-500" />
            </div>
            <span class="text-sm text-slate-600">用户 ID</span>
          </div>
          <code class="px-2.5 py-1 bg-blue-600/10 text-blue-600 rounded-md text-base font-bold font-mono">
            {{ profile?.id }}
          </code>
        </div>

        <div class="flex items-center justify-between px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
              <svg class="w-4 h-4 text-slate-500" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
              </svg>
            </div>
            <span class="text-sm text-slate-600">角色</span>
          </div>
          <span
            :class="cn(
              'px-2 py-0.5 rounded text-sm font-medium',
              profile?.role === 'admin' ? 'bg-purple-50 text-purple-600' : 'bg-slate-100 text-slate-600'
            )"
          >
            {{ profile?.role === 'admin' ? '管理员' : '用户' }}
          </span>
        </div>

        <div class="flex items-center justify-between px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
              <Calendar class="w-4 h-4 text-slate-500" />
            </div>
            <span class="text-sm text-slate-600">注册时间</span>
          </div>
          <span class="text-slate-900 text-sm">
            {{ profile?.created_at && new Date(profile.created_at).toLocaleDateString('zh-CN') }}
          </span>
        </div>
      </div>
    </div>

    <!-- Edit Profile Card -->
    <div v-if="editing" class="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div class="px-5 py-3 border-b border-slate-100">
        <h2 class="font-medium text-slate-900 text-sm flex items-center gap-2">
          <Pencil class="w-4 h-4 text-blue-600" />
          编辑资料
        </h2>
      </div>

      <form @submit.prevent="handleSave" class="p-5 space-y-4">
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
          <label class="block text-sm font-medium text-slate-700 mb-1.5">用户名</label>
          <input
            v-model="editUsername"
            type="text"
            placeholder="输入用户名"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
            required
            minlength="2"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1.5">头像 URL</label>
          <input
            v-model="editAvatar"
            type="url"
            placeholder="https://example.com/avatar.png"
            class="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder:text-slate-400 text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600 transition-all"
          />
          <div v-if="editAvatar" class="mt-2 flex items-center gap-3">
            <img
              :src="editAvatar"
              alt="头像预览"
              class="h-10 w-10 rounded-lg object-cover border border-slate-200"
              @error="($event.target as HTMLImageElement).style.display = 'none'"
              @load="($event.target as HTMLImageElement).style.display = ''"
            />
            <span class="text-xs text-slate-400">预览</span>
          </div>
        </div>

        <div class="flex gap-3 pt-1">
          <button
            type="button"
            class="flex-1 h-9 bg-slate-100 hover:bg-slate-200 active:bg-slate-300 text-slate-700 rounded-lg text-sm font-medium transition-all duration-150"
            @click="cancelEditing"
            :disabled="updateProfileMutation.isPending.value"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="updateProfileMutation.isPending.value || !hasChanges"
            class="flex-1 h-9 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-lg text-sm font-medium transition-all duration-150 disabled:opacity-50 flex items-center justify-center shadow-sm"
          >
            <Loader2 v-if="updateProfileMutation.isPending.value" class="h-4 w-4 animate-spin" />
            <template v-else>
              <Check class="h-4 w-4 mr-1.5" />
              保存修改
            </template>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>