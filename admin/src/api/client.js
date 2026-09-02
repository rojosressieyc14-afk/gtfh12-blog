import axios from "axios";

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api",
  timeout: 10000,
  withCredentials: true
});

function getCookie(name) {
  const match = document.cookie.match(new RegExp("(^|; )" + name + "=([^;]*)"));
  return match ? decodeURIComponent(match[2]) : "";
}

client.interceptors.request.use((config) => {
  const csrf = getCookie("blog_csrf");
  if (csrf) {
    config.headers["X-CSRF-Token"] = csrf;
  }
  return config;
});

function showToast(message, type = "error") {
  const existing = document.getElementById("pulse-admin-toast");
  if (existing) existing.remove();

  const toast = document.createElement("div");
  toast.id = "pulse-admin-toast";
  toast.textContent = message;
  Object.assign(toast.style, {
    position: "fixed",
    top: "20px",
    left: "50%",
    transform: "translateX(-50%)",
    padding: "12px 24px",
    borderRadius: "12px",
    fontSize: "14px",
    fontWeight: "500",
    zIndex: "9999",
    backdropFilter: "blur(12px)",
    transition: "opacity 0.3s ease",
    fontFamily: '"Outfit", "Noto Sans SC", sans-serif',
    ...(type === "error"
      ? { background: "rgba(255, 80, 80, 0.9)", color: "#fff" }
      : { background: "rgba(255, 209, 102, 0.9)", color: "#1d1a1a" })
  });
  document.body.appendChild(toast);
  setTimeout(() => {
    toast.style.opacity = "0";
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

client.interceptors.response.use(
  (res) => res,
  (err) => {
    const status = err.response?.status;
    const msg = err.response?.data?.message;

    if (status === 401) {
      if (window.location.pathname !== "/login") {
        showToast(msg || "登录已过期，请重新登录", "warning");
        setTimeout(() => {
          window.location.href = "/login";
        }, 800);
      }
    } else if (status === 403) {
      showToast(msg || "无权访问", "error");
    } else if (status >= 500) {
      showToast(msg || "服务器异常，请稍后重试", "error");
    }

    return Promise.reject(err);
  }
);

export default client;
