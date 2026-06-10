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
      </article>
      <article class="panel-card service-card">
        <p class="eyebrow">建议</p>
        <h4>优先补齐摘要、标签和分类</h4>
      </article>
    </div>

    <div v-if="articles.length" class="article-grid">
      <article v-for="item in articles" :key="item.id" class="article-card article-card--mine">
        <div class="article-card__meta">
          <span class="status-chip" :class="item.status">{{ labelMap[item.status] || "未知状态" }}</span>
          <span>{{ item.category?.name || "未分类" }} | {{ formatDate(item.updatedAt) }}</span>
        </div>

        <h3>{{ item.title }}</h3>
        <p>{{ item.summary || "暂无摘要" }}</p>

        <div v-if="item.tags?.length" class="tag-row">
          <span v-for="tag in item.tags" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
        </div>

        <footer class="my-article-footer">
          <router-link class="inline-link" :to="`/editor/${item.id}`">继续编辑</router-link>
          <router-link v-if="item.status === 'published'" class="inline-link" :to="`/article/${item.id}`">查看详情</router-link>
          <button class="inline-link delete-link" @click="handleDelete(item.id, item.title)">删除</button>
        </footer>

        <p v-if="item.rejectReason" class="reject-tip">驳回原因：{{ item.rejectReason }}</p>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>你还没有文章</h4>
      <p>创建你的第一篇文章。</p>
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
import { deleteArticle, listMyArticles } from "../api/article";

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

async function handleDelete(id, title) {
  if (!confirm(`确定要删除文章「${title}」吗？此操作不可撤销。`)) return;
  try {
    await deleteArticle(id);
    await loadMine();
  } catch (e) {
    alert("删除失败：" + (e.response?.data?.message || e.message));
  }
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

.delete-link {
  color: #f87171;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
}

.delete-link:hover {
  color: #ef4444;
  text-decoration: underline;
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
