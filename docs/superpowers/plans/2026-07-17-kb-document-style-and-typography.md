# Knowledge Base Document Style + Typography Upgrade Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform the flat knowledge base into a Feishu-style three-column document interface and establish a unified .markdown-body typography system across all article/note views.

**Architecture:** 
- Phase 1: CSS-only typography overhaul (global .markdown-body rules, remove component-scoped duplicates)
- Phase 2: Backend API for document tree (parent_id, sort_order, tree/move endpoints)
- Phase 3: Vue components for document tree (DocTree), table of contents (DocToc), and three-column KB detail layout

**Tech Stack:** Vue 3 (frontend), Go + Gin + GORM (backend), CSS custom properties (theming)

---

### Task 1: Rewrite global .markdown-body CSS

**Files:**
- Modify: `web/src/assets/main.css` (around line 959)

- [ ] **Step 1: Read current .markdown-body section**

Read lines 959-972 from `web/src/assets/main.css` to confirm the current state.

- [ ] **Step 2: Replace .markdown-body CSS block**

Replace the existing minimal `.markdown-body` rules (lines 959-972) with the full typography system:

```css
.markdown-body {
  max-width: 720px;
  margin: 0 auto;
  line-height: 1.8;
  font-size: 16px;
  color: rgba(250, 246, 239, 0.92);
}
.markdown-body:lang(en) {
  line-height: 1.6;
}
.markdown-body h1 {
  font-size: 2em;
  font-weight: 700;
  margin: 1.5em 0 0.5em;
  line-height: 1.3;
}
.markdown-body h2 {
  font-size: 1.6em;
  font-weight: 600;
  margin: 1.4em 0 0.4em;
  line-height: 1.35;
}
.markdown-body h3 {
  font-size: 1.3em;
  font-weight: 600;
  margin: 1.3em 0 0.3em;
  line-height: 1.4;
}
.markdown-body h4 {
  font-size: 1.1em;
  font-weight: 600;
  margin: 1.2em 0 0.2em;
}
.markdown-body p {
  margin-bottom: 1.2em;
}
.markdown-body p,
.markdown-body li {
  color: var(--text-soft);
}
.markdown-body pre {
  background: rgba(255,255,255,0.06);
  border-radius: 12px;
  padding: 20px;
  overflow-x: auto;
  font-size: 14px;
  line-height: 1.6;
  border: 1px solid rgba(255,255,255,0.08);
}
.markdown-body code {
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  background: rgba(255,255,255,0.08);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.9em;
}
.markdown-body pre code {
  background: none;
  padding: 0;
  border-radius: 0;
}
.markdown-body blockquote {
  border-left: 4px solid var(--accent);
  background: rgba(255,138,76,0.08);
  padding: 12px 20px;
  margin: 1.2em 0;
  border-radius: 0 8px 8px 0;
  color: var(--text-soft);
}
.markdown-body img {
  border-radius: 12px;
  max-width: 100%;
  margin: 1.5em auto;
  display: block;
  box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}
.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5em 0;
  border-radius: 8px;
  overflow: hidden;
}
.markdown-body th,
.markdown-body td {
  padding: 10px 16px;
  border-bottom: 1px solid var(--border);
  text-align: left;
}
.markdown-body th {
  background: rgba(255,255,255,0.06);
  font-weight: 600;
}
.markdown-body tr:nth-child(even) {
  background: rgba(255,255,255,0.03);
}
.markdown-body ul,
.markdown-body ol {
  padding-left: 1.5em;
  margin-bottom: 1.2em;
}
.markdown-body li {
  margin-bottom: 0.3em;
}
.markdown-body hr {
  border: none;
  height: 1px;
  background: var(--border);
  margin: 2em 0;
}
.markdown-body a {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 2px;
}
.markdown-body a:hover {
  color: var(--accent-soft);
}
```

- [ ] **Step 3: Verify the edit**

Run: `Select-String -Path "web/src/assets/main.css" -Pattern "\.markdown-body" | Measure-Object | Select-Object -ExpandProperty Count`
Expected: > 40 lines of markdown-body rules (was 4, now comprehensive)

---

### Task 2: Remove scoped markdown-body styles from KbNoteView.vue

**Files:**
- Modify: `web/src/views/KbNoteView.vue`

- [ ] **Step 1: Read KbNoteView.vue**

Read `web/src/views/KbNoteView.vue` to confirm current scoped styles.

- [ ] **Step 2: Delete the scoped .markdown-body block**

Remove lines 94-135 (the entire `.markdown-body` scoped style block within `<style scoped>`) from `web/src/views/KbNoteView.vue`.

The component already uses `class="markdown-body"` on the rendered content div (line 19), so the global CSS will take effect automatically.

- [ ] **Step 3: Verify the edit**

Run: `Select-String -Path "web/src/views/KbNoteView.vue" -Pattern "\.markdown-body"`
Expected: 0 matches (all removed)

---

### Task 3: Remove scoped markdown-body deep styles from KnowledgeBaseNoteEditor.vue

**Files:**
- Modify: `web/src/views/KnowledgeBaseNoteEditor.vue`

- [ ] **Step 1: Read KnowledgeBaseNoteEditor.vue**

Read `web/src/views/KnowledgeBaseNoteEditor.vue` to confirm current scoped styles.

- [ ] **Step 2: Delete the scoped deep .markdown-body block**

Remove lines 191-224 (the entire `.note-preview :deep(.markdown-body)` block) from `web/src/views/KnowledgeBaseNoteEditor.vue`.

- [ ] **Step 3: Verify the edit**

Run: `Select-String -Path "web/src/views/KnowledgeBaseNoteEditor.vue" -Pattern "\.markdown-body"`
Expected: 0 matches

---

### Task 4: Update DetailView.vue article reading width

**Files:**
- Modify: `web/src/views/DetailView.vue`
- Modify: `web/src/assets/main.css`

- [ ] **Step 1: Read DetailView.vue template and main.css detail styles**

Read lines 49-51 of DetailView.vue to see the `.article-detail-body` usage on line 50.
Read lines 1002-1027 of main.css for current `.article-detail-body` styles.

- [ ] **Step 2: Remove width restriction from .article-detail-body in DetailView**

In `web/src/views/DetailView.vue`, line 50, the article body already uses `.markdown-body.article-detail-body`. The global `.markdown-body` now has `max-width: 720px; margin: 0 auto`. 

In `web/src/assets/main.css`, the `.article-detail-body` at line 1002 already has `padding`, `border-radius`, `border`, and `background`. No change needed to it - the `.markdown-body` width restriction on the same element handles centering.

No code change needed in DetailView.vue - the article content will automatically use the new global `.markdown-body` styles.

- [ ] **Step 3: Ensure .article-detail-body doesn't conflict with .markdown-body width**

The `max-width: 720px` on `.markdown-body` combined with padding on `.article-detail-body` may make content too narrow. Fix by adding:

```css
/* main.css - add after .markdown-body block */
.article-detail-body.markdown-body {
  max-width: 760px;
}
```

This gives a bit more width to accommodate the article-detail-body's internal padding.

---

### Task 5: Add parent_id and sort_order to KnowledgeDocument model

**Files:**
- Modify: `server/internal/model/knowledge_base.go`

- [ ] **Step 1: Read current model**

Read `server/internal/model/knowledge_base.go` to confirm the KnowledgeDocument struct.

- [ ] **Step 2: Add parent_id and sort_order fields**

Add these fields to the KnowledgeDocument struct after `UserID`:
```go
ParentID  *uint  `gorm:"index" json:"parentId"`
SortOrder int    `gorm:"not null;default:0" json:"sortOrder"`
```

- [ ] **Step 3: Verify the edit**

Run: `Select-String -Path "server/internal/model/knowledge_base.go" -Pattern "ParentID|SortOrder"`
Expected: 2 matches

---

### Task 6: Add tree and move service methods

**Files:**
- Modify: `server/internal/service/knowledge_base_service.go`

- [ ] **Step 1: Read the current service file**

Read the full `server/internal/service/knowledge_base_service.go` to understand existing methods.

- [ ] **Step 2: Add DocTreeItem type**

Add a new type for the tree structure after the `Source` struct (around line 35):
```go
type DocTreeItem struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	ParentID  *uint          `json:"parentId"`
	SortOrder int            `json:"sortOrder"`
	Children  []*DocTreeItem `json:"children"`
}
```

- [ ] **Step 3: Add GetDocumentTree method**

Add after the `ListDocuments` method:
```go
func (s *KnowledgeBaseService) GetDocumentTree(kbID, userID uint) ([]*DocTreeItem, error) {
	if _, err := s.GetByID(kbID, userID); err != nil {
		return nil, err
	}
	var docs []model.KnowledgeDocument
	if err := s.db.Where("knowledge_base_id = ?", kbID).
		Order("sort_order asc, created_at desc").
		Find(&docs).Error; err != nil {
		return nil, err
	}
	docMap := make(map[uint]*DocTreeItem)
	var roots []*DocTreeItem
	for i := range docs {
		item := &DocTreeItem{
			ID:        docs[i].ID,
			Title:     docs[i].Title,
			ParentID:  docs[i].ParentID,
			SortOrder: docs[i].SortOrder,
			Children:  []*DocTreeItem{},
		}
		docMap[docs[i].ID] = item
		if docs[i].ParentID == nil {
			roots = append(roots, item)
		} else {
			if parent, ok := docMap[*docs[i].ParentID]; ok {
				parent.Children = append(parent.Children, item)
			} else {
				roots = append(roots, item)
			}
		}
	}
	return roots, nil
}
```

- [ ] **Step 4: Add MoveDocument method**

Add after `GetDocumentTree`:
```go
type MoveDocumentOpts struct {
	ParentID  *uint `json:"parentId"`
	SortOrder int   `json:"sortOrder"`
}

func (s *KnowledgeBaseService) MoveDocument(kbID, docID, userID uint, opts MoveDocumentOpts) error {
	if _, err := s.GetByID(kbID, userID); err != nil {
		return err
	}
	var doc model.KnowledgeDocument
	if err := s.db.Where("id = ? AND knowledge_base_id = ?", docID, kbID).First(&doc).Error; err != nil {
		return ErrDocNotFound
	}
	return s.db.Model(&doc).Updates(map[string]interface{}{
		"parent_id":  opts.ParentID,
		"sort_order": opts.SortOrder,
	}).Error
}
```

- [ ] **Step 5: Verify the file compiles**

Run: `cd server; go build ./...`
Expected: No errors

---

### Task 7: Add handler methods for tree and move

**Files:**
- Modify: `server/internal/handler/knowledge_base_handler.go`

- [ ] **Step 1: Read the current handler**

Read `server/internal/handler/knowledge_base_handler.go` to confirm existing methods.

- [ ] **Step 2: Add GetDocumentTree handler**

Add after the `Query` method:
```go
func (h *KnowledgeBaseHandler) GetDocumentTree(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}
	tree, err := h.svc.GetDocumentTree(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}
```

- [ ] **Step 3: Add MoveDocument handler**

Add after `GetDocumentTree`:
```go
func (h *KnowledgeBaseHandler) MoveDocument(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, _ := strconv.Atoi(c.Param("id"))
	docID, _ := strconv.Atoi(c.Param("docId"))
	if kbID <= 0 || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}
	var opts service.MoveDocumentOpts
	if err := safeBindJSON(c, &opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}
	if err := h.svc.MoveDocument(uint(kbID), uint(docID), authUser.ID, opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文档已移动"})
}
```

- [ ] **Step 4: Verify the file**

Run: `cd server; go build ./...`
Expected: No errors

---

### Task 8: Register new routes

**Files:**
- Modify: `server/internal/router/router.go`

- [ ] **Step 1: Read router.go**

Read `server/internal/router/router.go` lines 119-127 to confirm existing KB routes.

- [ ] **Step 2: Add two new routes**

Add after line 125 (`kbHandler.ListDocuments`):
```go
api.GET("/knowledge-bases/:id/tree", middleware.RequireAuth(), kbHandler.GetDocumentTree)
api.PUT("/knowledge-bases/:id/documents/:docId/move", middleware.RequireAuth(), kbHandler.MoveDocument)
```

- [ ] **Step 3: Verify build**

Run: `cd server; go build ./...`
Expected: No errors

---

### Task 9: Create DocTree.vue component

**Files:**
- Create: `web/src/components/DocTree.vue`

- [ ] **Step 1: Create DocTree.vue**

```vue
<template>
  <div class="doc-tree">
    <div
      v-for="item in tree"
      :key="item.id"
      class="doc-tree-node"
    >
      <div
        class="doc-tree-item"
        :class="{ 'doc-tree-item--active': activeId === item.id, 'doc-tree-item--folder': item.children?.length }"
        @click="$emit('select', item)"
      >
        <span class="doc-tree-icon">{{ item.children?.length ? '📁' : '📄' }}</span>
        <span class="doc-tree-title">{{ item.title || '无标题' }}</span>
      </div>
      <div v-if="item.children?.length" class="doc-tree-children">
        <DocTree
          :tree="item.children"
          :active-id="activeId"
          @select="(n) => $emit('select', n)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  tree: { type: Array, default: () => [] },
  activeId: { type: [Number, String], default: null },
});
defineEmits(['select']);
</script>

<style scoped>
.doc-tree {
  font-size: 13px;
}
.doc-tree-node {
  margin: 0;
}
.doc-tree-children {
  padding-left: 16px;
}
.doc-tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  color: rgba(246,241,234,0.7);
  transition: background 0.15s, color 0.15s;
  margin-bottom: 2px;
}
.doc-tree-item:hover {
  background: rgba(255,255,255,0.06);
  color: rgba(246,241,234,0.9);
}
.doc-tree-item--active {
  background: rgba(255,138,76,0.15);
  color: #ff8a4c;
}
.doc-tree-item--folder {
  font-weight: 500;
}
.doc-tree-icon {
  font-size: 14px;
  line-height: 1;
}
.doc-tree-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
```

- [ ] **Step 2: Verify file exists**

Run: `Test-Path "web/src/components/DocTree.vue"`
Expected: True

---

### Task 10: Create DocToc.vue component

**Files:**
- Create: `web/src/components/DocToc.vue`

- [ ] **Step 1: Create DocToc.vue**

```vue
<template>
  <div class="doc-toc">
    <p class="doc-toc-label">目录</p>
    <div v-if="items.length" class="doc-toc-list">
      <button
        v-for="item in items"
        :key="item.id"
        class="doc-toc-link"
        :class="[`doc-toc-link--h${item.level}`, { 'doc-toc-link--active': activeId === item.id }]"
        @click="$emit('jump', item.id)"
      >
        {{ item.text }}
      </button>
    </div>
    <p v-else class="doc-toc-empty">暂无目录</p>
  </div>
</template>

<script setup>
defineProps({
  items: { type: Array, default: () => [] },
  activeId: { type: [String, Number], default: null },
});
defineEmits(['jump']);
</script>

<style scoped>
.doc-toc {
  padding: 16px 12px;
}
.doc-toc-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: rgba(246,241,234,0.4);
  margin-bottom: 12px;
}
.doc-toc-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.doc-toc-link {
  display: block;
  width: 100%;
  text-align: left;
  border: none;
  background: none;
  padding: 5px 8px;
  border-radius: 4px;
  font-size: 13px;
  color: rgba(246,241,234,0.5);
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
  line-height: 1.4;
}
.doc-toc-link:hover {
  color: rgba(246,241,234,0.8);
  background: rgba(255,255,255,0.04);
}
.doc-toc-link--active {
  color: #ff8a4c;
  background: rgba(255,138,76,0.1);
}
.doc-toc-link--h3 {
  padding-left: 20px;
  font-size: 12px;
}
.doc-toc-link--h4 {
  padding-left: 28px;
  font-size: 11px;
}
.doc-toc-empty {
  font-size: 12px;
  color: rgba(246,241,234,0.3);
}
</style>
```

- [ ] **Step 2: Verify file exists**

Run: `Test-Path "web/src/components/DocToc.vue"`
Expected: True

---

### Task 11: Refactor KnowledgeBaseDetail.vue to three-column layout

**Files:**
- Modify: `web/src/views/KnowledgeBaseDetail.vue`
- Modify: `web/src/api/knowledgeBase.js`

- [ ] **Step 1: Add API client methods**

Add to `web/src/api/knowledgeBase.js`:
```js
export const getDocumentTree = (kbId) => client.get(`/knowledge-bases/${kbId}/tree`);
export const moveDocument = (kbId, docId, payload) => client.put(`/knowledge-bases/${kbId}/documents/${docId}/move`, payload);
```

- [ ] **Step 2: Rewrite KnowledgeBaseDetail.vue template**

Replace the entire template with a three-column layout:
```vue
<template>
  <div class="kb-detail-layout">
    <!-- Left: Document Tree -->
    <aside class="kb-detail-sidebar">
      <div class="kb-sidebar-header">
        <h4>{{ kb?.name || '知识库' }}</h4>
        <button class="kb-sidebar-add" @click="showNewFolder = true">+</button>
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

    <!-- Middle: Document Content -->
    <main class="kb-detail-content">
      <div v-if="currentDoc" class="kb-doc-view">
        <div class="kb-doc-meta">
          <span v-if="currentDoc.category" class="kb-doc-category">{{ currentDoc.category.name }}</span>
          <span class="kb-doc-date">{{ formatDate(currentDoc.createdAt) }}</span>
          <span class="kb-doc-views">{{ currentDoc.viewCount }} 次阅读</span>
          <span v-if="currentDoc.isPublic" class="badge badge-public">公开</span>
          <span v-else class="badge badge-private">私密</span>
        </div>
        <h1 class="kb-doc-title">{{ currentDoc.title || '无标题' }}</h1>
        <div class="markdown-body" v-html="renderedContent"></div>
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

    <!-- Right: Table of Contents -->
    <aside v-if="currentDoc" class="kb-detail-toc">
      <DocToc
        :items="docToc"
        :active-id="activeHeadingId"
        @jump="scrollToHeading"
      />
    </aside>
  </div>
</template>
```

- [ ] **Step 3: Rewrite the script section**

Replace the `<script setup>` block:
```vue
<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { getKnowledgeBase, getDocumentTree, listDocuments, deleteDocument } from "../api/knowledgeBase";
import DocTree from "../components/DocTree.vue";
import DocToc from "../components/DocToc.vue";

const route = useRoute();
const router = useRouter();
const kb = ref(null);
const docTree = ref([]);
const currentDoc = ref(null);
const activeHeadingId = ref(null);
const allDocs = ref([]);

const renderedContent = computed(() => {
  if (!currentDoc.value?.content) return "";
  return DOMPurify.sanitize(marked.parse(currentDoc.value.content));
});

const docToc = computed(() => {
  if (!currentDoc.value?.content) return [];
  const headingRe = /^(#{1,4})\s+(.+)$/gm;
  const items = [];
  let match;
  while ((match = headingRe.exec(currentDoc.value.content)) !== null) {
    const level = match[1].length;
    const text = match[2];
    const id = text.toLowerCase().replace(/[^\w\u4e00-\u9fff]+/g, "-").replace(/(^-|-$)/g, "");
    items.push({ id, text, level });
  }
  return items;
});

function scrollToHeading(id) {
  const el = document.getElementById(id);
  if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
}

function onSelectDoc(item) {
  if (!item.children) {
    const doc = allDocs.value.find((d) => d.id === item.id);
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
    if (currentDoc.value?.id === doc.id) currentDoc.value = null;
    await load();
  } catch (e) {
    alert("删除失败：" + (e?.response?.data?.message || e.message));
  }
}

async function load() {
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
    alert("加载失败：" + (e?.response?.data?.message || e.message));
  }
}

onMounted(load);
</script>
```

- [ ] **Step 4: Replace the style section**

Replace the entire `<style scoped>` block:
```css
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
```

---

### Task 12: Add JetBrains Mono font import

**Files:**
- Modify: `web/index.html`

- [ ] **Step 1: Read index.html**

Read `web/index.html` to see existing Google Fonts imports.

- [ ] **Step 2: Add JetBrains Mono font**

In the `<head>` section, add JetBrains Mono after the existing font imports:
```html
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
```

---

### Verification

After all tasks are complete, run these checks:

1. **Frontend build**: `cd web; npm run build`
2. **Backend build**: `cd server; go build ./...`
3. **CSS check**: Verify `.markdown-body` rules are comprehensive in `web/src/assets/main.css`
4. **Scoped style check**: Verify no `.markdown-body` styles remain in KbNoteView.vue or KnowledgeBaseNoteEditor.vue
5. **Route check**: Verify `GET /api/knowledge-bases/:id/tree` and `PUT /api/knowledge-bases/:id/documents/:docId/move` are registered
6. **Component check**: Verify DocTree.vue and DocToc.vue exist
