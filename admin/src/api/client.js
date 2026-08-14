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

export default client;