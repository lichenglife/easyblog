<template>
  <div class="image-uploader">
    <div
      class="upload-area"
      :class="{ 'is-dragover': isDragover }"
      @dragover.prevent="isDragover = true"
      @dragleave.prevent="isDragover = false"
      @drop.prevent="handleDrop"
      @click="handleClick"
    >
      <input
        ref="inputRef"
        type="file"
        accept="image/*"
        class="hidden-input"
        @change="handleFileSelect"
      />

      <div v-if="loading" class="uploading">
        <div class="spinner"></div>
        <p>上传中...</p>
      </div>

      <div v-else-if="imageUrl" class="preview">
        <img :src="imageUrl" alt="预览图" />
        <div class="preview-actions">
          <button type="button" @click.stop="handleRemove" class="btn-remove">
            删除
          </button>
        </div>
      </div>

      <div v-else class="upload-placeholder">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <polyline points="17 8 12 3 7 8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <line x1="12" y1="3" x2="12" y2="15" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <p>点击或拖拽上传图片</p>
        <p class="hint">支持 JPG、PNG、GIF、WebP 格式，最大 {{ maxSize }}MB</p>
      </div>
    </div>

    <p v-if="error" class="error-msg">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { imageApi } from '@/api/image'

interface Props {
  modelValue?: string
  maxSize?: number // MB
  limit?: number
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  maxSize: 5,
  limit: 1,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const isDragover = ref(false)
const loading = ref(false)
const error = ref('')
const imageUrl = ref(props.modelValue)

const maxSizeBytes = computed(() => props.maxSize * 1024 * 1024)

watch(() => props.modelValue, (val) => {
  if (val) {
    imageUrl.value = val
  }
})

const handleClick = () => {
  inputRef.value?.click()
}

const handleDrop = (e: DragEvent) => {
  isDragover.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    uploadFile(files[0])
  }
}

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    uploadFile(files[0])
  }
}

const uploadFile = async (file: File) => {
  // 检查文件类型
  const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']
  if (!validTypes.includes(file.type)) {
    error.value = '只支持图片文件（JPG、PNG、GIF、WebP、SVG）'
    return
  }

  // 检查文件大小
  if (file.size > maxSizeBytes.value) {
    error.value = `图片大小不能超过 ${props.maxSize}MB`
    return
  }

  error.value = ''
  loading.value = true

  try {
    const res = await imageApi.upload(file)
    imageUrl.value = res.url
    emit('update:modelValue', res.url)
  } catch (err) {
    console.error('上传图片失败:', err)
    error.value = '上传失败，请重试'
  } finally {
    loading.value = false
    if (inputRef.value) {
      inputRef.value.value = ''
    }
  }
}

const handleRemove = () => {
  imageUrl.value = ''
  emit('update:modelValue', '')
}
</script>

<style scoped>
.image-uploader {
  width: 100%;
}

.upload-area {
  position: relative;
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-area:hover,
.upload-area.is-dragover {
  border-color: #667eea;
}

.hidden-input {
  display: none;
}

.uploading {
  text-align: center;
  padding: 2rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.preview {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-actions {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.preview:hover .preview-actions {
  opacity: 1;
}

.btn-remove {
  padding: 0.5rem 1rem;
  background: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-remove:hover {
  background: #f78989;
}

.upload-placeholder {
  text-align: center;
  padding: 2rem;
  color: #909399;
}

.upload-placeholder .icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #c0c4cc;
}

.upload-placeholder .hint {
  font-size: 0.75rem;
  margin-top: 0.5rem;
}

.error-msg {
  color: #f56c6c;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}
</style>
