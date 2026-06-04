<template>
  <div v-if="visible" class="ai-review-mask" @click.self="close">
    <div class="ai-review-dialog">
      <div class="ai-review-header">
        <div>
          <p class="admin-label">AI 审核辅助</p>
          <h2>{{ targetTitle || "AI 审核分析" }}</h2>
          <p class="table-note">将内容复制到 AI 对话中分析，再把结果粘贴回来保存</p>
        </div>
        <button class="role-btn" @click="close">关闭</button>
      </div>

      <div v-if="loading" class="ai-review-loading">
        <p>加载内容中...</p>
      </div>

      <template v-if="!loading">
        <div class="ai-review-section">
          <div class="ai-review-section-head">
            <h3>待分析内容</h3>
            <button class="role-btn" @click="copyPrompt">
              {{ copied ? "已复制" : "复制内容到剪贴板" }}
            </button>
          </div>
          <pre class="ai-review-content">{{ promptText }}</pre>
        </div>

        <div v-if="existingResult" class="ai-review-section">
          <div class="ai-review-section-head">
            <h3>已有分析结果（{{ formatDate(existingResult.createdAt) }}）</h3>
          </div>
          <div class="ai-review-result-display">
            <div class="result-row">
              <strong>风险等级</strong>
              <span :class="riskLevelClass(existingResult.riskLevel)">{{ riskLevelLabel(existingResult.riskLevel) }}</span>
            </div>
            <div v-if="existingResult.riskLabels?.length" class="result-row">
              <strong>风险标签</strong>
              <span>{{ existingResult.riskLabels.join(", ") }}</span>
            </div>
            <div v-if="existingResult.summary" class="result-row">
              <strong>分析摘要</strong>
              <p>{{ existingResult.summary }}</p>
            </div>
            <div v-if="existingResult.suspiciousSegments?.length" class="result-row">
              <strong>可疑片段</strong>
              <ul>
                <li v-for="(seg, i) in existingResult.suspiciousSegments" :key="i">{{ seg }}</li>
              </ul>
            </div>
            <div v-if="existingResult.suggestion" class="result-row">
              <strong>审核建议</strong>
              <p>{{ existingResult.suggestion }}</p>
            </div>
          </div>
        </div>

        <div class="ai-review-section">
          <div class="ai-review-section-head">
            <h3>{{ existingResult ? "更新分析结果" : "粘贴分析结果" }}</h3>
            <p class="table-note">AI 仅提供辅助参考，不自动代替人工裁决</p>
          </div>
          <form class="ai-review-form" @submit.prevent="submitResult">
            <div class="form-row">
              <label>风险等级</label>
              <select v-model="form.riskLevel" class="filter-select">
                <option value="low">低风险</option>
                <option value="mid">中风险</option>
                <option value="high">高风险</option>
              </select>
            </div>
            <div class="form-row">
              <label>风险标签（逗号分隔）</label>
              <input v-model="form.riskLabelsInput" class="filter-select filter-input" placeholder="例如：内容质量低, 结构不完整" />
            </div>
            <div class="form-row">
              <label>分析摘要</label>
              <textarea v-model="form.summary" class="filter-select filter-input form-textarea" placeholder="AI 分析的整体摘要"></textarea>
            </div>
            <div class="form-row">
              <label>可疑片段（每行一个）</label>
              <textarea v-model="form.suspiciousInput" class="filter-select filter-input form-textarea" placeholder="可疑片段1&#10;可疑片段2"></textarea>
            </div>
            <div class="form-row">
              <label>审核建议</label>
              <textarea v-model="form.suggestion" class="filter-select filter-input form-textarea" placeholder="建议通过 / 建议驳回 / 需修改后重新提交..."></textarea>
            </div>
            <div v-if="saveError" class="form-error">{{ saveError }}</div>
            <div class="form-actions">
              <button type="submit" class="solid-btn" :disabled="saving">
                {{ saving ? "保存中..." : "保存 AI 审核结果" }}
              </button>
            </div>
          </form>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { getAIReviewContent, getAIReviewResult, saveAIReviewResult } from "../api/dashboard";

const props = defineProps({
  visible: Boolean,
  targetType: String,
  targetId: Number,
  targetTitle: String
});

const emit = defineEmits(["close", "saved"]);

const loading = ref(false);
const copied = ref(false);
const promptText = ref("");
const existingResult = ref(null);
const saving = ref(false);
const saveError = ref("");

const form = ref({
  riskLevel: "low",
  riskLabelsInput: "",
  summary: "",
  suspiciousInput: "",
  suggestion: ""
});

watch(() => props.visible, async (show) => {
  if (!show) return;
  loading.value = true;
  existingResult.value = null;
  saveError.value = "";
  form.value = { riskLevel: "low", riskLabelsInput: "", summary: "", suspiciousInput: "", suggestion: "" };
  try {
    const [contentRes, resultRes] = await Promise.all([
      getAIReviewContent(props.targetType, props.targetId),
      getAIReviewResult(props.targetType, props.targetId).catch(() => null)
    ]);
    promptText.value = contentRes.data.prompt || "";
    if (resultRes?.data?.item) {
      const item = resultRes.data.item;
      existingResult.value = item;
      form.value = {
        riskLevel: item.riskLevel || "low",
        riskLabelsInput: (item.riskLabels || []).join(", "),
        summary: item.summary || "",
        suspiciousInput: (item.suspiciousSegments || []).join("\n"),
        suggestion: item.suggestion || ""
      };
    }
  } catch (err) {
    promptText.value = "加载内容失败，请重试";
  } finally {
    loading.value = false;
  }
});

async function copyPrompt() {
  try {
    await navigator.clipboard.writeText(promptText.value);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  } catch {
    const textarea = document.createElement("textarea");
    textarea.value = promptText.value;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  }
}

async function submitResult() {
  saving.value = true;
  saveError.value = "";
  try {
    const payload = {
      riskLevel: form.value.riskLevel,
      riskLabels: form.value.riskLabelsInput.split(",").map((s) => s.trim()).filter(Boolean),
      summary: form.value.summary,
      suspiciousSegments: form.value.suspiciousInput.split("\n").map((s) => s.trim()).filter(Boolean),
      suggestion: form.value.suggestion
    };
    await saveAIReviewResult(props.targetType, props.targetId, payload);
    existingResult.value = null;
    const resultRes = await getAIReviewResult(props.targetType, props.targetId);
    if (resultRes?.data?.item) {
      existingResult.value = resultRes.data.item;
    }
    emit("saved", props.targetId);
  } catch (err) {
    saveError.value = err.response?.data?.message || "保存失败，请重试";
  } finally {
    saving.value = false;
  }
}

function close() {
  emit("close");
}

function riskLevelClass(level) {
  return { "risk-low": level === "low", "risk-mid": level === "mid", "risk-high": level === "high" };
}

function riskLevelLabel(level) {
  return { low: "低风险", mid: "中风险", high: "高风险" }[level] || level;
}

function formatDate(dateStr) {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")} ${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`;
}
</script>

<style scoped>
.ai-review-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.ai-review-dialog {
  background: var(--bg-card, #1a1a2e);
  border-radius: 16px;
  padding: 28px;
  max-width: 720px;
  width: 90%;
  max-height: 85vh;
  overflow-y: auto;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.ai-review-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.ai-review-header h2 {
  margin: 4px 0 0;
}

.ai-review-loading {
  text-align: center;
  padding: 40px 0;
}

.ai-review-section {
  margin-bottom: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding-top: 16px;
}

.ai-review-section-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.ai-review-section-head h3 {
  margin: 0;
  font-size: 1rem;
}

.ai-review-content {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  padding: 16px;
  max-height: 240px;
  overflow-y: auto;
  font-size: 0.85rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.ai-review-result-display {
  display: grid;
  gap: 10px;
}

.result-row {
  display: grid;
  gap: 4px;
}

.result-row strong {
  font-size: 0.8rem;
  color: var(--text-soft, #888);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.result-row p,
.result-row ul {
  margin: 0;
}

.result-row ul {
  padding-left: 20px;
}

.risk-low { color: #4caf50; }
.risk-mid { color: #ff9800; }
.risk-high { color: #f44336; }

.ai-review-form {
  display: grid;
  gap: 14px;
}

.form-row {
  display: grid;
  gap: 6px;
}

.form-row label {
  font-size: 0.85rem;
  font-weight: 600;
}

.form-textarea {
  min-height: 72px;
  resize: vertical;
  font-family: inherit;
}

.form-error {
  color: #f44336;
  font-size: 0.85rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .ai-review-dialog {
    max-width: 100%;
    width: 100%;
    max-height: 100vh;
    height: 100vh;
    border-radius: 0;
    padding: 20px;
  }
  .ai-review-mask {
    align-items: flex-start;
    padding-top: 0;
  }
}

@media (max-width: 480px) {
  .ai-review-dialog {
    padding: 14px;
  }
  .ai-review-header {
    flex-direction: column;
    gap: 12px;
  }
  .ai-review-section-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .ai-review-content {
    max-height: 160px;
    padding: 12px;
  }
  .form-actions {
    flex-direction: column;
  }
  .form-actions button {
    width: 100%;
  }
}
</style>
