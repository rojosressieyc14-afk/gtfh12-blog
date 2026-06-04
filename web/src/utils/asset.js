export function toAssetUrl(path) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  // Dev: Vite :5173 → Go :8080; Prod: Nginx serves everything on same origin
  if (import.meta.env.DEV) {
    return `http://localhost:8080${path}`;
  }
  const base = import.meta.env.VITE_BASE_PATH || "";
  return `${base}${path}`;
}
