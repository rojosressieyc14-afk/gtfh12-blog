<template>
  <section class="uc-collections">
    <section class="uc-hero uc-hero--collections">
      <div class="uc-hero__copy">
        <p class="eyebrow uc-hero__eyebrow">个人库</p>
        <h2 class="uc-hero__title">我的收藏与点赞</h2>
        <p class="uc-hero__text">收藏值得反复参考的内容，点赞标记当下认同的文章。</p>
      </div>
    </section>

    <div class="collections-grid">
      <section class="collection-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">收藏</p>
            <h3>收藏夹</h3>
          </div>
        </div>
        <div v-if="favorites.length" class="article-grid article-grid--single">
          <ArticleCard v-for="item in favorites" :key="`fav-${item.id}`" :item="item" />
        </div>
        <div v-else class="empty-panel empty-panel--compact">
          <h4>你还没有收藏文章</h4>
          <p>看到值得反复参考的内容，可以先收藏起来，后续查阅会更方便。</p>
        </div>
      </section>

      <section class="collection-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">点赞</p>
            <h3>我的点赞</h3>
          </div>
        </div>
        <div v-if="likes.length" class="article-grid article-grid--single">
          <ArticleCard v-for="item in likes" :key="`like-${item.id}`" :item="item" />
        </div>
        <div v-else class="empty-panel empty-panel--compact">
          <h4>你还没有点赞文章</h4>
          <p>点赞适合标记当下认同或觉得有价值的内容，方便后面回看。</p>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { listMyFavorites, listMyLikes } from "../api/article";
import ArticleCard from "../components/ArticleCard.vue";

const favorites = ref([]);
const likes = ref([]);

async function loadData() {
  const [favoritesRes, likesRes] = await Promise.all([listMyFavorites(), listMyLikes()]);
  favorites.value = favoritesRes.data.items || [];
  likes.value = likesRes.data.items || [];
}

onMounted(loadData);
</script>

<style scoped>
.uc-collections {
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

.collections-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.collection-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 24px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(255, 255, 255, 0.5);
}

.collection-panel :deep(.article-card:first-child) {
  background:
    radial-gradient(120% 120% at 80% 0%, rgba(255, 138, 76, 0.06), transparent 45%),
    radial-gradient(100% 100% at 0% 100%, rgba(255, 209, 102, 0.04), transparent 50%),
    #f4f4f6;
  color: #1d1d1f;
  border-color: rgba(0, 0, 0, 0.08);
}

.collection-panel :deep(.article-card:first-child::before) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), transparent 44%, rgba(255, 209, 102, 0.04));
}

.collection-panel :deep(.article-card:first-child h3) {
  color: #1d1d1f;
}

.collection-panel :deep(.article-card:first-child p) {
  color: #48484a;
}

.collection-panel :deep(.article-card:first-child .article-card__meta),
.collection-panel :deep(.article-card:first-child footer) {
  color: #6b7280;
}

.collection-panel :deep(.article-card:first-child .tag-chip) {
  background: rgba(180, 83, 10, 0.1);
  border-color: rgba(180, 83, 10, 0.2);
  color: #92400e;
}

.collection-panel :deep(.article-card:first-child .status-chip) {
  background: rgba(0, 0, 0, 0.06);
  color: #374151;
}

.collection-panel :deep(.article-card:first-child .status-chip.published) {
  background: rgba(22, 163, 74, 0.12);
  color: #166534;
}

.collection-panel :deep(.article-card:first-child .article-card__shine) {
  background: radial-gradient(circle, rgba(255, 138, 76, 0.15), transparent 68%);
}

.empty-panel--compact {
  min-height: 220px;
}

@media (max-width: 768px) {
  .collections-grid {
    grid-template-columns: 1fr;
  }

  .collection-panel {
    padding: 18px;
  }
}

@media (max-width: 480px) {
  .collection-panel {
    padding: 14px;
  }
}
</style>
