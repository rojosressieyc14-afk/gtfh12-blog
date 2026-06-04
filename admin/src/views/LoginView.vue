<template>
  <section class="admin-login admin-login--brand">
    <div class="login-panel login-panel--brand">
      <p class="login-tag">PulseBlog Admin</p>
      <h1>独立后台管理入口</h1>
      <p class="login-text">
        在这里统一处理文章审核、作品管理、评论维护、素材清理和运营日志。当前默认管理员账号为
        <code>admin</code>，密码为 <code>admin123</code>。
      </p>

      <form class="login-form" @submit.prevent="submit">
        <input v-model.trim="form.username" placeholder="管理员用户名" autocomplete="username" />
        <input
          v-model.trim="form.password"
          type="password"
          placeholder="管理员密码"
          autocomplete="current-password"
        />
        <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
        <button :disabled="store.loading">
          {{ store.loading ? "登录中..." : "进入后台" }}
        </button>
      </form>
    </div>
  </section>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useAdminStore } from "../stores/auth";

const router = useRouter();
const store = useAdminStore();
const errorMessage = ref("");
const form = reactive({
  username: "admin",
  password: "admin123"
});

async function submit() {
  errorMessage.value = "";

  try {
    await store.loginAction(form);
    router.push("/");
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || error.message || "登录失败，请检查账号或密码。";
  }
}
</script>

<style scoped>
.admin-login--brand {
  background:
    radial-gradient(circle at 18% 18%, rgba(255, 166, 77, 0.14), transparent 24%),
    radial-gradient(circle at 82% 24%, rgba(62, 255, 207, 0.1), transparent 24%);
}

.login-panel--brand {
  background:
    radial-gradient(circle at top right, rgba(255, 166, 77, 0.12), transparent 30%),
    rgba(255, 255, 255, 0.06);
}

@media (max-width: 768px) {
  .admin-login--brand {
    padding: 24px;
  }
  .login-panel--brand {
    padding: 28px;
  }
}

@media (max-width: 480px) {
  .admin-login--brand {
    padding: 12px;
  }
  .login-panel--brand {
    padding: 20px;
  }
  .login-panel--brand h1 {
    font-size: clamp(1.6rem, 6vw, 2rem);
  }
  .login-form input,
  .login-form button {
    padding: 12px 14px;
  }
}
</style>
