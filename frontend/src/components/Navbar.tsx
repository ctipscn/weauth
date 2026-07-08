import { Link } from 'react-router-dom'

// WeAuth Logo
function WeAuthLogo({ className = "w-5 h-5" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

interface NavbarProps {
  showMenu?: boolean
}

export default function Navbar({ showMenu = true }: NavbarProps) {
  return (
    <nav className="fixed top-0 w-full z-50 bg-white/80 backdrop-blur-md border-b border-slate-200/60">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-50 to-cyan-50 border border-blue-200/50 flex items-center justify-center">
              <WeAuthLogo className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <div className="flex items-baseline gap-2">
                <span className="text-xl font-bold text-slate-900">微验</span>
                <span className="text-sm text-slate-500">WeAuth</span>
              </div>
            </div>
          </Link>
          
          {showMenu && (
            <div className="hidden md:flex items-center gap-8">
              <a href="/#features" className="text-sm text-slate-600 hover:text-slate-900 transition-colors">
                特性
              </a>
              <a href="/#pricing" className="text-sm text-slate-600 hover:text-slate-900 transition-colors">
                定价
              </a>
              <Link to="/demo" className="text-sm text-slate-600 hover:text-slate-900 transition-colors">
                演示
              </Link>
              <Link to="/docs" className="text-sm text-slate-600 hover:text-slate-900 transition-colors">
                文档
              </Link>
              <a 
                href="https://cloud.emoera.com/" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-sm text-slate-600 hover:text-slate-900 transition-colors flex items-center gap-1"
              >
                E时代云服务
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </a>
            </div>
          )}

            <div className="flex items-center gap-3">
              <Link
                to="/login"
                className="text-sm font-medium text-slate-600 hover:text-slate-900 transition-colors"
              >
                登录
              </Link>
              <Link
                to="/console"
                className="px-4 py-2 bg-white border-2 border-blue-600 text-blue-600 text-sm font-medium rounded-lg hover:bg-blue-50 hover:border-blue-700 hover:text-blue-700 transition-all"
              >
                控制台
              </Link>
            </div>
        </div>
      </div>
    </nav>
  )
}
