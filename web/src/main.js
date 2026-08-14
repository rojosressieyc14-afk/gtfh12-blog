import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/main.css";
import "highlight.js/styles/github-dark.css";

const app = createApp(App);
app.directive("highlight", {
  mounted(el) {
    import("highlight.js").then(({ default: hljs }) => {
      el.querySelectorAll("pre code").forEach((block) => hljs.highlightElement(block));
    });
  },
  updated(el) {
    import("highlight.js").then(({ default: hljs }) => {
      el.querySelectorAll("pre code").forEach((block) => hljs.highlightElement(block));
    });
  },
});
app.use(createPinia()).use(router).mount("#app");
