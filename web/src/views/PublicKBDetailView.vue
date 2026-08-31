<template>
  <section class="content-section kb-pub-detail-page">
    <div v-if="loading" class="empty-panel">
      <h4>加载中...</h4>
    </div>
    <div v-else-if="!kb" class="empty-panel">
      <h4>知识库不存在</h4>
      <router-link class="ghost-btn" to="/knowledge-bases">返回知识库列表</router-link>
    </div>
    <template v-else>
      <section class="kb-pub-detail-hero">
        <div class="kb-pub-detail-hero__copy">
          <router-link class="kb-pub-back" to="/knowledge-bases">&larr; 知识库列表</router-link>
          <p class="eyebrow kb-pub-hero__eyebrow">知识库</p>
          <h2 class="kb-pub-hero__title">{{ kb.name }}</h2>
          <p v-if="kb.description" class="kb-pub-hero__text">{{ kb.description }}</p>
          <p class="kb-pub-hero__meta">{{ docs.length }} 篇公开文档</p>
        </div>
      </section>

      <section class="content-section content-section--compact">
        <div v-if="!docs.length" class="empty-panel">
          <h4>该知识库暂无公开文档</h4>
        </div>
        <div v-else class="kb-pub-doc-list">
          <article
            v-for="doc in docs"
            :key="doc.id"
            class="kb-pub-doc-card"
            @click="openDoc(doc)"
          >
            <div class="kb-pub-doc-card__meta">
              <span v-if="doc.category" class="tag-chip">{{ doc.category.name }}</span>
              <span class="table-note">{{ formatDate(doc.createdAt) }}</span>
              <span class="table-note">{{ doc.viewCount }} 次阅读</span>
            </div>
            <h3>{{ doc.title || '无标题' }}</h3>
            <p class="detail-summary">{{ getPreview(doc.content) }}</p>
            <div v-if="doc.tags?.length" class="tag-row">
              <span v-for="tag in doc.tags" :key="tag.id" class="tag-chip"># {{ tag.name }}</span>
            </div>
          </article>
        </div>
      </section>
    </template>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listPublicDocuments, getKnowledgeBase } from "../api/knowledgeBase";

const route = useRoute();
const router = useRouter();
const kb = ref(null);
const docs = ref([]);
const loading = ref(true);

function formatDate(value) {
  return new Date(value).toLocaleDateString("zh-CN");
}

function getPreview(content) {
  if (!content) return "";
  const plain = content.replace(/[#*_`>\[\]()]/g, "").replace(/\n+/g, " ").trim();
  return plain.length > 160 ? plain.slice(0, 160) + "..." : plain;
}

function openDoc(doc) {
  router.push(`/kb-note/${doc.id}`);
}

onMounted(async () => {
  try {
    const [kbRes, docRes] = await Promise.all([
      getKnowledgeBase(route.params.id),
      listPublicDocuments(route.params.id),
    ]);
    kb.value = kbRes.data.item;
    docs.value = docRes.data.items || [];
  } catch {
    kb.value = null;
    docs.value = [];
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.kb-pub-detail-page {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.kb-pub-detail-hero {
  border-radius: 28px;
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.18), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.14), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  padding: clamp(32px, 6vw, 56px) clamp(24px, 4vw, 48px);
}

.kb-pub-detail-hero__copy {
  display: grid;
  gap: 8px;
}

.kb-pub-back {
  font-size: 0.85rem;
  color: #b4530a;
  text-decoration: none;
  margin-bottom: 4px;
}
.kb-pub-back:hover {
  text-decoration: underline;
}

.kb-pub-hero__eyebrow {
  color: #b4530a;
  font-weight: 600;
}

.kb-pub-hero__title {
  margin: 0;
  font-size: clamp(1.6rem, 3vw, 2.4rem);
  line-height: 1.12;
  color: #1d1d1f;
}

.kb-pub-hero__text {
  margin: 0;
  font-size: clamp(0.9rem, 1.3vw, 1.05rem);
  line-height: 1.6;
  color: #48484a;
}

.kb-pub-hero__meta {
  margin: 0;
  font-size: 0.85rem;
  color: #6b7280;
}

.kb-pub-doc-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.kb-pub-doc-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.kb-pub-doc-card:hover {
  border-color: rgba(249, 115, 22, 0.25);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.05);
}

.kb-pub-doc-card__meta {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.kb-pub-doc-card h3 {
  margin: 0 0 8px;
  font-size: 1.05rem;
  color: #1d1d1f;
}

.empty-panel {
  text-align: center;
  padding: 60px 20px;
}

.empty-panel h4 {
  margin: 0 0 12px;
  color: #6b7280;
}

@media (max-width: 768px) {
  .kb-pub-doc-list {
    grid-template-columns: 1fr;
  }
}
</style>
