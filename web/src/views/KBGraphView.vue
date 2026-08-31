<template>
  <section class="kb-graph-page">
    <div class="kb-graph-page__header">
      <router-link class="ghost-btn" :to="`/knowledge-bases/${kbId}`">&larr; 返回知识库</router-link>
      <h2>文档关系图谱</h2>
      <span class="table-note" v-if="graphData">{{ graphData.nodes?.length || 0 }} 个节点 · {{ graphData.edges?.length || 0 }} 条关联</span>
    </div>
    <div class="kb-graph-page__body">
      <div v-if="loading" class="kb-graph-page__empty">加载中...</div>
      <div v-else-if="!graphData?.nodes?.length" class="kb-graph-page__empty">暂无文档关系数据</div>
      <svg v-else ref="graphSvg" class="kb-graph-page__svg"></svg>
    </div>
    <div class="kb-graph-page__legend">
      <span class="kb-graph-legend__item">
        <span class="kb-graph-legend__dot kb-graph-legend__dot--active"></span> 当前文档
      </span>
      <span class="kb-graph-legend__item">
        <span class="kb-graph-legend__dot"></span> 其他文档
      </span>
      <span class="kb-graph-legend__item">
        <svg width="20" height="10"><line x1="0" y1="5" x2="20" y2="5" stroke="rgba(246,241,234,0.3)" stroke-width="2"/><marker id="legend-arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill="rgba(246,241,234,0.3)"/></marker></svg>
        Wiki-link 关联
      </span>
    </div>
  </section>
</template>

<script setup>
import { nextTick, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getGraphData, listDocuments } from "../api/knowledgeBase";

const route = useRoute();
const router = useRouter();
const kbId = route.params.id;
const graphData = ref(null);
const allDocs = ref([]);
const loading = ref(true);
const graphSvg = ref(null);

function renderGraph(data) {
  const svg = graphSvg.value;
  if (!svg) return;
  const rect = svg.getBoundingClientRect();
  const W = rect.width || 900, H = rect.height || 560;
  svg.innerHTML = "";

  const nodeMap = {};
  data.nodes.forEach((n, i) => {
    const angle = (2 * Math.PI * i) / data.nodes.length;
    const radius = Math.min(W, H) * 0.32;
    nodeMap[n.id] = { ...n, x: W/2 + radius * Math.cos(angle), y: H/2 + radius * Math.sin(angle), vx: 0, vy: 0 };
  });

  for (let iter = 0; iter < 200; iter++) {
    data.nodes.forEach((n) => {
      const a = nodeMap[n.id];
      data.nodes.forEach((m) => {
        if (n.id === m.id) return;
        const b = nodeMap[m.id];
        let dx = a.x - b.x, dy = a.y - b.y, dist = Math.sqrt(dx*dx + dy*dy) || 1;
        let force = 5000 / (dist * dist);
        a.vx += (dx / dist) * force; a.vy += (dy / dist) * force;
      });
    });
    data.edges.forEach((e) => {
      const a = nodeMap[e.source], b = nodeMap[e.target];
      if (!a || !b) return;
      let dx = b.x - a.x, dy = b.y - a.y, dist = Math.sqrt(dx*dx + dy*dy) || 1;
      let force = (dist - 140) * 0.004;
      a.vx += (dx / dist) * force; a.vy += (dy / dist) * force;
      b.vx -= (dx / dist) * force; b.vy -= (dy / dist) * force;
    });
    data.nodes.forEach((n) => {
      const a = nodeMap[n.id];
      a.vx *= 0.82; a.vy *= 0.82; a.x += a.vx; a.y += a.vy;
      a.x = Math.max(60, Math.min(W - 60, a.x));
      a.y = Math.max(60, Math.min(H - 60, a.y));
    });
  }

  const defs = document.createElementNS("http://www.w3.org/2000/svg", "defs");
  defs.innerHTML = `<marker id="arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 0 L 10 5 L 0 10 z" fill="rgba(246,241,234,0.25)"/></marker>`;
  svg.appendChild(defs);

  data.edges.forEach((e) => {
    const a = nodeMap[e.source], b = nodeMap[e.target];
    if (!a || !b) return;
    const line = document.createElementNS("http://www.w3.org/2000/svg", "line");
    line.setAttribute("x1", a.x); line.setAttribute("y1", a.y);
    line.setAttribute("x2", b.x); line.setAttribute("y2", b.y);
    line.setAttribute("stroke", "rgba(246,241,234,0.15)");
    line.setAttribute("stroke-width", "1.5");
    line.setAttribute("marker-end", "url(#arrow)");
    svg.appendChild(line);
  });

  data.nodes.forEach((n) => {
    const a = nodeMap[n.id];
    const g = document.createElementNS("http://www.w3.org/2000/svg", "g");
    g.style.cursor = "pointer";
    g.addEventListener("click", () => {
      router.push({ name: "uc-knowledge-base-detail", params: { id: kbId } });
    });
    g.addEventListener("mouseenter", () => {
      circle.setAttribute("r", String(Math.max(8, Math.min(20, 8 + n.links * 2.5))));
      circle.setAttribute("fill", "#ff8a4c");
    });
    g.addEventListener("mouseleave", () => {
      circle.setAttribute("r", String(Math.max(6, Math.min(16, 6 + n.links * 2))));
      circle.setAttribute("fill", "rgba(255,138,76,0.6)");
    });

    const r = Math.max(6, Math.min(16, 6 + n.links * 2));
    const circle = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    circle.setAttribute("cx", a.x); circle.setAttribute("cy", a.y);
    circle.setAttribute("r", String(r));
    circle.setAttribute("fill", "rgba(255,138,76,0.6)");
    circle.setAttribute("stroke", "#ff8a4c");
    circle.setAttribute("stroke-width", "1.5");
    g.appendChild(circle);

    const text = document.createElementNS("http://www.w3.org/2000/svg", "text");
    text.setAttribute("x", a.x); text.setAttribute("y", a.y + r + 16);
    text.setAttribute("text-anchor", "middle");
    text.setAttribute("fill", "rgba(246,241,234,0.75)");
    text.setAttribute("font-size", "12");
    text.setAttribute("font-weight", "500");
    text.textContent = n.title.length > 16 ? n.title.slice(0, 16) + "…" : n.title;
    g.appendChild(text);

    svg.appendChild(g);
  });
}

onMounted(async () => {
  try {
    const [graphRes, docRes] = await Promise.all([
      getGraphData(kbId),
      listDocuments(kbId).catch(() => ({ data: { items: [] } })),
    ]);
    graphData.value = graphRes.data.graph || { nodes: [], edges: [] };
    allDocs.value = docRes.data.items || [];
    await nextTick();
    renderGraph(graphData.value);
  } catch { graphData.value = { nodes: [], edges: [] }; } finally { loading.value = false; }
});
</script>

<style scoped>
.kb-graph-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  gap: 16px;
}
.kb-graph-page__header {
  display: flex;
  align-items: center;
  gap: 16px;
}
.kb-graph-page__header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
}
.kb-graph-page__body {
  flex: 1;
  border-radius: 20px;
  border: 1px solid var(--border, rgba(255,255,255,0.1));
  background: rgba(255,255,255,0.03);
  overflow: hidden;
  position: relative;
}
.kb-graph-page__svg {
  width: 100%;
  height: 100%;
}
.kb-graph-page__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: rgba(246,241,234,0.3);
  font-size: 15px;
}
.kb-graph-page__legend {
  display: flex;
  gap: 20px;
  align-items: center;
  padding: 8px 0;
  font-size: 12px;
  color: rgba(246,241,234,0.5);
}
.kb-graph-legend__item {
  display: flex;
  align-items: center;
  gap: 6px;
}
.kb-graph-legend__dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255,138,76,0.6);
  border: 1.5px solid #ff8a4c;
}
.kb-graph-legend__dot--active {
  background: #ff8a4c;
  box-shadow: 0 0 8px rgba(255,138,76,0.5);
}
</style>
