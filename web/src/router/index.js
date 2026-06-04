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
import NotFoundView from "../views/NotFoundView.vue";

const routes = [
  { path: "/", name: "home", component: HomeView, meta: { title: "PulseBlog" } },
  { path: "/about", name: "about", component: AboutView, meta: { title: "关于 PulseBlog" } },
  { path: "/auth", name: "auth", component: AuthView, meta: { title: "登录 / 注册" } },
  { path: "/articles", name: "articles", component: ArticlesView, meta: { title: "文章" } },
  { path: "/editor", name: "editor", component: EditorView, meta: { title: "写文章" } },
  { path: "/editor/:id", name: "editor-edit", component: EditorView, meta: { title: "编辑文章" } },
  { path: "/article/:id", name: "detail", component: DetailView, meta: { title: "文章详情" } },
  { path: "/projects", name: "projects", component: ProjectsView, meta: { title: "项目" } },
  { path: "/projects/:id", name: "project-detail", component: ProjectDetailView, meta: { title: "项目详情" } },
  { path: "/project-editor", name: "project-editor", component: ProjectEditorView, meta: { title: "新建项目" } },
  { path: "/project-editor/:id", name: "project-editor-edit", component: ProjectEditorView, meta: { title: "编辑项目" } },
  { path: "/my-articles", name: "my-articles", component: MyArticlesView, meta: { title: "我的文章" } },
  { path: "/my-projects", name: "my-projects", component: MyProjectsView, meta: { title: "我的项目" } },
  { path: "/collections", name: "collections", component: CollectionsView, meta: { title: "我的收藏" } },
  { path: "/notifications", name: "notifications", component: NotificationsView, meta: { title: "通知中心" } },
  { path: "/profile", name: "profile", component: ProfileView, meta: { title: "个人资料" } },
  { path: "/interview", name: "interview", component: InterviewView, meta: { title: "AI 面试官" } },
  { path: "/author/:id", name: "author", component: AuthorView, meta: { title: "作者主页" } },
  { path: "/:pathMatch(.*)*", name: "not-found", component: NotFoundView, meta: { title: "页面不存在" } }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
});

router.beforeEach((to) => {
  const token = localStorage.getItem("blog_token");
  if (["editor", "editor-edit", "project-editor", "project-editor-edit", "my-articles", "my-projects", "profile", "collections", "notifications", "interview"].includes(to.name) && !token) {
    return { name: "auth" };
  }
  document.title = to.meta?.title && to.meta.title !== "PulseBlog" ? `${to.meta.title} | PulseBlog` : "PulseBlog";
  return true;
});

export default router;
