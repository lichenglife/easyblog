<template>
  <div class="profile-page">
    <div class="profile-container">
      <h1>个人中心</h1>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else class="profile-content">
        <div class="profile-header">
          <div class="avatar">
            <img :src="userInfo?.avatar || '/default-avatar.png'" alt="头像" />
          </div>
          <div class="user-info">
            <h2>{{ userInfo?.nickname || userInfo?.username }}</h2>
            <p class="username">@{{ userInfo?.username }}</p>
            <span :class="['role-badge', userInfo?.role === 2 ? 'role-admin' : 'role-user']">
              {{ userInfo?.role === 2 ? '管理员' : '普通用户' }}
            </span>
          </div>
        </div>

        <div class="profile-body">
          <div class="info-section">
            <h3>基本信息</h3>
            <div class="info-item">
              <label>邮箱</label>
              <span>{{ userInfo?.email || '未设置' }}</span>
            </div>
            <div class="info-item">
              <label>手机</label>
              <span>{{ userInfo?.phone || '未设置' }}</span>
            </div>
            <div class="info-item">
              <label>个人简介</label>
              <span>{{ userInfo?.bio || '暂无简介' }}</span>
            </div>
            <div class="info-item">
              <label>注册时间</label>
              <span>{{ formatDate(userInfo?.createAt || '') }}</span>
            </div>
          </div>

          <div class="actions-section">
            <h3>账号设置</h3>
            <button @click="showEditModal = true" class="btn-edit">
              编辑资料
            </button>
            <button @click="showPasswordModal = true" class="btn-password">
              修改密码
            </button>
            <button @click="handleLogout" class="btn-logout">
              退出登录
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑资料模态框 -->
    <div v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
      <div class="modal" @click.stop>
        <h3>编辑资料</h3>
        <form @submit.prevent="handleEditProfile">
          <div class="form-group">
            <label>昵称</label>
            <input v-model="editForm.nickname" type="text" placeholder="请输入昵称" />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="editForm.email" type="email" placeholder="请输入邮箱" />
          </div>
          <div class="form-group">
            <label>手机</label>
            <input v-model="editForm.phone" type="tel" placeholder="请输入手机号" />
          </div>
          <div class="form-group">
            <label>个人简介</label>
            <textarea v-model="editForm.bio" rows="3" placeholder="请输入个人简介"></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showEditModal = false" class="btn-cancel">
              取消
            </button>
            <button type="submit" :disabled="saving" class="btn-save">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 修改密码模态框 -->
    <div v-if="showPasswordModal" class="modal-overlay" @click="showPasswordModal = false">
      <div class="modal" @click.stop>
        <h3>修改密码</h3>
        <form @submit.prevent="handleChangePassword">
          <div class="form-group">
            <label>当前密码</label>
            <input v-model="passwordForm.oldPassword" type="password" placeholder="请输入当前密码" required />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" required />
          </div>
          <div class="form-group">
            <label>确认密码</label>
            <input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" required />
          </div>
          <div v-if="passwordError" class="error-text">{{ passwordError }}</div>
          <div class="modal-actions">
            <button type="button" @click="showPasswordModal = false" class="btn-cancel">
              取消
            </button>
            <button type="submit" :disabled="saving" class="btn-save">
              {{ saving ? '修改中...' : '确认修改' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(true)
const userInfo = ref(userStore.userInfo)

const showEditModal = ref(false)
const showPasswordModal = ref(false)
const saving = ref(false)
const passwordError = ref('')

const editForm = reactive({
  nickname: '',
  email: '',
  phone: '',
  bio: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const fetchUserInfo = async () => {
  loading.value = true
  try {
    const user = await userApi.getCurrentUser()
    userInfo.value = user
    userStore.setUserInfo(user)
    Object.assign(editForm, {
      nickname: user.nickname || '',
      email: user.email || '',
      phone: user.phone || '',
      bio: user.bio || '',
    })
  } catch (error) {
    console.error('获取用户信息失败:', error)
  } finally {
    loading.value = false
  }
}

const handleEditProfile = async () => {
  saving.value = true
  try {
    await userApi.updateUserInfo(userInfo.value!.username, {
      nickname: editForm.nickname,
      email: editForm.email,
      phone: editForm.phone,
    })
    showEditModal.value = false
    await fetchUserInfo()
    alert('资料更新成功')
  } catch (error) {
    console.error('更新资料失败:', error)
    alert('更新失败')
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async () => {
  passwordError.value = ''
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordError.value = '两次输入的新密码不一致'
    return
  }
  if (passwordForm.newPassword.length < 6) {
    passwordError.value = '密码长度至少 6 位'
    return
  }

  saving.value = true
  try {
    await userApi.changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    showPasswordModal.value = false
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    alert('密码修改成功，请重新登录')
    userStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('修改密码失败:', error)
    passwordError.value = '修改失败，请检查当前密码是否正确'
  } finally {
    saving.value = false
  }
}

const handleLogout = () => {
  if (!confirm('确定要退出登录吗？')) return
  userStore.logout()
  router.push('/')
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    fetchUserInfo()
  } else {
    router.push('/login')
  }
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 2rem 1rem;
}

.profile-container {
  max-width: 800px;
  margin: 0 auto;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.profile-container h1 {
  margin: 0;
  padding: 1.5rem 2rem;
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
  font-size: 1.25rem;
  color: #333;
}

.loading {
  padding: 3rem;
  text-align: center;
  color: #666;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 2rem;
  border-bottom: 1px solid #eee;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info h2 {
  margin: 0 0 0.5rem;
  font-size: 1.25rem;
  color: #333;
}

.username {
  margin: 0 0 0.5rem;
  color: #666;
  font-size: 0.875rem;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.role-admin {
  background: #fee2e2;
  color: #dc2626;
}

.role-user {
  background: #e0e7ff;
  color: #4f46e5;
}

.profile-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  padding: 2rem;
}

@media (max-width: 640px) {
  .profile-body {
    grid-template-columns: 1fr;
  }
}

.info-section h3,
.actions-section h3 {
  margin: 0 0 1rem;
  font-size: 1rem;
  color: #333;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid #eee;
}

.info-item label {
  color: #666;
  font-size: 0.875rem;
}

.info-item span {
  color: #333;
  font-size: 0.875rem;
}

.actions-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.actions-section button {
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-edit {
  background: #667eea;
  color: #fff;
}

.btn-password {
  background: #f59e0b;
  color: #fff;
}

.btn-logout {
  background: #6b7280;
  color: #fff;
}

.actions-section button:hover {
  opacity: 0.9;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  border-radius: 8px;
  padding: 1.5rem;
  width: 100%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal h3 {
  margin: 0 0 1.5rem;
  font-size: 1.125rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-size: 0.875rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.error-text {
  color: #dc2626;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

.modal-actions button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
}

.btn-cancel {
  background: #eee;
  color: #333;
}

.btn-save {
  background: #667eea;
  color: #fff;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
