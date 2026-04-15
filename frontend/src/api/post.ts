import { http } from './request'

// 文章信息
export interface PostInfo {
  id: number
  postId: string
  userId: number
  title: string
  content: string
  summary: string
  coverImage: string
  categoryId: number
  status: number // 0-草稿 1-已发布 2-已下架
  isTop: number
  viewCount: number
  likeCount: number
  commentCount: number
  createAt: string
  updateAt: string
}

// 创建文章请求
export interface CreatePostRequest {
  title: string
  content: string
  summary?: string
  coverImage?: string
  categoryId?: number
  status?: number
  isTop?: number
}

// 更新文章请求
export interface UpdatePostRequest {
  id: number
  title: string
  content: string
  summary?: string
  coverImage?: string
  categoryId?: number
  status?: number
  isTop?: number
}

// 文章列表响应
export interface PostListResponse {
  totalCount: number
  hasMore: boolean
  posts: PostInfo[]
}

// 文章 API
export const postApi = {
  // 创建文章
  create(data: CreatePostRequest): Promise<{ postId: string }> {
    return http.post('/post', data)
  },

  // 获取文章详情
  getById(id: number): Promise<PostInfo> {
    return http.get(`/post/${id}`)
  },

  // 获取文章列表
  getList(params?: { page?: number; pageSize?: number }): Promise<PostListResponse> {
    return http.get('/post/list', { params })
  },

  // 获取用户的文章列表
  getByUserId(userId: number, params?: { page?: number; pageSize?: number }): Promise<PostListResponse> {
    return http.get(`/post/user/${userId}`, { params })
  },

  // 更新文章
  update(id: number, data: UpdatePostRequest): Promise<void> {
    return http.put(`/post/${id}`, data)
  },

  // 删除文章
  delete(id: number): Promise<void> {
    return http.delete(`/post/${id}`)
  },
}

// 分类信息
export interface CategoryInfo {
  id: number
  name: string
  slug: string
  parentId: number
  level: number
  sort: number
  postCount: number
  status: number
  createAt: string
  updateAt: string
  children?: CategoryInfo[]
}

// 创建分类请求
export interface CreateCategoryRequest {
  name: string
  slug: string
  parentId?: number
  sort?: number
}

// 更新分类请求
export interface UpdateCategoryRequest {
  id: number
  name: string
  slug: string
  parentId?: number
  sort?: number
  status?: number
}

// 分类 API
export const categoryApi = {
  // 创建分类
  create(data: CreateCategoryRequest): Promise<{ categoryId: number }> {
    return http.post('/category', data)
  },

  // 获取分类详情
  getById(id: number): Promise<CategoryInfo> {
    return http.get(`/category/${id}`)
  },

  // 获取分类树
  getTree(): Promise<CategoryInfo[]> {
    return http.get('/categories/tree')
  },

  // 获取分类列表
  getList(params?: { page?: number; pageSize?: number }): Promise<{
    totalCount: number
    hasMore: boolean
    categories: CategoryInfo[]
  }> {
    return http.get('/categories', { params })
  },

  // 更新分类
  update(data: UpdateCategoryRequest): Promise<{ categoryId: number }> {
    return http.put('/category', data)
  },

  // 删除分类
  delete(id: number): Promise<void> {
    return http.delete(`/category/${id}`)
  },
}

// 标签信息
export interface TagInfo {
  id: number
  name: string
  slug: string
  postCount: number
  status: number
  createAt: string
  updateAt: string
}

// 创建标签请求
export interface CreateTagRequest {
  name: string
  slug: string
}

// 更新标签请求
export interface UpdateTagRequest {
  id: number
  name: string
  slug: string
  status?: number
}

// 标签 API
export const tagApi = {
  // 创建标签
  create(data: CreateTagRequest): Promise<{ tagId: number }> {
    return http.post('/tag', data)
  },

  // 获取标签详情
  getById(id: number): Promise<TagInfo> {
    return http.get(`/tag/${id}`)
  },

  // 获取标签列表
  getList(params?: { page?: number; pageSize?: number }): Promise<{
    totalCount: number
    hasMore: boolean
    tags: TagInfo[]
  }> {
    return http.get('/tags', { params })
  },

  // 更新标签
  update(data: UpdateTagRequest): Promise<{ tagId: number }> {
    return http.put('/tag', data)
  },

  // 删除标签
  delete(id: number): Promise<void> {
    return http.delete(`/tag/${id}`)
  },
}
