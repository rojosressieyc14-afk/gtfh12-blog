<template>
  <section class="manage-panel">
    <div class="review-head">
      <div>
        <p class="admin-label">项目审核</p>
        <h2>待审项目</h2>
        <p class="review-context">{{ projectReviewSliceLabel }}</p>
      </div>
      <button class="role-btn" :disabled="reviewQueuesRefreshing" @click="refreshProjectReviews">
        {{ reviewQueuesRefreshing ? "刷新中..." : "刷新队列" }}
      </button>
    </div>

    <div class="article-bulk-bar">
      <div class="review-toolbar-group review-toolbar-group--filters">
        <input
          v-model.trim="projectReviewKeyword"
          class="filter-select filter-input"
          placeholder="筛选待审项目标题 / 作者"
        />
        <select v-model="projectReviewSort" class="filter-select">
          <option value="updated_desc">最近更新优先</option>
          <option value="updated_asc">最早提交优先</option>
          <option value="title_asc">按标题排序</option>
        </select>
        <select v-model.number="projectReviewPageSize" class="filter-select">
          <option :value="5">每页 5 条</option>
          <option :value="10">每页 10 条</option>
        </select>
      </div>
      <div class="review-toolbar-group review-toolbar-group--actions">
        <div class="review-selection-row">
          <button class="role-btn" @click="toggleSelectVisibleProjectReviews">
            {{ allVisibleProjectReviewsSelected ? "取消全选" : "全选待审" }}
          </button>
          <span class="taxonomy-usage">已选 {{ selectedProjectReviewIds.length }} 项</span>
          <input
            v-model.trim="bulkProjectReviewReason"
            class="filter-select filter-input review-reason-inline"
            placeholder="填写共享驳回原因"
          />
        </div>
        <div class="review-action-row">
          <button class="btn-approve" :disabled="!selectedProjectReviewIds.length || bulkProjectReviewApproveSaving" @click="bulkApproveSelectedProjectReviews">
            {{ bulkProjectReviewApproveSaving ? "通过中..." : "批量通过" }}
          </button>
          <button class="btn-reject" :disabled="!selectedProjectReviewIds.length || bulkProjectReviewSaving" @click="bulkRejectSelectedProjectReviews">
            {{ bulkProjectReviewSaving ? "驳回中..." : "批量驳回" }}
          </button>
          <button class="role-btn" :disabled="(!selectedProjectReviewIds.length && !bulkProjectReviewReason) || bulkProjectReviewSaving || bulkProjectReviewApproveSaving" @click="resetProjectBulkReviewState">
            清空
          </button>
        </div>
      </div>
      <div class="review-toolbar-group review-toolbar-group--feedback">
        <div class="review-shortcuts">
          <button
            v-for="reason in reviewReasonPresets"
            :key="`bulk-project-review-reason-${reason}`"
            class="role-btn"
            type="button"
            @click="bulkProjectReviewReason = reason"
          >
            {{ reason }}
          </button>
        </div>
        <div v-if="projectReviewLastAction" class="review-summary-inline">
          <strong>最近操作</strong>
          <span>{{ projectReviewLastAction }}</span>
        </div>
        <p class="bulk-review-hint">批量通过仅处理待审项目；批量驳回需要共享原因，执行后会清空当前选择。</p>
      </div>
    </div>

    <div class="table-list">
      <article v-for="item in pagedPendingProjects" :key="`project-review-${item.id}`" class="table-card table-card--article">
        <div class="table-main">
          <div class="table-title">
            <label class="article-select-check">
              <input :checked="selectedProjectReviewIds.includes(item.id)" type="checkbox" @change="toggleProjectReviewSelection(item.id)" />
            </label>
            <h3>{{ item.title }}</h3>
            <span class="pill pill--pending">待审核</span>
          </div>
          <div class="review-row-meta">
            <span class="review-row-meta__item">
              <strong>作者</strong>
              <span>{{ item.author?.username || "未知作者" }}</span>
            </span>
            <span class="review-row-meta__item">
              <strong>更新</strong>
              <span>{{ formatDate(item.updatedAt) }}</span>
            </span>
          </div>
          <p :class="['table-note', 'pending-review-summary', { 'pending-review-summary--empty': !item.summary }]">
            {{ item.summary || "暂无摘要" }}
          </p>
        </div>
        <div class="table-actions pending-review-actions">
          <span class="pending-review-actions__label">操作</span>
          <button class="role-btn pending-review-actions__secondary" @click="$emit('open-project-preview', item.id, true)">进入审核</button>
          <button class="role-btn pending-review-actions__secondary" @click="$emit('open-ai-review', 'project', item.id, item.title)">AI 审核</button>
          <button class="btn-approve pending-review-actions__primary" @click="quickApproveProject(item.id)">快速通过</button>
        </div>
      </article>
    </div>

    <div class="pager-row">
      <button class="role-btn" :disabled="projectReviewPage <= 1" @click="projectReviewPage -= 1">上一页</button>
      <span>第 {{ projectReviewPage }} / {{ projectReviewTotalPages }} 页 · {{ filteredPendingProjects.length }} 项待处理</span>
      <button class="role-btn" :disabled="projectReviewPage >= projectReviewTotalPages" @click="projectReviewPage += 1">下一页</button>
    </div>

    <div v-if="!filteredPendingProjects.length" class="compact-empty compact-empty--actionable">
      <template v-if="pendingProjects.length">
        <p>当前筛选条件下没有待审项目。</p>
        <div class="compact-empty__actions">
          <button class="role-btn" type="button" @click="resetProjectReviewFilters">清除筛选</button>
          <button class="role-btn" type="button" :disabled="reviewQueuesRefreshing" @click="refreshProjectReviews">刷新队列</button>
        </div>
      </template>
      <template v-else>
        <p>当前没有待审项目，刷新后会重新拉取队列。</p>
        <div class="compact-empty__actions">
          <button class="role-btn" type="button" :disabled="reviewQueuesRefreshing" @click="refreshProjectReviews">刷新队列</button>
        </div>
      </template>
    </div>
  </section>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import {
  getPendingProjectReviews,
  reviewProject,
  bulkReviewAdminProjects
} from "../../api/dashboard";

const emit = defineEmits(["flash", "open-project-preview", "open-ai-review", "data-changed"]);

const pendingProjects = ref([]);
const projectReviewKeyword = ref("");
const projectReviewSort = ref("updated_desc");
const projectReviewPage = ref(1);
const projectReviewPageSize = ref(5);
const selectedProjectReviewIds = ref([]);
const bulkProjectReviewReason = ref("");
const bulkProjectReviewSaving = ref(false);
const bulkProjectReviewApproveSaving = ref(false);
const projectReviewLastAction = ref("");
const reviewQueuesRefreshing = ref(false);

const reviewReasonPresets = [
  "请补充事实依据与示例。",
  "请完善摘要、截图或成果说明。",
  "当前内容结构不完整，请补齐后再提交。",
  "存在明显排版或错别字问题，请修订后再提交。"
];

function say(message) {
  emit("flash", message);
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

const filteredPendingProjects = computed(() => {
  const keyword = projectReviewKeyword.value.trim().toLowerCase();
  const items = pendingProjects.value.filter((item) => {
    if (!keyword) return true;
    const haystack = [item.title, item.summary, item.author?.username]
      .filter(Boolean)
      .join(" ")
      .toLowerCase();
    return haystack.includes(keyword);
  });
  return [...items].sort((left, right) => compareQueueItems(left, right, projectReviewSort.value));
});

const projectReviewTotalPages = computed(() =>
  Math.max(1, Math.ceil(filteredPendingProjects.value.length / projectReviewPageSize.value))
);

const pagedPendingProjects = computed(() => {
  const start = (projectReviewPage.value - 1) * projectReviewPageSize.value;
  return filteredPendingProjects.value.slice(start, start + projectReviewPageSize.value);
});

const projectReviewSliceLabel = computed(() => {
  if (!filteredPendingProjects.value.length) return "当前无待审项目";
  const start = (projectReviewPage.value - 1) * projectReviewPageSize.value + 1;
  const end = Math.min(start + pagedPendingProjects.value.length - 1, filteredPendingProjects.value.length);
  return `当前显示 ${start}-${end} / ${filteredPendingProjects.value.length} 项`;
});

const allVisibleProjectReviewsSelected = computed(
  () => pagedPendingProjects.value.length > 0 &&
    pagedPendingProjects.value.every((item) => selectedProjectReviewIds.value.includes(item.id))
);

function compareQueueItems(left, right, sortMode) {
  if (sortMode === "updated_asc") {
    return new Date(left.updatedAt).getTime() - new Date(right.updatedAt).getTime();
  }
  if (sortMode === "title_asc") {
    return String(left.title || "").localeCompare(String(right.title || ""), "zh-CN");
  }
  return new Date(right.updatedAt).getTime() - new Date(left.updatedAt).getTime();
}

function resetProjectBulkReviewState() {
  selectedProjectReviewIds.value = [];
  bulkProjectReviewReason.value = "";
  bulkProjectReviewSaving.value = false;
  bulkProjectReviewApproveSaving.value = false;
}

function resetProjectReviewFilters() {
  projectReviewKeyword.value = "";
  projectReviewSort.value = "updated_desc";
  projectReviewPage.value = 1;
  projectReviewPageSize.value = 5;
}

function clampQueuePages() {
  projectReviewPage.value = Math.min(projectReviewPage.value, projectReviewTotalPages.value);
}

function toggleProjectReviewSelection(projectID) {
  if (selectedProjectReviewIds.value.includes(projectID)) {
    selectedProjectReviewIds.value = selectedProjectReviewIds.value.filter((item) => item !== projectID);
    return;
  }
  selectedProjectReviewIds.value = [...selectedProjectReviewIds.value, projectID];
}

function toggleSelectVisibleProjectReviews() {
  if (allVisibleProjectReviewsSelected.value) {
    selectedProjectReviewIds.value = [];
    return;
  }
  selectedProjectReviewIds.value = pagedPendingProjects.value.map((item) => item.id);
}

async function refreshProjectReviews() {
  reviewQueuesRefreshing.value = true;
  try {
    const projectRes = await getPendingProjectReviews();
    pendingProjects.value = projectRes.data.items || [];
    resetProjectBulkReviewState();
  } finally {
    reviewQueuesRefreshing.value = false;
  }
}

async function quickApproveProject(id) {
  try {
    await reviewProject(id, { action: "approve", reason: "" });
    await refreshProjectReviews();
    projectReviewLastAction.value = `已快速通过项目 #${id}。`;
    emit("data-changed");
    say("项目已审核通过。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "项目审核失败。");
  }
}

async function bulkApproveSelectedProjectReviews() {
  if (!selectedProjectReviewIds.value.length || bulkProjectReviewApproveSaving.value) return;
  bulkProjectReviewApproveSaving.value = true;
  try {
    const count = selectedProjectReviewIds.value.length;
    await bulkReviewAdminProjects({
      projectIds: selectedProjectReviewIds.value,
      action: "approve",
      reason: ""
    });
    await refreshProjectReviews();
    projectReviewLastAction.value = `已批量通过 ${count} 个项目。`;
    emit("data-changed");
    say("批量通过项目已完成。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "批量通过项目失败。");
  } finally {
    resetProjectBulkReviewState();
  }
}

async function bulkRejectSelectedProjectReviews() {
  if (!selectedProjectReviewIds.value.length || bulkProjectReviewSaving.value) return;
  if (!bulkProjectReviewReason.value.trim()) {
    say("请先填写项目批量驳回原因。");
    return;
  }
  bulkProjectReviewSaving.value = true;
  try {
    const count = selectedProjectReviewIds.value.length;
    await bulkReviewAdminProjects({
      projectIds: selectedProjectReviewIds.value,
      action: "reject",
      reason: bulkProjectReviewReason.value.trim()
    });
    await refreshProjectReviews();
    projectReviewLastAction.value = `已批量驳回 ${count} 个项目。`;
    emit("data-changed");
    say("批量驳回项目已完成。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "批量驳回项目失败。");
  } finally {
    resetProjectBulkReviewState();
  }
}

watch([projectReviewKeyword, projectReviewSort], () => {
  projectReviewPage.value = 1;
  resetProjectBulkReviewState();
  projectReviewLastAction.value = "";
});

watch([projectReviewPage, projectReviewPageSize], () => {
  clampQueuePages();
  resetProjectBulkReviewState();
  projectReviewLastAction.value = "";
});

async function loadData() {
  try {
    const projectRes = await getPendingProjectReviews();
    pendingProjects.value = projectRes.data.items || [];
  } catch {
    // silently fail on initial load
  }
}

loadData();
</script>

<style scoped>
.review-head--stack {
  align-items: flex-start;
}

.table-title {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.review-row-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 16px;
  color: var(--soft);
  font-size: 0.9rem;
}

.review-row-meta__item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.review-row-meta__item strong {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 600;
}

.pending-review-summary {
  margin: 6px 0 0;
  color: rgba(255, 255, 255, 0.78);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.45;
}

.pending-review-summary--empty {
  color: var(--soft);
  font-style: italic;
}

.pending-review-actions {
  min-width: 116px;
  display: grid;
  align-content: start;
  justify-items: stretch;
  gap: 8px;
}

.pending-review-actions__label {
  color: var(--soft);
  font-size: 0.82rem;
  text-align: right;
}

.pending-review-actions__secondary,
.pending-review-actions__primary {
  width: 100%;
}

.pill--pending {
  background: rgba(255, 166, 77, 0.2);
  color: #ffd79d;
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

.taxonomy-usage {
  color: var(--soft);
  font-size: 0.9rem;
  white-space: nowrap;
}

@media (max-width: 960px) {
  .review-toolbar-group--actions,
  .review-toolbar-group--feedback {
    padding: 10px;
  }
  .pending-review-actions {
    min-width: 0;
  }
  .pending-review-actions__label {
    text-align: left;
  }
}

@media (max-width: 768px) {
  .review-toolbar-group--filters {
    flex-direction: column;
  }
  .review-toolbar-group--filters .filter-select {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .review-head {
    flex-direction: column;
  }
  .review-selection-row,
  .review-action-row {
    flex-direction: column;
    align-items: stretch;
  }
  .review-reason-inline {
    flex: 1 1 auto;
  }
  .review-shortcuts {
    flex-direction: column;
  }
  .review-shortcuts button {
    width: 100%;
  }
  .table-card {
    flex-direction: column;
  }
  .table-main {
    width: 100%;
  }
  .pending-review-actions {
    width: 100%;
  }
  .pending-review-actions button {
    width: 100%;
  }
  .review-row-meta {
    flex-direction: column;
    gap: 4px;
  }
  .table-title {
    flex-direction: column;
    align-items: flex-start;
  }
  .article-select-check {
    align-self: flex-start;
  }
}
</style>
