import { defineStore } from "pinia";
import { getMe, login, logout } from "../api/auth";

export const useAdminStore = defineStore("admin-auth", {
  state: () => ({
    profile: null,
    loading: false,
    restoring: true
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.profile),
    isAdmin: (state) => state.profile?.role === "admin"
  },
  actions: {
    async loginAction(payload) {
      this.loading = true;
      try {
        const { data } = await login(payload);
        if (data.user?.role !== "admin") {
          throw new Error("当前账号不是管理员，无法进入后台。");
        }
        this.profile = data.user;
      } finally {
        this.loading = false;
      }
    },
    async fetchProfile() {
      this.restoring = true;
      try {
        const { data } = await getMe();
        if (data.user?.role !== "admin") {
          this.profile = null;
          return null;
        }
        this.profile = data.user;
        return data.user;
      } catch {
        this.profile = null;
        return null;
      } finally {
        this.restoring = false;
      }
    },
    async logout() {
      try {
        await logout();
      } catch {}
      this.profile = null;
    }
  }
});