<template>
  <div class="admin-shell">
    <button class="sidebar-toggle" :class="{ 'sidebar-toggle--active': sidebarOpen }" @click="sidebarOpen = !sidebarOpen" aria-label="菜单">
      <span></span><span></span><span></span>
    </button>

    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>
    <aside class="admin-sidebar admin-sidebar--brand" :class="{ 'admin-sidebar--open': sidebarOpen }">
      <div class="admin-sidebar__top">
        <p class="admin-label">运营后台</p>
        <h2>PulseBlog</h2>
        <p class="sidebar-text">集中处理审核、风控、用户与内容管理。</p>
      </div>

      <nav class="admin-nav dashboard-nav" @click="sidebarOpen = false">
        <a href="#overview">总览</a>
        <a href="#taxonomy">分类与标签</a>
        <a href="#review">审核</a>
        <a href="#users">用户</a>
        <a href="#articles">文章</a>
        <a href="#projects">项目</a>
        <router-link to="/comments">评论</router-link>
        <router-link to="/logs">日志</router-link>
        <router-link to="/uploads">资源</router-link>
        <router-link to="/settings">系统设置</router-link>
      </nav>

      <div class="admin-user">
        <span>{{ store.profile?.username || "管理员" }}</span>
        <button @click="logout">退出登录</button>
      </div>
    </aside>

    <main class="admin-main dashboard-page">
      <div v-if="flash" class="flash-banner">{{ flash }}</div>

      <DashboardOverview :stats="stats" @refresh="loadData" />
      <CategoryTagSection :categories="categories" :tags="tags" @flash="say" @data-changed="handleTaxonomyChanged" />
      <section id="review" class="manage-grid dashboard-grid">
        <ArticleReviewPanel @flash="say" @open-article-preview="openArticlePreview" @open-ai-review="openAIReview" @data-changed="handleArticleDataChanged" />
        <ProjectReviewPanel @flash="say" @open-project-preview="openProjectPreview" @open-ai-review="openAIReview" @data-changed="handleProjectDataChanged" />
      </section>
      <UserManagementSection @flash="say" />
      <ArticleManagementSection :categories="categories" :tags="tags" @flash="say" @open-article-preview="openArticlePreview" @data-changed="handleArticleDataChanged" />
      <ProjectManagementSection @flash="say" @open-project-preview="openProjectPreview" @data-changed="handleProjectDataChanged" />
    </main>

    <div v-if="previewVisible && previewArticle" class="preview-mask" @click.self="closePreview">
      <article class="preview-dialog">
        <header class="preview-dialog__head">
          <div>
            <p class="admin-label">{{ reviewMode ? "文章审核" : "文章预览" }}</p>
            <h2>{{ previewArticle.title }}</h2>
          </div>
          <button class="role-btn" @click="closePreview">关闭</button>
        </header>
        <p class="table-note">{{ previewArticle.author?.username || "未知作者" }} · {{ statusLabel(previewArticle.status) }}</p>
        <p v-if="previewArticle.summary" class="preview-summary">{{ previewArticle.summary }}</p>
        <div class="taxonomy-editor">
          <label class="taxonomy-editor__field">
            <span>分类</span>
            <select v-model="previewCategoryId" class="filter-select">
              <option value="">未分类</option>
              <option v-for="item in categories" :key="`preview-category-${item.id}`" :value="String(item.id)">
                {{ item.name }}
              </option>
            </select>
          </label>
          <div class="taxonomy-editor__field taxonomy-editor__field--stack">
            <span>标签</span>
            <div class="taxonomy-chip-list">
              <label v-for="item in tags" :key="`preview-tag-${item.id}`" class="taxonomy-chip">
                <input
                  :checked="previewTagIds.includes(item.id)"
                  type="checkbox"
                  @change="togglePreviewTag(item.id)"
                />
                <span>{{ item.name }}</span>
              </label>
            </div>
          </div>
          <div class="taxonomy-editor__actions">
            <button class="role-btn" :disabled="articleTaxonomySaving" @click="savePreviewTaxonomy">
              {{ articleTaxonomySaving ? "保存中..." : "保存分类与标签" }}
            </button>
          </div>
        </div>
        <pre class="preview-content">{{ previewArticle.content }}</pre>

        <div v-if="reviewMode" class="review-box">
          <div class="review-shortcuts">
            <button
              v-for="reason in reviewReasonPresets"
              :key="`article-review-reason-${reason}`"
              class="role-btn"
              type="button"
              @click="reviewReason = reason"
            >
              {{ reason }}
            </button>
            <button class="role-btn" type="button" @click="reviewReason = ''">清空原因</button>
          </div>
          <textarea
            v-model.trim="reviewReason"
            class="review-reason"
            placeholder="驳回时请填写原因，通过时可以留空？"
          ></textarea>
          <div class="review-actions">
            <button class="btn-approve" @click="submitReview('approve')">审核通过并发布</button>
            <button class="btn-reject" @click="submitReview('reject')">驳回文章</button>
          </div>
        </div>
      </article>
    </div>

    <div v-if="previewVisible && previewProject" class="preview-mask" @click.self="closePreview">
      <article class="preview-dialog">
        <header class="preview-dialog__head">
          <div>
            <p class="admin-label">{{ projectReviewMode ? "项目审核" : "项目预览" }}</p>
            <h2>{{ previewProject.title }}</h2>
          </div>
          <button class="role-btn" @click="closePreview">关闭</button>
        </header>
        <p class="table-note">{{ previewProject.author?.username || "未知作者" }} · {{ projectStatusLabel(previewProject.status) }}</p>
        <p v-if="previewProject.summary" class="preview-summary">{{ previewProject.summary }}</p>
        <pre class="preview-content">{{ previewProject.content }}</pre>

        <div v-if="projectReviewMode" class="review-box">
          <div class="review-shortcuts">
            <button
              v-for="reason in reviewReasonPresets"
              :key="`project-review-reason-${reason}`"
              class="role-btn"
              type="button"
              @click="projectReviewReason = reason"
            >
              {{ reason }}
            </button>
            <button class="role-btn" type="button" @click="projectReviewReason = ''">清空原因</button>
          </div>
          <textarea
            v-model.trim="projectReviewReason"
            class="review-reason"
            placeholder="驳回时请填写原因，通过时可以留空？"
          ></textarea>
          <div class="review-actions">
            <button class="btn-approve" @click="submitProjectReview(previewProject.id, 'approve')">审核通过并发布</button>
            <button class="btn-reject" @click="submitProjectReview(previewProject.id, 'reject')">驳回项目</button>
          </div>
        </div>
      </article>
    </div>

    <AIReviewPanel
      v-if="aiReviewVisible"
      :target-type="aiReviewTargetType"
      :target-id="aiReviewTargetId"
      :target-title="aiReviewTargetTitle"
      @close="closeAIReview"
      @saved="onAIReviewSaved"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  getDashboard,
  getAdminArticleDetail,
  getAdminProjectDetail,
  reviewArticle,
  reviewProject,
  updateAdminArticleTaxonomy
} from "../api/dashboard";
import { getMetadata } from "../api/meta";
import { useAdminStore } from "../stores/auth";

import AIReviewPanel from "../components/AIReviewPanel.vue";
import DashboardOverview from "../components/dashboard/DashboardOverview.vue";
import CategoryTagSection from "../components/dashboard/CategoryTagSection.vue";
import ArticleReviewPanel from "../components/dashboard/ArticleReviewPanel.vue";
import ProjectReviewPanel from "../components/dashboard/ProjectReviewPanel.vue";
import UserManagementSection from "../components/dashboard/UserManagementSection.vue";
import ArticleManagementSection from "../components/dashboard/ArticleManagementSection.vue";
import ProjectManagementSection from "../components/dashboard/ProjectManagementSection.vue";

const router = useRouter();
const store = useAdminStore();
const sidebarOpen = ref(false);

// ── Parent-level shared state ──────────────────────────────────────────
const stats = ref({});
const categories = ref([]);
const tags = ref([]);
const flash = ref("");

// ── AI review state ────────────────────────────────────────────────────
const aiReviewVisible = ref(false);
const aiReviewTargetType = ref("article");
const aiReviewTargetId = ref(0);
const aiReviewTargetTitle = ref("");

function openAIReview(type, id, title) {
  aiReviewTargetType.value = type;
  aiReviewTargetId.value = id;
  aiReviewTargetTitle.value = title;
  aiReviewVisible.value = true;
}

function closeAIReview() {
  aiReviewVisible.value = false;
}

function onAIReviewSaved() {
  closeAIReview();
}

// ── Preview dialog state ───────────────────────────────────────────────
const previewVisible = ref(false);
const previewArticle = ref(null);
const previewProject = ref(null);
const reviewMode = ref(false);
const projectReviewMode = ref(false);
const reviewReason = ref("");
const projectReviewReason = ref("");
const previewCategoryId = ref("");
const previewTagIds = ref([]);
const articleTaxonomySaving = ref(false);

const reviewReasonPresets = [
  "请补充事实依据与示例。",
  "请完善摘要、截图或成果说明。",
  "当前内容结构不完整，请补齐后再提交。",
  "存在明显排版或错别字问题，请修订后再提交。"
];

// ── Shared helpers ─────────────────────────────────────────────────────
function say(message) {
  flash.value = message;
  window.setTimeout(() => {
    if (flash.value === message) {
      flash.value = "";
    }
  }, 2200);
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

function statusLabel(status) {
  return {
    draft: "草稿",
    pending: "审核中",
    published: "已发布",
    rejected: "已驳回"
  }[status] || status;
}

function projectStatusLabel(status) {
  return statusLabel(status);
}

function statusClass(status) {
  return {
    "pill--pending": status === "pending",
    "pill--success": status === "published",
    "pill--reject": status === "rejected"
  };
}

function normalizeSortOrder(value) {
  const parsed = Number.parseInt(value, 10);
  return Number.isFinite(parsed) && parsed >= 0 ? parsed : 0;
}

// ── Metadata loading ───────────────────────────────────────────────────
async function loadMetadataOnly() {
  const { data } = await getMetadata();
  categories.value = data.categories || [];
  tags.value = data.tags || [];
}

// ── Data-changed coordination handlers ─────────────────────────────────
async function handleTaxonomyChanged() {
  await loadMetadataOnly();
}

async function handleArticleDataChanged() {
  try {
    const { data } = await getDashboard();
    stats.value = data.stats || {};
  } catch {
    // stats refresh is best-effort
  }
}

async function handleProjectDataChanged() {
  try {
    const { data } = await getDashboard();
    stats.value = data.stats || {};
  } catch {
    // stats refresh is best-effort
  }
}

// ── Main data loading ──────────────────────────────────────────────────
async function loadData() {
  try {
    const profile = await store.fetchProfile();
    if (!profile) {
      router.push("/login");
      return;
    }
    const [dashboardRes] = await Promise.all([
      getDashboard(),
      loadMetadataOnly()
    ]);
    stats.value = dashboardRes.data.stats || {};
    say("后台数据已刷新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "后台数据加载失败。");
  }
}

// ── Preview dialog helpers ─────────────────────────────────────────────
function setPreviewTaxonomy(item) {
  previewCategoryId.value = item?.categoryId ? String(item.categoryId) : "";
  previewTagIds.value = (item?.tags || []).map((tag) => tag.id).filter(Boolean);
}

function togglePreviewTag(tagID) {
  if (previewTagIds.value.includes(tagID)) {
    previewTagIds.value = previewTagIds.value.filter((item) => item !== tagID);
    return;
  }
  previewTagIds.value = [...previewTagIds.value, tagID];
}

async function savePreviewTaxonomy() {
  if (!previewArticle.value || articleTaxonomySaving.value) return;
  articleTaxonomySaving.value = true;
  try {
    const payload = {
      categoryId: previewCategoryId.value ? Number(previewCategoryId.value) : null,
      tagIds: previewTagIds.value
    };
    const { data } = await updateAdminArticleTaxonomy(previewArticle.value.id, payload);
    previewArticle.value = data.item;
    setPreviewTaxonomy(data.item);
    say("文章分类与标签已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新文章分类与标签失败。");
  } finally {
    articleTaxonomySaving.value = false;
  }
}

async function openArticlePreview(id, asReview = false) {
  try {
    const { data } = await getAdminArticleDetail(id);
    previewArticle.value = data.item;
    setPreviewTaxonomy(data.item);
    previewProject.value = null;
    previewVisible.value = true;
    reviewMode.value = asReview;
    projectReviewMode.value = false;
    reviewReason.value = data.item?.rejectReason || "";
    projectReviewReason.value = "";
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "加载文章详情失败。");
  }
}

async function openProjectPreview(id, asReview = false) {
  try {
    const { data } = await getAdminProjectDetail(id);
    previewProject.value = data.item;
    previewArticle.value = null;
    previewVisible.value = true;
    projectReviewMode.value = asReview;
    reviewMode.value = false;
    projectReviewReason.value = data.item?.rejectReason || "";
    reviewReason.value = "";
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "加载项目详情失败。");
  }
}

function closePreview() {
  previewVisible.value = false;
  previewArticle.value = null;
  previewProject.value = null;
  reviewMode.value = false;
  projectReviewMode.value = false;
  reviewReason.value = "";
  projectReviewReason.value = "";
  previewCategoryId.value = "";
  previewTagIds.value = [];
  articleTaxonomySaving.value = false;
}

async function submitReview(action) {
  if (!previewArticle.value) return;
  if (action === "reject" && !reviewReason.value) {
    say("驳回时请填写原因。");
    return;
  }
  try {
    const articleID = previewArticle.value.id;
    await reviewArticle(articleID, { action, reason: reviewReason.value });
    closePreview();
    say(action === "approve" ? "文章已审核通过。" : "文章已驳回。");
    handleArticleDataChanged();
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "提交审核结果失败。");
  }
}

async function submitProjectReview(id, action) {
  if (action === "reject" && !projectReviewReason.value) {
    say("驳回项目时请填写原因。");
    return;
  }
  try {
    await reviewProject(id, { action, reason: projectReviewReason.value });
    closePreview();
    say(action === "approve" ? "项目已审核通过。" : "项目已驳回。");
    handleProjectDataChanged();
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "提交项目审核结果失败。");
  }
}

// ── Auth ───────────────────────────────────────────────────────────────
function logout() {
  store.logout();
  router.push("/login");
}

onMounted(loadData);
</script>

<style scoped>
.dashboard-page {
  display: grid;
  gap: 20px;
}

.dashboard-nav {
  display: grid;
  gap: 10px;
}

.dashboard-nav a {
  padding: 10px 12px;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.03);
  color: inherit;
  text-decoration: none;
}

.dashboard-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.table-title {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.pill--pending {
  background: rgba(255, 166, 77, 0.2);
  color: #ffd79d;
}

.pill--success {
  background: rgba(131, 242, 143, 0.16);
  color: #c8f8cd;
}

.pill--reject {
  background: rgba(255, 139, 139, 0.16);
  color: #ffd2d2;
}

.reject-note {
  color: #ffd2d2;
}

.taxonomy-editor {
  margin: 16px 0;
  display: grid;
  gap: 12px;
}

.taxonomy-editor__field {
  display: grid;
  gap: 8px;
}

.taxonomy-editor__field--stack {
  align-items: flex-start;
}

.taxonomy-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.taxonomy-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
}

.taxonomy-editor__actions {
  display: flex;
  justify-content: flex-start;
}

.review-head--stack {
  align-items: flex-start;
}

.review-shortcuts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.review-context {
  margin: 6px 0 0;
  color: var(--soft);
  font-size: 0.92rem;
}

.review-summary-inline {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  color: var(--soft);
  font-size: 0.92rem;
}

.compact-empty--actionable {
  display: grid;
  gap: 10px;
}

.compact-empty--actionable p {
  margin: 0;
}

.compact-empty__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.article-select-check {
  display: inline-flex;
  align-items: center;
}

.article-bulk-bar {
  display: grid;
  gap: 10px;
}

.review-toolbar-group {
  display: grid;
  gap: 10px;
}

.review-toolbar-group--actions,
.review-toolbar-group--feedback {
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
}

.review-selection-row,
.review-action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.review-reason-inline {
  flex: 1 1 240px;
}

.bulk-review-hint {
  margin: 0;
  color: var(--soft);
  font-size: 0.9rem;
}

.taxonomy-usage {
  color: var(--soft);
  font-size: 0.9rem;
  white-space: nowrap;
}

/* Sidebar toggle button (mobile only) */
.sidebar-toggle {
  display: none;
  position: fixed;
  top: 12px;
  left: 12px;
  z-index: 110;
  flex-direction: column;
  gap: 5px;
  padding: 10px;
  background: rgba(10, 14, 19, 0.85);
  border: 1px solid var(--border);
  border-radius: 14px;
  cursor: pointer;
  backdrop-filter: blur(8px);
}

.sidebar-toggle span {
  display: block;
  width: 22px;
  height: 2px;
  background: var(--text, #f7f3ea);
  border-radius: 2px;
  transition: transform 0.3s, opacity 0.3s;
}

.sidebar-toggle--active span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}
.sidebar-toggle--active span:nth-child(2) {
  opacity: 0;
}
.sidebar-toggle--active span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

.sidebar-overlay {
  display: none;
}

@media (max-width: 960px) {
  .admin-shell {
    position: relative;
  }
  .sidebar-toggle {
    display: flex;
  }
  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.5);
    z-index: 100;
  }
  .admin-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: min(280px, 75vw);
    height: 100vh;
    z-index: 101;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    overflow-y: auto;
  }
  .admin-sidebar--open {
    transform: translateX(0);
  }
  .admin-main {
    margin-left: 0;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  .review-toolbar-group--actions,
  .review-toolbar-group--feedback {
    padding: 10px;
  }
}

@media (max-width: 768px) {
  .admin-main {
    padding: 16px;
  }
  .preview-dialog {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .admin-main {
    padding: 12px;
    padding-top: 56px;
  }
  .preview-dialog {
    padding: 14px;
    width: 100vw;
    height: 100vh;
    border-radius: 0;
  }
  .preview-dialog__head {
    flex-direction: column;
    gap: 12px;
  }
  .review-actions {
    flex-direction: column;
  }
  .review-actions button {
    width: 100%;
  }
  .review-shortcuts {
    flex-direction: column;
  }
  .review-shortcuts button {
    width: 100%;
  }
  .review-reason {
    width: 100%;
  }
}
</style>
