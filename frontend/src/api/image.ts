import { http } from './request'

// 图片上传响应
export interface UploadImageResponse {
  filename: string
  url: string
  size: number
}

// 图片信息
export interface ImageInfo {
  filename: string
  url: string
  size: number
  createdAt: string
}

// 图片列表响应
export interface ImageListResponse {
  total: number
  page: number
  size: number
  images: ImageInfo[]
}

// 图片 API
export const imageApi = {
  // 上传图片
  upload(file: File): Promise<UploadImageResponse> {
    const formData = new FormData()
    formData.append('file', file)
    return http.upload('/image/upload', formData)
  },

  // 获取图片列表
  getList(params?: { page?: number; pageSize?: number }): Promise<ImageListResponse> {
    return http.get('/images', { params })
  },

  // 删除图片
  delete(filename: string): Promise<void> {
    return http.delete(`/image/${filename}`)
  },
}
