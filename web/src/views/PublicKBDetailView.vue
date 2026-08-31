<template>
  <div class="kb-detail-layout">
    <aside class="kb-detail-sidebar">
      <div class="kb-sidebar-header">
        <router-link class="kb-pub-back-link" to="/knowledge-bases">&larr; 知识库</router-link>
        <h4>{{ kb?.name || '知识库' }}</h4>
      </div>
      <div v-if="isLoggedIn" class="kb-search-box">
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
    </aside>

    <main class="kb-detail-content" @click="onDocContentClick">
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
        </div>
        <h1 class="kb-doc-title">{{ currentDoc.title || '无标题' }}</h1>
        <div v-highlight class="markdown-body" v-html="renderedContent"></div>
        <div v-if="isLoggedIn" class="kb-doc-actions">
          <button class="ghost-btn" @click="$router.push({ name: 'uc-kb-graph', params: { id: $route.params.id } })">图谱</button>
        </div>
      </div>
      <div v-else-if="!docs.length" class="kb-doc-empty">
        <h4>该知识库暂无公开文档</h4>
        <router-link class="ghost-btn" to="/knowledge-bases">返回列表</router-link>
      </div>
      <div v-else class="kb-doc-empty">
        <h4>选择一篇文档</h4>
        <p>从左栏选择一个文档开始阅读</p>
      </div>
    </main>

    <aside v-if="currentDoc && isLoggedIn" class="kb-detail-toc">
      <DocToc
        :items="docToc"
        :active-id="activeHeadingId"
        @jump="scrollToHeading"
      />
      <div v-if="backlinks.length > 0" class="kb-backlinks">
        <h5 class="kb-backlinks-title">反向链接</h5>
        <div
          v-for="bl in backlinks"
          :key="bl.docId"
          class="kb-backlink-item"
          @click="onSelectDoc({ id: bl.docId })"
        >
          <div class="kb-backlink-title">{{ bl.title || '无标题' }}</div>
          <div class="kb-backlink-snippet">{{ bl.snippet }}</div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { useUserStore } from "../stores/user";
import {
  listPublicDocuments, getKnowledgeBase, getDocumentTree,
  searchDocuments, getBacklinks
} from "../api/knowledgeBase";
import DocTree from "../components/DocTree.vue";
import DocToc from "../components/DocToc.vue";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const isLoggedIn = computed(() => userStore.isLoggedIn);

const kb = ref(null);
const docTree = ref([]);
const currentDoc = ref(null);
const activeHeadingId = ref(null);
const allDocs = ref([]);
const docs = ref([]);
const loading = ref(true);
const errorMessage = ref("");
const searchKeyword = ref("");
const searchResults = ref([]);
const searchLoading = ref(false);
let searchTimer = null;
const backlinks = ref([]);

function buildHeadingId(text) {
  return text.toLowerCase().replace(/[^\w\u4e00-\u9fff]+/g, "-").replace(/(^-|-$)/g, "");
}

const renderedContent = computed(() => {
  if (!currentDoc.value?.content) return "";
  let md = currentDoc.value.content;
  md = md.replace(/\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g, (_, title, display) => {
    const text = display || title;
    const target = allDocs.value.find(d => d.title && d.title.toLowerCase() === title.toLowerCase());
    if (target) {
      return `<a class="wiki-link" data-doc-id="${target.id}" href="#">${text}</a>`;
    }
    return `<a class="wiki-link wiki-link-broken" href="#">${text}</a>`;
  });
  const renderer = new marked.Renderer();
  renderer.heading = (token) => {
    const text = token.text || "";
    const id = buildHeadingId(text);
    return `<h${token.depth} id="${id}">${text}</h${token.depth}>`;
  };
  return DOMPurify.sanitize(marked.parse(md, { renderer }));
});

const docToc = computed(() => {
  if (!currentDoc.value?.content) return [];
  const headingRe = /^(#{1,4})\s+(.+)$/gm;
  const items = [];
  let match;
  while ((match = headingRe.exec(currentDoc.value.content)) !== null) {
    items.push({ id: buildHeadingId(match[2]), text: match[2], level: match[1].length });
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
        if (e.isIntersecting) { activeHeadingId.value = e.target.id; break; }
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
    router.push({ name: 'kb-read', params: { id: route.params.id }, query: { doc: item.id } });
  }
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

function onDocContentClick(e) {
  const link = e.target.closest('.wiki-link');
  if (!link) return;
  e.preventDefault();
  const docId = Number(link.dataset.docId);
  if (!docId) return;
  router.push({ name: 'kb-read', params: { id: route.params.id }, query: { doc: docId } });
}

function onSearchInput() {
  clearTimeout(searchTimer);
  if (!searchKeyword.value.trim()) { searchResults.value = []; return; }
  searchTimer = setTimeout(performSearch, 300);
}

async function performSearch() {
  searchLoading.value = true;
  try {
    const res = await searchDocuments(route.params.id, searchKeyword.value.trim());
    searchResults.value = res.data.items || [];
  } catch { searchResults.value = []; } finally { searchLoading.value = false; }
}

function clearSearch() { searchKeyword.value = ""; searchResults.value = []; }

function selectSearchResult(r) {
  router.push({ name: 'kb-read', params: { id: route.params.id }, query: { doc: r.id } });
  clearSearch();
}

function highlightSnippet(snippet) {
  if (!searchKeyword.value.trim()) return snippet;
  const kw = searchKeyword.value.trim().replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return snippet.replace(new RegExp(`(${kw})`, "gi"), '<mark>$1</mark>');
}

async function loadBacklinks() {
  if (!currentDoc.value) { backlinks.value = []; return; }
  try {
    const res = await getBacklinks(route.params.id, currentDoc.value.id);
    backlinks.value = res.data.items || [];
  } catch { backlinks.value = []; }
}

watch(currentDoc, () => {
  nextTick(setupTocObserver);
  if (isLoggedIn.value) loadBacklinks();
}, { flush: "post" });

onMounted(async () => {
  try {
    const [kbRes, docRes, treeRes] = await Promise.all([
      getKnowledgeBase(route.params.id),
      listPublicDocuments(route.params.id),
      getDocumentTree(route.params.id).catch(() => ({ data: { tree: [] } })),
    ]);
    kb.value = kbRes.data.item;
    docs.value = docRes.data.items || [];
    allDocs.value = docRes.data.items || [];
    docTree.value = treeRes.data.tree || [];
  } catch {
    errorMessage.value = "加载失败";
  } finally {
    loading.value = false;
  }
});

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
  flex-direction: column;
  gap: 4px;
  padding: 16px;
  border-bottom: 1px solid var(--border, rgba(255,255,255,0.08));
}
.kb-sidebar-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}
.kb-pub-back-link {
  font-size: 12px;
  color: rgba(246,241,234,0.5);
  text-decoration: none;
}
.kb-pub-back-link:hover {
  color: var(--accent);
}
.kb-sidebar-tree {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
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
.kb-doc-date, .kb-doc-views { font-size: 12px; }
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
.kb-doc-empty h4 { margin: 0 0 8px; font-size: 16px; }
.kb-doc-empty p { margin: 0; font-size: 13px; }
.kb-detail-toc {
  width: 200px;
  min-width: 200px;
  border-left: 1px solid var(--border, rgba(255,255,255,0.08));
  overflow-y: auto;
  background: rgba(255,255,255,0.01);
  display: flex;
  flex-direction: column;
}
.kb-backlinks {
  padding: 12px;
  border-top: 1px solid var(--border, rgba(255,255,255,0.08));
}
.kb-backlinks-title {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: rgba(246,241,234,0.4);
}
.kb-backlink-item {
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.12s;
  margin-bottom: 4px;
}
.kb-backlink-item:hover { background: rgba(255,138,76,0.08); }
.kb-backlink-title {
  font-size: 12px;
  font-weight: 600;
  color: rgba(246,241,234,0.85);
  margin-bottom: 2px;
}
.kb-backlink-snippet {
  font-size: 11px;
  color: rgba(246,241,234,0.4);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

:deep(.wiki-link) {
  color: var(--accent);
  border-bottom: 1px dashed var(--accent);
  cursor: pointer;
  text-decoration: none;
  transition: color 0.15s, border-color 0.15s;
}
:deep(.wiki-link:hover) { color: #ffd166; border-color: #ffd166; }
:deep(.wiki-link-broken) {
  color: rgba(246,241,234,0.35);
  border-color: rgba(246,241,234,0.2);
  cursor: default;
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
.kb-search-input-wrap:focus-within { border-color: var(--accent); }
.kb-search-icon { width: 14px; height: 14px; color: rgba(246,241,234,0.35); flex-shrink: 0; }
.kb-search-input { flex: 1; background: none; border: none; outline: none; color: inherit; font-size: 13px; min-width: 0; }
.kb-search-input::placeholder { color: rgba(246,241,234,0.3); }
.kb-search-clear { background: none; border: none; color: rgba(246,241,234,0.4); font-size: 16px; cursor: pointer; padding: 0; line-height: 1; }
.kb-search-clear:hover { color: var(--accent); }
.kb-search-results {
  position: absolute; left: 12px; right: 12px; top: 100%; z-index: 20;
  background: #1a1a1a; border: 1px solid rgba(255,255,255,0.1); border-radius: 10px;
  max-height: 360px; overflow-y: auto; box-shadow: 0 12px 32px rgba(0,0,0,0.5);
}
.kb-search-result-item { padding: 10px 14px; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.05); transition: background 0.12s; }
.kb-search-result-item:last-child { border-bottom: none; }
.kb-search-result-item:hover { background: rgba(255,138,76,0.08); }
.kb-search-result-title { font-size: 13px; font-weight: 600; margin-bottom: 4px; color: rgba(246,241,234,0.9); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.kb-search-result-snippet { font-size: 12px; color: rgba(246,241,234,0.5); line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.kb-search-result-snippet :deep(mark) { background: rgba(255,138,76,0.3); color: #ffd166; border-radius: 2px; padding: 0 2px; }
.kb-search-result-meta { display: flex; gap: 8px; align-items: center; margin-top: 6px; font-size: 11px; color: rgba(246,241,234,0.35); }
.kb-search-empty {
  position: absolute; left: 12px; right: 12px; top: 100%; z-index: 20;
  background: #1a1a1a; border: 1px solid rgba(255,255,255,0.1); border-radius: 10px;
  padding: 20px; text-align: center; font-size: 13px; color: rgba(246,241,234,0.3);
  box-shadow: 0 12px 32px rgba(0,0,0,0.5);
}

@media (max-width: 960px) {
  .kb-detail-toc { display: none; }
}
@media (max-width: 768px) {
  .kb-detail-layout { flex-direction: column; height: auto; }
  .kb-detail-sidebar { width: 100%; min-width: unset; border-right: none; border-bottom: 1px solid var(--border, rgba(255,255,255,0.1)); max-height: 200px; }
  .kb-detail-content { padding: 20px; }
}
</style>
