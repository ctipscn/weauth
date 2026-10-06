<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/lib/auth'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

onMounted(async () => {
  if (auth.token) {
    await auth.checkAuth()
  } else {
    auth.isLoading = false
  }
})
</script>

<template>
  <div v-if="auth.isLoading && route.path.startsWith('/console')" class="min-h-screen flex items-center justify-center">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
  </div>
  <template v-else>
    <Navbar />
    <router-view />
  </template>
</template>