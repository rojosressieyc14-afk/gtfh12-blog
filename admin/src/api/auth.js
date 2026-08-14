import client from "./client";

export const login = (payload) => client.post("/auth/login", payload);
export const logout = () => client.post("/auth/logout");
export const getMe = () => client.get("/auth/me");