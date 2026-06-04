<template>
  <section class="content-section">
    <div class="section-head">
      <div>
        <p class="eyebrow">个人库</p>
        <h3>我的收藏与点赞</h3>
      </div>
    </div>

    <div class="collections-grid">
      <section class="panel-card collection-panel">
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

      <section class="panel-card collection-panel">
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
.empty-panel--compact {
  min-height: 220px;
}

@media (max-width: 768px) {
  .collection-panel {
    padding: 18px;
  }

  .collections-grid {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .collection-panel {
    padding: 14px;
  }
}
</style>
