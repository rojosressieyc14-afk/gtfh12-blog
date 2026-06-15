<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">设置</p>
        <h3>API Key 管理</h3>
      </div>
    </div>

    <div class="services-grid">
      <article class="panel-card service-card">
        <p class="eyebrow">安全</p>
        <h4>你的 API Key 经过 AES-256-GCM 加密存储，仅在你使用时解密，不会暴露</h4>
      </article>
    </div>

    <div v-if="keys.length" class="key-list" style="margin-top:20px">
      <article v-for="key in keys" :key="key.id" class="key-card panel-card">
        <div class="key-card__head">
          <div>
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

    <div v-else class="empty-panel" style="margin-top:20px">
      <h4>还没有配置 API Key</h4>
      <p>添加 API Key 后可在 AI 面试中使用你自己的密钥。</p>
    </div>

    <div class="action-strip" style="margin-top:20px">
      <button class="solid-btn" @click="showAdd = true">添加 API Key</button>
    </div>

    <Teleport to="body">
      <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
        <div class="modal-card panel-card">
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
    alert(e?.response?.data?.message || "添加失败");
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
    alert("删除失败：" + (e?.response?.data?.message || e.message));
  }
}

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}

onMounted(load);
</script>

<style scoped>
.key-list {
  display: grid;
  gap: 12px;
}
.key-card {
  padding: 18px;
  border-radius: 20px;
}
.key-card__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.key-card__head div {
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
  background: rgba(255, 217, 142, 0.16);
  color: #ffd98e;
  text-transform: uppercase;
}
.key-prefix {
  font-family: monospace;
  font-size: 0.92rem;
  color: var(--soft, rgba(242,239,232,0.7));
}
.delete-link {
  color: #f87171;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  font: inherit;
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
