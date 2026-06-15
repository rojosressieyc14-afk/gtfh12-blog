# User Center Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign personal center with two-level sidebar navigation (工作台/设置) and data-rich dashboard.

**Architecture:** Two new backend endpoints (stats + recent activity) feed the dashboard. Sidebar gets tab-based grouping with user card. Overview page shows 4 stat cards + recent activity timeline.

**Tech Stack:** Go/Gin, Vue 3, GORM, MySQL

---

### Task 1: Backend — UserCenterService

**Files:**
- Create: `server/internal/service/user_center_service.go`
- Modify: none

- [ ] **Step 1: Create user center service**

Create `server/internal/service/user_center_service.go`:

```go
package service

import (
	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserCenterService struct {
	db *gorm.DB
}

func NewUserCenterService(db *gorm.DB) *UserCenterService {
	return &UserCenterService{db: db}
}

type UserStats struct {
	ArticlesCount int   `json:"articlesCount"`
	ProjectsCount int   `json:"projectsCount"`
	KbDocsCount   int   `json:"kbDocsCount"`
	TotalViews    int64 `json:"totalViews"`
}

func (s *UserCenterService) GetStats(c *gin.Context) (*UserStats, error) {
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		return &UserStats{}, nil
	}
	userID := authUser.ID

	var articlesCount int64
	s.db.Model(&model.Article{}).Where("author_id = ?", userID).Count(&articlesCount)

	var projectsCount int64
	s.db.Model(&model.Project{}).Where("author_id = ?", userID).Count(&projectsCount)

	var kbDocsCount int64
	s.db.Model(&model.KnowledgeDocument{}).Where("user_id = ?", userID).Count(&kbDocsCount)

	var totalViews int64
	s.db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Where("author_id = ?", userID).Scan(&totalViews)

	return &UserStats{
		ArticlesCount: int(articlesCount),
		ProjectsCount: int(projectsCount),
		KbDocsCount:   int(kbDocsCount),
		TotalViews:    totalViews,
	}, nil
}

type ActivityItem struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *UserCenterService) GetRecentActivity(c *gin.Context) ([]ActivityItem, error) {
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		return nil, nil
	}
	userID := authUser.ID

	var articles []ActivityItem
	s.db.Model(&model.Article{}).
		Select("id, 'article' as type, title, updated_at").
		Where("author_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&articles)

	var projects []ActivityItem
	s.db.Model(&model.Project{}).
		Select("id, 'project' as type, title, updated_at").
		Where("author_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&projects)

	var kbDocs []ActivityItem
	s.db.Model(&model.KnowledgeDocument{}).
		Select("id, 'kb_doc' as type, title, updated_at").
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&kbDocs)

	all := append(append(articles, projects...), kbDocs...)

	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].UpdatedAt > all[i].UpdatedAt {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	if len(all) > 10 {
		all = all[:10]
	}

	return all, nil
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd D:\blog\server && go build ./...`
Expected: no errors

### Task 2: Backend — UserCenterHandler

**Files:**
- Create: `server/internal/handler/user_center_handler.go`

- [ ] **Step 1: Create user center handler**

Create `server/internal/handler/user_center_handler.go`:

```go
package handler

import (
	"net/http"

	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type UserCenterHandler struct {
	userCenterService *service.UserCenterService
}

func NewUserCenterHandler(userCenterService *service.UserCenterService) *UserCenterHandler {
	return &UserCenterHandler{userCenterService: userCenterService}
}

func (h *UserCenterHandler) GetStats(c *gin.Context) {
	stats, err := h.userCenterService.GetStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *UserCenterHandler) GetRecentActivity(c *gin.Context) {
	items, err := h.userCenterService.GetRecentActivity(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取最近动态失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd D:\blog\server && go build ./...`
Expected: no errors

### Task 3: Backend — Wire up router

**Files:**
- Modify: `server/internal/router/router.go`

- [ ] **Step 1: Add import and initialization**

In `server/internal/router/router.go`, after line 44 (`aiReviewService := service.NewAIReviewService(db)`), add:

```go
userCenterService := service.NewUserCenterService(db)
```

After line 54 (`aiReviewHandler := handler.NewAIReviewHandler(aiReviewService)`), add:

```go
userCenterHandler := handler.NewUserCenterHandler(userCenterService)
```

- [ ] **Step 2: Add routes**

Before the `admin := api.Group("/admin")` block (before line 133), add:

```go
		api.GET("/user-center/stats", middleware.RequireAuth(), userCenterHandler.GetStats)
		api.GET("/user-center/recent-activity", middleware.RequireAuth(), userCenterHandler.GetRecentActivity)
```

- [ ] **Step 3: Verify compilation**

Run: `cd D:\blog\server && go build ./...`
Expected: no errors

### Task 4: Frontend — API client

**Files:**
- Create: `web/src/api/userCenter.js`

- [ ] **Step 1: Create API module**

Create `web/src/api/userCenter.js`:

```js
import http from "./http";

export function getUserStats() {
  return http.get("/user-center/stats");
}

export function getRecentActivity() {
  return http.get("/user-center/recent-activity");
}
```

### Task 5: Frontend — UserCenterLayout sidebar redesign

**Files:**
- Modify: `web/src/components/UserCenterLayout.vue`

- [ ] **Step 1: Rewrite template with two-level tabs**

Replace the entire file content:

```vue
<template>
  <section class="uc-layout">
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
  const items = activeTab.value === "workspace" ? workspaceItems : settingsItems;
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
```

### Task 6: Frontend — UserCenterOverview dashboard redesign

**Files:**
- Modify: `web/src/views/UserCenterOverview.vue`

- [ ] **Step 1: Rewrite template with stat cards and recent activity**

Replace the entire file:

```vue
<template>
  <section class="uc-overview">
    <div class="uc-welcome">
      <p class="eyebrow">工作台概览</p>
      <h2>欢迎回来，{{ userStore.profile?.username }}</h2>
    </div>

    <div class="uc-stats">
      <div class="panel-card uc-stat-card">
        <div class="uc-stat-value">{{ stats.articlesCount }}</div>
        <p class="uc-stat-label">文章</p>
      </div>
      <div class="panel-card uc-stat-card">
        <div class="uc-stat-value">{{ stats.projectsCount }}</div>
        <p class="uc-stat-label">项目</p>
      </div>
      <div class="panel-card uc-stat-card">
        <div class="uc-stat-value">{{ stats.kbDocsCount }}</div>
        <p class="uc-stat-label">知识库文档</p>
      </div>
      <div class="panel-card uc-stat-card">
        <div class="uc-stat-value">{{ formatViews(stats.totalViews) }}</div>
        <p class="uc-stat-label">总浏览量</p>
      </div>
    </div>

    <div v-if="activities.length" class="uc-recent">
      <p class="eyebrow">最近动态</p>
      <div class="uc-timeline">
        <div v-for="item in activities" :key="item.type + item.id" class="uc-timeline__item">
          <span class="uc-timeline__dot" :class="'uc-timeline__dot--' + item.type"></span>
          <div class="uc-timeline__body">
            <router-link v-if="item.type === 'article'" :to="`/user-center/editor/${item.id}`" class="inline-link">
              {{ item.title }}
            </router-link>
            <router-link v-else-if="item.type === 'project'" :to="`/user-center/project-editor/${item.id}`" class="inline-link">
              {{ item.title }}
            </router-link>
            <span v-else>{{ item.title }}</span>
            <span class="uc-timeline__type">
              {{ typeLabel(item.type) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="uc-quick-actions">
      <router-link class="panel-card uc-action-card" to="/user-center/editor">
        <span class="uc-action-icon">&#9998;</span>
        <strong>写新文章</strong>
        <p>开始一篇新文章创作</p>
      </router-link>
      <router-link class="panel-card uc-action-card" to="/user-center/project-editor">
        <span class="uc-action-icon">&#10010;</span>
        <strong>新建项目</strong>
        <p>记录一个新项目作品</p>
      </router-link>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useUserStore } from "../stores/user";
import { getUserStats, getRecentActivity } from "../api/userCenter";

const userStore = useUserStore();
const stats = ref({ articlesCount: 0, projectsCount: 0, kbDocsCount: 0, totalViews: 0 });
const activities = ref([]);

function formatViews(val) {
  if (val >= 10000) return (val / 10000).toFixed(1) + "万";
  if (val >= 1000) return (val / 1000).toFixed(1) + "k";
  return String(val);
}

function typeLabel(type) {
  const map = { article: "更新了文章", project: "更新了项目", kb_doc: "编辑了笔记" };
  return map[type] || type;
}

onMounted(async () => {
  try {
    const [{ data: statsData }, { data: activityData }] = await Promise.all([
      getUserStats(),
      getRecentActivity(),
    ]);
    stats.value = statsData.stats;
    activities.value = activityData.items || [];
  } catch (e) {
    console.error("Failed to load user center data", e);
  }
});
</script>

<style scoped>
.uc-overview {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.uc-welcome h2 {
  margin: 4px 0 0;
}

.uc-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.uc-stat-card {
  padding: 22px 20px;
  text-align: center;
}

.uc-stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #f97316;
  line-height: 1.2;
}

.uc-stat-label {
  margin: 6px 0 0;
  font-size: 0.85rem;
  color: var(--text-soft, #a0aec0);
}

.uc-recent {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.uc-timeline {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0;
}

.uc-timeline__item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.04));
}

.uc-timeline__item:last-child {
  border-bottom: none;
}

.uc-timeline__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.uc-timeline__dot--article {
  background: #f97316;
}

.uc-timeline__dot--project {
  background: #60a5fa;
}

.uc-timeline__dot--kb_doc {
  background: #34d399;
}

.uc-timeline__body {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 0.92rem;
}

.uc-timeline__type {
  font-size: 0.78rem;
  color: var(--text-soft, #a0aec0);
}

.uc-quick-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.uc-action-card {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  text-decoration: none;
  color: inherit;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s;
}

.uc-action-card:hover {
  transform: translateY(-2px);
  border-color: rgba(249, 115, 22, 0.3);
}

.uc-action-icon {
  font-size: 1.4rem;
  margin-bottom: 4px;
}

.uc-action-card p {
  margin: 0;
  font-size: 0.85rem;
  color: var(--text-soft, #a0aec0);
}

@media (max-width: 768px) {
  .uc-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .uc-quick-actions {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .uc-stats {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
```

### Task 7: Verify build

- [ ] **Step 1: Verify Go build**

Run: `cd D:\blog\server && go build ./...`
Expected: no errors

- [ ] **Step 2: Verify frontend build**

Run: `cd D:\blog\web && npx vite build`
Expected: no errors

- [ ] **Step 3: Restart services and test**

Run: `cd D:\blog && .\start.ps1`
Expected: all services start without errors
