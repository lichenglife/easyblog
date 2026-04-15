<template>
  <div class="user-management">
    <div class="header">
      <h1>用户管理</h1>
      <div class="filters">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索用户名或邮箱"
          class="search-input"
        />
      </div>
    </div>

    <div class="table-container">
      <table class="user-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>昵称</th>
            <th>邮箱</th>
            <th>手机</th>
            <th>角色</th>
            <th>状态</th>
            <th>注册时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.userID">
            <td>{{ user.userID.slice(0, 8) }}...</td>
            <td>{{ user.username }}</td>
            <td>{{ user.nickname || '-' }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.phone || '-' }}</td>
            <td>
              <span :class="['role-badge', user.role === 2 ? 'role-admin' : 'role-user']">
                {{ user.role === 2 ? '管理员' : '普通用户' }}
              </span>
            </td>
            <td>
              <span :class="['status-badge', getStatusClass(user.status)]">
                {{ getStatusText(user.status) }}
              </span>
            </td>
            <td>{{ formatDate(user.createAt) }}</td>
            <td class="actions">
              <button
                v-if="user.status === 0"
                @click="handleAudit(user, 1)"
                class="btn-audit"
              >
                通过
              </button>
              <button
                v-if="user.status === 0"
                @click="handleAudit(user, 2)"
                class="btn-reject"
              >
                拒绝
              </button>
              <button
                v-if="user.status === 1"
                @click="handleBan(user)"
                class="btn-ban"
              >
                封禁
              </button>
              <button
                v-if="user.status === 2"
                @click="handleUnban(user)"
                class="btn-unban"
              >
                解封
              </button>
              <button @click="handleDelete(user)" class="btn-delete">
                删除
              </button>
            </td>
          </tr>
          <tr v-if="users.length === 0">
            <td colspan="9" class="empty-message">暂无用户数据</td>
          </tr>
        </tbody>
      </table>
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
import { userApi, type UserInfo } from '@/api/user'

const users = ref<UserInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 10
const searchQuery = ref('')

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(
    (user) =>
      user.username.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query)
  )
})

const totalPages = computed(() => Math.ceil(users.value.length / pageSize))

const getStatusClass = (status: number): string => {
  switch (status) {
    case 0:
      return 'status-pending'
    case 1:
      return 'status-normal'
    case 2:
      return 'status-banned'
    default:
      return ''
  }
}

const getStatusText = (status: number): string => {
  switch (status) {
    case 0:
      return '待审核'
    case 1:
      return '正常'
    case 2:
      return '已封禁'
    default:
      return '未知'
  }
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.getUserList({
      page: currentPage.value,
      pageSize,
    })
    users.value = res.users
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAudit = async (user: UserInfo, status: number) => {
  try {
    await userApi.auditUser({ userId: parseInt(user.userID), status })
    await fetchUsers()
  } catch (error) {
    console.error('审核用户失败:', error)
  }
}

const handleBan = async (user: UserInfo) => {
  if (!confirm(`确定要封禁用户 "${user.username}" 吗？`)) return
  try {
    await userApi.banUser(parseInt(user.userID))
    await fetchUsers()
  } catch (error) {
    console.error('封禁用户失败:', error)
  }
}

const handleUnban = async (user: UserInfo) => {
  if (!confirm(`确定要解封用户 "${user.username}" 吗？`)) return
  try {
    await userApi.auditUser({ userId: parseInt(user.userID), status: 1 })
    await fetchUsers()
  } catch (error) {
    console.error('解封用户失败:', error)
  }
}

const handleDelete = async (user: UserInfo) => {
  if (!confirm(`确定要删除用户 "${user.username}" 吗？此操作不可恢复！`)) return
  try {
    await userApi.deleteUser(parseInt(user.userID))
    await fetchUsers()
  } catch (error) {
    console.error('删除用户失败:', error)
  }
}

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-management {
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

.filters {
  display: flex;
  gap: 1rem;
}

.search-input {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  width: 250px;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.table-container {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.user-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.user-table th,
.user-table td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.user-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  font-size: 0.875rem;
}

.user-table td {
  font-size: 0.875rem;
  color: #666;
}

.user-table tbody tr:hover {
  background: #f8f9fa;
}

.role-badge,
.status-badge {
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

.status-pending {
  background: #fef3c7;
  color: #d97706;
}

.status-normal {
  background: #d1fae5;
  color: #059669;
}

.status-banned {
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

.btn-audit {
  background: #10b981;
  color: #fff;
}

.btn-reject {
  background: #f59e0b;
  color: #fff;
}

.btn-ban {
  background: #ef4444;
  color: #fff;
}

.btn-unban {
  background: #3b82f6;
  color: #fff;
}

.btn-delete {
  background: #6b7280;
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
