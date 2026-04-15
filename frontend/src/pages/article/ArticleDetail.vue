<template>
  <div class="article-detail-page">
    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="!post" class="error">
      <p>文章不存在或已被删除</p>
      <button @click="$router.push('/articles')" class="btn-back">返回列表</button>
    </div>

    <article v-else class="article-container">
      <header class="article-header">
        <h1 class="article-title">{{ post.title }}</h1>

        <div class="article-meta">
          <div class="meta-left">
            <span class="author-name">作者：{{ authorName }}</span>
            <span class="publish-date">{{ formatDate(post.createAt) }}</span>
          </div>

          <div class="meta-right">
            <span class="meta-item">
              <i class="icon-view"></i>
              {{ post.viewCount }} 阅读
            </span>
            <span class="meta-item meta-item-clickable" @click="handleLike">
              <i class="icon-like" :class="{ active: isLiked }"></i>
              {{ likeCount }} 点赞
            </span>
            <span class="meta-item">
              <i class="icon-comment"></i>
              {{ post.commentCount }} 评论
            </span>
          </div>
        </div>

        <div v-if="isLoggedIn && isAuthor" class="article-actions">
          <button @click="handleEdit" class="btn-edit">编辑</button>
          <button @click="handleDelete" class="btn-delete">删除</button>
        </div>
      </header>

      <div v-if="post.coverImage" class="article-cover">
        <img :src="post.coverImage" :alt="post.title" />
      </div>

      <div class="article-content">
        <div v-html="renderedContent" class="content-body"></div>
      </div>

      <footer class="article-footer">
        <div class="tags" v-if="post.categoryId">
          <span class="tag">分类 ID: {{ post.categoryId }}</span>
        </div>
      </footer>

      <!-- 评论区 -->
      <CommentList :post-id="postId" />
    </article>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { postApi, type PostInfo } from '@/api/post'
import { likeApi } from '@/api/like'
import { useUserStore } from '@/stores/user'
import CommentList from '@/components/business/CommentList.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const post = ref<PostInfo | null>(null)
const loading = ref(false)
const isLiked = ref(false)
const likeCount = ref(0)

const postId = computed(() => {
  const id = route.params.id as string
  return parseInt(id, 10)
})

const isLoggedIn = computed(() => userStore.isLoggedIn)
const currentUserId = computed(() => userStore.userInfo?.userID)

const isAuthor = computed(() => {
  if (!post.value || !currentUserId.value) return false
  return post.value.userId.toString() === currentUserId.value.toString()
})

const authorName = computed(() => {
  return '作者'
})

const renderedContent = computed(() => {
  if (!post.value) return ''
  return post.value.content.replace(/\n/g, '<br>')
})

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const fetchPost = async () => {
  loading.value = true
  try {
    const data = await postApi.getById(postId.value)
    post.value = data
  } catch (error) {
    console.error('获取文章详情失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchLikeStatus = async () => {
  if (!isLoggedIn.value) return
  try {
    const status = await likeApi.getPostStatus(postId.value)
    isLiked.value = status.isLiked
    likeCount.value = status.count
  } catch (error) {
    console.error('获取点赞状态失败:', error)
  }
}

const handleEdit = () => {
  router.push(`/editor/${postId.value}`)
}

const handleDelete = async () => {
  if (!confirm('确定要删除这篇文章吗？')) return

  try {
    await postApi.delete(postId.value)
    alert('文章已删除')
    router.push('/articles')
  } catch (error) {
    console.error('删除文章失败:', error)
    alert('删除失败，请重试')
  }
}

const handleLike = async () => {
  if (!isLoggedIn.value) {
    router.push('/login')
    return
  }

  try {
    if (isLiked.value) {
      await likeApi.unlikePost(postId.value)
      isLiked.value = false
      likeCount.value = Math.max(0, likeCount.value - 1)
    } else {
      await likeApi.likePost(postId.value)
      isLiked.value = true
      likeCount.value++
    }
  } catch (error) {
    console.error('点赞失败:', error)
  }
}

onMounted(() => {
  fetchPost()
  fetchLikeStatus()
})
</script>

<style scoped>
.article-detail-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.loading,
.error {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.btn-back {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-back:hover {
  opacity: 0.9;
}

.article-container {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.article-header {
  padding: 2rem 2rem 1rem;
  border-bottom: 1px solid #eee;
}

.article-title {
  margin: 0 0 1rem;
  font-size: 1.75rem;
  color: #333;
  line-height: 1.3;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.meta-left {
  display: flex;
  gap: 1rem;
  color: #666;
  font-size: 0.875rem;
}

.meta-right {
  display: flex;
  gap: 1rem;
  color: #999;
  font-size: 0.875rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.meta-item-clickable {
  cursor: pointer;
  transition: color 0.2s;
}

.meta-item-clickable:hover {
  color: #667eea;
}

.icon-like.active {
  color: #f56c6c;
  font-weight: bold;
}

.article-actions {
  margin-top: 1rem;
  display: flex;
  gap: 0.5rem;
}

.btn-edit,
.btn-delete {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-edit {
  background: #667eea;
  color: #fff;
}

.btn-delete {
  background: #f56c6c;
  color: #fff;
}

.article-cover {
  width: 100%;
  height: 300px;
  overflow: hidden;
  background: #f5f5f5;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  padding: 2rem;
}

.content-body {
  font-size: 1rem;
  line-height: 1.8;
  color: #333;
}

.content-body :deep(p) {
  margin-bottom: 1rem;
}

.article-footer {
  padding: 1rem 2rem;
  border-top: 1px solid #eee;
}

.tags {
  display: flex;
  gap: 0.5rem;
}

.tag {
  padding: 0.25rem 0.75rem;
  background: #f0f0f0;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #666;
}

.comments-section {
  margin-top: 2rem;
  padding: 2rem;
  background: #fafafa;
  border-radius: 8px;
}

.comments-title {
  margin: 0 0 1.5rem;
  font-size: 1.25rem;
  color: #333;
}

.comment-form {
  margin-bottom: 1.5rem;
}

.comment-input {
  width: 100%;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
  resize: vertical;
}

.comment-input:focus {
  outline: none;
  border-color: #667eea;
}

.btn-submit {
  margin-top: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.login-hint {
  color: #999;
  font-size: 0.875rem;
  margin-bottom: 1.5rem;
}

.login-hint a {
  color: #667eea;
}

.no-comments {
  text-align: center;
  color: #999;
  padding: 2rem;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-item {
  padding: 1rem;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #eee;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.comment-author {
  font-weight: 600;
  color: #333;
  font-size: 0.875rem;
}

.comment-date {
  color: #999;
  font-size: 0.75rem;
}

.comment-content {
  color: #333;
  font-size: 0.875rem;
  line-height: 1.6;
  margin-bottom: 0.75rem;
}

.edited-tag {
  color: #999;
  font-size: 0.75rem;
  font-style: italic;
}

.comment-actions {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
}

.comment-action {
  color: #666;
  cursor: pointer;
  transition: color 0.2s;
}

.comment-action:hover {
  color: #667eea;
}

.comment-action.delete {
  color: #999;
}

.comment-action.delete:hover {
  color: #f56c6c;
}

@media (max-width: 640px) {
  .article-title {
    font-size: 1.25rem;
  }

  .article-meta {
    flex-direction: column;
    align-items: flex-start;
  }

  .article-cover {
    height: 200px;
  }
}
</style>
