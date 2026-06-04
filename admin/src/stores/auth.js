import { defineStore } from "pinia";
import { getMe, login } from "../api/auth";

export const useAdminStore = defineStore("admin-auth", {
  state: () => ({
    token: localStorage.getItem("admin_token") || "",
    profile: null,
    loading: false
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
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

        this.token = data.token;
        this.profile = data.user;
        localStorage.setItem("admin_token", data.token);
      } finally {
        this.loading = false;
      }
    },
    async fetchProfile() {
      if (!this.token) return null;

      try {
        const { data } = await getMe();
        if (data.user?.role !== "admin") {
          this.logout();
          return null;
        }

        this.profile = data.user;
        return data.user;
      } catch {
        this.logout();
        return null;
      }
    },
    logout() {
      this.token = "";
      this.profile = null;
      localStorage.removeItem("admin_token");
    }
  }
});
