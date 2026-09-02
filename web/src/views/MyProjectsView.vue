<template>
  <section class="uc-projects">
    <section class="uc-hero uc-hero--projects">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">作品管理</p>
        <h2 class="uc-hero__title">我的项目</h2>
        <p class="uc-hero__text">统一查看草稿、审核中和已发布项目，整理出真正看得懂的案例。</p>
      </div>
      <div class="uc-hero__stats">
        <article class="uc-hero__stat">
          <strong>{{ total }}</strong>
          <span>全部项目</span>
        </article>
      </div>
    </section>

    <p v-if="errorMessage" class="error-text" style="margin-bottom:12px">{{ errorMessage }}</p>

    <div v-if="projects.length" class="project-grid">
      <article v-for="item in projects" :key="item.id" class="project-card uc-project-card">
        <div class="project-card__head">
          <div>
            <h3>{{ item.title }}</h3>
            <p class="table-note">{{ formatDate(item.updatedAt) }}</p>
          </div>
          <span class="status-chip" :class="item.status">{{ labelMap[item.status] || item.status }}</span>
          <span v-if="item.isPrivate" class="privacy-badge">私有</span>
        </div>

        <p class="detail-summary">{{ item.summary || "暂无摘要" }}</p>

        <div v-if="item.techStacks?.length" class="tag-row">
          <span v-for="stack in item.techStacks.slice(0, 6)" :key="stack" class="tag-chip"># {{ stack }}</span>
        </div>

        <p class="detail-summary project-status-copy">{{ statusCopy(item) }}</p>
        <p v-if="item.rejectReason" class="reject-tip">驳回原因：{{ item.rejectReason }}</p>

        <footer class="project-footer project-footer--mine">
          <span class="project-footer__label">{{ item.isFeatured ? "首页精选项目" : "普通展示项目" }}</span>
          <div class="project-footer__actions">
            <router-link v-if="item.status === 'published'" class="inline-link" :to="`/projects/${item.id}`">
              查看详情
            </router-link>
            <router-link class="inline-link" :to="`/project-editor/${item.id}`">继续编辑</router-link>
            <button class="inline-link delete-link" @click="handleDelete(item.id, item.title)">删除</button>
          </div>
        </footer>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>你还没有项目</h4>
      <p>创建你的第一个项目。</p>
    </div>

    <div class="pager-row">
      <button class="ghost-btn" :disabled="page <= 1" @click="goToPage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button class="ghost-btn" :disabled="page >= totalPages" @click="goToPage(page + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { deleteProject, listMyProjects } from "../api/project";
import { formatDate } from "../utils/date";

const projects = ref([]);
const page = ref(1);
const pageSize = 8;
const total = ref(0);
const errorMessage = ref("");

const labelMap = {
  draft: "草稿",
  pending: "审核中",
  published: "已发布",
  rejected: "已驳回"
};

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

async function loadMine() {
  const { data } = await listMyProjects(page.value, pageSize);
  projects.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

function statusCopy(item) {
  if (item.status === "draft") return "当前还是草稿，整理完整后再提交审核。";
  if (item.status === "pending") return "项目已提交审核，正在等待管理员处理。";
  if (item.status === "published") return "项目已经公开展示。再次编辑后会回到草稿状态。";
  if (item.status === "rejected") return "项目已被驳回，修改后可以重新提交审核。";
  return "";
}

async function handleDelete(id, title) {
  if (!confirm(`确定要删除项目「${title}」吗？此操作不可撤销。`)) return;
  try {
    await deleteProject(id);
    await loadMine();
  } catch (e) {
    errorMessage.value = "删除失败：" + (e.response?.data?.message || e.message);
  }
}

async function goToPage(nextPage) {
  page.value = nextPage;
  await loadMine();
}

onMounted(loadMine);
</script>

<style scoped>
.uc-projects {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.uc-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  padding: clamp(32px, 6vw, 56px) clamp(24px, 4vw, 48px);
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(200px, 0.6fr);
  gap: 32px;
  align-items: center;
}

.uc-hero__copy {
  display: grid;
  gap: 12px;
}

.uc-hero__eyebrow {
  color: #b4530a;
  font-weight: 600;
}

.uc-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
  color: #1d1d1f;
}

.uc-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: #48484a;
}

.uc-hero__stats {
  display: flex;
  gap: 12px;
}

.uc-hero__stat {
  padding: 22px 28px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.7);
  text-align: center;
  flex: 1;
}

.uc-hero__stat strong {
  display: block;
  font-size: 2rem;
  font-weight: 700;
  color: #1d1d1f;
  line-height: 1.2;
}

.uc-hero__stat span {
  display: block;
  margin-top: 4px;
  color: #6b7280;
  font-size: 0.85rem;
}

.project-grid {
  grid-template-columns: repeat(2, 1fr);
}

.uc-project-card {
  border-radius: 22px;
  transition: border-color 0.2s;
}

.project-grid :deep(.project-card:first-child),
.uc-project-card:first-child {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.uc-project-card:first-child::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.uc-project-card:first-child h3 {
  color: #1d1d1f;
}

.uc-project-card:first-child p,
.uc-project-card:first-child .detail-summary {
  color: #48484a;
}

.uc-project-card:first-child .table-note,
.uc-project-card:first-child .project-footer__label {
  color: #6b7280;
}

.uc-project-card:first-child .tag-chip {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.uc-project-card:first-child .status-chip {
  background: rgba(0, 0, 0, 0.06);
  color: #374151;
}

.uc-project-card:first-child .status-chip.published {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.project-status-copy {
  margin-top: 12px;
}

.project-footer--mine {
  display: flex;
  gap: 12px;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
}

.project-footer__label {
  font-size: 0.85rem;
  color: #9ca3af;
}

.project-footer__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.reject-tip {
  margin-top: 12px;
  color: #dc2626;
  font-size: 0.92rem;
}

.privacy-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  background: rgba(255, 217, 142, 0.16);
  color: #92400e;
}

.delete-link {
  color: #dc2626;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
}

.delete-link:hover {
  color: #b91c1c;
  text-decoration: underline;
}

@media (max-width: 768px) {
  .uc-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .project-grid {
    grid-template-columns: 1fr;
  }

  .project-footer--mine {
    flex-direction: column;
    align-items: flex-start;
  }

  .project-footer__actions {
    width: 100%;
  }

  .project-footer__actions a {
    flex: 1;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .project-card {
    padding: 16px;
  }

  .project-card__head {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
