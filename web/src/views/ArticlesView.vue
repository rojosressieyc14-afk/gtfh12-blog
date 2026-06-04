<template>
  <section class="content-section articles-page">
    <section class="portfolio-stage panel-card portfolio-stage--articles">
      <div class="portfolio-stage__copy">
        <p class="eyebrow">文章列表</p>
        <h2>把公开写作整理成可以被检索、筛选和持续阅读的内容库。</h2>
        <p class="detail-summary">
          这里适合放学习笔记、项目复盘、技术总结和方法沉淀。读者可以按主题、标签、作者和排序方式快速进入内容。
        </p>

        <div class="inline-actions">
          <router-link class="ghost-btn" to="/my-articles">我的文章</router-link>
          <router-link class="solid-btn" to="/editor">写新文章</router-link>
        </div>
      </div>

      <div class="portfolio-stage__stats">
        <article class="project-fact">
          <strong>{{ total }}</strong>
          <span>匹配文章</span>
        </article>
        <article class="project-fact">
          <strong>{{ categories.length }}</strong>
          <span>内容分类</span>
        </article>
        <article class="project-fact">
          <strong>{{ visibleTags.length }}</strong>
          <span>可用标签</span>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">筛选</p>
          <h3>文章检索</h3>
        </div>
        <input
          v-model.trim="keyword"
          class="field-input search-input"
          placeholder="搜索标题或摘要"
          @keyup.enter="syncRouteWithPage(1)"
        />
      </div>

      <div class="filter-row filter-row--articles">
        <select v-model="categoryId" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
          <option value="">全部分类</option>
          <option v-for="item in categories" :key="item.id" :value="String(item.id)">{{ item.name }}</option>
        </select>
        <select v-model="sortBy" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
          <option value="latest">最新发布</option>
          <option value="popular">最多阅读</option>
          <option value="oldest">最早发布</option>
        </select>
        <button class="ghost-btn" @click="resetFilters">重置</button>
      </div>

      <div v-if="authorId" class="filter-row filter-row--articles">
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
    </section>

    <section
      v-if="trendingArticles.length && !keyword && !categoryId && !selectedTag && !authorId && page === 1"
      class="content-section content-section--compact"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">推荐阅读</p>
          <h3>适合放在内容首页的文章</h3>
        </div>
      </div>

      <div class="article-grid article-grid--featured">
        <ArticleCard v-for="item in trendingArticles" :key="`trending-${item.id}`" :item="item" />
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">内容库</p>
          <h3>浏览全部公开文章</h3>
        </div>
      </div>

      <div class="article-grid">
        <ArticleCard v-for="item in articles" :key="item.id" :item="item" />
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
    </section>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
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

watch(
  () => route.query,
  async () => {
    applyRouteQuery();
    await loadArticles();
  },
  { immediate: true }
);

onMounted(async () => {
  await Promise.all([loadMetadata(), loadTrendingArticles()]);
});
</script>

<style scoped>
.articles-page {
  gap: 28px;
}

.content-section--compact {
  padding-top: 0;
}

.portfolio-stage--articles {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(280px, 0.75fr);
  gap: 24px;
  align-items: center;
}

.portfolio-stage__copy {
  display: grid;
  gap: 18px;
}

.portfolio-stage__copy h2 {
  margin: 0;
  font-size: clamp(2rem, 4vw, 3.4rem);
  line-height: 1.02;
}

.portfolio-stage__stats {
  display: grid;
  gap: 16px;
}

.filter-row--articles {
  align-items: center;
}

.tag-row--article-filters {
  margin-top: 16px;
}

.tag-chip--button {
  border: none;
  cursor: pointer;
}

.tag-chip--active,
.ghost-btn.active {
  background: #f97316;
  color: #fff7ed;
  border-color: rgba(249, 115, 22, 0.45);
}

.article-grid--featured {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

@media (max-width: 1100px) {
  .article-grid--featured {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .portfolio-stage--articles {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .portfolio-stage__copy h2 {
    font-size: clamp(1.6rem, 5vw, 2.2rem);
  }

  .filter-row--articles {
    flex-direction: column;
  }

  .filter-row--articles .field-input,
  .filter-row--articles .ghost-btn {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .articles-page {
    gap: 18px;
  }

  .portfolio-stage__stats {
    gap: 10px;
  }

  .portfolio-stage__copy h2 {
    font-size: 1.4rem;
  }
}
</style>
