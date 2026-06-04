import { defineStore } from "pinia";
import { getMe, login, register } from "../api/auth";
import { updateProfile } from "../api/profile";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: localStorage.getItem("blog_token") || "",
    profile: null,
    loading: false
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    isAdmin: (state) => state.profile?.role === "admin",
    isBanned: (state) => state.profile?.status === "banned"
  },
  actions: {
    async loginAction(payload) {
      this.loading = true;
      try {
        const { data } = await login(payload);
        this.token = data.token;
        this.profile = data.user;
        localStorage.setItem("blog_token", data.token);
      } finally {
        this.loading = false;
      }
    },
    async registerAction(payload) {
      this.loading = true;
      try {
        const { data } = await register(payload);
        this.token = data.token;
        this.profile = data.user;
        localStorage.setItem("blog_token", data.token);
      } finally {
        this.loading = false;
      }
    },
    async fetchProfile() {
      if (!this.token) return;
      try {
        const { data } = await getMe();
        this.profile = data.user;
        if (data.user.status === "banned") {
          this.logout();
        }
      } catch (error) {
        this.logout();
      }
    },
    async updateProfileAction(payload) {
      const { data } = await updateProfile(payload);
      this.profile = data.user;
    },
    logout() {
      this.token = "";
      this.profile = null;
      localStorage.removeItem("blog_token");
    }
  }
});
