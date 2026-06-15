package handler

import (
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ApiKeyHandler struct {
	svc *service.ApiKeyService
}

func NewApiKeyHandler(svc *service.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{svc: svc}
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	keys, err := h.svc.List(authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载 API Key 列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": keys})
}

func (h *ApiKeyHandler) Create(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	var payload struct {
		Provider string `json:"provider"`
		Key      string `json:"key"`
		BaseURL  string `json:"baseURL"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	item, err := h.svc.Create(authUser.ID, payload.Provider, payload.Key, payload.BaseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h *ApiKeyHandler) Delete(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "API Key 已删除"})
}
