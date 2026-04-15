<template>
  <div class="editor-page">
    <div class="editor-header">
      <h1>{{ isEditMode ? '编辑文章' : '创建新文章' }}</h1>
      <div class="header-actions">
        <button @click="handleSave" :disabled="saving" class="btn-save">
          {{ saving ? '保存中...' : '保存' }}
        </button>
        <button @click="handlePreview" class="btn-preview">预览</button>
        <button @click="handleCancel" class="btn-cancel">取消</button>
      </div>
    </div>

    <div class="editor-container">
      <div class="editor-main">
        <input
          v-model="form.title"
          type="text"
          placeholder="请输入文章标题..."
          class="title-input"
          maxlength="255"
        />

        <textarea
          v-model="form.content"
          placeholder="开始写作..."
          rows="20"
          class="content-textarea"
        ></textarea>
      </div>

      <div class="editor-sidebar">
        <div class="form-group">
          <label class="form-label">文章摘要</label>
          <textarea
            v-model="form.summary"
            placeholder="简要描述文章内容（可选）"
            rows="4"
            class="summary-textarea"
          ></textarea>
        </div>

        <div class="form-group">
          <label class="form-label">封面图片</label>
          <ImageUploader v-model="form.coverImage" />
        </div>

        <div class="form-group">
          <label class="form-label">分类 ID</label>
          <input
            v-model.number="form.categoryId"
            type="number"
            placeholder="输入分类 ID（可选）"
            class="input-field"
          />
        </div>

        <div class="form-group">
          <label class="form-label">文章状态</label>
          <select v-model="form.status" class="select-field">
            <option :value="0">草稿</option>
            <option :value="1">已发布</option>
            <option :value="2">已下架</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">是否置顶</label>
          <select v-model="form.isTop" class="select-field">
            <option :value="0">否</option>
            <option :value="1">是</option>
          </select>
        </div>
      </div>
    </div>

    <!-- 预览弹窗 -->
    <div v-if="showPreview" class="preview-modal" @click.self="showPreview = false">
      <div class="preview-content">
        <div class="preview-header">
          <h2>文章预览</h2>
          <button @click="showPreview = false" class="btn-close">×</button>
        </div>
        <div class="preview-body">
          <h1 class="preview-title">{{ form.title }}</h1>
          <div v-if="form.coverImage" class="preview-cover">
            <img :src="form.coverImage" :alt="form.title" />
          </div>
          <div class="preview-summary" v-if="form.summary">
            <h3>摘要</h3>
            <p>{{ form.summary }}</p>
          </div>
          <div class="preview-content-body">
            <div v-html="renderedContent"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { postApi, type CreatePostRequest } from '@/api/post'
import ImageUploader from '@/components/common/ImageUploader.vue'

const router = useRouter()
const route = useRoute()

const isEditMode = computed(() => {
  return !!route.params.id
})

const editingId = computed(() => {
  const id = route.params.id as string
  return id ? parseInt(id, 10) : null
})

const form = ref<CreatePostRequest>({
  title: '',
  content: '',
  summary: '',
  coverImage: '',
  categoryId: undefined,
  status: 0,
  isTop: 0,
})

const saving = ref(false)
const showPreview = ref(false)

const renderedContent = computed(() => {
  return form.value.content.replace(/\n/g, '<br>')
})

const validateForm = (): boolean => {
  if (!form.value.title.trim()) {
    alert('请输入文章标题')
    return false
  }
  if (!form.value.content.trim()) {
    alert('请输入文章内容')
    return false
  }
  return true
}

const handleSave = async () => {
  if (!validateForm()) return

  saving.value = true
  try {
    if (isEditMode.value && editingId.value) {
      await postApi.update(editingId.value, {
        id: editingId.value,
        ...form.value,
      })
      alert('文章已更新')
    } else {
      const result = await postApi.create(form.value)
      alert('文章已创建')
      router.push(`/articles/${result.postId}`)
    }
  } catch (error) {
    console.error('保存文章失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const handlePreview = () => {
  if (!form.value.title.trim() || !form.value.content.trim()) {
    alert('请输入标题和内容后再预览')
    return
  }
  showPreview.value = true
}

const handleCancel = () => {
  router.back()
}

const loadArticle = async () => {
  if (!isEditMode.value) return

  try {
    const post = await postApi.getById(editingId.value!)
    form.value = {
      title: post.title,
      content: post.content,
      summary: post.summary,
      coverImage: post.coverImage,
      categoryId: post.categoryId,
      status: post.status,
      isTop: post.isTop,
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    alert('文章不存在或加载失败')
    router.back()
  }
}

onMounted(() => {
  loadArticle()
})
</script>

<style scoped>
.editor-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem 1rem;
  background: #f5f5f5;
  min-height: 100vh;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  background: #fff;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.editor-header h1 {
  margin: 0;
  font-size: 1.25rem;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.btn-save,
.btn-preview,
.btn-cancel {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: opacity 0.2s;
}

.btn-save {
  background: #667eea;
  color: #fff;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-preview {
  background: #f0f0f0;
  color: #333;
}

.btn-cancel {
  background: #fff;
  color: #666;
  border: 1px solid #ddd;
}

.btn-save:hover,
.btn-preview:hover {
  opacity: 0.9;
}

.editor-container {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 1.5rem;
}

.editor-main {
  background: #fff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.title-input {
  width: 100%;
  padding: 0.75rem 0;
  border: none;
  border-bottom: 1px solid #eee;
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.title-input:focus {
  outline: none;
  border-bottom-color: #667eea;
}

.content-textarea {
  width: 100%;
  min-height: 400px;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 1rem;
  font-size: 1rem;
  line-height: 1.8;
  font-family: inherit;
  resize: vertical;
}

.content-textarea:focus {
  outline: none;
  border-color: #667eea;
}

.editor-sidebar {
  background: #fff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  height: fit-content;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: #666;
  font-weight: 500;
}

.input-field,
.select-field {
  width: 100%;
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

.summary-textarea {
  width: 100%;
  padding: 0.625rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
  resize: vertical;
}

.summary-textarea:focus {
  outline: none;
  border-color: #667eea;
}

.cover-preview {
  margin-top: 0.75rem;
  border-radius: 4px;
  overflow: hidden;
  max-height: 200px;
}

.cover-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 2rem 1rem;
  z-index: 1000;
  overflow-y: auto;
}

.preview-content {
  background: #fff;
  border-radius: 8px;
  width: 100%;
  max-width: 800px;
  margin: 2rem auto;
  overflow: hidden;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
}

.preview-header h2 {
  margin: 0;
  font-size: 1.125rem;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

.preview-body {
  padding: 2rem;
}

.preview-title {
  margin: 0 0 1rem;
  font-size: 1.5rem;
  color: #333;
}

.preview-cover {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 1.5rem;
  background: #f5f5f5;
}

.preview-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-summary h3 {
  margin: 0 0 0.5rem;
  font-size: 1rem;
  color: #666;
}

.preview-summary p {
  margin: 0;
  color: #666;
  line-height: 1.6;
}

.preview-content-body {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #eee;
  line-height: 1.8;
  color: #333;
}

@media (max-width: 900px) {
  .editor-container {
    grid-template-columns: 1fr;
  }

  .editor-sidebar {
    order: -1;
  }
}
</style>
