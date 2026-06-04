<template>
  <div>
    <nav v-if="showQuickNav" class="admin-quick-nav">
      <router-link to="/">后台总览</router-link>
      <router-link to="/comments">评论管理</router-link>
      <router-link to="/logs">操作日志</router-link>
      <router-link to="/uploads">资源库</router-link>
      <router-link to="/moderation-hits">风控命中</router-link>
      <router-link to="/sensitive-words">敏感词管理</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useAdminStore } from "./stores/auth";

const store = useAdminStore();
const route = useRoute();
const showQuickNav = computed(() => route.name !== "login" && localStorage.getItem("admin_token"));

onMounted(() => {
  store.fetchProfile();
});
</script>

<style scoped>
.admin-quick-nav {
  position: fixed;
  right: 20px;
  top: 18px;
  z-index: 50;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.admin-quick-nav a {
  padding: 10px 14px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(10, 14, 19, 0.7);
  color: #f7f3ea;
  text-decoration: none;
  backdrop-filter: blur(12px);
  transition: transform 0.22s ease, background 0.22s ease, border-color 0.22s ease;
}

.admin-quick-nav a:hover {
  transform: translateY(-1px);
  background: rgba(16, 22, 30, 0.9);
  border-color: rgba(255, 217, 142, 0.24);
}

.admin-quick-nav a.router-link-exact-active {
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
}

@media (max-width: 768px) {
  .admin-quick-nav {
    right: 10px;
    top: 10px;
    gap: 6px;
  }
  .admin-quick-nav a {
    padding: 8px 12px;
    font-size: 0.85rem;
  }
}

@media (max-width: 480px) {
  .admin-quick-nav {
    right: 0;
    top: auto;
    bottom: 0;
    left: 0;
    width: 100%;
    justify-content: center;
    gap: 4px;
    padding: 8px 6px;
    background: rgba(10, 14, 19, 0.95);
    backdrop-filter: blur(12px);
    border-top: 1px solid var(--border);
    border-radius: 16px 16px 0 0;
    z-index: 100;
    flex-wrap: nowrap;
    overflow-x: auto;
  }
  .admin-quick-nav a {
    padding: 6px 10px;
    font-size: 0.75rem;
    white-space: nowrap;
  }
}
</style>
