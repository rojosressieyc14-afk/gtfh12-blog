<div align="center">

# PulseBlog

**求职导向的个人品牌站与创作平台**

把你的技术积累整理成真正能对外展示的作品集。

<br/>

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vue.js&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-5-646CFF?style=flat-square&logo=vite&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat-square&logo=mysql&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

</div>

---

## Features

<table>
<tr>
<td width="50%">

### Content System
- Rich text editor + Markdown support
- Draft → Submit → Review → Publish workflow
- Categories / Tags / Private articles
- Likes / Favorites / Comments

### Project Showcase
- Independent project pages (tech stack, role, challenges, results)
- Same review workflow as articles

### Knowledge Base
- Markdown note editor with categories & tags
- Public / Private toggle
- AI vector search (Qdrant + DeepSeek Embedding)
- Graceful degradation without API key

</td>
<td width="50%">

### AI Interview Coach
- Input target position + upload resume
- DeepSeek-powered multi-round interview
- Per-question scoring + final report

### AI Moderation
- Sensitive word detection (trie-based)
- Content compliance review (auto-triggered on submit)
- Auto-ban on repeated violations

### Admin Panel
- Dashboard with global statistics
- Article / Project review queues
- User management, comments, audit logs
- Sensitive word library, upload management

</td>
</tr>
</table>

---

## Architecture

```
                    ┌─────────────────────────────────────┐
                    │           Client Layer               │
                    │  ┌──────────┐    ┌────────────────┐  │
                    │  │   Web    │    │     Admin      │  │
                    │  │ Vue 3    │    │     Vue 3      │  │
                    │  │ + Vite   │    │     + Vite     │  │
                    │  └────┬─────┘    └───────┬────────┘  │
                    └───────┼──────────────────┼───────────┘
                            │                  │
                            │    JWT (RSA256)  │
                            │    + CSRF Cookie │
                    ┌───────▼──────────────────▼───────────┐
                    │           API Server                  │
                    │     Go + Gin + GORM                   │
                    │                                       │
                    │  ┌─────────┐  ┌───────────────────┐  │
                    │  │  Auth   │  │  Moderation (Trie) │  │
                    │  │  Rate   │  │  Embedding (AI)    │  │
                    │  │  Limit  │  │  Vector Search     │  │
                    │  └─────────┘  └───────────────────┘  │
                    └───────┬──────────────┬───────────────┘
                            │              │
                    ┌───────▼──────┐ ┌─────▼──────┐
                    │   MySQL 8.0  │ │   Qdrant   │
                    │   (Primary)  │ │  (Vector)  │
                    └──────────────┘ └────────────┘
```

---

## Quick Start

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Backend API |
| Node.js | 18+ | Frontend build |
| MySQL | 8.0+ | Primary database |
| Qdrant | Latest | Vector search (optional) |

### One-Click Launch

```powershell
.\start.ps1
```

Auto-starts: Qdrant → MySQL → Backend (:8080) → Frontend (:5173) → Admin (:5174)

### Manual Setup

<details>
<summary><strong>Backend</strong></summary>

```bash
cd server
cp .env.example .env   # Edit with your DB config
go mod download
go run ./cmd/api
```

</details>

<details>
<summary><strong>Frontend</strong></summary>

```bash
cd web
npm install
npm run dev
```

</details>

<details>
<summary><strong>Admin Panel</strong></summary>

```bash
cd admin
npm install
npm run dev
```

</details>

<details>
<summary><strong>AI Features (Optional)</strong></summary>

Add to `server/.env`:

```env
DEEPSEEK_API_KEY=sk-your-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
```

</details>

---

## Project Structure

```
pulseblog/
├── server/                  # Go backend
│   ├── cmd/api/             # Entry point
│   └── internal/
│       ├── config/          # Configuration
│       ├── database/        # DB connection & migrations
│       ├── handler/         # HTTP handlers
│       ├── middleware/       # Auth, rate limit, logging
│       ├── model/           # Data models
│       └── service/         # Business logic
│
├── web/                     # Vue 3 frontend
│   └── src/
│       ├── api/             # API client
│       ├── views/           # Page components
│       ├── components/      # Shared components
│       └── router/          # Route definitions
│
├── admin/                   # Vue 3 admin panel
│   └── src/
│       ├── views/           # Admin pages
│       ├── components/      # Dashboard components
│       ├── composables/     # Shared composables
│       └── utils/           # Utility functions
│
├── qdrant/                  # Local vector DB
├── archive/                 # Archived docs & logs
└── docker-compose.yml
```

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Frontend** | Vue 3 + Pinia + Vue Router | SPA with state management |
| **Build** | Vite 5 | Fast dev server & bundler |
| **Backend** | Go + Gin | High-performance HTTP framework |
| **ORM** | GORM | Type-safe database operations |
| **Database** | MySQL 8.0 | Primary data store |
| **Vector DB** | Qdrant | AI-powered semantic search |
| **AI** | DeepSeek API | Chat + Embedding models |
| **Auth** | JWT (RSA256) + CSRF | Secure double-submit cookie pattern |

---

## Performance Optimizations

Recent audit-driven improvements:

| Optimization | Impact |
|-------------|--------|
| Trie-based sensitive word matching | O(n) single-pass scan vs O(n×m) per word |
| Batch tag synchronization | N queries → 1 IN query + bulk insert |
| Batch reaction summary loading | 4N queries → 4 queries per list page |
| Dashboard stats aggregation | 11 COUNT queries → 3聚合 queries |
| Rate limiter lifecycle | Proper goroutine cleanup with Stop() |

---

## Environment Variables

```env
# Database
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=blog_system

# Auth
JWT_SECRET=your-jwt-secret
DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_PASSWORD=your-admin-password

# CORS
WEB_ORIGIN=http://localhost:5173
ADMIN_ORIGIN=http://localhost:5174

# AI (Optional)
DEEPSEEK_API_KEY=sk-your-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
QDRANT_ADDR=http://localhost:6333

# App
GIN_MODE=debug
UPLOAD_DIR=./uploads
```

---

## License

[MIT](LICENSE)
