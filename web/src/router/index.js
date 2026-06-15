import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import AuthView from "../views/AuthView.vue";
import EditorView from "../views/EditorView.vue";
import DetailView from "../views/DetailView.vue";
import ArticlesView from "../views/ArticlesView.vue";
import MyArticlesView from "../views/MyArticlesView.vue";
import MyProjectsView from "../views/MyProjectsView.vue";
import ProjectEditorView from "../views/ProjectEditorView.vue";
import ProjectDetailView from "../views/ProjectDetailView.vue";
import ProjectsView from "../views/ProjectsView.vue";
import ProfileView from "../views/ProfileView.vue";
import AuthorView from "../views/AuthorView.vue";
import CollectionsView from "../views/CollectionsView.vue";
import NotificationsView from "../views/NotificationsView.vue";
import AboutView from "../views/AboutView.vue";
import InterviewView from "../views/InterviewView.vue";
import KnowledgeBaseView from "../views/KnowledgeBaseView.vue";
import KnowledgeBaseDetail from "../views/KnowledgeBaseDetail.vue";
import KnowledgeBaseNoteEditor from "../views/KnowledgeBaseNoteEditor.vue";
import KbNoteView from "../views/KbNoteView.vue";
import ApiKeysView from "../views/ApiKeysView.vue";
import NotFoundView from "../views/NotFoundView.vue";
import UserCenterLayout from "../components/UserCenterLayout.vue";
import UserCenterOverview from "../views/UserCenterOverview.vue";

const authGuard = { requiresAuth: true };

const routes = [
  { path: "/", name: "home", component: HomeView, meta: { title: "PulseBlog" } },
  { path: "/auth", name: "auth", component: AuthView, meta: { title: "登录 / 注册" } },
  { path: "/articles", name: "articles", component: ArticlesView, meta: { title: "文章" } },
  { path: "/article/:id", name: "detail", component: DetailView, meta: { title: "文章详情" } },
  { path: "/projects", name: "projects", component: ProjectsView, meta: { title: "项目" } },
  { path: "/projects/:id", name: "project-detail", component: ProjectDetailView, meta: { title: "项目详情" } },
  { path: "/interview", name: "interview", component: InterviewView, meta: { title: "AI 面试官", ...authGuard } },
  { path: "/author/:id", name: "author", component: AuthorView, meta: { title: "作者主页" } },

  {
    path: "/user-center",
    component: UserCenterLayout,
    meta: { title: "个人中心", ...authGuard },
    redirect: { name: "uc-overview" },
    children: [
      { path: "overview", name: "uc-overview", component: UserCenterOverview, meta: { title: "个人中心" } },
      { path: "articles", name: "uc-articles", component: MyArticlesView, meta: { title: "我的文章" } },
      { path: "projects", name: "uc-projects", component: MyProjectsView, meta: { title: "我的项目" } },
      { path: "knowledge-base", name: "uc-knowledge-base", component: KnowledgeBaseView, meta: { title: "知识库" } },
      { path: "knowledge-base/:id", name: "uc-knowledge-base-detail", component: KnowledgeBaseDetail, meta: { title: "知识库详情" } },
      { path: "knowledge-base/:id/editor", name: "uc-knowledge-base-editor", component: KnowledgeBaseNoteEditor, meta: { title: "新建笔记", sidebar: false } },
      { path: "knowledge-base/:id/editor/:noteId", name: "uc-knowledge-base-editor-edit", component: KnowledgeBaseNoteEditor, meta: { title: "编辑笔记", sidebar: false } },
      { path: "api-keys", name: "uc-api-keys", component: ApiKeysView, meta: { title: "API Key 管理" } },
      { path: "collections", name: "uc-collections", component: CollectionsView, meta: { title: "我的收藏" } },
      { path: "notifications", name: "uc-notifications", component: NotificationsView, meta: { title: "通知中心" } },
      { path: "profile", name: "uc-profile", component: ProfileView, meta: { title: "个人资料" } },
      { path: "about", name: "uc-about", component: AboutView, meta: { title: "关于 PulseBlog" } },
      { path: "editor", name: "uc-editor", component: EditorView, meta: { title: "写文章", sidebar: false } },
      { path: "editor/:id", name: "uc-editor-edit", component: EditorView, meta: { title: "编辑文章", sidebar: false } },
      { path: "project-editor", name: "uc-project-editor", component: ProjectEditorView, meta: { title: "新建项目", sidebar: false } },
      { path: "project-editor/:id", name: "uc-project-editor-edit", component: ProjectEditorView, meta: { title: "编辑项目", sidebar: false } },
    ]
  },

  { path: "/my-articles", redirect: "/user-center/articles" },
  { path: "/my-projects", redirect: "/user-center/projects" },
  { path: "/knowledge-base", redirect: "/user-center/knowledge-base" },
  { path: "/knowledge-base/:id", redirect: to => `/user-center/knowledge-base/${to.params.id}` },
  { path: "/api-keys", redirect: "/user-center/api-keys" },
  { path: "/collections", redirect: "/user-center/collections" },
  { path: "/notifications", redirect: "/user-center/notifications" },
  { path: "/profile", redirect: "/user-center/profile" },
  { path: "/about", redirect: "/user-center/about" },
  { path: "/editor", redirect: "/user-center/editor" },
  { path: "/editor/:id", redirect: to => `/user-center/editor/${to.params.id}` },
  { path: "/project-editor", redirect: "/user-center/project-editor" },
  { path: "/project-editor/:id", redirect: to => `/user-center/project-editor/${to.params.id}` },

  { path: "/kb-note/:id", name: "kb-note", component: KbNoteView, meta: { title: "公开笔记" } },

  { path: "/:pathMatch(.*)*", name: "not-found", component: NotFoundView, meta: { title: "页面不存在" } }
];

const router = createRouter({
  history: createWebHistory("/PulseBlog/"),
  routes,
  scrollBehavior: () => ({ top: 0 })
});

router.beforeEach((to) => {
  const token = localStorage.getItem("blog_token");
  if (to.meta?.requiresAuth && !token) {
    return { name: "auth" };
  }
  document.title = to.meta?.title && to.meta.title !== "PulseBlog" ? `${to.meta.title} | PulseBlog` : "PulseBlog";
  return true;
});

export default router;
