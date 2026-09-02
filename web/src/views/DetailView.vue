<template>
  <section v-if="article" class="article-detail-page">
    <div class="reading-progress-bar" :style="{ width: readProgress + '%' }"></div>
    <div class="detail-layout article-detail-shell">
      <div class="article-detail-header">
        <div class="detail-meta">
          <router-link class="status-chip published" :to="`/author/${article.author?.id}`">
            {{ article.author?.username || "匿名作者" }}
          </router-link>
          <span>{{ article.category?.name || "未分类" }}</span>
        </div>

        <h1 class="detail-title">{{ article.title }}</h1>
        <p class="detail-summary">{{ article.summary || "暂无摘要" }}</p>

        <div class="detail-stats">
          <span>{{ formatDate(article.publishedAt || article.createdAt) }}</span>
          <span>{{ commentCount }} 条评论</span>
          <span>{{ article.viewCount || 0 }} 次阅读</span>
          <span>{{ readingMinutes }} 分钟阅读</span>
        </div>

        <div v-if="article.tags?.length" class="tag-row tag-row--detail">
          <span v-for="tag in article.tags" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
        </div>

        <div v-if="coverUrl" class="detail-cover article-detail-cover">
          <img :src="coverUrl" :alt="article.title" />
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

      <section class="article-reading-layout">
        <article v-highlight class="markdown-body article-detail-body" v-html="html"></article>

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



          <section class="project-detail-panel">
            <p class="eyebrow">作者入口</p>
            <h3>继续了解这位作者</h3>
            <p class="detail-summary">你可以继续查看这位作者的项目案例、公开文章和完整个人资料。</p>
            <router-link class="ghost-btn" :to="`/author/${article.author?.id}`">打开作者页</router-link>
          </section>
        </aside>
      </section>

      <section class="share-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">分享</p>
            <h3>把这篇文章分享给更多人</h3>
          </div>
        </div>
        <button class="ghost-btn share-btn" @click="showShareCard = true">打开分享卡片</button>
      </section>

      <Teleport to="body">
        <div v-if="showShareCard" class="share-overlay" @click.self="showShareCard = false">
          <div class="share-card-modal">
            <button class="share-card-close" @click="showShareCard = false">&times;</button>

            <div class="share-card-preview">
              <div v-if="coverUrl" class="share-card-cover">
                <img :src="coverUrl" :alt="article?.title" />
              </div>
              <div class="share-card-body">
                <div class="share-card-meta">
                  <span>{{ article?.author?.username || "匿名作者" }}</span>
                  <span>{{ article?.category?.name || "未分类" }}</span>
                </div>
                <h3 class="share-card-title">{{ article?.title }}</h3>
                <p class="share-card-summary">{{ article?.summary || "这篇文章还没有摘要。" }}</p>
                <div v-if="article?.tags?.length" class="share-card-tags">
                  <span v-for="tag in article.tags.slice(0, 3)" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
                </div>
                <div class="share-card-url">
                  <input readonly :value="shareUrl" class="field-input share-card-url-input" @click="$event.target.select()" />
                  <button class="ghost-btn share-card-copy" @click="copyShareUrl">
                    {{ copied ? "已复制 ✓" : "复制链接" }}
                  </button>
                </div>
              </div>
            </div>

            <div class="share-card-actions">
              <button class="ghost-btn share-action-btn" @click="shareTwitter">
                <span class="share-action-icon">𝕏</span>
                <span>Twitter / X</span>
              </button>
              <button class="ghost-btn share-action-btn" @click="shareLinkedIn">
                <span class="share-action-icon">in</span>
                <span>LinkedIn</span>
              </button>
              <button class="ghost-btn share-action-btn" @click="shareWeChat">
                <span class="share-action-icon">微信</span>
                <span>保存图片分享</span>
              </button>
            </div>
          </div>
        </div>
      </Teleport>

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
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { useRoute } from "vue-router";
import { createComment, getArticle, listComments, toggleFavorite, toggleLike } from "../api/article";
import { useUserStore } from "../stores/user";
import { toAssetUrl } from "../utils/asset";
import { formatDate } from "../utils/date";

const route = useRoute();
const userStore = useUserStore();
const article = ref(null);
const comments = ref([]);
const commentText = ref("");
const replyTo = ref(null);
const replyText = ref("");
const readProgress = ref(0);
const copied = ref(false);
const showShareCard = ref(false);

const shareUrl = computed(() => window.location.href);

let scrollHandler = null;

watch(() => article.value?.title, (t) => {
  document.title = t ? `${t} — PulseBlog` : "PulseBlog";
}, { immediate: true });

function updateReadProgress() {
  const scrollTop = window.scrollY;
  const docHeight = document.documentElement.scrollHeight - window.innerHeight;
  readProgress.value = docHeight > 0 ? Math.min(100, Math.round((scrollTop / docHeight) * 100)) : 0;
}

function stripInlineMarkdown(text) {
  return text
    .replace(/\*\*(.+?)\*\*/g, "$1")
    .replace(/\*(.+?)\*/g, "$1")
    .replace(/`(.+?)`/g, "$1")
    .replace(/\[(.+?)\]\(.+?\)/g, "$1");
}

function buildHeadingId(text, index) {
  const safe = stripInlineMarkdown(text)
    .toLowerCase()
    .replace(/<[^>]+>/g, "")
    .replace(/[^\w\u4e00-\u9fa5\s-]/g, "")
    .trim()
    .replace(/\s+/g, "-");
  return `heading-${safe || "section"}-${index}`;
}

const articleOutline = computed(() => {
  const matches = [...(article.value?.content || "").matchAll(/^(#{1,3})\s+(.+)$/gm)];
  return matches.map((match, index) => {
    const text = stripInlineMarkdown(match[2].trim());
    return {
      level: match[1].length,
      text,
      id: buildHeadingId(text, index)
    };
  });
});

const html = computed(() => {
  const renderer = new marked.Renderer();
  renderer.heading = (token) => {
    const text = token.text || "";
    const id = articleOutline.value.find((item) => item.text === text)?.id || buildHeadingId(text, 0);
    return `<h${token.depth} id="${id}">${text}</h${token.depth}>`;
  };
  return DOMPurify.sanitize(marked.parse(article.value?.content || "", { renderer }));
});

const coverUrl = computed(() => toAssetUrl(article.value?.coverImage));
const commentCount = computed(() => comments.value.reduce((sum, item) => sum + 1 + (item.replies?.length || 0), 0));
const readingMinutes = computed(() => {
  const count = (article.value?.content || "").replace(/\s+/g, "").length;
  return Math.max(1, Math.ceil(count / 500));
});

async function loadDetail() {
  try {
    const { data } = await getArticle(route.params.id);
    article.value = data.item;
  } catch (e) {
    article.value = null;
  }
}

async function loadComments() {
  try {
    const { data } = await listComments(route.params.id);
    comments.value = data.items;
  } catch (e) {
    comments.value = [];
  }
}

async function submitComment() {
  if (!commentText.value) return;
  try {
    await createComment(route.params.id, { content: commentText.value });
    commentText.value = "";
    await loadComments();
  } catch { /* handled by axios interceptor */ }
}

function shareTwitter() {
  const url = encodeURIComponent(shareUrl.value);
  const text = encodeURIComponent(article.value?.title || "");
  window.open(`https://twitter.com/intent/tweet?text=${text}&url=${url}`, "_blank", "noopener");
}

function shareLinkedIn() {
  const url = encodeURIComponent(shareUrl.value);
  window.open(`https://www.linkedin.com/sharing/share-offsite/?url=${url}`, "_blank", "noopener");
}

function shareWeChat() {
  const canvas = document.createElement("canvas");
  canvas.width = 480;
  canvas.height = 320;
  const ctx = canvas.getContext("2d");

  ctx.fillStyle = "#f4f4f6";
  ctx.fillRect(0, 0, 480, 320);

  ctx.fillStyle = "#1d1d1f";
  ctx.font = "bold 22px sans-serif";
  const title = article.value?.title || "";
  const lines = title.match(/.{1,14}/g) || [title];
  lines.slice(0, 3).forEach((line, i) => {
    ctx.fillText(line, 32, 60 + i * 32);
  });

  ctx.fillStyle = "#6b7280";
  ctx.font = "14px sans-serif";
  ctx.fillText("PulseBlog · " + (article.value?.author?.username || ""), 32, 200);

  ctx.fillStyle = "#92400e";
  ctx.font = "12px monospace";
  ctx.fillText(shareUrl.value, 32, 240);

  const link = document.createElement("a");
  link.download = "share-card.png";
  link.href = canvas.toDataURL("image/png");
  link.click();
}

async function submitReply(parentId) {
  if (!replyText.value) return;
  try {
    await createComment(route.params.id, { content: replyText.value, parentId });
    replyText.value = "";
    replyTo.value = null;
    await loadComments();
  } catch { /* handled by axios interceptor */ }
}

async function onLike() {
  if (!userStore.isLoggedIn) return;
  try {
    const { data } = await toggleLike(route.params.id);
    article.value = { ...article.value, ...data.item };
  } catch { /* handled by axios interceptor */ }
}

async function onFavorite() {
  if (!userStore.isLoggedIn) return;
  try {
    const { data } = await toggleFavorite(route.params.id);
    article.value = { ...article.value, ...data.item };
  } catch { /* handled by axios interceptor */ }
}

function scrollToHeading(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: "smooth", block: "start" });
}

async function copyShareUrl() {
  try {
    await navigator.clipboard.writeText(shareUrl.value);
  } catch {
    const ta = document.createElement("textarea");
    ta.value = shareUrl.value;
    ta.style.position = "fixed";
    ta.style.opacity = "0";
    document.body.appendChild(ta);
    ta.select();
    document.execCommand("copy");
    document.body.removeChild(ta);
  }
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
}

function onEsc(e) {
  if (e.key === "Escape" && showShareCard.value) showShareCard.value = false;
}

onMounted(async () => {
  await loadDetail();
  await loadComments();
  scrollHandler = () => requestAnimationFrame(updateReadProgress);
  window.addEventListener("scroll", scrollHandler, { passive: true });
  window.addEventListener("keydown", onEsc);
});

onUnmounted(() => {
  if (scrollHandler) window.removeEventListener("scroll", scrollHandler);
  window.removeEventListener("keydown", onEsc);
});
</script>

<style scoped>
.reading-progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 2px;
  background: var(--accent);
  z-index: 9999;
  transition: width 0.1s linear;
}
.article-detail-shell {
  background: var(--panel);
  border-radius: 32px;
  padding: 42px;
  border: 1px solid var(--border);
}

.article-detail-header {
  max-width: 720px;
  margin: 0 auto 48px;
  text-align: center;
}

.article-reading-layout {
  display: grid;
  grid-template-columns: 1fr 220px;
  gap: 40px;
  align-items: start;
}

.article-reading-side {
  position: sticky;
  top: 112px;
  display: grid;
  gap: 16px;
}

.article-outline-list {
  display: grid;
  gap: 6px;
}

.article-outline-link {
  display: block;
  width: 100%;
  text-align: left;
  border: none;
  background: none;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 13px;
  color: rgba(246,241,234,0.5);
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
  line-height: 1.4;
}
.article-outline-link:hover {
  color: rgba(246,241,234,0.8);
  background: rgba(255,255,255,0.04);
}
.article-outline-link--2 {
  padding-left: 20px;
  font-size: 12px;
}
.article-outline-link--3 {
  padding-left: 28px;
  font-size: 11px;
}

.share-panel {
  margin-top: 48px;
}
.share-btn {
  font-size: 13px;
}

.share-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  animation: fadeIn 0.2s ease;
}

.share-card-modal {
  position: relative;
  width: min(480px, 100%);
  background: var(--panel, #1a1e26);
  border: 1px solid var(--border, rgba(255,255,255,0.1));
  border-radius: 24px;
  overflow: hidden;
  animation: slideUp 0.3s ease;
}

.share-card-close {
  position: absolute;
  top: 12px;
  right: 14px;
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255,255,255,0.08);
  color: var(--text, #f7f3ea);
  border-radius: 50%;
  font-size: 1.3rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  transition: background 0.15s;
}
.share-card-close:hover {
  background: rgba(255,255,255,0.15);
}

.share-card-cover {
  height: 180px;
  overflow: hidden;
}
.share-card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.share-card-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.share-card-meta {
  display: flex;
  gap: 12px;
  font-size: 0.85rem;
  color: var(--text-soft, rgba(246,241,234,0.72));
}

.share-card-title {
  margin: 0;
  font-size: 1.2rem;
  line-height: 1.35;
  color: var(--text, #f7f3ea);
}

.share-card-summary {
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.55;
  color: var(--text-soft, rgba(246,241,234,0.72));
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.share-card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.share-card-url {
  display: flex;
  gap: 8px;
  align-items: center;
}
.share-card-url-input {
  flex: 1;
  font-size: 0.82rem;
  border-radius: 10px;
  cursor: pointer;
}
.share-card-copy {
  white-space: nowrap;
  font-size: 0.82rem;
}

.share-card-actions {
  display: flex;
  border-top: 1px solid var(--border, rgba(255,255,255,0.1));
}
.share-action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 14px 8px;
  border-radius: 0;
  font-size: 0.85rem;
  border-right: 1px solid var(--border, rgba(255,255,255,0.1));
}
.share-action-btn:last-child {
  border-right: none;
}
.share-action-btn:hover {
  background: rgba(255,255,255,0.04);
}
.share-action-icon {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: rgba(255,255,255,0.08);
  font-size: 0.75rem;
  font-weight: 700;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@media (max-width: 1080px) {
  .article-reading-layout {
    grid-template-columns: 1fr;
  }
  .article-reading-side {
    position: static;
  }
}

@media (max-width: 768px) {
  .article-detail-shell {
    padding: 20px;
  }
  .detail-stats {
    flex-wrap: wrap;
  }
  .reaction-row {
    flex-direction: column;
  }
  .reaction-btn {
    width: 100%;
    justify-content: center;
  }
  .share-card-actions {
    flex-direction: column;
  }
  .share-action-btn {
    border-right: none;
    border-bottom: 1px solid var(--border, rgba(255,255,255,0.1));
  }
  .share-action-btn:last-child {
    border-bottom: none;
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
}
</style>
