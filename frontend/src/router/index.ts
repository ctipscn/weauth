import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/lib/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/pages/public/Landing.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/auth/Login.vue'),
    },
    {
      path: '/dashboard',
      name: 'auth-callback',
      component: () => import('@/pages/auth/AuthCallback.vue'),
    },
    {
      path: '/auth/callback',
      // OIDC 回调：后端验证后重定向到此地址并附带 ?token=<jwt>
      redirect: (to) => ({ path: '/dashboard', query: { token: to.query.token as string } }),
    },
    {
      path: '/demo',
      name: 'demo',
      component: () => import('@/pages/public/Demo.vue'),
    },
    {
      path: '/privacy',
      name: 'privacy',
      component: () => import('@/pages/public/Privacy.vue'),
    },
    {
      path: '/docs',
      name: 'docs',
      component: () => import('@/pages/public/Docs.vue'),
    },
    {
      path: '/console',
      component: () => import('@/components/Layout.vue'),
      redirect: '/console/sites',
      children: [
        {
          path: 'sites',
          name: 'admin-sites',
          component: () => import('@/pages/console/Sites.vue'),
        },
        {
          path: 'sites/:id',
          name: 'admin-site-detail',
          component: () => import('@/pages/console/SiteDetail.vue'),
        },
        {
          path: 'profile',
          name: 'admin-profile',
          component: () => import('@/pages/console/Profile.vue'),
        },
        {
          path: 'privacy',
          name: 'admin-privacy',
          component: () => import('@/pages/console/Privacy.vue'),
        },
        // ---- Admin-only routes ----
        {
          path: 'admin',
          redirect: '/console/admin/dashboard',
        },
        {
          path: 'admin/dashboard',
          name: 'admin-dashboard',
          component: () => import('@/pages/admin/AdminDashboard.vue'),
        },
        {
          path: 'admin/users',
          name: 'admin-users',
          component: () => import('@/pages/admin/AdminUsers.vue'),
        },
        {
          path: 'admin/sites',
          name: 'admin-all-sites',
          component: () => import('@/pages/admin/AdminSites.vue'),
        },
        {
          path: 'admin/sites/:id',
          name: 'admin-site-detail-admin',
          component: () => import('@/pages/admin/AdminSiteDetail.vue'),
        },
      ],
    },
    // Fallback redirects
    {
      path: '/sites',
      redirect: '/console/sites',
    },
    {
      path: '/profile',
      redirect: '/console/profile',
    },
  ],
})

// 导航守卫：阻止已登录用户进入登录页，阻止未登录用户进入控制台
router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()

  // 已登录用户不能访问登录页
  if (to.path === '/login' && auth.isAuthenticated) {
    next('/console/sites')
    return
  }

  // 未登录用户不能访问控制台
  if (to.path.startsWith('/console') && !auth.isAuthenticated) {
    next('/login')
    return
  }

  // 非管理员不能访问管理后台
  if (to.path.startsWith('/console/admin') && auth.user?.role !== 'admin') {
    next('/console/sites')
    return
  }

  next()
})

export default router