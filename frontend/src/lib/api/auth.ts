import { request, API_BASE, envelope } from './client'
import type { User } from './types'

export const authApi = {
  getOIDCUrl: () => `${API_BASE}/auth/oidc/authorize`,

  //login: (emailOrUsername: string, password: string) =>       // 密码登录暂未开放
  //  request<{ token: string; user: User }>('/auth/login', {
  //    method: 'POST',
  //    body: JSON.stringify(envelope({ email_or_username: emailOrUsername, password })),
  //  }),

  //register: (email: string, password: string) =>        // 邮箱注册暂未开放
  //  request<{ token: string; user: User }>('/auth/register', {
  //    method: 'POST',
  //    body: JSON.stringify(envelope({ email, password })),
  //  }),

  logout: () =>
    request<{ message: string }>('/auth/logout', { method: 'POST' }),

  me: () => request<User>('/auth/me'),
}