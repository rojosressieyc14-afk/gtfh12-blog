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

    <p v-if="errorMessage" class="error-text" style="margin-bottom:12px">{{ errorMessage }}</p>

    <div v-if="articles.length" class="article-grid">
      <article v-for="item in articles" :key="item.id" class="article-card article-card--mine">
        <div class="article-card__meta">
          <span class="status-chip" :class="item.status">{{ labelMap[item.status] || "未知状态" }}</span>
          <span v-if="item.isPrivate" class="privacy-badge">私有</span>
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
          <button class="inline-link" @click="toggleStats(item.id)">统计</button>
          <button class="inline-link delete-link" @click="handleDelete(item.id, item.title)">删除</button>
        </footer>

        <p v-if="item.rejectReason" class="reject-tip">驳回原因：{{ item.rejectReason }}</p>
      </article>
      <div v-if="statsArticleId === item.id && statsData.length" class="stats-chart-box">
        <p class="stats-chart-label">近 30 天阅读趋势</p>
        <canvas ref="statsCanvas" class="stats-canvas"></canvas>
      </div>
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
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { deleteArticle, getArticleStats, listMyArticles } from "../api/article";

const articles = ref([]);
const page = ref(1);
const pageSize = 8;
const total = ref(0);

const statsArticleId = ref(null);
const statsData = ref([]);
const statsCanvas = ref(null);
const errorMessage = ref("");

async function toggleStats(id) {
  if (statsArticleId.value === id) {
    statsArticleId.value = null;
    statsData.value = [];
    return;
  }
  statsArticleId.value = id;
  try {
    const { data } = await getArticleStats(id);
    statsData.value = data.items || [];
  } catch {
    statsData.value = [];
  }
}

function drawChart() {
  if (!statsCanvas.value || !statsData.value.length) return;
  const canvas = statsCanvas.value;
  const ctx = canvas.getContext("2d");
  const dpr = window.devicePixelRatio || 1;
  const rect = canvas.getBoundingClientRect();
  canvas.width = rect.width * dpr;
  canvas.height = rect.height * dpr;
  ctx.scale(dpr, dpr);
  const w = rect.width;
  const h = rect.height;
  const pad = { top: 12, right: 12, bottom: 24, left: 36 };
  const chartW = w - pad.left - pad.right;
  const chartH = h - pad.top - pad.bottom;

  const values = statsData.value.map((d) => d.count);
  const max = Math.max(...values, 1);
  const points = values.map((v, i) => ({
    x: pad.left + (i / Math.max(values.length - 1, 1)) * chartW,
    y: pad.top + chartH - (v / max) * chartH,
  }));

  ctx.clearRect(0, 0, w, h);

  // Grid lines
  ctx.strokeStyle = "rgba(255,255,255,0.06)";
  ctx.lineWidth = 1;
  for (let i = 0; i <= 4; i++) {
    const y = pad.top + (chartH / 4) * i;
    ctx.beginPath();
    ctx.moveTo(pad.left, y);
    ctx.lineTo(w - pad.right, y);
    ctx.stroke();
    ctx.fillStyle = "rgba(246,241,234,0.3)";
    ctx.font = "10px sans-serif";
    ctx.textAlign = "right";
    ctx.fillText(String(Math.round(max - (max / 4) * i)), pad.left - 4, y + 4);
  }

  // Line
  if (points.length < 2) return;
  ctx.beginPath();
  ctx.strokeStyle = "#ff8a4c";
  ctx.lineWidth = 2;
  ctx.lineJoin = "round";
  ctx.moveTo(points[0].x, points[0].y);
  for (let i = 1; i < points.length; i++) {
    ctx.lineTo(points[i].x, points[i].y);
  }
  ctx.stroke();

  // Dots
  for (const p of points) {
    ctx.beginPath();
    ctx.arc(p.x, p.y, 3, 0, Math.PI * 2);
    ctx.fillStyle = "#ff8a4c";
    ctx.fill();
  }

  // Date labels
  ctx.fillStyle = "rgba(246,241,234,0.3)";
  ctx.font = "10px sans-serif";
  ctx.textAlign = "center";
  const labels = statsData.value.map((d) => d.date.slice(5));
  const step = Math.max(1, Math.floor(labels.length / 6));
  for (let i = 0; i < labels.length; i += step) {
    ctx.fillText(labels[i], points[i].x, h - 6);
  }
}

watch(statsData, () => nextTick(drawChart), { flush: "post" });

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
    errorMessage.value = "删除失败：" + (e.response?.data?.message || e.message);
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

.privacy-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  background: rgba(255, 217, 142, 0.16);
  color: #ffd98e;
}

.stats-chart-box {
  margin: 0 0 20px;
  background: rgba(255,255,255,0.03);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid var(--border);
}
.stats-chart-label {
  font-size: 12px;
  font-weight: 600;
  color: rgba(246,241,234,0.5);
  margin-bottom: 12px;
}
.stats-canvas {
  width: 100%;
  height: 160px;
  display: block;
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
