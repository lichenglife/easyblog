import { http } from './request'

// 点赞信息
export interface LikeInfo {
  id: number
  userId: number
  targetId: number
  targetType: number // 1-文章 2-评论
  createAt: string
}

// 点赞状态响应
export interface LikeStatusResponse {
  isLiked: boolean
  count: number
}

// 点赞 API
export const likeApi = {
  // 点赞文章
  likePost(postId: number): Promise<void> {
    return http.post(`/like/post/${postId}`)
  },

  // 取消点赞文章
  unlikePost(postId: number): Promise<void> {
    return http.post(`/like/post/${postId}/unlike`)
  },

  // 获取文章点赞状态
  getPostStatus(postId: number): Promise<LikeStatusResponse> {
    return http.get(`/like/post/${postId}/status`)
  },
}
