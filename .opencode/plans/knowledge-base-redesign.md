# 知识库重建设计方案（已获用户确认）

## 1. 目标
- 知识库文档从纯文本 → 像文章一样的 Markdown 笔记
- 每篇笔记可设公开/私密，公开笔记有独立展示页
- 每篇笔记可加分类和标签（复用文章分类/标签）
- 保留 AI 向量检索功能
- 添加文档问题：需要配置 DEEPSEEK_API_KEY 或增加降级处理

## 2. 改动清单

### 2.1 后端 Model (`server/internal/model/knowledge_base.go`)
KnowledgeDocument 新增字段：
- `UserID uint` — 作者 ID（公开笔记需要）
- `User User` — 关联作者
- `IsPublic bool` — 是否公开
- `IsMarkdown bool` — 是否 Markdown 内容
- `CategoryID *uint` — 分类
- `Category *Category` — 关联分类
- `Tags []Tag gorm:"many2many:kb_document_tags;"` — 标签
- `ViewCount int64` — 浏览次数
- `PublishedAt *time.Time` — 公开时间

新增结构体 `KbDocumentTag` — 标签关联表

### 2.2 后端 DB Migration (`server/internal/database/db.go`)
- 在 AutoMigrate 中添加 `&model.KbDocumentTag{}`

### 2.3 后端 Service (`server/internal/service/knowledge_base_service.go`)
修改/新增方法：
- `AddDocument` — 增加 `userID`, `isPublic`, `categoryID`, `tagIDs`, `isMarkdown` 参数
- `UpdateDocument` — **新增**，支持编辑笔记内容、标题、分类、标签、公开状态
- `GetPublicNote` — **新增**，无需登录，返回公开笔记内容（+1 viewCount）
- 添加文档时如果 embedding 为 nil（无 API Key），跳过向量化，只存文本

### 2.4 后端 Handler (`server/internal/handler/knowledge_base_handler.go`)
修改/新增方法：
- `AddDocument` — 解析新字段
- `UpdateDocument` — **新增** PATCH 处理器
- `GetPublicNote` — **新增** GET 处理器（无需 auth）

### 2.5 后端 Router (`server/internal/router/router.go`)
新增路由：
- `PUT /api/knowledge-bases/:id/documents/:docId` — 编辑笔记
- `GET /api/kb-notes/:id` — 公开笔记详情（无需 auth，放在 api 分组外）

### 2.6 前端 API (`web/src/api/knowledgeBase.js`)
- `updateDocument(kbId, docId, payload)` — **新增**
- `getPublicNote(id)` — **新增**
- `addDocument` — 更新 payload 结构

### 2.7 前端 KnowledgeBaseDetail.vue（全面重做）
- Tab 切换：「笔记列表」|「AI 检索」
- 笔记卡片显示：标题、内容预览（Markdown 截断）、公开/私密标签、分类、标签、日期
- 筛选栏：按分类、标签筛选
- 操作：编辑、删除、切换公开/私密
- "新建笔记" 按钮

### 2.8 前端 KnowledgeBaseNoteEditor.vue（**新建**）
- Markdown 编辑器（textarea + 预览，参考 EditorView）
- 标题输入
- 分类下拉（从 API 加载）
- 标签输入（逗号分隔）
- 公开/私密开关
- 保存按钮

### 2.9 前端 KbNoteView.vue（**新建**）
- 公开笔记展示页
- Markdown 渲染内容
- 显示作者、分类、标签、日期
- 路径：`/kb-note/:id`

### 2.10 前端 Router (`web/src/router/index.js`)
新增路由：
- `/user-center/knowledge-base/:id/editor` → 新建笔记
- `/user-center/knowledge-base/:id/editor/:noteId` → 编辑笔记
- `/kb-note/:id` → 公开笔记展示（无需 auth）

## 3. 执行顺序

1. 后端 model 更新 + db migration
2. 后端 service 更新（AddDocument 扩展 + UpdateDocument + GetPublicNote）
3. 后端 handler 更新（AddDocument 参数 + UpdateDocument + GetPublicNote）
4. 后端 router 更新
5. `go build ./...` 验证编译
6. 前端 API 客户端更新
7. 前端 KnowledgeBaseDetail.vue 重做
8. 前端 KnowledgeBaseNoteEditor.vue 新建
9. 前端 KbNoteView.vue 新建
10. 前端 router 更新
11. `npx vite build` 验证编译
12. 重启服务验证
