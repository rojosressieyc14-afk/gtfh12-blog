<template>
  <section v-if="author" class="content-section author-page">
    <div class="author-hero author-hero--portfolio panel-card author-page-hero">
      <div class="author-hero__copy">
        <p class="eyebrow">个人名片</p>
        <div class="author-profile-head">
          <div class="author-avatar">
            <img v-if="author.avatar" :src="toAssetUrl(author.avatar)" :alt="author.username" />
            <span v-else>{{ author.username.slice(0, 1).toUpperCase() }}</span>
          </div>
          <div>
            <h2>{{ author.username }}</h2>
            <p class="author-headline">
              {{ author.headline || "持续做项目、写文章、记录学习，把成长过程公开展示出来。" }}
            </p>
          </div>
        </div>

        <p class="detail-summary">
          {{ author.bio || "这位作者还没有补充完整简介，不过作品和文章已经开始替他表达。" }}
        </p>

        <div v-if="author.skills?.length" class="tag-row">
          <span v-for="skill in author.skills" :key="skill" class="tag-chip"># {{ skill }}</span>
        </div>

        <div class="profile-meta">
          <span v-if="author.currentRole">当前角色：{{ author.currentRole }}</span>
          <span v-if="author.yearsLabel">经验阶段：{{ author.yearsLabel }}</span>
          <span v-if="author.location">所在城市：{{ author.location }}</span>
          <span v-if="author.email">联系邮箱：{{ author.email }}</span>
        </div>

        <div v-if="socialLinks.length" class="inline-actions">
          <a
            v-for="item in socialLinks"
            :key="item.label"
            class="ghost-btn"
            :href="item.url"
            target="_blank"
            rel="noreferrer"
          >
            {{ item.label }}
          </a>
        </div>
        <router-link class="ghost-btn" :to="authorArticlesLink">查看全部文章</router-link>
      </div>

      <div class="author-stats author-stats--grid">
        <div class="author-stat">
          <strong>{{ stats.projects || 0 }}</strong>
          <span>项目案例</span>
        </div>
        <div class="author-stat">
          <strong>{{ stats.articles || 0 }}</strong>
          <span>公开文章</span>
        </div>
        <div class="author-stat">
          <strong>{{ stats.comments || 0 }}</strong>
          <span>互动评论</span>
        </div>
        <div class="author-stat">
          <strong>{{ stats.likes || 0 }}</strong>
          <span>获得点赞</span>
        </div>
      </div>
    </div>

    <section v-if="author.focusAreas?.length || author.motto || author.resumeUrl" class="content-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">人物侧写</p>
          <h3>更完整的个人展示</h3>
        </div>
        <a v-if="author.resumeUrl" class="ghost-btn" :href="author.resumeUrl" target="_blank" rel="noreferrer">打开简历</a>
      </div>

      <div class="profile-snapshot-grid">
        <article v-if="author.motto" class="panel-card snapshot-card">
          <p class="eyebrow">个人主张</p>
          <h4>{{ author.motto }}</h4>
        </article>

        <article v-if="author.focusAreas?.length" class="panel-card snapshot-card">
          <p class="eyebrow">关注方向</p>
          <div class="tag-row">
            <span v-for="item in author.focusAreas" :key="item" class="tag-chip">{{ item }}</span>
          </div>
        </article>
      </div>
    </section>

    <section class="content-section content-section--split">
      <article class="panel-card author-capability-card">
        <div class="section-head section-head--compact">
          <div>
            <p class="eyebrow">能力画像</p>
            <h3>别人会怎样快速理解这位作者</h3>
          </div>
        </div>

        <div class="author-capability-list">
          <article v-for="item in capabilityCards" :key="item.title" class="author-capability-item">
            <span>{{ item.kicker }}</span>
            <strong>{{ item.title }}</strong>
            <p>{{ item.description }}</p>
          </article>
        </div>
      </article>

      <article class="panel-card author-capability-card author-capability-card--accent">
        <div class="section-head section-head--compact">
          <div>
            <p class="eyebrow">合作提示</p>
            <h3>这一页已经能替作者说明什么</h3>
          </div>
        </div>

        <div class="author-capability-list">
          <article v-for="item in connectionCards" :key="item.title" class="author-capability-item">
            <span>{{ item.kicker }}</span>
            <strong>{{ item.title }}</strong>
            <p>{{ item.description }}</p>
          </article>
        </div>
      </article>
    </section>

    <section v-if="socialCards.length" class="content-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">站外延伸</p>
          <h3>继续从外部链接认识这位作者</h3>
        </div>
      </div>

      <div class="services-grid">
        <a
          v-for="item in socialCards"
          :key="item.label"
          class="panel-card service-card author-link-card"
          :href="item.url"
          target="_blank"
          rel="noreferrer"
        >
          <p class="eyebrow">{{ item.label }}</p>
          <h4>{{ item.title }}</h4>
          <p>{{ item.description }}</p>
        </a>
      </div>
    </section>

    <section class="content-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">作品集</p>
          <h3>已发布项目</h3>
        </div>
        <router-link class="ghost-btn" :to="authorProjectsLink">查看全部项目</router-link>
      </div>

      <div v-if="projects.length" class="project-grid">
        <article v-for="item in projects" :key="item.id" class="project-card panel-card">
          <div v-if="item.coverImage" class="project-card__cover">
            <img :src="toAssetUrl(item.coverImage)" :alt="item.title" />
          </div>
          <div class="project-card__head">
            <div>
              <h3>{{ item.title }}</h3>
              <p class="table-note">{{ formatDate(item.publishedAt || item.createdAt) }}</p>
            </div>
            <span v-if="item.isFeatured" class="tag-chip">精选</span>
          </div>
          <p class="detail-summary">{{ item.summary || "这个项目暂时还没有摘要。" }}</p>
          <div v-if="item.techStacks?.length" class="tag-row">
            <span v-for="stack in item.techStacks.slice(0, 4)" :key="stack" class="tag-chip"># {{ stack }}</span>
          </div>
          <router-link class="solid-btn" :to="`/projects/${item.id}`">查看项目</router-link>
        </article>
      </div>
      <div v-else class="empty-panel">
        <h4>暂时还没有已发布项目</h4>
        <p>这里很适合放作品案例、过程拆解和上线复盘，慢慢就会变成完整的作品集页面。</p>
      </div>
    </section>

    <section class="content-section">
      <div class="section-head">
        <div>
          <p class="eyebrow">写作输出</p>
          <h3>已发布文章</h3>
        </div>
      </div>
      <div v-if="articles.length" class="article-grid">
        <ArticleCard v-for="item in articles" :key="item.id" :item="item" />
      </div>
      <div v-else class="empty-panel">
        <h4>暂时还没有已发布文章</h4>
        <p>等内容逐渐丰富后，这里会把作品和写作串起来，形成更完整的个人表达。</p>
      </div>
    </section>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { getAuthorProfile } from "../api/profile";
import ArticleCard from "../components/ArticleCard.vue";
import { toAssetUrl } from "../utils/asset";

const route = useRoute();
const author = ref(null);
const articles = ref([]);
const projects = ref([]);
const stats = ref({});

const authorProjectsLink = computed(() => ({
  path: "/projects",
  query: author.value?.id
    ? {
        authorId: String(author.value.id),
        authorName: author.value.username || "",
        sort: "featured"
      }
    : { sort: "featured" }
}));

const authorArticlesLink = computed(() => ({
  path: "/articles",
  query: author.value?.id
    ? {
        authorId: String(author.value.id),
        authorName: author.value.username || "",
        sort: "latest"
      }
    : { sort: "latest" }
}));

const socialLinks = computed(() => {
  if (!author.value) return [];
  return [
    { label: "个人网站", url: author.value.websiteUrl },
    { label: "个人简历", url: author.value.resumeUrl },
    { label: "GitHub", url: author.value.githubUrl },
    { label: "Gitee", url: author.value.giteeUrl },
    { label: "稀土掘金", url: author.value.juejinUrl },
    { label: "CSDN", url: author.value.csdnUrl }
  ].filter((item) => item.url);
});

const socialCards = computed(() =>
  socialLinks.value.map((item) => ({
    ...item,
    title:
      {
        个人网站: "查看完整线上名片",
        个人简历: "快速了解履历背景",
        GitHub: "浏览代码与仓库",
        Gitee: "查看国内镜像仓库",
        稀土掘金: "继续阅读更多文章",
        CSDN: "浏览更多技术记录"
      }[item.label] || item.label,
    description:
      {
        个人网站: "适合承接更多项目介绍、联系信息和更完整的个人表达。",
        个人简历: "如果你想快速了解背景、经历和技能结构，这里会是最直接的入口。",
        GitHub: "仓库、代码风格和工程习惯会让这位作者的能力更具体。",
        Gitee: "如果作者有国内镜像或协作仓库，这里也值得点进去看。",
        稀土掘金: "适合继续了解作者的技术输出和长期写作轨迹。",
        CSDN: "适合查看更多学习记录、教程和总结型文章。"
      }[item.label] || "继续从站外入口了解这位作者。"
  }))
);

const capabilityCards = computed(() => [
  {
    kicker: "角色",
    title: author.value?.currentRole || "当前角色仍在持续完善中",
    description:
      author.value?.headline || "这一栏最适合快速说明作者现在在做什么、擅长什么、希望被怎样认识。"
  },
  {
    kicker: "技能",
    title: (author.value?.skills || []).slice(0, 3).join(" / ") || "核心技能仍可继续补充",
    description:
      author.value?.bio || "结合技能、简介和项目案例，这一页会逐渐像一份更完整的线上作品集名片。"
  },
  {
    kicker: "方向",
    title: (author.value?.focusAreas || []).join(" / ") || "建议补充 3 到 5 个重点方向",
    description: "关注方向越清晰，别人越容易理解这位作者最想承接什么类型的内容和项目。"
  }
]);

const connectionCards = computed(() => [
  {
    kicker: "联系",
    title: author.value?.email ? `可通过 ${author.value.email} 联系` : "建议补一个直接联系入口",
    description: author.value?.email
      ? "邮箱已经可以承接项目交流、合作沟通和机会连接。"
      : "如果补上邮箱或个人网站，这个页面会更像可直接转发的线上名片。"
  },
  {
    kicker: "案例",
    title: projects.value.length ? `已公开 ${projects.value.length} 个项目案例` : "项目案例还可以继续补充",
    description: projects.value.length
      ? "项目案例越完整，这位作者的能力边界和做事方式就越容易被理解。"
      : "先放 2 到 3 个有代表性的项目，会让这个页面的说服力明显提高。"
  },
  {
    kicker: "输出",
    title: articles.value.length ? `已公开 ${articles.value.length} 篇文章` : "写作输出还可以继续积累",
    description: articles.value.length
      ? "文章会把项目结果背后的思考路径补充完整，让作者形象更立体。"
      : "如果后续补上学习笔记、复盘和总结，这一页会更像长期成长记录。"
  }
]);

async function loadAuthor() {
  const { data } = await getAuthorProfile(route.params.id);
  author.value = data.user;
  articles.value = data.articles || [];
  projects.value = data.projects || [];
  stats.value = data.stats || {};
}

function formatDate(value) {
  return new Date(value).toLocaleDateString("zh-CN");
}

watch(() => route.params.id, loadAuthor);
onMounted(loadAuthor);
</script>

<style scoped>
.author-page {
  margin-top: 28px;
}

.author-page-hero {
  background:
    radial-gradient(circle at 18% 20%, rgba(255, 138, 76, 0.12), transparent 24%),
    radial-gradient(circle at 82% 24%, rgba(255, 209, 102, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.07);
}

.section-head--compact {
  margin-bottom: 0;
}

.author-capability-card {
  display: grid;
  gap: 18px;
}

.author-capability-card--accent {
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.06);
}

.author-capability-list {
  display: grid;
  gap: 12px;
}

.author-capability-item {
  padding: 16px 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.author-capability-item span {
  color: var(--accent);
  font-size: 0.84rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.author-capability-item strong,
.author-capability-item p {
  display: block;
  margin: 0;
}

.author-capability-item strong {
  margin-top: 8px;
}

.author-capability-item p {
  margin-top: 8px;
  color: var(--text-soft);
}

.author-link-card {
  text-decoration: none;
  color: inherit;
  transition: transform 0.22s ease, border-color 0.22s ease, background 0.22s ease;
}

.author-link-card:hover {
  transform: translateY(-3px);
  border-color: rgba(255, 209, 102, 0.22);
  background: rgba(255, 255, 255, 0.08);
}

@media (max-width: 768px) {
  .author-page {
    margin-top: 14px;
  }

  .author-hero--portfolio {
    padding: 20px;
  }

  .author-profile-head {
    flex-direction: column;
    text-align: center;
  }

  .author-stats--grid {
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }

  .author-stat {
    padding: 12px;
    min-width: auto;
  }

  .author-stat strong {
    font-size: 1.4rem;
  }

  .profile-meta {
    flex-direction: column;
    gap: 6px;
  }
}

@media (max-width: 480px) {
  .author-stats--grid {
    grid-template-columns: 1fr;
  }

  .author-capability-item {
    padding: 12px 14px;
  }

  .author-capability-item p {
    font-size: 0.92rem;
  }
}
</style>
