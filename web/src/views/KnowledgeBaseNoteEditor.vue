<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">知识库</p>
        <h3>{{ isEdit ? "编辑笔记" : "新建笔记" }}</h3>
      </div>
      <router-link class="ghost-btn" :to="`/user-center/knowledge-base/${kbId}`">返回</router-link>
    </div>

    <form class="stack-form" @submit.prevent="handleSave">
      <label>
        标题
        <input v-model.trim="form.title" class="field-input" placeholder="笔记标题" />
      </label>

      <div class="note-toolbar">
        <label class="note-toggle">
          <input type="checkbox" v-model="form.isPublic" />
          <span>公开笔记</span>
        </label>
        <label class="note-category">
          分类
          <select v-model="form.categoryId" class="field-select">
            <option :value="null">无分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </label>
      </div>

      <label>
        标签（逗号分隔）
        <input v-model.trim="tagsText" class="field-input" placeholder="Vue, Go, 笔记" />
      </label>

      <label>
        内容（支持 Markdown）
        <textarea v-model="form.content" class="field-area field-area--editor" rows="16" placeholder="使用 Markdown 编写笔记内容..."></textarea>
      </label>

      <div v-if="form.content" class="note-preview">
        <p class="eyebrow">预览</p>
        <div class="markdown-body" v-html="renderedContent"></div>
      </div>

      <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>

      <div class="inline-actions">
        <button class="solid-btn" :disabled="saving || !form.content.trim()">
          {{ saving ? "保存中..." : "保存笔记" }}
        </button>
        <router-link class="ghost-btn" :to="`/user-center/knowledge-base/${kbId}`">取消</router-link>
      </div>
    </form>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { addDocument, updateDocument, listDocuments } from "../api/knowledgeBase";
import { getMetadata } from "../api/meta";

const route = useRoute();
const router = useRouter();
const kbId = route.params.id;
const noteId = route.params.noteId;
const isEdit = Boolean(noteId);

const categories = ref([]);
const saving = ref(false);
const errorMessage = ref("");
const tagsText = ref("");

const form = reactive({
  title: "",
  content: "",
  isPublic: false,
  isMarkdown: true,
  categoryId: null,
});

const renderedContent = computed(() => {
  if (!form.content) return "";
  return DOMPurify.sanitize(marked.parse(form.content));
});

async function loadCategories() {
  try {
    const { data } = await getMetadata();
    categories.value = data.categories || [];
  } catch {}
}

async function loadExisting() {
  if (!isEdit) return;
  try {
    const { data } = await listDocuments(kbId);
    const doc = (data.items || []).find((d) => String(d.id) === noteId);
    if (!doc) {
      errorMessage.value = "笔记不存在";
      return;
    }
    form.title = doc.title || "";
    form.content = doc.content || "";
    form.isPublic = doc.isPublic || false;
    form.isMarkdown = doc.isMarkdown !== false;
    form.categoryId = doc.categoryId || null;
    tagsText.value = (doc.tags || []).map((t) => t.name).join(", ");
  } catch (e) {
    errorMessage.value = "加载笔记失败";
  }
}

async function handleSave() {
  if (!form.content.trim()) return;
  saving.value = true;
  errorMessage.value = "";

  const payload = {
    title: form.title,
    content: form.content,
    isPublic: form.isPublic,
    isMarkdown: form.isMarkdown,
    categoryId: form.categoryId || null,
    tagNames: tagsText.value.split(/[,\n]/).map((s) => s.trim()).filter(Boolean),
  };

  try {
    if (isEdit) {
      await updateDocument(kbId, noteId, payload);
    } else {
      await addDocument(kbId, payload);
    }
    router.push(`/user-center/knowledge-base/${kbId}`);
  } catch (e) {
    errorMessage.value = e?.response?.data?.message || "保存失败";
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  loadCategories();
  loadExisting();
});
</script>

<style scoped>
.note-toolbar {
  display: flex;
  gap: 20px;
  align-items: center;
  flex-wrap: wrap;
}

.note-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 0.95rem;
}

.note-toggle input {
  width: 18px;
  height: 18px;
  accent-color: #f97316;
}

.field-select {
  padding: 10px 14px;
  border-radius: 14px;
  border: 1px solid var(--border, rgba(255,255,255,0.12));
  background: var(--panel, rgba(10,14,19,0.8));
  color: var(--text, #f7f3ea);
  font-size: 0.95rem;
  min-width: 160px;
}

.note-preview {
  margin-top: 16px;
  padding: 20px;
  border-radius: 18px;
  border: 1px solid var(--border, rgba(255,255,255,0.1));
  background: var(--panel, rgba(10,14,19,0.6));
}

.note-preview :deep(.markdown-body) {
  line-height: 1.7;
}

.note-preview :deep(.markdown-body h1),
.note-preview :deep(.markdown-body h2),
.note-preview :deep(.markdown-body h3) {
  margin-top: 1.2em;
  margin-bottom: 0.4em;
}

.note-preview :deep(.markdown-body p) {
  margin: 0.6em 0;
}

.note-preview :deep(.markdown-body code) {
  background: rgba(255,255,255,0.08);
  padding: 2px 6px;
  border-radius: 6px;
  font-size: 0.88em;
}

.note-preview :deep(.markdown-body pre code) {
  display: block;
  padding: 14px;
  overflow-x: auto;
  border-radius: 12px;
  background: rgba(0,0,0,0.3);
}

.note-preview :deep(.markdown-body img) {
  max-width: 100%;
  border-radius: 12px;
}
</style>
