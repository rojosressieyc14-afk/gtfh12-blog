import client from "./client";

export const getMetadata = () => client.get("/metadata");
