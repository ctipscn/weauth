export interface User {
  id: number
  username: string
  email: string
  avatar: string
  role: string
}

export interface AdminUser extends User {
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface AdminSite extends Site {
  user_id: number
  secret: string
  user?: User
  created_at: string
  updated_at: string
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
  difficulty_increment: number
  batch_increment: number
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

export interface SiteVerifyResponse {
  hostname: string
  action: string
  challenge_ts: string
}