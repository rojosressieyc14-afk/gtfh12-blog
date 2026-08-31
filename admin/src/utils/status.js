const ARTICLE_STATUS_MAP = {
  draft: { label: "草稿", class: "pill--draft" },
  pending: { label: "待审核", class: "pill--pending" },
  published: { label: "已发布", class: "pill--published" },
  rejected: { label: "已拒绝", class: "pill--rejected" },
};

const PROJECT_STATUS_MAP = {
  draft: { label: "草稿", class: "pill--draft" },
  pending: { label: "待审核", class: "pill--pending" },
  published: { label: "已发布", class: "pill--published" },
  rejected: { label: "已拒绝", class: "pill--rejected" },
};

export function statusLabel(status) {
  return ARTICLE_STATUS_MAP[status]?.label || status || "未知";
}

export function statusClass(status) {
  return ARTICLE_STATUS_MAP[status]?.class || "pill--default";
}

export function projectStatusLabel(status) {
  return PROJECT_STATUS_MAP[status]?.label || status || "未知";
}

export function projectStatusClass(status) {
  return PROJECT_STATUS_MAP[status]?.class || "pill--default";
}
