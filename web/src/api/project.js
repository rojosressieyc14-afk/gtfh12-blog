import client from "./client";

export const listProjects = ({
  keyword = "",
  page = 1,
  pageSize = 9,
  featured = false,
  stack = "",
  authorId = "",
  sort = "featured"
} = {}) =>
  client.get("/projects", { params: { keyword, page, pageSize, featured, stack, authorId, sort } });

export const getProject = (id) => client.get(`/projects/${id}`);
export const listMyProjects = (page = 1, pageSize = 8) => client.get("/my/projects", { params: { page, pageSize } });
export const getMyProject = (id) => client.get(`/my/projects/${id}`);
export const createProject = (payload) => client.post("/projects", payload);
export const updateProject = (id, payload) => client.put(`/projects/${id}`, payload);
export const submitProject = (id) => client.post(`/projects/${id}/submit`);
export const deleteProject = (id) => client.delete(`/projects/${id}`);
