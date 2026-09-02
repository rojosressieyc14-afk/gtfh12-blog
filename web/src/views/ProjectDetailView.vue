<template>
  <section v-if="project" class="project-detail-page">
    <section class="project-hero">
      <div class="project-hero__copy">
        <div class="detail-meta">
          <router-link class="status-chip published" :to="`/author/${project.author?.id}`">
            {{ project.author?.username || "匿名作者" }}
          </router-link>
          <span>{{ project.isFeatured ? "精选项目" : "项目案例" }}</span>
        </div>

        <h1 class="detail-title">{{ project.title }}</h1>
        <p class="detail-summary">{{ project.summary || "这个项目暂时还没有补充摘要。" }}</p>

        <div class="detail-stats">
          <span>{{ formatDate(project.publishedAt || project.createdAt) }}</span>
          <span>{{ project.techStacks?.length || 0 }} 项技术</span>
          <span>{{ project.highlights?.length || 0 }} 个亮点</span>
        </div>

        <div v-if="project.techStacks?.length" class="tag-row tag-row--detail">
          <span v-for="stack in project.techStacks" :key="stack" class="tag-chip"># {{ stack }}</span>
        </div>

        <div class="project-link-row project-link-row--hero">
          <a v-if="project.demoUrl" class="solid-btn" :href="project.demoUrl" target="_blank" rel="noreferrer">在线演示</a>
          <a v-if="project.repoUrl" class="ghost-btn" :href="project.repoUrl" target="_blank" rel="noreferrer">代码仓库</a>
          <router-link class="ghost-btn" to="/projects">返回作品集</router-link>
        </div>
      </div>

      <div v-if="coverUrl" class="project-hero__visual">
        <img :src="coverUrl" :alt="project.title" />
      </div>
    </section>

    <div class="detail-layout project-detail-shell">
      <section class="project-reading-layout">
        <article class="project-detail-body">
          <div class="project-detail-panel">
            <p class="eyebrow">项目快照</p>
            <div class="project-facts project-facts--grid">
              <div v-if="project.roleLabel" class="project-fact">
                <strong>{{ project.roleLabel }}</strong>
                <span>角色</span>
              </div>
              <div v-if="project.duration" class="project-fact">
                <strong>{{ project.duration }}</strong>
                <span>周期</span>
              </div>
              <div v-if="project.teamLabel" class="project-fact">
                <strong>{{ project.teamLabel }}</strong>
                <span>协作方式</span>
              </div>
            </div>
          </div>

          <div v-if="project.highlights?.length" class="project-detail-panel">
            <p class="eyebrow">亮点</p>
            <h3>这个项目最值得看的部分</h3>
            <div class="project-highlight-list project-highlight-list--detail">
              <span v-for="point in project.highlights" :key="point">· {{ point }}</span>
            </div>
          </div>

          <div v-if="project.process?.length" class="project-detail-panel">
            <p class="eyebrow">过程</p>
            <h3>项目是怎样一步步推进的</h3>
            <div class="project-highlight-list project-highlight-list--detail">
              <span v-for="point in project.process" :key="point">{{ point }}</span>
            </div>
          </div>

          <div v-if="caseTimeline.length" class="project-detail-panel">
            <p class="eyebrow">案例时间线</p>
            <h3>把项目从目标到结果串成完整故事</h3>
            <div class="case-timeline">
              <article v-for="item in caseTimeline" :key="item.title" class="case-timeline__item">
                <span class="case-timeline__badge">{{ item.kicker }}</span>
                <h4>{{ item.title }}</h4>
                <p>{{ item.description }}</p>
              </article>
            </div>
          </div>

          <div v-if="project.challenges?.length || project.solutions?.length" class="case-grid">
            <article v-if="project.challenges?.length" class="project-detail-panel">
              <p class="eyebrow">挑战</p>
              <h3>项目里的难点</h3>
              <div class="project-highlight-list project-highlight-list--detail">
                <span v-for="point in project.challenges" :key="point">· {{ point }}</span>
              </div>
            </article>

            <article v-if="project.solutions?.length" class="project-detail-panel">
              <p class="eyebrow">解决方案</p>
              <h3>我是怎么把问题拆开的</h3>
              <div class="project-highlight-list project-highlight-list--detail">
                <span v-for="point in project.solutions" :key="point">· {{ point }}</span>
              </div>
            </article>
          </div>

          <div v-if="project.results?.length" class="project-detail-panel">
            <p class="eyebrow">结果</p>
            <h3>项目最终带来了什么</h3>
            <div class="project-highlight-list project-highlight-list--detail">
              <span v-for="point in project.results" :key="point">· {{ point }}</span>
            </div>
          </div>

          <div class="case-grid">
            <article class="project-detail-panel">
              <p class="eyebrow">职责</p>
              <h3>我在这个项目里承担了什么</h3>
              <div class="project-highlight-list project-highlight-list--detail">
                <span v-for="point in roleHighlights" :key="point">· {{ point }}</span>
              </div>
            </article>

            <article class="project-detail-panel">
              <p class="eyebrow">复盘</p>
              <h3>这个案例最适合被怎样理解</h3>
              <div class="project-highlight-list project-highlight-list--detail">
                <span v-for="point in caseReflections" :key="point">· {{ point }}</span>
              </div>
            </article>
          </div>

          <article class="project-detail-panel markdown-body project-detail-body-content" v-html="html"></article>
        </article>

        <aside class="project-reading-side">
          <section class="project-detail-panel project-score-panel">
            <p class="eyebrow">案例强度</p>
            <h3>这个项目的展示完成度</h3>
            <div class="project-score-ring" :style="scoreRingStyle">
              <strong>{{ completenessScore }}%</strong>
              <span>案例完整度</span>
            </div>
            <p class="detail-summary">{{ completenessHint }}</p>
          </section>

          <section class="project-detail-panel">
            <p class="eyebrow">信息</p>
            <h3>项目信息卡</h3>
            <div class="project-facts">
              <div class="project-fact">
                <strong>{{ formatDate(project.publishedAt || project.createdAt) }}</strong>
                <span>发布时间</span>
              </div>
              <div class="project-fact">
                <strong>{{ project.techStacks?.length || 0 }}</strong>
                <span>技术标签</span>
              </div>
              <div class="project-fact">
                <strong>{{ project.highlights?.length || 0 }}</strong>
                <span>项目亮点</span>
              </div>
            </div>
          </section>

          <section class="project-detail-panel">
            <p class="eyebrow">作者</p>
            <h3>{{ project.author?.username || "匿名作者" }}</h3>
            <p class="detail-summary">进入作者页可以继续查看更多项目、文章和完整资料。</p>
            <router-link class="ghost-btn" :to="`/author/${project.author?.id}`">打开作者页</router-link>
          </section>
        </aside>
      </section>

      <section class="share-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">分享</p>
            <h3>把这个项目分享给更多人</h3>
          </div>
        </div>
        <button class="ghost-btn share-btn" @click="showShareCard = true">打开分享卡片</button>
      </section>
    </div>
  </section>

  <Teleport to="body">
    <div v-if="showShareCard" class="share-overlay" @click.self="showShareCard = false">
      <div class="share-card-modal">
        <button class="share-card-close" @click="showShareCard = false">&times;</button>

        <div class="share-card-preview">
          <div v-if="coverUrl" class="share-card-cover">
            <img :src="coverUrl" :alt="project?.title" />
          </div>
          <div class="share-card-body">
            <div class="share-card-meta">
              <span>{{ project?.author?.username || "匿名作者" }}</span>
              <span>{{ project?.isFeatured ? "精选项目" : "项目案例" }}</span>
            </div>
            <h3 class="share-card-title">{{ project?.title }}</h3>
            <p class="share-card-summary">{{ project?.summary || "这个项目暂时还没有补充摘要。" }}</p>
            <div v-if="project?.techStacks?.length" class="share-card-tags">
              <span v-for="stack in project.techStacks.slice(0, 4)" :key="stack" class="tag-chip"># {{ stack }}</span>
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
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from "vue";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { useRoute } from "vue-router";
import { getProject } from "../api/project";
import { toAssetUrl } from "../utils/asset";
import { formatDate } from "../utils/date";

const route = useRoute();
const project = ref(null);
const copied = ref(false);
const showShareCard = ref(false);

const shareUrl = computed(() => window.location.href);

const html = computed(() => DOMPurify.sanitize(marked.parse(project.value?.content || "")));
const coverUrl = computed(() => toAssetUrl(project.value?.coverImage));
const completenessScore = computed(() => {
  if (!project.value) return 0;
  const checks = [
    project.value.summary,
    project.value.roleLabel,
    project.value.duration,
    project.value.teamLabel,
    project.value.techStacks?.length,
    project.value.highlights?.length,
    project.value.process?.length,
    project.value.challenges?.length,
    project.value.solutions?.length,
    project.value.results?.length,
    project.value.content
  ];
  const done = checks.filter(Boolean).length;
  return Math.round((done / checks.length) * 100);
});

const completenessHint = computed(() => {
  if (completenessScore.value >= 90) return "这个案例已经很完整了，足够作为作品集代表项目来展示。";
  if (completenessScore.value >= 70) return "主体已经很清楚，再补一点过程细节会更有说服力。";
  return "案例还可以继续补充挑战、方案和结果，这会更能体现你的能力。";
});

const caseTimeline = computed(() => {
  if (!project.value) return [];
  return [
    {
      kicker: "目标",
      title: project.value.summary || "先把这个项目要解决的问题讲清楚",
      description: "一个好案例通常先说明为什么要做它，以及它服务的是谁、想解决什么。"
    },
    {
      kicker: "推进",
      title: project.value.process?.[0] || project.value.highlights?.[0] || "把推进过程拆成阶段",
      description: project.value.process?.[1] || "把调研、设计、实现和迭代拆开后，别人会更容易看懂你的工作方式。"
    },
    {
      kicker: "结果",
      title: project.value.results?.[0] || "最后留下可被外界理解的结果",
      description: project.value.results?.[1] || "项目结果不一定非得是数字，也可以是体验提升、结构更清晰或交付更完整。"
    }
  ];
});

const roleHighlights = computed(() => {
  if (!project.value) return [];
  const items = [];
  if (project.value.roleLabel) items.push(`主要角色：${project.value.roleLabel}`);
  if (project.value.teamLabel) items.push(`协作方式：${project.value.teamLabel}`);
  if (project.value.duration) items.push(`项目周期：${project.value.duration}`);
  if (project.value.techStacks?.length) items.push(`技术栈覆盖：${project.value.techStacks.join(" / ")}`);
  return items.length ? items : ["这个项目的职责描述还可以继续补充，比如你负责了哪些关键模块和决策。"];
});

const caseReflections = computed(() => {
  if (!project.value) return [];
  const items = [];
  if (project.value.challenges?.length) items.push("这个案例最有说服力的部分，是你如何拆解复杂问题。");
  if (project.value.solutions?.length) items.push("解决方案部分可以帮助别人理解你的判断方式，而不只是看到最终界面。");
  if (project.value.results?.length) items.push("结果部分越具体，这个项目就越像真正能代表你的作品案例。");
  if (!items.length) items.push("继续补充挑战、方案和结果，这个案例页会更像一份成熟的项目展示。");
  return items;
});

const scoreRingStyle = computed(() => ({
  background: `radial-gradient(circle at center, rgba(255, 255, 255, 0.06) 0 46%, transparent 48%), conic-gradient(from 180deg, rgba(255, 209, 102, 0.08) 0deg, rgba(255, 138, 76, 0.7) ${completenessScore.value * 3.6}deg, rgba(255, 255, 255, 0.06) 0deg)`
}));

async function loadDetail() {
  const { data } = await getProject(route.params.id);
  project.value = data.item;
}

function shareTwitter() {
  const url = encodeURIComponent(shareUrl.value);
  const text = encodeURIComponent(project.value?.title || "");
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
  const title = project.value?.title || "";
  const lines = title.match(/.{1,14}/g) || [title];
  lines.slice(0, 3).forEach((line, i) => {
    ctx.fillText(line, 32, 60 + i * 32);
  });

  ctx.fillStyle = "#6b7280";
  ctx.font = "14px sans-serif";
  ctx.fillText("PulseBlog · " + (project.value?.author?.username || ""), 32, 200);

  ctx.fillStyle = "#92400e";
  ctx.font = "12px monospace";
  ctx.fillText(shareUrl.value, 32, 240);

  const link = document.createElement("a");
  link.download = "share-card.png";
  link.href = canvas.toDataURL("image/png");
  link.click();
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

onMounted(() => {
  loadDetail();
  window.addEventListener("keydown", onEsc);
});

onUnmounted(() => {
  window.removeEventListener("keydown", onEsc);
});
</script>

<style scoped>
.project-detail-page {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.project-hero {
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

.project-hero__copy {
  display: grid;
  gap: 16px;
}

.project-hero__copy .detail-meta {
  color: #6b7280;
  font-size: 0.88rem;
}

.project-hero__copy .detail-title {
  margin: 0;
  font-size: clamp(1.8rem, 3.5vw, 3rem);
  line-height: 1.08;
  color: #1d1d1f;
}

.project-hero__copy .detail-summary {
  margin: 0;
  font-size: clamp(0.95rem, 1.4vw, 1.15rem);
  line-height: 1.6;
  color: #48484a;
}

.project-hero__copy .detail-stats {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  font-size: 0.88rem;
  color: #6b7280;
}

.project-hero__copy .tag-row {
  color: #92400e;
}

.project-hero__copy .tag-chip {
  background: rgba(180, 83, 10, 0.08);
  border-color: rgba(180, 83, 10, 0.18);
  color: #92400e;
}

.project-link-row--hero {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 8px;
}

.project-hero__visual {
  border-radius: 22px;
  overflow: hidden;
  aspect-ratio: 4 / 3;
}

.project-hero__visual img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.project-detail-shell {
  background:
    radial-gradient(circle at top right, rgba(255, 209, 102, 0.08), transparent 22%),
    rgba(255, 255, 255, 0.07);
}

.project-reading-layout {
  display: grid;
  grid-template-columns: 1fr 260px;
  gap: 40px;
  align-items: start;
}

.project-reading-side {
  position: sticky;
  top: 112px;
  display: grid;
  gap: 16px;
}

.project-detail-body {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.project-detail-body-content {
  padding: 0;
  border: none;
  background: none;
  border-radius: 0;
}

.project-score-panel {
  overflow: hidden;
}

.project-score-ring {
  display: grid;
  place-items: center;
  width: 160px;
  height: 160px;
  margin-bottom: 14px;
  border-radius: 999px;
  border: 1px solid rgba(255, 209, 102, 0.22);
}

.project-score-ring strong {
  font-size: 1.8rem;
  color: #fff1da;
}

.project-score-ring span {
  color: var(--text-soft);
  font-size: 0.85rem;
}

.case-timeline {
  display: grid;
  gap: 14px;
}

.case-timeline__item {
  padding: 18px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.03);
}

.case-timeline__item h4,
.case-timeline__item p {
  margin: 0;
}

.case-timeline__item h4 {
  margin-bottom: 10px;
}

.case-timeline__item p {
  color: var(--text-soft);
}

.case-timeline__badge {
  display: inline-flex;
  margin-bottom: 12px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 138, 76, 0.14);
  border: 1px solid rgba(255, 138, 76, 0.2);
  color: #ffd9c2;
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
  .project-reading-layout {
    grid-template-columns: 1fr;
  }
  .project-reading-side {
    position: static;
  }
}

@media (max-width: 900px) {
  .project-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .project-hero {
    padding: 28px 20px;
  }

  .project-detail-shell {
    padding: 20px;
  }

  .project-detail-panel {
    padding: 18px;
  }

  .case-timeline__item {
    padding: 14px 16px;
  }

  .project-score-ring {
    width: 140px;
    height: 140px;
  }

  .project-score-ring strong {
    font-size: 1.6rem;
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

  .project-link-row--hero {
    flex-direction: column;
  }

  .project-link-row--hero .solid-btn,
  .project-link-row--hero .ghost-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .project-hero .detail-title {
    font-size: clamp(1.4rem, 6vw, 2rem);
  }

  .project-facts--grid {
    gap: 8px;
  }

  .project-fact {
    padding: 12px 14px;
  }
}
</style>
