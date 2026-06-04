# PulseBlog 部署指南

## 环境要求

- Docker 24+
- Docker Compose v2+
- Node.js 20+（仅本地开发/调试时需要）
- Go 1.23+（仅本地开发/调试时需要）

## 快速启动（Docker 一键部署）

### 1. 构建前端

```bash
cd web && npm ci && npm run build && cd ..
cd admin && npm ci && npm run build && cd ..
```

### 2. 配置环境变量

复制环境变量模板并修改：

```bash
cp .env.prod .env
```

必须修改的配置项：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_PASSWORD` | MySQL 密码 | `please_change_this_password` |
| `JWT_SECRET` | JWT 签名密钥（至少 32 位随机字符串） | `change-me-secret` |
| `DEFAULT_ADMIN_PASSWORD` | 管理员密码 | `please_change_admin_password` |
| `WEB_ORIGIN` | Web 前端部署域名 | `http://localhost:3000` |
| `ADMIN_ORIGIN` | Admin 前端部署域名 | `http://localhost:3001` |

### 3. 启动服务

```bash
docker-compose up --build -d
```

### 4. 访问

- Web 前端：`http://localhost`
- Admin 后台：`http://localhost/admin`
- API 接口：`http://localhost/api`

### 5. 停止

```bash
docker-compose down
```

## 裸机手动部署

### 1. 构建前端

```bash
cd web && npm ci && npm run build && cd ..
cd admin && npm ci && npm run build && cd ..
```

### 2. 配置数据库

确保 MySQL 8.0+ 已运行，创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS blog_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置后端

```bash
cd server
cp .env.example .env
# 编辑 .env，修改数据库连接信息、JWT_SECRET 等
```

### 4. 启动后端

```bash
cd server && go run ./cmd/api
```

### 5. 配置 Nginx

参考 `deploy/nginx.conf`，将静态文件目录指向 `web/dist` 和 `admin/dist`。

## 安全 Checklist

部署上线前请逐项确认：

- [ ] `JWT_SECRET` 已改为 32 位以上随机字符串
- [ ] `DB_PASSWORD` 已改为强密码
- [ ] `DEFAULT_ADMIN_PASSWORD` 已修改
- [ ] 上传文件大小限制已确认（默认 10MB）
- [ ] CORS 已配置生产域名（`CORS_ORIGINS` 环境变量）
- [ ] HTTPS 已开启（建议使用 Let's Encrypt）
- [ ] 数据库端口未对外暴露（Docker 容器间通信）
- [ ] 日志格式已切换为 JSON 便于日志收集

## 环境变量参考

| 变量 | 说明 | 开发默认值 | 生产建议 |
|------|------|-----------|---------|
| `SERVER_PORT` | API 监听端口 | `8080` | `8080` |
| `DB_HOST` | 数据库地址 | `127.0.0.1` | Docker 内使用 `db` |
| `DB_PORT` | 数据库端口 | `3306` | `3306` |
| `DB_USER` | 数据库用户 | `root` | 创建专用用户 |
| `DB_PASSWORD` | 数据库密码 | 空 | 强密码 |
| `DB_NAME` | 数据库名 | `blog_system` | `blog_system` |
| `JWT_SECRET` | JWT 密钥 | `change-me-secret` | 32+ 位随机串 |
| `DEFAULT_ADMIN_USERNAME` | 默认管理员 | `admin` | `admin` |
| `DEFAULT_ADMIN_PASSWORD` | 管理员密码 | `admin123` | 强密码 |
| `WEB_ORIGIN` | Web 来源 | `http://localhost:5173` | 实际域名 |
| `ADMIN_ORIGIN` | Admin 来源 | `http://localhost:5174` | 实际域名 |
| `CORS_ORIGINS` | 额外 CORS 来源 | 空 | 逗号分隔域名 |
| `GIN_MODE` | Gin 运行模式 | `release` | `release` |
| `LOG_LEVEL` | 日志级别 | `info` | `warn` |
| `LOG_FORMAT` | 日志格式 | `text` | `json` |
| `UPLOAD_DIR` | 上传目录 | `./uploads` | `/app/uploads` |

## 架构

```
                        ┌─────────────┐
                        │   Browser   │
                        └──────┬──────┘
                               │
                        ┌──────▼──────┐
                        │    Nginx    │
                        │  :80        │
                        └──┬───┬──────┘
                           │   │
              ┌────────────┘   └────────────┐
              │                              │
       ┌──────▼──────┐              ┌───────▼───────┐
       │  Static     │              │  API Proxy    │
       │  /  & /admin │              │  /api/*       │
       └─────────────┘              └───────┬───────┘
                                            │
                                     ┌──────▼──────┐
                                     │  Go Service │
                                     │  :8080      │
                                     └──┬───┬──────┘
                                        │   │
                               ┌────────┘   └────────┐
                               │                      │
                        ┌──────▼──────┐      ┌───────▼───────┐
                        │   MySQL     │      │  Uploads     │
                        │   :3306     │      │  Storage     │
                        └─────────────┘      └───────────────┘
```

## CI/CD

项目包含 GitHub Actions 配置：

- **CI**（`.github/workflows/ci.yml`）：push/PR 时自动构建 + 测试
- **Deploy**（`.github/workflows/deploy.yml`）：push 到 main 时构建 Docker 镜像

## 日志

启动后日志输出到 stdout，格式由 `LOG_FORMAT` 控制：

```json
{"level":"info","time":"2026-01-01T12:00:00+08:00","caller":"main.go:42","msg":"server listening","port":"8080"}
```

请求日志包含 `request_id` 字段，可按 ID 串联单次请求。

## 常见问题

**Q: docker-compose up 启动后服务无法访问？**
A: 检查前端 dist 目录是否存在（需先执行 `npm run build`），检查 nginx 配置是否正确。

**Q: 数据库连接失败？**
A: 确认 MySQL 容器健康，检查 `DB_HOST` 是否设为 `db`（Docker 内部），检查密码配置。

**Q: 上传文件超过大小限制？**
A: 后端限制 10MB，nginx 限制 20MB。如需修改，同时调整 `upload_handler.go` 中的 `MaxBytesReader` 和 `deploy/nginx.conf` 中的 `client_max_body_size`。
