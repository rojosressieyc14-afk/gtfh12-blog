<template>
  <div class="doc-toc">
    <p class="doc-toc-label">目录</p>
    <div v-if="items.length" class="doc-toc-list">
      <button
        v-for="item in items"
        :key="item.id"
        class="doc-toc-link"
        :class="[`doc-toc-link--h${item.level}`, { 'doc-toc-link--active': activeId === item.id }]"
        @click="$emit('jump', item.id)"
      >
        {{ item.text }}
      </button>
    </div>
    <p v-else class="doc-toc-empty">暂无目录</p>
  </div>
</template>

<script setup>
defineProps({
  items: { type: Array, default: () => [] },
  activeId: { type: [String, Number], default: null },
});
defineEmits(['jump']);
</script>

<style scoped>
.doc-toc {
  padding: 16px 12px;
}
.doc-toc-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: rgba(246,241,234,0.4);
  margin-bottom: 12px;
}
.doc-toc-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.doc-toc-link {
  display: block;
  width: 100%;
  text-align: left;
  border: none;
  background: none;
  padding: 5px 8px;
  border-radius: 4px;
  font-size: 13px;
  color: rgba(246,241,234,0.5);
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
  line-height: 1.4;
}
.doc-toc-link:hover {
  color: rgba(246,241,234,0.8);
  background: rgba(255,255,255,0.04);
}
.doc-toc-link--active {
  color: #ff8a4c;
  background: rgba(255,138,76,0.1);
}
.doc-toc-link--h3 {
  padding-left: 20px;
  font-size: 12px;
}
.doc-toc-link--h4 {
  padding-left: 28px;
  font-size: 11px;
}
.doc-toc-empty {
  font-size: 12px;
  color: rgba(246,241,234,0.3);
}
</style>
