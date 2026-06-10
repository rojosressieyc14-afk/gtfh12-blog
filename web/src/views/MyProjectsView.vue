<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">作品管理</p>
        <h3>我的项目</h3>
      </div>
      <router-link class="solid-btn" to="/project-editor">新建项目</router-link>
    </div>

    <div class="services-grid">
      <article class="panel-card service-card">
        <p class="eyebrow">管理</p>
        <h4>统一查看草稿、审核中和已发布项目</h4>
      </article>
      <article class="panel-card service-card">
        <p class="eyebrow">建议</p>
        <h4>优先补齐亮点、难点和结果</h4>
      </article>
    </div>

    <div v-if="projects.length" class="project-grid">
      <article v-for="item in projects" :key="item.id" class="project-card panel-card">
        <div class="project-card__head">
          <div>
            <h3>{{ item.title }}</h3>
            <p class="table-note">{{ formatDate(item.updatedAt) }}</p>
          </div>
          <span class="status-chip" :class="item.status">{{ labelMap[item.status] || item.status }}</span>
        </div>

        <p class="detail-summary">{{ item.summary || "暂无摘要" }}</p>

        <div v-if="item.techStacks?.length" class="tag-row">
          <span v-for="stack in item.techStacks.slice(0, 6)" :key="stack" class="tag-chip"># {{ stack }}</span>
        </div>

        <p class="detail-summary project-status-copy">{{ statusCopy(item) }}</p>
        <p v-if="item.rejectReason" class="reject-tip">驳回原因：{{ item.rejectReason }}</p>

        <footer class="project-footer project-footer--mine">
          <span>{{ item.isFeatured ? "首页精选项目" : "普通展示项目" }}</span>
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

const projects = ref([]);
const page = ref(1);
const pageSize = 8;
const total = ref(0);

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
    alert("删除失败：" + (e.response?.data?.message || e.message));
  }
}

async function goToPage(nextPage) {
  page.value = nextPage;
  await loadMine();
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(loadMine);
</script>

<style scoped>
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

.project-footer__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.reject-tip {
  margin-top: 12px;
  color: #fecaca;
  font-size: 0.92rem;
}

.delete-link {
  color: #f87171;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
}

.delete-link:hover {
  color: #ef4444;
  text-decoration: underline;
}

@media (max-width: 768px) {
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
