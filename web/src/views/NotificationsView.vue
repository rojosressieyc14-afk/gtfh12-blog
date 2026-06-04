<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">站内信箱</p>
        <h3>通知中心</h3>
      </div>
      <button class="ghost-btn" @click="markAllRead">全部标为已读</button>
    </div>

    <div class="services-grid">
      <article class="panel-card service-card">
        <p class="eyebrow">范围</p>
        <h4>审核结果、评论互动和系统提醒都会集中显示</h4>
        <p>这里适合快速查看内容审核进度，以及文章和评论带来的后续互动。</p>
      </article>
      <article class="panel-card service-card">
        <p class="eyebrow">处理</p>
        <h4>点开通知后直接跳到对应页面</h4>
        <p>如果通知带有操作入口，点击后会自动标记已读，并进入对应页面继续处理。</p>
      </article>
    </div>

    <div class="filter-row">
      <button class="ghost-btn" :class="{ active: selectedType === '' }" @click="setType('')">全部</button>
      <button class="ghost-btn" :class="{ active: selectedType === 'article_review' }" @click="setType('article_review')">文章审核</button>
      <button class="ghost-btn" :class="{ active: selectedType === 'project_review' }" @click="setType('project_review')">项目审核</button>
      <button class="ghost-btn" :class="{ active: selectedType === 'article_comment' }" @click="setType('article_comment')">新评论</button>
      <button class="ghost-btn" :class="{ active: selectedType === 'comment_reply' }" @click="setType('comment_reply')">评论回复</button>
      <button class="ghost-btn" :class="{ active: unreadOnly }" @click="toggleUnreadOnly">
        {{ unreadOnly ? "只看未读" : "显示全部" }}
      </button>
    </div>

    <div v-if="items.length" class="notification-list">
      <article
        v-for="item in items"
        :key="item.id"
        class="notification-card"
        :class="{ unread: !item.isRead }"
        @click="openNotification(item)"
      >
        <div class="notification-head">
          <div class="notification-title">
            <span class="status-chip" :class="notificationTypeClass(item.type)">{{ notificationTypeLabel(item.type) }}</span>
            <strong>{{ item.title }}</strong>
          </div>
          <span>{{ formatDate(item.createdAt) }}</span>
        </div>
        <p>{{ item.content }}</p>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>暂时还没有通知</h4>
      <p>评论回复、审核结果和系统消息后续都会显示在这里。</p>
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
import { useRouter } from "vue-router";
import { listNotifications, markAllNotificationsRead, markNotificationRead } from "../api/notification";

const router = useRouter();
const items = ref([]);
const selectedType = ref("");
const unreadOnly = ref(false);
const page = ref(1);
const pageSize = 12;
const total = ref(0);

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

async function loadData() {
  const { data } = await listNotifications({
    type: selectedType.value || undefined,
    unreadOnly: unreadOnly.value || undefined,
    page: page.value,
    pageSize
  });
  items.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function markAllRead() {
  await markAllNotificationsRead();
  await loadData();
}

async function openNotification(item) {
  if (!item.isRead) {
    await markNotificationRead(item.id);
    item.isRead = true;
  }
  if (item.actionUrl) {
    await router.push(item.actionUrl);
  }
}

async function setType(type) {
  selectedType.value = type;
  page.value = 1;
  await loadData();
}

async function toggleUnreadOnly() {
  unreadOnly.value = !unreadOnly.value;
  page.value = 1;
  await loadData();
}

async function goToPage(nextPage) {
  page.value = nextPage;
  await loadData();
}

function notificationTypeLabel(type) {
  return {
    article_review: "文章审核",
    project_review: "项目审核",
    article_comment: "新评论",
    comment_reply: "评论回复",
    moderation: "风控提醒",
    system: "系统消息"
  }[type] || type;
}

function notificationTypeClass(type) {
  return {
    pending: type === "article_review" || type === "project_review",
    success: type === "article_comment" || type === "comment_reply",
    reject: type === "moderation"
  };
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(loadData);
</script>

<style scoped>
.notification-title {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.notification-card {
  cursor: pointer;
}

.notification-card.unread {
  border-color: rgba(255, 209, 102, 0.22);
  background: rgba(255, 209, 102, 0.06);
}

.ghost-btn.active {
  background: rgba(255, 209, 102, 0.12);
  border-color: rgba(255, 209, 102, 0.28);
}

.status-chip.pending {
  background: rgba(255, 166, 77, 0.18);
}

.status-chip.success {
  background: rgba(74, 222, 128, 0.18);
}

.status-chip.reject {
  background: rgba(248, 113, 113, 0.18);
}

@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
  }

  .filter-row .ghost-btn {
    width: 100%;
    justify-content: center;
  }

  .notification-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .notification-card {
    padding: 14px;
  }
}

@media (max-width: 480px) {
  .notification-title {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }

  .filter-row {
    gap: 8px;
  }
}
</style>
