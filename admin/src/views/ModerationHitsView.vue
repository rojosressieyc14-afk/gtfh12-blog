<template>
  <div class="hits-page">
    <section class="hits-hero">
      <div>
        <p class="login-tag">Moderation Log</p>
        <h1>违禁词命中日志</h1>
        <p class="login-text">
          这里记录被系统拦截的文本提交，包含用户、触发场景、字段、命中词和内容片段，方便你持续补词、排查误伤并跟进自动封禁情况。
        </p>
      </div>
      <article class="hits-stat">
        <span>当前结果总数</span>
        <strong>{{ total }}</strong>
      </article>
    </section>

    <section class="hits-layout">
      <article class="hits-panel hits-panel--settings">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">风控设置</p>
            <h2>自动封禁阈值</h2>
          </div>
        </div>
        <div class="settings-card">
          <label>
            24 小时内累计命中次数
            <input v-model.number="banThreshold" class="filter-select" type="number" min="1" max="100" />
          </label>
          <p>普通用户在 24 小时内累计达到该次数后，系统会自动封禁并通知管理员。</p>
          <button @click="saveSettings">保存阈值</button>
        </div>
      </article>

      <article class="hits-panel">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">筛选</p>
            <h2>查看拦截明细</h2>
          </div>
          <div class="toolbar-row">
            <input
              v-model.trim="keyword"
              class="filter-select filter-input"
              placeholder="搜索用户 / 字段 / 违禁词 / 内容片段"
              @keyup.enter="changePage(1)"
            />
            <select v-model="scene" class="filter-select" @change="changePage(1)">
              <option value="">全部场景</option>
              <option value="register">注册</option>
              <option value="profile">个人资料</option>
              <option value="article_create">文章创建</option>
              <option value="article_update">文章编辑</option>
              <option value="article_submit">文章提交审核</option>
              <option value="article_review">文章审核</option>
              <option value="project_create">项目创建</option>
              <option value="project_update">项目编辑</option>
              <option value="project_submit">项目提交审核</option>
              <option value="comment_create">评论发布</option>
            </select>
            <label class="hits-toggle">
              <input v-model="autoBannedOnly" type="checkbox" @change="changePage(1)" />
              <span>只看已自动封禁用户</span>
            </label>
            <button @click="changePage(1)">搜索</button>
          </div>
        </div>

        <div class="hits-list">
          <article v-for="item in items" :key="item.id" class="hits-card">
            <div class="hits-card__head">
              <div>
                <strong>{{ item.user?.username || `用户 #${item.userId}` }}</strong>
                <p>{{ sceneLabel(item.scene) }} · {{ item.field }}</p>
              </div>
              <span>{{ formatDate(item.createdAt) }}</span>
            </div>
            <div class="hits-card__word">
              <label>命中词</label>
              <strong>{{ item.matchedWord }}</strong>
            </div>
            <div class="hits-card__snippet">
              <label>内容片段</label>
              <p>{{ item.snippet || "暂无内容片段" }}</p>
            </div>
            <p v-if="item.user?.banReason" class="hits-ban-reason">封禁原因：{{ item.user.banReason }}</p>
          </article>
        </div>

        <div v-if="!items.length" class="compact-empty">当前还没有违禁词命中记录。</div>

        <div class="pager-row">
          <button :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <span>第 {{ page }} / {{ totalPages }} 页</span>
          <button :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
        </div>
      </article>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { getModerationHits, getModerationSettings, updateModerationSettings } from "../api/dashboard";

const items = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;
const keyword = ref("");
const scene = ref("");
const autoBannedOnly = ref(false);
const banThreshold = ref(5);

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

function sceneLabel(value) {
  return {
    register: "注册",
    profile: "个人资料",
    article_create: "文章创建",
    article_update: "文章编辑",
    article_submit: "文章提交审核",
    article_review: "文章审核",
    project_create: "项目创建",
    project_update: "项目编辑",
    project_submit: "项目提交审核",
    comment_create: "评论发布"
  }[value] || value;
}

async function loadHits() {
  const { data } = await getModerationHits({
    page: page.value,
    pageSize,
    keyword: keyword.value,
    scene: scene.value,
    autoBanned: autoBannedOnly.value
  });
  items.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function loadSettings() {
  const { data } = await getModerationSettings();
  banThreshold.value = data.banThreshold || 5;
}

async function saveSettings() {
  await updateModerationSettings({ banThreshold: Number(banThreshold.value || 5) });
  await loadSettings();
}

async function changePage(nextPage) {
  page.value = nextPage;
  await loadHits();
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

onMounted(async () => {
  await Promise.all([loadHits(), loadSettings()]);
});
</script>

<style scoped>
.hits-page {
  min-height: 100vh;
  padding: 28px;
  color: inherit;
}

.hits-hero,
.hits-panel {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.hits-hero {
  padding: 28px;
  border-radius: 32px;
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(240px, 0.75fr);
  gap: 20px;
}

.hits-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.6rem);
  line-height: 1;
}

.hits-layout {
  margin-top: 20px;
  display: grid;
  grid-template-columns: minmax(280px, 0.7fr) minmax(0, 1.3fr);
  gap: 20px;
}

.hits-panel {
  padding: 24px;
  border-radius: 28px;
}

.hits-stat {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 19, 0.52);
  display: grid;
  gap: 10px;
}

.hits-stat span,
.hits-card p,
.settings-card p {
  color: var(--soft);
}

.hits-stat strong {
  font-size: 2rem;
}

.settings-card {
  display: grid;
  gap: 14px;
  padding: 18px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.03);
}

.settings-card label {
  display: grid;
  gap: 8px;
}

.settings-card button {
  border-radius: 16px;
  border: 1px solid var(--border);
  padding: 12px 14px;
  cursor: pointer;
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
}

.hits-toggle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
}

.hits-list {
  display: grid;
  gap: 14px;
}

.hits-card {
  padding: 18px;
  border-radius: 22px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.03);
}

.hits-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.hits-card__head strong {
  font-size: 1.08rem;
}

.hits-card__head p {
  margin: 8px 0 0;
}

.hits-card__word,
.hits-card__snippet {
  margin-top: 14px;
  display: grid;
  gap: 8px;
}

.hits-card__word label,
.hits-card__snippet label {
  color: #ffd79d;
  font-size: 0.88rem;
}

.hits-card__word strong,
.hits-ban-reason {
  color: #ffbcbc;
}

.hits-card__snippet p,
.hits-ban-reason {
  margin: 0;
  line-height: 1.6;
  word-break: break-word;
}

.hits-ban-reason {
  margin-top: 14px;
}

@media (max-width: 960px) {
  .hits-page {
    padding: 20px;
  }

  .hits-hero,
  .hits-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .hits-page {
    padding: 16px;
  }
  .hits-hero {
    padding: 20px;
  }
  .hits-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .hits-page {
    padding: 12px;
  }
  .hits-hero {
    padding: 16px;
  }
  .hits-panel {
    padding: 14px;
  }
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .hits-card__head {
    flex-direction: column;
  }
  .hits-toggle {
    width: 100%;
    justify-content: center;
  }
  .settings-card {
    padding: 14px;
  }
}
</style>
