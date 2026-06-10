<template>
  <section v-if="project" class="project-detail-layout panel-card project-detail-layout--brand">
    <div class="project-detail-hero">
      <div class="project-detail-copy">
        <div class="detail-meta">
          <router-link class="status-chip published" :to="`/author/${project.author?.id}`">
            {{ project.author?.username || "匿名作者" }}
          </router-link>
          <span>{{ project.isFeatured ? "精选项目" : "项目案例" }}</span>
        </div>

        <h2 class="detail-title">{{ project.title }}</h2>
        <p class="detail-summary">{{ project.summary || "这个项目暂时还没有补充摘要。" }}</p>

        <div v-if="project.techStacks?.length" class="tag-row tag-row--detail">
          <span v-for="stack in project.techStacks" :key="stack" class="tag-chip"># {{ stack }}</span>
        </div>

        <div class="project-link-row project-link-row--hero">
          <a v-if="project.demoUrl" class="solid-btn" :href="project.demoUrl" target="_blank" rel="noreferrer">在线演示</a>
          <a v-if="project.repoUrl" class="ghost-btn" :href="project.repoUrl" target="_blank" rel="noreferrer">代码仓库</a>
          <router-link class="ghost-btn" to="/projects">返回作品集</router-link>
        </div>
      </div>

      <div class="project-detail-visual">
        <div class="project-detail-constellation">
          <div v-if="coverUrl" class="project-detail-cover">
            <img :src="coverUrl" :alt="project.title" />
          </div>
        </div>
      </div>
    </div>

    <section class="project-detail-grid">
      <article class="project-detail-main">
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

        <article class="project-detail-panel markdown-body" v-html="html"></article>
      </article>

      <aside class="project-detail-side">
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
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { useRoute } from "vue-router";
import { getProject } from "../api/project";
import { toAssetUrl } from "../utils/asset";

const route = useRoute();
const project = ref(null);

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

function formatDate(value) {
  return new Date(value).toLocaleDateString("zh-CN");
}

onMounted(loadDetail);
</script>

<style scoped>
.project-detail-layout--brand {
  background:
    radial-gradient(circle at 84% 18%, rgba(255, 209, 102, 0.12), transparent 24%),
    radial-gradient(circle at 12% 76%, rgba(255, 138, 76, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.07);
}

.project-detail-visual {
  position: relative;
}

.project-detail-constellation {
  position: relative;
  min-height: 380px;
}

.project-score-panel {
  overflow: hidden;
}

.project-score-ring {
  display: grid;
  place-items: center;
  width: 180px;
  height: 180px;
  margin-bottom: 18px;
  border-radius: 999px;
  border: 1px solid rgba(255, 209, 102, 0.22);
}

.project-score-ring strong {
  font-size: 2rem;
  color: #fff1da;
}

.project-score-ring span {
  color: var(--text-soft);
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

@media (max-width: 768px) {
  .project-detail-constellation {
    min-height: 260px;
  }

  .project-score-ring {
    width: 140px;
    height: 140px;
  }

  .project-score-ring strong {
    font-size: 1.6rem;
  }

  .case-timeline__item {
    padding: 14px 16px;
  }

  .project-detail-panel {
    padding: 18px;
  }

  .project-detail-layout--brand {
    padding: 20px;
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
  .project-detail-hero .detail-title {
    font-size: clamp(1.4rem, 6vw, 2rem);
  }

  .project-detail-visual {
    display: none;
  }

  .project-facts--grid {
    gap: 8px;
  }

  .project-fact {
    padding: 12px 14px;
  }

  .project-detail-constellation {
    min-height: auto;
  }
}
</style>
