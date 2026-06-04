import { createRouter, createWebHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";
import ModerationHitsView from "../views/ModerationHitsView.vue";
import SensitiveWordsView from "../views/SensitiveWordsView.vue";
import CommentsView from "../views/CommentsView.vue";
import LogsView from "../views/LogsView.vue";
import UploadsView from "../views/UploadsView.vue";
import SettingsView from "../views/SettingsView.vue";

const router = createRouter({
  history: createWebHistory("/PulseBlog/admin/"),
  routes: [
    { path: "/login", name: "login", component: LoginView, meta: { title: "管理后台登录" } },
    { path: "/", name: "dashboard", component: DashboardView, meta: { title: "后台总览" } },
    { path: "/moderation-hits", name: "moderation-hits", component: ModerationHitsView, meta: { title: "风控命中" } },
    { path: "/sensitive-words", name: "sensitive-words", component: SensitiveWordsView, meta: { title: "敏感词管理" } },
    { path: "/comments", name: "comments", component: CommentsView, meta: { title: "评论管理" } },
    { path: "/logs", name: "logs", component: LogsView, meta: { title: "操作日志" } },
    { path: "/uploads", name: "uploads", component: UploadsView, meta: { title: "资源库" } },
    { path: "/settings", name: "settings", component: SettingsView, meta: { title: "系统设置" } }
  ]
});

router.beforeEach((to) => {
  const token = localStorage.getItem("admin_token");
  if (to.name !== "login" && !token) {
    return { name: "login" };
  }
  if (to.name === "login" && token) {
    return { name: "dashboard" };
  }
  document.title = to.meta?.title ? `${to.meta.title} | PulseBlog` : "PulseBlog 管理后台";
  return true;
});

export default router;
