<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/lib/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

onMounted(async () => {
  const token = route.query.token as string
  if (token) {
    localStorage.setItem('token', token)
    auth.login(token, { id: 0, username: '', email: '', avatar: '', role: '' })
    await auth.checkAuth()
    router.replace('/console/sites')
  } else {
    router.replace('/login')
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
  </div>
</template>