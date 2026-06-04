<template>
  <section class="profile-page">
    <div class="profile-layout">
      <article class="panel-card profile-main-card">
        <p class="eyebrow">个人资料工作台</p>
        <h2>把你的个人页面整理成真正能对外展示的品牌名片。</h2>
        <p class="detail-summary">
          这些字段会影响首页、作者页和关于页的展示效果，也会决定别人第一次看到你时留下什么印象。
        </p>

        <section class="profile-progress-card">
          <div>
            <strong>{{ completionRate }}%</strong>
            <p>资料完整度</p>
          </div>
          <div class="profile-progress-bar">
            <span :style="{ width: `${completionRate}%` }"></span>
          </div>
          <p class="detail-summary">把身份、联系方式和作品集相关字段补齐后，整个站点会更像一份完整的线上简历。</p>
        </section>

        <section class="profile-preset-grid">
          <button class="preset-card" type="button" @click="applyPreset('frontend')">
            <span class="eyebrow">模板</span>
            <strong>前端开发者</strong>
            <p>更适合 UI 工程、交互设计和界面细节打磨型作品。</p>
          </button>
          <button class="preset-card" type="button" @click="applyPreset('fullstack')">
            <span class="eyebrow">模板</span>
            <strong>全栈开发者</strong>
            <p>适合独立项目、端到端交付和完整产品落地展示。</p>
          </button>
          <button class="preset-card" type="button" @click="applyPreset('student')">
            <span class="eyebrow">模板</span>
            <strong>学习成长型</strong>
            <p>适合把站点做成作品集、成长记录和学习归档的结合体。</p>
          </button>
        </section>

        <form class="stack-form" @submit.prevent="saveProfile">
          <section class="profile-avatar-block">
            <div class="profile-avatar-shell">
              <img v-if="form.avatar" :src="toAssetUrl(form.avatar)" alt="avatar" class="profile-avatar-image" />
              <span v-else>{{ avatarFallback }}</span>
            </div>
            <div class="profile-avatar-actions">
              <strong>头像与品牌卡片</strong>
              <p>上传一张清晰头像或个人品牌图片，用于资料页和作者页展示。</p>
              <label class="ghost-btn profile-upload-btn">
                <input type="file" accept="image/*" hidden @change="handleAvatarUpload" />
                {{ uploading ? "上传中..." : "上传头像" }}
              </label>
            </div>
          </section>

          <div class="profile-form-grid">
            <label>
              标题文案
              <input v-model.trim="form.headline" class="field-input" placeholder="例如：专注于产品体验和工程实现的前端开发者" />
            </label>

            <label>
              当前身份
              <input v-model.trim="form.currentRole" class="field-input" placeholder="例如：前端工程师 / 全栈开发者 / 学生开发者" />
            </label>

            <label>
              经历标签
              <input v-model.trim="form.yearsLabel" class="field-input" placeholder="例如：3 年经验 / 2024 至今 / 校园开发者" />
            </label>

            <label>
              所在地
              <input v-model.trim="form.location" class="field-input" placeholder="例如：上海 / 杭州 / 远程" />
            </label>

            <label>
              邮箱
              <input v-model.trim="form.email" class="field-input" placeholder="hello@example.com" />
            </label>

            <label>
              简历链接
              <input v-model.trim="form.resumeUrl" class="field-input" placeholder="https://your-site.com/resume.pdf" />
            </label>

            <label>
              个人网站
              <input v-model.trim="form.websiteUrl" class="field-input" placeholder="https://your-site.com" />
            </label>

            <label>
              GitHub
              <input v-model.trim="form.githubUrl" class="field-input" placeholder="https://github.com/your-name" />
            </label>

            <label>
              Gitee
              <input v-model.trim="form.giteeUrl" class="field-input" placeholder="https://gitee.com/your-name" />
            </label>

            <label>
              Juejin
              <input v-model.trim="form.juejinUrl" class="field-input" placeholder="https://juejin.cn/user/xxx" />
            </label>

            <label>
              CSDN
              <input v-model.trim="form.csdnUrl" class="field-input" placeholder="https://blog.csdn.net/xxx" />
            </label>
          </div>

          <label>
            技能栏
            <input v-model.trim="skillsText" class="field-input" placeholder="Vue 3, Go, Gin, MySQL, UI 设计" />
          </label>

          <div v-if="parsedSkills.length" class="tag-row">
            <span v-for="skill in parsedSkills" :key="skill" class="tag-chip"># {{ skill }}</span>
          </div>

          <label>
            关注方向
            <input v-model.trim="focusAreasText" class="field-input" placeholder="前端架构、界面系统、Go 后端、AI 应用" />
          </label>

          <div v-if="parsedFocusAreas.length" class="tag-row">
            <span v-for="item in parsedFocusAreas" :key="item" class="tag-chip">{{ item }}</span>
          </div>

          <label>
            个性签名
            <input v-model.trim="form.motto" class="field-input" placeholder="写一句能代表你的短句。" />
          </label>

          <label>
            个人简介
            <textarea
              v-model.trim="form.bio"
              class="field-area"
              placeholder="介绍你的方向、优势、学习路径，以及你希望别人怎样记住你。"
            ></textarea>
          </label>

          <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
        <div class="inline-actions">
            <button class="solid-btn profile-btn" :disabled="saving || uploading">
              {{ saving ? "保存中..." : "保存资料" }}
            </button>
            <router-link class="ghost-btn profile-link" :to="`/author/${userStore.profile?.id}`">预览作者页</router-link>
            <router-link class="ghost-btn profile-link" to="/about">打开关于页</router-link>
          </div>
        </form>
      </article>

      <aside class="panel-card profile-side-card">
        <p class="eyebrow">预览</p>
        <div class="brand-preview">
          <div class="brand-preview__avatar">
            <img v-if="form.avatar" :src="toAssetUrl(form.avatar)" alt="avatar" />
            <span v-else>{{ avatarFallback }}</span>
          </div>
          <h3>{{ userStore.profile?.username }}</h3>
          <p v-if="form.currentRole" class="brand-preview__meta">{{ form.currentRole }}</p>
          <p class="brand-preview__headline">{{ form.headline || "写一句清晰的话，定义你的方向。" }}</p>
          <p v-if="form.yearsLabel" class="brand-preview__meta">{{ form.yearsLabel }}</p>
          <p v-if="form.motto" class="brand-preview__motto">"{{ form.motto }}"</p>
          <p class="brand-preview__bio">{{ form.bio || "你的简短介绍会显示在这里。" }}</p>

          <div v-if="parsedSkills.length" class="tag-row">
            <span v-for="skill in parsedSkills.slice(0, 6)" :key="skill" class="tag-chip">{{ skill }}</span>
          </div>

          <div class="profile-meta">
            <span>角色：{{ userRoleText }}</span>
            <span>状态：{{ userStatusText }}</span>
          </div>
        </div>

        <section class="profile-side-section">
          <p class="eyebrow">公开链接</p>
          <div class="profile-link-list">
            <a v-if="form.websiteUrl" class="inline-link" :href="form.websiteUrl" target="_blank" rel="noreferrer">网站</a>
            <a v-if="form.resumeUrl" class="inline-link" :href="form.resumeUrl" target="_blank" rel="noreferrer">简历</a>
            <a v-if="form.githubUrl" class="inline-link" :href="form.githubUrl" target="_blank" rel="noreferrer">GitHub</a>
            <a v-if="form.giteeUrl" class="inline-link" :href="form.giteeUrl" target="_blank" rel="noreferrer">Gitee</a>
            <a v-if="form.juejinUrl" class="inline-link" :href="form.juejinUrl" target="_blank" rel="noreferrer">Juejin</a>
            <a v-if="form.csdnUrl" class="inline-link" :href="form.csdnUrl" target="_blank" rel="noreferrer">CSDN</a>
            <span v-if="!hasPublicLinks" class="detail-summary">补几个公开链接后，这里会更像真正的作品集名片。</span>
          </div>
        </section>

        <section class="profile-side-section">
          <p class="eyebrow">检查清单</p>
          <div class="profile-checklist">
            <span :class="{ complete: Boolean(form.avatar) }">头像</span>
            <span :class="{ complete: Boolean(form.headline) }">标题文案</span>
            <span :class="{ complete: Boolean(form.currentRole) }">当前身份</span>
            <span :class="{ complete: Boolean(form.bio) }">简介</span>
            <span :class="{ complete: parsedSkills.length > 0 }">技能栏</span>
            <span :class="{ complete: parsedFocusAreas.length > 0 }">关注方向</span>
            <span :class="{ complete: Boolean(form.email || form.websiteUrl) }">联系方式</span>
          </div>
        </section>
      </aside>
    </div>
  </section>
</template>

<script setup>
import { computed, reactive, ref, watchEffect } from "vue";
import { uploadImage } from "../api/upload";
import { useUserStore } from "../stores/user";
import { toAssetUrl } from "../utils/asset";

const userStore = useUserStore();
const uploading = ref(false);
const saving = ref(false);
const errorMessage = ref("");
const skillsText = ref("");
const focusAreasText = ref("");
const form = reactive({
  avatar: "",
  headline: "",
  currentRole: "",
  yearsLabel: "",
  motto: "",
  location: "",
  email: "",
  resumeUrl: "",
  websiteUrl: "",
  githubUrl: "",
  giteeUrl: "",
  juejinUrl: "",
  csdnUrl: "",
  bio: ""
});

watchEffect(() => {
  const profile = userStore.profile || {};
  form.avatar = profile.avatar || "";
  form.headline = profile.headline || "";
  form.currentRole = profile.currentRole || "";
  form.yearsLabel = profile.yearsLabel || "";
  form.motto = profile.motto || "";
  form.location = profile.location || "";
  form.email = profile.email || "";
  form.resumeUrl = profile.resumeUrl || "";
  form.websiteUrl = profile.websiteUrl || "";
  form.githubUrl = profile.githubUrl || "";
  form.giteeUrl = profile.giteeUrl || "";
  form.juejinUrl = profile.juejinUrl || "";
  form.csdnUrl = profile.csdnUrl || "";
  form.bio = profile.bio || "";
  skillsText.value = Array.isArray(profile.skills) ? profile.skills.join(", ") : "";
  focusAreasText.value = Array.isArray(profile.focusAreas) ? profile.focusAreas.join(", ") : "";
});

const parsedSkills = computed(() =>
  skillsText.value
    .split(/[,\n]/)
    .map((item) => item.trim())
    .filter((item, index, list) => item && list.indexOf(item) === index)
);

const parsedFocusAreas = computed(() =>
  focusAreasText.value
    .split(/[,\n]/)
    .map((item) => item.trim())
    .filter((item, index, list) => item && list.indexOf(item) === index)
);

const avatarFallback = computed(() => (userStore.profile?.username || "P").slice(0, 1).toUpperCase());
const userRoleText = computed(() => (userStore.profile?.role === "admin" ? "管理员" : "普通用户"));
const userStatusText = computed(() => (userStore.profile?.status === "banned" ? "已封禁" : "正常"));
const hasPublicLinks = computed(() =>
  Boolean(form.websiteUrl || form.resumeUrl || form.githubUrl || form.giteeUrl || form.juejinUrl || form.csdnUrl)
);
const completionRate = computed(() => {
  const checks = [
    form.avatar,
    form.headline,
    form.currentRole,
    form.yearsLabel,
    form.location,
    form.email || form.websiteUrl,
    form.motto,
    form.bio,
    parsedSkills.value.length,
    parsedFocusAreas.value.length
  ];
  const done = checks.filter(Boolean).length;
  return Math.round((done / checks.length) * 100);
});

async function handleAvatarUpload(event) {
  const [file] = event.target.files || [];
  if (!file) return;

  uploading.value = true;
  try {
    const { data } = await uploadImage(file);
    form.avatar = data.url;
  } finally {
    uploading.value = false;
    event.target.value = "";
  }
}

async function saveProfile() {
  saving.value = true;
  errorMessage.value = "";
  try {
    await userStore.updateProfileAction({
      ...form,
      skills: parsedSkills.value,
      focusAreas: parsedFocusAreas.value
    });
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || "保存失败，请稍后再试";
  } finally {
    saving.value = false;
  }
}

function applyPreset(kind) {
  if (kind === "frontend") {
    form.headline = form.headline || "专注于界面质量、交互体验和组件工程的前端开发者。";
    form.currentRole = form.currentRole || "前端工程师";
    form.yearsLabel = form.yearsLabel || "持续通过项目打磨界面与体验能力";
    form.motto = form.motto || "先想清楚，再把体验做细。";
    form.bio =
      form.bio ||
      "我更关注把想法落成清晰、可用、耐看的界面，同时不断积累组件化、交互设计和前端工程能力。";
    skillsText.value = skillsText.value || "Vue 3, TypeScript, 动效设计, 响应式界面, 组件系统";
    focusAreasText.value = focusAreasText.value || "前端架构, 交互设计, 设计系统";
    return;
  }

  if (kind === "fullstack") {
    form.headline = form.headline || "能够从界面到后端完整落地产品的全栈开发者。";
    form.currentRole = form.currentRole || "全栈开发者";
    form.yearsLabel = form.yearsLabel || "持续交付端到端项目";
    form.motto = form.motto || "做有用的东西，快速上线，持续迭代。";
    form.bio =
      form.bio ||
      "我喜欢完整负责一个产品闭环，从界面设计、后端实现、数据建模到上线后的持续迭代，都会主动参与。";
    skillsText.value = skillsText.value || "Vue 3, Go, Gin, MySQL, REST API, 产品交付";
    focusAreasText.value = focusAreasText.value || "全栈应用, 个人项目, 产品工程";
    return;
  }

  form.headline = form.headline || "通过项目、笔记和公开输出持续成长的开发者。";
  form.currentRole = form.currentRole || "成长型开发者";
  form.yearsLabel = form.yearsLabel || "在项目与写作中持续进步";
  form.motto = form.motto || "认真学习，清楚表达，持续构建。";
  form.bio =
    form.bio ||
    "这个站点会记录我的项目实践、技术笔记和成长轨迹，让学习与输出都变成能被看见的成果。";
  skillsText.value = skillsText.value || "Vue 3, Go, 笔记写作, 问题排查, 持续学习";
  focusAreasText.value = focusAreasText.value || "学习归档, 项目作品集, 公开写作";
}
</script>

<style scoped>
.profile-progress-card {
  display: grid;
  gap: 10px;
  margin: 18px 0 22px;
  padding: 18px;
  border-radius: 22px;
  border: 1px solid rgba(249, 115, 22, 0.2);
  background: rgba(249, 115, 22, 0.08);
}

.profile-progress-card strong {
  font-size: 2rem;
}

.profile-progress-card p {
  margin: 0;
}

.profile-progress-bar {
  height: 10px;
  border-radius: 999px;
  overflow: hidden;
  background: rgba(15, 23, 42, 0.16);
}

.profile-progress-bar span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #f97316, #fb7185);
}

.profile-preset-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 24px;
}

.preset-card {
  display: grid;
  gap: 8px;
  text-align: left;
  padding: 18px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 20px;
  background: rgba(15, 23, 42, 0.38);
  color: inherit;
  cursor: pointer;
  transition: transform 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.preset-card:hover {
  transform: translateY(-2px);
  border-color: rgba(249, 115, 22, 0.34);
  background: rgba(15, 23, 42, 0.52);
}

.preset-card p,
.preset-card strong {
  margin: 0;
}

.profile-side-section {
  margin-top: 24px;
  display: grid;
  gap: 10px;
}

.profile-link-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 14px;
}

.profile-checklist {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.profile-checklist span {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.12);
  color: rgba(226, 232, 240, 0.72);
}

.profile-checklist span.complete {
  background: rgba(249, 115, 22, 0.16);
  color: #fff7ed;
}

@media (max-width: 980px) {
  .profile-preset-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .profile-progress-card {
    padding: 14px;
  }

  .profile-progress-card strong {
    font-size: 1.6rem;
  }

  .profile-avatar-block {
    flex-direction: column;
    text-align: center;
    padding: 16px;
  }

  .profile-avatar-actions {
    align-items: center;
  }

  .profile-upload-btn {
    width: 100%;
    justify-content: center;
  }

  .preset-card {
    padding: 14px;
  }

  .profile-btn,
  .profile-link {
    width: 100%;
    justify-content: center;
  }

  .inline-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .profile-side-card {
    padding: 18px;
  }

  .profile-checklist span {
    font-size: 0.82rem;
    padding: 6px 10px;
  }

  .profile-link-list {
    gap: 8px;
  }
}
</style>
