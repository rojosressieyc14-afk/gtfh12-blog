<template>
  <section v-if="route.meta?.sidebar !== false" class="uc-layout">
    <aside class="uc-sidebar">
      <div class="uc-sidebar__user">
        <div class="uc-avatar">{{ avatarFallback }}</div>
        <div class="uc-user-info">
          <strong>{{ userStore.profile?.username }}</strong>
          <span class="uc-role-badge" :class="userStore.isAdmin ? 'uc-role-badge--admin' : 'uc-role-badge--user'">
            {{ userStore.isAdmin ? "管理员" : "用户" }}
          </span>
        </div>
      </div>

      <div class="uc-tab-bar">
        <button class="uc-tab" :class="{ 'uc-tab--active': activeTab === 'workspace' }" @click="activeTab = 'workspace'">工作台</button>
        <button class="uc-tab" :class="{ 'uc-tab--active': activeTab === 'settings' }" @click="activeTab = 'settings'">设置</button>
      </div>

      <nav class="uc-nav">
        <template v-for="item in visibleItems" :key="item.path">
          <a
            v-if="item.external"
            :href="adminUrl"
            target="_blank"
            rel="noreferrer"
            class="uc-nav__link"
          >
            <span class="uc-nav__indicator"></span>
            <span v-html="item.icon" class="uc-nav__icon"></span>
            <span>{{ item.label }}</span>
          </a>
          <router-link
            v-else
            :to="item.path"
            class="uc-nav__link"
            :class="{ 'uc-nav__link--active': isActive(item.path) }"
          >
            <span class="uc-nav__indicator"></span>
            <span v-html="item.icon" class="uc-nav__icon"></span>
            <span>{{ item.label }}</span>
          </router-link>
        </template>
      </nav>

      <div v-if="activeTab === 'workspace'" class="uc-sidebar__footer">
        <router-link class="uc-write-btn" to="/user-center/editor">
          <span class="uc-write-icon">+</span>
          写文章
        </router-link>
      </div>
    </aside>
    <main class="uc-content">
      <router-view />
    </main>
  </section>
  <router-view v-else />
</template>

<script setup>
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { useUserStore } from "../stores/user";

const route = useRoute();
const userStore = useUserStore();
const activeTab = ref("workspace");

const adminUrl = computed(() => {
  return import.meta.env.VITE_ADMIN_URL || "/PulseBlog/admin/";
});

const workspaceItems = [
  { path: "/user-center/overview", icon: "&#9675;", label: "概览" },
  { path: "/user-center/articles", icon: "&#9632;", label: "文章" },
  { path: "/user-center/projects", icon: "&#9650;", label: "项目" },
  { path: "/user-center/knowledge-base", icon: "&#9679;", label: "知识库" },
  { path: "/user-center/collections", icon: "&#9733;", label: "收藏" },
];

const settingsItems = [
  { path: "/user-center/profile", icon: "&#9671;", label: "个人资料" },
  { path: "/user-center/api-keys", icon: "&#9670;", label: "API Key" },
  { path: "/user-center/notifications", icon: "&#9702;", label: "通知" },
  { path: "/user-center/about", icon: "&#9432;", label: "关于" },
];

const visibleItems = computed(() => {
  const items = activeTab.value === "workspace" ? [...workspaceItems] : [...settingsItems];
  if (activeTab.value === "settings" && userStore.isAdmin) {
    items.push({
      path: "/user-center/admin",
      icon: "&#9878;",
      label: "后台管理",
      external: true,
    });
  }
  return items;
});

const avatarFallback = computed(() => {
  const name = userStore.profile?.username || "U";
  return name.charAt(0).toUpperCase();
});

function isActive(path) {
  if (path === "/user-center/overview") {
    return route.path === path;
  }
  return route.path.startsWith(path + "/") || route.path === path;
}
</script>

<style scoped>
.uc-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 0;
  min-height: calc(100vh - 100px);
}

.uc-sidebar {
  position: sticky;
  top: 20px;
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
  height: fit-content;
  min-height: calc(100vh - 120px);
  border-right: 1px solid var(--border, rgba(255,255,255,0.08));
}

.uc-sidebar__user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 6px 20px;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.06));
  margin-bottom: 16px;
}

.uc-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(249, 115, 22, 0.4));
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 1rem;
  color: #f97316;
  flex-shrink: 0;
}

.uc-user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.uc-user-info strong {
  font-size: 0.95rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.uc-role-badge {
  font-size: 0.72rem;
  padding: 2px 8px;
  border-radius: 999px;
  width: fit-content;
}

.uc-role-badge--admin {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.uc-role-badge--user {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.uc-tab-bar {
  display: flex;
  background: rgba(255,255,255,0.04);
  border-radius: 10px;
  padding: 3px;
  margin-bottom: 16px;
}

.uc-tab {
  flex: 1;
  padding: 8px 0;
  border: none;
  background: none;
  color: var(--text-soft, #a0aec0);
  cursor: pointer;
  font-size: 0.85rem;
  border-radius: 8px;
  transition: all 0.2s;
}

.uc-tab--active {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
  font-weight: 600;
}

.uc-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.uc-nav__link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 10px;
  color: var(--text-soft, #a0aec0);
  text-decoration: none;
  font-size: 0.92rem;
  transition: all 0.2s;
  position: relative;
}

.uc-nav__link:hover {
  background: rgba(255,255,255,0.05);
  color: var(--text, #f7f3ea);
}

.uc-nav__link--active {
  background: rgba(249, 115, 22, 0.1);
  color: #f97316;
  font-weight: 600;
}

.uc-nav__indicator {
  width: 3px;
  height: 18px;
  border-radius: 2px;
  background: transparent;
  transition: all 0.2s;
  position: absolute;
  left: 0;
}

.uc-nav__link--active .uc-nav__indicator {
  background: #f97316;
}

.uc-nav__icon {
  width: 20px;
  text-align: center;
  font-size: 0.9rem;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.uc-sidebar__footer {
  padding-top: 16px;
  border-top: 1px solid var(--border, rgba(255,255,255,0.06));
  margin-top: auto;
}

.uc-write-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.15), rgba(249, 115, 22, 0.25));
  color: #f97316;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 600;
  transition: all 0.2s;
}

.uc-write-btn:hover {
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.25), rgba(249, 115, 22, 0.35));
}

.uc-write-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(249, 115, 22, 0.3);
  font-size: 1rem;
  font-weight: 700;
  line-height: 1;
}

.uc-content {
  padding: 32px 40px;
  min-width: 0;
}

@media (max-width: 768px) {
  .uc-layout {
    grid-template-columns: 1fr;
  }

  .uc-sidebar {
    position: static;
    min-height: auto;
    border-right: none;
    border-bottom: 1px solid var(--border, rgba(255,255,255,0.08));
    padding: 16px;
  }

  .uc-sidebar__user {
    padding-bottom: 12px;
    margin-bottom: 12px;
  }

  .uc-nav {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 4px;
  }

  .uc-nav__link {
    padding: 8px 12px;
    font-size: 0.85rem;
  }

  .uc-nav__indicator {
    display: none;
  }

  .uc-sidebar__footer {
    margin-top: 12px;
  }

  .uc-content {
    padding: 20px 16px;
  }
}
</style>
