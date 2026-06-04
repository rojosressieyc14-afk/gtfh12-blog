import client from "./client";

export const startInterview = (payload) => client.post("/interview/start", payload);
export const submitAnswer = (id, payload) => client.post(`/interview/${id}/answer`, payload);
export const getSession = (id) => client.get(`/interview/${id}`);
export const endSession = (id) => client.post(`/interview/${id}/end`);
