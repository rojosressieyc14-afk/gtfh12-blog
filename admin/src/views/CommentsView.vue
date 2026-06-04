<template>
  <div class="comments-page">
    <section class="comments-hero">
      <div>
        <p class="login-tag">Comment Hub</p>
        <h1>评论管理</h1>
        <p class="login-text">集中查看文章评论与回复串，支持按内容、用户和文章标题搜索，并直接删除异常评论。</p>
      </div>
      <article class="comments-stat">
        <span>当前结果总数</span>
        <strong>{{ total }}</strong>
      </article>
    </section>

    <section class="comments-panel">
      <div class="review-head review-head--stack">
        <div>
          <p class="admin-label">评论列表</p>
          <h2>评论与回复串</h2>
        </div>
        <div class="toolbar-row">
          <input
            v-model.trim="keyword"
            class="filter-select filter-input"
            placeholder="搜索评论内容 / 用户 / 文章标题"
            @keyup.enter="changePage(1)"
          />
          <button @click="changePage(1)">搜索</button>
          <button class="role-btn" @click="resetFilters">清空</button>
        </div>
      </div>

      <div v-if="flash" class="page-flash">{{ flash }}</div>

      <div class="table-list">
        <article v-for="item in items" :key="item.id" class="table-card comment-admin-card">
          <div class="comment-admin-body">
            <div class="reply-item__head">
              <strong>{{ item.user?.username || `用户 #${item.userId}` }}</strong>
              <span>{{ formatDate(item.createdAt) }}</span>
            </div>
            <p class="table-note">文章：{{ item.article?.title || `文章 #${item.articleId}` }}</p>
            <p>{{ item.content }}</p>

            <div v-if="item.replies?.length" class="reply-thread">
              <article v-for="reply in item.replies" :key="reply.id" class="reply-item">
                <div class="reply-item__head">
                  <strong>{{ reply.user?.username || `用户 #${reply.userId}` }}</strong>
                  <div class="reply-item__meta">
                    <span>{{ formatDate(reply.createdAt) }}</span>
                    <button class="reply-delete" @click="removeReply(reply.id)">删除回复</button>
                  </div>
                </div>
                <p>{{ reply.content }}</p>
              </article>
            </div>
          </div>

          <div class="table-actions">
            <button class="btn-reject" @click="removeComment(item.id)">删除整串</button>
          </div>
        </article>
      </div>

      <div v-if="!items.length" class="compact-empty">当前没有匹配到评论记录。</div>

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
import { deleteAdminComment, getAdminComments } from "../api/dashboard";

const items = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = 10;
const keyword = ref("");
const flash = ref("");

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

function say(message) {
  flash.value = message;
  window.setTimeout(() => {
    if (flash.value === message) flash.value = "";
  }, 2200);
}

async function loadComments() {
  const { data } = await getAdminComments({ page: page.value, pageSize, keyword: keyword.value });
  items.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function changePage(nextPage) {
  page.value = nextPage;
  try {
    await loadComments();
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "加载评论失败。");
  }
}

async function removeComment(id) {
  if (!window.confirm("确定删除这条评论及其全部回复吗？")) return;
  try {
    await deleteAdminComment(id);
    await loadComments();
    say("评论已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除评论失败。");
  }
}

async function removeReply(id) {
  if (!window.confirm("确定只删除这条回复吗？")) return;
  try {
    await deleteAdminComment(id);
    await loadComments();
    say("回复已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除回复失败。");
  }
}

function resetFilters() {
  keyword.value = "";
  changePage(1);
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

onMounted(() => {
  changePage(1);
});
</script>

<style scoped>
.comments-page {
  min-height: 100vh;
  padding: 28px;
}

.comments-hero,
.comments-panel {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.comments-hero {
  padding: 28px;
  border-radius: 32px;
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(220px, 0.8fr);
  gap: 20px;
}

.comments-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.4rem);
  line-height: 1;
}

.comments-stat {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 19, 0.52);
  display: grid;
  gap: 10px;
}

.comments-stat span,
.comments-stat p {
  color: var(--soft);
}

.comments-stat strong {
  font-size: 2rem;
}

.comments-panel {
  margin-top: 20px;
  padding: 24px;
  border-radius: 28px;
}

.page-flash {
  margin-bottom: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 166, 77, 0.14);
  color: #ffe3c2;
}

.reply-item__meta {
  display: flex;
  gap: 10px;
  align-items: center;
}

@media (max-width: 960px) {
  .comments-page {
    padding: 20px;
  }

  .comments-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .comments-page {
    padding: 16px;
  }
  .comments-hero {
    padding: 20px;
  }
  .comments-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .comments-page {
    padding: 12px;
  }
  .comments-hero {
    padding: 16px;
  }
  .comments-panel {
    padding: 14px;
  }
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .reply-item__head {
    flex-direction: column;
    align-items: flex-start;
  }
  .reply-item__meta {
    width: 100%;
    justify-content: space-between;
  }
  .table-card {
    flex-direction: column;
    align-items: flex-start;
  }
  .table-actions {
    width: 100%;
  }
  .reply-thread {
    padding-left: 0;
  }
  .comment-admin-card {
    flex-direction: column;
  }
  .comment-admin-body {
    width: 100%;
  }
}
</style>
