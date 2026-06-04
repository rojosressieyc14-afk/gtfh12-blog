<template>
  <div class="sensitive-page">
    <section class="sensitive-hero">
      <div>
        <p class="login-tag">Sensitive Guard</p>
        <h1>违禁词词库管理</h1>
        <p class="login-text">
          系统内置词库始终生效，这里维护的是部署后的自定义补充词。新增后会立即进入后台拦截规则。
        </p>
      </div>
      <div class="sensitive-hero__stats">
        <article class="sensitive-stat">
          <span>当前页词条</span>
          <strong>{{ items.length }}</strong>
        </article>
        <article class="sensitive-stat">
          <span>总词条数</span>
          <strong>{{ total }}</strong>
        </article>
      </div>
    </section>

    <section class="sensitive-layout">
      <article class="sensitive-panel">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">新增词条</p>
            <h2>补充自定义违禁词</h2>
          </div>
        </div>

        <form class="sensitive-form" @submit.prevent="submitWord">
          <label>
            违禁词
            <input v-model.trim="form.word" placeholder="例如：你想追加屏蔽的词或变体" />
          </label>
          <label>
            分类
            <input v-model.trim="form.category" placeholder="例如：spam / politics / custom" />
          </label>
          <label>
            备注
            <textarea v-model.trim="form.note" placeholder="记录补充原因，方便后续排查"></textarea>
          </label>
          <div class="sensitive-actions">
            <button type="submit" :disabled="submitting">{{ submitting ? "保存中..." : "新增词条" }}</button>
            <button type="button" class="ghost-action" @click="resetForm">清空</button>
          </div>
        </form>
      </article>

      <article class="sensitive-panel">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">词库列表</p>
            <h2>自定义词条</h2>
          </div>
          <div class="toolbar-row">
            <input
              v-model.trim="keyword"
              class="filter-select filter-input"
              placeholder="搜索词条 / 分类 / 备注"
              @keyup.enter="changePage(1)"
            />
            <button @click="changePage(1)">搜索</button>
          </div>
        </div>

        <div v-if="message" class="sensitive-flash">{{ message }}</div>

        <div class="sensitive-list">
          <article v-for="item in items" :key="item.id" class="sensitive-card">
            <div class="sensitive-card__head">
              <strong>{{ item.word }}</strong>
              <span>{{ item.category || "custom" }}</span>
            </div>
            <p>{{ item.note || "暂无备注" }}</p>
            <footer>
              <span>{{ formatDate(item.createdAt) }}</span>
              <button class="delete-btn" @click="removeWord(item)">删除</button>
            </footer>
          </article>
        </div>

        <div v-if="!items.length" class="compact-empty">当前没有自定义违禁词。</div>

        <div class="pager-row">
          <button :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <span>第 {{ page }} / {{ totalPages }} 页</span>
          <button :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
        </div>
      </article>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { createSensitiveWord, deleteSensitiveWord, getSensitiveWords } from "../api/dashboard";

const items = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;
const keyword = ref("");
const submitting = ref(false);
const message = ref("");
const form = reactive({
  word: "",
  category: "custom",
  note: ""
});

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

function say(text) {
  message.value = text;
  window.setTimeout(() => {
    if (message.value === text) {
      message.value = "";
    }
  }, 2200);
}

async function loadWords() {
  const { data } = await getSensitiveWords({ page: page.value, pageSize, keyword: keyword.value });
  items.value = data.items || [];
  total.value = data.pagination?.total || 0;
}

async function changePage(nextPage) {
  page.value = nextPage;
  await loadWords();
}

function resetForm() {
  form.word = "";
  form.category = "custom";
  form.note = "";
}

async function submitWord() {
  if (!form.word) {
    say("请先填写违禁词。");
    return;
  }
  submitting.value = true;
  try {
    await createSensitiveWord(form);
    resetForm();
    await changePage(1);
    say("违禁词已加入词库。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "新增违禁词失败。");
  } finally {
    submitting.value = false;
  }
}

async function removeWord(item) {
  if (!window.confirm(`确定删除违禁词“${item.word}”吗？`)) return;
  try {
    await deleteSensitiveWord(item.id);
    await loadWords();
    say("违禁词已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除违禁词失败。");
  }
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

onMounted(loadWords);
</script>

<style scoped>
.sensitive-page {
  min-height: 100vh;
  padding: 28px;
  color: inherit;
}

.sensitive-hero,
.sensitive-panel {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.sensitive-hero {
  padding: 28px;
  border-radius: 32px;
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(240px, 0.75fr);
  gap: 20px;
}

.sensitive-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.6rem);
  line-height: 1;
}

.sensitive-hero__stats {
  display: grid;
  gap: 14px;
}

.sensitive-stat {
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 14, 19, 0.52);
  display: grid;
  gap: 10px;
}

.sensitive-stat span,
.sensitive-card p {
  color: var(--soft);
}

.sensitive-stat strong {
  font-size: 2rem;
}

.sensitive-layout {
  margin-top: 20px;
  display: grid;
  grid-template-columns: minmax(320px, 0.82fr) minmax(0, 1.18fr);
  gap: 20px;
}

.sensitive-panel {
  padding: 24px;
  border-radius: 28px;
}

.sensitive-form {
  display: grid;
  gap: 14px;
}

.sensitive-form label {
  display: grid;
  gap: 8px;
}

.sensitive-form input,
.sensitive-form textarea,
.sensitive-actions button,
.pager-row button {
  border-radius: 16px;
  border: 1px solid var(--border);
  padding: 12px 14px;
  background: rgba(10, 14, 19, 0.7);
  color: inherit;
}

.sensitive-form textarea {
  min-height: 120px;
  resize: vertical;
}

.sensitive-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.sensitive-actions button {
  cursor: pointer;
}

.sensitive-actions button:first-child {
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
}

.ghost-action {
  background: rgba(255, 255, 255, 0.05) !important;
  color: inherit !important;
}

.sensitive-flash {
  margin-bottom: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 166, 77, 0.14);
  color: #ffe3c2;
}

.sensitive-list {
  display: grid;
  gap: 12px;
}

.sensitive-card {
  padding: 16px 18px;
  border-radius: 20px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.03);
}

.sensitive-card__head,
.sensitive-card footer {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.sensitive-card strong {
  font-size: 1.08rem;
}

.sensitive-card span {
  color: #ffd79d;
}

.sensitive-card p {
  margin: 10px 0 14px;
}

.delete-btn {
  cursor: pointer;
  border: 1px solid rgba(255, 139, 139, 0.3);
  background: rgba(255, 139, 139, 0.12);
  color: #ffbcbc;
  border-radius: 14px;
  padding: 8px 12px;
}

@media (max-width: 960px) {
  .sensitive-page {
    padding: 20px;
  }

  .sensitive-hero,
  .sensitive-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .sensitive-page {
    padding: 16px;
  }
  .sensitive-hero {
    padding: 20px;
  }
  .sensitive-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .sensitive-page {
    padding: 12px;
  }
  .sensitive-hero {
    padding: 16px;
  }
  .sensitive-panel {
    padding: 14px;
  }
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .sensitive-card__head,
  .sensitive-card footer {
    flex-direction: column;
    align-items: flex-start;
  }
  .sensitive-actions {
    flex-direction: column;
  }
  .sensitive-actions button {
    width: 100%;
  }
}
</style>
