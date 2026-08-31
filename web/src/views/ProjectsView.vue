<template>
  <section class="content-section projects-page">
    <section class="projects-hero">
      <div class="projects-hero__copy">
        <p class="eyebrow projects-hero__eyebrow">作品集</p>
        <h2 class="projects-hero__title">把已经完成的项目，整理成别人真正看得懂的案例。</h2>
        <p class="projects-hero__text">
          这里不是简单的项目列表，而是一组可以被阅读、被比较、被理解的案例。你做了什么、为什么做、怎么做、最后做成了什么，都应该在这里被看见。
        </p>

        <div class="projects-hero__actions">
          <router-link class="ghost-btn projects-hero__ghost" to="/my-projects">我的项目</router-link>
          <router-link class="solid-btn" to="/project-editor">新建项目</router-link>
        </div>
      </div>

      <div class="projects-hero__stats">
        <article class="projects-hero__stat">
          <strong>{{ total }}</strong>
          <span>已发布项目</span>
        </article>
        <article class="projects-hero__stat">
          <strong>{{ featuredCount }}</strong>
          <span>精选案例</span>
        </article>
        <article class="projects-hero__stat">
          <strong>{{ visibleStacks.length }}</strong>
          <span>常见技术栈</span>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">精选案例</p>
        </div>
        <button class="ghost-btn projects-show-all-btn" @click="toggleAllProjects">
          <span>所有项目</span>
          <span class="projects-show-all-btn__arrow" :class="{ 'projects-show-all-btn__arrow--open': showAll }">&#8594;</span>
        </button>
      </div>

      <div
        v-if="featuredProjects.length && !featuredOnly && !stackFilter && !keyword && !authorId && page === 1"
        class="project-grid project-grid--featured"
      >
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

      <div v-if="showAll" ref="allProjectsPanel" class="all-projects-panel">
        <div class="all-projects-filters">
          <div class="all-projects-filters__row">
            <input
              v-model.trim="keyword"
              class="field-input filter-search"
              placeholder="搜索标题、摘要或技术栈"
              @keyup.enter="syncRouteWithPage(1)"
            />
            <div class="all-projects-filters__toggle-row">
              <button class="ghost-btn" :class="{ active: !featuredOnly }" @click="setFeatured(false)">全部项目</button>
              <button class="ghost-btn" :class="{ active: featuredOnly }" @click="setFeatured(true)">只看精选</button>
            </div>
            <select v-model="sortBy" class="field-input filter-select-home" @change="syncRouteWithPage(1)">
              <option value="featured">推荐排序</option>
              <option value="latest">最新发布</option>
              <option value="oldest">最早发布</option>
            </select>
            <button class="ghost-btn" @click="resetFilters">重置</button>
          </div>

          <div v-if="authorId" class="all-projects-filters__row">
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
        </div>

        <div class="section-head">
          <div>
            <p class="eyebrow">案例库</p>
            <h3>按作品集方式浏览全部项目</h3>
          </div>
        </div>

        <div class="project-grid project-grid--list">
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
      </div>
    </section>
  </section>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from "vue";
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

const showAll = ref(false);
const allProjectsPanel = ref(null);

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
  return Array.from(stackSet).slice(0, 12);
});

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

function toggleAllProjects() {
  showAll.value = !showAll.value;
  if (showAll.value) {
    nextTick(() => {
      allProjectsPanel.value?.scrollIntoView({ behavior: "smooth", block: "start" });
    });
  }
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

watch(
  [keyword, featuredOnly, stackFilter, authorId, page],
  () => {
    if (keyword.value || featuredOnly.value || stackFilter.value || authorId.value || page.value > 1) {
      showAll.value = true;
    }
  },
  { immediate: true }
);

onMounted(async () => {
  await loadFeaturedProjects();
});
</script>

<style scoped>
.projects-page {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.content-section--compact {
  padding-top: 0;
}

.projects-hero {
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

.projects-hero__copy {
  display: grid;
  gap: 16px;
}

.projects-hero__eyebrow {
  color: #b4530a;
  font-weight: 600;
}

.projects-hero__title {
  margin: 0;
  font-size: clamp(1.8rem, 3.5vw, 3rem);
  line-height: 1.08;
  color: #1d1d1f;
}

.projects-hero__text {
  margin: 0;
  font-size: clamp(0.95rem, 1.4vw, 1.15rem);
  line-height: 1.6;
  color: #48484a;
}

.projects-hero__actions {
  margin-top: 8px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.projects-hero__ghost {
  color: #374151;
  border-color: rgba(0, 0, 0, 0.15);
  background: transparent;
}

.projects-hero__ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}

.projects-hero__stats {
  display: grid;
  gap: 14px;
}

.projects-hero__stat {
  padding: 18px 20px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
}

.projects-hero__stat strong {
  display: block;
  font-size: 1.6rem;
  margin-bottom: 4px;
  color: #1d1d1f;
}

.projects-hero__stat span {
  color: #6b7280;
  font-size: 0.92rem;
}

.projects-show-all-btn {
  gap: 6px;
  align-items: center;
  border-radius: 999px;
}

.projects-show-all-btn__arrow {
  display: inline-block;
  transition: transform 0.3s ease;
  font-size: 1.1em;
}

.projects-show-all-btn__arrow--open {
  transform: rotate(90deg);
}

.all-projects-panel {
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

.all-projects-filters {
  padding: 20px;
  border-radius: 22px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.all-projects-filters__row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.all-projects-filters__toggle-row {
  display: flex;
  gap: 8px;
}

.filter-search {
  flex: 1;
  min-width: 180px;
}

.filter-select-home {
  width: auto;
  min-width: 150px;
}

.tag-row--project-filters {
  margin-top: 4px;
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
  grid-auto-rows: 1fr;
}

.project-card--featured {
  border-color: rgba(249, 115, 22, 0.28);
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 32%),
    rgba(255, 255, 255, 0.06);
}

.project-grid--list {
  grid-template-columns: repeat(3, 1fr);
  grid-auto-rows: 1fr;
}

.project-card--library {
  min-height: 100%;
}

.project-grid--featured :deep(.project-card:first-child),
.project-grid--list :deep(.project-card:first-child) {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.project-grid--featured :deep(.project-card:first-child::before),
.project-grid--list :deep(.project-card:first-child::before) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.project-grid--featured :deep(.project-card:first-child h3),
.project-grid--list :deep(.project-card:first-child h3) {
  color: #1d1d1f;
}

.project-grid--featured :deep(.project-card:first-child p),
.project-grid--featured :deep(.project-card:first-child .detail-summary),
.project-grid--list :deep(.project-card:first-child p),
.project-grid--list :deep(.project-card:first-child .detail-summary) {
  color: #48484a;
}

.project-grid--featured :deep(.project-card:first-child .table-note),
.project-grid--featured :deep(.project-card:first-child .project-meta-list),
.project-grid--list :deep(.project-card:first-child .table-note),
.project-grid--list :deep(.project-card:first-child .project-meta-list) {
  color: #6b7280;
}

.project-grid--featured :deep(.project-card:first-child .project-link-row .solid-btn),
.project-grid--list :deep(.project-card:first-child .project-link-row .solid-btn) {
  color: #b4530a;
}

.project-grid--featured :deep(.project-card:first-child .tag-chip),
.project-grid--list :deep(.project-card:first-child .tag-chip) {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.project-meta-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 0.92rem;
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
  .project-grid--featured {
    grid-template-columns: 1fr;
  }

  .project-grid--list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 900px) {
  .projects-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .projects-hero__title {
    font-size: clamp(1.6rem, 5vw, 2.2rem);
  }

  .all-projects-filters__row {
    flex-direction: column;
  }

  .all-projects-filters__row .field-input,
  .all-projects-filters__row .ghost-btn {
    width: 100%;
  }

  .all-projects-filters__toggle-row {
    width: 100%;
  }

  .all-projects-filters__toggle-row .ghost-btn {
    flex: 1;
  }

  .project-grid--list {
    grid-template-columns: 1fr;
  }

  .all-projects-panel {
    padding: 18px;
    border-radius: 22px;
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

  .projects-hero__stats {
    gap: 10px;
  }

  .projects-hero__title {
    font-size: 1.4rem;
  }
}
</style>
