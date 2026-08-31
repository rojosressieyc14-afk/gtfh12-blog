import client from "./client";

export const listKnowledgeBases = () => client.get("/knowledge-bases");
export const getKnowledgeBase = (id) => client.get(`/knowledge-bases/${id}`);
export const createKnowledgeBase = (payload) => client.post("/knowledge-bases", payload);
export const deleteKnowledgeBase = (id) => client.delete(`/knowledge-bases/${id}`);

export const listDocuments = (kbId) => client.get(`/knowledge-bases/${kbId}/documents`);
export const addDocument = (kbId, payload) => client.post(`/knowledge-bases/${kbId}/documents`, payload);
export const updateDocument = (kbId, docId, payload) => client.put(`/knowledge-bases/${kbId}/documents/${docId}`, payload);
export const deleteDocument = (kbId, docId) => client.delete(`/knowledge-bases/${kbId}/documents/${docId}`);

export const queryKnowledgeBase = (kbId, payload) => client.post(`/knowledge-bases/${kbId}/query`, payload);

export const searchDocuments = (kbId, keyword) => client.get(`/knowledge-bases/${kbId}/search`, { params: { q: keyword } });

export const getDocumentTree = (kbId) => client.get(`/knowledge-bases/${kbId}/tree`);
export const moveDocument = (kbId, docId, payload) => client.put(`/knowledge-bases/${kbId}/documents/${docId}/move`, payload);

export const getPublicNote = (id) => client.get(`/kb-notes/${id}`);

export const getBacklinks = (kbId, docId) => client.get(`/knowledge-bases/${kbId}/backlinks/${docId}`);
export const getGraphData = (kbId) => client.get(`/knowledge-bases/${kbId}/graph`);
export const listKbTags = (kbId) => client.get(`/knowledge-bases/${kbId}/tags`);
