<template>
  <footer class="site-footer panel-card">
    <div class="site-footer__main">
      <div class="site-footer__brand">
        <div class="brand-badge brand-badge--cluster site-footer__badge">
          <span class="brand-shape brand-shape--square"></span>
          <span class="brand-shape brand-shape--circle"></span>
          <span class="brand-shape brand-shape--triangle"></span>
        </div>
        <div>
          <p class="eyebrow">PulseBlog</p>
          <h3>把项目、写作和长期学习沉淀成真正可展示的个人资产。</h3>
          <p class="site-footer__text">
            这个站点既是博客，也是作品集、学习归档和对外展示的线上名片。
          </p>
        </div>
      </div>

      <div class="site-footer__grid">
        <section class="site-footer__column">
          <p class="site-footer__label">站点入口</p>
          <router-link to="/">首页</router-link>
          <router-link to="/about">关于我</router-link>
          <router-link to="/projects">作品集</router-link>
          <router-link to="/auth">登录 / 注册</router-link>
        </section>

        <section class="site-footer__column">
          <p class="site-footer__label">内容方向</p>
          <span>项目案例</span>
          <span>学习文档</span>
          <span>技术博客</span>
          <span>长期公开输出</span>
        </section>

        <section class="site-footer__column">
          <p class="site-footer__label">创作建议</p>
          <span>先发 2 到 3 个代表项目</span>
          <span>把学习内容整理成系列</span>
          <span>补齐个人资料与公开链接</span>
          <span>持续更新案例和复盘</span>
        </section>
      </div>
    </div>

    <div class="site-footer__bottom">
      <span>© {{ year }} PulseBlog</span>
      <span>个人品牌站 · 作品集 · 学习归档</span>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>

<style scoped>
.site-footer {
  margin: 28px auto 0;
  width: min(1200px, 100%);
  padding: 28px;
}

.site-footer__main {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr);
  gap: 28px;
}

.site-footer__brand {
  display: flex;
  gap: 18px;
  align-items: flex-start;
}

.site-footer__badge {
  flex-shrink: 0;
}

.site-footer__brand h3,
.site-footer__text,
.site-footer__label {
  margin: 0;
}

.site-footer__brand h3 {
  font-size: clamp(1.4rem, 3vw, 2.2rem);
  line-height: 1.15;
}

.site-footer__text {
  margin-top: 12px;
  color: var(--text-soft);
}

.site-footer__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.site-footer__column {
  display: grid;
  gap: 10px;
  padding: 18px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
}

.site-footer__label {
  color: var(--accent-soft);
  font-size: 0.82rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.site-footer__column a,
.site-footer__column span {
  color: var(--text-soft);
}

.site-footer__bottom {
  margin-top: 22px;
  padding-top: 18px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--text-soft);
  flex-wrap: wrap;
}

@media (max-width: 960px) {
  .site-footer__main,
  .site-footer__grid {
    grid-template-columns: 1fr;
  }

  .site-footer__brand {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .site-footer {
    padding: 18px;
    border-radius: 16px;
  }
  .site-footer__grid {
    gap: 12px;
  }
  .site-footer__column {
    padding: 14px;
  }
  .site-footer__bottom {
    flex-direction: column;
    text-align: center;
  }
}
</style>
