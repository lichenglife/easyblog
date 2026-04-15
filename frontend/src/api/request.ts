import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { useUserStore } from '@/stores/user'

// API 基础 URL
const API_BASE_URL = 'http://localhost:8080/v1'

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    // 后端返回格式：{ code: 200, message: '', data: {} }
    if (res.code === 200 || res.code === 0) {
      return res.data
    }
    // 处理业务错误
    if (res.code === 401 || res.code === 10002) {
      // 未授权，跳转登录
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    // 处理 HTTP 错误
    if (error.response) {
      switch (error.response.status) {
        case 401:
          console.error('未授权，请重新登录')
          break
        case 403:
          console.error('拒绝访问')
          break
        case 404:
          console.error('请求资源不存在')
          break
        case 500:
          console.error('服务器内部错误')
          break
        default:
          console.error(`请求失败：${error.response.status}`)
      }
    }
    return Promise.reject(error)
  }
)

// 导出请求方法
export const http = {
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return request.get(url, config).then((res) => res as T)
  },

  post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return request.post(url, data, config).then((res) => res as T)
  },

  put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return request.put(url, data, config).then((res) => res as T)
  },

  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return request.delete(url, config).then((res) => res as T)
  },

  upload<T = unknown>(url: string, formData: FormData): Promise<T> {
    return request.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }).then((res) => res as T)
  },
}

export default request
