import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi } from '@/lib/api'
import { User, Key, Check, Loader2, Calendar, AlertCircle } from 'lucide-react'

// WeAuth Logo
function WeAuthLogo({ className = "w-4 h-4" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

export default function Profile() {
  const queryClient = useQueryClient()
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  const { data: profile, isLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: userApi.getProfile,
  })

  const updatePasswordMutation = useMutation({
    mutationFn: userApi.updatePassword,
    onSuccess: () => {
      setMessage({ type: 'success', text: '密码更新成功' })
      setPassword('')
      setConfirmPassword('')
      queryClient.invalidateQueries({ queryKey: ['profile'] })
    },
    onError: (error) => {
      setMessage({ type: 'error', text: error instanceof Error ? error.message : '更新失败' })
    },
  })

  const handleUpdatePassword = (e: React.FormEvent) => {
    e.preventDefault()
    setMessage(null)

    if (password.length < 6) {
      setMessage({ type: 'error', text: '密码至少需要 6 个字符' })
      return
    }

    if (password !== confirmPassword) {
      setMessage({ type: 'error', text: '两次输入的密码不一致' })
      return
    }

    updatePasswordMutation.mutate(password)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-6 w-6 border-2 border-[#0066ff] border-t-transparent"></div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-slate-900">个人中心</h1>
        <p className="text-slate-500 text-sm mt-0.5">管理您的账户信息</p>
      </div>

      {/* Profile Card */}
      <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
        {/* Avatar Section */}
        <div className="p-5 border-b border-slate-100 bg-slate-50">
          <div className="flex items-center gap-4">
            {profile?.avatar ? (
              <img
                src={profile.avatar}
                alt={profile.username}
                className="h-14 w-14 rounded-xl"
              />
                ) : (
                  <div className="h-14 w-14 rounded-xl bg-blue-600 flex items-center justify-center text-xl font-semibold text-white shadow-sm">
                    {profile?.username?.charAt(0) || 'U'}
                  </div>
                )}
            <div>
              <h2 className="text-lg font-semibold text-slate-900">{profile?.username}</h2>
              <p className="text-slate-500 text-sm">{profile?.email}</p>
            </div>
          </div>
        </div>

        {/* Info Grid */}
        <div className="divide-y divide-slate-100">
          <div className="flex items-center justify-between px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                <User className="w-4 h-4 text-slate-500" />
              </div>
              <span className="text-sm text-slate-600">用户 ID</span>
            </div>
            <code className="px-2.5 py-1 bg-[#0066ff]/10 text-[#0066ff] rounded-md text-base font-bold font-mono">
              {profile?.id}
            </code>
          </div>

          <div className="flex items-center justify-between px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                <WeAuthLogo className="w-4 h-4 text-slate-500" />
              </div>
              <span className="text-sm text-slate-600">角色</span>
            </div>
            <span className={`px-2 py-0.5 rounded text-sm font-medium ${
              profile?.role === 'admin' 
                ? 'bg-purple-50 text-purple-600' 
                : 'bg-slate-100 text-slate-600'
            }`}>
              {profile?.role === 'admin' ? '管理员' : '用户'}
            </span>
          </div>

          <div className="flex items-center justify-between px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                <Calendar className="w-4 h-4 text-slate-500" />
              </div>
              <span className="text-sm text-slate-600">注册时间</span>
            </div>
            <span className="text-slate-900 text-sm">
              {profile?.created_at && new Date(profile.created_at).toLocaleDateString('zh-CN')}
            </span>
          </div>

          <div className="flex items-center justify-between px-5 py-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center">
                <Key className="w-4 h-4 text-slate-500" />
              </div>
              <span className="text-sm text-slate-600">备用密码</span>
            </div>
            <span className={`px-2 py-0.5 rounded text-sm font-medium ${
              profile?.has_password 
                ? 'bg-emerald-50 text-emerald-600' 
                : 'bg-slate-100 text-slate-500'
            }`}>
              {profile?.has_password ? '已设置' : '未设置'}
            </span>
          </div>
        </div>
      </div>

      {/* Password Settings */}
      <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
        <div className="px-5 py-3 border-b border-slate-100">
          <h2 className="font-medium text-slate-900 text-sm flex items-center gap-2">
            <Key className="w-4 h-4 text-[#0066ff]" />
            {profile?.has_password ? '修改密码' : '设置密码'}
          </h2>
          <p className="text-xs text-slate-500 mt-0.5">
            设置后可使用 用户ID + 密码 登录
          </p>
        </div>

        <form onSubmit={handleUpdatePassword} className="p-5 space-y-4">
          {message && (
            <div className={`flex items-center gap-2 p-2.5 rounded-lg text-sm ${
              message.type === 'success'
                ? 'bg-emerald-50 border border-emerald-100 text-emerald-600'
                : 'bg-red-50 border border-red-100 text-red-600'
            }`}>
              {message.type === 'success' ? (
                <Check className="w-4 h-4" />
              ) : (
                <AlertCircle className="w-4 h-4" />
              )}
              {message.text}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1.5">
              {profile?.has_password ? '新密码' : '密码'}
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="至少 6 个字符"
              className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-[#0066ff] focus:ring-1 focus:ring-[#0066ff] transition-all"
              required
              minLength={6}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1.5">
              确认密码
            </label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="再次输入"
              className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-[#0066ff] focus:ring-1 focus:ring-[#0066ff] transition-all"
              required
              minLength={6}
            />
          </div>

          <button
            type="submit"
            disabled={updatePasswordMutation.isPending}
            className="w-full h-9 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-lg text-sm font-medium transition-all duration-150 disabled:opacity-50 flex items-center justify-center shadow-sm"
          >
            {updatePasswordMutation.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <>
                <Check className="h-4 w-4 mr-1.5" />
                {profile?.has_password ? '更新密码' : '设置密码'}
              </>
            )}
          </button>
        </form>
      </div>

      {/* Usage Guide */}
      <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow p-5">
        <h2 className="font-medium text-slate-900 text-sm mb-3">备用登录说明</h2>
        <div className="space-y-2 text-sm text-slate-600">
          <div className="flex items-start gap-2">
            <span className="w-5 h-5 rounded-full bg-[#0066ff]/10 text-[#0066ff] flex items-center justify-center text-xs font-semibold flex-shrink-0">1</span>
            <p>记住用户 ID: <code className="px-1.5 py-0.5 bg-slate-100 rounded text-[#0066ff] font-bold">{profile?.id}</code></p>
          </div>
          <div className="flex items-start gap-2">
            <span className="w-5 h-5 rounded-full bg-[#0066ff]/10 text-[#0066ff] flex items-center justify-center text-xs font-semibold flex-shrink-0">2</span>
            <p>在登录页展开「用户 ID 登录」</p>
          </div>
          <div className="flex items-start gap-2">
            <span className="w-5 h-5 rounded-full bg-[#0066ff]/10 text-[#0066ff] flex items-center justify-center text-xs font-semibold flex-shrink-0">3</span>
            <p>输入用户 ID 和密码</p>
          </div>
        </div>
      </div>
    </div>
  )
}
