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
          <strong>{{ kbs.length }}</strong>
          <span>公开知识库</span>
        </article>
        <article class="kb-public-hero__stat">
          <strong>{{ totalDocs }}</strong>
          <span>公开文档</span>
        </article>
      </div>
    </section>

    <section class="content-section content-section--compact">
      <div class="section-head">
        <div>
          <p class="eyebrow">知识库列表</p>
        </div>
      </div>

      <div v-if="loading" class="empty-panel">
        <h4>加载中...</h4>
      </div>
      <div v-else-if="!kbs.length" class="empty-panel">
        <h4>暂无公开知识库</h4>
        <p>作者还没有公开任何知识库文档。</p>
      </div>
      <div v-else class="kb-public-grid">
        <article
          v-for="kb in kbs"
          :key="kb.id"
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
import { computed, onMounted, ref } from "vue";
import { listPublicKBs } from "../api/knowledgeBase";

const kbs = ref([]);
const loading = ref(true);

const totalDocs = computed(() => kbs.value.reduce((sum, kb) => sum + (kb.publicCount || 0), 0));

function formatDate(value) {
  return new Date(value).toLocaleDateString("zh-CN");
}

onMounted(async () => {
  try {
    const { data } = await listPublicKBs();
    kbs.value = data.items || [];
  } catch {
    kbs.value = [];
  } finally {
    loading.value = false;
  }
});
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
    #f4f4f6;
  color: #1d1d1f;
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
  color: #b4530a;
  font-weight: 600;
}

.kb-public-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
  color: #1d1d1f;
}

.kb-public-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: #48484a;
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
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.kb-public-hero__stat strong {
  font-size: 1.8rem;
  font-weight: 700;
  color: #b4530a;
}

.kb-public-hero__stat span {
  font-size: 0.8rem;
  color: #6b7280;
  margin-top: 4px;
}

.kb-public-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.kb-public-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.kb-public-card:hover {
  border-color: rgba(249, 115, 22, 0.25);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.05);
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
  color: #1d1d1f;
}

.kb-public-card__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.kb-public-card__footer .ghost-btn {
  border-color: rgba(0, 0, 0, 0.12);
  color: #b4530a;
}

.kb-public-card__footer .ghost-btn:hover {
  background: rgba(255, 138, 76, 0.08);
  border-color: #b4530a;
}

.empty-panel {
  text-align: center;
  padding: 60px 20px;
}

.empty-panel h4 {
  margin: 0 0 8px;
  color: #6b7280;
}

.empty-panel p {
  margin: 0;
  color: #9ca3af;
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
