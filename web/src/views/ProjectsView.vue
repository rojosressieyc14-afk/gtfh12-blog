<template>
  <section class="content-section projects-page">
    <section class="portfolio-stage panel-card portfolio-stage--projects">
      <div class="portfolio-stage__copy">
        <p class="eyebrow">作品集</p>
        <h2>把已经完成的项目，整理成别人真正看得懂的案例。</h2>
        <p class="detail-summary">
          这里不是简单的项目列表，而是一组可以被阅读、被比较、被理解的案例。你做了什么、为什么做、怎么做、最后做成了什么，都应该在这里被看见。
        </p>

        <div class="inline-actions">
          <router-link class="ghost-btn" to="/my-projects">我的项目</router-link>
          <router-link class="solid-btn" to="/project-editor">新建项目</router-link>
        </div>
      </div>

      <div class="portfolio-stage__stats">
        <article class="project-fact">
          <strong>{{ total }}</strong>
          <span>已发布项目</span>
        </article>
        <article class="project-fact">
          <strong>{{ featuredCount }}</strong>
          <span>精选案例</span>
        </article>
        <article class="project-fact">
          <strong>{{ visibleStacks.length }}</strong>
          <span>常见技术栈</span>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">筛选</p>
          <h3>案例检索</h3>
        </div>
        <input
          v-model.trim="keyword"
          class="field-input search-input"
          placeholder="搜索标题、摘要或技术栈"
          @keyup.enter="syncRouteWithPage(1)"
        />
      </div>

      <div class="filter-row filter-row--projects">
        <button class="ghost-btn" :class="{ active: !featuredOnly }" @click="setFeatured(false)">全部项目</button>
        <button class="ghost-btn" :class="{ active: featuredOnly }" @click="setFeatured(true)">只看精选</button>
        <select v-model="sortBy" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
          <option value="featured">推荐排序</option>
          <option value="latest">最新发布</option>
          <option value="oldest">最早发布</option>
        </select>
        <button class="ghost-btn" @click="resetFilters">重置</button>
      </div>

      <div v-if="authorId" class="filter-row filter-row--projects">
        <span class="tag-chip">作者过滤：{{ authorName || `#${authorId}` }}</span>
        <button class="ghost-btn" @click="clearAuthorFilter">清除作者</button>
      </div>

      <div v-if="visibleStacks.length" class="tag-row tag-row--project-filters">
        <button
          v-for="stack in visibleStacks"
          :key="stack"
          class="tag-chip tag-chip--button"
          :class="{ 'tag-chip--active': stackFilter === stack }"
          @click="toggleStackFilter(stack)"
        >
          # {{ stack }}
        </button>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="services-grid">
        <article v-for="item in portfolioPrompts" :key="item.title" class="panel-card service-card">
          <p class="eyebrow">{{ item.kicker }}</p>
          <h4>{{ item.title }}</h4>
          <p>{{ item.description }}</p>
        </article>
      </div>
    </section>

    <section v-if="featuredProjects.length && !featuredOnly && page === 1" class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">首页级案例</p>
          <h3>最适合放在第一屏的项目</h3>
        </div>
      </div>

      <div class="project-grid project-grid--featured">
        <article v-for="item in featuredProjects" :key="`featured-${item.id}`" class="project-card panel-card project-card--featured">
          <div v-if="item.coverImage" class="project-card__cover">
            <img :src="toAssetUrl(item.coverImage)" :alt="item.title" />
          </div>
          <div class="project-card__head">
            <div>
              <h3>{{ item.title }}</h3>
              <p class="table-note">{{ item.roleLabel || item.duration || "项目案例" }}</p>
            </div>
            <span class="tag-chip">精选</span>
          </div>
          <p class="detail-summary">{{ item.summary || "这个项目暂时还没有摘要。" }}</p>
          <div v-if="item.techStacks?.length" class="tag-row">
            <span v-for="stack in item.techStacks.slice(0, 4)" :key="stack" class="tag-chip"># {{ stack }}</span>
          </div>
          <div class="project-link-row">
            <a v-if="item.demoUrl" class="ghost-btn" :href="item.demoUrl" target="_blank" rel="noreferrer">在线演示</a>
            <a v-if="item.repoUrl" class="ghost-btn" :href="item.repoUrl" target="_blank" rel="noreferrer">代码仓库</a>
            <router-link class="solid-btn" :to="`/projects/${item.id}`">打开案例</router-link>
          </div>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">案例库</p>
          <h3>按作品集方式浏览全部项目</h3>
        </div>
      </div>

      <div class="project-grid">
        <article v-for="item in projects" :key="item.id" class="project-card panel-card project-card--library">
          <div v-if="item.coverImage" class="project-card__cover">
            <img :src="toAssetUrl(item.coverImage)" :alt="item.title" />
          </div>
          <div class="project-card__head">
            <div>
              <h3>{{ item.title }}</h3>
              <p class="table-note">{{ item.author?.username || "匿名作者" }} · {{ formatDate(item.publishedAt || item.createdAt) }}</p>
            </div>
            <span v-if="item.isFeatured" class="tag-chip">精选</span>
          </div>

          <p class="detail-summary">{{ item.summary || "这个项目暂时还没有摘要。" }}</p>

          <div class="project-meta-list">
            <span v-if="item.roleLabel">{{ item.roleLabel }}</span>
            <span v-if="item.duration">{{ item.duration }}</span>
            <span v-if="item.teamLabel">{{ item.teamLabel }}</span>
          </div>

          <div v-if="item.techStacks?.length" class="tag-row">
            <span v-for="stack in item.techStacks.slice(0, 5)" :key="stack" class="tag-chip"># {{ stack }}</span>
          </div>

          <div v-if="item.highlights?.length" class="project-highlight-list">
            <span v-for="point in item.highlights.slice(0, 3)" :key="point">· {{ point }}</span>
          </div>

          <div class="project-link-row">
            <a v-if="item.demoUrl" class="ghost-btn" :href="item.demoUrl" target="_blank" rel="noreferrer">在线演示</a>
            <a v-if="item.repoUrl" class="ghost-btn" :href="item.repoUrl" target="_blank" rel="noreferrer">代码仓库</a>
            <router-link class="solid-btn" :to="`/projects/${item.id}`">查看详情</router-link>
          </div>
        </article>
      </div>

      <div v-if="!projects.length" class="empty-panel">
        <h4>暂时没有匹配的项目</h4>
        <p>可以换个关键词或技术栈筛选，或者先去发布更多项目案例。</p>
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
import { listProjects } from "../api/project";
import { toAssetUrl } from "../utils/asset";

const route = useRoute();
const router = useRouter();

const projects = ref([]);
const featuredProjects = ref([]);
const page = ref(1);
const pageSize = 6;
const total = ref(0);
const keyword = ref("");
const featuredOnly = ref(false);
const stackFilter = ref("");
const sortBy = ref("featured");
const authorId = ref("");
const authorName = ref("");

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));
const featuredCount = computed(() => {
  const currentPageFeatured = projects.value.filter((item) => item.isFeatured).length;
  const highlightFeatured = featuredProjects.value.filter((item) => item.isFeatured).length;
  return Math.max(currentPageFeatured, highlightFeatured);
});
const visibleStacks = computed(() => {
  const stackSet = new Set();
  [...projects.value, ...featuredProjects.value].forEach((item) => {
    (item.techStacks || []).forEach((stack) => {
      if (stack) stackSet.add(stack);
    });
  });
  return Array.from(stackSet).slice(0, 8);
});
const portfolioPrompts = computed(() => [
  {
    kicker: "叙事",
    title: "先讲清楚项目价值",
    description: "好的作品集项目，不只是在展示界面，而是能让别人快速理解你为什么做它。"
  },
  {
    kicker: "过程",
    title: "把推进过程讲清楚",
    description: "把需求、方案、权衡和实现过程写出来，案例就会比普通项目列表更有说服力。"
  },
  {
    kicker: "结果",
    title: "让结果可以被外界理解",
    description: "无论是上线效果、体验提升还是结构优化，只要说明白结果，案例就会更完整。"
  }
]);

function applyRouteQuery() {
  keyword.value = typeof route.query.keyword === "string" ? route.query.keyword : "";
  featuredOnly.value = route.query.featured === "true";
  stackFilter.value = typeof route.query.stack === "string" ? route.query.stack : "";
  sortBy.value = typeof route.query.sort === "string" ? route.query.sort : "featured";
  authorId.value = typeof route.query.authorId === "string" ? route.query.authorId : "";
  authorName.value = typeof route.query.authorName === "string" ? route.query.authorName : "";
  const nextPage = Number(route.query.page || 1);
  page.value = Number.isFinite(nextPage) && nextPage > 0 ? nextPage : 1;
}

function buildQuery(nextPage = page.value) {
  const query = {};
  if (keyword.value) query.keyword = keyword.value;
  if (featuredOnly.value) query.featured = "true";
  if (stackFilter.value) query.stack = stackFilter.value;
  if (sortBy.value && sortBy.value !== "featured") query.sort = sortBy.value;
  if (authorId.value) query.authorId = authorId.value;
  if (authorName.value) query.authorName = authorName.value;
  if (nextPage > 1) query.page = String(nextPage);
  return query;
}

async function loadProjects() {
  const { data } = await listProjects({
    keyword: keyword.value,
    page: page.value,
    pageSize,
    featured: featuredOnly.value,
    stack: stackFilter.value,
    authorId: authorId.value,
    sort: sortBy.value
  });
  projects.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function loadFeaturedProjects() {
  const { data } = await listProjects({ page: 1, pageSize: 3, featured: true });
  featuredProjects.value = data.items || [];
}

function syncRouteWithPage(nextPage = 1) {
  router.replace({ query: buildQuery(nextPage) });
}

function goToPage(nextPage) {
  syncRouteWithPage(nextPage);
}

function setFeatured(value) {
  featuredOnly.value = value;
  stackFilter.value = "";
  syncRouteWithPage(1);
}

function resetFilters() {
  keyword.value = "";
  featuredOnly.value = false;
  stackFilter.value = "";
  sortBy.value = "featured";
  authorId.value = "";
  authorName.value = "";
  syncRouteWithPage(1);
}

function toggleStackFilter(stack) {
  stackFilter.value = stackFilter.value === stack ? "" : stack;
  syncRouteWithPage(1);
}

function clearAuthorFilter() {
  authorId.value = "";
  authorName.value = "";
  syncRouteWithPage(1);
}

function formatDate(value) {
  return new Date(value).toLocaleDateString("zh-CN");
}

watch(
  () => route.query,
  async () => {
    applyRouteQuery();
    await loadProjects();
  },
  { immediate: true }
);

onMounted(async () => {
  await loadFeaturedProjects();
});
</script>

<style scoped>
.projects-page {
  gap: 28px;
}

.content-section--compact {
  padding-top: 0;
}

.portfolio-stage--projects {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(280px, 0.8fr);
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

.filter-row--projects {
  align-items: center;
}

.tag-row--project-filters {
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

.project-grid--featured {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.project-card--featured {
  border-color: rgba(249, 115, 22, 0.28);
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 32%),
    rgba(255, 255, 255, 0.06);
}

.project-card--library {
  min-height: 100%;
}

.project-meta-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 0.92rem;
}

@media (max-width: 1100px) {
  .project-grid--featured {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .portfolio-stage--projects {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .portfolio-stage__copy h2 {
    font-size: clamp(1.6rem, 5vw, 2.2rem);
  }

  .filter-row--projects {
    flex-direction: column;
  }

  .filter-row--projects .ghost-btn,
  .filter-row--projects .field-input {
    width: 100%;
  }

  .project-meta-list {
    flex-direction: column;
    gap: 6px;
  }
}

@media (max-width: 480px) {
  .projects-page {
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
