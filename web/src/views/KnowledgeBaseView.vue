<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">知识库</p>
        <h3>我的知识库</h3>
      </div>
    </div>

    <div class="services-grid">
      <article class="panel-card service-card">
        <p class="eyebrow">AI 检索知识库</p>
        <h4>录入文本、构建向量知识库，用自然语言检索你的私有知识</h4>
      </article>
      <article class="panel-card service-card">
        <p class="eyebrow">配置</p>
        <h4>请先在设置中配置 API Key</h4>
      </article>
    </div>

    <div v-if="kbs.length" class="project-grid">
      <article v-for="kb in kbs" :key="kb.id" class="project-card panel-card" @click="$router.push(`/knowledge-base/${kb.id}`)">
        <div class="project-card__head">
          <div>
            <h3>{{ kb.name }}</h3>
            <p class="table-note">{{ kb.docCount }} 篇文档 · {{ formatDate(kb.updatedAt) }}</p>
          </div>
        </div>
        <p v-if="kb.description" class="detail-summary">{{ kb.description }}</p>
        <footer class="project-footer">
          <div></div>
          <button class="inline-link delete-link" @click.stop="handleDelete(kb)">删除</button>
        </footer>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>还没有知识库</h4>
      <p>创建一个知识库来管理你的私有知识。</p>
    </div>

    <div class="action-strip" style="margin-top:20px">
      <button class="solid-btn" @click="showCreate = true">新建知识库</button>
    </div>

    <Teleport to="body">
      <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
        <div class="modal-card panel-card">
          <h3>新建知识库</h3>
          <div class="stack-form" style="margin-top:16px">
            <label>
              名称
              <input v-model.trim="newName" class="field-input" placeholder="例如：Go 学习笔记" />
            </label>
            <label>
              描述（可选）
              <textarea v-model.trim="newDesc" class="field-area field-area--small" placeholder="这个知识库主要记录什么内容？"></textarea>
            </label>
          </div>
          <div class="action-strip" style="margin-top:20px">
            <button class="ghost-btn" @click="showCreate = false">取消</button>
            <button class="solid-btn" :disabled="!newName || creating" @click="handleCreate">
              {{ creating ? "创建中..." : "创建" }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { createKnowledgeBase, deleteKnowledgeBase, listKnowledgeBases } from "../api/knowledgeBase";

const kbs = ref([]);
const showCreate = ref(false);
const newName = ref("");
const newDesc = ref("");
const creating = ref(false);

async function load() {
  const { data } = await listKnowledgeBases();
  kbs.value = data.items || [];
}

async function handleCreate() {
  if (!newName.value.trim()) return;
  creating.value = true;
  try {
    await createKnowledgeBase({ name: newName.value.trim(), description: newDesc.value.trim() });
    newName.value = "";
    newDesc.value = "";
    showCreate.value = false;
    await load();
  } catch (e) {
    alert(e?.response?.data?.message || "创建失败");
  } finally {
    creating.value = false;
  }
}

async function handleDelete(kb) {
  if (!confirm(`确定要删除知识库「${kb.name}」及其所有文档吗？此操作不可撤销。`)) return;
  try {
    await deleteKnowledgeBase(kb.id);
    await load();
  } catch (e) {
    alert("删除失败：" + (e?.response?.data?.message || e.message));
  }
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(load);
</script>

<style scoped>
.project-card {
  cursor: pointer;
  transition: border-color 0.2s;
}
.project-card:hover {
  border-color: rgba(255, 217, 142, 0.4);
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
  text-decoration: underline;
}
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.modal-card {
  width: min(460px, 90vw);
  padding: 28px;
  border-radius: 26px;
}
</style>
