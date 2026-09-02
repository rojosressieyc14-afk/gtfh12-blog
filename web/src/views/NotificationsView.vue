<template>
  <section class="uc-notifications">
    <section class="uc-hero uc-hero--notifications">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">站内信箱</p>
        <h2 class="uc-hero__title">通知中心</h2>
        <p class="uc-hero__text">审核结果、评论互动和系统提醒都会集中显示。</p>
      </div>
      <div class="uc-hero__actions">
        <button class="ghost-btn uc-hero__ghost" @click="markAllRead">全部标为已读</button>
      </div>
    </section>

    <div class="uc-notifications__filters">
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
          <span class="notification-date">{{ formatDate(item.createdAt) }}</span>
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
import { formatDate } from "../utils/date";

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

onMounted(loadData);
</script>

<style scoped>
.uc-notifications {
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
  grid-template-columns: minmax(0, 1.4fr) auto;
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

.uc-hero__ghost {
  color: #374151;
  border-color: rgba(0, 0, 0, 0.15);
  background: rgba(255, 255, 255, 0.7);
}

.uc-hero__ghost:hover {
  background: rgba(255, 255, 255, 0.9);
}

.uc-notifications__filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.uc-notifications__filters .ghost-btn.active {
  background: rgba(249, 115, 22, 0.1);
  border-color: rgba(249, 115, 22, 0.25);
  color: #f97316;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-card {
  padding: 18px 20px;
  border-radius: 18px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.notification-card:hover {
  border-color: rgba(0, 0, 0, 0.12);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
}

.notification-card.unread {
  border-color: rgba(249, 115, 22, 0.2);
  background: rgba(255, 255, 255, 0.8);
}

.notification-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.notification-title {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.notification-title strong {
  color: #1d1d1f;
}

.notification-date {
  font-size: 0.82rem;
  color: #9ca3af;
  white-space: nowrap;
}

.notification-card p {
  margin: 0;
  font-size: 0.92rem;
  color: #48484a;
  line-height: 1.5;
}

.status-chip.pending {
  background: rgba(249, 115, 22, 0.12);
  color: #c2410c;
}

.status-chip.success {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.status-chip.reject {
  background: rgba(220, 38, 38, 0.12);
  color: #b91c1c;
}

@media (max-width: 768px) {
  .uc-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .uc-notifications__filters {
    flex-direction: column;
  }

  .uc-notifications__filters .ghost-btn {
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

  .uc-notifications__filters {
    gap: 6px;
  }
}
</style>
