<template>
  <article class="article-card article-card--localized" @click="$router.push(`/article/${item.id}`)">
    <div v-if="coverUrl" class="article-cover" :style="{ backgroundImage: `url(${coverUrl})` }"></div>
    <div class="article-card__shine"></div>

    <div class="article-card__meta">
      <span class="status-chip" :class="item.status">{{ statusLabel }}</span>
      <span>{{ item.category?.name || "未分类" }}</span>
    </div>

    <h3>{{ item.title }}</h3>
    <p>{{ item.summary || "这篇文章还没有摘要，但标题已经概括了主要内容。" }}</p>

    <div v-if="item.tags?.length" class="tag-row">
      <span v-for="tag in item.tags.slice(0, 3)" :key="tag.id || tag.name" class="tag-chip"># {{ tag.name }}</span>
    </div>

    <footer>
      <span>{{ item.author?.username || "匿名作者" }} · {{ formatDate(item.publishedAt || item.createdAt) }}</span>
      <span class="article-card__more">阅读全文</span>
    </footer>
  </article>
</template>

<script setup>
import { computed } from "vue";
import { toAssetUrl } from "../utils/asset";

const props = defineProps({
  item: {
    type: Object,
    required: true
  }
});

const statusLabel = computed(() => {
  const map = {
    draft: "草稿",
    pending: "审核中",
    published: "已发布",
    rejected: "已驳回"
  };
  return map[props.item.status] || "未知状态";
});

const coverUrl = computed(() => toAssetUrl(props.item.coverImage));

function formatDate(value) {
  return new Date(value).toLocaleString("zh-CN");
}
</script>

<style scoped>
.article-card--localized {
  isolation: isolate;
}

.article-card__more {
  color: #ffe3b9;
}

@media (max-width: 768px) {
  .article-card--localized {
    padding: 18px;
  }

  .article-card--localized h3 {
    font-size: 1.15rem;
  }

  .article-card--localized p {
    font-size: 0.92rem;
  }

  .article-cover {
    height: 120px;
    margin: -18px -18px 14px;
  }
}

@media (max-width: 480px) {
  .article-card--localized {
    padding: 14px;
    border-radius: 22px;
  }

  .article-card--localized h3 {
    font-size: 1rem;
  }

  .article-card--localized footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .article-cover {
    height: 100px;
    margin: -14px -14px 12px;
  }
}
</style>
