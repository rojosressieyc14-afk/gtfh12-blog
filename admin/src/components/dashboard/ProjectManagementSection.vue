<template>
  <section id="projects" class="manage-panel">
    <div class="review-head review-head--stack">
      <div>
        <p class="admin-label">项目管理</p>
        <h2>项目列表</h2>
      </div>
      <div class="toolbar-row">
        <input
          v-model.trim="projectKeyword"
          class="filter-select filter-input"
          placeholder="搜索项目标题 / 摘要"
          @keyup.enter="changeProjectPage(1)"
        />
        <select v-model="projectStatusFilter" class="filter-select" @change="changeProjectPage(1)">
          <option value="">全部状态</option>
          <option value="draft">草稿</option>
          <option value="pending">审核中</option>
          <option value="published">已发布</option>
          <option value="rejected">已驳回</option>
        </select>
        <button class="role-btn" @click="changeProjectPage(1)">搜索</button>
      </div>
    </div>

    <div class="table-list">
      <article v-for="item in projects" :key="item.id" class="table-card table-card--article">
        <div class="table-main">
          <div class="table-title">
            <h3>{{ item.title }}</h3>
            <span class="pill" :class="statusClass(item.status)">{{ projectStatusLabel(item.status) }}</span>
          </div>
          <p>{{ item.author?.username || "未知作者" }} · {{ formatDate(item.updatedAt) }}</p>
          <p class="table-note">{{ item.summary || "暂无摘要" }}</p>
          <p v-if="item.rejectReason" class="reject-note">驳回原因：{{ item.rejectReason }}</p>
          <div class="project-meta-tools">
            <label class="project-meta-check">
              <input v-model="item.isFeatured" :disabled="item.status !== 'published'" type="checkbox" />
              <span>首页精选</span>
            </label>
            <label class="project-meta-field">
              <span>排序优先级</span>
              <input
                v-model.number="item.sortOrder"
                class="filter-select project-sort-input"
                :disabled="item.status !== 'published'"
                type="number"
                min="0"
              />
            </label>
          </div>
        </div>
        <div class="table-actions">
          <button
            class="role-btn"
            :disabled="item.status !== 'published' || Boolean(projectMetaSaving[item.id])"
            @click="saveProjectMeta(item)"
          >
            {{ projectMetaSaving[item.id] ? "保存中..." : "保存展示设置" }}
          </button>
          <button class="role-btn" @click="$emit('open-project-preview', item.id, item.status === 'pending')">
            {{ item.status === "pending" ? "审核" : "预览" }}
          </button>
          <button v-if="item.status !== 'published'" class="btn-approve" @click="publishProject(item.id)">发布</button>
          <button class="btn-reject" @click="deleteProject(item.id)">删除</button>
        </div>
      </article>
    </div>

    <div class="pager-row">
      <button class="role-btn" :disabled="projectPage <= 1" @click="changeProjectPage(projectPage - 1)">上一页</button>
      <span>第 {{ projectPage }} / {{ projectTotalPages }} 页</span>
      <button class="role-btn" :disabled="projectPage >= projectTotalPages" @click="changeProjectPage(projectPage + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import {
  getAdminProjects,
  publishAdminProject,
  deleteAdminProject,
  updateAdminProjectMeta,
  getPendingReviews,
  getPendingProjectReviews
} from "../../api/dashboard";

const emit = defineEmits(["flash", "open-project-preview", "data-changed"]);

const projects = ref([]);
const projectStatusFilter = ref("");
const projectKeyword = ref("");
const projectPage = ref(1);
const projectTotal = ref(0);
const pageSize = 6;
const projectMetaSaving = ref({});

function say(message) {
  emit("flash", message);
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

function projectStatusLabel(status) {
  return {
    draft: "草稿",
    pending: "审核中",
    published: "已发布",
    rejected: "已驳回"
  }[status] || status;
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

const projectTotalPages = computed(() => Math.max(1, Math.ceil(projectTotal.value / pageSize)));

async function loadProjects() {
  const { data } = await getAdminProjects({
    page: projectPage.value,
    pageSize,
    status: projectStatusFilter.value,
    keyword: projectKeyword.value
  });
  projects.value = data.items || [];
  projectTotal.value = data.pagination?.total || 0;
}

async function changeProjectPage(nextPage) {
  projectPage.value = nextPage;
  await loadProjects();
}

async function refreshReviews() {
  try {
    await Promise.all([getPendingReviews(), getPendingProjectReviews()]);
  } catch {
    // silently refresh caches
  }
}

async function publishProject(id) {
  try {
    await publishAdminProject(id);
    await loadProjects();
    await refreshReviews();
    emit("data-changed");
    say("项目已发布。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "发布项目失败。");
  }
}

async function deleteProject(id) {
  if (!window.confirm("确定删除这个项目吗？")) return;
  try {
    await deleteAdminProject(id);
    await loadProjects();
    await refreshReviews();
    emit("data-changed");
    say("项目已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除项目失败。");
  }
}

async function saveProjectMeta(item) {
  const payload = {
    isFeatured: Boolean(item.isFeatured),
    sortOrder: normalizeSortOrder(item.sortOrder)
  };
  projectMetaSaving.value = { ...projectMetaSaving.value, [item.id]: true };
  try {
    const { data } = await updateAdminProjectMeta(item.id, payload);
    item.isFeatured = Boolean(data.item?.isFeatured);
    item.sortOrder = normalizeSortOrder(data.item?.sortOrder);
    say("项目展示信息已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新项目展示信息失败。");
  } finally {
    projectMetaSaving.value = { ...projectMetaSaving.value, [item.id]: false };
  }
}

loadProjects();
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

.project-meta-tools {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.project-meta-check,
.project-meta-field {
  display: flex;
  align-items: center;
  gap: 10px;
}

.project-meta-check span,
.project-meta-field span {
  color: var(--soft);
}

.project-sort-input {
  width: 120px;
}

@media (max-width: 768px) {
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
}

@media (max-width: 480px) {
  .table-card {
    flex-direction: column;
    align-items: flex-start;
  }
  .table-actions {
    flex-direction: column;
    width: 100%;
  }
  .table-actions button {
    width: 100%;
  }
  .project-meta-tools {
    width: 100%;
  }
  .project-meta-check,
  .project-meta-field {
    flex-direction: column;
    align-items: flex-start;
  }
  .project-sort-input {
    width: 100%;
  }
  .table-title {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
