<template>
  <section class="about-page">
    <section class="about-hero panel-card about-hero--brand">
      <div class="about-hero__main">
        <p class="eyebrow">关于</p>
        <div class="author-profile-head">
          <div class="author-avatar about-avatar">
            <img v-if="profile?.avatar" :src="toAssetUrl(profile.avatar)" :alt="profile?.username || 'avatar'" />
            <span v-else>{{ avatarFallback }}</span>
          </div>
          <div>
            <h2>{{ profile?.username || "你的个人站点" }}</h2>
            <p class="author-headline">
              {{ profile?.headline || "一个用于展示项目、写作输出和长期学习记录的个人空间。" }}
            </p>
          </div>
        </div>

        <p class="detail-summary">
          {{ profile?.bio || "补齐个人资料后，这里会逐渐变成一页真正能对外展示的自我介绍和线上简历。" }}
        </p>

        <div v-if="profile?.motto" class="about-motto">
          "{{ profile.motto }}"
        </div>

        <div class="inline-actions">
          <router-link v-if="userStore.isLoggedIn" class="solid-btn" to="/profile">编辑资料</router-link>
          <router-link v-else class="solid-btn" to="/auth">登录后定制</router-link>
          <a v-if="profile?.resumeUrl" class="ghost-btn" :href="profile.resumeUrl" target="_blank" rel="noreferrer">查看简历</a>
          <router-link v-if="profile?.id" class="ghost-btn" :to="`/author/${profile.id}`">作者主页</router-link>
        </div>
      </div>

      <div class="about-hero__side about-hero__side--constellation">
        <div class="about-constellation">
          <article class="about-node about-node--main">
            <strong>{{ profile?.currentRole || "创作者" }}</strong>
            <p>当前身份</p>
          </article>
          <article class="about-node about-node--metric">
            <strong>{{ profile?.yearsLabel || "持续积累中" }}</strong>
            <p>经历标签</p>
          </article>
          <article class="about-node about-node--links">
            <strong>{{ socialLinks.length }}</strong>
            <p>公开链接</p>
          </article>
        </div>
      </div>
    </section>

    <section class="content-section">
      <div class="profile-snapshot-grid">
        <article class="panel-card snapshot-card">
          <p class="eyebrow">方向</p>
          <div v-if="profile?.focusAreas?.length" class="tag-row">
            <span v-for="item in profile.focusAreas" :key="item" class="tag-chip">{{ item }}</span>
          </div>
          <p v-else class="detail-summary">在个人资料里补充你的关注方向，让别人更快理解你想做什么。</p>
        </article>

        <article class="panel-card snapshot-card">
          <p class="eyebrow">技能栏</p>
          <div v-if="profile?.skills?.length" class="tag-row">
            <span v-for="item in profile.skills" :key="item" class="tag-chip"># {{ item }}</span>
          </div>
          <p v-else class="detail-summary">这里会展示你的技能标签，形成简洁清晰的能力概览。</p>
        </article>

        <article class="panel-card snapshot-card">
          <p class="eyebrow">链接</p>
          <div v-if="socialLinks.length" class="about-links">
            <a v-for="item in socialLinks" :key="item.label" class="inline-link" :href="item.url" target="_blank" rel="noreferrer">
              {{ item.label }}
            </a>
          </div>
          <p v-else class="detail-summary">官网、GitHub、CSDN 等链接都可以在个人资料页里补充。</p>
        </article>
      </div>
    </section>

    <section class="content-section content-section--split">
      <article class="panel-card about-capability-card">
        <div class="section-head section-head--compact">
          <div>
            <p class="eyebrow">项目</p>
            <h3>精选作品</h3>
          </div>
          <router-link class="ghost-btn" :to="projectsLink">全部项目</router-link>
        </div>

        <div v-if="projects.length" class="project-grid">
          <article v-for="item in projects.slice(0, 3)" :key="item.id" class="project-card panel-card">
            <div v-if="item.coverImage" class="project-card__cover">
              <img :src="toAssetUrl(item.coverImage)" :alt="item.title" />
            </div>
            <div class="project-card__head">
              <div>
                <h3>{{ item.title }}</h3>
                <p class="table-note">{{ item.roleLabel || item.duration || "项目案例" }}</p>
              </div>
              <span v-if="item.isFeatured" class="tag-chip">精选</span>
            </div>
            <p class="detail-summary">{{ item.summary || "这个项目暂时还没有摘要。" }}</p>
            <router-link class="solid-btn" :to="`/projects/${item.id}`">打开案例详情</router-link>
          </article>
        </div>
        <div v-else class="empty-panel">
          <h4>还没有项目</h4>
          <p>发布几个项目后，这一页就会开始像一个真正的作品集介绍页。</p>
        </div>
      </article>

      <article class="panel-card about-capability-card about-capability-card--contact">
        <div class="section-head section-head--compact">
          <div>
            <p class="eyebrow">写作</p>
            <h3>最近文章</h3>
          </div>
          <router-link class="ghost-btn" :to="articlesLink">全部文章</router-link>
        </div>

        <div v-if="articles.length" class="article-grid">
          <ArticleCard v-for="item in articles.slice(0, 3)" :key="item.id" :item="item" />
        </div>
        <div v-else class="empty-panel">
          <h4>还没有文章</h4>
          <p>文章会让这个页面更完整，因为它不只展示结果，也展示你的思考过程。</p>
        </div>
      </article>
    </section>

    <section class="content-section">
      <div class="contact-cta panel-card">
        <div>
          <p class="eyebrow">联系</p>
          <h3>把这一页变成真正可联系的个人页面</h3>
          <p class="detail-summary">{{ contactText }}</p>
        </div>
        <div class="inline-actions">
          <a v-if="profile?.email" class="solid-btn" :href="`mailto:${profile.email}`">给我发邮件</a>
          <a v-if="profile?.websiteUrl" class="ghost-btn" :href="profile.websiteUrl" target="_blank" rel="noreferrer">访问网站</a>
          <router-link v-if="userStore.isLoggedIn" class="ghost-btn" to="/profile">继续完善资料</router-link>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { useUserStore } from "../stores/user";
import { getAuthorProfile } from "../api/profile";
import { toAssetUrl } from "../utils/asset";
import ArticleCard from "../components/ArticleCard.vue";

const userStore = useUserStore();
const projects = ref([]);
const articles = ref([]);

const profile = computed(() => userStore.profile);
const avatarFallback = computed(() => (profile.value?.username || "P").slice(0, 1).toUpperCase());
const socialLinks = computed(() => {
  if (!profile.value) return [];
  return [
    { label: "网站", url: profile.value.websiteUrl },
    { label: "简历", url: profile.value.resumeUrl },
    { label: "GitHub", url: profile.value.githubUrl },
    { label: "Gitee", url: profile.value.giteeUrl },
    { label: "Juejin", url: profile.value.juejinUrl },
    { label: "CSDN", url: profile.value.csdnUrl }
  ].filter((item) => item.url);
});
const projectsLink = computed(() =>
  profile.value?.id
    ? {
        path: "/projects",
        query: {
          authorId: String(profile.value.id),
          authorName: profile.value.username || "",
          sort: "featured"
        }
      }
    : { path: "/projects", query: { sort: "featured" } }
);
const articlesLink = computed(() =>
  profile.value?.id
    ? {
        path: "/articles",
        query: {
          authorId: String(profile.value.id),
          authorName: profile.value.username || "",
          sort: "latest"
        }
      }
    : { path: "/articles", query: { sort: "latest" } }
);

const contactText = computed(() => {
  if (profile.value?.email) {
    return "如果你想聊合作、项目、实习机会或有意思的想法，邮件会是联系我最快的方式。";
  }
  return "你可以在个人资料页补充邮箱或网站链接，让这里真正变成一个可联系的页面。";
});

watch(
  () => userStore.profile?.id,
  async (id) => {
    if (!id) {
      projects.value = [];
      articles.value = [];
      return;
    }
    const { data } = await getAuthorProfile(id);
    projects.value = data.projects || [];
    articles.value = data.articles || [];
  },
  { immediate: true }
);
</script>

<style scoped>
.about-hero--brand {
  background:
    radial-gradient(circle at 84% 18%, rgba(255, 209, 102, 0.14), transparent 24%),
    radial-gradient(circle at 12% 76%, rgba(255, 138, 76, 0.12), transparent 26%),
    rgba(255, 255, 255, 0.06);
}

.about-hero__side--constellation {
  display: grid;
  place-items: center;
}

.about-constellation {
  display: grid;
  gap: 14px;
  width: min(100%, 280px);
}

.about-node {
  display: grid;
  gap: 8px;
  padding: 18px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.about-node strong,
.about-node p {
  margin: 0;
}

.about-node p {
  color: var(--text-soft);
  font-size: 0.92rem;
}

.section-head--compact {
  margin-bottom: 0;
}

.about-capability-card {
  display: grid;
  gap: 18px;
}

.about-capability-card--contact {
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.06);
}

@media (max-width: 768px) {
  .about-hero__side--constellation {
    display: none;
  }

  .about-hero--brand {
    padding: 20px;
  }

  .profile-snapshot-grid {
    gap: 14px;
  }

  .contact-cta {
    flex-direction: column;
    align-items: flex-start;
    padding: 20px;
  }

  .contact-cta .inline-actions {
    flex-direction: column;
    width: 100%;
  }

  .contact-cta .inline-actions .solid-btn,
  .contact-cta .inline-actions .ghost-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .about-avatar {
    width: 72px;
    height: 72px;
  }

  .about-node {
    padding: 12px 14px;
  }

  .about-constellation {
    gap: 10px;
    width: 100%;
  }
}
</style>
