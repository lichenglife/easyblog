<template>
  <div class="comment-component">
    <!-- 发表评论表单 -->
    <div v-if="isLoggedIn" class="comment-form">
      <textarea
        v-model="commentContent"
        placeholder="写下你的评论..."
        rows="4"
        class="comment-input"
      ></textarea>
      <div class="form-actions">
        <button
          v-if="isReplyMode"
          @click="cancelReply"
          class="btn-cancel"
        >
          取消
        </button>
        <button
          @click="submitComment"
          :disabled="!commentContent.trim() || submitting"
          class="btn-submit"
        >
          {{ isReplyMode ? '回复' : '发表评论' }}
        </button>
      </div>
    </div>

    <div v-else class="login-hint">
      <router-link to="/login">登录</router-link> 后发表评论
    </div>

    <!-- 评论列表 -->
    <div class="comments-list">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="comments.length === 0" class="no-comments">暂无评论</div>
      <div v-else class="comment-list">
        <div
          v-for="comment in comments"
          :key="comment.id"
          class="comment-item"
        >
          <div class="comment-header">
            <span class="comment-author">用户 #{{ comment.userId }}</span>
            <span class="comment-date">{{ formatDate(comment.createAt) }}</span>
          </div>
          <div class="comment-content">
            {{ comment.content }}
            <span v-if="comment.isEdited" class="edited-tag">(已编辑)</span>
          </div>
          <div class="comment-meta">
            <span class="meta-item meta-item-clickable" @click="likeComment(comment)">
              👍 {{ comment.likeCount }}
            </span>
            <span
              v-if="isLoggedIn"
              class="meta-item meta-item-clickable"
              @click="startReply(comment)"
            >
              回复
            </span>
            <span
              v-if="isLoggedIn && canDelete(comment)"
              class="meta-item meta-item-clickable delete"
              @click="deleteComment(comment.id)"
            >
              删除
            </span>
          </div>

          <!-- 回复列表 -->
          <div v-if="comment.replyCount > 0" class="replies">
            <div
              v-for="reply in comment.replies"
              :key="reply.id"
              class="reply-item"
            >
              <div class="reply-header">
                <span class="reply-author">用户 #{{ reply.userId }}</span>
                <span class="reply-date">{{ formatDate(reply.createAt) }}</span>
              </div>
              <div class="reply-content">
                {{ reply.content }}
                <span v-if="reply.isEdited" class="edited-tag">(已编辑)</span>
              </div>
              <div class="reply-meta">
                <span class="meta-item meta-item-clickable" @click="likeComment(reply)">
                  👍 {{ reply.likeCount }}
                </span>
                <span
                  v-if="isLoggedIn && canDelete(reply)"
                  class="meta-item meta-item-clickable delete"
                  @click="deleteComment(reply.id)"
                >
                  删除
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
        class="btn-page"
      >
        上一页
      </button>
      <span class="page-info">第 {{ currentPage }} 页 / 共 {{ totalPages }} 页</span>
      <button
        :disabled="currentPage >= totalPages"
        @click="handlePageChange(currentPage + 1)"
        class="btn-page"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { commentApi, type CommentInfo } from '@/api/comment'
import { useUserStore } from '@/stores/user'

interface Props {
  postId: number
}

const props = defineProps<Props>()

const router = useRouter()
const userStore = useUserStore()

const comments = ref<CommentInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 20
const totalPages = ref(1)
const commentContent = ref('')
const submitting = ref(false)
const isReplyMode = ref(false)
const parentCommentId = ref<number | null>(null)

const isLoggedIn = computed(() => userStore.isLoggedIn)
const currentUserId = computed(() => userStore.userInfo?.userID)

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const canDelete = (comment: CommentInfo): boolean => {
  if (!currentUserId.value) return false
  return currentUserId.value.toString() === comment.userId.toString()
}

const fetchComments = async () => {
  loading.value = true
  try {
    const res = await commentApi.getByPostId(props.postId, {
      page: currentPage.value,
      pageSize,
    })
    comments.value = res.comments
    totalPages.value = Math.ceil(res.totalCount / pageSize)
  } catch (error) {
    console.error('获取评论列表失败:', error)
  } finally {
    loading.value = false
  }
}

const submitComment = async () => {
  if (!commentContent.value.trim()) return

  submitting.value = true
  try {
    await commentApi.create({
      postId: props.postId,
      parentId: isReplyMode.value ? (parentCommentId.value || 0) : 0,
      content: commentContent.value,
    })
    commentContent.value = ''
    isReplyMode.value = false
    parentCommentId.value = null
    await fetchComments()
  } catch (error) {
    console.error('发表评论失败:', error)
    alert('评论失败，请重试')
  } finally {
    submitting.value = false
  }
}

const startReply = (comment: CommentInfo) => {
  isReplyMode.value = true
  parentCommentId.value = comment.id
  commentContent.value = `回复 @用户 #${comment.userId}: `
}

const cancelReply = () => {
  isReplyMode.value = false
  parentCommentId.value = null
  commentContent.value = ''
}

const deleteComment = async (commentId: number) => {
  if (!confirm('确定要删除这条评论吗？')) return

  try {
    await commentApi.delete(commentId)
    await fetchComments()
  } catch (error) {
    console.error('删除评论失败:', error)
    alert('删除失败，请重试')
  }
}

const likeComment = async (comment: CommentInfo) => {
  if (!isLoggedIn.value) {
    router.push('/login')
    return
  }

  try {
    await commentApi.like(comment.id)
    comment.likeCount++
  } catch (error) {
    console.error('点赞评论失败:', error)
  }
}

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchComments()
}

watch(() => props.postId, () => {
  currentPage.value = 1
  fetchComments()
})

onMounted(() => {
  fetchComments()
})
</script>

<style scoped>
.comment-component {
  margin-top: 2rem;
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

.form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn-submit,
.btn-cancel {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-submit {
  background: #667eea;
  color: #fff;
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-cancel {
  background: #f0f0f0;
  color: #666;
}

.login-hint {
  color: #999;
  font-size: 0.875rem;
  margin-bottom: 1.5rem;
}

.login-hint a {
  color: #667eea;
}

.comments-list {
  background: #fafafa;
  border-radius: 8px;
  padding: 1.5rem;
}

.loading,
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
  font-size: 0.875rem;
}

.comment-author {
  font-weight: 600;
  color: #333;
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

.comment-meta,
.reply-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: #666;
}

.meta-item {
  cursor: pointer;
  transition: color 0.2s;
}

.meta-item:hover {
  color: #667eea;
}

.meta-item.delete:hover {
  color: #f56c6c;
}

.replies {
  margin-top: 1rem;
  padding-left: 1rem;
  border-left: 2px solid #f0f0f0;
}

.reply-item {
  padding: 0.75rem;
  background: #f9f9f9;
  border-radius: 4px;
  margin-top: 0.5rem;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  font-size: 0.75rem;
}

.reply-author {
  font-weight: 600;
  color: #333;
}

.reply-date {
  color: #999;
  font-size: 0.7rem;
}

.reply-content {
  color: #333;
  font-size: 0.875rem;
  line-height: 1.6;
  margin-bottom: 0.5rem;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-page {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-page:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: #666;
  font-size: 0.875rem;
}
</style>
