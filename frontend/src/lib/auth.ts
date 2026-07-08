import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi, User } from './api'

interface AuthState {
  token: string | null
  user: User | null
  isLoading: boolean
  setToken: (token: string | null) => void
  setUser: (user: User | null) => void
  login: (token: string, user: User) => void
  logout: () => Promise<void>
  checkAuth: () => Promise<boolean>
}

export const useAuth = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      isLoading: true,

      setToken: (token) => {
        set({ token })
        if (token) {
          localStorage.setItem('token', token)
        } else {
          localStorage.removeItem('token')
        }
      },

      setUser: (user) => set({ user }),

      login: (token, user) => {
        set({ token, user, isLoading: false })
        localStorage.setItem('token', token)
      },

      logout: async () => {
        try {
          await authApi.logout()
        } catch {
          // Ignore errors
        }
        set({ token: null, user: null })
        localStorage.removeItem('token')
      },

      checkAuth: async () => {
        const { token } = get()
        if (!token) {
          set({ isLoading: false })
          return false
        }

        try {
          const user = await authApi.me()
          set({ user, isLoading: false })
          return true
        } catch {
          set({ token: null, user: null, isLoading: false })
          localStorage.removeItem('token')
          return false
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token }),
    }
  )
)
