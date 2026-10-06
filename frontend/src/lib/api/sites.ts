import { request, envelope } from './client'
import type { Site, SiteDomain, IpPolicy, IpRecord, SiteStats } from './types'

export const sitesApi = {
  list: () => request<Site[]>('/console/sites'),

  create: (name: string, domains: string[], options?: {
    batch_mode_enabled?: boolean
    batch_difficulty?: number
    min_batch_count?: number
    max_batch_count?: number
  }) =>
    request<Site & { secret: string }>('/console/sites', {
      method: 'POST',
      body: JSON.stringify(envelope({ name, domains, ...options })),
    }),

  get: (id: number) => request<Site>(`/console/sites/${id}`),

  update: (id: number, data: Partial<Site>) =>
    request<Site>(`/console/sites/${id}`, {
      method: 'PUT',
      body: JSON.stringify(envelope(data)),
    }),

  delete: (id: number) =>
    request<{ message: string }>(`/console/sites/${id}`, { method: 'DELETE' }),

  // Secret 管理
  getSecret: (id: number) =>
    request<{ secret: string }>(`/console/sites/${id}/secret`),

  regenerateSecret: (id: number) =>
    request<{ secret: string }>(`/console/sites/${id}/regenerate-secret`, {
      method: 'POST',
    }),

  // Domains
  listDomains: (siteId: number) =>
    request<SiteDomain[]>(`/console/sites/${siteId}/domains`),

  addDomain: (siteId: number, domain: string) =>
    request<SiteDomain>(`/console/sites/${siteId}/domains`, {
      method: 'POST',
      body: JSON.stringify(envelope({ domain })),
    }),

  deleteDomain: (siteId: number, domainId: number) =>
    request<{ message: string }>(`/console/sites/${siteId}/domains/${domainId}`, {
      method: 'DELETE',
    }),

  // IP Policy
  getIpPolicy: (siteId: number) =>
    request<IpPolicy>(`/console/sites/${siteId}/ip-policy`),

  updateIpPolicy: (siteId: number, data: Partial<IpPolicy>) =>
    request<IpPolicy>(`/console/sites/${siteId}/ip-policy`, {
      method: 'PUT',
      body: JSON.stringify(envelope(data)),
    }),

  // IP Records
  listIpRecords: (siteId: number) =>
    request<IpRecord[]>(`/console/sites/${siteId}/ip-records`),

  deleteIpRecord: (siteId: number, ip: string) =>
    request<{ message: string }>(`/console/sites/${siteId}/ip-records/${ip}`, {
      method: 'DELETE',
    }),

  // Stats
  getStats: (siteId: number) =>
    request<SiteStats>(`/console/sites/${siteId}/stats`),
}