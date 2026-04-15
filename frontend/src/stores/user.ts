import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi, type UserInfo } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 2)

  // 设置用户信息
  const setUserInfo = (user: UserInfo) => {
    userInfo.value = user
  }

  // 设置 Token
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 登录
  const login = async (username: string, password: string) => {
    const response = await userApi.login({ username, password })
    setToken(response.token)
    setUserInfo(response.user)
    localStorage.setItem('userInfo', JSON.stringify(response.user))
  }

  // 注册
  const register = async (data: {
    username: string
    password: string
    email: string
    nick_name: string
    phone?: string
  }) => {
    return userApi.register(data)
  }

  // 获取当前用户信息
  const fetchCurrentUser = async () => {
    try {
      const user = await userApi.getCurrentUser()
      setUserInfo(user)
      localStorage.setItem('userInfo', JSON.stringify(user))
      return user
    } catch (error) {
      logout()
      throw error
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  // 初始化
  const initFromStorage = () => {
    const storedToken = localStorage.getItem('token')
    const storedUserInfo = localStorage.getItem('userInfo')

    if (storedToken) {
      token.value = storedToken
    }

    if (storedUserInfo) {
      try {
        userInfo.value = JSON.parse(storedUserInfo)
      } catch {
        localStorage.removeItem('userInfo')
      }
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    setUserInfo,
    setToken,
    login,
    register,
    fetchCurrentUser,
    logout,
    initFromStorage,
  }
})
