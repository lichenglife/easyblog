<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <h2>管理后台</h2>
      <nav>
        <router-link to="/admin/users">用户管理</router-link>
        <router-link to="/admin/articles">文章管理</router-link>
        <router-link to="/admin/categories">分类管理</router-link>
        <router-link to="/admin/tags">标签管理</router-link>
      </nav>
    </aside>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
// 管理后台布局
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 200px;
  background: #333;
  color: #fff;
  padding: 1rem;
}

.sidebar h2 {
  margin-bottom: 1.5rem;
}

.sidebar nav {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.sidebar a {
  color: #ccc;
  text-decoration: none;
  padding: 0.5rem;
  border-radius: 4px;
}

.sidebar a:hover,
.sidebar a.router-link-active {
  background: #444;
  color: #fff;
}

.main {
  flex: 1;
  padding: 2rem;
  background: #f5f5f5;
}
</style>
