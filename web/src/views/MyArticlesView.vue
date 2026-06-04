<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">创作空间</p>
        <h3>我的文章</h3>
      </div>
      <router-link class="solid-btn" to="/editor">写新文章</router-link>
    </div>

    <div class="services-grid">
      <article class="panel-card service-card">
        <p class="eyebrow">管理</p>
        <h4>统一整理草稿、已发布和审核中的内容</h4>
        <p>这里适合持续维护你的学习笔记、项目复盘和公开输出，不需要在多个地方来回切换。</p>
      </article>
      <article class="panel-card service-card">
        <p class="eyebrow">建议</p>
        <h4>优先补齐摘要、标签和分类</h4>
        <p>这些信息会直接影响首页展示、文章库检索和审核时的可读性。</p>
      </article>
    </div>

    <div v-if="articles.length" class="article-grid">
      <article v-for="item in articles" :key="item.id" class="article-card article-card--mine">
        <div class="article-card__meta">
          <span class="status-chip" :class="item.status">{{ labelMap[item.status] || "未知状态" }}</span>
          <span>{{ item.category?.name || "未分类" }} | {{ formatDate(item.updatedAt) }}</span>
        </div>

        <h3>{{ item.title }}</h3>
        <p>{{ item.summary || "这篇文章还没有摘要，建议补一小段，让前台展示和审核过程都更清楚。" }}</p>

        <div v-if="item.tags?.length" class="tag-row">
          <span v-for="tag in item.tags" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
        </div>

        <footer class="my-article-footer">
          <router-link class="inline-link" :to="`/editor/${item.id}`">继续编辑</router-link>
          <router-link v-if="item.status === 'published'" class="inline-link" :to="`/article/${item.id}`">查看详情</router-link>
        </footer>

        <p v-if="item.rejectReason" class="reject-tip">驳回原因：{{ item.rejectReason }}</p>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>你还没有文章</h4>
      <p>先写一篇自我介绍、项目复盘或学习笔记，这里会逐渐变成你的个人内容后台。</p>
    </div>

    <div class="pager-row">
      <button class="ghost-btn" :disabled="page <= 1" @click="goToPage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button class="ghost-btn" :disabled="page >= totalPages" @click="goToPage(page + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { listMyArticles } from "../api/article";

const articles = ref([]);
const page = ref(1);
const pageSize = 8;
const total = ref(0);

const labelMap = {
  draft: "草稿",
  pending: "审核中",
  published: "已发布",
  rejected: "已驳回"
};

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

async function loadMine() {
  const { data } = await listMyArticles(page.value, pageSize);
  articles.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function goToPage(nextPage) {
  page.value = nextPage;
  await loadMine();
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(loadMine);
</script>

<style scoped>
.my-article-footer {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

.reject-tip {
  margin-top: 12px;
  color: #fecaca;
  font-size: 0.92rem;
}

@media (max-width: 768px) {
  .my-article-footer {
    flex-direction: column;
    gap: 10px;
  }

  .my-article-footer a {
    width: 100%;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .article-card--mine {
    padding: 18px;
  }

  .article-card--mine h3 {
    font-size: 1.1rem;
  }
}
</style>
