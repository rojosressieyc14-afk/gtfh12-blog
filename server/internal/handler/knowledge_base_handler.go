package handler

import (
	"net/http"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type KnowledgeBaseHandler struct {
	svc *service.KnowledgeBaseService
	llm service.LLMProvider
}

func NewKnowledgeBaseHandler(svc *service.KnowledgeBaseService, llm service.LLMProvider) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{svc: svc, llm: llm}
}

func (h *KnowledgeBaseHandler) Create(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	kb, err := h.svc.Create(authUser.ID, payload.Name, payload.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": kb})
}

func (h *KnowledgeBaseHandler) List(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbs, err := h.svc.List(authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载知识库列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": kbs})
}

func (h *KnowledgeBaseHandler) ListPublicKBs(c *gin.Context) {
	items, err := h.svc.ListPublicKBs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载知识库列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *KnowledgeBaseHandler) ListPublicDocuments(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	docs, err := h.svc.ListPublicDocuments(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载文档列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": docs})
}

func (h *KnowledgeBaseHandler) Get(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	kb, err := h.svc.GetByID(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": kb})
}

func (h *KnowledgeBaseHandler) Delete(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(uint(id), authUser.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "知识库已删除"})
}

func (h *KnowledgeBaseHandler) AddDocument(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var payload struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		IsPublic   bool     `json:"isPublic"`
		IsMarkdown bool     `json:"isMarkdown"`
		CategoryID *uint    `json:"categoryId"`
		TagNames   []string `json:"tagNames"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	doc, err := h.svc.AddDocument(uint(id), authUser.ID, service.AddDocumentOpts{
		Title:      payload.Title,
		Content:    payload.Content,
		IsPublic:   payload.IsPublic,
		IsMarkdown: payload.IsMarkdown,
		CategoryID: payload.CategoryID,
		TagNames:   payload.TagNames,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": doc})
}

func (h *KnowledgeBaseHandler) UpdateDocument(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}
	docID, ok := parseID(c, "docId")
	if !ok {
		return
	}

	var payload struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		IsPublic   bool     `json:"isPublic"`
		IsMarkdown bool     `json:"isMarkdown"`
		CategoryID *uint    `json:"categoryId"`
		TagNames   []string `json:"tagNames"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	doc, err := h.svc.UpdateDocument(uint(docID), uint(kbID), authUser.ID, service.AddDocumentOpts{
		Title:      payload.Title,
		Content:    payload.Content,
		IsPublic:   payload.IsPublic,
		IsMarkdown: payload.IsMarkdown,
		CategoryID: payload.CategoryID,
		TagNames:   payload.TagNames,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": doc})
}

func (h *KnowledgeBaseHandler) ListDocuments(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	docs, err := h.svc.ListDocuments(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": docs})
}

func (h *KnowledgeBaseHandler) DeleteDocument(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}
	docID, ok := parseID(c, "docId")
	if !ok {
		return
	}

	if err := h.svc.DeleteDocument(uint(kbID), uint(docID), authUser.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文档已删除"})
}

func (h *KnowledgeBaseHandler) Query(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var payload struct {
		Question string `json:"question"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	result, err := h.svc.Query(uint(id), authUser.ID, payload.Question, h.llm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *KnowledgeBaseHandler) GetDocumentTree(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	tree, err := h.svc.GetDocumentTree(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

func (h *KnowledgeBaseHandler) MoveDocument(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}
	docID, ok := parseID(c, "docId")
	if !ok {
		return
	}
	var opts service.MoveDocumentOpts
	if err := safeBindJSON(c, &opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}
	if err := h.svc.MoveDocument(uint(kbID), uint(docID), authUser.ID, opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文档已移动"})
}

func (h *KnowledgeBaseHandler) GetPublicNote(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	doc, err := h.svc.GetPublicNote(uint(noteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "笔记不存在或未公开"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": doc})
}

func (h *KnowledgeBaseHandler) Search(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "搜索关键词不能为空"})
		return
	}

	results, err := h.svc.SearchDocuments(uint(id), authUser.ID, keyword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": results})
}

func (h *KnowledgeBaseHandler) Backlinks(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}
	docID, ok := parseID(c, "docId")
	if !ok {
		return
	}

	items, err := h.svc.GetBacklinks(uint(kbID), uint(docID), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *KnowledgeBaseHandler) Graph(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}

	data, err := h.svc.GetGraphData(uint(kbID), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"graph": data})
}

func (h *KnowledgeBaseHandler) ListTags(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	kbID, ok := parseID(c, "id")
	if !ok {
		return
	}

	tags, err := h.svc.ListAllTags(uint(kbID), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": tags})
}
