import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { authApi } from '@/lib/api'
import { motion } from 'framer-motion'
import { Loader2, ChevronDown, ChevronUp } from 'lucide-react'
import Navbar from '@/components/Navbar'

// WeAuth Logo Component
function WeAuthLogo({ className = "w-4 h-4" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

export default function Login() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { login, token } = useAuth()

  const [showBackupLogin, setShowBackupLogin] = useState(false)
  const [userId, setUserId] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(searchParams.get('error') || '')

  // 如果已登录，跳转到控制台
  useEffect(() => {
    if (token) {
      navigate('/console/sites')
    }
  }, [token, navigate])

  const handleOIDCLogin = () => {
    window.location.href = authApi.getOIDCUrl()
  }

  const handlePasswordLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      const data = await authApi.login(parseInt(userId), password)
      login(data.token, data.user)
      navigate('/console/sites')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-violet-50/30 to-purple-50/30 flex flex-col">
      {/* Top Bar */}
      <Navbar showMenu={false} />

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="w-full max-w-md"
        >
          {/* Decorative elements - 更轻微的装饰 */}
          <div className="relative">
            <div className="absolute -top-32 -left-32 w-64 h-64 bg-violet-100/20 rounded-full blur-3xl"></div>
            <div className="absolute -bottom-32 -right-32 w-64 h-64 bg-purple-100/20 rounded-full blur-3xl"></div>
            
            {/* Card */}
            <div className="relative bg-white/90 backdrop-blur-sm rounded-3xl border border-slate-200 shadow-xl overflow-hidden">
              {/* Card Header */}
              <div className="px-8 py-10 text-center">
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ delay: 0.2, duration: 0.5 }}
                  className="w-16 h-16 mx-auto mb-6 rounded-2xl bg-violet-50 border border-violet-100 flex items-center justify-center"
                >
                  <WeAuthLogo className="w-8 h-8 text-violet-600" />
                </motion.div>
                
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.3 }}
                >
                  <h1 className="text-2xl font-bold text-slate-900 tracking-tight mb-2">
                    登录到控制台
                  </h1>
                  <p className="text-sm text-slate-500">
                    管理您的 PoW 验证服务
                  </p>
                </motion.div>
              </div>

              {/* Card Body */}
              <div className="px-8 pb-8">
                {/* Error Message */}
                {error && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 text-sm"
                  >
                    {error}
                  </motion.div>
                )}

                {/* Primary OIDC Login */}
                <motion.button
                  whileHover={{ scale: 1.02, y: -2 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={handleOIDCLogin}
                  className="w-full h-12 rounded-xl font-medium transition-all bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 shadow-sm hover:shadow-md"
                >
                  <div className="flex items-center justify-center gap-2.5">
                    <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z" />
                    </svg>
                    使用 E时代通行证 登录
                  </div>
                </motion.button>

                <p className="text-center text-xs text-slate-500 mt-4">
                  首次登录将自动创建账户
                </p>

                {/* Divider with Backup Login Toggle */}
                <div className="mt-6 pt-6 border-t border-slate-200/60">
                  <button
                    onClick={() => setShowBackupLogin(!showBackupLogin)}
                    className="w-full flex items-center justify-center gap-1.5 text-xs text-slate-500 hover:text-slate-700 transition-colors group"
                  >
                    {showBackupLogin ? (
                      <ChevronUp className="w-4 h-4 transition-transform group-hover:-translate-y-0.5" />
                    ) : (
                      <ChevronDown className="w-4 h-4 transition-transform group-hover:translate-y-0.5" />
                    )}
                    使用用户 ID 登录（备用）
                  </button>

                  {/* Backup Password Login */}
                  {showBackupLogin && (
                    <motion.form
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      exit={{ opacity: 0, height: 0 }}
                      onSubmit={handlePasswordLogin}
                      className="mt-5 space-y-4"
                    >
                      <div>
                        <label className="block text-xs font-medium text-slate-700 mb-2">
                          用户 ID
                        </label>
                        <input
                          type="number"
                          value={userId}
                          onChange={(e) => setUserId(e.target.value)}
                          placeholder="输入您的用户 ID"
                          className="w-full h-11 px-4 rounded-xl bg-slate-50/50 border border-slate-200/60 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-violet-400 focus:ring-4 focus:ring-violet-100 transition-all backdrop-blur-sm"
                          required
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-medium text-slate-700 mb-2">
                          密码
                        </label>
                        <input
                          type="password"
                          value={password}
                          onChange={(e) => setPassword(e.target.value)}
                          placeholder="输入您的密码"
                          className="w-full h-11 px-4 rounded-xl bg-slate-50/50 border border-slate-200/60 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-violet-400 focus:ring-4 focus:ring-violet-100 transition-all backdrop-blur-sm"
                          required
                          minLength={1}
                        />
                      </div>
                      <motion.button
                        whileHover={{ scale: 1.02, y: -2 }}
                        whileTap={{ scale: 0.98 }}
                        type="submit"
                        disabled={loading}
                        className="w-full h-11 rounded-xl text-sm font-medium transition-all disabled:opacity-50 disabled:cursor-not-allowed bg-slate-900 hover:bg-slate-800 text-white shadow-sm"
                      >
                        <div className="flex items-center justify-center">
                          {loading ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                          ) : (
                            '登录'
                          )}
                        </div>
                      </motion.button>
                      <p className="text-center text-xs text-slate-500">
                        需要先使用 E时代通行证 登录并设置密码
                      </p>
                    </motion.form>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Footer */}
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="text-center text-xs text-slate-500 mt-8"
          >
            © {new Date().getFullYear()} 微验 WeAuth
          </motion.p>
        </motion.div>
      </main>
    </div>
  )
}
