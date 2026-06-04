<template>
  <header class="site-header site-header--brand">
    <div class="site-header__main">
      <div class="site-brand" @click="$router.push('/')">
        <div class="brand-badge brand-badge--cluster">
          <span class="brand-shape brand-shape--square"></span>
          <span class="brand-shape brand-shape--circle"></span>
          <span class="brand-shape brand-shape--triangle"></span>
        </div>
        <div>
          <h1>PulseBlog</h1>
          <p>作品集、学习归档与长期写作空间</p>
        </div>
      </div>
      <p class="site-header__tagline">把博客做成真正能对外展示的个人品牌站</p>
    </div>

    <nav class="site-nav">
      <router-link to="/">首页</router-link>
      <router-link to="/about">关于</router-link>
      <router-link to="/articles">文章</router-link>
      <router-link to="/projects">项目</router-link>
      <router-link to="/interview">AI 面试官</router-link>
      <router-link to="/editor">写文章</router-link>
      <router-link v-if="userStore.isLoggedIn" to="/my-articles">我的文章</router-link>
      <router-link v-if="userStore.isLoggedIn" to="/my-projects">我的项目</router-link>
      <router-link v-if="userStore.isLoggedIn" to="/collections">收藏夹</router-link>
      <router-link v-if="userStore.isLoggedIn" to="/notifications">通知</router-link>
      <router-link v-if="userStore.isLoggedIn" to="/profile">个人资料</router-link>
    </nav>

    <div class="site-header__actions">
      <a
        v-if="userStore.isLoggedIn && userStore.profile?.role === 'admin'"
        class="ghost-btn"
        href="http://localhost:3001"
        target="_blank"
        rel="noreferrer"
      >
        后台管理
      </a>
      <router-link v-if="!userStore.isLoggedIn" class="solid-btn" to="/auth">登录 / 注册</router-link>
      <button v-else class="ghost-btn" @click="logout">{{ userStore.profile?.username }} | 退出登录</button>
    </div>

    <button class="hamburger-btn" :class="{ 'hamburger-btn--active': mobileOpen }" @click="toggleMobile" aria-label="菜单">
      <span></span><span></span><span></span>
    </button>

    <Teleport to="body">
      <div v-if="mobileOpen" class="mobile-overlay" @click="mobileOpen = false"></div>
      <aside class="mobile-drawer" :class="{ 'mobile-drawer--open': mobileOpen }">
        <nav class="mobile-drawer__nav" @click="mobileOpen = false">
          <router-link to="/">首页</router-link>
          <router-link to="/about">关于</router-link>
          <router-link to="/articles">文章</router-link>
          <router-link to="/projects">项目</router-link>
          <router-link to="/interview">AI 面试官</router-link>
          <router-link to="/editor">写文章</router-link>
          <router-link v-if="userStore.isLoggedIn" to="/my-articles">我的文章</router-link>
          <router-link v-if="userStore.isLoggedIn" to="/my-projects">我的项目</router-link>
          <router-link v-if="userStore.isLoggedIn" to="/collections">收藏夹</router-link>
          <router-link v-if="userStore.isLoggedIn" to="/notifications">通知</router-link>
          <router-link v-if="userStore.isLoggedIn" to="/profile">个人资料</router-link>
          <hr v-if="userStore.isLoggedIn" />
          <a v-if="userStore.isLoggedIn && userStore.profile?.role === 'admin'" href="http://localhost:3001" target="_blank" rel="noreferrer">后台管理</a>
          <router-link v-if="!userStore.isLoggedIn" to="/auth">登录 / 注册</router-link>
          <button v-else class="mobile-logout" @click="logout">退出登录</button>
        </nav>
      </aside>
    </Teleport>
  </header>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "../stores/user";

const router = useRouter();
const userStore = useUserStore();
const mobileOpen = ref(false);

function toggleMobile() {
  mobileOpen.value = !mobileOpen.value;
}

function logout() {
  userStore.logout();
  mobileOpen.value = false;
  router.push("/");
}
</script>

<style scoped>
.site-header {
  display: grid;
  grid-template-columns: minmax(260px, 0.9fr) minmax(0, 1.2fr) auto auto;
  align-items: center;
}

.site-header__main {
  display: grid;
  gap: 8px;
}

.site-header__tagline {
  margin: 0;
  color: var(--text-soft);
  font-size: 0.88rem;
}

.site-header__actions {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.site-header--brand .brand-badge--cluster {
  position: relative;
  overflow: visible;
  background: linear-gradient(135deg, rgba(255, 209, 102, 0.18), rgba(255, 138, 76, 0.3));
}

.brand-shape {
  position: absolute;
  display: block;
  border: 2px solid rgba(23, 18, 15, 0.85);
  background: rgba(255, 248, 240, 0.5);
}

.brand-shape--square {
  width: 18px;
  height: 18px;
  left: 8px;
  top: 12px;
  border-radius: 6px;
}

.brand-shape--circle {
  width: 14px;
  height: 14px;
  right: 8px;
  top: 9px;
  border-radius: 999px;
}

.brand-shape--triangle {
  width: 16px;
  height: 16px;
  left: 18px;
  bottom: 8px;
  clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
  border-radius: 4px;
}

/* Hamburger button */
.hamburger-btn {
  display: none;
  flex-direction: column;
  gap: 5px;
  padding: 10px;
  background: none;
  border: none;
  cursor: pointer;
  z-index: 101;
}

.hamburger-btn span {
  display: block;
  width: 24px;
  height: 2px;
  background: var(--text, #f7f3ea);
  border-radius: 2px;
  transition: transform 0.3s, opacity 0.3s;
}

.hamburger-btn--active span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}
.hamburger-btn--active span:nth-child(2) {
  opacity: 0;
}
.hamburger-btn--active span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

/* Mobile drawer */
:global(.mobile-overlay) {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 200;
}

:global(.mobile-drawer) {
  position: fixed;
  top: 0;
  right: 0;
  width: min(300px, 80vw);
  height: 100vh;
  background: var(--panel, #1a1e26);
  border-left: 1px solid var(--border, rgba(255,255,255,0.1));
  z-index: 201;
  padding: 80px 20px 20px;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  overflow-y: auto;
}

:global(.mobile-drawer--open) {
  transform: translateX(0);
}

.mobile-drawer__nav {
  display: grid;
  gap: 6px;
}

.mobile-drawer__nav a,
.mobile-logout {
  padding: 12px 14px;
  border-radius: 12px;
  color: var(--text, #f7f3ea);
  text-decoration: none;
  font-size: 1rem;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
}

.mobile-drawer__nav a:hover,
.mobile-drawer__nav a.router-link-exact-active {
  background: rgba(255,255,255,0.06);
}

.mobile-drawer__nav hr {
  border: none;
  border-top: 1px solid var(--border, rgba(255,255,255,0.1));
  margin: 8px 0;
}

.mobile-logout {
  color: #ffbcbc;
}

/* === Responsive === */
@media (max-width: 1024px) {
  .site-nav {
    display: none;
  }
  .site-header__actions {
    display: none;
  }
  .hamburger-btn {
    display: flex;
  }
}

@media (max-width: 960px) {
  .site-header {
    grid-template-columns: 1fr auto;
  }
  .site-header__tagline {
    display: none;
  }
  .site-header__main .site-brand p {
    display: none;
  }
}
</style>
