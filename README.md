# PulseBlog

求职导向的个人品牌站与创作平台。把你的技术积累整理成真正能对外展示的作品集。

> Go + Vue 3 全栈，支持文章/项目管理、AI 模拟面试、知识库笔记、向量检索、AI 审校。

## 功能一览

### 📝 文章系统
- 富文本编辑器 + Markdown 支持
- 草稿 → 提交 → 审核 → 发布 工作流
- 分类 / 标签 / 私有文章 / 公开浏览
- 收藏 / 点赞 / 评论

### 📦 项目管理
- 独立作品展示页（技术栈、角色、难点、成果）
- 与文章相同的审核工作流

### 🧠 知识库笔记
- Markdown 笔记编辑器，支持分类和标签
- 公开/私密切换，公开笔记可独立访问 `/kb-note/:id`
- AI 向量检索（基于 Qdrant + DeepSeek Embedding）
- 无 API Key 时自动降级，不影响笔记创建和编辑

### 🤖 AI 面试官
- 输入目标职位 + 上传简历
- DeepSeek 驱动多轮面试
- 逐题评分 + 最终评分报告

### 🛡️ AI 审校
- 敏感词检测
- 内容合规审查（文章/项目提交自动触发）

### 👤 个人中心
- 创作工作台 + 设置分区
- 创作数据统计（文章/项目/知识库/浏览量）
- 最近动态时间线

### 🔧 后台管理面板
- Dashboard 全局统计
- 文章/项目审核队列
- 用户管理、评论管理、操作日志
- 敏感词库、上传文件管理
- AI 审校结果查看

## 技术栈

| 层 | 技术 |
|---|---|
| **前端** | Vue 3 + Pinia + Vue Router + Vite |
| **后台管理** | Vue 3 + Pinia + Vue Router + Vite |
| **后端** | Go + Gin + GORM |
| **数据库** | MySQL 8.0 |
| **向量数据库** | Qdrant |
| **AI 集成** | DeepSeek API (Chat + Embedding) |
| **认证** | JWT (RSA256) |

## 架构

```
┌──────────────┐   ┌───────────────┐
│  Web 前端     │   │  Admin 后台    │
│  Vue 3 + Vite │   │  Vue 3 + Vite │
└──────┬───────┘   └───────┬───────┘
       └────────┬──────────┘
                │ JWT Auth
       ┌────────▼──────────┐   ┌──────────────┐
       │  API Server        │   │  Qdrant       │
       │  Go + Gin + GORM   │   │  向量数据库    │
       └────────┬──────────┘   └──────────────┘
                │
       ┌────────▼──────────┐
       │  MySQL 8.0         │
       └───────────────────┘
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Qdrant（可选，知识库 AI 检索需要）

### 一键启动

```powershell
.\start.ps1
```

脚本自动启动：Qdrant → MySQL → 后端(8080) → 前端(5173) → 管理后台(5174)

### 手动启动

#### 1. 后端

```bash
cp server/.env.example server/.env
# 编辑 .env 填入数据库连接
cd server
go mod download
go run cmd/api/main.go
```

#### 2. 前端

```bash
cd web
npm install
npm run dev
```

#### 3. 后台管理

```bash
cd admin
npm install
npm run dev
```

#### 4. Qdrant（可选，知识库需要）

```bash
cd qdrant
.\qdrant.exe
```

#### 5. AI 功能（可选）

在 `server/.env` 中配置：

```env
DEEPSEEK_API_KEY=sk-your-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
```

## 目录结构

```
├── server/              Go 后端 API
│   ├── cmd/api/         入口
│   ├── internal/
│   │   ├── config/      配置
│   │   ├── database/    数据库连接与迁移
│   │   ├── handler/     HTTP 处理器
│   │   ├── middleware/  中间件（JWT、限流、日志）
│   │   ├── model/       数据模型
│   │   ├── service/     业务逻辑
│   │   └── utils/       工具
│   └── Dockerfile
├── web/                 Vue 3 前端
│   └── src/
│       ├── api/         API 客户端
│       ├── views/       页面组件
│       ├── components/  公共组件
│       └── router/      路由
├── admin/               Vue 3 后台管理
│   └── src/
│       ├── views/       管理页面
│       └── components/  面板组件
├── qdrant/              Qdrant 向量数据库（本地）
├── docs/                设计和实现文档
└── docker-compose.yml
```

## License

MIT
