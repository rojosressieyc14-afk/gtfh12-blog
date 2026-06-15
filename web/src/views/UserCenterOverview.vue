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
