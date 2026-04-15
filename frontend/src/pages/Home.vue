<template>
  <div class="home">
    <nav class="navbar">
      <div class="container">
        <router-link to="/" class="logo">Easyblog</router-link>
        <div class="nav-links">
          <router-link to="/articles">文章</router-link>
          <template v-if="userStore.isLoggedIn">
            <router-link to="/user/profile">个人中心</router-link>
            <router-link v-if="userStore.isAdmin" to="/admin">管理后台</router-link>
            <button @click="handleLogout">退出</button>
          </template>
          <template v-else>
            <router-link to="/login">登录</router-link>
            <router-link to="/register">注册</router-link>
          </template>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <div class="hero">
        <h1>欢迎来到 Easyblog</h1>
        <p>记录技术，分享生活</p>
        <router-link to="/articles" class="btn-primary">开始阅读</router-link>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.navbar {
  background: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
}

.navbar .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-links a {
  color: #666;
  text-decoration: none;
}

.nav-links a:hover {
  color: #007bff;
}

.nav-links button {
  background: none;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.main-content {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 64px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.hero {
  text-align: center;
  color: #fff;
}

.hero h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.25rem;
  margin-bottom: 2rem;
}

.btn-primary {
  display: inline-block;
  background: #fff;
  color: #667eea;
  padding: 0.75rem 2rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: 600;
  transition: transform 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
}
</style>
