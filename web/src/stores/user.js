import { defineStore } from "pinia";
import { getMe, login, logout, register, sendVerifyCode, verifyRegister } from "../api/auth";
import { updateProfile } from "../api/profile";

export const useUserStore = defineStore("user", {
  state: () => ({
    profile: null,
    loading: false,
    restoring: true
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.profile),
    isAdmin: (state) => state.profile?.role === "admin",
    isBanned: (state) => state.profile?.status === "banned"
  },
  actions: {
    async loginAction(payload) {
      this.loading = true;
      try {
        const { data } = await login(payload);
        this.profile = data.user;
      } finally {
        this.loading = false;
      }
    },
    async registerAction(payload) {
      this.loading = true;
      try {
        const { data } = await verifyRegister(payload);
        this.profile = data.user;
      } finally {
        this.loading = false;
      }
    },
    async sendCodeAction(email) {
      await sendVerifyCode(email);
    },
    async fetchProfile() {
      this.restoring = true;
      try {
        const { data } = await getMe();
        this.profile = data.user;
        if (data.user.status === "banned") {
          this.profile = null;
        }
      } catch (error) {
        this.profile = null;
      } finally {
        this.restoring = false;
      }
    },
    async updateProfileAction(payload) {
      const { data } = await updateProfile(payload);
      this.profile = data.user;
    },
    async logout() {
      try {
        await logout();
      } catch {}
      this.profile = null;
    }
  }
});