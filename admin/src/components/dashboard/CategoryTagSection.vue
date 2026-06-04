<template>
  <section id="taxonomy" class="manage-grid dashboard-grid">
    <section class="manage-panel">
      <div class="review-head review-head--stack">
        <div>
          <p class="admin-label">内容分类</p>
          <h2>分类管理</h2>
        </div>
        <form class="category-form" @submit.prevent="submitCategory">
          <input
            v-model.trim="categoryFormName"
            type="text"
            maxlength="50"
            placeholder="新增分类名称，例如：架构复盘"
          />
          <button :disabled="categorySubmitting || !categoryFormName">
            {{ categorySubmitting ? "保存中..." : "新增分类" }}
          </button>
        </form>
      </div>

      <div v-if="categories.length" class="category-list">
        <article v-for="item in categories" :key="item.id" class="category-row">
          <input
            v-model.trim="categoryDrafts[item.id]"
            class="filter-select filter-input category-row__input"
            type="text"
            maxlength="50"
            :disabled="Boolean(categorySaving[item.id]) || Boolean(categoryDeleting[item.id])"
            placeholder="分类名称"
            @keyup.enter="renameCategory(item)"
          />
          <span class="taxonomy-usage">文章 {{ item.articleCount ?? 0 }}</span>
          <div class="category-row__actions">
            <button
              class="role-btn"
              :disabled="!canRenameCategory(item) || Boolean(categorySaving[item.id]) || Boolean(categoryDeleting[item.id])"
              @click="renameCategory(item)"
            >
              {{ categorySaving[item.id] ? "保存中..." : "重命名" }}
            </button>
            <button
              class="btn-reject"
              :disabled="Boolean(categorySaving[item.id]) || Boolean(categoryDeleting[item.id])"
              @click="removeCategory(item)"
            >
              {{ categoryDeleting[item.id] ? "删除中..." : "删除" }}
            </button>
          </div>
        </article>
      </div>
      <p v-else class="compact-empty">当前还没有可用分类。</p>
    </section>

    <section class="manage-panel">
      <div class="review-head review-head--stack">
        <div>
          <p class="admin-label">内容标签</p>
          <h2>标签管理</h2>
        </div>
        <form class="category-form" @submit.prevent="submitTag">
          <input
            v-model.trim="tagFormName"
            type="text"
            maxlength="30"
            placeholder="新增标签名称，例如：Golang"
          />
          <button :disabled="tagSubmitting || !tagFormName">
            {{ tagSubmitting ? "保存中..." : "新增标签" }}
          </button>
        </form>
      </div>

      <div v-if="tags.length" class="category-list">
        <article v-for="item in tags" :key="item.id" class="category-row">
          <input
            v-model.trim="tagDrafts[item.id]"
            class="filter-select filter-input category-row__input"
            type="text"
            maxlength="30"
            :disabled="Boolean(tagSaving[item.id]) || Boolean(tagDeleting[item.id])"
            placeholder="标签名称"
            @keyup.enter="renameTag(item)"
          />
          <span class="taxonomy-usage">文章 {{ item.articleCount ?? 0 }}</span>
          <div class="category-row__actions">
            <button
              class="role-btn"
              :disabled="!canRenameTag(item) || Boolean(tagSaving[item.id]) || Boolean(tagDeleting[item.id])"
              @click="renameTag(item)"
            >
              {{ tagSaving[item.id] ? "保存中..." : "重命名" }}
            </button>
            <button
              class="btn-reject"
              :disabled="Boolean(tagSaving[item.id]) || Boolean(tagDeleting[item.id])"
              @click="removeTag(item)"
            >
              {{ tagDeleting[item.id] ? "删除中..." : "删除" }}
            </button>
          </div>
        </article>
      </div>
      <p v-else class="compact-empty">当前还没有可用标签。</p>
    </section>
  </section>
</template>

<script setup>
import { ref, watch } from "vue";
import {
  createCategory,
  createTag,
  deleteCategory,
  deleteTag,
  updateCategory,
  updateTag
} from "../../api/meta";

const emit = defineEmits(["flash", "data-changed"]);

const props = defineProps({
  categories: { type: Array, required: true },
  tags: { type: Array, required: true }
});

const categoryFormName = ref("");
const categorySubmitting = ref(false);
const categoryDrafts = ref({});
const categorySaving = ref({});
const categoryDeleting = ref({});

const tagFormName = ref("");
const tagSubmitting = ref(false);
const tagDrafts = ref({});
const tagSaving = ref({});
const tagDeleting = ref({});

function syncDrafts() {
  categoryDrafts.value = Object.fromEntries(props.categories.map((item) => [item.id, item.name || ""]));
  tagDrafts.value = Object.fromEntries(props.tags.map((item) => [item.id, item.name || ""]));
}

watch(() => [props.categories, props.tags], syncDrafts, { immediate: true });

function say(message) {
  emit("flash", message);
}

function canRenameCategory(item) {
  const nextName = categoryDrafts.value[item.id]?.trim();
  return Boolean(nextName) && nextName !== item.name;
}

function canRenameTag(item) {
  const nextName = tagDrafts.value[item.id]?.trim();
  return Boolean(nextName) && nextName !== item.name;
}

async function submitCategory() {
  if (!categoryFormName.value || categorySubmitting.value) return;
  categorySubmitting.value = true;
  try {
    await createCategory({ name: categoryFormName.value });
    categoryFormName.value = "";
    emit("data-changed");
    say("分类已新增。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "新增分类失败。");
  } finally {
    categorySubmitting.value = false;
  }
}

async function renameCategory(item) {
  const nextName = categoryDrafts.value[item.id]?.trim();
  if (!nextName || nextName === item.name || categorySaving.value[item.id] || categoryDeleting.value[item.id]) return;
  categorySaving.value = { ...categorySaving.value, [item.id]: true };
  try {
    await updateCategory(item.id, { name: nextName });
    emit("data-changed");
    say("分类已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新分类失败。");
  } finally {
    categorySaving.value = { ...categorySaving.value, [item.id]: false };
  }
}

async function removeCategory(item) {
  if (categorySaving.value[item.id] || categoryDeleting.value[item.id]) return;
  if (!window.confirm(`确定删除分类"${item.name}"吗？`)) return;
  categoryDeleting.value = { ...categoryDeleting.value, [item.id]: true };
  try {
    await deleteCategory(item.id);
    emit("data-changed");
    say("分类已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除分类失败。");
  } finally {
    categoryDeleting.value = { ...categoryDeleting.value, [item.id]: false };
  }
}

async function submitTag() {
  if (!tagFormName.value || tagSubmitting.value) return;
  tagSubmitting.value = true;
  try {
    await createTag({ name: tagFormName.value });
    tagFormName.value = "";
    emit("data-changed");
    say("标签已新增。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "新增标签失败。");
  } finally {
    tagSubmitting.value = false;
  }
}

async function renameTag(item) {
  const nextName = tagDrafts.value[item.id]?.trim();
  if (!nextName || nextName === item.name || tagSaving.value[item.id] || tagDeleting.value[item.id]) return;
  tagSaving.value = { ...tagSaving.value, [item.id]: true };
  try {
    await updateTag(item.id, { name: nextName });
    emit("data-changed");
    say("标签已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新标签失败。");
  } finally {
    tagSaving.value = { ...tagSaving.value, [item.id]: false };
  }
}

async function removeTag(item) {
  if (tagSaving.value[item.id] || tagDeleting.value[item.id]) return;
  if (!window.confirm(`确定删除标签"${item.name}"吗？`)) return;
  tagDeleting.value = { ...tagDeleting.value, [item.id]: true };
  try {
    await deleteTag(item.id);
    emit("data-changed");
    say("标签已删除。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "删除标签失败。");
  } finally {
    tagDeleting.value = { ...tagDeleting.value, [item.id]: false };
  }
}
</script>

<style scoped>
.dashboard-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.review-head--stack {
  align-items: flex-start;
}

.category-list {
  display: grid;
  gap: 12px;
}

.category-row {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.category-row__input {
  flex: 1 1 240px;
}

.category-row__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.taxonomy-usage {
  color: var(--soft);
  font-size: 0.9rem;
  white-space: nowrap;
}

.compact-empty--actionable {
  display: grid;
  gap: 10px;
}

.compact-empty--actionable p {
  margin: 0;
}

.compact-empty__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 960px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .category-form {
    flex-direction: column;
    width: 100%;
  }
  .category-form input {
    width: 100%;
  }
  .category-form button {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .category-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .category-row__input {
    flex: 1 1 auto;
    width: 100%;
  }
  .category-row__actions {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
