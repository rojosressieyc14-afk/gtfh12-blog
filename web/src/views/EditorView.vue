<template>
  <section class="editor-workbench">
    <div class="editor-main panel-card">
      <header class="editor-hero">
        <div class="editor-hero__copy">
          <p class="eyebrow">创作台</p>
          <h2>{{ isEdit ? "继续打磨这篇文章" : "开始一篇新的公开输出" }}</h2>
          <p class="detail-summary">
            把它写成一篇真正能代表你的内容。可以是学习笔记、项目复盘、源码阅读，或者一篇值得长期被引用的经验总结。
          </p>
        </div>

        <div class="editor-hero__status">
          <article class="status-panel">
            <span>当前状态</span>
            <strong>{{ draftStateLabel }}</strong>
          </article>
          <article class="status-panel">
            <span>发布方式</span>
            <strong>{{ submitModeLabel }}</strong>
          </article>
        </div>
      </header>

      <div v-if="showRestoreNotice" class="draft-banner">
        <div>
          <strong>已恢复本地草稿</strong>
          <p>刚刚帮你找回了尚未同步的内容。你可以继续完善，也可以直接丢弃。</p>
        </div>
        <button class="ghost-btn" type="button" @click="dismissRestoreNotice">知道了</button>
      </div>

      <div class="action-strip">
        <button
          v-if="hasLocalDraft"
          class="ghost-btn"
          type="button"
          :disabled="saving || submitting || autosaving"
          @click="discardLocalDraft"
        >
          丢弃本地草稿
        </button>
        <button class="ghost-btn" type="button" :disabled="saving || submitting" @click="saveArticle">
          {{ saving ? "正在保存..." : "保存草稿" }}
        </button>
        <button class="solid-btn" type="button" :disabled="saving || submitting" @click="submitForReview">
          {{ submitting ? "正在提交..." : submitLabel }}
        </button>
      </div>
      <p v-if="errorMessage" class="error-text" style="margin-top:10px">{{ errorMessage }}</p>

      <div class="editor-body">
        <div class="stack-form">
          <label>
            文章标题
            <input
              v-model.trim="form.title"
              class="field-input"
              placeholder="例如：从 0 到 1 搭建一个作品集内容平台"
            />
          </label>

          <label>
            内容摘要
            <textarea
              v-model.trim="form.summary"
              class="field-area field-area--small"
              placeholder="用一两句话说明这篇文章解决了什么问题，为什么值得读。"
            ></textarea>
          </label>

          <div class="editor-meta-grid">
            <label>
              所属分类
              <select v-model="form.categoryId" class="field-input">
                <option :value="null">请选择分类</option>
                <option v-for="item in categories" :key="item.id" :value="item.id">{{ item.name }}</option>
              </select>
            </label>

            <label>
              文章标签
              <input v-model.trim="tagsInput" class="field-input" placeholder="例如：Go, Gin, Vue3, 学习记录" />
            </label>
          </div>

          <div v-if="parsedTags.length" class="tag-row">
            <span v-for="tag in parsedTags" :key="tag" class="tag-chip"># {{ tag }}</span>
          </div>

          <div class="editor-meta-grid">
            <label>
              封面图片
              <input class="field-input" type="file" accept="image/*" @change="handleUpload" />
            </label>

            <article class="hint-card">
              <span>写作建议</span>
              <strong>标题清楚、摘要具体、正文分层，审核和展示效果都会更好。</strong>
            </article>
          </div>

          <div v-if="coverUrl" class="editor-cover-preview">
            <img :src="coverUrl" alt="文章封面" />
          </div>

          <label>
            Markdown 正文
            <textarea
              v-model="form.content"
              class="field-area field-area--editor"
              placeholder="# 先写结论&#10;&#10;## 背景&#10;&#10;## 方案&#10;&#10;## 过程&#10;&#10;## 总结"
            ></textarea>
          </label>
        </div>
      </div>
    </div>

    <aside class="preview-card panel-card editor-preview">
      <div class="preview-head">
        <p class="eyebrow">实时预览</p>
        <h3>{{ form.title || "这篇文章还没有标题" }}</h3>
        <p class="detail-summary">
          {{ form.summary || "写一段简洁有力的摘要，让读者在首页和详情页一眼知道这篇内容的价值。" }}
        </p>
      </div>

      <div class="preview-stats">
        <article class="preview-stat">
          <span>分类</span>
          <strong>{{ currentCategoryName }}</strong>
        </article>
        <article class="preview-stat">
          <span>标签数</span>
          <strong>{{ parsedTags.length }}</strong>
        </article>
        <article class="preview-stat">
          <span>字数</span>
          <strong>{{ wordCount }}</strong>
        </article>
      </div>

      <div v-if="coverUrl" class="preview-cover">
        <img :src="coverUrl" alt="封面预览" />
      </div>

      <div v-if="parsedTags.length" class="tag-row">
        <span v-for="tag in parsedTags.slice(0, 8)" :key="tag" class="tag-chip"># {{ tag }}</span>
      </div>

      <div class="preview-note">
        <span class="preview-note__dot"></span>
        <p>{{ previewGuide }}</p>
      </div>

      <div class="markdown-body article-detail-body" v-html="compiledMarkdown"></div>
    </aside>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from "vue";
import { onBeforeRouteLeave, useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { createArticle, getArticle, submitArticle, updateArticle } from "../api/article";
import { getMetadata } from "../api/meta";
import { uploadImage } from "../api/upload";
import { useUserStore } from "../stores/user";
import { toAssetUrl } from "../utils/asset";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const isEdit = computed(() => Boolean(route.params.id));
const saving = ref(false);
const submitting = ref(false);
const autosaving = ref(false);
const errorMessage = ref("");
const categories = ref([]);
const tagsInput = ref("");
const isBootstrapping = ref(true);
const isDirty = ref(false);
const lastSavedAt = ref("");
const lastServerSavedAt = ref("");
const hasLocalDraft = ref(false);
const showRestoreNotice = ref(false);
let autosaveTimer = null;

const form = reactive({
  title: "",
  summary: "",
  content: "",
  coverImage: "",
  categoryId: null
});

const parsedTags = computed(() =>
  tagsInput.value
    .split(/[,\n]/)
    .map((item) => item.trim())
    .filter((item, index, list) => item && list.indexOf(item) === index)
);

const compiledMarkdown = computed(() => DOMPurify.sanitize(marked.parse(form.content || "## 右侧会实时显示你的 Markdown 预览")));
const submitLabel = computed(() => (userStore.isAdmin ? "直接发布" : "提交审核"));
const submitModeLabel = computed(() => (userStore.isAdmin ? "管理员可直接发布" : "普通用户提交后进入审核"));
const coverUrl = computed(() => toAssetUrl(form.coverImage));
const currentCategoryName = computed(() => {
  const current = categories.value.find((item) => item.id === form.categoryId);
  return current?.name || "未分类";
});
const wordCount = computed(() => `${form.title} ${form.summary} ${form.content}`.replace(/\s+/g, "").length);
const previewGuide = computed(() => {
  if (!form.content.trim()) return "从结论、背景、方案、过程、总结五段开始，通常最容易写出结构清楚的文章。";
  if (parsedTags.value.length < 2) return "可以再补几个标签，方便后续归档、检索和前台展示。";
  if (!form.summary.trim()) return "摘要还空着，建议补上一段前言，首页展示会更完整。";
  return "结构已经成型了，继续把案例细节和你的思考写进去。";
});
const draftStateLabel = computed(() => {
  if (saving.value) return "正在保存到服务器";
  if (submitting.value) return "正在提交内容";
  if (autosaving.value) return "正在自动同步云端草稿";
  if (lastServerSavedAt.value) return `云端草稿已同步：${lastServerSavedAt.value}`;
  if (lastSavedAt.value) return `本地草稿已保存：${lastSavedAt.value}`;
  if (isDirty.value) return "检测到未同步修改";
  return "支持本地草稿与云端自动保存";
});

function draftStorageKey() {
  const userId = userStore.profile?.id || "guest";
  const articleId = route.params.id || "new";
  return `pulseblog:draft:${userId}:${articleId}`;
}

function buildPayload() {
  return {
    ...form,
    categoryId: form.categoryId || null,
    tags: parsedTags.value
  };
}

function hasMeaningfulContent() {
  return Boolean(form.title.trim() || form.summary.trim() || form.content.trim() || form.coverImage || tagsInput.value.trim() || form.categoryId);
}

function markClean() {
  isDirty.value = false;
}

function syncDraftState() {
  hasLocalDraft.value = Boolean(localStorage.getItem(draftStorageKey()));
}

function persistLocalDraft() {
  if (isBootstrapping.value) return;
  const payload = { ...buildPayload(), savedAt: new Date().toISOString() };
  localStorage.setItem(draftStorageKey(), JSON.stringify(payload));
  lastSavedAt.value = new Date(payload.savedAt).toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
  hasLocalDraft.value = true;
}

function restoreLocalDraft() {
  const raw = localStorage.getItem(draftStorageKey());
  if (!raw) {
    hasLocalDraft.value = false;
    return;
  }
  try {
    const payload = JSON.parse(raw);
    Object.assign(form, {
      title: payload.title || "",
      summary: payload.summary || "",
      content: payload.content || "",
      coverImage: payload.coverImage || "",
      categoryId: payload.categoryId || null
    });
    tagsInput.value = Array.isArray(payload.tags) ? payload.tags.join(", ") : "";
    if (payload.savedAt) {
      lastSavedAt.value = new Date(payload.savedAt).toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
    }
    hasLocalDraft.value = true;
    showRestoreNotice.value = true;
    markClean();
  } catch {
    localStorage.removeItem(draftStorageKey());
    hasLocalDraft.value = false;
  }
}

function clearLocalDraft() {
  localStorage.removeItem(draftStorageKey());
  hasLocalDraft.value = false;
  lastSavedAt.value = "";
}

function updateServerSavedTime() {
  lastServerSavedAt.value = new Date().toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
}

function dismissRestoreNotice() {
  showRestoreNotice.value = false;
}

function discardLocalDraft() {
  clearLocalDraft();
  showRestoreNotice.value = false;
  markClean();
}

function confirmLeave() {
  return window.confirm("当前还有未同步的修改，确定离开吗？");
}

function handleBeforeUnload(event) {
  if (!isDirty.value) return;
  event.preventDefault();
  event.returnValue = "";
}

async function loadDetail() {
  if (!isEdit.value) return;
  const { data } = await getArticle(route.params.id);
  Object.assign(form, {
    title: data.item.title || "",
    summary: data.item.summary || "",
    content: data.item.content || "",
    coverImage: data.item.coverImage || "",
    categoryId: data.item.categoryId || null
  });
  tagsInput.value = data.item.tags?.map((item) => item.name).join(", ") || "";
  markClean();
}

async function loadMetadata() {
  const { data } = await getMetadata();
  categories.value = data.categories || [];
}

async function saveArticle() {
  saving.value = true;
  errorMessage.value = "";
  try {
    if (isEdit.value) {
      await updateArticle(route.params.id, buildPayload());
      clearLocalDraft();
    } else {
      const oldKey = draftStorageKey();
      const { data } = await createArticle(buildPayload());
      localStorage.removeItem(oldKey);
      await router.replace(`/editor/${data.item.id}`);
      clearLocalDraft();
    }
    updateServerSavedTime();
    markClean();
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
    if (!isEdit.value) {
      const oldKey = draftStorageKey();
      const { data } = await createArticle(buildPayload());
      localStorage.removeItem(oldKey);
      await router.replace(`/editor/${data.item.id}`);
      await submitArticle(data.item.id);
    } else {
      await updateArticle(route.params.id, buildPayload());
      await submitArticle(route.params.id);
    }
    clearLocalDraft();
    lastServerSavedAt.value = "";
    markClean();
    router.push("/my-articles");
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || "提交失败，请稍后再试";
  } finally {
    submitting.value = false;
  }
}

async function autosaveToServer() {
  if (isBootstrapping.value || saving.value || submitting.value || autosaving.value) return;
  if (!hasMeaningfulContent()) return;
  autosaving.value = true;
  try {
    if (isEdit.value) {
      await updateArticle(route.params.id, buildPayload());
    } else {
      const oldKey = draftStorageKey();
      const { data } = await createArticle(buildPayload());
      localStorage.removeItem(oldKey);
      await router.replace(`/editor/${data.item.id}`);
    }
    updateServerSavedTime();
    markClean();
  } catch {
    // 网络异常时保留本地草稿作为兜底
  } finally {
    autosaving.value = false;
  }
}

async function handleUpload(event) {
  const [file] = event.target.files || [];
  if (!file) return;
  const { data } = await uploadImage(file);
  form.coverImage = data.url;
}

watch(
  () => [form.title, form.summary, form.content, form.coverImage, form.categoryId, tagsInput.value],
  () => {
    if (isBootstrapping.value) return;
    isDirty.value = true;
    if (autosaveTimer) clearTimeout(autosaveTimer);
    autosaveTimer = setTimeout(() => {
      persistLocalDraft();
      autosaveToServer();
    }, 800);
  }
);

watch(
  () => route.params.id,
  () => {
    syncDraftState();
  }
);

onBeforeRouteLeave(() => {
  if (!isDirty.value) return true;
  return confirmLeave();
});

onMounted(async () => {
  window.addEventListener("beforeunload", handleBeforeUnload);
  await loadMetadata();
  await loadDetail();
  restoreLocalDraft();
  isBootstrapping.value = false;
});

onBeforeUnmount(() => {
  window.removeEventListener("beforeunload", handleBeforeUnload);
  if (autosaveTimer) clearTimeout(autosaveTimer);
});
</script>

<style scoped>
.editor-workbench {
  margin-top: 28px;
  display: grid;
  grid-template-columns: minmax(0, 1.12fr) minmax(320px, 0.88fr);
  gap: 24px;
}

.editor-main,
.editor-preview {
  min-width: 0;
}

.editor-main {
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(circle at top left, rgba(255, 138, 76, 0.12), transparent 24%),
    radial-gradient(circle at 82% 18%, rgba(255, 209, 102, 0.1), transparent 22%),
    rgba(255, 255, 255, 0.07);
}

.editor-main::before {
  content: "";
  position: absolute;
  inset: auto -120px -120px auto;
  width: 280px;
  height: 280px;
  border-radius: 42% 58% 50% 50% / 46% 42% 58% 54%;
  background: radial-gradient(circle, rgba(255, 209, 102, 0.12), transparent 68%);
  pointer-events: none;
}

.editor-hero {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  margin-bottom: 22px;
}

.editor-hero__copy h2 {
  margin: 0 0 14px;
  font-size: clamp(2rem, 4vw, 3rem);
}

.editor-hero__status {
  width: min(280px, 100%);
  display: grid;
  gap: 12px;
}

.status-panel {
  padding: 16px 18px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(8, 11, 17, 0.38);
}

.status-panel span,
.hint-card span,
.preview-stat span {
  display: block;
  margin-bottom: 6px;
  color: var(--text-soft);
  font-size: 0.84rem;
}

.status-panel strong,
.hint-card strong,
.preview-stat strong {
  font-size: 1rem;
  line-height: 1.5;
}

.action-strip {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 22px;
}

.editor-body {
  position: relative;
  z-index: 1;
}

.draft-banner,
.hint-card {
  padding: 16px 18px;
  border-radius: 22px;
}

.draft-banner {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  margin-bottom: 20px;
  border: 1px solid rgba(255, 209, 102, 0.26);
  background: rgba(255, 209, 102, 0.08);
}

.draft-banner p {
  margin: 6px 0 0;
  color: var(--text-soft);
}

.hint-card {
  border: 1px dashed rgba(255, 209, 102, 0.26);
  background: rgba(255, 209, 102, 0.06);
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

.preview-head {
  margin-bottom: 20px;
}

.preview-head h3 {
  font-size: 1.8rem;
}

.preview-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.preview-stat {
  padding: 14px 16px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
}

.preview-note {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 16px 18px;
  margin-bottom: 20px;
  border-radius: 22px;
  border: 1px solid rgba(255, 209, 102, 0.2);
  background: rgba(255, 209, 102, 0.08);
}

.preview-note p {
  margin: 0;
  color: #fff0d5;
}

.preview-note__dot {
  width: 12px;
  height: 12px;
  margin-top: 6px;
  border-radius: 999px;
  background: linear-gradient(135deg, var(--accent-soft), var(--accent));
  box-shadow: 0 0 18px rgba(255, 138, 76, 0.42);
  flex-shrink: 0;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
  transform: none;
}

@media (max-width: 960px) {
  .editor-workbench {
    grid-template-columns: 1fr;
  }

  .editor-hero,
  .preview-stats {
    grid-template-columns: 1fr;
    display: grid;
  }

  .editor-hero__status {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .editor-workbench {
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

  .editor-hero__copy h2 {
    font-size: clamp(1.5rem, 5vw, 2rem);
  }

  .preview-head h3 {
    font-size: 1.4rem;
  }
}

@media (max-width: 480px) {
  .editor-main {
    padding: 18px;
  }

  .draft-banner {
    flex-direction: column;
    text-align: left;
  }

  .draft-banner button {
    width: 100%;
    justify-content: center;
  }

  .status-panel {
    padding: 12px 14px;
  }

  .preview-stat {
    padding: 10px 12px;
  }
}
</style>
