<template>
  <div class="category-management">
    <div class="header">
      <h1>分类管理</h1>
      <button @click="showCreateModal = true" class="btn-create">新建分类</button>
    </div>

    <div class="tree-container">
      <div class="category-tree">
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="categories.length === 0" class="empty">暂无分类数据</div>
        <div v-else class="tree-list">
          <div v-for="category in categories" :key="category.id" class="category-node">
            <div class="node-header">
              <span class="node-name">{{ category.name }}</span>
              <span class="node-slug">({{ category.slug }})</span>
              <div class="node-actions">
                <button @click="handleEdit(category)" class="btn-edit">编辑</button>
                <button @click="handleDelete(category)" class="btn-delete">删除</button>
              </div>
            </div>
            <!-- 子分类 -->
            <div v-if="category.children && category.children.length > 0" class="node-children">
              <div
                v-for="child in category.children"
                :key="child.id"
                class="category-node child-node"
              >
                <div class="node-header">
                  <span class="node-name">{{ child.name }}</span>
                  <span class="node-slug">({{ child.slug }})</span>
                  <div class="node-actions">
                    <button @click="handleEdit(child)" class="btn-edit">编辑</button>
                    <button @click="handleDelete(child)" class="btn-delete">删除</button>
                  </div>
                </div>
                <!-- 三级分类 -->
                <div v-if="child.children && child.children.length > 0" class="node-children">
                  <div
                    v-for="grandchild in child.children"
                    :key="grandchild.id"
                    class="category-node grandchild-node"
                  >
                    <div class="node-header">
                      <span class="node-name">{{ grandchild.name }}</span>
                      <span class="node-slug">({{ grandchild.slug }})</span>
                      <div class="node-actions">
                        <button @click="handleEdit(grandchild)" class="btn-edit">编辑</button>
                        <button @click="handleDelete(grandchild)" class="btn-delete">删除</button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑弹窗 -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ showEditModal ? '编辑分类' : '新建分类' }}</h2>
          <button @click="closeModal" class="btn-close">×</button>
        </div>
        <form @submit.prevent="handleSubmit" class="modal-body">
          <div class="form-group">
            <label class="form-label">分类名称</label>
            <input
              v-model="formData.name"
              type="text"
              placeholder="输入分类名称"
              class="input-field"
              required
              maxlength="50"
            />
          </div>
          <div class="form-group">
            <label class="form-label">分类别名</label>
            <input
              v-model="formData.slug"
              type="text"
              placeholder="用于 URL 的别名（可选）"
              class="input-field"
              maxlength="50"
            />
          </div>
          <div class="form-group">
            <label class="form-label">父级分类</label>
            <select v-model="formData.parentId" class="select-field">
              <option :value="0">无（顶级分类）</option>
              <option
                v-for="cat in flatCategories"
                :key="cat.id"
                :value="cat.id"
                :disabled="cat.level >= 2"
              >
                {{ ' '.repeat(cat.level * 2) }}{{ cat.name }}
              </option>
            </select>
            <p class="form-hint">最多支持 3 级分类</p>
          </div>
          <div class="form-group">
            <label class="form-label">排序</label>
            <input
              v-model.number="formData.sort"
              type="number"
              placeholder="数字越小越靠前"
              class="input-field"
              min="0"
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
import { ref, onMounted } from 'vue'
import { categoryApi, type CategoryInfo, type CreateCategoryRequest } from '@/api/post'

const categories = ref<CategoryInfo[]>([])
const flatCategories = ref<CategoryInfo[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)

const formData = ref<CreateCategoryRequest & { id?: number }>({
  name: '',
  slug: '',
  parentId: 0,
  sort: 0,
})

const flattenCategories = (cats: CategoryInfo[], level = 0): CategoryInfo[] => {
  const result: CategoryInfo[] = []
  for (const cat of cats) {
    result.push({ ...cat, level })
    if (cat.children && cat.children.length > 0) {
      result.push(...flattenCategories(cat.children, level + 1))
    }
  }
  return result
}

const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await categoryApi.getTree()
    categories.value = res
    flatCategories.value = flattenCategories(res)
  } catch (error) {
    console.error('获取分类树失败:', error)
  } finally {
    loading.value = false
  }
}

const handleEdit = (category: CategoryInfo) => {
  formData.value = {
    id: category.id,
    name: category.name,
    slug: category.slug,
    parentId: category.parentId,
    sort: category.sort,
  }
  showEditModal.value = true
}

const handleDelete = async (category: CategoryInfo) => {
  if (!confirm(`确定要删除分类 "${category.name}" 吗？此操作不可恢复！`)) return

  try {
    await categoryApi.delete(category.id)
    await fetchCategories()
    alert('分类已删除')
  } catch (error) {
    console.error('删除分类失败:', error)
    alert('删除失败，请重试')
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  formData.value = {
    name: '',
    slug: '',
    parentId: 0,
    sort: 0,
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (showEditModal.value && formData.value.id) {
      await categoryApi.update({
        id: formData.value.id,
        name: formData.value.name,
        slug: formData.value.slug,
        parentId: formData.value.parentId,
        sort: formData.value.sort,
      })
      alert('分类已更新')
    } else {
      await categoryApi.create({
        name: formData.value.name,
        slug: formData.value.slug,
        parentId: formData.value.parentId,
        sort: formData.value.sort,
      })
      alert('分类已创建')
    }
    closeModal()
    await fetchCategories()
  } catch (error) {
    console.error('保存分类失败:', error)
    alert('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-management {
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

.tree-container {
  background: #fff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.loading,
.empty {
  text-align: center;
  color: #999;
  padding: 2rem;
}

.category-node {
  border-bottom: 1px solid #f0f0f0;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 0;
}

.node-name {
  font-weight: 600;
  color: #333;
  font-size: 1rem;
}

.node-slug {
  color: #999;
  font-size: 0.875rem;
}

.node-actions {
  margin-left: auto;
  display: flex;
  gap: 0.5rem;
}

.btn-edit,
.btn-delete {
  padding: 0.25rem 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
}

.btn-edit {
  background: #667eea;
  color: #fff;
}

.btn-delete {
  background: #f56c6c;
  color: #fff;
}

.node-children {
  padding-left: 1.5rem;
  border-left: 2px solid #f0f0f0;
}

.child-node > .node-header {
  background: #fafafa;
}

.grandchild-node > .node-header {
  background: #f5f5f5;
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
  max-width: 500px;
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

.form-hint {
  font-size: 0.75rem;
  color: #999;
  margin: 0.25rem 0 0;
}

.input-field,
.select-field {
  padding: 0.625rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
}

.input-field:focus,
.select-field:focus {
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
