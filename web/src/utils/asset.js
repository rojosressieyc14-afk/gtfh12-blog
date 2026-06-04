export function toAssetUrl(path) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return `http://localhost:8080${path}`;
}
