# PulseBlog

求职导向的个人作品集与博客平台。支持文章/项目管理、AI 审校、AI 模拟面试、后台管理面板。

## 架构

```
┌─────────────┐  ┌──────────────┐
│  Web 前端    │  │  Admin 后台  │  Vue 3 + Vite
│  localhost   │  │  localhost   │
│  :5173       │  │  :5174       │
└──────┬───────┘  └──────┬───────┘
       └────────┬────────┘
                │ JWT Auth
       ┌────────▼────────┐
       │  API Server      │  Go / Gin + GORM
       │  localhost:8080  │
       └────────┬────────┘
                │
       ┌────────▼────────┐
       │  MySQL          │
       └─────────────────┘
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 后端

```bash
cp server/.env.example server/.env
# 编辑 .env 填入数据库连接
cd server
go mod download
go run cmd/api/main.go
```

### 2. 前端

```bash
cd web
npm install
npm run dev        # → http://localhost:5173
```

### 3. 后台管理

```bash
cd admin
npm install
npm run dev        # → http://localhost:5174
```

### 4. AI 面试官（可选）

在 `server/.env` 中配置：

```env
DEEPSEEK_API_KEY=sk-your-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
```

## 功能

| 功能 | 说明 |
|---|---|
| **文章管理** | 富文本编辑器、草稿/提交/审核/发布工作流、标签分类 |
| **项目管理** | 项目展示页、同文章工作流 |
| **用户系统** | 注册/登录、JWT 认证、个人资料、收藏/点赞 |
| **AI 审校** | 自动敏感词检测、内容合规审查 |
| **AI 面试官** | 职位输入 + 简历上传 → DeepSeek 驱动的一问一答语音面试 → 评分报告 |
| **通知中心** | 审核结果、系统通知 |
| **后台面板** | Dashboard 统计、内容审核队列、用户/评论/日志管理、敏感词库、上传管理 |

## 部署

```bash
docker-compose up -d
```

项目包含 GitHub Actions CI/CD 配置（`.github/workflows/`）。

## 目录结构

```
├── server/           Go 后端 API
│   ├── cmd/api/      入口
│   ├── internal/
│   │   ├── config/   配置
│   │   ├── database/ 数据库连接与迁移
│   │   ├── handler/  HTTP 处理器
│   │   ├── middleware/ 中间件（JWT、限流、日志）
│   │   ├── model/    数据模型
│   │   ├── service/  业务逻辑
│   │   └── utils/    工具（JWT、密码）
│   └── Dockerfile
├── web/              Vue 3 前端
│   └── src/
│       ├── api/      API 客户端
│       ├── views/    页面组件
│       └── router/   路由
├── admin/            Vue 3 后台管理
│   └── src/
│       ├── views/    管理页面
│       └── components/ 面板组件
├── deploy/           部署配置（Nginx）
├── scripts/          脚本
└── docker-compose.yml
```
