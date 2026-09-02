# P2 Code Quality + P3 UX Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce initial bundle size via lazy loading + code splitting, and add client-side image compression.

**Architecture:** Convert all Vue router imports to dynamic `() => import(...)` for on-demand chunking. Configure Vite `manualChunks` to separate vendor libraries into shared chunks. Add `browser-image-compression` to compress images before upload.

**Tech Stack:** Vue 3, Vue Router, Vite, browser-image-compression

---

### Task 1: Web Router Lazy Loading

**Files:**
- Modify: `web/src/router/index.js`

- [ ] **Step 1: Convert all static imports to dynamic imports**

Replace the 29 static `import` statements (lines 3-30) with lazy-loaded versions. The `useUserStore` import stays static (needed in `beforeEach`).

```javascript
import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "../stores/user";

const HomeView = () => import("../views/HomeView.vue");
const AuthView = () => import("../views/AuthView.vue");
const EditorView = () => import("../views/EditorView.vue");
const DetailView = () => import("../views/DetailView.vue");
const ArticlesView = () => import("../views/ArticlesView.vue");
const MyArticlesView = () => import("../views/MyArticlesView.vue");
const MyProjectsView = () => import("../views/MyProjectsView.vue");
const ProjectEditorView = () => import("../views/ProjectEditorView.vue");
const ProjectDetailView = () => import("../views/ProjectDetailView.vue");
const ProjectsView = () => import("../views/ProjectsView.vue");
const ProfileView = () => import("../views/ProfileView.vue");
const AuthorView = () => import("../views/AuthorView.vue");
const CollectionsView = () => import("../views/CollectionsView.vue");
const NotificationsView = () => import("../views/NotificationsView.vue");
const AboutView = () => import("../views/AboutView.vue");
const InterviewView = () => import("../views/InterviewView.vue");
const KnowledgeBaseView = () => import("../views/KnowledgeBaseView.vue");
const KnowledgeBaseDetail = () => import("../views/KnowledgeBaseDetail.vue");
const KnowledgeBaseNoteEditor = () => import("../views/KnowledgeBaseNoteEditor.vue");
const KbNoteView = () => import("../views/KbNoteView.vue");
const PublicKnowledgeBases = () => import("../views/PublicKnowledgeBases.vue");
const PublicKBDetailView = () => import("../views/PublicKBDetailView.vue");
const KBGraphView = () => import("../views/KBGraphView.vue");
const KBReadView = () => import("../views/KBReadView.vue");
const ApiKeysView = () => import("../views/ApiKeysView.vue");
const NotFoundView = () => import("../views/NotFoundView.vue");
const UserCenterLayout = () => import("../components/UserCenterLayout.vue");
const UserCenterOverview = () => import("../views/UserCenterOverview.vue");
```

Keep the route definitions and everything else exactly as-is.

- [ ] **Step 2: Verify build**

Run: `npx vite build` in `web/`
Expected: Build succeeds, multiple JS chunks produced instead of one monolithic bundle.

- [ ] **Step 3: Commit**

```bash
git add web/src/router/index.js
git commit -m "perf(web): lazy-load all route components for code splitting"
```

---

### Task 2: Admin Router Lazy Loading

**Files:**
- Modify: `admin/src/router/index.js`

- [ ] **Step 1: Convert all static imports to dynamic imports**

Replace lines 3-10 with lazy-loaded versions. The `useAdminStore` import stays static.

```javascript
import { createRouter, createWebHistory } from "vue-router";
import { useAdminStore } from "../stores/auth";

const LoginView = () => import("../views/LoginView.vue");
const DashboardView = () => import("../views/DashboardView.vue");
const ModerationHitsView = () => import("../views/ModerationHitsView.vue");
const SensitiveWordsView = () => import("../views/SensitiveWordsView.vue");
const CommentsView = () => import("../views/CommentsView.vue");
const LogsView = () => import("../views/LogsView.vue");
const UploadsView = () => import("../views/UploadsView.vue");
const SettingsView = () => import("../views/SettingsView.vue");
```

Keep route definitions and guards unchanged.

- [ ] **Step 2: Verify build**

Run: `npx vite build` in `admin/`
Expected: Build succeeds, multiple JS chunks.

- [ ] **Step 3: Commit**

```bash
git add admin/src/router/index.js
git commit -m "perf(admin): lazy-load all route components for code splitting"
```

---

### Task 3: Vite manualChunks Configuration

**Files:**
- Modify: `web/vite.config.js`
- Modify: `admin/vite.config.js`

- [ ] **Step 1: Add manualChunks to web vite.config.js**

```javascript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  base: "/PulseBlog/",
  plugins: [vue()],
  server: {
    port: 5173,
    strictPort: false,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          "vendor-vue": ["vue", "vue-router", "pinia"],
          "vendor-tiptap": [
            "@tiptap/vue-3",
            "@tiptap/starter-kit",
            "@tiptap/extension-placeholder",
            "@tiptap/extension-link",
            "@tiptap/extension-image",
            "@tiptap/extension-code-block-lowlight",
            "@tiptap/extension-task-list",
            "@tiptap/extension-task-item",
            "@tiptap/extension-table",
            "@tiptap/extension-table-row",
            "@tiptap/extension-table-cell",
            "@tiptap/extension-table-header",
            "@tiptap/extension-highlight",
            "@tiptap/extension-typography",
            "@tiptap/extension-color",
            "@tiptap/extension-text-style",
          ],
          "vendor-highlight": ["highlight.js", "lowlight"],
        }
      }
    }
  }
});
```

- [ ] **Step 2: Add manualChunks to admin vite.config.js**

```javascript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  base: "/PulseBlog/admin/",
  plugins: [vue()],
  server: {
    port: 5174,
    strictPort: false,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          "vendor-vue": ["vue", "vue-router", "pinia"],
        }
      }
    }
  }
});
```

- [ ] **Step 3: Verify both builds**

Run in `web/`: `npx vite build`
Run in `admin/`: `npx vite build`
Expected: Both succeed, multiple vendor chunks visible in output.

- [ ] **Step 4: Commit**

```bash
git add web/vite.config.js admin/vite.config.js
git commit -m "perf: configure Vite manualChunks for vendor code splitting"
```

---

### Task 4: Client-Side Image Compression

**Files:**
- Modify: `web/package.json`
- Modify: `web/src/components/RichEditor.vue`
- Create: `web/src/utils/compressImage.js`

- [ ] **Step 1: Install browser-image-compression**

Run in `web/`: `npm install browser-image-compression`

- [ ] **Step 2: Create compressImage utility**

Create `web/src/utils/compressImage.js`:

```javascript
import imageCompression from "browser-image-compression";

const OPTIONS = {
  maxSizeMB: 1,
  maxWidthOrHeight: 2048,
  useWebWorker: true,
  initialQuality: 0.8,
};

export async function compressImage(file) {
  if (!file.type.startsWith("image/")) return file;
  if (file.size <= 200 * 1024) return file;
  try {
    return await imageCompression(file, OPTIONS);
  } catch {
    return file;
  }
}
```

- [ ] **Step 3: Integrate compression into RichEditor image upload**

In `web/src/components/RichEditor.vue`, find the `addImage` method (around line 168-193) and add compression before the upload call:

```javascript
import { compressImage } from "../utils/compressImage";
```

In the `addImage` function, after `const file = event.target.files[0];` and before the `if (!file) return;`, add:

```javascript
const compressed = await compressImage(file);
```

Then change `formData.append("file", file)` to `formData.append("file", compressed)`.

- [ ] **Step 4: Verify build**

Run: `npx vite build` in `web/`
Expected: Build succeeds.

- [ ] **Step 5: Commit**

```bash
git add web/package.json web/package-lock.json web/src/utils/compressImage.js web/src/components/RichEditor.vue
git commit -m "feat: client-side image compression before upload (browser-image-compression)"
```

---

### Task 5: Final Verification

- [ ] **Step 1: Full build check**

```bash
cd web && npx vite build
cd admin && npx vite build
```

- [ ] **Step 2: Go backend build**

```bash
cd server && go build ./cmd/api
```

- [ ] **Step 3: Go tests**

```bash
cd server && go test ./internal/handler -short -count=1 -v
```

- [ ] **Step 4: Commit any fixes, push**

```bash
git add -A
git commit -m "chore: P2+P3 verification pass"
git push
```
