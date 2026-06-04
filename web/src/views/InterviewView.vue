<template>
  <div class="interview-page">
    <!-- Phase: Setup -->
    <section v-if="phase === 'setup'" class="interview-setup">
      <div class="setup-hero">
        <p class="page-label">AI Interview Agent</p>
        <h1>模拟面试官</h1>
        <p class="setup-desc">
          输入目标职位和简历，AI 面试官会模拟真实面试场景，逐题提问并评分。
        </p>
      </div>

      <div class="setup-form">
        <label class="setup-field">
          <span>应聘职位</span>
          <input v-model="position" class="setup-input" placeholder="例如：前端工程师、Go 后端开发、产品经理..." />
        </label>

        <label class="setup-field">
          <span>简历内容 <span class="field-hint">（粘贴文本或上传 .txt 文件）</span></span>
          <textarea v-model="resumeText" class="setup-textarea" placeholder="在此粘贴简历内容..." rows="6"></textarea>
          <label class="file-btn">
            <input type="file" accept=".txt" @change="handleFileUpload" />
            <span>上传 .txt 文件</span>
          </label>
        </label>

        <label class="setup-field">
          <span>面试题数量</span>
          <div class="qty-selector">
            <button class="qty-btn" :class="{ active: totalQuestions === n }" v-for="n in [3, 5, 8]" :key="n" @click="totalQuestions = n">{{ n }} 题</button>
          </div>
        </label>

        <button class="start-btn" :disabled="!position || starting" @click="startInterview">
          {{ starting ? "准备中..." : "开始面试" }}
        </button>
      </div>
    </section>

    <!-- Phase: Interview -->
    <section v-if="phase === 'interview'" class="interview-chat">
      <header class="chat-header">
        <div>
          <p class="page-label">AI Interview Agent</p>
          <h2>{{ position }}</h2>
          <p class="chat-progress">第 {{ currentRound }} / {{ totalQuestions }} 题</p>
        </div>
        <button class="end-btn" @click="confirmEnd">结束面试</button>
      </header>

      <div class="chat-messages" ref="chatRef">
        <div v-for="msg in messages" :key="msg.key" class="chat-msg" :class="msg.role">
          <div class="msg-avatar">{{ msg.role === 'agent' ? '🤖' : '👤' }}</div>
          <div class="msg-bubble">
            <div class="msg-text">{{ msg.text }}</div>
            <div v-if="msg.score !== undefined" class="msg-score">
              <span class="score-badge" :class="scoreClass(msg.score)">{{ msg.score }} 分</span>
              <span v-if="msg.feedback" class="score-feedback">{{ msg.feedback }}</span>
            </div>
          </div>
        </div>

        <div v-if="agentThinking" class="chat-msg agent">
          <div class="msg-avatar">🤖</div>
          <div class="msg-bubble msg-thinking">
            <span class="thinking-dot"></span><span class="thinking-dot"></span><span class="thinking-dot"></span>
          </div>
        </div>
      </div>

      <div class="chat-input-area">
        <div v-if="voiceSupported" class="voice-row">
          <button class="voice-btn" :class="{ recording }" @click="toggleVoice" :title="recording ? '点击停止录音' : '点击开始语音回答'">
            {{ recording ? '⏹ 停止录音' : '🎤 语音回答' }}
          </button>
          <span v-if="voiceText" class="voice-preview">{{ voiceText }}</span>
        </div>
        <div class="input-row">
          <textarea
            v-model="answerText"
            class="answer-input"
            placeholder="输入你的回答..."
            rows="2"
            @keydown.enter.ctrl="sendAnswer"
          ></textarea>
          <button class="send-btn" :disabled="!answerText.trim() || submitting" @click="sendAnswer">
            {{ submitting ? "发送中..." : "发送" }}
          </button>
        </div>
      </div>
    </section>

    <!-- Phase: Results -->
    <section v-if="phase === 'results'" class="interview-results">
      <div class="results-hero">
        <p class="page-label">面试完成</p>
        <h1>{{ position }}</h1>
        <p class="results-subtitle">共 {{ rounds.length }} 题 · AI 面试评估报告</p>
      </div>

      <div class="results-summary">
        <div class="summary-stat">
          <span class="stat-value" :class="scoreClass(overallScore)">{{ overallScore }}</span>
          <span class="stat-label">综合评分</span>
        </div>
        <div class="summary-stat">
          <span class="stat-value">{{ rounds.length }}</span>
          <span class="stat-label">答题数</span>
        </div>
        <div class="summary-stat">
          <span class="stat-value" :class="scoreClass(highestScore)">{{ highestScore }}</span>
          <span class="stat-label">最高分</span>
        </div>
        <div class="summary-stat">
          <span class="stat-value" :class="scoreClass(lowestScore)">{{ lowestScore }}</span>
          <span class="stat-label">最低分</span>
        </div>
      </div>

      <div class="results-detail">
        <article v-for="(r, i) in rounds" :key="r.id" class="result-card">
          <div class="result-head">
            <span class="result-num">第 {{ i + 1 }} 题</span>
            <span class="score-badge" :class="scoreClass(r.score)">{{ r.score }} 分</span>
          </div>
          <p class="result-q"><strong>问题：</strong>{{ r.question }}</p>
          <p class="result-a"><strong>回答：</strong>{{ r.answer }}</p>
          <p v-if="r.feedback" class="result-feedback"><strong>反馈：</strong>{{ r.feedback }}</p>
        </article>
      </div>

      <div class="results-actions">
        <button class="restart-btn" @click="resetAll">再来一次</button>
        <button class="share-btn" @click="copyResults">复制报告</button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { startInterview as apiStartInterview, submitAnswer, getSession, endSession } from "../api/interview";

const router = useRouter();
const chatRef = ref(null);

let sessionId = ref(null);
const phase = ref("setup");
const position = ref("");
const resumeText = ref("");
const totalQuestions = ref(5);
const starting = ref(false);
const submitting = ref(false);
const agentThinking = ref(false);
const messages = ref([]);
const rounds = ref([]);
const answerText = ref("");
const voiceText = ref("");
const recording = ref(false);

let recognition = null;
const voiceSupported = ref(false);

const currentRound = computed(() => {
  const answered = rounds.value.filter((r) => r.answer).length;
  return Math.min(answered + 1, totalQuestions.value);
});

const overallScore = computed(() => {
  const scored = rounds.value.filter((r) => r.score > 0);
  if (!scored.length) return 0;
  return Math.round(scored.reduce((s, r) => s + r.score, 0) / scored.length);
});

const highestScore = computed(() => {
  const scored = rounds.value.filter((r) => r.score > 0);
  if (!scored.length) return 0;
  return Math.max(...scored.map((r) => r.score));
});

const lowestScore = computed(() => {
  const scored = rounds.value.filter((r) => r.score > 0);
  if (!scored.length) return 0;
  return Math.min(...scored.map((r) => r.score));
});

function scoreClass(score) {
  if (score >= 80) return "score-high";
  if (score >= 60) return "score-mid";
  return "score-low";
}

function handleFileUpload(e) {
  const file = e.target.files?.[0];
  if (!file) return;
  if (file.size > 1024 * 1024) {
    alert("文件超过 1MB 限制，请粘贴文本内容");
    e.target.value = "";
    return;
  }
  const reader = new FileReader();
  reader.onload = () => {
    resumeText.value = reader.result;
  };
  reader.readAsText(file);
}

async function startInterview() {
  if (!position.value.trim()) return;
  starting.value = true;
  try {
    const { data } = await apiStartInterview({
      position: position.value.trim(),
      resumeText: resumeText.value.trim(),
      totalQuestions: totalQuestions.value
    });
    sessionId.value = data.item.id;
    rounds.value = data.item.rounds || [];
    messages.value = data.item.rounds.map((r) => ({
      key: `q-${r.id}`,
      role: "agent",
      text: r.question
    }));
    phase.value = "interview";
    speakText(data.item.rounds[0]?.question || "");
  } catch (err) {
    alert(err?.response?.data?.message || "启动面试失败");
  } finally {
    starting.value = false;
  }
}

async function sendAnswer() {
  const text = answerText.value.trim();
  if (!text || submitting.value || !sessionId.value) return;
  submitting.value = true;
  answerText.value = "";
  voiceText.value = "";

  // Track exact index so score matching is position-based, not text-based
  const userMsgIndex = messages.value.length;
  messages.value.push({ key: `a-${Date.now()}`, role: "user", text });
  agentThinking.value = true;
  scrollChat();

  try {
    const { data } = await submitAnswer(sessionId.value, { answer: text });
    rounds.value = data.item.rounds || [];

    if (data.completed) {
      const lastRound = rounds.value[rounds.value.length - 1];
      if (lastRound) {
        // Replace temp message in-place with full round data (includes score/feedback)
        messages.value[userMsgIndex] = {
          key: `a-${lastRound.id}`,
          role: "user",
          text: lastRound.answer,
          score: lastRound.score,
          feedback: lastRound.feedback
        };
      } else {
        messages.value.splice(userMsgIndex, 1);
      }
      agentThinking.value = false;
      phase.value = "results";
      return;
    }

    const newRound = rounds.value[rounds.value.length - 1];
    if (newRound && !newRound.answer) {
      messages.value.push({
        key: `q-${newRound.id}`,
        role: "agent",
        text: newRound.question
      });
      speakText(newRound.question);
    }

    // Attach score/feedback to the correct user message by index
    const lastAnswered = rounds.value[rounds.value.length - 2];
    if (lastAnswered && lastAnswered.score > 0) {
      const msg = messages.value[userMsgIndex];
      if (msg && msg.role === "user") {
        msg.score = lastAnswered.score;
        msg.feedback = lastAnswered.feedback;
      }
    }
  } catch (err) {
    messages.value.push({
      key: `err-${Date.now()}`,
      role: "agent",
      text: "抱歉，处理回答时出错，请重试。"
    });
  } finally {
    submitting.value = false;
    agentThinking.value = false;
    scrollChat();
  }
}

async function confirmEnd() {
  if (!confirm("确定要结束本次面试吗？已完成的问题仍会评分。")) return;
  if (!sessionId.value) return;
  try {
    await endSession(sessionId.value);
    const { data } = await getSession(sessionId.value);
    rounds.value = data.item.rounds || [];
    phase.value = "results";
  } catch {
    phase.value = "results";
  }
}

function resetAll() {
  sessionId.value = null;
  phase.value = "setup";
  position.value = "";
  resumeText.value = "";
  messages.value = [];
  rounds.value = [];
  answerText.value = "";
  voiceText.value = "";
}

async function copyResults() {
  const lines = [
    `🤖 AI 面试报告 - ${position.value}`,
    `综合评分：${overallScore.value}/100`,
    `共 ${rounds.value.length} 题\n`,
    ...rounds.value.map((r, i) =>
      `第 ${i + 1} 题 (${r.score}分)\n问题：${r.question}\n回答：${r.answer}${r.feedback ? `\n反馈：${r.feedback}` : ""}\n`
    )
  ];
  try {
    await navigator.clipboard.writeText(lines.join("\n"));
    alert("报告已复制到剪贴板");
  } catch {
    alert("复制失败，请手动复制");
  }
}

// ── Speech ──────────────────────────────────────────────────────────────

function initSpeech() {
  const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
  voiceSupported.value = Boolean(SpeechRecognition);
  if (!SpeechRecognition) return;

  recognition = new SpeechRecognition();
  recognition.lang = "zh-CN";
  recognition.continuous = false;
  recognition.interimResults = true;

  recognition.onresult = (event) => {
    let transcript = "";
    for (let i = event.resultIndex; i < event.results.length; i++) {
      transcript += event.results[i][0].transcript;
    }
    voiceText.value = transcript;
    answerText.value = transcript;
  };

  recognition.onerror = () => {
    recording.value = false;
  };

  recognition.onend = () => {
    recording.value = false;
  };
}

function toggleVoice() {
  if (!recognition) return;
  if (recording.value) {
    recognition.stop();
    recording.value = false;
  } else {
    voiceText.value = "";
    recording.value = true;
    recognition.start();
  }
}

function speakText(text) {
  if (!window.speechSynthesis) return;
  window.speechSynthesis.cancel();
  const utterance = new SpeechSynthesisUtterance(text);
  utterance.lang = "zh-CN";
  utterance.rate = 1.0;
  utterance.pitch = 1.0;
  window.speechSynthesis.speak(utterance);
}

function scrollChat() {
  nextTick(() => {
    if (chatRef.value) {
      chatRef.value.scrollTop = chatRef.value.scrollHeight;
    }
  });
}

watch(phase, () => {
  if (phase.value === "interview") {
    nextTick(scrollChat);
  }
});

// Also watch messages for scrolling
watch(messages, () => {
  nextTick(scrollChat);
}, { deep: true });

initSpeech();
onBeforeUnmount(() => {
  if (recognition) {
    recognition.stop();
  }
  window.speechSynthesis?.cancel();
});
</script>

<style scoped>
.interview-page {
  min-height: 100vh;
  padding: 28px;
  max-width: 800px;
  margin: 0 auto;
}

.page-label {
  color: var(--accent, #ffd98e);
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-size: 0.76rem;
  margin: 0;
}

/* ── Setup ── */
.setup-hero {
  text-align: center;
  margin-bottom: 32px;
}

.setup-hero h1 {
  margin: 10px 0 12px;
  font-size: clamp(2rem, 4vw, 3.2rem);
  line-height: 1;
}

.setup-desc {
  color: var(--soft, rgba(242, 239, 232, 0.7));
  max-width: 520px;
  margin: 0 auto;
}

.setup-form {
  max-width: 520px;
  margin: 0 auto;
  display: grid;
  gap: 20px;
}

.setup-field {
  display: grid;
  gap: 8px;
}

.setup-field span {
  font-weight: 600;
}

.field-hint {
  font-weight: 400;
  font-size: 0.85rem;
  color: var(--soft);
}

.setup-input,
.setup-textarea {
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid var(--border, rgba(255, 255, 255, 0.12));
  background: rgba(10, 14, 19, 0.7);
  color: inherit;
  font: inherit;
  resize: vertical;
}

.file-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.file-btn input {
  display: none;
}

.file-btn span {
  padding: 10px 14px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.05);
  font-size: 0.88rem;
  font-weight: 400;
}

.qty-selector {
  display: flex;
  gap: 10px;
}

.qty-btn {
  flex: 1;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.04);
  color: inherit;
  cursor: pointer;
  font: inherit;
  transition: all 0.2s;
}

.qty-btn.active {
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
  border-color: transparent;
}

.start-btn {
  padding: 16px;
  border-radius: 18px;
  border: none;
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
  font-size: 1.1rem;
  cursor: pointer;
}

.start-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ── Chat ── */
.interview-chat {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 56px);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 16px;
}

.chat-header h2 {
  margin: 6px 0 4px;
  font-size: 1.3rem;
}

.chat-progress {
  margin: 0;
  color: var(--soft);
  font-size: 0.9rem;
}

.end-btn {
  padding: 10px 16px;
  border-radius: 14px;
  border: 1px solid rgba(255, 139, 139, 0.3);
  background: rgba(255, 139, 139, 0.12);
  color: #ffbcbc;
  cursor: pointer;
  font: inherit;
  font-size: 0.88rem;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  display: grid;
  gap: 16px;
  padding: 8px 0;
}

.chat-msg {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.chat-msg.user {
  flex-direction: row-reverse;
}

.msg-avatar {
  font-size: 1.5rem;
  flex-shrink: 0;
  width: 36px;
  text-align: center;
}

.msg-bubble {
  max-width: min(480px, 80%);
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid var(--border);
  background: var(--panel, rgba(255, 255, 255, 0.06));
}

.chat-msg.user .msg-bubble {
  background: linear-gradient(135deg, rgba(255, 217, 142, 0.12), rgba(255, 166, 77, 0.08));
}

.msg-text {
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.msg-score {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--border);
  display: grid;
  gap: 6px;
}

.score-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 0.82rem;
  font-weight: 700;
}

.score-high {
  background: rgba(131, 242, 143, 0.16);
  color: #c8f8cd;
}

.score-mid {
  background: rgba(255, 217, 142, 0.16);
  color: #ffd98e;
}

.score-low {
  background: rgba(255, 139, 139, 0.16);
  color: #ffd2d2;
}

.score-feedback {
  font-size: 0.88rem;
  color: var(--soft);
}

.msg-thinking {
  display: flex;
  gap: 4px;
  padding: 18px 24px;
}

.thinking-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--soft);
  animation: think 1.4s infinite;
}

.thinking-dot:nth-child(2) { animation-delay: 0.2s; }
.thinking-dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes think {
  0%, 80%, 100% { opacity: 0.3; transform: scale(0.8); }
  40% { opacity: 1; transform: scale(1); }
}

.chat-input-area {
  margin-top: 12px;
  display: grid;
  gap: 8px;
  border-top: 1px solid var(--border);
  padding-top: 12px;
}

.voice-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.voice-btn {
  padding: 10px 16px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.05);
  color: inherit;
  cursor: pointer;
  font: inherit;
  font-size: 0.88rem;
  transition: all 0.2s;
}

.voice-btn.recording {
  background: rgba(255, 139, 139, 0.2);
  border-color: rgba(255, 139, 139, 0.4);
  color: #ffbcbc;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(255, 139, 139, 0.3); }
  50% { box-shadow: 0 0 0 8px rgba(255, 139, 139, 0); }
}

.voice-preview {
  font-size: 0.85rem;
  color: var(--soft);
  font-style: italic;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.input-row {
  display: flex;
  gap: 10px;
}

.answer-input {
  flex: 1;
  padding: 12px 14px;
  border-radius: 16px;
  border: 1px solid var(--border);
  background: rgba(10, 14, 19, 0.7);
  color: inherit;
  font: inherit;
  resize: none;
}

.send-btn {
  padding: 12px 20px;
  border-radius: 16px;
  border: none;
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ── Results ── */
.results-hero {
  text-align: center;
  margin-bottom: 28px;
}

.results-hero h1 {
  margin: 10px 0 8px;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
}

.results-subtitle {
  color: var(--soft);
  margin: 0;
}

.results-summary {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 28px;
}

.summary-stat {
  padding: 20px;
  border-radius: 20px;
  border: 1px solid var(--border);
  background: var(--panel);
  text-align: center;
  display: grid;
  gap: 6px;
}

.stat-value {
  font-size: 2.2rem;
  font-weight: 700;
}

.stat-label {
  font-size: 0.85rem;
  color: var(--soft);
}

.results-detail {
  display: grid;
  gap: 16px;
}

.result-card {
  padding: 20px;
  border-radius: 22px;
  border: 1px solid var(--border);
  background: var(--panel);
}

.result-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.result-num {
  font-size: 0.85rem;
  color: var(--soft);
}

.result-card p {
  margin: 8px 0;
  line-height: 1.6;
}

.result-q strong,
.result-a strong,
.result-feedback strong {
  font-size: 0.85rem;
  color: var(--soft);
}

.result-feedback {
  color: var(--soft);
  font-size: 0.92rem;
}

.results-actions {
  margin-top: 28px;
  display: flex;
  gap: 12px;
  justify-content: center;
}

.restart-btn,
.share-btn {
  padding: 14px 24px;
  border-radius: 16px;
  border: 1px solid var(--border);
  font: inherit;
  cursor: pointer;
  font-weight: 600;
}

.restart-btn {
  background: linear-gradient(135deg, #ffd98e, #ffa64d);
  color: #24170c;
  border: none;
}

.share-btn {
  background: rgba(255, 255, 255, 0.05);
  color: inherit;
}

@media (max-width: 600px) {
  .interview-page {
    padding: 16px;
  }
  .results-summary {
    grid-template-columns: repeat(2, 1fr);
  }
  .chat-msg .msg-bubble {
    max-width: 85%;
  }
}
</style>
