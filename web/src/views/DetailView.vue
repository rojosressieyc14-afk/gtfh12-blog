<template>
  <section v-if="article" class="article-detail-page">
    <div class="detail-layout panel-card article-detail-shell article-detail-shell--brand">
      <div class="article-detail-top">
        <div class="article-detail-main-copy">
          <div class="detail-meta">
            <router-link class="status-chip published" :to="`/author/${article.author?.id}`">
              {{ article.author?.username || "匿名作者" }}
            </router-link>
            <span>{{ article.category?.name || "未分类" }}</span>
          </div>

          <h2 class="detail-title">{{ article.title }}</h2>
          <p class="detail-summary">{{ article.summary || "这篇文章暂时还没有补充摘要。" }}</p>

          <div class="detail-stats article-detail-stats">
            <span>{{ formatDate(article.publishedAt || article.createdAt) }}</span>
            <span>{{ commentCount }} 条评论</span>
            <span>{{ article.viewCount || 0 }} 次阅读</span>
            <span>{{ readingMinutes }} 分钟阅读</span>
          </div>

          <div v-if="article.tags?.length" class="tag-row tag-row--detail">
            <span v-for="tag in article.tags" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
          </div>

          <div class="reaction-row">
            <button class="ghost-btn reaction-btn" @click="onLike">
              {{ article.isLiked ? "已点赞" : "点赞" }} · {{ article.likesCount || 0 }}
            </button>
            <button class="ghost-btn reaction-btn" @click="onFavorite">
              {{ article.isFavorited ? "已收藏" : "收藏" }} · {{ article.favoritesCount || 0 }}
            </button>
          </div>
        </div>

        <div class="article-detail-side-visual">
          <div class="article-detail-orbit">
            <div class="article-detail-node article-detail-node--one"></div>
            <div class="article-detail-node article-detail-node--two"></div>
            <div class="article-detail-node article-detail-node--three"></div>
            <div v-if="coverUrl" class="detail-cover article-detail-cover">
              <img :src="coverUrl" :alt="article.title" />
            </div>
          </div>
        </div>
      </div>

      <section class="article-reading-layout">
        <article class="markdown-body article-detail-body" v-html="html"></article>

        <aside class="article-reading-side">
          <section class="project-detail-panel article-outline-panel">
            <p class="eyebrow">阅读导航</p>
            <h3>文章结构</h3>
            <div v-if="articleOutline.length" class="article-outline-list">
              <button
                v-for="item in articleOutline"
                :key="item.id"
                class="article-outline-link"
                :class="`article-outline-link--${item.level}`"
                @click="scrollToHeading(item.id)"
              >
                {{ item.text }}
              </button>
            </div>
            <p v-else class="detail-summary">继续补充小标题后，这里会自动生成文章目录。</p>
          </section>

          <section class="project-detail-panel article-reading-panel">
            <p class="eyebrow">阅读提示</p>
            <h3>这篇内容适合怎么读</h3>
            <div class="project-highlight-list project-highlight-list--detail">
              <span>先看摘要和目录，快速判断这篇内容是不是你要找的。</span>
              <span>如果是案例型文章，重点看过程、难点和最后结论。</span>
              <span>如果是学习型文章，建议连同标签一起看，方便后续继续检索。</span>
            </div>
          </section>

          <section class="project-detail-panel">
            <p class="eyebrow">作者入口</p>
            <h3>继续了解这位作者</h3>
            <p class="detail-summary">你可以继续查看这位作者的项目案例、公开文章和完整个人资料。</p>
            <router-link class="ghost-btn" :to="`/author/${article.author?.id}`">打开作者页</router-link>
          </section>
        </aside>
      </section>

      <section class="comment-panel article-comment-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">讨论区</p>
            <h3>评论与回复</h3>
          </div>
        </div>

        <form v-if="userStore.isLoggedIn" class="stack-form" @submit.prevent="submitComment">
          <textarea
            v-model.trim="commentText"
            class="field-area field-area--small"
            placeholder="写下你对这篇文章的想法。"
          ></textarea>
          <button class="solid-btn comment-btn">发表评论</button>
        </form>
        <p v-else class="detail-summary">登录后可以参与讨论。</p>

        <div class="comment-list">
          <article v-for="item in comments" :key="item.id" class="comment-item">
            <header>
              <strong>{{ item.user?.username || "匿名用户" }}</strong>
              <span>{{ formatDate(item.createdAt) }}</span>
            </header>
            <p>{{ item.content }}</p>
            <button v-if="userStore.isLoggedIn" class="reply-link" @click="replyTo = item.id">回复</button>

            <form v-if="replyTo === item.id" class="stack-form reply-form" @submit.prevent="submitReply(item.id)">
              <textarea
                v-model.trim="replyText"
                class="field-area field-area--small"
                placeholder="写下你的回复。"
              ></textarea>
              <button class="ghost-btn comment-btn">发送回复</button>
            </form>

            <div v-if="item.replies?.length" class="reply-list">
              <article v-for="reply in item.replies" :key="reply.id" class="reply-item">
                <header>
                  <strong>{{ reply.user?.username || "匿名用户" }}</strong>
                  <span>{{ formatDate(reply.createdAt) }}</span>
                </header>
                <p>{{ reply.content }}</p>
              </article>
            </div>
          </article>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { marked } from "marked";
import { useRoute } from "vue-router";
import { createComment, getArticle, listComments, toggleFavorite, toggleLike } from "../api/article";
import { useUserStore } from "../stores/user";
import { toAssetUrl } from "../utils/asset";

const route = useRoute();
const userStore = useUserStore();
const article = ref(null);
const comments = ref([]);
const commentText = ref("");
const replyTo = ref(null);
const replyText = ref("");

function buildHeadingId(text, index) {
  const safe = text
    .toLowerCase()
    .replace(/<[^>]+>/g, "")
    .replace(/[^\w\u4e00-\u9fa5\s-]/g, "")
    .trim()
    .replace(/\s+/g, "-");
  return `heading-${safe || "section"}-${index}`;
}

const articleOutline = computed(() => {
  const matches = [...(article.value?.content || "").matchAll(/^(#{1,3})\s+(.+)$/gm)];
  return matches.map((match, index) => ({
    level: match[1].length,
    text: match[2].trim(),
    id: buildHeadingId(match[2], index)
  }));
});

const html = computed(() => {
  const renderer = new marked.Renderer();
  renderer.heading = (token) => {
    const text = token.text || "";
    const id = articleOutline.value.find((item) => item.text === text)?.id || buildHeadingId(text, 0);
    return `<h${token.depth} id="${id}">${text}</h${token.depth}>`;
  };
  return marked.parse(article.value?.content || "", { renderer });
});

const coverUrl = computed(() => toAssetUrl(article.value?.coverImage));
const commentCount = computed(() => comments.value.reduce((sum, item) => sum + 1 + (item.replies?.length || 0), 0));
const readingMinutes = computed(() => {
  const count = (article.value?.content || "").replace(/\s+/g, "").length;
  return Math.max(1, Math.ceil(count / 500));
});

async function loadDetail() {
  const { data } = await getArticle(route.params.id);
  article.value = data.item;
}

async function loadComments() {
  const { data } = await listComments(route.params.id);
  comments.value = data.items;
}

async function submitComment() {
  if (!commentText.value) return;
  await createComment(route.params.id, { content: commentText.value });
  commentText.value = "";
  await loadComments();
}

async function submitReply(parentId) {
  if (!replyText.value) return;
  await createComment(route.params.id, { content: replyText.value, parentId });
  replyText.value = "";
  replyTo.value = null;
  await loadComments();
}

async function onLike() {
  if (!userStore.isLoggedIn) return;
  const { data } = await toggleLike(route.params.id);
  article.value = { ...article.value, ...data.item };
}

async function onFavorite() {
  if (!userStore.isLoggedIn) return;
  const { data } = await toggleFavorite(route.params.id);
  article.value = { ...article.value, ...data.item };
}

function scrollToHeading(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(async () => {
  await loadDetail();
  await loadComments();
});
</script>

<style scoped>
.article-detail-shell--brand {
  background:
    radial-gradient(circle at 84% 18%, rgba(255, 209, 102, 0.1), transparent 24%),
    radial-gradient(circle at 14% 76%, rgba(255, 138, 76, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.07);
}

.article-detail-top {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(320px, 0.85fr);
  gap: 24px;
  align-items: center;
}

.article-detail-side-visual {
  position: relative;
}

.article-detail-orbit {
  position: relative;
  min-height: 360px;
}

.article-reading-layout {
  margin-top: 28px;
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(260px, 0.85fr);
  gap: 22px;
  align-items: start;
}

.article-reading-side {
  display: grid;
  gap: 18px;
  position: sticky;
  top: 112px;
}

.article-outline-list {
  display: grid;
  gap: 10px;
}

.article-outline-link {
  text-align: left;
  color: inherit;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 12px 14px;
  cursor: pointer;
  transition: transform 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.article-outline-link:hover {
  transform: translateY(-2px);
  border-color: rgba(255, 209, 102, 0.28);
  background: rgba(255, 209, 102, 0.08);
}

.article-outline-link--2 {
  margin-left: 14px;
}

.article-outline-link--3 {
  margin-left: 28px;
}

.article-detail-node {
  position: absolute;
  border: 2px solid rgba(255, 243, 231, 0.74);
  background: rgba(255, 204, 153, 0.14);
  animation: articleNodeFloat 5s ease-in-out infinite;
}

.article-detail-node--one {
  width: 60px;
  height: 60px;
  left: 22px;
  top: 18px;
  border-radius: 18px;
}

.article-detail-node--two {
  width: 54px;
  height: 54px;
  right: 26px;
  top: 86px;
  border-radius: 999px;
  animation-delay: 0.7s;
}

.article-detail-node--three {
  width: 64px;
  height: 54px;
  left: 34px;
  bottom: 26px;
  clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
  border-radius: 10px;
  animation-delay: 1.1s;
}

@keyframes articleNodeFloat {
  0%,
  100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-7px) rotate(5deg);
  }
}

@media (max-width: 1080px) {
  .article-reading-layout {
    grid-template-columns: 1fr;
  }

  .article-reading-side {
    position: static;
  }
}

@media (max-width: 960px) {
  .article-detail-top {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .article-detail-orbit {
    min-height: 240px;
  }

  .article-detail-stats {
    width: 100%;
    flex-wrap: wrap;
  }

  .article-detail-side-visual {
    display: none;
  }

  .article-detail-shell--brand {
    padding: 20px;
  }

  .reaction-row {
    flex-direction: column;
  }

  .reaction-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .detail-title {
    font-size: clamp(1.4rem, 6vw, 2rem);
  }

  .detail-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .article-detail-body {
    padding: 16px 18px;
  }

  .article-outline-link {
    padding: 10px 12px;
  }
}
</style>
