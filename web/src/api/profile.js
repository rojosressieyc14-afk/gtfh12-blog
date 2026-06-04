import client from "./client";

export const updateProfile = (payload) => client.put("/auth/me", payload);
export const getAuthorProfile = (id) => client.get(`/authors/${id}`);
export const getRecommendedAuthors = () => client.get("/authors/recommended");
