<template>
  <div class="kb-detail-layout">
    <aside class="kb-detail-sidebar">
      <div class="kb-sidebar-header">
        <h4>{{ kb?.name || '知识库' }}</h4>
        <button class="kb-sidebar-add" @click="createNewDoc">+</button>
      </div>
      <div class="kb-search-box">
        <div class="kb-search-input-wrap">
          <svg class="kb-search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
          <input
            v-model="searchKeyword"
            type="text"
            class="kb-search-input"
            placeholder="搜索文档..."
            @input="onSearchInput"
            @keydown.escape="clearSearch"
          />
          <button v-if="searchKeyword" class="kb-search-clear" @click="clearSearch">×</button>
        </div>
        <div v-if="searchResults.length > 0" class="kb-search-results">
          <div
            v-for="r in searchResults"
            :key="r.id"
            class="kb-search-result-item"
            @click="selectSearchResult(r)"
          >
            <div class="kb-search-result-title">{{ r.title || '无标题' }}</div>
            <div class="kb-search-result-snippet" v-html="highlightSnippet(r.snippet)"></div>
            <div class="kb-search-result-meta">
              <span v-if="r.isPublic" class="badge badge-public">公开</span>
              <span v-else class="badge badge-private">私密</span>
              <span>{{ r.createdAt }}</span>
            </div>
          </div>
        </div>
        <div v-else-if="searchKeyword && !searchLoading" class="kb-search-empty">
          没有找到匹配的文档
        </div>
      </div>
      <div class="kb-sidebar-tree">
        <DocTree
          :tree="docTree"
          :active-id="currentDoc?.id"
          @select="onSelectDoc"
        />
      </div>
      <div class="kb-sidebar-footer">
        <button class="kb-sidebar-btn" @click="createNewDoc">新建文档</button>
      </div>
    </aside>

    <main class="kb-detail-content">
      <div v-if="loading" class="kb-doc-empty">
        <h4>加载中...</h4>
      </div>
      <div v-else-if="errorMessage" class="kb-doc-empty">
        <h4>出错了</h4>
        <p>{{ errorMessage }}</p>
      </div>
      <div v-else-if="currentDoc" class="kb-doc-view">
        <div class="kb-doc-meta">
          <span v-if="currentDoc.category" class="kb-doc-category">{{ currentDoc.category.name }}</span>
          <span class="kb-doc-date">{{ formatDate(currentDoc.createdAt) }}</span>
          <span class="kb-doc-views">{{ currentDoc.viewCount }} 次阅读</span>
          <span v-if="currentDoc.isPublic" class="badge badge-public">公开</span>
          <span v-else class="badge badge-private">私密</span>
        </div>
        <h1 class="kb-doc-title">{{ currentDoc.title || '无标题' }}</h1>
        <div v-highlight class="markdown-body" v-html="renderedContent"></div>
        <div class="kb-doc-actions">
          <router-link class="ghost-btn" :to="`/user-center/knowledge-base/${$route.params.id}/editor/${currentDoc.id}`">编辑</router-link>
          <button class="ghost-btn" @click="handleDeleteDoc(currentDoc)">删除</button>
        </div>
      </div>
      <div v-else class="kb-doc-empty">
        <h4>选择一篇文档</h4>
        <p>从左栏选择一个文档，或创建一个新文档</p>
      </div>
    </main>

    <aside v-if="currentDoc" class="kb-detail-toc">
      <DocToc
        :items="docToc"
        :active-id="activeHeadingId"
        @jump="scrollToHeading"
      />
    </aside>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { getKnowledgeBase, getDocumentTree, listDocuments, deleteDocument, searchDocuments } from "../api/knowledgeBase";
import DocTree from "../components/DocTree.vue";
import DocToc from "../components/DocToc.vue";

const route = useRoute();
const router = useRouter();
const kb = ref(null);
const docTree = ref([]);
const currentDoc = ref(null);
const activeHeadingId = ref(null);
const allDocs = ref([]);
const loading = ref(true);
const errorMessage = ref("");
const searchKeyword = ref("");
const searchResults = ref([]);
const searchLoading = ref(false);
let searchTimer = null;

function onSearchInput() {
  clearTimeout(searchTimer);
  if (!searchKeyword.value.trim()) {
    searchResults.value = [];
    return;
  }
  searchTimer = setTimeout(performSearch, 300);
}

async function performSearch() {
  searchLoading.value = true;
  try {
    const res = await searchDocuments(route.params.id, searchKeyword.value.trim());
    searchResults.value = res.data.items || [];
  } catch (e) {
    searchResults.value = [];
  } finally {
    searchLoading.value = false;
  }
}

function clearSearch() {
  searchKeyword.value = "";
  searchResults.value = [];
}

function selectSearchResult(r) {
  const doc = allDocs.value.find((d) => Number(d.id) === Number(r.id));
  if (doc) currentDoc.value = doc;
  clearSearch();
}

function highlightSnippet(snippet) {
  if (!searchKeyword.value.trim()) return snippet;
  const kw = searchKeyword.value.trim().replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return snippet.replace(new RegExp(`(${kw})`, "gi"), '<mark>$1</mark>');
}

function buildHeadingId(text) {
  return text.toLowerCase().replace(/[^\w\u4e00-\u9fff]+/g, "-").replace(/(^-|-$)/g, "");
}

const renderedContent = computed(() => {
  if (!currentDoc.value?.content) return "";
  const renderer = new marked.Renderer();
  renderer.heading = (token) => {
    const text = token.text || "";
    const id = buildHeadingId(text);
    return `<h${token.depth} id="${id}">${text}</h${token.depth}>`;
  };
  return DOMPurify.sanitize(marked.parse(currentDoc.value.content, { renderer }));
});

const docToc = computed(() => {
  if (!currentDoc.value?.content) return [];
  const headingRe = /^(#{1,4})\s+(.+)$/gm;
  const items = [];
  let match;
  while ((match = headingRe.exec(currentDoc.value.content)) !== null) {
    const level = match[1].length;
    const text = match[2];
    const id = buildHeadingId(text);
    items.push({ id, text, level });
  }
  return items;
});

let tocObserver = null;

function scrollToHeading(id) {
  const el = document.getElementById(id);
  if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
}

function setupTocObserver() {
  tocObserver?.disconnect();
  const ids = docToc.value.map((i) => i.id);
  if (!ids.length) return;
  tocObserver = new IntersectionObserver(
    (entries) => {
      for (const e of entries) {
        if (e.isIntersecting) {
          activeHeadingId.value = e.target.id;
          break;
        }
      }
    },
    { rootMargin: "-80px 0px -60% 0px" }
  );
  for (const id of ids) {
    const el = document.getElementById(id);
    if (el) tocObserver.observe(el);
  }
}

function onSelectDoc(item) {
  const isFolder = item.children && item.children.length > 0;
  if (!isFolder) {
    const doc = allDocs.value.find((d) => Number(d.id) === Number(item.id));
    if (doc) currentDoc.value = doc;
  }
}

function createNewDoc() {
  router.push(`/user-center/knowledge-base/${route.params.id}/editor`);
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

async function handleDeleteDoc(doc) {
  if (!confirm(`确定要删除笔记「${doc.title || "无标题"}」吗？`)) return;
  try {
    await deleteDocument(route.params.id, doc.id);
    if (Number(currentDoc.value?.id) === Number(doc.id)) currentDoc.value = null;
    await load();
  } catch (e) {
    errorMessage.value = "删除失败：" + (e?.response?.data?.message || e.message);
  }
}

async function load() {
  loading.value = true;
  errorMessage.value = "";
  try {
    const [kbRes, treeRes, docRes] = await Promise.all([
      getKnowledgeBase(route.params.id),
      getDocumentTree(route.params.id),
      listDocuments(route.params.id),
    ]);
    kb.value = kbRes.data.item;
    docTree.value = treeRes.data.tree || [];
    allDocs.value = docRes.data.items || [];
  } catch (e) {
    errorMessage.value = "加载失败：" + (e?.response?.data?.message || e.message);
  } finally {
    loading.value = false;
  }
}

watch(currentDoc, () => nextTick(setupTocObserver), { flush: "post" });
onMounted(load);
onUnmounted(() => tocObserver?.disconnect());
</script>

<style scoped>
.kb-detail-layout {
  display: flex;
  height: calc(100vh - 140px);
  gap: 0;
  border-radius: 20px;
  border: 1px solid var(--border, rgba(255,255,255,0.1));
  overflow: hidden;
  background: rgba(255,255,255,0.03);
}

.kb-detail-sidebar {
  width: 260px;
  min-width: 260px;
  border-right: 1px solid var(--border, rgba(255,255,255,0.1));
  display: flex;
  flex-direction: column;
  background: rgba(255,255,255,0.02);
}

.kb-sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.08));
}
.kb-sidebar-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}
.kb-sidebar-add {
  background: none;
  border: none;
  color: rgba(246,241,234,0.5);
  font-size: 20px;
  cursor: pointer;
  padding: 0 4px;
  line-height: 1;
}
.kb-sidebar-add:hover {
  color: var(--accent);
}

.kb-sidebar-tree {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.kb-sidebar-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border, rgba(255,255,255,0.08));
}
.kb-sidebar-btn {
  width: 100%;
  padding: 8px;
  border-radius: 8px;
  border: 1px dashed var(--border, rgba(255,255,255,0.2));
  background: none;
  color: rgba(246,241,234,0.5);
  cursor: pointer;
  font-size: 13px;
  transition: color 0.15s, border-color 0.15s;
}
.kb-sidebar-btn:hover {
  color: rgba(246,241,234,0.8);
  border-color: var(--accent);
}

.kb-search-box {
  padding: 0 12px 8px;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.08));
  position: relative;
}
.kb-search-input-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 8px;
  padding: 6px 10px;
  transition: border-color 0.2s;
}
.kb-search-input-wrap:focus-within {
  border-color: var(--accent);
}
.kb-search-icon {
  width: 14px;
  height: 14px;
  color: rgba(246,241,234,0.35);
  flex-shrink: 0;
}
.kb-search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  color: inherit;
  font-size: 13px;
  min-width: 0;
}
.kb-search-input::placeholder {
  color: rgba(246,241,234,0.3);
}
.kb-search-clear {
  background: none;
  border: none;
  color: rgba(246,241,234,0.4);
  font-size: 16px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}
.kb-search-clear:hover {
  color: var(--accent);
}
.kb-search-results {
  position: absolute;
  left: 12px;
  right: 12px;
  top: 100%;
  z-index: 20;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 10px;
  max-height: 360px;
  overflow-y: auto;
  box-shadow: 0 12px 32px rgba(0,0,0,0.5);
}
.kb-search-result-item {
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid rgba(255,255,255,0.05);
  transition: background 0.12s;
}
.kb-search-result-item:last-child {
  border-bottom: none;
}
.kb-search-result-item:hover {
  background: rgba(255,138,76,0.08);
}
.kb-search-result-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
  color: rgba(246,241,234,0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.kb-search-result-snippet {
  font-size: 12px;
  color: rgba(246,241,234,0.5);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.kb-search-result-snippet :deep(mark) {
  background: rgba(255,138,76,0.3);
  color: #ffd166;
  border-radius: 2px;
  padding: 0 2px;
}
.kb-search-result-meta {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-top: 6px;
  font-size: 11px;
  color: rgba(246,241,234,0.35);
}
.kb-search-empty {
  position: absolute;
  left: 12px;
  right: 12px;
  top: 100%;
  z-index: 20;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 10px;
  padding: 20px;
  text-align: center;
  font-size: 13px;
  color: rgba(246,241,234,0.3);
  box-shadow: 0 12px 32px rgba(0,0,0,0.5);
}

.kb-detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px 48px;
  min-width: 0;
}

.kb-doc-view {
  max-width: 720px;
  margin: 0 auto;
}

.kb-doc-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  font-size: 13px;
  color: rgba(246,241,234,0.5);
  flex-wrap: wrap;
}
.kb-doc-category {
  background: rgba(255,209,102,0.15);
  color: #ffd166;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 11px;
}
.kb-doc-date,
.kb-doc-views {
  font-size: 12px;
}
.kb-doc-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 24px;
  line-height: 1.3;
}

.kb-doc-actions {
  display: flex;
  gap: 12px;
  margin-top: 48px;
  padding-top: 20px;
  border-top: 1px solid var(--border, rgba(255,255,255,0.08));
}

.kb-doc-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: rgba(246,241,234,0.3);
  text-align: center;
}
.kb-doc-empty h4 {
  margin: 0 0 8px;
  font-size: 16px;
}
.kb-doc-empty p {
  margin: 0;
  font-size: 13px;
}

.kb-detail-toc {
  width: 200px;
  min-width: 200px;
  border-left: 1px solid var(--border, rgba(255,255,255,0.08));
  overflow-y: auto;
  background: rgba(255,255,255,0.01);
}

@media (max-width: 960px) {
  .kb-detail-toc {
    display: none;
  }
}
@media (max-width: 768px) {
  .kb-detail-layout {
    flex-direction: column;
    height: auto;
  }
  .kb-detail-sidebar {
    width: 100%;
    min-width: unset;
    border-right: none;
    border-bottom: 1px solid var(--border, rgba(255,255,255,0.1));
    max-height: 200px;
  }
  .kb-detail-content {
    padding: 20px;
  }
}
</style>
