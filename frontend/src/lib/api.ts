const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api'

interface ApiError {
  error: string
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('token')
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error((data as ApiError).error || 'Request failed')
  }

  return data as T
}

// Auth API
export const authApi = {
  getOIDCUrl: () => `${API_BASE}/auth/oidc/authorize`,
  
  login: (userId: number, password: string) =>
    request<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, password }),
    }),

  logout: () =>
    request<{ message: string }>('/auth/logout', { method: 'POST' }),

  me: () => request<User>('/auth/me'),
}

// User API
export const userApi = {
  getProfile: () =>
    request<UserProfile>('/user/profile'),

  updatePassword: (password: string) =>
    request<{ message: string }>('/user/password', {
      method: 'PUT',
      body: JSON.stringify({ password }),
    }),
}

// Sites API
export const sitesApi = {
  list: () => request<Site[]>('/admin/sites'),

  create: (name: string, domains: string[], options?: {
    batch_mode_enabled?: boolean
    batch_difficulty?: number
    min_batch_count?: number
    max_batch_count?: number
  }) =>
    request<Site & { secret: string }>('/admin/sites', {
      method: 'POST',
      body: JSON.stringify({ name, domains, ...options }),
    }),

  get: (id: number) => request<Site>(`/admin/sites/${id}`),

  update: (id: number, data: Partial<Site>) =>
    request<Site>(`/admin/sites/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<{ message: string }>(`/admin/sites/${id}`, { method: 'DELETE' }),

  // Secret 管理
  getSecret: (id: number) =>
    request<{ secret: string }>(`/admin/sites/${id}/secret`),

  regenerateSecret: (id: number) =>
    request<{ secret: string }>(`/admin/sites/${id}/regenerate-secret`, {
      method: 'POST',
    }),

  // Domains
  listDomains: (siteId: number) =>
    request<SiteDomain[]>(`/admin/sites/${siteId}/domains`),

  addDomain: (siteId: number, domain: string) =>
    request<SiteDomain>(`/admin/sites/${siteId}/domains`, {
      method: 'POST',
      body: JSON.stringify({ domain }),
    }),

  deleteDomain: (siteId: number, domainId: number) =>
    request<{ message: string }>(`/admin/sites/${siteId}/domains/${domainId}`, {
      method: 'DELETE',
    }),

  // IP Policy
  getIpPolicy: (siteId: number) =>
    request<IpPolicy>(`/admin/sites/${siteId}/ip-policy`),

  updateIpPolicy: (siteId: number, data: Partial<IpPolicy>) =>
    request<IpPolicy>(`/admin/sites/${siteId}/ip-policy`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // IP Records
  listIpRecords: (siteId: number) =>
    request<IpRecord[]>(`/admin/sites/${siteId}/ip-records`),

  deleteIpRecord: (siteId: number, ip: string) =>
    request<{ message: string }>(`/admin/sites/${siteId}/ip-records/${ip}`, {
      method: 'DELETE',
    }),

  // Stats
  getStats: (siteId: number) =>
    request<SiteStats>(`/admin/sites/${siteId}/stats`),
}

// PoW API (for widget)
export const powApi = {
  createChallenge: (sitekey: string, origin: string, action?: string, difficulty?: number) =>
    request<ChallengeResponse>('/pow/challenge', {
      method: 'POST',
      body: JSON.stringify({ sitekey, origin, action, difficulty }),
    }),

  verify: (challengeId: string, nonce: string, iterations?: number, solveTimeMs?: number) =>
    request<{ token: string }>('/pow/verify', {
      method: 'POST',
      body: JSON.stringify({
        challenge_id: challengeId,
        nonce,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: false,
      }),
    }),

  // 批量模式验证
  verifyBatch: (challengeId: string, nonces: string[], iterations?: number, solveTimeMs?: number) =>
    request<{ token: string }>('/pow/verify', {
      method: 'POST',
      body: JSON.stringify({
        challenge_id: challengeId,
        nonces,
        iterations,
        solve_time_ms: solveTimeMs,
        batch_mode: true,
      }),
    }),
}

// Types
export interface User {
  id: number
  username: string
  email: string
  avatar: string
  role: string
}

export interface UserProfile extends User {
  has_password: boolean
  created_at: string
}

export interface Site {
  id: number
  name: string
  sitekey: string
  base_difficulty: number
  max_difficulty: number
  challenge_timeout: number
  token_timeout: number
  enabled: boolean
  // 批量模式配置
  batch_mode_enabled: boolean
  batch_difficulty: number
  min_batch_count: number
  max_batch_count: number
  domains: SiteDomain[]
  created_at: string
}

export interface SiteDomain {
  id: number
  site_id: number
  domain: string
  created_at: string
}

export interface IpPolicy {
  id: number
  site_id: number
  time_window_seconds: number
  request_threshold: number
  difficulty_increment: number  // 动态难度模式：难度增量
  batch_increment: number       // 批量模式：批次增量
  failure_weight: number
  enabled: boolean
}

export interface IpRecord {
  id: number
  site_id: number
  client_ip: string
  request_count: number
  failure_count: number
  current_difficulty: number
  window_start: string
  last_request_at: string
}

export interface SiteStats {
  total_challenges: number
  solved_challenges: number
  total_verifications: number
  success_verifications: number
  last_24h_challenges: number
  last_24h_solved: number
}

export interface ChallengeResponse {
  challenge_id: string
  challenge?: string       // 单个模式
  challenges?: string[]    // 批量模式
  difficulty: number
  batch_mode: boolean
  batch_total: number
}
