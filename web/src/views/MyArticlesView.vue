<template>
  <section class="uc-articles">
    <section class="uc-hero uc-hero--articles">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">创作空间</p>
        <h2 class="uc-hero__title">我的文章</h2>
        <p class="uc-hero__text">统一管理草稿、已发布和审核中的内容，持续输出有价值的文字。</p>
      </div>
      <div class="uc-hero__stats">
        <article class="uc-hero__stat">
          <strong>{{ total }}</strong>
          <span>全部文章</span>
        </article>
      </div>
    </section>

    <div v-if="errorMessage" class="error-text" style="margin-bottom:12px">{{ errorMessage }}</div>

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

  ctx.strokeStyle = "rgba(0,0,0,0.06)";
  ctx.lineWidth = 1;
  for (let i = 0; i <= 4; i++) {
    const y = pad.top + (chartH / 4) * i;
    ctx.beginPath();
    ctx.moveTo(pad.left, y);
    ctx.lineTo(w - pad.right, y);
    ctx.stroke();
    ctx.fillStyle = "rgba(0,0,0,0.3)";
    ctx.font = "10px sans-serif";
    ctx.textAlign = "right";
    ctx.fillText(String(Math.round(max - (max / 4) * i)), pad.left - 4, y + 4);
  }

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

  for (const p of points) {
    ctx.beginPath();
    ctx.arc(p.x, p.y, 3, 0, Math.PI * 2);
    ctx.fillStyle = "#ff8a4c";
    ctx.fill();
  }

  ctx.fillStyle = "rgba(0,0,0,0.3)";
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
.uc-articles {
  display: flex;
  flex-direction: column;
  gap: 24px;
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
  grid-template-columns: minmax(0, 1.4fr) minmax(200px, 0.6fr);
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
  display: flex;
  gap: 12px;
}

.uc-hero__stat {
  padding: 22px 28px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
  text-align: center;
  flex: 1;
}

.uc-hero__stat strong {
  display: block;
  font-size: 2rem;
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

.article-grid :deep(.article-card:first-child) {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.article-grid :deep(.article-card:first-child::before) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.article-grid :deep(.article-card:first-child h3) {
  color: #1d1d1f;
}

.article-grid :deep(.article-card:first-child p) {
  color: #48484a;
}

.article-grid :deep(.article-card:first-child .article-card__meta),
.article-grid :deep(.article-card:first-child footer) {
  color: #6b7280;
}

.article-grid :deep(.article-card:first-child .tag-chip) {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.article-grid :deep(.article-card:first-child .status-chip) {
  background: rgba(0, 0, 0, 0.06);
  color: #374151;
}

.article-grid :deep(.article-card:first-child .status-chip.published) {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.article-grid :deep(.article-card:first-child .article-card__shine) {
  background: radial-gradient(circle, rgba(255, 138, 76, 0.15), transparent 68%);
}

.article-grid {
  grid-template-columns: repeat(2, 1fr);
}

.article-grid :deep(.article-card) {
  min-height: 175px;
}

.my-article-footer {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

.reject-tip {
  margin-top: 12px;
  color: #dc2626;
  font-size: 0.92rem;
}

.privacy-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  background: rgba(255, 217, 142, 0.16);
  color: #92400e;
}

.stats-chart-box {
  margin: 0 0 20px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 16px;
  padding: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.stats-chart-label {
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  margin-bottom: 12px;
}

.stats-canvas {
  width: 100%;
  height: 160px;
  display: block;
}

.delete-link {
  color: #dc2626;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
}

.delete-link:hover {
  color: #b91c1c;
  text-decoration: underline;
}

@media (max-width: 768px) {
  .uc-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }

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
