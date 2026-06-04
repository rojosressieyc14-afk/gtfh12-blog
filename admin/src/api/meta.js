import client from "./client";

export const getMetadata = () => client.get("/metadata");
export const createCategory = (payload) => client.post("/admin/categories", payload);
export const updateCategory = (id, payload) => client.put(`/admin/categories/${id}`, payload);
export const deleteCategory = (id) => client.delete(`/admin/categories/${id}`);
export const createTag = (payload) => client.post("/admin/tags", payload);
export const updateTag = (id, payload) => client.put(`/admin/tags/${id}`, payload);
export const deleteTag = (id) => client.delete(`/admin/tags/${id}`);
