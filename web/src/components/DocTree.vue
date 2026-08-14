<template>
  <div class="doc-tree">
    <div
      v-for="item in tree"
      :key="item.id"
      class="doc-tree-node"
    >
      <div
        class="doc-tree-item"
        :class="{ 'doc-tree-item--active': activeId === item.id, 'doc-tree-item--folder': item.children?.length }"
        @click="$emit('select', item)"
      >
        <span class="doc-tree-icon">{{ item.children?.length ? '📁' : '📄' }}</span>
        <span class="doc-tree-title">{{ item.title || '无标题' }}</span>
      </div>
      <div v-if="item.children?.length" class="doc-tree-children">
        <DocTree
          :tree="item.children"
          :active-id="activeId"
          @select="(n) => $emit('select', n)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({ name: 'DocTree' });
defineProps({
  tree: { type: Array, default: () => [] },
  activeId: { type: [Number, String], default: null },
});
defineEmits(['select']);
</script>

<style scoped>
.doc-tree {
  font-size: 13px;
}
.doc-tree-node {
  margin: 0;
}
.doc-tree-children {
  padding-left: 16px;
}
.doc-tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  color: rgba(246,241,234,0.7);
  transition: background 0.15s, color 0.15s;
  margin-bottom: 2px;
}
.doc-tree-item:hover {
  background: rgba(255,255,255,0.06);
  color: rgba(246,241,234,0.9);
}
.doc-tree-item--active {
  background: rgba(255,138,76,0.15);
  color: #ff8a4c;
}
.doc-tree-item--folder {
  font-weight: 500;
}
.doc-tree-icon {
  font-size: 14px;
  line-height: 1;
}
.doc-tree-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
