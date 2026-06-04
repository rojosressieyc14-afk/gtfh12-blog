<template>
  <section class="not-found-page panel-card">
    <div class="not-found-copy">
      <p class="eyebrow">404</p>
      <h2>这个页面暂时没有出现在你的浏览路径里。</h2>
      <p class="detail-summary">
        你访问的地址可能已经变更、输入有误，或者这里的内容还在整理中。先回到首页、作品集或关于页面继续浏览会更顺畅。
      </p>

      <div class="inline-actions">
        <router-link class="solid-btn" to="/">返回首页</router-link>
        <router-link class="ghost-btn" to="/projects">查看作品集</router-link>
        <router-link class="ghost-btn" to="/about">关于我</router-link>
      </div>
    </div>

    <div class="not-found-orbit">
      <div class="not-found-ring not-found-ring--outer"></div>
      <div class="not-found-ring not-found-ring--mid"></div>
      <div class="not-found-core">404</div>
      <article class="not-found-note not-found-note--one">
        <strong>作品集</strong>
        <p>代表项目、案例页面和完整过程说明都在这里。</p>
      </article>
      <article class="not-found-note not-found-note--two">
        <strong>文章内容</strong>
        <p>学习记录、复盘文章和长期输出会持续沉淀。</p>
      </article>
      <article class="not-found-note not-found-note--three">
        <strong>关于我</strong>
        <p>也可以从个人介绍和对外名片入口继续了解这个站点。</p>
      </article>
    </div>
  </section>
</template>

<style scoped>
.not-found-page {
  margin-top: 28px;
  padding: 36px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(320px, 0.9fr);
  gap: 28px;
  min-height: 480px;
  align-items: center;
  background:
    radial-gradient(circle at 18% 20%, rgba(255, 138, 76, 0.14), transparent 24%),
    radial-gradient(circle at 82% 24%, rgba(255, 209, 102, 0.12), transparent 24%),
    rgba(255, 255, 255, 0.07);
}

.not-found-copy h2 {
  margin: 0 0 18px;
  font-size: clamp(2rem, 4vw, 3.8rem);
  line-height: 1.05;
}

.not-found-orbit {
  position: relative;
  min-height: 360px;
}

.not-found-ring,
.not-found-core {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 999px;
}

.not-found-ring {
  border: 1px dashed rgba(255, 255, 255, 0.16);
}

.not-found-ring--outer {
  width: 280px;
  height: 280px;
}

.not-found-ring--mid {
  width: 180px;
  height: 180px;
}

.not-found-core {
  width: 112px;
  height: 112px;
  display: grid;
  place-items: center;
  font-size: 2rem;
  font-weight: 700;
  color: #1b140e;
  background: linear-gradient(135deg, var(--accent-soft), var(--accent));
  box-shadow: 0 24px 54px rgba(255, 138, 76, 0.28);
}

.not-found-note {
  position: absolute;
  width: 180px;
  padding: 16px 18px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.06);
}

.not-found-note strong,
.not-found-note p {
  margin: 0;
}

.not-found-note p {
  margin-top: 8px;
  color: var(--text-soft);
}

.not-found-note--one {
  left: 0;
  top: 24px;
}

.not-found-note--two {
  right: 0;
  top: 116px;
}

.not-found-note--three {
  left: 56px;
  bottom: 12px;
}

@media (max-width: 960px) {
  .not-found-page {
    grid-template-columns: 1fr;
    padding: 28px;
  }

  .not-found-orbit {
    min-height: auto;
    display: grid;
    gap: 14px;
  }

  .not-found-ring,
  .not-found-core,
  .not-found-note {
    position: relative;
    inset: auto;
    left: auto;
    top: auto;
    right: auto;
    bottom: auto;
    transform: none;
    width: auto;
    height: auto;
  }

  .not-found-core {
    width: 112px;
    height: 112px;
    margin: 0 auto;
  }
}

@media (max-width: 768px) {
  .inline-actions {
    flex-direction: column;
  }

  .inline-actions .solid-btn,
  .inline-actions .ghost-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .not-found-page {
    padding: 18px;
    gap: 18px;
  }

  .not-found-copy h2 {
    font-size: clamp(1.4rem, 6vw, 2rem);
  }
}
</style>
