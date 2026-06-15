import client from "./client";

export const getUserStats = () => client.get("/user-center/stats");
export const getRecentActivity = () => client.get("/user-center/recent-activity");
