<template>
  <div class="tag-management">
    <div class="header">
      <h1>标签管理</h1>
      <button @click="showCreateModal = true" class="btn-create">新建标签</button>
    </div>

    <div class="table-container">
      <table class="tag-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>标签名称</th>
            <th>别名</th>
            <th>文章数</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tag in tags" :key="tag.id">
            <td>{{ tag.id }}</td>
            <td>
              <span class="tag-name">{{ tag.name }}</span>
            </td>
            <td>{{ tag.slug }}</td>
            <td>{{ tag.postCount }}</td>
            <td>
              <span :class="['status-badge', tag.status === 1 ? 'status-normal' : 'status-hidden']">
                {{ tag.status === 1 ? '显示' : '隐藏' }}
              </span>
            </td>
            <td>{{ formatDate(tag.createAt) }}</td>
            <td class="actions">
              <button @click="handleEdit(tag)" class="btn-edit">编辑</button>
              <button @click="handleDelete(tag)" class="btn-delete">删除</button>
            </td>
          </tr>
          <tr v-if="tags.length === 0">
            <td colspan="7" class="empty-message">暂无标签数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="totalPages > 1">
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

    <!-- 创建/编辑弹窗 -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ showEditModal ? '编辑标签' : '新建标签' }}</h2>
          <button @click="closeModal" class="btn-close">×</button>
        </div>
        <form @submit.prevent="handleSubmit" class="modal-body">
          <div class="form-group">
            <label class="form-label">标签名称</label>
            <input
              v-model="formData.name"
              type="text"
              placeholder="输入标签名称"
              class="input-field"
              required
              maxlength="50"
            />
          </div>
          <div class="form-group">
            <label class="form-label">标签别名</label>
            <input
              v-model="formData.slug"
              type="text"
              placeholder="用于 URL 的别名（可选）"
              class="input-field"
              maxlength="50"
            />
          </div>
          <div class="modal-footer">
            <button type="button" @click="closeModal" class="btn-cancel">取消</button>
            <button type="submit" :disabled="submitting" class="btn-submit">
              {{ submitting ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { tagApi, type TagInfo } from '@/api/post'

const tags = ref<TagInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 20
const totalCount = ref(0)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)

const formData = ref<{ id?: number; name: string; slug: string }>({
  name: '',
  slug: '',
})

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize))

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const fetchTags = async () => {
  loading.value = true
  try {
    const res = await tagApi.getList({
      page: currentPage.value,
      pageSize,
    })
    tags.value = res.tags
    totalCount.value = res.totalCount
  } catch (error) {
    console.error('获取标签列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleEdit = (tag: TagInfo) => {
  formData.value = {
    id: tag.id,
    name: tag.name,
    slug: tag.slug,
  }
  showEditModal.value = true
}

const handleDelete = async (tag: TagInfo) => {
  if (!confirm(`确定要删除标签 "${tag.name}" 吗？此操作不可恢复！`)) return

  try {
    await tagApi.delete(tag.id)
    await fetchTags()
    alert('标签已删除')
  } catch (error) {
    console.error('删除标签失败:', error)
    alert('删除失败，请重试')
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  formData.value = {
    name: '',
    slug: '',
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (showEditModal.value && formData.value.id) {
      await tagApi.update({
        id: formData.value.id,
        name: formData.value.name,
        slug: formData.value.slug,
      })
      alert('标签已更新')
    } else {
      await tagApi.create({
        name: formData.value.name,
        slug: formData.value.slug,
      })
      alert('标签已创建')
    }
    closeModal()
    await fetchTags()
  } catch (error) {
    console.error('保存标签失败:', error)
    alert('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchTags()
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
.tag-management {
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

.btn-create {
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

.table-container {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.tag-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.tag-table th,
.tag-table td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.tag-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  font-size: 0.875rem;
}

.tag-table td {
  font-size: 0.875rem;
  color: #666;
}

.tag-table tbody tr:hover {
  background: #f8f9fa;
}

.tag-name {
  font-weight: 500;
  color: #667eea;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-normal {
  background: #d1fae5;
  color: #059669;
}

.status-hidden {
  background: #f3f4f6;
  color: #6b7280;
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

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: #fff;
  border-radius: 8px;
  width: 100%;
  max-width: 450px;
  padding: 1.5rem;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.125rem;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #999;
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #666;
}

.input-field {
  padding: 0.625rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
}

.input-field:focus {
  outline: none;
  border-color: #667eea;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn-cancel,
.btn-submit {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-cancel {
  background: #f0f0f0;
  color: #666;
}

.btn-submit {
  background: #667eea;
  color: #fff;
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
