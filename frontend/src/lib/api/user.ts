import { request, envelope } from './client'
import type { UserProfile } from './types'

export const userApi = {
  getProfile: () =>
    request<UserProfile>('/user/profile'),

  updatePassword: (data: { old_password: string; new_password: string }) =>
    request<string>('/auth/password/update', {
      method: 'PUT',
      body: JSON.stringify(envelope(data)),
    }),

  deactivate: () =>
    request<{ message: string }>('/user/deactivate', {
      method: 'POST',
    }),

  /** 设置密码（OIDC 用户首次设置，无需旧密码） */
  setPassword: (newPassword: string) =>
    request<string>('/auth/password/set', {
      method: 'PUT',
      body: JSON.stringify(envelope({ new_password: newPassword })),
    }),

  updateProfile: (data: { username?: string; avatar?: string }) =>
    request<UserProfile>('/user/profile', {
      method: 'PUT',
      body: JSON.stringify(envelope(data)),
    }),
}