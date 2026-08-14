import client from "./client";

export const register = (payload) => client.post("/auth/register", payload);
export const login = (payload) => client.post("/auth/login", payload);
export const logout = () => client.post("/auth/logout");
export const getMe = () => client.get("/auth/me");
