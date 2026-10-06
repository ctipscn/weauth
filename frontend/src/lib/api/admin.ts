import { request, envelope } from './client'
import type { AdminUser, AdminSite } from './types'

export const adminUsersApi = {
  /** 获取所有用户列表（管理员） */
  list: () => request<AdminUser[]>('/admin/users'),

  /** 删除用户（可选级联删除站点） */
  delete: (id: number, cascade_sites?: boolean) =>
    request<string>(`/admin/users/${id}`, {
      method: 'DELETE',
      ...(cascade_sites !== undefined
        ? { body: JSON.stringify(envelope({ cascade_sites })) }
        : {}),
    }),
}

export const adminSitesApi = {
  /** 获取所有站点列表（管理员） */
  list: () => request<AdminSite[]>('/admin/sites'),

  /** 获取站点详情（管理员） */
  get: (id: number) => request<AdminSite>('/admin/sites/' + id),

  /** 删除站点 */
  delete: (id: number) =>
    request<string>('/admin/sites/' + id, { method: 'DELETE' }),

  /** 删除站点域名 */
  deleteDomain: (siteId: number, domainId: number) =>
    request<string>('/admin/sites/' + siteId + '/domains/' + domainId, {
      method: 'DELETE',
    }),
}