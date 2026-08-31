<template>
  <section class="content-section articles-page">
    <section class="articles-hero">
      <div class="articles-hero__copy">
        <p class="eyebrow articles-hero__eyebrow">文章列表</p>
        <h2 class="articles-hero__title">把公开写作整理成可以被检索、筛选和持续阅读的内容库。</h2>
        <p class="articles-hero__text">
          这里适合放学习笔记、项目复盘、技术总结和方法沉淀。读者可以按主题、标签、作者和排序方式快速进入内容。
        </p>

        <div class="articles-hero__actions">
          <router-link class="ghost-btn articles-hero__ghost" to="/my-articles">我的文章</router-link>
          <router-link class="solid-btn" to="/editor">写新文章</router-link>
        </div>
      </div>

      <div class="articles-hero__stats">
        <article class="articles-hero__stat">
          <strong>{{ total }}</strong>
          <span>匹配文章</span>
        </article>
        <article class="articles-hero__stat">
          <strong>{{ categories.length }}</strong>
          <span>内容分类</span>
        </article>
        <article class="articles-hero__stat">
          <strong>{{ visibleTags.length }}</strong>
          <span>可用标签</span>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">推荐阅读</p>
        </div>
        <button class="ghost-btn articles-show-all-btn" @click="toggleAllArticles">
          <span>所有文章</span>
          <span class="articles-show-all-btn__arrow" :class="{ 'articles-show-all-btn__arrow--open': showAll }">&#8594;</span>
        </button>
      </div>

      <div
        v-if="trendingArticles.length && !keyword && !categoryId && !selectedTag && !authorId && page === 1"
        class="article-grid article-grid--featured"
      >
        <ArticleCard v-for="item in trendingArticles" :key="`trending-${item.id}`" :item="item" :keyword="keyword" />
      </div>

      <div v-if="showAll" ref="allArticlesPanel" class="all-articles-panel">
        <div class="all-articles-filters">
          <div class="all-articles-filters__row">
            <select v-model="categoryId" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
              <option value="">全部分类</option>
              <option v-for="item in categories" :key="item.id" :value="String(item.id)">{{ item.name }}</option>
            </select>
            <select v-model="sortBy" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
              <option value="latest">最新发布</option>
              <option value="popular">最多阅读</option>
              <option value="oldest">最早发布</option>
            </select>
            <input
              v-model.trim="keyword"
              class="field-input filter-search"
              placeholder="搜索标题或摘要"
              @keyup.enter="syncRouteWithPage(1)"
            />
            <button class="ghost-btn" @click="resetFilters">重置</button>
          </div>

          <div v-if="authorId" class="all-articles-filters__row">
            <span class="tag-chip">作者过滤：{{ authorName || `#${authorId}` }}</span>
            <button class="ghost-btn" @click="clearAuthorFilter">清除作者</button>
          </div>

          <div v-if="visibleTags.length" class="tag-row tag-row--article-filters">
            <button
              v-for="item in visibleTags"
              :key="item"
              class="tag-chip tag-chip--button"
              :class="{ 'tag-chip--active': selectedTag === item }"
              @click="toggleTag(item)"
            >
              # {{ item }}
            </button>
          </div>
        </div>

        <div class="section-head">
          <div>
            <p class="eyebrow">内容库</p>
            <h3>浏览全部公开文章</h3>
          </div>
        </div>

        <div class="article-grid article-grid--list">
          <ArticleCard v-for="item in articles" :key="item.id" :item="item" :keyword="keyword" />
        </div>

        <div v-if="!articles.length" class="empty-panel">
          <h4>暂时没有匹配的文章</h4>
          <p>可以换一个关键词、分类或标签，或者先去发布一篇新文章。</p>
        </div>

        <div class="pager-row">
          <button class="ghost-btn" :disabled="page <= 1" @click="goToPage(page - 1)">上一页</button>
          <span>第 {{ page }} / {{ totalPages }} 页</span>
          <button class="ghost-btn" :disabled="page >= totalPages" @click="goToPage(page + 1)">下一页</button>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import ArticleCard from "../components/ArticleCard.vue";
import { listArticles, listTrendingArticles } from "../api/article";
import { getMetadata } from "../api/meta";

const route = useRoute();
const router = useRouter();

const articles = ref([]);
const trendingArticles = ref([]);
const categories = ref([]);
const tags = ref([]);
const keyword = ref("");
const categoryId = ref("");
const selectedTag = ref("");
const sortBy = ref("latest");
const authorId = ref("");
const authorName = ref("");
const page = ref(1);
const pageSize = 9;
const total = ref(0);

const showAll = ref(false);
const allArticlesPanel = ref(null);

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));
const visibleTags = computed(() => {
  const fromMetadata = tags.value.map((item) => item.name).filter(Boolean);
  if (fromMetadata.length) return fromMetadata.slice(0, 12);

  const set = new Set();
  [...articles.value, ...trendingArticles.value].forEach((item) => {
    (item.tags || []).forEach((tag) => {
      if (tag?.name) set.add(tag.name);
    });
  });
  return Array.from(set).slice(0, 12);
});

function applyRouteQuery() {
  keyword.value = typeof route.query.keyword === "string" ? route.query.keyword : "";
  categoryId.value = typeof route.query.categoryId === "string" ? route.query.categoryId : "";
  selectedTag.value = typeof route.query.tag === "string" ? route.query.tag : "";
  sortBy.value = typeof route.query.sort === "string" ? route.query.sort : "latest";
  authorId.value = typeof route.query.authorId === "string" ? route.query.authorId : "";
  authorName.value = typeof route.query.authorName === "string" ? route.query.authorName : "";
  const nextPage = Number(route.query.page || 1);
  page.value = Number.isFinite(nextPage) && nextPage > 0 ? nextPage : 1;
}

function buildQuery(nextPage = page.value) {
  const query = {};
  if (keyword.value) query.keyword = keyword.value;
  if (categoryId.value) query.categoryId = categoryId.value;
  if (selectedTag.value) query.tag = selectedTag.value;
  if (sortBy.value && sortBy.value !== "latest") query.sort = sortBy.value;
  if (authorId.value) query.authorId = authorId.value;
  if (authorName.value) query.authorName = authorName.value;
  if (nextPage > 1) query.page = String(nextPage);
  return query;
}

async function loadArticles() {
  const { data } = await listArticles({
    keyword: keyword.value,
    page: page.value,
    pageSize,
    categoryId: categoryId.value,
    tag: selectedTag.value,
    authorId: authorId.value,
    sort: sortBy.value
  });
  articles.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function loadTrendingArticles() {
  const { data } = await listTrendingArticles();
  trendingArticles.value = data.items || [];
}

async function loadMetadata() {
  const { data } = await getMetadata();
  categories.value = data.categories || [];
  tags.value = data.tags || [];
}

function syncRouteWithPage(nextPage = 1) {
  router.replace({ query: buildQuery(nextPage) });
}

function goToPage(nextPage) {
  syncRouteWithPage(nextPage);
}

function toggleTag(tag) {
  selectedTag.value = selectedTag.value === tag ? "" : tag;
  syncRouteWithPage(1);
}

function clearAuthorFilter() {
  authorId.value = "";
  authorName.value = "";
  syncRouteWithPage(1);
}

function resetFilters() {
  keyword.value = "";
  categoryId.value = "";
  selectedTag.value = "";
  sortBy.value = "latest";
  authorId.value = "";
  authorName.value = "";
  syncRouteWithPage(1);
}

function toggleAllArticles() {
  showAll.value = !showAll.value;
  if (showAll.value) {
    nextTick(() => {
      allArticlesPanel.value?.scrollIntoView({ behavior: "smooth", block: "start" });
    });
  }
}

watch(
  () => route.query,
  async () => {
    applyRouteQuery();
    await loadArticles();
  },
  { immediate: true }
);

watch(
  [keyword, categoryId, selectedTag, authorId, page],
  () => {
    if (keyword.value || categoryId.value || selectedTag.value || authorId.value || page.value > 1) {
      showAll.value = true;
    }
  },
  { immediate: true }
);

onMounted(async () => {
  await Promise.all([loadMetadata(), loadTrendingArticles()]);
});
</script>

<style scoped>
.articles-page {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.content-section--compact {
  padding-top: 0;
}

.articles-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  padding: clamp(40px, 8vw, 80px) clamp(24px, 5vw, 60px);
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(260px, 0.75fr);
  gap: 32px;
  align-items: center;
}

.articles-hero__copy {
  display: grid;
  gap: 16px;
}

.articles-hero__eyebrow {
  color: #b4530a;
  font-weight: 600;
}

.articles-hero__title {
  margin: 0;
  font-size: clamp(1.8rem, 3.5vw, 3rem);
  line-height: 1.08;
  color: #1d1d1f;
}

.articles-hero__text {
  margin: 0;
  font-size: clamp(0.95rem, 1.4vw, 1.15rem);
  line-height: 1.6;
  color: #48484a;
}

.articles-hero__actions {
  margin-top: 8px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.articles-hero__ghost {
  color: #374151;
  border-color: rgba(0, 0, 0, 0.15);
  background: transparent;
}

.articles-hero__ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}

.articles-hero__stats {
  display: grid;
  gap: 14px;
}

.articles-hero__stat {
  padding: 18px 20px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
}

.articles-hero__stat strong {
  display: block;
  font-size: 1.6rem;
  margin-bottom: 4px;
  color: #1d1d1f;
}

.articles-hero__stat span {
  color: #6b7280;
  font-size: 0.92rem;
}

.tag-chip--button {
  border: none;
  cursor: pointer;
}

.tag-chip--active {
  background: #f97316;
  color: #fff7ed;
  border-color: rgba(249, 115, 22, 0.45);
}

.article-grid--featured {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-auto-rows: 1fr;
}

.article-grid--featured :deep(.article-card) {
  min-height: 280px;
  display: flex;
  flex-direction: column;
}

.article-grid--featured :deep(.article-card:nth-child(1)) {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.article-grid--featured :deep(.article-card:nth-child(1)::before) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.article-grid--featured :deep(.article-card:nth-child(1) h3) {
  color: #1d1d1f;
}

.article-grid--featured :deep(.article-card:nth-child(1) p) {
  color: #48484a;
}

.article-grid--featured :deep(.article-card:nth-child(1) footer),
.article-grid--featured :deep(.article-card:nth-child(1) .article-card__meta) {
  color: #6b7280;
}

.article-grid--featured :deep(.article-card:nth-child(1) .article-card__more) {
  color: #b4530a;
}

.article-grid--featured :deep(.article-card:nth-child(1) .tag-chip) {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.article-grid--featured :deep(.article-card:nth-child(1) .status-chip) {
  background: rgba(0, 0, 0, 0.06);
  color: #374151;
}

.article-grid--featured :deep(.article-card:nth-child(1) .status-chip.published) {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.article-grid--featured :deep(.article-card:nth-child(1) .article-card__shine) {
  background: radial-gradient(circle, rgba(255, 138, 76, 0.15), transparent 68%);
}

.articles-show-all-btn {
  gap: 6px;
  align-items: center;
  border-radius: 999px;
}

.articles-show-all-btn__arrow {
  display: inline-block;
  transition: transform 0.3s ease;
  font-size: 1.1em;
}

.articles-show-all-btn__arrow--open {
  transform: rotate(90deg);
}

.all-articles-panel {
  margin-top: 24px;
  padding: 28px;
  border-radius: 28px;
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(14px);
  display: flex;
  flex-direction: column;
  gap: 20px;
  animation: panelSlideIn 0.3s ease;
}

.all-articles-panel :deep(.article-card) {
  display: flex;
  flex-direction: column;
  min-height: 175px;
}

.all-articles-filters {
  padding: 20px;
  border-radius: 22px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.all-articles-filters__row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.filter-select-home {
  width: auto;
  min-width: 150px;
}

.filter-search {
  flex: 1;
  min-width: 180px;
}

.article-grid--list {
  grid-template-columns: repeat(3, 1fr);
  grid-auto-rows: 1fr;
}

.article-grid--list :deep(.article-card) {
  display: flex;
  flex-direction: column;
  min-height: 175px;
}

.article-grid--list :deep(.article-card:nth-child(2n)) {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.article-grid--list :deep(.article-card:nth-child(2n)::before) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.article-grid--list :deep(.article-card:nth-child(2n) h3) {
  color: #1d1d1f;
}

.article-grid--list :deep(.article-card:nth-child(2n) p) {
  color: #48484a;
}

.article-grid--list :deep(.article-card:nth-child(2n) footer),
.article-grid--list :deep(.article-card:nth-child(2n) .article-card__meta) {
  color: #6b7280;
}

.article-grid--list :deep(.article-card:nth-child(2n) .article-card__more) {
  color: #b4530a;
}

.article-grid--list :deep(.article-card:nth-child(2n) .tag-chip) {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.article-grid--list :deep(.article-card:nth-child(2n) .status-chip) {
  background: rgba(0, 0, 0, 0.06);
  color: #374151;
}

.article-grid--list :deep(.article-card:nth-child(2n) .status-chip.published) {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.article-grid--list :deep(.article-card:nth-child(2n) .article-card__shine) {
  background: radial-gradient(circle, rgba(255, 138, 76, 0.15), transparent 68%);
}

@keyframes panelSlideIn {
  from {
    opacity: 0;
    transform: translateY(-12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 1100px) {
  .article-grid--featured {
    grid-template-columns: 1fr;
  }

  .article-grid--list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 900px) {
  .articles-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .articles-hero__title {
    font-size: clamp(1.6rem, 5vw, 2.2rem);
  }

  .all-articles-filters__row {
    flex-direction: column;
  }

  .all-articles-filters__row .field-input,
  .all-articles-filters__row .ghost-btn {
    width: 100%;
  }

  .article-grid--list {
    grid-template-columns: 1fr;
  }

  .all-articles-panel {
    padding: 18px;
    border-radius: 22px;
  }
}

@media (max-width: 480px) {
  .articles-page {
    gap: 18px;
  }

  .articles-hero__stats {
    gap: 10px;
  }

  .articles-hero__title {
    font-size: 1.4rem;
  }
}
</style>
