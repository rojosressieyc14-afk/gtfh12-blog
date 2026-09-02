package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseID(c *gin.Context, param string) (int, bool) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的资源 ID"})
		return 0, false
	}
	return id, true
}
