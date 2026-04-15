import { http } from './request'

// 用户登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 用户注册请求参数
export interface RegisterParams {
  username: string
  password: string
  email: string
  nick_name: string
  phone?: string
}

// 用户信息
export interface UserInfo {
  userID: string
  username: string
  nickname: string
  email: string
  phone?: string
  avatar?: string
  bio?: string
  role: number
  status: number
  createAt: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: UserInfo
}

// 更新用户请求参数
export interface UpdateUserParams {
  nickname?: string
  email?: string
  phone?: string
}

// 修改密码请求参数
export interface ChangePasswordParams {
  oldPassword: string
  newPassword: string
}

// 审核用户请求参数
export interface AuditUserParams {
  userId: number
  status: number // 1-通过 2-拒绝
}

// 用户 API
export const userApi = {
  // 登录
  login(data: LoginParams): Promise<LoginResponse> {
    return http.post('/user/login', data)
  },

  // 注册
  register(data: RegisterParams): Promise<{ user_id: string }> {
    return http.post('/user', data)
  },

  // 获取当前用户信息
  getCurrentUser(): Promise<UserInfo> {
    return http.get('/user/info')
  },

  // 更新用户信息
  updateUserInfo(username: string, data: UpdateUserParams): Promise<void> {
    return http.put(`/user/${username}`, data)
  },

  // 修改密码
  changePassword(data: ChangePasswordParams): Promise<void> {
    return http.post('/user/change-password', data)
  },

  // 登出
  logout(): Promise<void> {
    return http.post('/user/logout')
  },

  // 获取用户列表
  getUserList(params?: { page?: number; pageSize?: number }): Promise<{
    totalCount: number
    hasMore: boolean
    users: UserInfo[]
  }> {
    return http.get('/user/list', { params })
  },

  // 根据 ID 获取用户
  getUserById(id: number): Promise<UserInfo> {
    return http.get(`/user/${id}`)
  },

  // 删除用户
  deleteUser(id: number): Promise<void> {
    return http.delete(`/user/${id}`)
  },

  // 审核用户（管理员）
  auditUser(data: AuditUserParams): Promise<void> {
    return http.post('/admin/user/audit', data)
  },

  // 封禁用户（管理员）
  banUser(userId: number): Promise<void> {
    return http.post('/admin/user/ban', { userId })
  },
}
