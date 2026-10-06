import { defineStore } from 'pinia'
import { authApi } from './api'
import type { User } from './api'

/** 简易 JWT payload 解码（仅读取，不验签） */
function decodeToken(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = parts[1]!
    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch {
    return null
  }
}



export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token'),
    user: null as User | null,
    isLoading: true,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    /** 从 JWT token 中判断是否为 OIDC 登录 */
    isOidc: (state): boolean => {
      if (!state.token) return false
      const payload = decodeToken(state.token)
      if (!payload) return false
      return !!(payload as any).is_oidc_login
    },
  },
  actions: {
    setToken(token: string | null) {
      this.token = token
      if (token) {
        localStorage.setItem('token', token)
      } else {
        localStorage.removeItem('token')
      }
    },
    setUser(user: User | null) {
      this.user = user
    },
    login(token: string, user: User) {
      this.token = token
      this.user = user
      this.isLoading = false
      localStorage.setItem('token', token)
    },
    async logout() {
      try {
        await authApi.logout()
      } catch {
        // Ignore errors
      }
      this.token = null
      this.user = null
      localStorage.removeItem('token')
    },
    async checkAuth(): Promise<boolean> {
      if (!this.token) {
        this.isLoading = false
        return false
      }
      try {
        const user = await authApi.me()
        this.user = user
        this.isLoading = false
        return true
      } catch {
        this.token = null
        this.user = null
        this.isLoading = false
        localStorage.removeItem('token')
        return false
      }
    },
  },
})