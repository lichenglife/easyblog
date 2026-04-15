# 前端设计规范

## 1. 项目结构

### 1.1 目录组织

```
frontend/
├── public/                 # 静态资源
│   ├── favicon.ico
│   └── logo.png
├── src/
│   ├── api/                # API 请求封装
│   │   ├── request.ts      # Axios 实例配置
│   │   ├── user.ts         # 用户相关 API
│   │   ├── post.ts         # 文章相关 API
│   │   ├── comment.ts      # 评论相关 API
│   │   └── category.ts     # 分类相关 API
│   ├── assets/             # 资源文件
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/         # 全局样式
│   │       ├── variables.scss   # SCSS 变量
│   │       ├── mixins.scss      # SCSS 混入
│   │       └── global.scss      # 全局样式
│   ├── components/         # 通用组件
│   │   ├── common/         # 基础组件
│   │   │   ├── PageHeader.vue
│   │   │   ├── ImageUploader.vue
│   │   │   └── MarkdownEditor.vue
│   │   └── business/       # 业务组件
│   │       ├── PostCard.vue
│   │       ├── CommentItem.vue
│   │       └── UserAvatar.vue
│   ├── composables/        # 组合式函数
│   │   ├── usePagination.ts
│   │   ├── usePermission.ts
│   │   └── useUpload.ts
│   ├── constants/          # 常量定义
│   │   ├── routes.ts
│   │   ├── permissions.ts
│   │   └── config.ts
│   ├── layouts/            # 布局组件
│   │   ├── DefaultLayout.vue
│   │   ├── AdminLayout.vue
│   │   └── AuthLayout.vue
│   ├── router/             # 路由配置
│   │   ├── index.ts
│   │   ├── routes.ts
│   │   └── guards.ts
│   ├── stores/             # Pinia 状态管理
│   │   ├── user.ts
│   │   ├── post.ts
│   │   └── app.ts
│   ├── types/              # TypeScript 类型定义
│   │   ├── user.ts
│   │   ├── post.ts
│   │   ├── comment.ts
│   │   └── api.ts
│   ├── utils/              # 工具函数
│   │   ├── request.ts
│   │   ├── validate.ts
│   │   └── storage.ts
│   ├── directives/         # 自定义指令
│   │   └── permission.ts
│   ├── views/              # 页面组件
│   │   ├── home/
│   │   │   └── index.vue
│   │   ├── user/
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   └── Profile.vue
│   │   ├── post/
│   │   │   ├── List.vue
│   │   │   ├── Detail.vue
│   │   │   └── Editor.vue
│   │   └── admin/
│   │       ├── Dashboard.vue
│   │       └── Users.vue
│   ├── App.vue
│   └── main.ts
├── tests/
│   ├── unit/               # 单元测试
│   └── e2e/                # E2E 测试
├── .eslintrc.cjs
├── .prettierrc.json
├── tsconfig.json
├── vite.config.ts
└── package.json
```

---

## 2. Vue 组件规范

### 2.1 单文件组件结构

```vue
<script setup lang="ts">
// 1. 导入语句按顺序排列
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import type { User } from '@/types/user'

// 2. 类型定义
interface Props {
  title: string
  visible?: boolean
}

// 3. Props 定义
const props = withDefaults(defineProps<Props>(), {
  visible: false
})

// 4. Emits 定义
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'submit', data: FormData): void
}>()

// 5. 状态管理
const userStore = useUserStore()
const router = useRouter()
const loading = ref(false)
const formData = ref<FormData>({} as FormData)

// 6. 计算属性
const isValid = computed(() => {
  return formData.value.title !== ''
})

// 7. 监听器
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm()
  }
})

// 8. 生命周期
onMounted(() => {
  initData()
})

// 9. 方法定义
const handleSubmit = async () => {
  if (!isValid.value) return
  loading.value = true
  try {
    emit('submit', formData.value)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.value = {} as FormData
}

const initData = () => {
  // 初始化数据
}
</script>

<template>
  <div class="component-root">
    <!-- 模板内容 -->
  </div>
</template>

<style lang="scss" scoped>
.component-root {
  // 样式定义
}
</style>
```

### 2.2 组件命名规范

| 类型 | 命名方式 | 示例 |
|------|---------|------|
| 页面组件 | PascalCase + index.vue | `views/user/index.vue` |
| 通用组件 | PascalCase | `PostCard.vue`, `UserAvatar.vue` |
| 布局组件 | PascalCase + Layout | `AdminLayout.vue` |
| 组合式函数 | useXxx | `usePagination`, `usePermission` |

---

## 3. TypeScript 规范

### 3.1 类型定义

```typescript
// types/user.ts
export interface User {
  id: number
  username: string
  email: string
  avatar?: string
  role: 'admin' | 'user'
  status: 'pending' | 'active' | 'banned'
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// 使用类型别名区分不同用途
export type UserRole = User['role']
export type UserStatus = User['status']
```

### 3.2 API 响应类型

```typescript
// types/api.ts
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
```

---

## 4. API 请求规范

### 4.1 Axios 实例配置

```typescript
// api/request.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

const config = {
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
}

const createAxiosInstance = (): AxiosInstance => {
  const instance = axios.create(config)

  // 请求拦截器
  instance.interceptors.request.use(
    (config) => {
      const userStore = useUserStore()
      if (userStore.token) {
        config.headers.Authorization = `Bearer ${userStore.token}`
      }
      return config
    },
    (error) => Promise.reject(error)
  )

  // 响应拦截器
  instance.interceptors.response.use(
    (response: AxiosResponse<ApiResponse>) => {
      const { code, message, data } = response.data
      
      if (code !== 200) {
        ElMessage.error(message)
        // Token 过期处理
        if (code === 401) {
          const userStore = useUserStore()
          userStore.logout()
          router.push('/login')
        }
        return Promise.reject(new Error(message))
      }
      
      return data
    },
    (error) => {
      ElMessage.error(error.message || '请求失败')
      return Promise.reject(error)
    }
  )

  return instance
}

export const request = createAxiosInstance()
```

### 4.2 API 模块封装

```typescript
// api/user.ts
import { request } from './request'
import type { User, LoginRequest, LoginResponse } from '@/types/user'
import type { PageResponse } from '@/types/api'

export const userApi = {
  // 用户登录
  login: (data: LoginRequest): Promise<LoginResponse> => {
    return request.post('/auth/login', data)
  },

  // 用户注册
  register: (data: Partial<User>): Promise<User> => {
    return request.post('/users/register', data)
  },

  // 获取当前用户信息
  getCurrentUser: (): Promise<User> => {
    return request.get('/user/me')
  },

  // 更新用户信息
  updateProfile: (data: Partial<User>): Promise<User> => {
    return request.put('/user/profile', data)
  },

  // 修改密码
  changePassword: (oldPwd: string, newPwd: string): Promise<void> => {
    return request.post('/user/change-password', { oldPwd, newPwd })
  },

  // 用户列表（管理员）
  getUserList: (page: number, pageSize: number): Promise<PageResponse<User>> => {
    return request.get('/users', { params: { page, pageSize } })
  },

  // 审核用户
  auditUser: (id: number, approved: boolean): Promise<void> => {
    return request.post(`/users/${id}/audit`, { approved })
  },

  // 封禁用户
  banUser: (id: number): Promise<void> => {
    return request.post(`/users/${id}/ban`)
  },
}
```

---

## 5. 状态管理规范

### 5.1 Pinia Store 定义

```typescript
// stores/user.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api/user'
import type { User } from '@/types/user'
import { storage } from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string>(storage.get('token') || '')

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const avatar = computed(() => user.value?.avatar || '/default-avatar.png')

  // Actions
  const setToken = (newToken: string) => {
    token.value = newToken
    storage.set('token', newToken)
  }

  const fetchCurrentUser = async () => {
    try {
      user.value = await userApi.getCurrentUser()
    } catch (error) {
      console.error('Failed to fetch current user:', error)
      throw error
    }
  }

  const login = async (username: string, password: string) => {
    const { token: newToken, user: userData } = await userApi.login({ username, password })
    setToken(newToken)
    user.value = userData
  }

  const logout = () => {
    token.value = ''
    user.value = null
    storage.remove('token')
  }

  return {
    // State
    user,
    token,
    // Getters
    isLoggedIn,
    isAdmin,
    avatar,
    // Actions
    setToken,
    fetchCurrentUser,
    login,
    logout,
  }
})
```

---

## 6. 路由规范

### 6.1 路由配置

```typescript
// router/routes.ts
import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
  },
  {
    path: '/auth',
    component: () => import('@/layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'Login',
        component: () => import('@/views/user/Login.vue'),
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('@/views/user/Register.vue'),
      },
    ],
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/views/admin/Users.vue'),
      },
    ],
  },
]
```

### 6.2 路由守卫

```typescript
// router/guards.ts
import type { Router } from 'vue-router'
import { useUserStore } from '@/stores/user'

export const setupRouterGuards = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()

    // 需要认证的路由
    if (to.meta.requiresAuth && !userStore.isLoggedIn) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    // 需要管理员权限的路由
    if (to.meta.requiresAdmin && !userStore.isAdmin) {
      next({ name: 'Home' })
      return
    }

    // 已登录用户访问登录/注册页
    if (userStore.isLoggedIn && ['Login', 'Register'].includes(to.name as string)) {
      next({ name: 'Home' })
      return
    }

    next()
  })
}
```

---

## 7. 组件设计规范

### 7.1 基础组件设计原则

1. **单一职责**: 每个组件只负责一个功能
2. **Props 单向数据流**: 使用 Props 接收数据，使用 Emits 通知父组件
3. **可组合性**: 组件应该易于组合和复用
4. **可访问性**: 支持键盘操作和屏幕阅读器

### 7.2 图片上传组件

```vue
<!-- components/common/ImageUploader.vue -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElUpload, type UploadUserFile, type UploadProps } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

interface Props {
  modelValue?: string
  limit?: number
  maxSize?: number // MB
}

const props = withDefaults(defineProps<Props>(), {
  limit: 1,
  maxSize: 2,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const fileList = ref<UploadUserFile[]>([])

const uploadProps: Partial<UploadProps> = {
  limit: props.limit,
  accept: 'image/*',
  beforeUpload: (file) => {
    const isImage = file.type.startsWith('image/')
    const isLt2M = file.size / 1024 / 1024 < props.maxSize
    
    if (!isImage) {
      ElMessage.error('只能上传图片文件')
      return false
    }
    if (!isLt2M) {
      ElMessage.error(`图片大小不能超过 ${props.maxSize}MB`)
      return false
    }
    return true
  },
  onExceed: () => {
    ElMessage.warning(`最多只能上传 ${props.limit} 张图片`)
  },
}
</script>

<template>
  <el-upload
    v-model:file-list="fileList"
    v-bind="uploadProps"
    class="image-uploader"
  >
    <el-icon><Plus /></el-icon>
  </el-upload>
</template>

<style lang="scss" scoped>
.image-uploader {
  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    width: 148px;
    height: 148px;
    display: flex;
    justify-content: center;
    align-items: center;
    
    &:hover {
      border-color: var(--el-color-primary);
    }
  }
}
</style>
```

### 7.3 Markdown 编辑器组件

```vue
<!-- components/common/MarkdownEditor.vue -->
<script setup lang="ts">
import { ref, watch } from 'vue'
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

interface Props {
  modelValue: string
  height?: string
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  placeholder: '请输入内容...',
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const editorValue = ref(props.modelValue)

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== editorValue.value) {
      editorValue.value = newVal
    }
  }
)

const onUpdate = (value: string) => {
  emit('update:modelValue', value)
}
</script>

<template>
  <MdEditor
    v-model="editorValue"
    :height="height"
    :placeholder="placeholder"
    @update="onUpdate"
  />
</template>
```

---

## 8. 样式规范

### 8.1 SCSS 变量

```scss
// assets/styles/variables.scss

// 主题色
$primary-color: #409eff;
$success-color: #67c23a;
$warning-color: #e6a23c;
$danger-color: #f56c6c;
$info-color: #909399;

// 字体
$font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
  'Microsoft YaHei', Arial, sans-serif;
$font-size-base: 14px;
$font-size-small: 12px;
$font-size-large: 16px;

// 间距
$spacing-xs: 4px;
$spacing-sm: 8px;
$spacing-md: 16px;
$spacing-lg: 24px;
$spacing-xl: 32px;

// 边框
$border-color: #dcdfe6;
$border-radius: 4px;

// 阴影
$box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
```

### 8.2 样式命名规范

使用 BEM 命名规范：

```scss
// 块 (Block)
.post-card {}

// 元素 (Element)
.post-card__header {}
.post-card__title {}
.post-card__content {}
.post-card__footer {}

// 修饰符 (Modifier)
.post-card--featured {}
.post-card--sticky {}
```

---

## 9. 工具函数规范

### 9.1 存储工具

```typescript
// utils/storage.ts
const STORAGE_PREFIX = 'easyblog_'

export const storage = {
  get: <T>(key: string): T | null => {
    try {
      const value = localStorage.getItem(STORAGE_PREFIX + key)
      return value ? JSON.parse(value) : null
    } catch (error) {
      console.error('Storage get error:', error)
      return null
    }
  },

  set: <T>(key: string, value: T): void => {
    try {
      localStorage.setItem(STORAGE_PREFIX + key, JSON.stringify(value))
    } catch (error) {
      console.error('Storage set error:', error)
    }
  },

  remove: (key: string): void => {
    try {
      localStorage.removeItem(STORAGE_PREFIX + key)
    } catch (error) {
      console.error('Storage remove error:', error)
    }
  },

  clear: (): void => {
    try {
      Object.keys(localStorage).forEach((key) => {
        if (key.startsWith(STORAGE_PREFIX)) {
          localStorage.removeItem(key)
        }
      })
    } catch (error) {
      console.error('Storage clear error:', error)
    }
  },
}
```

### 9.2 验证工具

```typescript
// utils/validate.ts
export const validate = {
  // 邮箱验证
  email: (email: string): boolean => {
    const reg = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    return reg.test(email)
  },

  // 密码验证（6-20 位，包含字母和数字）
  password: (password: string): boolean => {
    const reg = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,20}$/
    return reg.test(password)
  },

  // 用户名验证（3-20 位字母、数字、下划线）
  username: (username: string): boolean => {
    const reg = /^[a-zA-Z0-9_]{3,20}$/
    return reg.test(username)
  },

  // URL 验证
  url: (url: string): boolean => {
    const reg = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/
    return reg.test(url)
  },
}
```

---

## 10. 测试规范

### 10.1 组件测试

```typescript
// tests/unit/components/PostCard.test.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PostCard from '@/components/business/PostCard.vue'
import { createTestingPinia } from '@pinia/testing'

describe('PostCard', () => {
  const mockPost = {
    id: 1,
    title: 'Test Post',
    summary: 'Test summary',
    coverImage: '/test.jpg',
    author: { id: 1, username: 'author' },
    viewCount: 100,
    likeCount: 10,
    createdAt: '2024-01-01',
  }

  it('renders post title correctly', () => {
    const wrapper = mount(PostCard, {
      props: { post: mockPost },
      global: { plugins: [createTestingPinia()] },
    })

    expect(wrapper.find('.post-card__title').text()).toBe('Test Post')
  })

  it('emits click event when clicked', async () => {
    const wrapper = mount(PostCard, {
      props: { post: mockPost },
      global: { plugins: [createTestingPinia()] },
    })

    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(1)
  })

  it('displays view count and like count', () => {
    const wrapper = mount(PostCard, {
      props: { post: mockPost },
      global: { plugins: [createTestingPinia()] },
    })

    expect(wrapper.text()).toContain('100')
    expect(wrapper.text()).toContain('10')
  })
})
```

### 10.2 E2E 测试

```typescript
// tests/e2e/login.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Login Flow', () => {
  test('should login successfully with valid credentials', async ({ page }) => {
    await page.goto('/login')

    await page.fill('input[name="username"]', 'testuser')
    await page.fill('input[name="password"]', 'Test123456')
    await page.click('button[type="submit"]')

    await expect(page).toHaveURL('/')
    await expect(page.locator('.user-avatar')).toBeVisible()
  })

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/login')

    await page.fill('input[name="username"]', 'invalid')
    await page.fill('input[name="password"]', 'wrong')
    await page.click('button[type="submit"]')

    await expect(page.locator('.el-message--error')).toBeVisible()
  })
})
```

---

## 11. 代码质量工具

### 11.1 ESLint 配置

```javascript
// .eslintrc.cjs
module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-recommended',
    'plugin:@typescript-eslint/recommended',
    'prettier',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 'latest',
    parser: '@typescript-eslint/parser',
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint', 'vue'],
  rules: {
    'vue/multi-word-component-names': 'off',
    '@typescript-eslint/no-explicit-any': 'warn',
    '@typescript-eslint/explicit-module-boundary-types': 'off',
  },
}
```

### 11.2 Prettier 配置

```json
// .prettierrc.json
{
  "semi": false,
  "singleQuote": true,
  "printWidth": 100,
  "tabWidth": 2,
  "trailingComma": "es5",
  "arrowParens": "always"
}
```

---

## 12. 提交前检查清单

### 代码质量
- [ ] 代码已 `prettier` 格式化
- [ ] 通过 `eslint` 检查
- [ ] 无 TypeScript 类型错误

### 测试
- [ ] 组件测试通过
- [ ] 核心组件覆盖率≥80%
- [ ] E2E 测试通过（核心流程）

### 功能
- [ ] 页面响应式布局正常
- [ ] 跨浏览器测试通过（Chrome/Firefox/Safari）
- [ ] 无控制台警告和错误

### 性能
- [ ] 首屏加载时间 < 2s
- [ ] 组件按需加载
- [ ] 图片资源压缩优化

---

## 13. 常用命令

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建
npm run build

# 预览构建结果
npm run preview

# 运行单元测试
npm run test:unit

# 运行 E2E 测试
npm run test:e2e

# 代码检查
npm run lint
npm run lint:fix

# 类型检查
npm run type-check
```

---

*最后更新：2026-04-01*
