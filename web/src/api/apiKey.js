import client from "./client";

export const listApiKeys = () => client.get("/user/api-keys");
export const createApiKey = (payload) => client.post("/user/api-keys", payload);
export const deleteApiKey = (id) => client.delete(`/user/api-keys/${id}`);
