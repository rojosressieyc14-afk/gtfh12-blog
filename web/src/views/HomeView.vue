<template>
  <section class="hero-panel hero-panel--portfolio home-hero">
    <div class="hero-copy home-hero__copy">
      <p class="eyebrow">个人作品集平台</p>
      <h2>{{ heroTitle }}</h2>
      <p class="hero-text">{{ heroText }}</p>

      <div class="hero-actions">
        <router-link class="solid-btn" :to="userStore.isLoggedIn ? '/editor' : '/auth'">
          {{ userStore.isLoggedIn ? "开始写作" : "注册开始" }}
        </router-link>
        <router-link class="ghost-btn" :to="articlesLibraryLink">文章库</router-link>
        <router-link class="ghost-btn" :to="browseProjectsLink">项目库</router-link>
        <router-link class="ghost-btn" :to="featuredAuthorLink">作者页</router-link>
      </div>

      <div class="hero-metrics">
        <article class="hero-metric">
          <strong>{{ total }}</strong>
          <span>公开文章</span>
        </article>
        <article class="hero-metric">
          <strong>{{ featuredProjects.length }}</strong>
          <span>精选项目</span>
        </article>
        <article class="hero-metric">
          <strong>{{ tags.length }}</strong>
          <span>内容标签</span>
        </article>
      </div>

      <div v-if="siteOwner.skills?.length" class="tag-row tag-row--hero">
        <span v-for="skill in siteOwner.skills.slice(0, 6)" :key="skill" class="tag-chip">{{ skill }}</span>
      </div>
    </div>

    <div class="hero-orbit hero-orbit--portfolio">
      <article class="orbit-card">
        <strong>项目案例</strong>
        <p>已完成工作的可读、可比较、可证明的公开案例。</p>
      </article>
      <article class="orbit-card">
        <strong>学习沉淀</strong>
        <p>笔记、总结和复盘的长期归档，持续增长的内容资产。</p>
      </article>
      <article class="orbit-card">
        <strong>公开表达</strong>
        <p>项目、文章和个人品牌汇聚在同一入口的能力画像。</p>
      </article>
    </div>
  </section>

  <section class="portfolio-intro">
    <article class="panel-card portfolio-card portfolio-card--wide">
      <p class="eyebrow">定位</p>
      <h3>项目案例、学习笔记与公开输出的长期作品展示面。</h3>
    </article>

    <article class="panel-card portfolio-card">
      <p class="eyebrow">项目</p>
      <h3>展示你真正做过什么</h3>
      <router-link class="ghost-btn" :to="ownerProjectsLink">查看项目案例</router-link>
    </article>

    <article class="panel-card portfolio-card">
      <p class="eyebrow">学习</p>
      <h3>沉淀你的成长轨迹</h3>
      <router-link class="ghost-btn" :to="studyArticlesLink">查看学习内容</router-link>
    </article>
  </section>

  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">精选项目</p>
        <h3>精选项目</h3>
      </div>
      <router-link class="ghost-btn" :to="browseProjectsLink">查看完整作品集</router-link>
    </div>

    <div v-if="featuredProjects.length" class="project-grid project-grid--hero">
      <article v-for="item in featuredProjects" :key="`project-${item.id}`" class="project-card panel-card project-card--hero">
        <div v-if="item.coverImage" class="project-card__cover">
          <img :src="toAssetUrl(item.coverImage)" :alt="item.title" />
        </div>
        <div class="project-card__head">
          <div>
            <h3>{{ item.title }}</h3>
            <p class="table-note">{{ item.summary || "这个项目暂时还没有摘要。" }}</p>
          </div>
          <span v-if="item.isFeatured" class="tag-chip">精选</span>
        </div>
        <div v-if="item.techStacks?.length" class="tag-row">
          <span v-for="stack in item.techStacks.slice(0, 4)" :key="stack" class="tag-chip"># {{ stack }}</span>
        </div>
        <div class="project-link-row">
          <a v-if="item.demoUrl" class="ghost-btn" :href="item.demoUrl" target="_blank" rel="noreferrer">在线演示</a>
          <a v-if="item.repoUrl" class="ghost-btn" :href="item.repoUrl" target="_blank" rel="noreferrer">代码仓库</a>
          <router-link class="solid-btn" :to="`/projects/${item.id}`">查看案例</router-link>
        </div>
      </article>
    </div>
    <div v-else class="empty-panel">
      <h4>还没有已发布项目</h4>
      <p>发布代表性项目后这里会展示你的作品集。</p>
    </div>
  </section>

  <section v-if="recommendedAuthors.length" class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">推荐作者</p>
        <h3>平台上的活跃创作者</h3>
      </div>
    </div>

    <div class="author-grid">
      <article v-for="item in recommendedAuthors" :key="`author-${item.id}`" class="author-card panel-card">
        <div class="author-card__avatar">
          <img v-if="item.avatar" :src="toAssetUrl(item.avatar)" :alt="item.username" />
          <div v-else class="author-card__avatar-placeholder">{{ item.username?.charAt(0) || '?' }}</div>
        </div>
        <div class="author-card__body">
          <h4>{{ item.username }}</h4>
          <p v-if="item.headline" class="author-card__role">{{ item.headline }}</p>
          <p v-else class="author-card__role author-card__role--empty">暂无简介</p>
          <div v-if="item.skills?.length" class="tag-row">
            <span v-for="skill in item.skills.slice(0, 3)" :key="skill" class="tag-chip">{{ skill }}</span>
          </div>
        </div>
        <router-link class="solid-btn" :to="`/author/${item.id}`">查看主页</router-link>
      </article>
    </div>
  </section>

  <section class="content-section content-section--split">
    <div class="content-split-card panel-card">
      <div class="section-head">
        <div>
          <p class="eyebrow">推荐内容</p>
          <h3>推荐内容</h3>
        </div>
        <router-link class="ghost-btn" :to="articlesLibraryLink">全部文章</router-link>
      </div>
      <div class="article-grid article-grid--single">
        <ArticleCard v-for="item in featuredArticles" :key="`featured-${item.id}`" :item="item" />
      </div>
    </div>

    <div class="content-split-card panel-card">
      <div class="section-head">
        <div>
          <p class="eyebrow">学习笔记</p>
          <h3>最近的学习记录</h3>
        </div>
        <router-link class="ghost-btn" :to="studyArticlesLink">只看学习内容</router-link>
      </div>
      <div class="article-grid article-grid--single">
        <ArticleCard v-for="item in learningArticles" :key="`learning-${item.id}`" :item="item" />
      </div>
    </div>
  </section>

  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">探索</p>
        <h3>全部文章</h3>
      </div>
      <router-link class="ghost-btn" :to="articlesLibraryLink">打开文章库</router-link>
    </div>

    <div class="filter-row">
      <input
        v-model="keyword"
        class="field-input search-input"
        placeholder="搜索标题或摘要"
        @keyup.enter="goToPage(1)"
      />
      <select v-model="selectedCategory" class="field-input filter-select-home" @change="goToPage(1)">
        <option value="">全部分类</option>
        <option v-for="item in categories" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <select v-model="selectedTag" class="field-input filter-select-home" @change="goToPage(1)">
        <option value="">全部标签</option>
        <option v-for="item in tags" :key="item.id" :value="item.name">{{ item.name }}</option>
      </select>
      <button class="ghost-btn" @click="resetFilters">重置筛选</button>
    </div>

    <div class="article-grid">
      <ArticleCard v-for="item in articles" :key="item.id" :item="item" />
    </div>

    <div v-if="!articles.length" class="empty-panel">
      <h4>还没有已发布文章</h4>
      <p>写一篇文章来开始你的内容积累。</p>
    </div>

    <div class="pager-row">
      <button class="ghost-btn" :disabled="page <= 1" @click="goToPage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button class="ghost-btn" :disabled="page >= totalPages" @click="goToPage(page + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { listArticles, listTrendingArticles } from "../api/article";
import { getMetadata } from "../api/meta";
import { listProjects } from "../api/project";
import { useUserStore } from "../stores/user";
import ArticleCard from "../components/ArticleCard.vue";
import { toAssetUrl } from "../utils/asset";
import { getAuthorProfile, getRecommendedAuthors } from "../api/profile";

const userStore = useUserStore();
const router = useRouter();
const keyword = ref("");
const articles = ref([]);
const trendingArticles = ref([]);
const featuredProjects = ref([]);
const recommendedAuthors = ref([]);
const categories = ref([]);
const tags = ref([]);
const selectedCategory = ref("");
const selectedTag = ref("");
const page = ref(1);
const pageSize = 9;
const total = ref(0);
const ownerProfile = ref(null);

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));
const siteOwner = computed(() => ownerProfile.value || userStore.profile || {});
const featuredAuthorLink = computed(() => (siteOwner.value?.id ? `/author/${siteOwner.value.id}` : "/auth"));
const browseProjectsLink = computed(() => ({ path: "/projects", query: { sort: "featured" } }));
const ownerProjectsLink = computed(() => ({
  path: "/projects",
  query: siteOwner.value?.id
    ? {
        authorId: String(siteOwner.value.id),
        authorName: siteOwner.value.username || "",
        sort: "featured"
      }
    : { sort: "featured" }
}));
const articlesLibraryLink = computed(() => ({ path: "/articles", query: { sort: "latest" } }));
const studyArticlesLink = computed(() => {
  const matchedTag = tags.value.find((item) => isStudyText(item.name || ""));
  const matchedCategory = categories.value.find((item) => isStudyText(item.name || ""));

  return {
    path: "/articles",
    query: {
      ...(matchedCategory?.id ? { categoryId: String(matchedCategory.id) } : {}),
      ...(matchedTag?.name ? { tag: matchedTag.name } : {}),
      sort: "latest"
    }
  };
});
const heroTitle = computed(() => siteOwner.value?.headline || "把你的博客做成作品集、学习归档和长期表达入口");
const heroText = computed(() => {
  if (siteOwner.value?.bio) return siteOwner.value.bio;
  return "这个平台适合把项目案例、学习笔记和公开写作放在同一个站点里，既能用于求职展示，也能长期沉淀内容资产。";
});

const featuredArticles = computed(() => {
  if (trendingArticles.value.length >= 3) return trendingArticles.value.slice(0, 3);
  return articles.value.slice(0, 3);
});

const learningArticles = computed(() => {
  const source = [...articles.value, ...trendingArticles.value];
  const unique = source.filter((item, index, list) => list.findIndex((entry) => entry.id === item.id) === index);
  const matched = unique.filter((item) => isStudyArticle(item));
  if (matched.length) return matched.slice(0, 3);
  return unique.slice(0, 3);
});

function isStudyText(value) {
  const haystack = String(value).toLowerCase();
  return ["study", "note", "docs", "source", "review", "summary", "学习", "笔记", "源码", "总结", "复盘"].some((word) =>
    haystack.includes(word)
  );
}

function isStudyArticle(item) {
  const haystack = [item.title, item.summary, item.category?.name, ...(item.tags || []).map((tag) => tag.name)]
    .filter(Boolean)
    .join(" ")
    .toLowerCase();

  return isStudyText(haystack);
}

async function fetchArticles() {
  const { data } = await listArticles({
    keyword: keyword.value,
    page: page.value,
    pageSize,
    categoryId: selectedCategory.value,
    tag: selectedTag.value
  });
  articles.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function fetchTrending() {
  const { data } = await listTrendingArticles();
  trendingArticles.value = data.items || [];
}

async function fetchProjects() {
  const { data } = await listProjects({ page: 1, pageSize: 3, featured: true });
  featuredProjects.value = data.items || [];
}

async function fetchRecommendedAuthors() {
  try {
    const { data } = await getRecommendedAuthors();
    recommendedAuthors.value = data.items || [];
  } catch {
    recommendedAuthors.value = [];
  }
}

async function fetchMetadata() {
  const { data } = await getMetadata();
  categories.value = data.categories || [];
  tags.value = data.tags || [];
}

async function fetchOwnerProfile(id) {
  if (!id) {
    ownerProfile.value = null;
    return;
  }
  try {
    const { data } = await getAuthorProfile(id);
    ownerProfile.value = data.user || null;
  } catch {
    ownerProfile.value = userStore.profile || null;
  }
}

async function goToPage(nextPage) {
  page.value = nextPage;
  await fetchArticles();
}

async function resetFilters() {
  keyword.value = "";
  selectedCategory.value = "";
  selectedTag.value = "";
  await goToPage(1);
}

watch(
  () => userStore.profile?.id,
  (id) => {
    fetchOwnerProfile(id);
  },
  { immediate: true }
);

onMounted(async () => {
  await Promise.all([fetchArticles(), fetchTrending(), fetchMetadata(), fetchProjects(), fetchRecommendedAuthors()]);
});
</script>

<style scoped>
.home-hero__copy {
  position: relative;
  z-index: 1;
}

.hero-orbit {
  position: relative;
  min-height: 320px;
}

.orbit-card {
  position: absolute;
  display: grid;
  gap: 10px;
  padding: 18px 20px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.orbit-card strong,
.orbit-card p {
  margin: 0;
}

.orbit-card p {
  color: var(--text-soft);
}

.orbit-card:nth-child(1) {
  top: 28px;
  left: 10%;
}

.orbit-card:nth-child(2) {
  right: 8%;
  top: 128px;
}

.orbit-card:nth-child(3) {
  left: 22%;
  bottom: 18px;
}

.content-section--split {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.content-split-card {
  min-width: 0;
}

.project-grid--hero .project-card--hero {
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 32%),
    rgba(255, 255, 255, 0.06);
}

.author-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.author-card {
  display: grid;
  gap: 14px;
  padding: 20px;
  text-align: center;
  justify-items: center;
}

.author-card__avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
}

.author-card__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-card__avatar-placeholder {
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--soft);
}

.author-card__body {
  display: grid;
  gap: 6px;
}

.author-card__body h4 {
  margin: 0;
  font-size: 1.1rem;
}

.author-card__role {
  color: var(--soft);
  font-size: 0.9rem;
  margin: 0;
}

.author-card__role--empty {
  font-style: italic;
}

.author-card .tag-row {
  justify-content: center;
}

@media (max-width: 960px) {
  .content-section--split {
    grid-template-columns: 1fr;
  }

  .orbit-card {
    position: relative;
    inset: auto;
  }

  .hero-orbit {
    display: grid;
    gap: 14px;
    min-height: auto;
  }
}

@media (max-width: 768px) {
  .author-grid {
    grid-template-columns: 1fr;
  }

  .orbit-card {
    padding: 14px 16px;
  }
}

@media (max-width: 480px) {
  .hero-actions {
    flex-direction: column;
  }

  .hero-actions .solid-btn,
  .hero-actions .ghost-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
