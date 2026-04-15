<template>
  <div class="register-page">
    <div class="register-container">
      <h1>用户注册</h1>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="formData.username"
            type="text"
            placeholder="4-30 位字母、数字、下划线"
            required
          />
        </div>
        <div class="form-group">
          <label for="email">邮箱</label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            placeholder="请输入邮箱"
            required
          />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            placeholder="6-30 位"
            required
          />
        </div>
        <div class="form-group">
          <label for="nick_name">昵称</label>
          <input
            id="nick_name"
            v-model="formData.nick_name"
            type="text"
            placeholder="2-30 位"
            required
          />
        </div>
        <div class="form-group">
          <label for="phone">手机号（可选）</label>
          <input
            id="phone"
            v-model="formData.phone"
            type="tel"
            placeholder="请输入手机号"
          />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn-submit" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </div>
        <div class="form-footer">
          <span>已有账号？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const error = ref('')
const success = ref('')

const formData = reactive({
  username: '',
  email: '',
  password: '',
  nick_name: '',
  phone: '',
})

const handleRegister = async () => {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    await userStore.register({
      username: formData.username,
      email: formData.email,
      password: formData.password,
      nick_name: formData.nick_name,
      phone: formData.phone || undefined,
    })
    success.value = '注册成功，请等待管理员审核'
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-container {
  background: #fff;
  padding: 2.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 400px;
}

.register-container h1 {
  text-align: center;
  color: #333;
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-actions {
  margin: 1.5rem 0;
}

.btn-submit {
  width: 100%;
  padding: 0.75rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-submit:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-footer {
  text-align: center;
  color: #666;
}

.form-footer a {
  color: #667eea;
  text-decoration: none;
}

.form-footer a:hover {
  text-decoration: underline;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fee;
  color: #c00;
  border-radius: 4px;
  text-align: center;
}

.success-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #efe;
  color: #0a0;
  border-radius: 4px;
  text-align: center;
}
</style>
