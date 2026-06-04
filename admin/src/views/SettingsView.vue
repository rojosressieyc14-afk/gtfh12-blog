<template>
  <div class="settings-page">
    <section class="settings-hero">
      <div>
        <p class="login-tag">System Config</p>
        <h1>系统设置</h1>
        <p class="login-text">
          配置平台风控阈值与系统参数。修改后立即生效。
        </p>
      </div>
    </section>

    <section class="settings-grid">
      <article class="settings-card">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">风控阈值</p>
            <h2>自动封禁设置</h2>
          </div>
        </div>

        <div class="settings-field">
          <label>24 小时内触发风控次数上限</label>
          <p class="settings-hint">普通用户达到此次数后将被自动封禁。设置为 0 表示不限制。</p>
          <div class="settings-field-row">
            <input
              v-model.number="banThreshold"
              class="filter-select filter-input settings-input"
              type="number"
              min="0"
              max="100"
            />
            <span class="settings-unit">次</span>
          </div>
        </div>

        <div v-if="message" class="settings-flash">{{ message }}</div>

        <div class="settings-actions">
          <button
            :disabled="saving || !isDirty"
            @click="saveSettings"
          >
            {{ saving ? "保存中..." : "保存设置" }}
          </button>
          <button v-if="isDirty" class="ghost-action" @click="resetForm">撤销更改</button>
        </div>
      </article>

      <article class="settings-card">
        <div class="review-head review-head--stack">
          <div>
            <p class="admin-label">系统信息</p>
            <h2>运行状态</h2>
          </div>
        </div>

        <div class="settings-info-list">
          <div class="settings-info-row">
            <span>风控阈值</span>
            <strong>{{ banThreshold }} 次 / 24h</strong>
          </div>
          <div class="settings-info-row">
            <span>当前状态</span>
            <strong :class="banThreshold > 0 ? 'status-active' : 'status-disabled'">
              {{ banThreshold > 0 ? "自动封禁已启用" : "自动封禁已关闭" }}
            </strong>
          </div>
          <div class="settings-info-row">
            <span>说明</span>
            <strong>修改后立即生效，无需重启服务</strong>
          </div>
        </div>
      </article>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { getModerationSettings, updateModerationSettings } from "../api/dashboard";

const banThreshold = ref(0);
const originalThreshold = ref(0);
const saving = ref(false);
const message = ref("");

const isDirty = computed(() => banThreshold.value !== originalThreshold.value);

function say(text) {
  message.value = text;
  window.setTimeout(() => {
    if (message.value === text) {
      message.value = "";
    }
  }, 2200);
}

async function loadSettings() {
  try {
    const { data } = await getModerationSettings();
    banThreshold.value = data.banThreshold ?? 0;
    originalThreshold.value = banThreshold.value;
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "加载设置失败。");
  }
}

function resetForm() {
  banThreshold.value = originalThreshold.value;
}

async function saveSettings() {
  if (!isDirty.value || saving.value) return;
  saving.value = true;
  try {
    const { data } = await updateModerationSettings({ banThreshold: banThreshold.value });
    banThreshold.value = data.banThreshold ?? banThreshold.value;
    originalThreshold.value = banThreshold.value;
    say("设置已保存。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "保存设置失败。");
  } finally {
    saving.value = false;
  }
}

onMounted(loadSettings);
</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  padding: 28px;
  color: inherit;
}

.settings-hero,
.settings-card {
  border: 1px solid var(--border);
  background: var(--panel);
  backdrop-filter: blur(18px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.settings-hero {
  padding: 28px;
  border-radius: 32px;
}

.settings-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.6rem);
  line-height: 1;
}

.settings-grid {
  margin-top: 20px;
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(280px, 0.8fr);
  gap: 20px;
}

.settings-card {
  padding: 24px;
  border-radius: 28px;
}

.settings-field {
  display: grid;
  gap: 8px;
  margin-top: 16px;
}

.settings-field label {
  font-weight: 600;
}

.settings-hint {
  color: var(--soft);
  font-size: 0.9rem;
  margin: 0;
}

.settings-field-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.settings-input {
  width: 120px;
}

.settings-unit {
  color: var(--soft);
}

.settings-flash {
  margin-top: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 166, 77, 0.14);
  color: #ffe3c2;
}

.settings-actions {
  margin-top: 18px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.settings-actions button {
  padding: 12px 16px;
  border-radius: 16px;
  border: 1px solid var(--border);
  cursor: pointer;
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
}

.settings-actions button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.ghost-action {
  background: rgba(255, 255, 255, 0.05) !important;
  color: inherit !important;
}

.settings-info-list {
  display: grid;
  gap: 14px;
  margin-top: 16px;
}

.settings-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.03);
}

.settings-info-row span {
  color: var(--soft);
}

.status-active {
  color: #83f28f;
}

.status-disabled {
  color: #ffbcbc;
}

@media (max-width: 960px) {
  .settings-page {
    padding: 20px;
  }

  .settings-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .settings-page {
    padding: 16px;
  }
  .settings-hero {
    padding: 20px;
  }
  .settings-card {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .settings-page {
    padding: 12px;
  }
  .settings-hero {
    padding: 16px;
  }
  .settings-card {
    padding: 14px;
  }
  .settings-field-row {
    flex-direction: column;
    align-items: flex-start;
  }
  .settings-info-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }
  .settings-actions {
    flex-direction: column;
  }
  .settings-actions button {
    width: 100%;
  }
}
</style>
