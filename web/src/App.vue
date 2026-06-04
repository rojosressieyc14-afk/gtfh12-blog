<template>
  <div class="shell">
    <SiteHeader />
    <main class="page">
      <router-view />
    </main>
    <SiteFooter />
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useUserStore } from "./stores/user";
import SiteHeader from "./components/SiteHeader.vue";
import SiteFooter from "./components/SiteFooter.vue";

const userStore = useUserStore();

onMounted(() => {
  userStore.fetchProfile();
});
</script>
