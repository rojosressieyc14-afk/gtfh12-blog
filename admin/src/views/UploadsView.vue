<template>
  <div class="uploads-page">
    <section class="uploads-hero">
      <div>
        <p class="login-tag">Asset Room</p>
        <h1>上传资源库</h1>
        <p class="login-text">集中查看后台上传的图片资源，方便检查封面、清理冗余文件和排查引用问题。</p>
      </div>
      <article class="uploads-stat">
        <span>当前资源数量</span>
        <strong>{{ items.length }}</strong>
      </article>
    </section>

    <section class="uploads-panel">
      <div class="review-head review-head--stack">
        <div>
          <p class="admin-label">资源管理</p>
          <h2>图片与静态资源</h2>
        </div>
        <div class="toolbar-row">
          <input v-model.trim="keyword" class="filter-select filter-input" placeholder="搜索资源文件名" />
          <button @click="loadUploads">刷新资源</button>
        </div>
      </div>

      <div v-if="flash" class="page-flash">{{ flash }}</div>

      <div class="upload-grid">
        <article v-for="item in filteredItems" :key="item.name" class="upload-card upload-card--panel">
          <img :src="fileUrl(item.url)" :alt="item.name" />
          <strong>{{ item.name }}</strong>
          <span class="upload-meta">{{ formatSize(item.size) }} · {{ formatDate(item.time) }}</span>
          <div class="table-actions">
            <a class="role-btn link-btn" :href="fileUrl(item.url)" target="_blank" rel="noreferrer">查看原图</a>
            <button class="role-btn" @click="copyLink(item.url)">复制链接</button>
            <button class="btn-reject" @click="removeUpload(item.name)">删除</button>
          </div>
        </article>
      </div>

      <div v-if="!filteredItems.length" class="compact-empty">当前没有匹配到上传资源。</div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { deleteAdminUpload, getAdminUploads } from "../api/dashboard";

const items = ref([]);
const flash = ref("");
const keyword = ref("");
const filteredItems = computed(() =>
  items.value.filter((item) => item.name.toLowerCase().includes(keyword.value.toLowerCase()))
);

function say(message) {
  flash.value = message;
  window.setTimeout(() => {
    if (flash.value === message) flash.value = "";
  }, 2200);
}

async function loadUploads() {
  try {
    const { data } = await getAdminUploads();
    items.value = data.items || [];
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "加载资源失败。");
  }
}

async function removeUpload(name) {
  if (!window.confirm(`确定删除资源“${name}”吗？`)) return;
  try {
    await deleteAdminUpload(name);
    await loadUploads();
    say("资源已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除资源失败。");
  }
}

async function copyLink(url) {
  const link = fileUrl(url);
  try {
    await navigator.clipboard.writeText(link);
    say("资源链接已复制。");
  } catch (error) {
    say(error?.message || "复制链接失败。");
  }
}

function fileUrl(url) {
  if (import.meta.env.DEV) {
    return `http://localhost:8080${url}`;
  }
  const base = import.meta.env.VITE_BASE_PATH || "";
  return `${base}${url}`;
}

function formatSize(size) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

onMounted(loadUploads);
</script>

<style scoped>
.uploads-page {
  min-height: 100vh;
  padding: 28px;
}

.uploads-hero,
.uploads-panel {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.uploads-hero {
  padding: 28px;
  border-radius: 32px;
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(220px, 0.8fr);
  gap: 20px;
}

.uploads-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.4rem);
  line-height: 1;
}

.uploads-stat {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 19, 0.52);
  display: grid;
  gap: 10px;
}

.uploads-stat span {
  color: var(--soft);
}

.uploads-stat strong {
  font-size: 2rem;
}

.uploads-panel {
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

@media (max-width: 960px) {
  .uploads-page {
    padding: 20px;
  }

  .uploads-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .uploads-page {
    padding: 16px;
  }
  .uploads-hero {
    padding: 20px;
  }
  .uploads-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .uploads-page {
    padding: 12px;
  }
  .uploads-hero {
    padding: 16px;
  }
  .uploads-panel {
    padding: 14px;
  }
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .upload-grid {
    grid-template-columns: 1fr;
  }
  .upload-card {
    padding: 12px;
  }
  .table-actions {
    width: 100%;
    flex-direction: column;
  }
  .table-actions a,
  .table-actions button {
    width: 100%;
    text-align: center;
  }
}
</style>
