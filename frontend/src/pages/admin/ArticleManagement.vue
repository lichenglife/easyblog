<template>
  <div class="article-management">
    <div class="header">
      <h1>文章管理</h1>
      <div class="actions">
        <button @click="handleCreate" class="btn-create">新建文章</button>
      </div>
    </div>

    <div class="filters">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索文章标题..."
        class="search-input"
        @input="handleSearch"
      />
      <select v-model="statusFilter" class="select-filter" @change="fetchArticles">
        <option :value="-1">全部状态</option>
        <option :value="0">草稿</option>
        <option :value="1">已发布</option>
        <option :value="2">已下架</option>
      </select>
    </div>

    <div class="table-container">
      <table class="article-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>标题</th>
            <th>作者</th>
            <th>分类</th>
            <th>状态</th>
            <th>阅读量</th>
            <th>点赞数</th>
            <th>评论数</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in articles" :key="article.id">
            <td>{{ article.id }}</td>
            <td class="title-cell">
              <span class="article-title">{{ article.title }}</span>
              <span v-if="article.isTop === 1" class="top-tag">置顶</span>
            </td>
            <td>{{ article.userId }}</td>
            <td>{{ article.categoryId || '-' }}</td>
            <td>
              <span :class="['status-badge', getStatusClass(article.status)]">
                {{ getStatusText(article.status) }}
              </span>
            </td>
            <td>{{ article.viewCount }}</td>
            <td>{{ article.likeCount }}</td>
            <td>{{ article.commentCount }}</td>
            <td>{{ formatDate(article.createAt) }}</td>
            <td class="actions">
              <button @click="handleEdit(article)" class="btn-edit">编辑</button>
              <button @click="handleView(article)" class="btn-view">查看</button>
              <button
                v-if="article.status === 1"
                @click="handleToggleTop(article)"
                class="btn-top"
              >
                {{ article.isTop === 1 ? '取消置顶' : '置顶' }}
              </button>
              <button @click="handleDelete(article)" class="btn-delete">删除</button>
            </td>
          </tr>
          <tr v-if="articles.length === 0">
            <td colspan="10" class="empty-message">暂无文章数据</td>
          </tr>
        </tbody>
      </table>
    </div>

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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { postApi, type PostInfo } from '@/api/post'

const router = useRouter()

const articles = ref<PostInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 10
const totalCount = ref(0)
const searchQuery = ref('')
const statusFilter = ref(-1)

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize))

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const getStatusClass = (status: number): string => {
  switch (status) {
    case 0:
      return 'status-draft'
    case 1:
      return 'status-published'
    case 2:
      return 'status-offline'
    default:
      return ''
  }
}

const getStatusText = (status: number): string => {
  switch (status) {
    case 0:
      return '草稿'
    case 1:
      return '已发布'
    case 2:
      return '已下架'
    default:
      return '未知'
  }
}

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await postApi.getList({
      page: currentPage.value,
      pageSize,
    })
    articles.value = res.posts
    totalCount.value = res.totalCount
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleCreate = () => {
  router.push('/editor')
}

const handleEdit = (article: PostInfo) => {
  router.push(`/editor/${article.id}`)
}

const handleView = (article: PostInfo) => {
  router.push(`/articles/${article.id}`)
}

const handleToggleTop = async (article: PostInfo) => {
  const action = article.isTop === 1 ? '取消置顶' : '置顶'
  if (!confirm(`确定要${action}这篇文章吗？`)) return

  try {
    // 调用更新接口
    await postApi.update(article.id, {
      id: article.id,
      title: article.title,
      content: article.content,
      summary: article.summary,
      coverImage: article.coverImage,
      categoryId: article.categoryId,
      status: article.status,
      isTop: article.isTop === 1 ? 0 : 1,
    })
    await fetchArticles()
    alert(`${action}成功`)
  } catch (error) {
    console.error('操作失败:', error)
    alert('操作失败，请重试')
  }
}

const handleDelete = async (article: PostInfo) => {
  if (!confirm(`确定要删除文章 "${article.title}" 吗？此操作不可恢复！`)) return

  try {
    await postApi.delete(article.id)
    await fetchArticles()
    alert('文章已删除')
  } catch (error) {
    console.error('删除文章失败:', error)
    alert('删除失败，请重试')
  }
}

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchArticles()
}

onMounted(() => {
  fetchArticles()
})
</script>

<style scoped>
.article-management {
  padding: 1.5rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.header h1 {
  margin: 0;
  font-size: 1.5rem;
  color: #333;
}

.actions .btn-create {
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-create:hover {
  opacity: 0.9;
}

.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.search-input,
.select-filter {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
}

.search-input {
  width: 250px;
}

.search-input:focus,
.select-filter:focus {
  outline: none;
  border-color: #667eea;
}

.table-container {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.article-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.article-table th,
.article-table td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.article-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  font-size: 0.875rem;
}

.article-table td {
  font-size: 0.875rem;
  color: #666;
}

.article-table tbody tr:hover {
  background: #f8f9fa;
}

.title-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.article-title {
  font-weight: 500;
  color: #333;
}

.top-tag {
  padding: 0.125rem 0.5rem;
  background: #ef4444;
  color: #fff;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 500;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-draft {
  background: #f3f4f6;
  color: #6b7280;
}

.status-published {
  background: #d1fae5;
  color: #059669;
}

.status-offline {
  background: #fee2e2;
  color: #dc2626;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.actions button {
  padding: 0.25rem 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.actions button:hover {
  opacity: 0.8;
}

.btn-edit {
  background: #667eea;
  color: #fff;
}

.btn-view {
  background: #3b82f6;
  color: #fff;
}

.btn-top {
  background: #f59e0b;
  color: #fff;
}

.btn-delete {
  background: #f56c6c;
  color: #fff;
}

.empty-message {
  text-align: center;
  padding: 2rem;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
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

.btn-page:not(:disabled):hover {
  opacity: 0.9;
}

.page-info {
  color: #666;
  font-size: 0.875rem;
}
</style>
