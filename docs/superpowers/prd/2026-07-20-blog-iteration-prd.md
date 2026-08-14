# PulseBlog 迭代 PRD

## 背景

博客已完成知识库文档化改造和文章排版升级。经审查发现当前改动存在若干 Bug，且项目整体缺少一些关键功能。本 PRD 分阶段覆盖 Bug 修复和新功能。

---

## Phase 1：Bug 修复（当前改动遗留问题）

### P1.1 知识库 TOC 跳转无效

**问题：** `docToc` 用自定义正则提取 heading ID，而 `marked` 渲染时用自己的 ID 生成逻辑，两者不一致，点击目录不会跳转。

**方案：** 把 `KnowledgeBaseDetail.vue` 的 `renderedContent` 改为使用自定义 `marked.Renderer`，确保 heading ID 与 `docToc` 生成的 ID 一致（复用 DetailView 的 `buildHeadingId` 逻辑）。

**涉及文件：** `web/src/views/KnowledgeBaseDetail.vue`

---

### P1.2 知识库 "+" 按钮无反应

**问题：** `showNewDoc` 变量声明但模板未绑定任何交互。

**方案：** 点击 "+" 时直接弹出"新建文档"对话框（内联轻弹窗）或直接导航到编辑器。

**涉及文件：** `web/src/views/KnowledgeBaseDetail.vue`

---

### P1.3 文章详情 heading ID 匹配失败

**问题：** 当 heading 包含内联格式（如 `# Hello **World**`）时，`articleOutline` 提取的是原始 Markdown 文本，而 `marked.Renderer` 中的 `token.text` 是纯文本，两者不匹配。

**方案：** 在 `articleOutline` 生成时用 `marked.parseInline` 预处理文本，或改用正则去除内联标记，保证与 `token.text` 一致。

**涉及文件：** `web/src/views/DetailView.vue`

---

### P1.4 MoveDocument 无父节点校验

**问题：** 后端 `MoveDocument` 接受任意 `ParentID`，不验证目标是否存在且属于同一知识库，可能导致孤儿文档。

**方案：** 在 `MoveDocument` 中增加 `ParentID` 存在性校验：若 `ParentID != nil`，则查询对应文档是否存在且 `knowledge_base_id` 匹配。

**涉及文件：** `server/internal/service/knowledge_base_service.go`

---

### P1.5 知识库无加载状态

**问题：** 页面渲染时 API 请求中无任何加载反馈。

**方案：** 在内容区显示骨架屏或 "加载中..." 提示。

**涉及文件：** `web/src/views/KnowledgeBaseDetail.vue`

---

### P1.6 错误处理用 alert()

**问题：** 多个页面在 API 失败时使用 `alert()`，体验差。

**方案：** 改为内联错误提示文本（已有 `errorMessage` 模式可用），统一使用内联 `.error-text` 样式。

**涉及文件：** `web/src/views/KnowledgeBaseDetail.vue`, `web/src/views/KnowledgeBaseView.vue`

---

### P1.7 文章详情缺少 API 错误处理

**问题：** `loadDetail`/`loadComments` 没有 try/catch，失败时页面白屏。

**方案：** 给两个函数加 try/catch，失败时显示错误状态。

**涉及文件：** `web/src/views/DetailView.vue`

---

### P1.8 图片样式破坏行内图片

**问题：** `.markdown-body img { display: block }` 让行内图片（如表情、小图标）变成块级。

**方案：** 改为只对独立图片（非 p 内图片）应用 `display: block`，或使用 `img:not([class])` 选择器区分。

**涉及文件：** `web/src/assets/main.css`

---

### P1.9 DocToc 滚动高亮不工作

**问题：** `activeHeadingId` 声明后从未更新。

**方案：** 用 `IntersectionObserver` 监听 heading 元素进入视口，更新 `activeHeadingId`。

**涉及文件：** `web/src/views/KnowledgeBaseDetail.vue`

---

### P1.10 DocTree 递归组件兼容性

**问题：** Vue <3.3 中 `<script setup>` 组件不能自动递归引用自身。

**方案：** 添加 `defineOptions({ name: 'DocTree' })` 或使用 `import` 显式注册。

**涉及文件：** `web/src/components/DocTree.vue`

---

## Phase 2：新功能

### P2.1 文章全文搜索

**需求：** 目前搜索只匹配标题/摘要，缺少正文内容搜索。

**方案：**
- 后端：`GET /api/articles` 的 `keyword` 参数扩展为搜索 title + summary + content（MySQL LIKE 或 GORM 全文索引）
- 前端：搜索框保持现有位置，搜索结果高亮关键词

**涉及文件：**
- `server/internal/service/article_service.go`
- `web/src/views/ArticlesView.vue`（搜索结果高亮）

---

### P2.2 代码语法高亮

**需求：** 代码块目前是纯色样式，缺少语言语法高亮。

**方案：** 引入 `highlight.js` 或 `prism.js`，在 Markdown 渲染后对 `<pre><code>` 元素应用高亮。通过 `marked` 的 `renderer.code` 注入语言类名。

**涉及文件：**
- `web/package.json`（新增 highlight.js 依赖）
- `web/src/assets/main.css`（新增代码高亮主题类）
- `web/src/views/DetailView.vue`（渲染时初始化高亮）
- `web/src/views/KnowledgeBaseDetail.vue`（渲染时初始化高亮）
- `web/src/views/KbNoteView.vue`（渲染时初始化高亮）

---

### P2.3 阅读进度条

**需求：** 文章阅读时顶部显示阅读进度条，提升沉浸感。

**方案：** 在 `DetailView.vue` 顶部添加固定定位的进度条（1px 高，橙色），监听 `scroll` 事件计算进度百分比。

**涉及文件：** `web/src/views/DetailView.vue`

---

### P2.4 SEO 元标签

**需求：** 文章详情页缺少 Open Graph / Twitter Card 元标签，社交分享预览差。

**方案：**
- 使用 Vue 的 `useHead`（如 `@unhead/vue`）动态设置 `<meta>` 标签
- 文章页面：`og:title`、`og:description`、`og:image`（封面图）、`og:url`、`twitter:card`
- 首页：`og:site_name`、`description`

**涉及文件：**
- `web/src/main.js`（全局注册 @unhead/vue）
- `web/src/views/DetailView.vue`（文章 meta 标签）
- `web/index.html`（全局默认 meta）

---

### P2.5 社交分享按钮

**需求：** 文章底部缺少分享到社交平台的能力。

**方案：** 在 `DetailView.vue` 文章底部添加分享按钮组（Twitter/X、微信（复制链接）、LinkedIn），使用 URL 参数拼接分享链接。

**涉及文件：** `web/src/views/DetailView.vue`

---

### P2.6 RSS / Atom 订阅

**需求：** 博客缺少 RSS 端点，读者无法通过订阅器跟踪更新。

**方案：** 后端新增 `GET /api/feed.xml` 或 `GET /api/feed` 端点，返回 Atom 格式的最近 20 篇文章。

**涉及文件：**
- `server/internal/handler/feed_handler.go`（新增）
- `server/internal/router/router.go`（注册路由）
- 前端首页添加 RSS 订阅图标链接

---

### P2.7 文章阅读统计图表

**需求：** 作者无法看到文章随时间变化的阅读趋势。

**方案：** 后端记录每日阅读量（`article_daily_views` 表），提供 `GET /api/my/articles/:id/stats` 返回每日趋势数据。前端用简单折线图展示（Canvas 或 Chart.js）。

**涉及文件：**
- `server/internal/model/article.go`（新增 DailyView 模型）
- `server/internal/service/article_service.go`（统计服务）
- `web/src/views/MyArticlesView.vue`（文章统计入口）

---

### P2.8 管理员批量分类/标签编辑 UI

**需求：** 后端已有批量分类/标签 API，但管理面板缺少对应 UI。

**方案：** 在管理面板文章列表增加批量操作下拉菜单，选中文章后批量修改分类/标签。

**涉及文件：** `admin/src/views/ArticlesView.vue`（或管理面板对应视图）

---

## Phase 3：安全与基础设施

### P3.1 Content-Security-Policy 头

**需求：** 目前安全头中缺少 CSP，存在 XSS 风险。

**方案：** 在 `middleware/security.go` 中添加 `Content-Security-Policy` 头，允许加载的资源包括自域名、Google Fonts、Unpkg 等。

**涉及文件：** `server/internal/middleware/security.go`

---

### P3.2 JWT 改用 HttpOnly Cookie

**需求：** 当前 JWT 存储在 `localStorage`，易受 XSS 攻击。

**方案：** 后端在登录/注册时写入 HttpOnly Cookie（`Set-Cookie: token=<jwt>; HttpOnly; Secure; SameSite=Strict`），前端不再手动管理 token。CSRF 保护需同步实现。

**涉及文件：**
- `server/internal/handler/auth_handler.go`（写入 Cookie）
- `web/src/api/client.js`（移除手动 token 逻辑）
- `web/src/stores/user.js`（调整认证流程）

---

## 实施优先级

| 优先级 | 内容 | 预估工时 |
|--------|------|----------|
| P0 | Phase 1 全部 Bug 修复 | 高 |
| P1 | P2.1 全文搜索 + P2.2 语法高亮 | 中 |
| P1 | P2.3 阅读进度条 + P2.4 SEO 标签 | 中 |
| P2 | P2.5 社交分享 + P2.6 RSS | 低 |
| P3 | P2.7 统计图表 + P2.8 管理 UI | 低 |
| P4 | Phase 3 安全加固 | 低（需评估影响） |

---

## 非功能性要求

- 所有新功能向后兼容，不破坏现有路由和 API
- 前端构建后体积增加不超过 100KB（gzip）
- 后端新增端点响应时间 < 200ms
