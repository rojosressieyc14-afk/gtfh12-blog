<template>
  <section id="overview" class="admin-hero dashboard-hero">
    <div>
      <p class="admin-label">后台总览</p>
      <h1>内容发布与审核面板</h1>
      <p>优先处理待审内容，再统一管理用户、文章、项目以及内容元数据。</p>
      <div class="hero-actions">
        <button @click="$emit('refresh')">刷新数据</button>
        <router-link class="hero-link" to="/moderation-hits">风控命中</router-link>
        <router-link class="hero-link" to="/sensitive-words">敏感词</router-link>
        <router-link class="hero-link" to="/logs">操作日志</router-link>
        <router-link class="hero-link" to="/uploads">资源库</router-link>
      </div>
    </div>

    <div class="hero-stats">
      <article v-for="item in statCards" :key="item.key" class="hero-stat-card">
        <span>{{ item.label }}</span>
        <strong>{{ stats[item.key] ?? 0 }}</strong>
      </article>
    </div>
  </section>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  stats: { type: Object, required: true }
});

defineEmits(["refresh"]);

const statCards = computed(() => [
  { key: "users", label: "用户总数" },
  { key: "draft", label: "文章草稿" },
  { key: "pending", label: "待审文章" },
  { key: "published", label: "已发文章" },
  { key: "pendingProjects", label: "待审项目" },
  { key: "projects", label: "已发项目" },
  { key: "draftProjects", label: "项目草稿" },
  { key: "moderationHits", label: "风控命中" }
]);
</script>

<style scoped>
.dashboard-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(300px, 0.85fr);
  gap: 20px;
}

.hero-actions {
  margin-top: 18px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-actions button,
.hero-link {
  padding: 12px 16px;
  border-radius: 16px;
  border: 1px solid var(--border);
  text-decoration: none;
}

.hero-actions button {
  cursor: pointer;
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
}

.hero-link {
  color: inherit;
  background: rgba(255, 255, 255, 0.05);
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.hero-stat-card {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(9, 14, 20, 0.38);
}

.hero-stat-card span {
  color: var(--soft);
}

.hero-stat-card strong {
  display: block;
  margin-top: 10px;
  font-size: 2rem;
}

@media (max-width: 960px) {
  .dashboard-hero {
    grid-template-columns: 1fr;
  }
  .hero-stats {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .hero-stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 8px;
  }
  .hero-stat-card {
    padding: 12px;
  }
  .hero-stat-card strong {
    font-size: 1.3rem;
  }
  .hero-actions {
    flex-direction: column;
  }
  .hero-actions button,
  .hero-link {
    width: 100%;
    text-align: center;
  }
}
</style>
