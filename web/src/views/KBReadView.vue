<template>
  <div class="kb-read-layout">
    <aside class="kb-read-sidebar">
      <div class="kb-sidebar-header">
        <router-link class="kb-pub-back-link" :to="`/knowledge-bases`">&larr; 知识库</router-link>
        <h4>{{ kb?.name || '知识库' }}</h4>
      </div>
      <div class="kb-sidebar-tree">
        <DocTree
          :tree="docTree"
          :active-id="currentDoc?.id"
          @select="onSelectDoc"
        />
      </div>
    </aside>

    <main class="kb-read-content">
      <div v-if="loading" class="kb-doc-empty">
        <h4>加载中...</h4>
      </div>
      <div v-else-if="errorMessage" class="kb-doc-empty">
        <h4>出错了</h4>
        <p>{{ errorMessage }}</p>
      </div>
      <div v-else-if="currentDoc" class="kb-read-article">
        <div class="kb-doc-meta">
          <span v-if="currentDoc.category" class="kb-doc-category">{{ currentDoc.category.name }}</span>
          <span class="kb-doc-date">{{ formatDate(currentDoc.createdAt) }}</span>
          <span class="kb-doc-views">{{ currentDoc.viewCount }} 次阅读</span>
          <span v-if="currentDoc.isPublic" class="badge badge-public">公开</span>
          <span v-else class="badge badge-private">私密</span>
        </div>
        <h1 class="kb-doc-title">{{ currentDoc.title || '无标题' }}</h1>
        <div v-if="currentDoc.tags?.length" class="tag-row" style="margin-bottom: 24px;">
          <span v-for="tag in currentDoc.tags" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
        </div>
        <div v-highlight class="markdown-body" v-html="renderedContent"></div>
        <div class="kb-read-nav">
          <router-link class="ghost-btn" :to="`/knowledge-bases/${$route.params.id}`">← 返回知识库</router-link>
          <router-link v-if="userStore.isLoggedIn" class="ghost-btn" :to="`/user-center/knowledge-base/${$route.params.id}/editor/${currentDoc.id}`">编辑此文档</router-link>
        </div>
      </div>
      <div v-else class="kb-doc-empty">
        <h4>选择一篇文档</h4>
        <p>从左栏选择一个文档开始阅读</p>
      </div>
    </main>

    <aside class="kb-read-toc">
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
import { useUserStore } from "../stores/user";
import { getKnowledgeBase, getDocumentTree, listDocuments, listPublicDocuments } from "../api/knowledgeBase";
import DocTree from "../components/DocTree.vue";
import DocToc from "../components/DocToc.vue";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const kb = ref(null);
const docs = ref([]);
const allDocs = ref([]);
const docTree = ref([]);
const currentDoc = ref(null);
const loading = ref(true);
const errorMessage = ref("");
const docToc = ref([]);
const activeHeadingId = ref("");
let tocObserver = null;

const renderedContent = computed(() => {
  if (!currentDoc.value?.content) return "";
  let html = currentDoc.value.content;
  if (currentDoc.value.isMarkdown !== false) {
    html = DOMPurify.sanitize(marked.parse(html));
  }
  // Render wiki-links
  html = html.replace(/\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g, (_, title, alias) => {
    const display = alias || title;
    const doc = allDocs.value.find(d => d.title?.toLowerCase() === title.toLowerCase());
    if (doc) {
      return `<span class="wiki-link" data-doc-id="${doc.id}">${display}</span>`;
    }
    return `<span class="wiki-link wiki-link--broken">${display}</span>`;
  });
  return html;
});

function onSelectDoc(item) {
  const isFolder = item.children && item.children.length > 0;
  if (!isFolder) {
    const doc = allDocs.value.find((d) => Number(d.id) === Number(item.id));
    if (doc) {
      currentDoc.value = doc;
      router.replace({ query: { doc: doc.id } });
    }
  }
}

function scrollToHeading(id) {
  const el = document.getElementById(id);
  if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
}

function setupTocObserver() {
  if (tocObserver) tocObserver.disconnect();
  docToc.value = [];
  activeHeadingId.value = "";
  nextTick(() => {
    const contentEl = document.querySelector(".kb-read-content");
    if (!contentEl) return;
    const headings = contentEl.querySelectorAll(".markdown-body h1, .markdown-body h2, .markdown-body h3");
    headings.forEach((h, i) => {
      if (!h.id) h.id = `doc-heading-${i}`;
      docToc.value.push({ id: h.id, text: h.textContent, level: Number(h.tagName[1]) });
    });
    tocObserver = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            activeHeadingId.value = entry.target.id;
            break;
          }
        }
      },
      { rootMargin: "-80px 0px -60% 0px", threshold: 0 }
    );
    docToc.value.forEach(({ id }) => {
      const el = document.getElementById(id);
      if (el) tocObserver.observe(el);
    });
  });
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

async function load() {
  loading.value = true;
  errorMessage.value = "";
  try {
    const docApi = userStore.isLoggedIn ? listDocuments : listPublicDocuments;
    const [kbRes, docRes, treeRes] = await Promise.all([
      getKnowledgeBase(route.params.id),
      docApi(route.params.id).catch(() => ({ data: { items: [] } })),
      getDocumentTree(route.params.id).catch(() => ({ data: { tree: [] } })),
    ]);
    kb.value = kbRes.data.item;
    docs.value = docRes.data.items || [];
    allDocs.value = docRes.data.items || [];
    docTree.value = treeRes.data.tree || [];

    // Select doc from query or first doc
    const docId = route.query.doc;
    if (docId) {
      const found = allDocs.value.find(d => Number(d.id) === Number(docId));
      if (found) currentDoc.value = found;
    }
    if (!currentDoc.value && allDocs.value.length) {
      currentDoc.value = allDocs.value[0];
    }
  } catch {
    errorMessage.value = "加载失败";
  } finally {
    loading.value = false;
  }
}

watch(currentDoc, () => {
  nextTick(setupTocObserver);
}, { flush: "post" });

watch(() => route.query.doc, (newDocId) => {
  if (newDocId) {
    const found = allDocs.value.find(d => Number(d.id) === Number(newDocId));
    if (found) currentDoc.value = found;
  }
});

onMounted(load);
onUnmounted(() => { if (tocObserver) tocObserver.disconnect(); });
</script>

<style scoped>
.kb-read-layout {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr) 200px;
  gap: 0;
  min-height: calc(100vh - 80px);
}

.kb-read-sidebar {
  border-right: 1px solid var(--border);
  padding: 20px 0;
  overflow-y: auto;
  position: sticky;
  top: 80px;
  height: calc(100vh - 80px);
}

.kb-sidebar-header {
  padding: 0 16px 16px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 12px;
}

.kb-sidebar-header h4 {
  margin: 8px 0 0;
  font-size: 1rem;
}

.kb-pub-back-link {
  font-size: 0.85rem;
  color: var(--text-soft);
  text-decoration: none;
  transition: color 0.2s;
}
.kb-pub-back-link:hover {
  color: var(--accent);
}

.kb-sidebar-tree {
  padding: 0 8px;
}

.kb-read-content {
  padding: 40px clamp(24px, 5vw, 72px);
  max-width: 780px;
  margin: 0 auto;
}

.kb-read-article {
  animation: fadeUp 0.35s ease;
}

.kb-doc-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 16px;
  font-size: 0.85rem;
  color: var(--text-soft);
}

.kb-doc-category {
  background: rgba(255, 138, 76, 0.12);
  color: var(--accent);
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 0.78rem;
}

.kb-doc-title {
  margin: 0 0 24px;
  font-size: clamp(1.6rem, 3vw, 2.2rem);
  line-height: 1.2;
}

.kb-read-nav {
  display: flex;
  gap: 12px;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}

.kb-read-toc {
  border-left: 1px solid var(--border);
  padding: 20px 16px;
  position: sticky;
  top: 80px;
  height: calc(100vh - 80px);
  overflow-y: auto;
}

.kb-doc-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
  color: var(--text-soft);
}
.kb-doc-empty h4 { margin: 0 0 8px; }

:deep(.wiki-link) {
  color: var(--accent);
  cursor: pointer;
  border-bottom: 1px dashed rgba(255, 138, 76, 0.4);
  transition: border-color 0.2s;
}
:deep(.wiki-link:hover) {
  border-bottom-color: var(--accent);
}
:deep(.wiki-link--broken) {
  color: rgba(246, 241, 234, 0.35);
  border-bottom-color: rgba(246, 241, 234, 0.15);
}

@media (max-width: 960px) {
  .kb-read-layout {
    grid-template-columns: 1fr;
  }
  .kb-read-sidebar,
  .kb-read-toc {
    display: none;
  }
}
</style>
