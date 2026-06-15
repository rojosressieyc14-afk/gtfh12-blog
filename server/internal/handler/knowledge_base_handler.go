package handler

import (
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type KnowledgeBaseHandler struct {
	svc *service.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(svc *service.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{svc: svc}
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

func (h *KnowledgeBaseHandler) Get(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	kbID, _ := strconv.Atoi(c.Param("id"))
	docID, _ := strconv.Atoi(c.Param("docId"))
	if kbID <= 0 || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	kbID, _ := strconv.Atoi(c.Param("id"))
	docID, _ := strconv.Atoi(c.Param("docId"))
	if kbID <= 0 || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
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
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}

	var payload struct {
		Question string `json:"question"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	result, err := h.svc.Query(uint(id), authUser.ID, payload.Question, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *KnowledgeBaseHandler) GetPublicNote(c *gin.Context) {
	noteID, _ := strconv.Atoi(c.Param("id"))
	if noteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}

	doc, err := h.svc.GetPublicNote(uint(noteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "笔记不存在或未公开"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": doc})
}
