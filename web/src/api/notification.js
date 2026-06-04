import client from "./client";

export const listNotifications = (params = {}) => client.get("/notifications", { params });
export const markNotificationRead = (id) => client.post(`/notifications/${id}/read`);
export const markAllNotificationsRead = () => client.post("/notifications/read-all");
