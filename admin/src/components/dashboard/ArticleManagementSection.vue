<template>
  <section id="articles" class="manage-panel">
    <div class="review-head review-head--stack">
      <div>
        <p class="admin-label">文章管理</p>
        <h2>文章列表</h2>
      </div>
      <div class="toolbar-row">
        <input
          v-model.trim="articleKeyword"
          class="filter-select filter-input"
          placeholder="搜索标题 / 作者 / 摘要"
          @keyup.enter="changeArticlePage(1)"
        />
        <select v-model="articleStatusFilter" class="filter-select" @change="changeArticlePage(1)">
          <option value="">全部状态</option>
          <option value="draft">草稿</option>
          <option value="pending">审核中</option>
          <option value="published">已发布</option>
          <option value="rejected">已驳回</option>
        </select>
        <select v-model="articleCategoryFilter" class="filter-select" @change="changeArticlePage(1)">
          <option value="">全部分类</option>
          <option v-for="item in categories" :key="`article-category-${item.id}`" :value="String(item.id)">
            {{ item.name }}
          </option>
        </select>
        <select v-model="articleTagFilter" class="filter-select" @change="changeArticlePage(1)">
          <option value="">全部标签</option>
          <option v-for="item in tags" :key="`article-tag-${item.id}`" :value="String(item.id)">
            {{ item.name }}
          </option>
        </select>
        <button class="role-btn" @click="changeArticlePage(1)">搜索</button>
      </div>
    </div>

    <div class="article-bulk-bar">
      <button class="role-btn" @click="toggleSelectVisibleArticles">
        {{ allVisibleArticlesSelected ? "取消本页全选" : "全选本页" }}
      </button>
      <span class="taxonomy-usage">已选 {{ selectedArticleIds.length }} 篇</span>
      <select v-model="bulkArticleCategoryId" class="filter-select">
        <option value="">批量调整分类</option>
        <option v-for="item in categories" :key="`bulk-article-category-${item.id}`" :value="String(item.id)">
          {{ item.name }}
        </option>
      </select>
      <div class="taxonomy-chip-list">
        <label v-for="item in tags" :key="`bulk-article-tag-${item.id}`" class="taxonomy-chip">
          <input :checked="articleBulkHasTag(item.id)" type="checkbox" @change="toggleBulkArticleTag(item.id)" />
          <span>{{ item.name }}</span>
        </label>
      </div>
      <button class="role-btn" :disabled="!selectedArticleIds.length || bulkArticleTaxonomySaving" @click="saveBulkArticleTaxonomy">
        {{ bulkArticleTaxonomySaving ? "保存中..." : "批量保存元数据" }}
      </button>
      <button class="btn-approve" :disabled="!selectedArticleIds.length || bulkArticlePublishSaving" @click="bulkPublishSelectedArticles">
        {{ bulkArticlePublishSaving ? "发布中..." : "批量发布" }}
      </button>
      <button class="btn-reject" :disabled="!selectedArticleIds.length || bulkArticleDeleteSaving" @click="bulkDeleteSelectedArticles">
        {{ bulkArticleDeleteSaving ? "删除中..." : "批量删除" }}
      </button>
    </div>

    <div class="table-list">
      <article v-for="item in articles" :key="item.id" class="table-card table-card--article">
        <div class="table-main">
          <div class="table-title">
            <label class="article-select-check">
              <input :checked="selectedArticleIds.includes(item.id)" type="checkbox" @change="toggleArticleSelection(item.id)" />
            </label>
            <h3>{{ item.title }}</h3>
            <span class="pill" :class="statusClass(item.status)">{{ statusLabel(item.status) }}</span>
          </div>
          <p>{{ item.author?.username || "未知作者" }} · {{ item.category?.name || "未分类" }} · {{ formatDate(item.updatedAt) }}</p>
          <p class="table-note">{{ item.summary || "暂无摘要" }}</p>
          <p v-if="item.rejectReason" class="reject-note">驳回原因：{{ item.rejectReason }}</p>
          <div class="article-quick-taxonomy">
            <label class="article-quick-taxonomy__field">
              <span>分类</span>
              <select v-model="articleCategoryDrafts[item.id]" class="filter-select">
                <option value="">未分类</option>
                <option v-for="meta in categories" :key="`quick-category-${item.id}-${meta.id}`" :value="String(meta.id)">
                  {{ meta.name }}
                </option>
              </select>
            </label>
            <div class="article-quick-taxonomy__field article-quick-taxonomy__field--stack">
              <span>标签</span>
              <div class="taxonomy-chip-list">
                <label v-for="meta in tags" :key="`quick-tag-${item.id}-${meta.id}`" class="taxonomy-chip">
                  <input
                    :checked="articleHasTagDraft(item, meta.id)"
                    type="checkbox"
                    @change="toggleArticleTagDraft(item, meta.id)"
                  />
                  <span>{{ meta.name }}</span>
                </label>
              </div>
            </div>
          </div>
          <div class="table-actions">
            <button
              class="role-btn"
              :disabled="!canSaveArticleQuickTaxonomy(item) || Boolean(articleQuickTaxonomySaving[item.id])"
              @click="saveArticleQuickTaxonomy(item)"
            >
              {{ articleQuickTaxonomySaving[item.id] ? "保存中..." : "保存元数据" }}
            </button>
            <button class="role-btn" @click="$emit('open-article-preview', item.id, item.status === 'pending')">
              {{ item.status === "pending" ? "审核" : "预览" }}
            </button>
            <button v-if="item.status !== 'published'" class="btn-approve" @click="publishArticle(item.id)">发布</button>
            <button class="btn-reject" @click="deleteArticle(item.id)">删除</button>
          </div>
        </div>
      </article>
    </div>

    <div class="pager-row">
      <button class="role-btn" :disabled="articlePage <= 1" @click="changeArticlePage(articlePage - 1)">上一页</button>
      <span>第 {{ articlePage }} / {{ articleTotalPages }} 页</span>
      <button class="role-btn" :disabled="articlePage >= articleTotalPages" @click="changeArticlePage(articlePage + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import {
  getAdminArticles,
  updateAdminArticleTaxonomy,
  bulkUpdateAdminArticleTaxonomy,
  bulkPublishAdminArticles,
  bulkDeleteAdminArticles,
  publishAdminArticle,
  deleteAdminArticle,
  getPendingReviews,
  getPendingProjectReviews
} from "../../api/dashboard";

const emit = defineEmits(["flash", "open-article-preview", "data-changed"]);

const props = defineProps({
  categories: { type: Array, required: true },
  tags: { type: Array, required: true }
});

const articles = ref([]);
const articleStatusFilter = ref("");
const articleKeyword = ref("");
const articleCategoryFilter = ref("");
const articleTagFilter = ref("");
const articlePage = ref(1);
const articleTotal = ref(0);
const pageSize = 6;
const selectedArticleIds = ref([]);
const bulkArticleCategoryId = ref("");
const bulkArticleTagIds = ref([]);
const bulkArticleTaxonomySaving = ref(false);
const bulkArticlePublishSaving = ref(false);
const bulkArticleDeleteSaving = ref(false);
const articleCategoryDrafts = ref({});
const articleTagDrafts = ref({});
const articleQuickTaxonomySaving = ref({});

function say(message) {
  emit("flash", message);
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString("zh-CN") : "暂无时间";
}

function statusLabel(status) {
  return {
    draft: "草稿",
    pending: "审核中",
    published: "已发布",
    rejected: "已驳回"
  }[status] || status;
}

function statusClass(status) {
  return {
    "pill--pending": status === "pending",
    "pill--success": status === "published",
    "pill--reject": status === "rejected"
  };
}

const articleTotalPages = computed(() => Math.max(1, Math.ceil(articleTotal.value / pageSize)));

const allVisibleArticlesSelected = computed(
  () => articles.value.length > 0 && articles.value.every((item) => selectedArticleIds.value.includes(item.id))
);

function syncArticleTaxonomyDrafts(items) {
  articleCategoryDrafts.value = Object.fromEntries(
    items.map((item) => [item.id, item.categoryId ? String(item.categoryId) : ""])
  );
  articleTagDrafts.value = Object.fromEntries(
    items.map((item) => [item.id, (item.tags || []).map((tag) => tag.id).filter(Boolean)])
  );
}

function resetBulkArticleDrafts() {
  selectedArticleIds.value = [];
  bulkArticleCategoryId.value = "";
  bulkArticleTagIds.value = [];
  bulkArticleTaxonomySaving.value = false;
  bulkArticlePublishSaving.value = false;
  bulkArticleDeleteSaving.value = false;
}

function articleHasTagDraft(item, tagID) {
  return (articleTagDrafts.value[item.id] || []).includes(tagID);
}

function toggleArticleTagDraft(item, tagID) {
  const current = articleTagDrafts.value[item.id] || [];
  articleTagDrafts.value = {
    ...articleTagDrafts.value,
    [item.id]: current.includes(tagID) ? current.filter((value) => value !== tagID) : [...current, tagID]
  };
}

function articleBulkHasTag(tagID) {
  return bulkArticleTagIds.value.includes(tagID);
}

function toggleBulkArticleTag(tagID) {
  if (bulkArticleTagIds.value.includes(tagID)) {
    bulkArticleTagIds.value = bulkArticleTagIds.value.filter((item) => item !== tagID);
    return;
  }
  bulkArticleTagIds.value = [...bulkArticleTagIds.value, tagID];
}

function toggleArticleSelection(articleID) {
  if (selectedArticleIds.value.includes(articleID)) {
    selectedArticleIds.value = selectedArticleIds.value.filter((item) => item !== articleID);
    return;
  }
  selectedArticleIds.value = [...selectedArticleIds.value, articleID];
}

function toggleSelectVisibleArticles() {
  if (allVisibleArticlesSelected.value) {
    selectedArticleIds.value = [];
    return;
  }
  selectedArticleIds.value = articles.value.map((item) => item.id);
}

function canSaveArticleQuickTaxonomy(item) {
  const nextCategory = articleCategoryDrafts.value[item.id] || "";
  const currentCategory = item.categoryId ? String(item.categoryId) : "";
  const nextTags = [...(articleTagDrafts.value[item.id] || [])].sort((left, right) => left - right);
  const currentTags = [...((item.tags || []).map((tag) => tag.id).filter(Boolean))].sort((left, right) => left - right);
  return nextCategory !== currentCategory || JSON.stringify(nextTags) !== JSON.stringify(currentTags);
}

async function loadArticles() {
  const { data } = await getAdminArticles({
    page: articlePage.value,
    pageSize,
    status: articleStatusFilter.value,
    keyword: articleKeyword.value,
    categoryId: articleCategoryFilter.value,
    tagId: articleTagFilter.value
  });
  articles.value = data.items || [];
  syncArticleTaxonomyDrafts(articles.value);
  resetBulkArticleDrafts();
  articleTotal.value = data.pagination?.total || 0;
}

async function changeArticlePage(nextPage) {
  articlePage.value = nextPage;
  await loadArticles();
}

async function refreshReviews() {
  try {
    await Promise.all([getPendingReviews(), getPendingProjectReviews()]);
  } catch {
    // silently refresh caches
  }
}

async function publishArticle(id) {
  try {
    await publishAdminArticle(id);
    await loadArticles();
    await refreshReviews();
    emit("data-changed");
    say("文章已发布。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "发布文章失败。");
  }
}

async function deleteArticle(id) {
  if (!window.confirm("确定删除这篇文章吗？")) return;
  try {
    await deleteAdminArticle(id);
    await loadArticles();
    await refreshReviews();
    emit("data-changed");
    say("文章已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除文章失败。");
  }
}

async function saveArticleQuickTaxonomy(item) {
  if (!canSaveArticleQuickTaxonomy(item) || articleQuickTaxonomySaving.value[item.id]) return;
  articleQuickTaxonomySaving.value = { ...articleQuickTaxonomySaving.value, [item.id]: true };
  try {
    const payload = {
      categoryId: articleCategoryDrafts.value[item.id] ? Number(articleCategoryDrafts.value[item.id]) : null,
      tagIds: articleTagDrafts.value[item.id] || []
    };
    await updateAdminArticleTaxonomy(item.id, payload);
    await loadArticles();
    await refreshReviews();
    say("文章元数据已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新文章元数据失败。");
  } finally {
    articleQuickTaxonomySaving.value = { ...articleQuickTaxonomySaving.value, [item.id]: false };
  }
}

async function saveBulkArticleTaxonomy() {
  if (!selectedArticleIds.value.length || bulkArticleTaxonomySaving.value) return;
  bulkArticleTaxonomySaving.value = true;
  try {
    await bulkUpdateAdminArticleTaxonomy({
      articleIds: selectedArticleIds.value,
      categoryId: bulkArticleCategoryId.value ? Number(bulkArticleCategoryId.value) : null,
      tagIds: bulkArticleTagIds.value,
      replaceTags: true
    });
    await loadArticles();
    await refreshReviews();
    say("批量文章元数据已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "批量更新文章元数据失败。");
  } finally {
    bulkArticleTaxonomySaving.value = false;
  }
}

async function bulkPublishSelectedArticles() {
  if (!selectedArticleIds.value.length || bulkArticlePublishSaving.value) return;
  bulkArticlePublishSaving.value = true;
  try {
    await bulkPublishAdminArticles({ articleIds: selectedArticleIds.value });
    await loadArticles();
    await refreshReviews();
    emit("data-changed");
    say("批量发布已完成。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "批量发布失败。");
  } finally {
    bulkArticlePublishSaving.value = false;
  }
}

async function bulkDeleteSelectedArticles() {
  if (!selectedArticleIds.value.length || bulkArticleDeleteSaving.value) return;
  if (!window.confirm("确定批量删除选中的文章吗？仅草稿和已驳回文章允许删除。")) return;
  bulkArticleDeleteSaving.value = true;
  try {
    await bulkDeleteAdminArticles({ articleIds: selectedArticleIds.value });
    selectedArticleIds.value = [];
    await loadArticles();
    await refreshReviews();
    emit("data-changed");
    say("批量删除已完成。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "批量删除失败。");
  } finally {
    bulkArticleDeleteSaving.value = false;
  }
}

loadArticles();
</script>

<style scoped>
.review-head--stack {
  align-items: flex-start;
}

.table-title {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.pill--pending {
  background: rgba(255, 166, 77, 0.2);
  color: #ffd79d;
}

.pill--success {
  background: rgba(131, 242, 143, 0.16);
  color: #c8f8cd;
}

.pill--reject {
  background: rgba(255, 139, 139, 0.16);
  color: #ffd2d2;
}

.reject-note {
  color: #ffd2d2;
}

.article-bulk-bar {
  display: grid;
  gap: 10px;
}

.taxonomy-usage {
  color: var(--soft);
  font-size: 0.9rem;
  white-space: nowrap;
}

.taxonomy-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.taxonomy-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
}

.article-quick-taxonomy {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.article-quick-taxonomy__field {
  display: grid;
  gap: 8px;
}

.article-quick-taxonomy__field--stack {
  align-items: flex-start;
}

.article-select-check {
  display: inline-flex;
  align-items: center;
}

@media (max-width: 768px) {
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
  .review-head--stack {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .table-card {
    flex-direction: column;
    align-items: flex-start;
  }
  .table-actions {
    flex-direction: column;
    width: 100%;
  }
  .table-actions button {
    width: 100%;
  }
  .article-quick-taxonomy {
    width: 100%;
  }
  .article-bulk-bar {
    gap: 8px;
  }
  .taxonomy-chip-list {
    flex-wrap: wrap;
  }
  .article-bulk-bar .filter-select {
    width: 100%;
  }
  .table-title {
    flex-direction: column;
    align-items: flex-start;
  }
  .article-select-check {
    align-self: flex-start;
  }
  .table-main p {
    word-break: break-word;
  }
}
</style>
