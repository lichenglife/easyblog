<template>
  <div class="article-list-page">
    <div class="header">
      <h1>文章列表</h1>
      <button v-if="isLoggedIn" @click="$router.push('/editor')" class="btn-create">
        写文章
      </button>
    </div>

    <div class="filters">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索文章标题..."
        class="search-input"
      />
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="posts.length === 0" class="empty">
      暂无文章
    </div>

    <div v-else class="posts-container">
      <article
        v-for="post in filteredPosts"
        :key="post.id"
        class="post-card"
        @click="goToDetail(post.id)"
      >
        <div class="post-cover" v-if="post.coverImage">
          <img :src="post.coverImage" :alt="post.title" />
        </div>
        <div class="post-content">
          <h2 class="post-title">{{ post.title }}</h2>
          <p class="post-summary">{{ post.summary || post.content.slice(0, 150) }}...</p>
          <div class="post-meta">
            <span class="post-date">{{ formatDate(post.createAt) }}</span>
            <span class="post-views">阅读 {{ post.viewCount }}</span>
            <span class="post-likes">点赞 {{ post.likeCount }}</span>
            <span class="post-comments">评论 {{ post.commentCount }}</span>
          </div>
        </div>
      </article>
    </div>

    <div class="pagination">
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
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const posts = ref<PostInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 10
const searchQuery = ref('')

const isLoggedIn = computed(() => userStore.isLoggedIn)

const filteredPosts = computed(() => {
  if (!searchQuery.value) return posts.value
  const query = searchQuery.value.toLowerCase()
  return posts.value.filter(
    (post) => post.title.toLowerCase().includes(query)
  )
})

const totalPages = computed(() => Math.ceil(posts.value.length / pageSize))

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const fetchPosts = async () => {
  loading.value = true
  try {
    const res = await postApi.getList({
      page: currentPage.value,
      pageSize,
    })
    posts.value = res.posts
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: number) => {
  router.push(`/articles/${id}`)
}

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchPosts()
}

onMounted(() => {
  fetchPosts()
})
</script>

<style scoped>
.article-list-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1rem;
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

.btn-create {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-create:hover {
  opacity: 0.9;
}

.filters {
  margin-bottom: 1.5rem;
}

.search-input {
  width: 100%;
  max-width: 400px;
  padding: 0.75rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.loading,
.empty {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.posts-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.post-card {
  display: flex;
  gap: 1.5rem;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.post-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.post-cover {
  width: 200px;
  flex-shrink: 0;
  background: #f5f5f5;
}

.post-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-content {
  flex: 1;
  padding: 1.25rem;
}

.post-title {
  margin: 0 0 0.75rem;
  font-size: 1.125rem;
  color: #333;
}

.post-summary {
  margin: 0 0 1rem;
  color: #666;
  font-size: 0.875rem;
  line-height: 1.6;
}

.post-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
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

@media (max-width: 640px) {
  .post-card {
    flex-direction: column;
  }

  .post-cover {
    width: 100%;
    height: 160px;
  }
}
</style>
