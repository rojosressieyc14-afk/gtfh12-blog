<template>
  <div class="rich-editor">
    <div class="rich-editor__toolbar" v-if="editor">
      <div class="toolbar-group">
        <button
          type="button"
          :class="{ active: editor.isActive('bold') }"
          title="加粗 (Ctrl+B)"
          @click="editor.chain().focus().toggleBold().run()"
        >
          <strong>B</strong>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('italic') }"
          title="斜体 (Ctrl+I)"
          @click="editor.chain().focus().toggleItalic().run()"
        >
          <em>I</em>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('underline') }"
          title="下划线 (Ctrl+U)"
          @click="editor.chain().focus().toggleUnderline().run()"
        >
          <u>U</u>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('strike') }"
          title="删除线"
          @click="editor.chain().focus().toggleStrike().run()"
        >
          <s>S</s>
        </button>
      </div>

      <span class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          type="button"
          :class="{ active: editor.isActive('heading', { level: 2 }) }"
          title="标题 2"
          @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
        >
          H2
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('heading', { level: 3 }) }"
          title="标题 3"
          @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
        >
          H3
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('heading', { level: 4 }) }"
          title="标题 4"
          @click="editor.chain().focus().toggleHeading({ level: 4 }).run()"
        >
          H4
        </button>
      </div>

      <span class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          type="button"
          :class="{ active: editor.isActive('bulletList') }"
          title="无序列表"
          @click="editor.chain().focus().toggleBulletList().run()"
        >
          <span class="toolbar-icon">•≡</span>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('orderedList') }"
          title="有序列表"
          @click="editor.chain().focus().toggleOrderedList().run()"
        >
          <span class="toolbar-icon">1.</span>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('taskList') }"
          title="任务列表"
          @click="editor.chain().focus().toggleTaskList().run()"
        >
          <span class="toolbar-icon">☑</span>
        </button>
      </div>

      <span class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          type="button"
          :class="{ active: editor.isActive('codeBlock') }"
          title="代码块"
          @click="editor.chain().focus().toggleCodeBlock().run()"
        >
          <span class="toolbar-icon">&lt;/&gt;</span>
        </button>
        <button
          type="button"
          :class="{ active: editor.isActive('blockquote') }"
          title="引用"
          @click="editor.chain().focus().toggleBlockquote().run()"
        >
          <span class="toolbar-icon">"</span>
        </button>
        <button
          type="button"
          title="分割线"
          @click="editor.chain().focus().setHorizontalRule().run()"
        >
          <span class="toolbar-icon">—</span>
        </button>
      </div>

      <span class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          type="button"
          title="插入链接"
          @click="setLink"
        >
          <span class="toolbar-icon">🔗</span>
        </button>
        <button
          type="button"
          title="插入图片"
          @click="triggerImageUpload"
        >
          <span class="toolbar-icon">🖼</span>
        </button>
        <button
          type="button"
          title="插入表格"
          @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()"
        >
          <span class="toolbar-icon">▦</span>
        </button>
      </div>

      <input
        ref="imageInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleImageUpload"
      />
    </div>

    <EditorContent :editor="editor" class="rich-editor__content" />
  </div>
</template>

<script setup>
import { watch, onBeforeUnmount, ref } from "vue";
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import { Markdown } from "@tiptap/markdown";
import Image from "@tiptap/extension-image";
import { TableKit } from "@tiptap/extension-table";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import { common, createLowlight } from "lowlight";
import Placeholder from "@tiptap/extension-placeholder";
import Link from "@tiptap/extension-link";
import Underline from "@tiptap/extension-underline";
import TaskList from "@tiptap/extension-task-list";
import TaskItem from "@tiptap/extension-task-item";
import { uploadImage } from "../api/upload";
import { compressImage } from "../utils/compressImage";

const props = defineProps({
  modelValue: { type: String, default: "" },
  placeholder: { type: String, default: "开始写作..." },
});

const emit = defineEmits(["update:modelValue"]);

const imageInput = ref(null);
const lowlight = createLowlight(common);

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      codeBlock: false,
    }),
    Markdown.configure({
      transformPastedText: true,
      transformCopiedText: true,
    }),
    Image.configure({ inline: true, allowBase64: true }),
    TableKit.configure({
      table: { resizable: true },
    }),
    CodeBlockLowlight.configure({
      lowlight,
    }),
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: { class: "editor-link" },
    }),
    Underline,
    TaskList,
    TaskItem.configure({ nested: true }),
  ],
  onUpdate: ({ editor: e }) => {
    emit("update:modelValue", e.getMarkdown());
  },
  editorProps: {
    handlePaste: (view, event) => {
      const files = event.clipboardData?.files;
      if (!files?.length) return false;
      for (const file of files) {
        if (file.type.startsWith("image/")) {
          uploadAndInsertImage(file);
          return true;
        }
      }
      return false;
    },
    handleDrop: (view, event) => {
      const files = event.dataTransfer?.files;
      if (!files?.length) return false;
      for (const file of files) {
        if (file.type.startsWith("image/")) {
          uploadAndInsertImage(file);
          return true;
        }
      }
      return false;
    },
  },
});

watch(
  () => props.modelValue,
  (value) => {
    const currentMarkdown = editor.value?.getMarkdown();
    if (value !== currentMarkdown) {
      editor.value?.commands.setContent(value, { contentType: "markdown" });
    }
  }
);

function setLink() {
  const url = window.prompt("输入链接地址：");
  if (url === null) return;
  if (url === "") {
    editor.value.chain().focus().extendMarkRange("link").unsetLink().run();
    return;
  }
  const safe = url.replace(/^(javascript|data|vbscript):/i, "#");
  editor.value.chain().focus().extendMarkRange("link").setLink({ href: safe }).run();
}

function triggerImageUpload() {
  imageInput.value?.click();
}

async function handleImageUpload(event) {
  const file = event.target.files?.[0];
  if (!file) return;
  await uploadAndInsertImage(file);
  event.target.value = "";
}

async function uploadAndInsertImage(file) {
  try {
    const compressed = await compressImage(file);
    const { data } = await uploadImage(compressed);
    editor.value.chain().focus().setImage({ src: data.url }).run();
  } catch {
    window.alert("图片上传失败，请稍后重试。");
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

<style scoped>
.rich-editor {
  border: 1px solid var(--border);
  border-radius: 18px;
  overflow: hidden;
  background: rgba(8, 11, 17, 0.4);
}

.rich-editor__toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
}

.toolbar-group {
  display: flex;
  gap: 2px;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border);
  margin: 0 4px;
}

.rich-editor__toolbar button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--text-soft);
  cursor: pointer;
  font-size: 0.85rem;
  font-family: inherit;
  transition: background 0.15s, color 0.15s;
}

.rich-editor__toolbar button:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f6f1ea;
}

.rich-editor__toolbar button.active {
  background: rgba(255, 138, 76, 0.2);
  color: var(--accent);
}

.toolbar-icon {
  font-size: 0.9rem;
  line-height: 1;
}

.rich-editor__content {
  min-height: 400px;
  max-height: 70vh;
  overflow-y: auto;
}

.rich-editor__content :deep(.tiptap) {
  padding: 20px 24px;
  min-height: 400px;
  outline: none;
}

.rich-editor__content :deep(.tiptap p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  float: left;
  color: var(--text-soft);
  pointer-events: none;
  height: 0;
  font-style: italic;
}

.rich-editor__content :deep(.tiptap h1) {
  font-size: 1.8rem;
  margin: 1.2em 0 0.6em;
}

.rich-editor__content :deep(.tiptap h2) {
  font-size: 1.45rem;
  margin: 1em 0 0.5em;
}

.rich-editor__content :deep(.tiptap h3) {
  font-size: 1.2rem;
  margin: 0.8em 0 0.4em;
}

.rich-editor__content :deep(.tiptap pre) {
  background: #0d0d0d;
  color: #fff;
  font-family: "JetBrains Mono", "Fira Code", monospace;
  padding: 14px 18px;
  border-radius: 10px;
  overflow-x: auto;
}

.rich-editor__content :deep(.tiptap pre code) {
  background: none;
  color: inherit;
  padding: 0;
  font-size: 0.9rem;
}

.rich-editor__content :deep(.tiptap code) {
  background: rgba(255, 138, 76, 0.12);
  color: var(--accent-soft);
  padding: 2px 6px;
  border-radius: 5px;
  font-size: 0.88em;
}

.rich-editor__content :deep(.tiptap blockquote) {
  border-left: 3px solid var(--accent);
  margin: 1em 0;
  padding-left: 16px;
  color: var(--text-soft);
}

.rich-editor__content :deep(.tiptap ul),
.rich-editor__content :deep(.tiptap ol) {
  padding-left: 24px;
}

.rich-editor__content :deep(.tiptap li) {
  margin: 4px 0;
}

.rich-editor__content :deep(.tiptap ul[data-type="taskList"]) {
  list-style: none;
  padding-left: 0;
}

.rich-editor__content :deep(.tiptap ul[data-type="taskList"] li) {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.rich-editor__content :deep(.tiptap ul[data-type="taskList"] li label input[type="checkbox"]) {
  margin-top: 5px;
  accent-color: var(--accent);
}

.rich-editor__content :deep(.tiptap img) {
  max-width: 100%;
  height: auto;
  border-radius: 10px;
  margin: 12px 0;
}

.rich-editor__content :deep(.tiptap table) {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.rich-editor__content :deep(.tiptap td),
.rich-editor__content :deep(.tiptap th) {
  border: 1px solid var(--border);
  padding: 8px 10px;
  min-width: 80px;
  position: relative;
}

.rich-editor__content :deep(.tiptap th) {
  background: rgba(255, 255, 255, 0.06);
  font-weight: 600;
}

.rich-editor__content :deep(.tiptap hr) {
  border: none;
  border-top: 1px solid var(--border);
  margin: 1.5em 0;
}

.rich-editor__content :deep(.tiptap .editor-link) {
  color: var(--accent-soft);
  text-decoration: underline;
  cursor: pointer;
}
</style>
