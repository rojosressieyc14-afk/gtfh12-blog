<template>
  <section class="project-editor-page">
    <div class="project-editor-main panel-card">
      <header class="editor-head">
        <div>
          <p class="eyebrow">项目案例编辑</p>
          <h2>{{ isEdit ? "继续完善项目案例" : "新建项目案例" }}</h2>
          <p class="detail-summary">
            先保存草稿，再在内容完整后提交审核。审核通过后，项目才会公开展示在作品列表里。
          </p>
        </div>

        <div class="editor-status">
          <span class="status-chip" :class="currentStatus">{{ statusLabel(currentStatus) }}</span>
          <p class="table-note">{{ statusHint }}</p>
        </div>
      </header>

      <div v-if="rejectReason" class="reject-banner">
        <strong>驳回原因</strong>
        <p>{{ rejectReason }}</p>
      </div>

      <div class="action-strip">
        <button class="ghost-btn" type="button" :disabled="saving || submitting" @click="saveProject">
          {{ saving ? "保存中..." : "保存草稿" }}
        </button>
        <button class="solid-btn" type="button" :disabled="saving || submitting" @click="submitForReview">
          {{ submitting ? "提交中..." : submitButtonLabel }}
        </button>
      </div>
      <p v-if="errorMessage" class="error-text" style="margin-top:10px">{{ errorMessage }}</p>

      <div class="stack-form">
        <label>
          项目名称
          <input v-model.trim="form.title" class="field-input" placeholder="例如：个人作品集与内容平台" />
        </label>

        <label>
          项目摘要
          <textarea
            v-model.trim="form.summary"
            class="field-area field-area--small"
            placeholder="用一两句话说明项目定位、亮点和解决的问题。"
          ></textarea>
        </label>

        <div class="editor-meta-grid">
          <label>
            你的角色
            <input v-model.trim="form.roleLabel" class="field-input" placeholder="例如：全栈开发 / 独立开发" />
          </label>
          <label>
            项目周期
            <input v-model.trim="form.duration" class="field-input" placeholder="例如：6 周 / 2026 Q2" />
          </label>
        </div>

        <div class="editor-meta-grid">
          <label>
            团队信息
            <input v-model.trim="form.teamLabel" class="field-input" placeholder="例如：个人项目 / 双人协作" />
          </label>
          <label>
            封面图片
            <input class="field-input" type="file" accept="image/*" @change="handleUpload" />
          </label>
        </div>

        <div v-if="coverUrl" class="editor-cover-preview">
          <img :src="coverUrl" alt="项目封面" />
        </div>

        <div class="editor-meta-grid">
          <label>
            技术栈
            <input v-model.trim="techStacksInput" class="field-input" placeholder="Go, Gin, Vue 3, MySQL" />
          </label>
          <label>
            项目亮点
            <input v-model.trim="highlightsInput" class="field-input" placeholder="审核流、实时预览、双端后台" />
          </label>
        </div>

        <div class="editor-meta-grid">
          <label>
            推进过程
            <textarea
              v-model.trim="processInput"
              class="field-area field-area--small"
              placeholder="需求梳理&#10;信息架构&#10;开发实现"
            ></textarea>
          </label>
          <label>
            遇到的挑战
            <textarea
              v-model.trim="challengesInput"
              class="field-area field-area--small"
              placeholder="权限控制&#10;审核流转&#10;编码问题"
            ></textarea>
          </label>
        </div>

        <div class="editor-meta-grid">
          <label>
            解决方案
            <textarea
              v-model.trim="solutionsInput"
              class="field-area field-area--small"
              placeholder="拆分模块&#10;统一状态流&#10;补齐后台审核"
            ></textarea>
          </label>
          <label>
            最终结果
            <textarea
              v-model.trim="resultsInput"
              class="field-area field-area--small"
              placeholder="形成完整作品案例&#10;支持审核闭环"
            ></textarea>
          </label>
        </div>

        <div class="editor-meta-grid">
          <label>
            在线地址
            <input v-model.trim="form.demoUrl" class="field-input" placeholder="https://demo.example.com" />
          </label>
          <label>
            仓库地址
            <input v-model.trim="form.repoUrl" class="field-input" placeholder="https://github.com/your/repo" />
          </label>
        </div>

        <div class="editor-meta-grid">
          <label>
            排序优先级
            <input v-model.number="form.sortOrder" class="field-input" type="number" min="0" />
          </label>
          <label class="project-switch">
            <span>设为首页精选</span>
            <input v-model="form.isFeatured" type="checkbox" />
          </label>
        </div>

        <div class="privacy-toggle">
          <label class="toggle-label">
            <span>可见范围</span>
            <div class="toggle-group">
              <button type="button" class="toggle-btn" :class="{ active: !form.isPrivate }" @click="form.isPrivate = false">公开发布</button>
              <button type="button" class="toggle-btn" :class="{ active: form.isPrivate }" @click="form.isPrivate = true">仅自己可见</button>
            </div>
          </label>
        </div>

        <label>
          项目正文
          <textarea
            v-model="form.content"
            class="field-area field-area--editor"
            placeholder="# 项目背景&#10;&#10;## 目标与约束&#10;&#10;## 设计与实现&#10;&#10;## 难点拆解&#10;&#10;## 最终结果"
          ></textarea>
        </label>
      </div>
    </div>

    <aside class="preview-card panel-card case-preview">
      <p class="eyebrow">案例预览</p>
      <h3>{{ form.title || "这里会显示你的项目标题" }}</h3>
      <p class="detail-summary">
        {{ form.summary || "把这个项目整理成完整案例，别人会更容易理解你的能力边界和思考方式。" }}
      </p>

      <div class="project-facts project-facts--preview">
        <div v-if="form.roleLabel" class="project-fact">
          <strong>{{ form.roleLabel }}</strong>
          <span>角色定位</span>
        </div>
        <div v-if="form.duration" class="project-fact">
          <strong>{{ form.duration }}</strong>
          <span>项目周期</span>
        </div>
        <div v-if="form.teamLabel" class="project-fact">
          <strong>{{ form.teamLabel }}</strong>
          <span>团队信息</span>
        </div>
      </div>

      <div v-if="coverUrl" class="preview-cover">
        <img :src="coverUrl" alt="项目封面预览" />
      </div>

      <div v-if="techStacks.length" class="tag-row">
        <span v-for="stack in techStacks" :key="stack" class="tag-chip"># {{ stack }}</span>
      </div>

      <section class="case-summary-grid">
        <article class="case-summary-card">
          <span>亮点</span>
          <strong>{{ highlights.length }}</strong>
        </article>
        <article class="case-summary-card">
          <span>挑战</span>
          <strong>{{ challenges.length }}</strong>
        </article>
        <article class="case-summary-card">
          <span>结果</span>
          <strong>{{ results.length }}</strong>
        </article>
      </section>

      <div class="case-checklist">
        <div>
          <span class="case-checklist__label">案例完整度</span>
          <strong>{{ checklistText }}</strong>
        </div>
        <p>{{ checklistHint }}</p>
      </div>

      <div class="markdown-body article-detail-body" v-html="compiledMarkdown"></div>
    </aside>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { createProject, getMyProject, submitProject, updateProject } from "../api/project";
import { uploadImage } from "../api/upload";
import { toAssetUrl } from "../utils/asset";
import { useUserStore } from "../stores/user";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const isEdit = computed(() => Boolean(route.params.id));
const saving = ref(false);
const submitting = ref(false);
const errorMessage = ref("");
const currentStatus = ref("draft");
const rejectReason = ref("");
const techStacksInput = ref("");
const highlightsInput = ref("");
const processInput = ref("");
const challengesInput = ref("");
const solutionsInput = ref("");
const resultsInput = ref("");

const form = reactive({
  title: "",
  summary: "",
  roleLabel: "",
  duration: "",
  teamLabel: "",
  content: "",
  coverImage: "",
  demoUrl: "",
  repoUrl: "",
  sortOrder: 0,
  isPrivate: false,
  isFeatured: false
});

const techStacks = computed(() => parseInputList(techStacksInput.value));
const highlights = computed(() => parseInputList(highlightsInput.value));
const process = computed(() => parseInputList(processInput.value));
const challenges = computed(() => parseInputList(challengesInput.value));
const solutions = computed(() => parseInputList(solutionsInput.value));
const results = computed(() => parseInputList(resultsInput.value));
const compiledMarkdown = computed(() => DOMPurify.sanitize(marked.parse(form.content || "## 这里会实时预览你的项目正文")));
const coverUrl = computed(() => toAssetUrl(form.coverImage));
const submitButtonLabel = computed(() => (userStore.isAdmin ? "直接发布" : "提交审核"));

const completenessScore = computed(() => {
  let score = 0;
  if (form.title.trim()) score += 1;
  if (form.summary.trim()) score += 1;
  if (form.roleLabel.trim()) score += 1;
  if (techStacks.value.length) score += 1;
  if (highlights.value.length) score += 1;
  if (process.value.length) score += 1;
  if (challenges.value.length) score += 1;
  if (solutions.value.length) score += 1;
  if (results.value.length) score += 1;
  if (form.content.trim()) score += 1;
  return score;
});

const checklistText = computed(() => {
  if (completenessScore.value >= 9) return "已经比较完整，可以提交审核了。";
  if (completenessScore.value >= 6) return "主体已经成型，再补一些细节会更强。";
  return "还在搭骨架，建议继续完善关键信息。";
});

const checklistHint = computed(() => {
  if (!process.value.length) return "建议补充推进过程，让别人看到你是怎么一步步做出来的。";
  if (!challenges.value.length) return "挑战与难点还没写，这一段通常最能体现你的判断力。";
  if (!results.value.length) return "再补上结果与收益，案例会更有说服力。";
  return "现在已经不只是项目记录了，更像一篇成熟的作品案例。";
});

const statusHint = computed(() => {
  if (currentStatus.value === "draft") return "保存后仍是草稿，提交审核后才会进入后台审核队列。";
  if (currentStatus.value === "pending") return "项目正在等待管理员审核。";
  if (currentStatus.value === "published") return "项目已公开展示。再次编辑后会回到草稿。";
  if (currentStatus.value === "rejected") return "项目已被驳回，修改后可以重新提交。";
  return "";
});

function statusLabel(status) {
  return {
    draft: "草稿",
    pending: "审核中",
    published: "已发布",
    rejected: "已驳回"
  }[status] || status;
}

function buildPayload() {
  return {
    ...form,
    techStacks: techStacks.value,
    highlights: highlights.value,
    process: process.value,
    challenges: challenges.value,
    solutions: solutions.value,
    results: results.value
  };
}

function privacyLabel() {
  return form.isPrivate ? "仅自己可见" : "公开发布";
}

function parseInputList(value) {
  return value
    .split(/[,\n]/)
    .map((item) => item.trim())
    .filter((item, index, list) => item && list.indexOf(item) === index);
}

function syncStatus(item) {
  currentStatus.value = item.status || "draft";
  rejectReason.value = item.rejectReason || "";
}

async function loadDetail() {
  if (!isEdit.value) return;
  const { data } = await getMyProject(route.params.id);
  Object.assign(form, {
    title: data.item.title || "",
    summary: data.item.summary || "",
    roleLabel: data.item.roleLabel || "",
    duration: data.item.duration || "",
    teamLabel: data.item.teamLabel || "",
    content: data.item.content || "",
    coverImage: data.item.coverImage || "",
    demoUrl: data.item.demoUrl || "",
    repoUrl: data.item.repoUrl || "",
    sortOrder: data.item.sortOrder || 0,
    isPrivate: Boolean(data.item.isPrivate),
    isFeatured: Boolean(data.item.isFeatured)
  });
  syncStatus(data.item);
  techStacksInput.value = data.item.techStacks?.join(", ") || "";
  highlightsInput.value = data.item.highlights?.join(", ") || "";
  processInput.value = data.item.process?.join("\n") || "";
  challengesInput.value = data.item.challenges?.join("\n") || "";
  solutionsInput.value = data.item.solutions?.join("\n") || "";
  resultsInput.value = data.item.results?.join("\n") || "";
}

async function saveProject() {
  saving.value = true;
  errorMessage.value = "";
  try {
    if (isEdit.value) {
      const { data } = await updateProject(route.params.id, buildPayload());
      syncStatus(data.item);
    } else {
      const { data } = await createProject(buildPayload());
      syncStatus(data.item);
      await router.replace(`/project-editor/${data.item.id}`);
    }
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || "保存失败，请稍后再试";
  } finally {
    saving.value = false;
  }
}

async function submitForReview() {
  submitting.value = true;
  errorMessage.value = "";
  try {
    let projectId = route.params.id;
    if (isEdit.value) {
      const { data } = await updateProject(route.params.id, buildPayload());
      syncStatus(data.item);
    } else {
      const { data } = await createProject(buildPayload());
      syncStatus(data.item);
      projectId = data.item.id;
    }
    const { data } = await submitProject(projectId);
    syncStatus(data.item);
    router.push("/my-projects");
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || "提交失败，请稍后再试";
  } finally {
    submitting.value = false;
  }
}

async function handleUpload(event) {
  const [file] = event.target.files || [];
  if (!file) return;
  const { data } = await uploadImage(file);
  form.coverImage = data.url;
}

onMounted(loadDetail);
</script>

<style scoped>
.project-editor-page {
  margin-top: 28px;
  display: grid;
  grid-template-columns: minmax(0, 1.12fr) minmax(320px, 0.88fr);
  gap: 24px;
}

.editor-head {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  margin-bottom: 22px;
}

.editor-status {
  width: min(260px, 100%);
}

.reject-banner {
  margin-bottom: 20px;
  padding: 16px 18px;
  border-radius: 22px;
  border: 1px solid rgba(248, 113, 113, 0.22);
  background: rgba(127, 29, 29, 0.22);
}

.reject-banner p {
  margin: 8px 0 0;
  color: #fecaca;
}

.action-strip {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 22px;
}

.project-switch {
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid var(--border);
  background: rgba(8, 11, 17, 0.48);
  justify-content: space-between;
}

.editor-cover-preview,
.preview-cover {
  overflow: hidden;
  border-radius: 26px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
}

.editor-cover-preview img,
.preview-cover img {
  width: 100%;
  display: block;
  object-fit: cover;
}

.editor-cover-preview img {
  max-height: 320px;
}

.case-summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin: 18px 0;
}

.case-summary-card {
  padding: 16px 18px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(8, 11, 17, 0.38);
}

.case-summary-card span,
.case-checklist__label {
  display: block;
  margin-bottom: 6px;
  color: var(--text-soft);
  font-size: 0.84rem;
}

.case-checklist {
  padding: 18px 20px;
  border-radius: 24px;
  border: 1px solid rgba(255, 209, 102, 0.2);
  background: rgba(255, 209, 102, 0.08);
  margin-bottom: 20px;
}

.case-checklist strong {
  display: block;
  font-size: 1rem;
  color: #fff0d5;
}

.case-checklist p {
  margin: 10px 0 0;
  color: #fff0d5;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
  transform: none;
}

.privacy-toggle {
  margin-bottom: 16px;
}
.toggle-label {
  display: flex;
  align-items: center;
  gap: 12px;
}
.toggle-group {
  display: flex;
  border-radius: 12px;
  border: 1px solid var(--border, rgba(255,255,255,0.12));
  overflow: hidden;
}
.toggle-btn {
  padding: 8px 16px;
  border: none;
  background: rgba(255,255,255,0.04);
  color: var(--soft, rgba(242,239,232,0.7));
  cursor: pointer;
  font: inherit;
  font-size: 0.88rem;
  transition: all 0.2s;
}
.toggle-btn.active {
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 600;
}

@media (max-width: 960px) {
  .project-editor-page {
    grid-template-columns: 1fr;
  }

  .case-summary-grid,
  .editor-head {
    grid-template-columns: 1fr;
    display: grid;
  }
}

@media (max-width: 768px) {
  .project-editor-page {
    margin-top: 14px;
    gap: 16px;
  }

  .action-strip {
    flex-direction: column;
  }

  .action-strip button {
    width: 100%;
    justify-content: center;
  }

  .case-summary-card {
    padding: 12px 14px;
  }

  .case-checklist {
    padding: 14px 16px;
  }

  .project-switch {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .project-editor-main {
    padding: 18px;
  }

  .reject-banner {
    padding: 12px 14px;
  }
}
</style>
