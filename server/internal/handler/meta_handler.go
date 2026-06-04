package handler

import (
	"errors"
	"net/http"
	"strconv"

	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type MetaHandler struct {
	metaService *service.MetaService
}

func NewMetaHandler(metaService *service.MetaService) *MetaHandler {
	return &MetaHandler{metaService: metaService}
}

func (h *MetaHandler) List(c *gin.Context) {
	data, err := h.metaService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载元数据失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *MetaHandler) CreateCategory(c *gin.Context) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	category, err := h.metaService.CreateCategory(payload.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": category})
}

func (h *MetaHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的分类 ID"})
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	category, err := h.metaService.UpdateCategory(uint(id), payload.Name)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrCategoryNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": category})
}

func (h *MetaHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的分类 ID"})
		return
	}

	if err := h.metaService.DeleteCategory(uint(id)); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrCategoryNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *MetaHandler) CreateTag(c *gin.Context) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tag, err := h.metaService.CreateTag(payload.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": tag})
}

func (h *MetaHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的标签 ID"})
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tag, err := h.metaService.UpdateTag(uint(id), payload.Name)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrTagNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": tag})
}

func (h *MetaHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的标签 ID"})
		return
	}

	if err := h.metaService.DeleteTag(uint(id)); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrTagNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
