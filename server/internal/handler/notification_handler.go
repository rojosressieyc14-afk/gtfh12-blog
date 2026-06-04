package handler

import (
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.notificationService.List(authUser.ID, service.NotificationListFilter{
		UnreadOnly: c.Query("unreadOnly") == "true",
		Type:       c.Query("type"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载通知失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.notificationService.MarkRead(authUser.ID, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	if err := h.notificationService.MarkAllRead(authUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "批量标记已读失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已全部标记为已读"})
}
