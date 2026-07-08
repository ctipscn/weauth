import { useEffect } from 'react'
import { Routes, Route, Navigate, useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from './lib/auth'
import Landing from './pages/Landing'
import Login from './pages/Login'
import Sites from './pages/Sites'
import SiteDetail from './pages/SiteDetail'
import Profile from './pages/Profile'
import Demo from './pages/Demo'
import Privacy from './pages/Privacy'
import Docs from './pages/Docs'
import Layout from './components/Layout'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token, isLoading, checkAuth } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) {
      navigate('/login')
    }
  }, [token, navigate])

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  if (!token) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

function AuthCallback() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { login, checkAuth } = useAuth()

  useEffect(() => {
    const token = searchParams.get('token')
    if (token) {
      localStorage.setItem('token', token)
      login(token, { id: 0, username: '', email: '', avatar: '', role: '' })
      checkAuth().then(() => {
        navigate('/console/sites', { replace: true })
      })
    } else {
      navigate('/login', { replace: true })
    }
  }, [searchParams, navigate, login, checkAuth])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>
  )
}

export default function App() {
  return (
    <Routes>
      {/* Public routes */}
      <Route path="/" element={<Landing />} />
      <Route path="/login" element={<Login />} />
      <Route path="/dashboard" element={<AuthCallback />} />
      <Route path="/demo" element={<Demo />} />
      <Route path="/privacy" element={<Privacy />} />
      <Route path="/docs" element={<Docs />} />
      
      {/* Protected routes */}
      <Route
        path="/console"
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="/console/sites" replace />} />
        <Route path="sites" element={<Sites />} />
        <Route path="sites/:id" element={<SiteDetail />} />
        <Route path="profile" element={<Profile />} />
      </Route>
      
      {/* Fallback: redirect old paths to new console paths for backward compatibility */}
      <Route path="/sites" element={<Navigate to="/console/sites" replace />} />
      <Route path="/profile" element={<Navigate to="/console/profile" replace />} />
    </Routes>
  )
}
