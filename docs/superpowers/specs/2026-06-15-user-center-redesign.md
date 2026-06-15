# User Center Redesign

Date: 2026-06-15

## Background

用户反馈个人中心布局拥挤、不够专业。当前侧边栏 11 项扁平排列（含 emoji 图标），概览页信息单薄。需要重新布局，提升专业感。

## Current Problems

- 侧边栏 11 项无分组，视觉拥挤
- emoji 图标不统一，显得不专业
- 概览页只有角色/状态两个 stat + 快捷卡片（与侧边栏重复）
- 缺少创作数据展示

## Design

### Two-level Sidebar Navigation

```ascii
┌──────────────────┐
│ [avatar] username │  ← user card
│ 角色标签          │
│                   │
│ ── ◎ 工作台 ──  │  ← tab
│   📊 概览        │
│   📄 文章        │
│   📦 项目        │
│   📚 知识库      │
│   ⭐ 收藏        │
│                   │
│ ── ◎ 设置 ────  │  ← tab
│   👤 个人资料    │
│   🔑 API Key    │
│   🔔 通知        │
│   ℹ️ 关于        │
│                   │
│ [+ 写文章]       │  ← sticky button
└──────────────────┘
```

- **Tab 切换**：侧边栏顶部两个 tab（工作台 / 设置），点击切换下方菜单项
- **工作台**（5 项）：概览 / 文章 / 项目 / 知识库 / 收藏
- **设置**（4 项）：个人资料 / API Key / 通知 / 关于
- **顶部**：用户头像 + 用户名 + 角色标签（管理员/普通用户）
- **底部**：固定「写文章」按钮
- **间距**：整体 padding 加大，nav 项间 gap 增加
- **激活态**：左侧指示条 + 背景高亮

### Overview Page (Dashboard)

- **欢迎语**：保持现有
- **4 指标卡**：文章数 / 项目数 / 知识库文档数 / 总浏览量（新增后端接口）
- **最近动态**：用户最近操作时间线
- **快捷操作**：写文章 / 新建项目 两个快捷按钮

### Backend New Endpoints

1. `GET /api/user-center/stats` — 用户创作统计（需要新增 handler + service）
   - articles_count (all statuses)
   - projects_count (all statuses)
   - kb_docs_count
   - total_views (sum of article view_count)

2. `GET /api/user-center/recent-activity` — 最近动态（从 articles/projects/kb 的 updatedAt 排序取最近 10 条）

### Files to Modify

| File | Change |
|------|--------|
| `web/src/components/UserCenterLayout.vue` | Sidebar redesign: tabs, user card, grouping, larger spacing, SVG icons |
| `web/src/views/UserCenterOverview.vue` | Redesign: 4 stat cards, recent activity, quick actions |
| `server/internal/handler/user_handler.go` | Add Stats + RecentActivity handlers |
| `server/internal/service/user_service.go` | Add Stats + RecentActivity service methods |
| `server/internal/router/router.go` | Add new routes |

### Icons

Use simple inline SVG icons or consistent Unicode symbols. No emoji.

### Not In Scope

- 子页面内容改造（MyArticlesView, KnowledgeBaseView 等仅改边距适应新布局）
- 移动端侧边栏行为保持现有（横向滚动 tab）

## Spec Review

- No placeholders or TODOs
- Scope focused on layout + overview + API
- No contradictions
- Requirements unambiguous
