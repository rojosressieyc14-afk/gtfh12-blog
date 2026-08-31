<template>
  <section class="uc-overview">
    <section class="uc-hero">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">工作台概览</p>
        <h2 class="uc-hero__title">欢迎回来，{{ userStore.profile?.username }}</h2>
        <p class="uc-hero__text">管理你的文章、项目和知识库，持续输出有价值的内容。</p>
      </div>
      <div class="uc-hero__stats">
        <article class="uc-hero__stat">
          <strong>{{ stats.articlesCount }}</strong>
          <span>文章</span>
        </article>
        <article class="uc-hero__stat">
          <strong>{{ stats.projectsCount }}</strong>
          <span>项目</span>
        </article>
        <article class="uc-hero__stat">
          <strong>{{ stats.kbDocsCount }}</strong>
          <span>知识库文档</span>
        </article>
        <article class="uc-hero__stat">
          <strong>{{ formatViews(stats.totalViews) }}</strong>
          <span>总浏览量</span>
        </article>
      </div>
    </section>

    <section v-if="activities.length" class="uc-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">最近动态</p>
          <h3>更新记录</h3>
        </div>
      </div>
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
    </section>

    <section class="uc-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">快速开始</p>
          <h3>创作入口</h3>
        </div>
      </div>
      <div class="uc-quick-actions">
        <router-link class="uc-action-card" to="/user-center/editor">
          <span class="uc-action-icon">&#9998;</span>
          <strong>写新文章</strong>
          <p>开始一篇新文章创作</p>
        </router-link>
        <router-link class="uc-action-card" to="/user-center/project-editor">
          <span class="uc-action-icon">&#10010;</span>
          <strong>新建项目</strong>
          <p>记录一个新项目作品</p>
        </router-link>
      </div>
    </section>
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

.uc-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  padding: clamp(32px, 6vw, 56px) clamp(24px, 4vw, 48px);
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(240px, 0.7fr);
  gap: 32px;
  align-items: center;
}

.uc-hero__copy {
  display: grid;
  gap: 12px;
}

.uc-hero__eyebrow {
  color: #b4530a;
  font-weight: 600;
}

.uc-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
  color: #1d1d1f;
}

.uc-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: #48484a;
}

.uc-hero__stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.uc-hero__stat {
  padding: 18px 16px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
  text-align: center;
}

.uc-hero__stat strong {
  display: block;
  font-size: 1.8rem;
  font-weight: 700;
  color: #1d1d1f;
  line-height: 1.2;
}

.uc-hero__stat span {
  display: block;
  margin-top: 4px;
  color: #6b7280;
  font-size: 0.85rem;
}

.uc-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.uc-timeline {
  display: flex;
  flex-direction: column;
}

.uc-timeline__item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
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
  color: #9ca3af;
}

.uc-quick-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.uc-action-card {
  padding: 24px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-decoration: none;
  color: #1d1d1f;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}

.uc-action-card:hover {
  transform: translateY(-2px);
  border-color: rgba(249, 115, 22, 0.3);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
}

.uc-action-icon {
  font-size: 1.5rem;
  margin-bottom: 4px;
}

.uc-action-card strong {
  font-size: 1.05rem;
  color: #1d1d1f;
}

.uc-action-card p {
  margin: 0;
  font-size: 0.85rem;
  color: #6b7280;
}

@media (max-width: 768px) {
  .uc-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .uc-hero__stats {
    grid-template-columns: repeat(4, 1fr);
  }

  .uc-quick-actions {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .uc-hero__stats {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
