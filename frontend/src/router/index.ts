import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/Home.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/user/LoginPage.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/pages/user/RegisterPage.vue'),
    },
    {
      path: '/articles',
      name: 'articles',
      component: () => import('@/pages/article/ArticleList.vue'),
    },
    {
      path: '/articles/:id',
      name: 'article-detail',
      component: () => import('@/pages/article/ArticleDetail.vue'),
    },
    {
      path: '/editor',
      name: 'editor-create',
      component: () => import('@/pages/article/ArticleEditor.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/editor/:id',
      name: 'editor-edit',
      component: () => import('@/pages/article/ArticleEditor.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/pages/user/ProfilePage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/pages/admin/AdminLayout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('@/pages/admin/UserManagement.vue'),
        },
        {
          path: 'articles',
          name: 'admin-articles',
          component: () => import('@/pages/admin/ArticleManagement.vue'),
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: () => import('@/pages/admin/CategoryManagement.vue'),
        },
        {
          path: 'tags',
          name: 'admin-tags',
          component: () => import('@/pages/admin/TagManagement.vue'),
        },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, _, next) => {
  const userStore = useUserStore()

  // 检查是否需要认证
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'home' })
    return
  }

  // 已登录用户访问登录/注册页，重定向到首页
  if (userStore.isLoggedIn && ['login', 'register'].includes(to.name as string)) {
    next({ name: 'home' })
    return
  }

  next()
})

export default router
