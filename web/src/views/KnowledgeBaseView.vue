<template>
  <section class="uc-kb">
    <section class="uc-hero uc-hero--kb">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">知识库</p>
        <h2 class="uc-hero__title">我的知识库</h2>
        <p class="uc-hero__text">录入文本、构建向量知识库，用自然语言检索你的私有知识。请先在设置中配置 API Key。</p>
      </div>
      <div class="uc-hero__actions">
        <button class="solid-btn" @click="showCreate = true">新建知识库</button>
      </div>
    </section>

    <div v-if="kbs.length" class="kb-grid">
      <article v-for="kb in kbs" :key="kb.id" class="kb-card" @click="$router.push({ name: 'uc-knowledge-base-detail', params: { id: kb.id } })">
        <div class="kb-card__head">
          <div>
            <h3>{{ kb.name }}</h3>
            <p class="table-note">{{ kb.docCount }} 篇文档 · {{ formatDate(kb.updatedAt) }}</p>
          </div>
        </div>
        <p v-if="kb.description" class="detail-summary">{{ kb.description }}</p>
        <footer class="kb-card__footer">
          <div></div>
          <button class="inline-link delete-link" @click.stop="handleDelete(kb)">删除</button>
        </footer>
      </article>
    </div>

    <div v-if="errorMessage" class="empty-panel">
      <h4>出错了</h4>
      <p>{{ errorMessage }}</p>
    </div>
    <div v-else-if="!kbs.length" class="empty-panel">
      <h4>还没有知识库</h4>
      <p>创建一个知识库来管理你的私有知识。</p>
    </div>

    <Teleport to="body">
      <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
        <div class="modal-card">
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
import { formatDate } from "../utils/date";

const kbs = ref([]);
const showCreate = ref(false);
const newName = ref("");
const newDesc = ref("");
const creating = ref(false);
const errorMessage = ref("");

async function load() {
  errorMessage.value = "";
  try {
    const { data } = await listKnowledgeBases();
    kbs.value = data.items || [];
  } catch (e) {
    errorMessage.value = e?.response?.data?.message || "加载失败";
  }
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
    errorMessage.value = e?.response?.data?.message || "创建失败";
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
    errorMessage.value = "删除失败：" + (e?.response?.data?.message || e.message);
  }
}

onMounted(load);
</script>

<style scoped>
.uc-kb {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.uc-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    var(--panel);
  padding: clamp(32px, 6vw, 56px) clamp(24px, 4vw, 48px);
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) auto;
  gap: 32px;
  align-items: center;
}

.uc-hero__copy {
  display: grid;
  gap: 12px;
}

.uc-hero__eyebrow {
  color: var(--accent);
  font-weight: 600;
}

.uc-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
}

.uc-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: var(--text-soft);
}

.kb-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.kb-card {
  padding: 22px;
  border-radius: 28px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: transform 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
}

.kb-card:hover {
  transform: translateY(-5px);
  border-color: rgba(255, 209, 102, 0.45);
  box-shadow: 0 24px 50px rgba(0, 0, 0, 0.22);
}

.kb-card:first-child {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.14), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.10), transparent 50%),
    rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 209, 102, 0.2);
}

.kb-card:first-child h3 {
  color: inherit;
}

.kb-card:first-child .detail-summary {
  color: var(--text-soft);
}

.kb-card:first-child .table-note {
  color: var(--text-soft);
}

.kb-card__head {
  margin-bottom: 8px;
}

.kb-card__head h3 {
  margin: 0;
  font-size: 1.1rem;
}

.kb-card__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.delete-link {
  color: #dc2626;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
}

.delete-link:hover {
  text-decoration: underline;
  color: #b91c1c;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal-card {
  width: min(460px, 90vw);
  padding: 28px;
  border-radius: 26px;
  background: rgba(30, 30, 30, 0.95);
  border: 1px solid var(--border);
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(20px);
}

.modal-card h3 {
  margin: 0;
}

@media (max-width: 768px) {
  .uc-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .uc-hero__actions {
    width: 100%;
  }

  .uc-hero__actions .solid-btn {
    width: 100%;
    justify-content: center;
  }

  .kb-grid {
    grid-template-columns: 1fr;
  }
}
</style>
