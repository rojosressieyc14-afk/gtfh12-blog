<template>
  <div class="logs-page">
    <section class="logs-hero">
      <div>
        <p class="login-tag">Ops Trail</p>
        <h1>操作日志</h1>
        <p class="login-text">记录后台关键动作，便于复盘审核、删文、删评、调整角色和风控设置等操作。</p>
      </div>
      <article class="logs-stat">
        <span>当前结果总数</span>
        <strong>{{ total }}</strong>
      </article>
    </section>

    <section class="logs-panel">
      <div class="review-head review-head--stack">
        <div>
          <p class="admin-label">检索日志</p>
          <h2>后台操作轨迹</h2>
        </div>
        <div class="toolbar-row">
          <input
            v-model.trim="keyword"
            class="filter-select filter-input"
            placeholder="搜索管理员 / 动作 / 描述"
            @keyup.enter="changePage(1)"
          />
          <select v-model="action" class="filter-select" @change="changePage(1)">
            <option value="">全部动作</option>
            <option v-for="item in actionOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
          <input v-model="dateFrom" class="filter-select" type="date" @change="changePage(1)" />
          <input v-model="dateTo" class="filter-select" type="date" @change="changePage(1)" />
          <button @click="changePage(1)">搜索</button>
          <button class="role-btn" @click="resetFilters">清空</button>
        </div>
      </div>

      <div class="logs-list">
        <article v-for="item in items" :key="item.id" class="logs-card">
          <div class="logs-card__head">
            <div>
              <strong>{{ item.operator?.username || `管理员 #${item.operatorId}` }}</strong>
              <p>{{ actionLabel(item.action) }} · {{ targetLabel(item.targetType, item.targetId) }}</p>
            </div>
            <span class="logs-time">{{ formatDate(item.createdAt) }}</span>
          </div>
          <p class="logs-desc">{{ item.description }}</p>
        </article>
      </div>

      <div v-if="!items.length" class="compact-empty">当前没有匹配到操作日志。</div>

      <div class="pager-row">
        <button class="role-btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <span>第 {{ page }} / {{ totalPages }} 页</span>
        <button class="role-btn" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { getAdminLogs } from "../api/dashboard";

const items = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;
const keyword = ref("");
const action = ref("");
const dateFrom = ref("");
const dateTo = ref("");

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));
const actionOptions = [
  { value: "review_article", label: "审核文章" },
  { value: "review_project", label: "审核项目" },
  { value: "delete_article", label: "删除文章" },
  { value: "delete_project", label: "删除项目" },
  { value: "delete_comment", label: "删除评论" },
  { value: "delete_upload", label: "删除资源" },
  { value: "publish_article", label: "直接发布文章" },
  { value: "publish_project", label: "直接发布项目" },
  { value: "update_user_role", label: "修改用户角色" },
  { value: "update_user_status", label: "修改用户状态" },
  { value: "create_sensitive_word", label: "新增敏感词" },
  { value: "delete_sensitive_word", label: "删除敏感词" },
  { value: "update_moderation_settings", label: "更新风控设置" },
  { value: "update_project_meta", label: "调整项目排序" }
];

async function loadLogs() {
  const { data } = await getAdminLogs({
    page: page.value,
    pageSize,
    keyword: keyword.value,
    action: action.value || undefined,
    dateFrom: dateFrom.value || undefined,
    dateTo: dateTo.value || undefined
  });
  items.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function changePage(nextPage) {
  page.value = nextPage;
  await loadLogs();
}

function actionLabel(value) {
  return {
    review_article: "审核文章",
    review_project: "审核项目",
    delete_article: "删除文章",
    delete_project: "删除项目",
    delete_comment: "删除评论",
    delete_upload: "删除资源",
    publish_article: "直接发布文章",
    publish_project: "直接发布项目",
    update_user_role: "修改用户角色",
    update_user_status: "修改用户状态",
    create_sensitive_word: "新增敏感词",
    delete_sensitive_word: "删除敏感词",
    update_moderation_settings: "更新风控设置",
    update_project_meta: "调整项目排序"
  }[value] || value;
}

function targetLabel(type, id) {
  const map = {
    article: "文章",
    project: "项目",
    comment: "评论",
    user: "用户",
    upload: "资源",
    sensitive_word: "敏感词",
    system_setting: "系统设置"
  };
  return `${map[type] || type} #${id || 0}`;
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

function resetFilters() {
  keyword.value = "";
  action.value = "";
  dateFrom.value = "";
  dateTo.value = "";
  changePage(1);
}

onMounted(() => {
  changePage(1);
});
</script>

<style scoped>
.logs-page {
  min-height: 100vh;
  padding: 28px;
}

.logs-hero,
.logs-panel {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.logs-hero {
  padding: 28px;
  border-radius: 32px;
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(220px, 0.8fr);
  gap: 20px;
}

.logs-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.4rem);
  line-height: 1;
}

.logs-stat {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 19, 0.52);
  display: grid;
  gap: 10px;
}

.logs-stat span,
.logs-desc,
.logs-card__head p {
  color: var(--soft);
}

.logs-stat strong {
  font-size: 2rem;
}

.logs-panel {
  margin-top: 20px;
  padding: 24px;
  border-radius: 28px;
}

.logs-list {
  display: grid;
  gap: 14px;
}

.logs-card {
  padding: 18px;
  border-radius: 22px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.03);
}

.logs-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.logs-card__head strong {
  font-size: 1.05rem;
}

.logs-card__head p,
.logs-desc {
  margin: 8px 0 0;
}

.logs-time {
  color: var(--soft);
  white-space: nowrap;
}

@media (max-width: 960px) {
  .logs-page {
    padding: 20px;
  }

  .logs-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .logs-page {
    padding: 16px;
  }
  .logs-hero {
    padding: 20px;
  }
  .logs-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .logs-page {
    padding: 12px;
  }
  .logs-hero {
    padding: 16px;
  }
  .logs-panel {
    padding: 14px;
  }
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .logs-card__head {
    flex-direction: column;
  }
}
</style>
