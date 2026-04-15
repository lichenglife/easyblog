import { http } from './request'

// 评论信息
export interface CommentInfo {
  id: number
  postId: number
  userId: number
  parentId: number
  content: string
  isEdited: number
  likeCount: number
  replyCount: number
  status: number // 0-待审核 1-正常 2-已删除
  createAt: string
  updateAt: string
  replies?: CommentInfo[] // 楼中楼回复
}

// 创建评论请求
export interface CreateCommentRequest {
  postId: number
  parentId?: number
  content: string
}

// 更新评论请求
export interface UpdateCommentRequest {
  id: number
  content: string
}

// 评论列表响应
export interface CommentListResponse {
  totalCount: number
  hasMore: boolean
  comments: CommentInfo[]
}

// 评论 API
export const commentApi = {
  // 创建评论
  create(data: CreateCommentRequest): Promise<CommentInfo> {
    return http.post('/comment', data)
  },

  // 获取评论详情
  getById(id: number): Promise<CommentInfo> {
    return http.get(`/comment/${id}`)
  },

  // 获取文章评论列表（顶级评论）
  getByPostId(postId: number, params?: { page?: number; pageSize?: number }): Promise<CommentListResponse> {
    return http.get(`/comment/post/${postId}`, { params })
  },

  // 获取楼中楼回复列表
  getReplies(parentId: number, params?: { page?: number; pageSize?: number }): Promise<CommentListResponse> {
    return http.get(`/comment/reply/${parentId}`, { params })
  },

  // 更新评论
  update(id: number, data: UpdateCommentRequest): Promise<CommentInfo> {
    return http.put(`/comment/${id}`, data)
  },

  // 删除评论
  delete(id: number): Promise<void> {
    return http.delete(`/comment/${id}`)
  },

  // 点赞评论
  like(id: number): Promise<void> {
    return http.post(`/comment/${id}/like`)
  },

  // 取消点赞评论
  unlike(id: number): Promise<void> {
    return http.post(`/comment/${id}/unlike`)
  },
}
