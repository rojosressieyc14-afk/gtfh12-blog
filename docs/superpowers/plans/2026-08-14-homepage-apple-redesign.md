# 首页品牌化改版（苹果式极简，暗色为主+浅色 hero）

Date: 2026-08-14
Status: Draft
Author: opencode

## Overview

按用户确认的方向重构 `web/src/views/HomeView.vue`：

- **内容大幅精简**：只保留「大 hero + 精选项目 + 最近文章 + 收尾 CTA」。删除：hero 指标区、技能标签、hero-orbit 三卡片、定位/项目/学习卡片区、推荐作者区、推荐内容/学习笔记双栏、全部文章列表（含搜索/筛选/分页）。
- **视觉**：整体暗色体系不变；hero 改为浅色大卡片（#f4f4f6 + 橙色径向光晕），苹果式大标题（clamp 2.75rem–5rem）、大留白、居中排版、单一主 CTA + 文本链接。
- 不动 `main.css` 全局样式与其他视图（`hero-panel` 等全局类继续供其他视图使用）。

## 相关文件

- `web/src/views/HomeView.vue`（唯一改动文件）

## Tasks

### T1: template 重写

替换整个 `<template>`（结构如下，项目卡片复用现有 `.project-card` 标记）：

```html
<template>
  <section class="hero-apple">
    <div class="hero-apple__inner">
      <p class="eyebrow hero-apple__eyebrow">PulseBlog</p>
      <h2 class="hero-apple__title">{{ heroTitle }}</h2>
      <p class="hero-apple__text">{{ heroText }}</p>
      <div class="hero-apple__actions">
        <router-link class="solid-btn" :to="userStore.isLoggedIn ? '/editor' : '/auth'">
          {{ userStore.isLoggedIn ? "开始写作" : "注册开始" }}
        </router-link>
        <router-link class="text-link" :to="articlesLibraryLink">打开文章库</router-link>
      </div>
    </div>
  </section>

  <section class="content-section home-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">精选项目</p>
        <h2>精选项目</h2>
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

  <section class="content-section home-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">最近文章</p>
        <h2>最近文章</h2>
      </div>
      <router-link class="ghost-btn" :to="articlesLibraryLink">打开文章库</router-link>
    </div>

    <div v-if="featuredArticles.length" class="article-grid">
      <ArticleCard v-for="item in featuredArticles" :key="`featured-${item.id}`" :item="item" />
    </div>
    <div v-else class="empty-panel">
      <h4>还没有已发布文章</h4>
      <p>写一篇文章来开始你的内容积累。</p>
    </div>
  </section>

  <section class="home-outro">
    <router-link class="text-link" :to="userStore.isLoggedIn ? '/editor' : '/auth'">
      {{ userStore.isLoggedIn ? "开始写作 →" : "注册一个账号 →" }}
    </router-link>
  </section>
</template>
```

**验收**

- 模板不再引用被删状态：`recommendedAuthors`、`categories`、`tags`、`keyword`、`total`、`totalPages`、`page`、`selectedCategory`、`selectedTag`、`hero-metrics`、`tag-row--hero`、`hero-orbit`、`portfolio-intro`。

### T2: script 重写

替换 `<script setup>` 为：

```js
<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { listArticles, listTrendingArticles } from "../api/article";
import { listProjects } from "../api/project";
import { useUserStore } from "../stores/user";
import ArticleCard from "../components/ArticleCard.vue";
import { toAssetUrl } from "../utils/asset";
import { getAuthorProfile } from "../api/profile";

const userStore = useUserStore();
const articles = ref([]);
const trendingArticles = ref([]);
const featuredProjects = ref([]);
const ownerProfile = ref(null);

const siteOwner = computed(() => ownerProfile.value || userStore.profile || {});
const browseProjectsLink = computed(() => ({ path: "/projects", query: { sort: "featured" } }));
const articlesLibraryLink = computed(() => ({ path: "/articles", query: { sort: "latest" } }));
const heroTitle = computed(() => siteOwner.value?.headline || "把你的博客做成作品集、学习归档和长期表达入口");
const heroText = computed(() => {
  if (siteOwner.value?.bio) return siteOwner.value.bio;
  return "这个平台适合把项目案例、学习笔记和公开写作放在同一个站点里，既能用于求职展示，也能长期沉淀内容资产。";
});

const featuredArticles = computed(() => {
  if (trendingArticles.value.length >= 3) return trendingArticles.value.slice(0, 3);
  return articles.value.slice(0, 3);
});

async function fetchArticles() {
  const { data } = await listArticles({ page: 1, pageSize: 9 });
  articles.value = data.items || [];
}

async function fetchTrending() {
  const { data } = await listTrendingArticles();
  trendingArticles.value = data.items || [];
}

async function fetchProjects() {
  const { data } = await listProjects({ page: 1, pageSize: 3, featured: true });
  featuredProjects.value = data.items || [];
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

watch(
  () => userStore.profile?.id,
  (id) => {
    fetchOwnerProfile(id);
  },
  { immediate: true }
);

onMounted(async () => {
  await Promise.all([fetchArticles(), fetchTrending(), fetchProjects()]);
});
</script>
```

**删除清单**：`useRouter`/`router`、`keyword`、`categories`、`tags`、`selectedCategory`、`selectedTag`、`page`、`pageSize`、`total`、`totalPages`、`recommendedAuthors`、`featuredAuthorLink`、`ownerProjectsLink`、`studyArticlesLink`、`isStudyText`、`isStudyArticle`、`goToPage`、`resetFilters`、`fetchRecommendedAuthors`、`fetchMetadata`；import 移除 `getMetadata`、`getRecommendedAuthors`。

**验收**

- 无未使用变量/import（`router` 已删，页面不再用跳转编程式导航）。
- 数据加载仍覆盖 hero 标题、精选项目、最近文章三处。

### T3: style 重写

替换整个 `<style scoped>`：

```css
<style scoped>
.home-section {
  padding-top: 28px;
}

.hero-apple {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  padding: clamp(56px, 10vw, 120px) clamp(24px, 6vw, 80px);
  text-align: center;
  display: grid;
  justify-items: center;
}

.hero-apple__eyebrow {
  color: var(--accent);
  font-weight: 600;
}

.hero-apple__title {
  max-width: 18ch;
  margin: 14px 0 0;
  font-size: clamp(2.75rem, 6vw, 5rem);
  line-height: 1.05;
  letter-spacing: -0.02em;
  font-weight: 700;
  color: #1d1d1f;
}

.hero-apple__text {
  max-width: 52ch;
  margin: 18px 0 0;
  font-size: clamp(1.05rem, 1.6vw, 1.35rem);
  line-height: 1.6;
  color: #48484a;
}

.hero-apple__actions {
  margin-top: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 22px;
  flex-wrap: wrap;
}

.text-link {
  color: var(--accent);
  font-weight: 600;
  text-decoration: none;
}

.text-link:hover {
  text-decoration: underline;
}

.project-grid--hero .project-card--hero {
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 32%),
    rgba(255, 255, 255, 0.06);
}

.home-outro {
  padding: 48px 0 8px;
  text-align: center;
}

@media (max-width: 480px) {
  .hero-apple__actions {
    flex-direction: column;
  }

  .hero-apple__actions .solid-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
```

**删除**：`.home-hero__copy`、`.hero-orbit`、`.orbit-card` 系列、`.content-section--split`、`.content-split-card`、`.author-grid`、`.author-card` 系列及对应媒体查询。

**验收**

- scoped 样式内不再有 orbit/author/split 相关选择器。
- 浅色 hero 在暗色页面上形成品牌对比；桌面与移动端（480px 以下按钮纵向堆叠）均正常。

### T4: 验证

- `cd web; npx vite build` 通过。
- `grep -rn "recommendedAuthors\|hero-orbit\|portfolio-intro\|studyArticlesLink" web/src` 确认无残留引用。
- `grep -n "hero-panel" web/src`：确认 `hero-panel` 仍被其他视图使用（如 AuthorView/ProjectsView），本页删除不影响全局类定义（`main.css` 不动）。
- 手动冒烟：无数据（空库）与有数据两条路径；未登录/已登录 hero CTA 文案切换。

## Suggested commits

1. `refactor(web): strip homepage to hero + featured + latest`（T1+T2）
2. `style(web): apple-style light hero and brand typography`（T3）

## Exit criteria

- 首页仅剩三大区块 + 收尾 CTA；hero 为浅色苹果式大标题；无死代码；`vite build` 通过；其他视图无回归（`hero-panel` 等全局类保留）。
