# PulseBlog V2.0 — 部署上线工程化 PRD

## 1. 文档定义

| 字段 | 内容 |
|------|------|
| 产品名称 | PulseBlog |
| 文档版本 | V2.0 |
| 文档类型 | 工程化 PRD — 聚焦部署上线与生产就绪度 |
| 前置文档 | `PRD_V1_个人作品集求职导向内容平台.md` — V1.0 功能闭环已完成 |
| 当前状态 | 功能完整，本地可运行，但缺少生产环境适配层 |
| 文档目标 | 让项目可从"本地 dev 可运行"升级为"生产环境可部署、可运维、可迭代" |

## 2. 当前状态摘要

### 2.1 已完成（V1.0 功能闭环）

- **前端 Web**：16 个页面全部实现（首页、文章/项目列表与详情、编辑器、作者页、个人资料、通知、收藏、404 等），构建通过
- **前端 Admin**：7 个页面/视图，核心管理功能全部实现在 Dashboard 中（文章/项目管理、审核队列、用户/评论/敏感词/日志/资源管理、分类/标签增删改），构建通过
- **后端 Go**：约 55+ 个 API 接口，涵盖认证鉴权、文章/项目 CRUD 与审核流、评论树、通知、风控敏感词、批量操作、分类标签管理。编译通过，`go vet` 通过
- **测试**：3 个 service 层测试文件（review_flow, meta_service, notification_service），全部通过
- **管理工作**：31 轮任务已完成，包括全量 mojibake 清理、批量审核操作、分类标签治理、审核队列体验优化

### 2.2 当前的上线阻塞点

1. **硬编码 API 地址**：web/admin 的 `client.js` 中 `baseURL: "http://localhost:8080/api"` 写死，无环境变量替换机制
2. **无部署基础设施**：没有 Dockerfile、nginx 配置、静态文件服务策略、CI/CD
3. **安全默认值未改**：`JWT_SECRET=change-me-secret` 在生产环境是严重风险
4. **文件上传无限制**：缺少文件大小和速率限制
5. **CORS 仅限 localhost**：无法在生产域名下工作
6. **Admin 大单体**：DashboardView.vue 2249 行，长期维护困难
7. **测试覆盖薄弱**：仅 service 层 3 个测试，无 handler 集成测试、无 e2e
8. **错误处理脆弱**：使用字符串比较 `err.Error() == "xxx"` 做逻辑分支

## 3. 总体方案

三阶段推进，每阶段可独立交付，互不阻塞。

```
Phase 1 — 生产硬化 (生产可运行)
  ├─ DEP-01: 环境变量替换 API baseURL
  ├─ DEP-02: Docker + docker-compose 部署方案
  ├─ DEP-03: Go 静态文件嵌入（可选方案）
  ├─ DEP-04: Nginx 反向代理配置
  ├─ DEP-05: 安全加固（JWT 密钥 / 上传限制 / CORS）
  └─ 验收：项目可在 Docker 中一键启动，浏览器访问可用

Phase 2 — 工程健壮性 (可维护)
  ├─ DEP-06: Admin Dashboard 组件拆分
  ├─ DEP-07: 错误处理改用 sentinel errors
  ├─ DEP-08: 补充集成测试与 handler 测试
  └─ 验收：构建通过，测试覆盖率明显提升

Phase 3 — 运维与体验 (可运维)
  ├─ DEP-09: 结构化日志接入
  ├─ DEP-10: 前端构建时环境变量支持
  ├─ DEP-11: CI/CD 最小配置（GitHub Actions）
  ├─ DEP-12: 部署文档与 quick-start 脚本
  └─ 验收：可走通完整 CI → build → deploy 流程
```

## 4. 详细任务包

---

### Phase 1 — 生产硬化

#### DEP-01：API baseURL 环境变量化

**目标**：消除 web 和 admin 中硬编码的 `localhost:8080`，使构建产物可指向任意后端地址。

**范围**：前端

**输入文件**：
- `web/src/api/client.js`
- `admin/src/api/client.js`
- `web/vite.config.js`
- `admin/vite.config.js`

**改动内容**：
1. 在 `web/vite.config.js` 和 `admin/vite.config.js` 中定义 `VITE_API_URL` 环境变量，默认值保持 `http://localhost:8080/api`
2. 修改 `client.js`：`baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api"`
3. 提供 `.env.development` / `.env.production` 模板文件

**验收标准**：
- 不传 `VITE_API_URL` 时构建仍使用默认 localhost
- 传 `VITE_API_URL=https://api.example.com` 构建后产物请求正确地址
- `cd web && npm run build` 通过
- `cd admin && npm run build` 通过

**依赖**：无

---

#### DEP-02：Docker + docker-compose 部署方案

**目标**：提供一键部署能力，包含 MySQL + Go API + Web 静态文件分发。

**范围**：基础设施

**输入文件**：
- `server/Dockerfile`（新建）
- `docker-compose.yml`（项目根目录，新建）
- `.dockerignore`（新建）

**改动内容**：
1. **Go 服务 Dockerfile**：多阶段构建
   - 阶段一（build）：`golang:1.26-alpine`，编译 `cmd/api/main.go`
   - 阶段二（runtime）：`alpine:latest`，复制二进制 + uploads 目录
   - 暴露 8080 端口
   - CMD 运行编译后的二进制
2. **docker-compose.yml**：
   - `db` 服务：MySQL 8.0，挂载 volume 持久化，初始化数据库
   - `api` 服务：Go 二进制，依赖 db，映射 8080 端口
   - 通过环境变量注入 DB 配置、JWT_SECRET 等
   - 网络：`app-network` 内部通信
3. 创建 `.env.prod` 模板，docker-compose 引用

**验收标准**：
- `docker-compose up --build` 一键启动成功
- Go 服务启动后自动执行 migration + seed
- API 可通过 `http://localhost:8080` 访问
- `docker-compose down` 正常停止
- 日志输出可查看

**依赖**：无

---

#### DEP-03：Go 静态文件嵌入（整合方案）

**目标**：为不需要额外 nginx 的场景提供单二进制部署方案，Go 服务直接 serve web/admin 的构建产物。

> 注：与 DEP-04 二选一，或同时提供让用户自己选。推荐二选一，默认为 DEP-04 nginx 方案更灵活。

**范围**：后端

**输入文件**：
- `server/cmd/api/main.go`
- `server/internal/router/router.go`
- `web/dist/`（构建产物）
- `admin/dist/`（构建产物）

**改动内容**：
1. 使用 Go 1.16+ 的 `embed` 特性，将 web/dist 和 admin/dist 嵌入二进制
2. 在 router 中添加路由：
   - `/admin/*` → 从嵌入的 admin 文件系统 serve
   - `/` → 从嵌入的 web 文件系统 serve
   - API 路径保留 `/api/*` 前缀
3. SPA fallback：对非 API 路径返回 index.html（处理前端路由刷新）

**验收标准**：
- 单二进制文件可直接运行，无需外部静态文件目录
- 访问 `/` 展示 Web 前端
- 访问 `/admin/` 展示 Admin 前端
- 前端路由刷新（如 `/articles`）正常返回页面
- API 请求正常响应

**依赖**：DEP-01（确保 baseURL 在构建时已正确配置）

---

#### DEP-04：Nginx 反向代理配置

**目标**：提供标准生产架构——nginx 在前端处理 TLS 终止、静态文件分发、API 反向代理。

**范围**：基础设施

**输入文件**：
- `deploy/nginx.conf`（新建）
- `deploy/start.sh`（新建）

**改动内容**：
1. 编写 nginx 配置，包含：
   - 监听 443（HTTPS），HTTP 301 跳转 HTTPS
   - `/api/` 反向代理到 Go 服务 `http://api:8080`
   - `/admin/` 指向 admin 构建产物
   - `/` 指向 web 构建产物
   - 静态文件缓存头（`Cache-Control: max-age=3600`）
   - SPA fallback：`try_files $uri $uri/ /index.html`
   - 上传文件访问：`/uploads/` → Go 服务
2. 提供 `start.sh` 启动脚本
3. 支持通过环境变量替换域名

**验收标准**：
- nginx 启动后所有路由正常
- /api 请求正确转发到 Go 服务
- SPA 路由刷新正常
- 证书可通过环境变量注入或使用 Let's Encrypt

**依赖**：DEP-02（docker 环境）

---

#### DEP-05：安全加固

**目标**：修复生产环境安全风险。

**范围**：后端 + 配置

**输入文件**：
- `server/.env.example`
- `server/internal/handler/upload_handler.go`
- `server/internal/middleware/auth.go`
- `server/internal/router/router.go`

**改动内容**：
1. **JWT_SECRET**：在 `.env.example` 中添加注释 `# 生产环境请替换为至少 32 位随机字符串`；添加 `main.go` 启动时对 weak secret 的警告
2. **上传限制**：在 `UploadImage` 中增加 `c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)`（10MB）
3. **CORS 生产域名**：将允许来源从硬编码改为从环境变量读取逗号分隔的多个 origin；默认保持 localhost 供开发
4. **上传路径穿越防护**：验证 `filepath.Clean(dst)` 以 `/uploads/` 为前缀

**验收标准**：
- 上传超过 10MB 的文件返回明确的 413 错误
- CORS 配置可通过环境变量指定生产域名
- 生成环境中 nginx 或启动日志会提示 JWT 为默认值时的安全警告
- `cd server && go build ./cmd/api/` 通过

**依赖**：无

---

### Phase 2 — 工程健壮性

#### DEP-06：Admin Dashboard 组件拆分

**目标**：将 2249 行的 DashboardView.vue 拆分为可独立维护的子组件，保持功能不变。

**范围**：前端 Admin

**输入文件**：
- `admin/src/views/DashboardView.vue`
- `admin/src/api/dashboard.js`

**改动内容**：
1. 拆分为以下组件（每个文件 < 400 行）：
   - `DashboardOverview.vue` — 顶部统计卡片 + 快捷入口
   - `ArticleReviewPanel.vue` — 文章待审核面板（当前 >800 行）
   - `ProjectReviewPanel.vue` — 项目待审核面板
   - `ArticleManagementSection.vue` — 文章管理列表 + 批量操作
   - `ProjectManagementSection.vue` — 项目管理列表
   - `UserManagementSection.vue` — 用户管理列表
   - `CategoryTagSection.vue` — 分类/标签管理面板
2. DashboardView.vue 降级为容器组件，只负责布局和状态协调
3. 使用 Pinia store（或 provide/inject）在子组件间共享状态

**验收标准**：
- 拆分后所有功能与原版完全一致
- Admin 构建通过
- 审核、管理、批量操作等均正常工作
- 无样式回归

**依赖**：无

---

#### DEP-07：错误处理改用 sentinel errors

**目标**：消除 `err.Error() == "字符串"` 的脆弱模式，改用类型安全的方式。

**范围**：后端 Go

**输入文件**：
- `server/internal/service/article_service.go`
- `server/internal/service/project_service.go`
- `server/internal/service/auth_service.go`
- `server/internal/handler/*.go`（所有使用字符串比较的地方）

**改动内容**：
1. 在 `service` 包中定义 sentinel errors：

```go
var (
    ErrNotPublished     = errors.New("项目暂未公开")
    ErrNoPermission     = errors.New("没有权限操作该项目")
    ErrNotPending       = errors.New("项目当前不在待审核状态")
    ErrUserBanned       = errors.New("已被封禁的用户不能管理项目")
    ErrRejectReasonReq  = errors.New("驳回时必须填写原因")
    ErrInvalidAction    = errors.New("无效的审核操作")
    ErrTitleContentReq  = errors.New("项目标题和正文不能为空")
    // ... 其他
)
```

2. 将所有 `return errors.New("xxx")` / `err.Error() == "xxx"` 替换为 sentinel
3. handler 层使用 `errors.Is()` / `errors.As()` 判断错误类型

**验收标准**：
- 代码中不再出现 `err.Error() == "xxx"` 模式
- 所有 sentinel 错误在 handler 中正确映射到 HTTP 状态码
- `go vet ./...` 通过
- `go test ./...` 通过

**依赖**：无

---

#### DEP-08：补充集成测试

**目标**：将测试覆盖从 service 层扩展到 handler 和关键流程。

**范围**：后端 Go

**输入文件**：
- `server/internal/handler/*_test.go`（新建）
- `server/internal/service/*_test.go`（扩展）

**改动内容**：
1. **Handler 集成测试**（`handler_test.go`）：
   - 使用 `httptest.NewRecorder()` + `gin.SetMode(gin.TestMode)`
   - 测试场景：
     - 注册 → 登录 → 获取 me
     - 创建文章 → 提交审核 → 管理员审核通过
     - 创建项目 → 提交审核 → 管理员驳回
     - 评论 → 回复 → 通知触发
     - 未登录访问受保护接口 → 401
     - 普通用户访问管理接口 → 403
2. **Service 层扩展测试**：
   - `article_service_test.go`：文章状态流转测试（draft→pending→published/rejected→draft）
   - 权限边界：被封禁用户不能创建/提交内容

**验收标准**：
- 新增至少 3 个 handler 测试文件
- service 层新增至少 2 个测试文件
- 核心审核流程 + 权限校验被测试覆盖
- `go test ./...` 全部通过

**依赖**：DEP-07（推荐，sentinel 错误让 handler 测试更容易）或可并行进行

---

### Phase 3 — 运维与体验

#### DEP-09：结构化日志接入

**目标**：从 `log.Printf` 升级到结构化日志，支持 JSON 输出、日志级别、请求追踪。

**范围**：后端 Go

**输入文件**：
- `server/cmd/api/main.go`
- `server/internal/handler/*.go`
- `server/internal/middleware/auth.go`

**改动内容**：
1. 接入 `go.uber.org/zap` 或 `github.com/rs/zerolog`
2. 添加 request ID middleware（`X-Request-ID`），在日志中传递
3. 关键日志点：
   - 请求入口（method, path, duration, status, request_id）
   - 审核操作（operator, action, target_type, target_id）
   - 风控命中（user, scene, field, word）
   - 错误（带 stack trace，仅 error 级别）
4. config 中增加 `LOG_LEVEL` 和 `LOG_FORMAT`（json/text）

**验收标准**：
- 启动日志、请求日志、审核日志、错误日志均为 JSON 格式
- 日志中包含 request_id，可按 ID 串联一次请求
- 可通过环境变量切换日志级别
- `go build` 通过，测试通过

**依赖**：无

---

#### DEP-10：前端构建时环境变量完整支持

**目标**：使前端构建支持开发/测试/生产三套配置，且构建时不暴露密钥。

**范围**：前端 Web + Admin

**输入文件**：
- `web/.env.development`（新建）
- `web/.env.production`（新建）
- `admin/.env.development`（新建）
- `admin/.env.production`（新建）
- `web/vite.config.js`
- `admin/vite.config.js`

**改动内容**：
1. 创建 `.env.development`（开发默认值）和 `.env.production`（生产默认值）
2. `VITE_API_URL` 根据环境自动切换
3. 添加 `VITE_APP_TITLE` 用于页面标题
4. 添加 `VITE_SENTRY_DSN`（可选，预留）
5. 在构建日志中输出当前环境标识

**验收标准**：
- `npm run dev` 使用 development 配置，API 指向 localhost
- `npm run build` 使用 production 配置
- 可通过环境变量覆盖构建时配置
- 构建产物不包含明文密钥（VITE_* 前缀变量已内联）

**依赖**：DEP-01

---

#### DEP-11：CI/CD 最小配置

**目标**：提供开箱即用的 GitHub Actions 配置，PR 自动 lint/build/test，push 到 main 自动部署。

**范围**：基础设施

**输入文件**：
- `.github/workflows/ci.yml`（新建）
- `.github/workflows/deploy.yml`（新建）

**改动内容**：
1. **CI workflow**（pull_request / push 触发）：
   - 并行 job：
     - web: `npm ci` → `npm run build`
     - admin: `npm ci` → `npm run build`
     - server: `go build ./...` → `go vet ./...` → `go test ./...`
   - 缓存 node_modules 和 Go module cache
2. **Deploy workflow**（可选，push 到 main 触发）：
   - 构建 Docker 镜像
   - push 到 Docker Hub / GitHub Container Registry
   - 通过 SSH / webhook 触发远程部署

**验收标准**：
- CI 在 push 和 PR 时自动触发
- 三个项目全部构建通过
- 测试通过
- 缓存生效，总运行时间 < 3 分钟

**依赖**：DEP-02（Docker 镜像构建）

---

#### DEP-12：部署文档与 Quick-Start 脚本

**目标**：让一个不熟悉项目的人能在 15 分钟内从 git clone 到 running。

**范围**：文档 + 脚本

**输入文件**：
- `DEPLOYMENT.md`（新建，项目根目录）
- `scripts/setup.sh`（新建）
- `scripts/start-prod.sh`（新建）

**改动内容**：
1. **DEPLOYMENT.md** 包含：
   - 架构说明（nginx → Go API → MySQL）
   - 两种部署方式：Docker 一键部署 / 裸机手动部署
   - 环境变量对照表与说明
   - 安全 checklist（改密码、改 JWT secret、开启 HTTPS）
   - 监控与日志查看方式
   - 常见问题排查
2. **setup.sh**：环境检查（docker、node、go 版本）
3. **start-prod.sh**：构建前端 → docker-compose up 的快捷脚本

**验收标准**：
- 按文档操作可成功部署
- setup.sh 可检查所有必备依赖
- 文档已包含 TLS 证书获取指导

**依赖**：DEP-02、DEP-04

## 5. 工程化要求

### 5.1 编码规范

- 前端代码保持 Vue 3 `<script setup>` + Options API 风格
- Go 代码保持现有分层（handler → service → model）
- 所有新增文件使用 UTF-8 编码
- 新增组件文件名使用 PascalCase

### 5.2 接口规范

- 保持现有统一响应结构：`{items, pagination}` / `{item}` / `{message}`
- 新增 API 遵循相同模式
- 错误响应保持：`{message, code?, details?}`

### 5.3 安全规范

- 禁止在 frontend 环境变量中存放任何密钥
- 所有敏感配置通过后端环境变量注入
- 上传文件类型白名单在前端 + 后端双重校验
- CORS 生产环境必须指定具体域名，不允许 `*`

### 5.4 测试要求

- 新增 handler 测试尽量覆盖正常 + 异常路径
- 状态机测试覆盖所有合法转换
- 测试不得依赖外部服务（数据库使用 SQLite 内存模式或 mock）

## 6. 里程碑与排期

| 里程碑 | 包含任务 | 预估工时 | 交付物 |
|--------|---------|---------|--------|
| M1: 本地可部署 | DEP-01, DEP-02, DEP-05 | 2-3 天 | docker-compose 一键运行 |
| M2: 架构完备 | DEP-03 或 DEP-04, DEP-10 | 1-2 天 | nginx/embed 方案就绪 |
| M3: 工程健壮 | DEP-06, DEP-07, DEP-08 | 3-4 天 | 代码质量达标 |
| M4: 可运维 | DEP-09, DEP-11, DEP-12 | 2 天 | CI 绿色 + 文档完备 |

**总预估**：8-11 天（1.5-2 周，单人全职）

## 7. 风险与应对

| 风险 | 概率 | 影响 | 应对方案 |
|------|------|------|---------|
| Go embed 导致二进制体积过大 | 低 | 中 | 默认推荐 nginx 方案，embed 作为备选 |
| Admin 拆组件导致 Dashboard 功能回归 | 中 | 高 | 逐个组件替换，每步构建验证 + 手动冒烟 |
| 测试数据库环境搭建复杂 | 中 | 低 | 使用 SQLite 内存模式做 handler 测试 |
| Docker 在 Windows 上兼容问题 | 中 | 低 | 提供 裸机部署指南作为 fallback |
| JWT secret 泄露 | 低 | 高 | 启动时检查 + 部署文档红色警告 |

## 8. 交付标准

### 8.1 必须满足

- 三项目构建均通过（web, admin, server）
- `go vet ./...` 无 error
- `go test ./...` 全部通过
- 可通过 `docker-compose up` 一键启动并在浏览器中完整使用

### 8.2 建议满足

- 测试覆盖率 service 层 > 60%，handler 层 > 40%
- CI 运行时间 < 3 分钟
- Admin Dashboard 拆分后每个组件 < 400 行

### 8.3 最终判定

当满足以下所有条件时，项目达到"可部署上线"标准：

1. ✅ 一个基础环境（Docker 或 裸机）可在 15 分钟内从 git clone 到浏览器可用
2. ✅ API baseURL 通过环境变量配置，构建产物不包含 localhost 硬编码
3. ✅ 安全默认值已处理（JWT secret、上传限制、CORS）
4. ✅ CI 在 PR 提交时自动执行 lint + build + test
5. ✅ 有清晰的部署文档，包含安全 checklist
