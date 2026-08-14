<template>
  <section class="content-section">
    <article v-if="note" class="kb-note-article">
      <div class="kb-note-header">
        <p class="eyebrow">公开笔记</p>
        <div class="kb-note-meta">
          <span v-if="note.user">作者: {{ note.user.username }}</span>
          <span>{{ formatDate(note.createdAt) }}</span>
          <span v-if="note.category">分类: {{ note.category.name }}</span>
          <span>{{ note.viewCount }} 次阅读</span>
        </div>
        <div v-if="note.tags?.length" class="tag-row">
          <span v-for="tag in note.tags" :key="tag.id" class="tag-chip"># {{ tag.name }}</span>
        </div>
      </div>

      <h1 class="kb-note-title">{{ note.title }}</h1>

      <div v-highlight class="markdown-body" v-html="renderedContent"></div>

      <div class="kb-note-footer">
        <router-link class="ghost-btn" to="/">返回首页</router-link>
      </div>
    </article>

    <div v-else-if="loading" class="empty-panel">
      <h4>加载中...</h4>
    </div>

    <div v-else class="empty-panel">
      <h4>笔记不存在或未公开</h4>
      <router-link class="solid-btn" to="/">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { getPublicNote } from "../api/knowledgeBase";

const route = useRoute();
const note = ref(null);
const loading = ref(true);

const renderedContent = computed(() => {
  if (!note.value?.content) return "";
  return DOMPurify.sanitize(marked.parse(note.value.content));
});

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(async () => {
  try {
    const { data } = await getPublicNote(route.params.id);
    note.value = data.item;
  } catch {
    note.value = null;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.kb-note-article {
  max-width: 800px;
  margin: 0 auto;
}

.kb-note-header {
  margin-bottom: 32px;
}

.kb-note-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.88rem;
  color: var(--text-soft, #a0aec0);
  margin: 8px 0;
}

.kb-note-title {
  font-size: 2rem;
  line-height: 1.3;
  margin-bottom: 24px;
}

.kb-note-footer {
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border, rgba(255,255,255,0.08));
}
</style>
