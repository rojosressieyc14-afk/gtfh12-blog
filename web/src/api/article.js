import client from "./client";

export const listArticles = ({
  keyword = "",
  page = 1,
  pageSize = 9,
  categoryId = "",
  tag = "",
  authorId = "",
  sort = "latest"
} = {}) =>
  client.get("/articles", { params: { keyword, page, pageSize, categoryId, tag, authorId, sort } });
export const listTrendingArticles = () => client.get("/articles/trending");
export const getArticle = (id) => client.get(`/articles/${id}`);
export const createArticle = (payload) => client.post("/articles", payload);
export const updateArticle = (id, payload) => client.put(`/articles/${id}`, payload);
export const submitArticle = (id) => client.post(`/articles/${id}/submit`);
export const listMyArticles = (page = 1, pageSize = 8) => client.get("/my/articles", { params: { page, pageSize } });
export const listComments = (id) => client.get(`/articles/${id}/comments`);
export const createComment = (id, payload) => client.post(`/articles/${id}/comments`, payload);
export const toggleLike = (id) => client.post(`/articles/${id}/like`);
export const toggleFavorite = (id) => client.post(`/articles/${id}/favorite`);
export const deleteArticle = (id) => client.delete(`/articles/${id}`);
export const listMyLikes = () => client.get("/my/likes");
export const listMyFavorites = () => client.get("/my/favorites");
