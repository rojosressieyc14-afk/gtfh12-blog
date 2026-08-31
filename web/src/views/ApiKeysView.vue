<template>
  <section class="uc-api-keys">
    <section class="uc-hero uc-hero--api-keys">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">设置</p>
        <h2 class="uc-hero__title">API Key 管理</h2>
        <p class="uc-hero__text">你的 API Key 经过 AES-256-GCM 加密存储，仅在你使用时解密，不会暴露。</p>
      </div>
      <div class="uc-hero__actions">
        <button class="solid-btn" @click="showAdd = true">添加 API Key</button>
      </div>
    </section>

    <div v-if="keys.length" class="key-list">
      <article v-for="key in keys" :key="key.id" class="key-card">
        <div class="key-card__head">
          <div class="key-card__info">
            <span class="provider-badge">{{ key.provider }}</span>
            <code class="key-prefix">{{ key.keyPrefix }}</code>
          </div>
          <button class="inline-link delete-link" @click="handleDelete(key)">删除</button>
        </div>
        <p class="table-note">
          {{ key.baseURL }} · {{ key.lastUsedAt ? "上次使用: " + formatDate(key.lastUsedAt) : "尚未使用" }}
        </p>
      </article>
    </div>

    <div v-else class="empty-panel">
      <h4>还没有配置 API Key</h4>
      <p>添加 API Key 后可在 AI 面试中使用你自己的密钥。</p>
    </div>

    <p v-if="errorMessage" class="error-text" style="margin-top:12px">{{ errorMessage }}</p>

    <Teleport to="body">
      <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
        <div class="modal-card">
          <h3>添加 API Key</h3>
          <div class="stack-form" style="margin-top:16px">
            <label>
              提供商
              <select v-model="newProvider" class="field-input">
                <option value="deepseek">DeepSeek</option>
                <option value="openai">OpenAI</option>
              </select>
            </label>
            <label>
              API Key
              <input v-model="newKey" class="field-input" type="password" :placeholder="newProvider === 'deepseek' ? 'sk-...' : 'sk-...'" />
            </label>
            <label>
              自定义 API 地址（可选）
              <input v-model="newBaseURL" class="field-input" :placeholder="defaultBaseURL" />
            </label>
          </div>
          <div class="action-strip" style="margin-top:20px">
            <button class="ghost-btn" @click="showAdd = false">取消</button>
            <button class="solid-btn" :disabled="!newKey.trim() || adding" @click="handleAdd">
              {{ adding ? "添加中..." : "添加" }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { createApiKey, deleteApiKey, listApiKeys } from "../api/apiKey";

const keys = ref([]);
const showAdd = ref(false);
const errorMessage = ref("");
const newProvider = ref("deepseek");
const newKey = ref("");
const newBaseURL = ref("");
const adding = ref(false);

const defaultBaseURL = computed(() => {
  return newProvider.value === "deepseek" ? "https://api.deepseek.com/v1" : "https://api.openai.com/v1";
});

async function load() {
  try {
    const { data } = await listApiKeys();
    keys.value = data.items || [];
  } catch {
    // API key endpoints may not be available if encryption key not configured
  }
}

async function handleAdd() {
  if (!newKey.value.trim()) return;
  adding.value = true;
  try {
    await createApiKey({
      provider: newProvider.value,
      key: newKey.value.trim(),
      baseURL: newBaseURL.value.trim() || defaultBaseURL.value
    });
    newKey.value = "";
    newBaseURL.value = "";
    showAdd.value = false;
    await load();
  } catch (e) {
    errorMessage.value = e?.response?.data?.message || "添加失败";
  } finally {
    adding.value = false;
  }
}

async function handleDelete(key) {
  if (!confirm(`确定要删除此 API Key (${key.keyPrefix}) 吗？`)) return;
  try {
    await deleteApiKey(key.id);
    await load();
  } catch (e) {
    errorMessage.value = "删除失败：" + (e?.response?.data?.message || e.message);
  }
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(load);
</script>

<style scoped>
.uc-api-keys {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.uc-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
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
  color: #b4530a;
  font-weight: 600;
}

.uc-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
  color: #1d1d1f;
}

.uc-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: #48484a;
}

.key-list {
  display: grid;
  gap: 12px;
}

.key-card {
  padding: 18px 20px;
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.5);
  transition: border-color 0.2s;
}

.key-card:first-child {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.key-card__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.key-card__info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.provider-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 600;
  background: rgba(249, 115, 22, 0.1);
  color: #c2410c;
  text-transform: uppercase;
}

.key-prefix {
  font-family: monospace;
  font-size: 0.92rem;
  color: #48484a;
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
  color: #b91c1c;
  text-decoration: underline;
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
  background: #fff;
  color: #1d1d1f;
  border: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.12);
}

.modal-card h3 {
  margin: 0;
  color: #1d1d1f;
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
}
</style>
