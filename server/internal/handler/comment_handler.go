package handler

import (
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) List(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	items, err := h.commentService.List(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *CommentHandler) Create(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload service.CommentPayload
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	item, err := h.commentService.Create(uint(id), authUser.ID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}
