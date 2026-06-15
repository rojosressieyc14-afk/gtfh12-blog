package handler

import (
	"net/http"

	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type UserCenterHandler struct {
	userCenterService *service.UserCenterService
}

func NewUserCenterHandler(userCenterService *service.UserCenterService) *UserCenterHandler {
	return &UserCenterHandler{userCenterService: userCenterService}
}

func (h *UserCenterHandler) GetStats(c *gin.Context) {
	stats, err := h.userCenterService.GetStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *UserCenterHandler) GetRecentActivity(c *gin.Context) {
	items, err := h.userCenterService.GetRecentActivity(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取最近动态失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
