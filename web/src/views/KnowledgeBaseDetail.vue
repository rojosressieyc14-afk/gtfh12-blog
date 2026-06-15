<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">知识库</p>
        <h3>{{ kb?.name || "加载中..." }}</h3>
        <p v-if="kb?.description" class="detail-summary">{{ kb.description }}</p>
      </div>
      <router-link class="ghost-btn" :to="`/user-center/knowledge-base`">返回列表</router-link>
    </div>

    <div class="kb-tabs">
      <button class="kb-tab" :class="{ 'kb-tab--active': tab === 'notes' }" @click="tab = 'notes'">笔记列表 ({{ docs.length }})</button>
      <button class="kb-tab" :class="{ 'kb-tab--active': tab === 'query' }" @click="tab = 'query'">AI 检索</button>
    </div>

    <div v-if="tab === 'notes'">
      <div class="kb-filters">
        <select v-model="filterCategory" class="field-select field-select--sm">
          <option value="">全部分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <input v-model="filterTag" class="field-input field-input--sm" placeholder="按标签筛选..." />
        <router-link class="solid-btn" :to="`/user-center/knowledge-base/${$route.params.id}/editor`">新建笔记</router-link>
      </div>

      <div v-if="filteredDocs.length" class="note-grid">
        <article v-for="doc in filteredDocs" :key="doc.id" class="note-card panel-card">
          <div class="note-card__head">
            <div>
              <div class="note-card__title-row">
                <h4>{{ doc.title || "无标题" }}</h4>
                <span v-if="doc.isPublic" class="badge badge-public">公开</span>
                <span v-else class="badge badge-private">私密</span>
              </div>
              <div class="note-card__meta">
                <span v-if="doc.category">{{ doc.category.name }}</span>
                <span>{{ formatDate(doc.createdAt) }}</span>
                <span v-if="doc.viewCount">{{ doc.viewCount }} 次阅读</span>
              </div>
            </div>
            <div class="note-card__actions">
              <router-link class="inline-link" :to="`/user-center/knowledge-base/${$route.params.id}/editor/${doc.id}`">编辑</router-link>
              <button class="inline-link delete-link" @click="handleDeleteDoc(doc)">删除</button>
            </div>
          </div>

          <p class="note-card__preview">{{ truncate(doc.content, 260) }}</p>

          <div v-if="doc.tags?.length" class="tag-row">
            <span v-for="tag in doc.tags" :key="tag.id" class="tag-chip"># {{ tag.name }}</span>
          </div>
        </article>
      </div>

      <div v-else class="empty-panel">
        <h4>{{ docs.length ? "没有匹配的笔记" : "还没有笔记" }}</h4>
        <p>{{ docs.length ? "尝试调整筛选条件" : "创建你的第一篇笔记。" }}</p>
      </div>
    </div>

    <div v-if="tab === 'query'" class="kb-query-panel">
      <div class="query-box">
        <textarea v-model="question" class="field-area field-area--small" placeholder="输入问题，AI 会基于知识库内容回答..." rows="3"></textarea>
        <button class="solid-btn" style="margin-top:8px" :disabled="!question.trim() || querying" @click="handleQuery">
          {{ querying ? "检索中..." : "检索" }}
        </button>
      </div>

      <div v-if="queryResult" class="query-result">
        <div class="query-answer">
          <p class="eyebrow">回答</p>
          <p>{{ queryResult.answer }}</p>
        </div>
        <div v-if="queryResult.sources?.length" class="query-sources">
          <p class="eyebrow">参考来源 ({{ queryResult.sources.length }})</p>
          <div v-for="(src, i) in queryResult.sources" :key="i" class="source-item">
            <p class="source-title">{{ src.title || "片段 " + (i + 1) }}</p>
            <p class="source-text">{{ truncate(src.content, 150) }}</p>
            <p class="source-score">相关度: {{ (src.score * 100).toFixed(1) }}%</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { deleteDocument, getKnowledgeBase, listDocuments, queryKnowledgeBase } from "../api/knowledgeBase";
import { getMetadata } from "../api/meta";

const route = useRoute();
const kb = ref(null);
const docs = ref([]);
const categories = ref([]);
const tab = ref("notes");
const filterCategory = ref("");
const filterTag = ref("");
const question = ref("");
const querying = ref(false);
const queryResult = ref(null);

const filteredDocs = computed(() => {
  let list = docs.value;
  if (filterCategory.value) {
    list = list.filter((d) => String(d.categoryId) === filterCategory.value);
  }
  if (filterTag.value) {
    const q = filterTag.value.toLowerCase();
    list = list.filter((d) => (d.tags || []).some((t) => t.name.toLowerCase().includes(q)));
  }
  return list;
});

async function load() {
  try {
    const [{ data: kbData }, { data: docData }, { data: metaData }] = await Promise.all([
      getKnowledgeBase(route.params.id),
      listDocuments(route.params.id),
      getMetadata().catch(() => ({ data: { categories: [] } })),
    ]);
    kb.value = kbData.item;
    docs.value = docData.items || [];
    categories.value = metaData.categories || [];
  } catch (e) {
    alert("加载失败：" + (e?.response?.data?.message || e.message));
  }
}

async function handleDeleteDoc(doc) {
  if (!confirm(`确定要删除笔记「${doc.title || "无标题"}」吗？`)) return;
  try {
    await deleteDocument(route.params.id, doc.id);
    await load();
  } catch (e) {
    alert("删除失败：" + (e?.response?.data?.message || e.message));
  }
}

async function handleQuery() {
  if (!question.value.trim()) return;
  querying.value = true;
  queryResult.value = null;
  try {
    const { data } = await queryKnowledgeBase(route.params.id, { question: question.value.trim() });
    queryResult.value = data.result;
  } catch (e) {
    alert(e?.response?.data?.message || "查询失败");
  } finally {
    querying.value = false;
  }
}

function truncate(text, max) {
  if (!text) return "";
  return text.length > max ? text.slice(0, max) + "..." : text;
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(load);
</script>

<style scoped>
.kb-tabs {
  display: flex;
  gap: 4px;
  margin: 20px 0;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.08));
  padding-bottom: 0;
}

.kb-tab {
  padding: 10px 20px;
  border: none;
  background: none;
  color: var(--text-soft, #a0aec0);
  cursor: pointer;
  font-size: 0.95rem;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: color 0.2s, border-color 0.2s;
}

.kb-tab:hover {
  color: var(--text, #f7f3ea);
}

.kb-tab--active {
  color: #f97316;
  border-bottom-color: #f97316;
  font-weight: 600;
}

.kb-filters {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.field-select--sm,
.field-input--sm {
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--border, rgba(255,255,255,0.12));
  background: var(--panel, rgba(10,14,19,0.8));
  color: var(--text, #f7f3ea);
  font-size: 0.88rem;
}

.field-select--sm {
  min-width: 130px;
}

.field-input--sm {
  min-width: 160px;
}

.note-grid {
  display: grid;
  gap: 14px;
}

.note-card {
  padding: 20px;
  border-radius: 20px;
}

.note-card__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 10px;
}

.note-card__title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.note-card__title-row h4 {
  margin: 0;
  font-size: 1.1rem;
}

.note-card__meta {
  display: flex;
  gap: 12px;
  font-size: 0.82rem;
  color: var(--text-soft, #a0aec0);
  margin-top: 4px;
}

.note-card__actions {
  display: flex;
  gap: 12px;
  white-space: nowrap;
}

.note-card__preview {
  color: var(--text-soft, #a0aec0);
  font-size: 0.92rem;
  line-height: 1.5;
  margin-bottom: 10px;
}

.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
}

.badge-public {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
}

.badge-private {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
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

.kb-query-panel {
  margin-top: 16px;
}

.query-box {
  max-width: 600px;
}

.query-result {
  margin-top: 20px;
}

.query-answer {
  padding: 16px;
  border-radius: 18px;
  border: 1px solid var(--border, rgba(255,255,255,0.12));
  background: var(--panel, rgba(255,255,255,0.06));
  margin-bottom: 16px;
}

.query-sources {
  display: grid;
  gap: 10px;
}

.source-item {
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--border, rgba(255,255,255,0.08));
  background: rgba(255,255,255,0.03);
}

.source-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.source-text {
  font-size: 0.88rem;
  color: var(--soft, rgba(242,239,232,0.7));
}

.source-score {
  font-size: 0.8rem;
  color: var(--accent, #ffd98e);
  margin-top: 4px;
}

@media (max-width: 768px) {
  .kb-filters {
    flex-direction: column;
    align-items: stretch;
  }

  .note-card__head {
    flex-direction: column;
  }
}
</style>
