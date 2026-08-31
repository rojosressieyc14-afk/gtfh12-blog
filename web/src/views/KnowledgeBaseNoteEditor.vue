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

      <div class="tag-input-wrap">
        <label>
          标签（逗号分隔）
          <input
            v-model.trim="tagsText"
            class="field-input"
            placeholder="Vue, Go, 笔记"
            @input="onTagInput"
            @focus="onTagInput"
            @blur="setTimeout(() => showTagDropdown = false, 150)"
          />
        </label>
        <div v-if="showTagDropdown" class="tag-dropdown">
          <div
            v-for="tag in tagSuggestions"
            :key="tag"
            class="tag-dropdown-item"
            @mousedown.prevent="pickTag(tag)"
          >{{ tag }}</div>
        </div>
      </div>

      <label class="editor-content-label">
        内容
        <RichEditor v-model="form.content" placeholder="开始编写笔记内容..." />
      </label>

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
import { addDocument, updateDocument, listDocuments, listKbTags } from "../api/knowledgeBase";
import { getMetadata } from "../api/meta";
import RichEditor from "../components/RichEditor.vue";

const route = useRoute();
const router = useRouter();
const kbId = route.params.id;
const noteId = route.params.noteId;
const isEdit = Boolean(noteId);

const categories = ref([]);
const saving = ref(false);
const errorMessage = ref("");
const tagsText = ref("");
const existingTags = ref([]);
const tagSuggestions = ref([]);
const showTagDropdown = ref(false);

function onTagInput() {
  const parts = tagsText.value.split(/[,\n]/);
  const current = parts[parts.length - 1]?.trim().toLowerCase() || "";
  if (!current) { showTagDropdown.value = false; return; }
  tagSuggestions.value = existingTags.value.filter(t =>
    t.toLowerCase().includes(current) && !tagsText.value.toLowerCase().includes(t.toLowerCase())
  ).slice(0, 8);
  showTagDropdown.value = tagSuggestions.value.length > 0;
}

function pickTag(tag) {
  const parts = tagsText.value.split(/[,\n]/);
  parts[parts.length - 1] = tag;
  tagsText.value = parts.join(", ") + ", ";
  showTagDropdown.value = false;
}

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
  listKbTags(kbId).then(res => { existingTags.value = res.data.items || []; }).catch(() => {});
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

.editor-content-label {
  display: block;
}

.editor-content-label :deep(.rich-editor) {
  margin-top: 8px;
}

.tag-input-wrap {
  position: relative;
}
.tag-dropdown {
  position: absolute;
  left: 0;
  right: 0;
  top: 100%;
  z-index: 30;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 10px;
  max-height: 200px;
  overflow-y: auto;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}
.tag-dropdown-item {
  padding: 8px 14px;
  font-size: 13px;
  cursor: pointer;
  color: rgba(246,241,234,0.7);
  transition: background 0.12s;
}
.tag-dropdown-item:hover {
  background: rgba(255,138,76,0.1);
  color: #ff8a4c;
}

</style>
