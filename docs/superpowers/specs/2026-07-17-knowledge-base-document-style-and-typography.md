# 知识库文档化改造 + 文章排版升级

## 背景

博客现有知识库为扁平列表式笔记管理，缺乏层级组织能力；文章排版单一（字体、行距、代码块等均无差异化），阅读体验不佳。

## 目标

1. 知识库改为飞书文档风格的三栏布局（文档树 + 阅读区 + 目录）
2. 建立统一的 .markdown-body 排版系统，覆盖文章详情/知识库笔记/公开笔记

## 1. 知识库文档化

### 1.1 后端改造

- KnowledgeDocument 模型新增：
  - `parent_id *uint` — 父文档 ID（nil 表示根级）
  - `sort_order int` — 同级排序序号
- 新增 API：
  - `GET /api/knowledge-bases/:id/tree` — 返回文档树（递归结构）
  - `PUT /api/knowledge-bases/:id/documents/:docId/move` — 拖拽移动 {parent_id, sort_order}
- 更新现有 API 适配树结构

### 1.2 前端—知识库列表页

保持卡片式布局，点击进入三栏视图。

### 1.3 前端—三栏详情页

- **左栏 (260px)：** DocTree.vue — 递归树组件，文件夹折叠展开，当前高亮，新建/删除操作
- **中栏 (flex:1)：** 文档阅读区 — 面包屑 + 元数据 + Markdown 渲染 + 底部操作
- **右栏 (200px)：** DocToc.vue — 自动提取 h1-h3 目录，滚动高亮 + 点击跳转

### 1.4 前端—编辑器

保持全屏模式，左侧编辑 + 右侧预览（对齐全局 CSS）。

## 2. 文章/笔记排版升级

### 2.1 全局 .markdown-body 样式

- 阅读宽度 720px 居中
- 行高中文 1.8 / 英文 1.6
- 标题 h1-h4 字重梯级
- 代码块 JetBrains Mono + 深色背景 + 圆角
- 引用块左侧橙色竖线
- 图片圆角 + 居中 + 阴影
- 表格隔行变色 + 圆角

### 2.2 移除组件内联样式

KbNoteView.vue、KnowledgeBaseNoteEditor.vue、DetailView.vue 中的本地 .markdown-body 样式移至全局 main.css。

## 实施顺序

| # | 步骤 | 文件 |
|---|------|------|
| 1 | 全局 .markdown-body CSS | web/src/assets/main.css |
| 2 | 提取组件内联样式至全局 | KbNoteView, Editor, DetailView |
| 3 | DetailView 阅读区宽度 | DetailView.vue |
| 4 | 后端 parent_id + tree API | server/internal |
| 5 | DocTree 组件 | web/src/components/DocTree.vue |
| 6 | DocToc 组件 | web/src/components/DocToc.vue |
| 7 | 重构 KB Detail 三栏布局 | KnowledgeBaseDetail.vue |
| 8 | 编辑器预览对齐 | KnowledgeBaseNoteEditor.vue |
