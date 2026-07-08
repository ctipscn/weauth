import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { cn } from '@/lib/utils'
import {
  Globe,
  User,
  LogOut,
  Menu,
  X,
  ChevronDown,
  Settings,
  BarChart3,
  ExternalLink,
} from 'lucide-react'
import { useState, useRef, useEffect } from 'react'

// WeAuth Logo Component
function WeAuthLogo({ className = "w-4 h-4" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

const navItems = [
  { 
    path: '/console/sites', 
    label: '站点管理', 
    icon: Globe,
    description: '管理您的验证站点'
  },
  { 
    path: '/console/profile', 
    label: '个人中心', 
    icon: User,
    description: '账户设置与安全'
  },
]

export default function Layout() {
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuth()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
  const userMenuRef = useRef<HTMLDivElement>(null)

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  // Close user menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setUserMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30">
      {/* Top Navigation */}
      <header className="fixed top-0 left-0 right-0 h-14 bg-white/90 backdrop-blur-md border-b border-slate-200/60 z-50">
        <div className="h-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex items-center justify-between">
          {/* Logo */}
          <div className="flex items-center gap-8">
            <Link to="/" className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-lg bg-blue-50 border border-blue-200 flex items-center justify-center">
                <WeAuthLogo className="w-4 h-4 text-blue-600" />
              </div>
              <div className="flex items-baseline gap-1.5">
                <span className="text-slate-900 font-semibold hidden sm:block tracking-tight">微验</span>
                <span className="text-slate-500 text-sm hidden sm:block">WeAuth</span>
              </div>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center gap-1">
              {navItems.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  className={cn(
                    'px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                    location.pathname.startsWith(item.path)
                      ? 'bg-[#0066ff]/10 text-[#0066ff]'
                      : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
                  )}
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </div>

          {/* Right Side */}
          <div className="flex items-center gap-4">
            {/* Demo Link */}
            <Link
              to="/demo"
              className="hidden sm:flex items-center gap-1 text-sm text-slate-500 hover:text-slate-700 transition-colors"
            >
              Demo
              <ExternalLink className="w-3 h-3" />
            </Link>

            {/* User Menu */}
            <div className="relative" ref={userMenuRef}>
              <button
                onClick={() => setUserMenuOpen(!userMenuOpen)}
                className="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-slate-100 transition-colors"
              >
                {user?.avatar ? (
                  <img
                    src={user.avatar}
                    alt={user.username}
                    className="w-7 h-7 rounded-full"
                  />
                ) : (
                  <div className="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-xs font-medium text-white shadow-sm">
                    {user?.username?.charAt(0) || 'U'}
                  </div>
                )}
                <span className="text-sm text-slate-700 hidden sm:block max-w-[100px] truncate">
                  {user?.username || '用户'}
                </span>
                <ChevronDown className={cn(
                  "w-4 h-4 text-slate-400 transition-transform hidden sm:block",
                  userMenuOpen && "rotate-180"
                )} />
              </button>

              {/* Dropdown */}
              {userMenuOpen && (
                <div className="absolute right-0 mt-1 w-52 bg-white rounded-lg border border-slate-200 shadow-lg overflow-hidden">
                  <div className="px-3 py-2.5 border-b border-slate-100 bg-slate-50">
                    <p className="text-sm font-medium text-slate-900 truncate">{user?.username}</p>
                    <p className="text-xs text-slate-500 truncate">{user?.email}</p>
                    <p className="text-xs text-slate-400 mt-0.5">ID: {user?.id}</p>
                  </div>
                  <div className="py-1">
                    <Link
                      to="/console/profile"
                      onClick={() => setUserMenuOpen(false)}
                      className="flex items-center gap-2 px-3 py-2 text-sm text-slate-600 hover:text-slate-900 hover:bg-slate-50 transition-colors"
                    >
                      <Settings className="w-4 h-4" />
                      账户设置
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-500 hover:text-red-600 hover:bg-red-50 transition-colors"
                    >
                      <LogOut className="w-4 h-4" />
                      退出登录
                    </button>
                  </div>
                </div>
              )}
            </div>

            {/* Mobile Menu Button */}
            <button
              className="md:hidden p-1.5 text-slate-500 hover:text-slate-700 transition-colors"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </header>

      {/* Mobile Navigation */}
      {mobileMenuOpen && (
        <div className="fixed inset-0 top-14 bg-white z-40 md:hidden">
          <nav className="p-4 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                onClick={() => setMobileMenuOpen(false)}
                className={cn(
                  'flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors',
                  location.pathname.startsWith(item.path)
                    ? 'bg-[#0066ff]/10 text-[#0066ff]'
                    : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                )}
              >
                <item.icon className="w-5 h-5" />
                <div>
                  <p className="font-medium">{item.label}</p>
                  <p className="text-xs text-slate-400">{item.description}</p>
                </div>
              </Link>
            ))}
            <Link
              to="/demo"
              onClick={() => setMobileMenuOpen(false)}
              className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-600 hover:bg-slate-50 hover:text-slate-900 transition-colors"
            >
              <BarChart3 className="w-5 h-5" />
              <div>
                <p className="font-medium">Demo</p>
                <p className="text-xs text-slate-400">查看验证组件演示</p>
              </div>
            </Link>
          </nav>
        </div>
      )}

      {/* Main Content */}
      <main className="pt-14 min-h-screen">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
