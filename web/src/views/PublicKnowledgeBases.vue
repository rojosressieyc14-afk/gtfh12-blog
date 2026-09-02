<template>
  <section class="content-section kb-public-page">
    <section class="kb-public-hero">
      <div class="kb-public-hero__copy">
        <p class="eyebrow kb-public-hero__eyebrow">知识库</p>
        <h2 class="kb-public-hero__title">结构化的知识沉淀与公开分享空间。</h2>
        <p class="kb-public-hero__text">
          浏览公开的知识库文档，按主题、标签和分类快速进入内容。知识库支持双向链接和图谱视图，帮助你发现内容之间的关联。
        </p>
      </div>
      <div class="kb-public-hero__stats">
        <article class="kb-public-hero__stat">
          <strong>{{ publicKbs.length }}</strong>
          <span>公开知识库</span>
        </article>
        <article class="kb-public-hero__stat">
          <strong>{{ totalDocs }}</strong>
          <span>公开文档</span>
        </article>
      </div>
    </section>

    <section v-if="myKbs.length" class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">我的知识库</p>
        </div>
        <router-link class="ghost-btn" to="/user-center/knowledge-base">管理知识库</router-link>
      </div>
      <div class="kb-public-grid">
        <article
          v-for="kb in myKbs"
          :key="`my-${kb.id}`"
          class="kb-public-card kb-public-card--mine"
          @click="$router.push({ name: 'uc-knowledge-base-detail', params: { id: kb.id } })"
        >
          <div class="kb-public-card__head">
            <div>
              <h3>{{ kb.name }}</h3>
              <p class="table-note">{{ kb.docCount }} 篇文档 · {{ formatDate(kb.updatedAt) }}</p>
            </div>
            <span class="tag-chip tag-chip--accent">我的</span>
          </div>
          <p v-if="kb.description" class="detail-summary">{{ kb.description }}</p>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">公开知识库</p>
        </div>
      </div>

      <div v-if="loading" class="empty-panel">
        <h4>加载中...</h4>
      </div>
      <div v-else-if="!publicKbs.length" class="empty-panel">
        <h4>暂无公开知识库</h4>
        <p>还没有公开的知识库文档。</p>
      </div>
      <div v-else class="kb-public-grid">
        <article
          v-for="kb in publicKbs"
          :key="`pub-${kb.id}`"
          class="kb-public-card"
          @click="$router.push(`/knowledge-bases/${kb.id}`)"
        >
          <div class="kb-public-card__head">
            <div>
              <h3>{{ kb.name }}</h3>
              <p class="table-note">{{ kb.userName }} · {{ kb.publicCount }} 篇公开文档</p>
            </div>
            <span class="tag-chip">{{ kb.docCount }} 篇</span>
          </div>
          <p v-if="kb.description" class="detail-summary">{{ kb.description }}</p>
          <footer class="kb-public-card__footer">
            <span class="table-note">更新于 {{ formatDate(kb.updatedAt) }}</span>
            <router-link class="ghost-btn" :to="`/knowledge-bases/${kb.id}`">浏览文档</router-link>
          </footer>
        </article>
      </div>
    </section>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useUserStore } from "../stores/user";
import { listPublicKBs, listKnowledgeBases } from "../api/knowledgeBase";
import { formatDate } from "../utils/date";

const route = useRoute();
const userStore = useUserStore();
const myKbs = ref([]);
const publicKbs = ref([]);
const loading = ref(true);

const totalDocs = computed(() => publicKbs.value.reduce((sum, kb) => sum + (kb.publicCount || 0), 0));

async function loadData() {
  loading.value = true;
  try {
    const tasks = [listPublicKBs()];
    if (userStore.isLoggedIn) tasks.push(listKnowledgeBases());
    const results = await Promise.all(tasks);
    publicKbs.value = results[0].data.items || [];
    if (userStore.isLoggedIn && results[1]) myKbs.value = results[1].data.items || [];
  } catch { publicKbs.value = []; myKbs.value = []; } finally { loading.value = false; }
}

watch(() => route.fullPath, loadData);
onMounted(loadData);
</script>

<style scoped>
.kb-public-page {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.kb-public-hero {
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

.kb-public-hero__copy {
  display: grid;
  gap: 12px;
}

.kb-public-hero__eyebrow {
  color: var(--accent);
  font-weight: 600;
}

.kb-public-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
}

.kb-public-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: var(--text-soft);
}

.kb-public-hero__stats {
  display: flex;
  gap: 24px;
}

.kb-public-hero__stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 24px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid var(--border);
}

.kb-public-hero__stat strong {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--accent);
}

.kb-public-hero__stat span {
  font-size: 0.8rem;
  color: var(--text-soft);
  margin-top: 4px;
}

.kb-public-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.kb-public-card {
  padding: 22px;
  border-radius: 28px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: transform 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
}

.kb-public-card:hover {
  transform: translateY(-5px);
  border-color: rgba(255, 209, 102, 0.45);
  box-shadow: 0 24px 50px rgba(0, 0, 0, 0.22);
}

.kb-public-card--mine {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.14), transparent 45%),
    rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 209, 102, 0.2);
}

.kb-public-card__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.kb-public-card__head h3 {
  margin: 0;
  font-size: 1.1rem;
}

.kb-public-card__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.kb-public-card__footer .ghost-btn {
  color: var(--accent);
}

.kb-public-card__footer .ghost-btn:hover {
  background: rgba(255, 138, 76, 0.12);
  border-color: var(--accent);
}

.tag-chip--accent {
  background: rgba(255, 138, 76, 0.15);
  color: var(--accent);
}

.empty-panel {
  text-align: center;
  padding: 60px 20px;
}

.empty-panel h4 {
  margin: 0 0 8px;
}

.empty-panel p {
  margin: 0;
  color: var(--text-soft);
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .kb-public-hero {
    grid-template-columns: 1fr;
    padding: 28px 20px;
  }

  .kb-public-grid {
    grid-template-columns: 1fr;
  }
}
</style>
