package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/server/internal/config"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	cfg config.Config
}

func NewUploadHandler(cfg config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": "文件大小不能超过 10MB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少上传文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不支持的图片格式"})
		return
	}

	if file.Size > 10<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": "文件大小不能超过 10MB"})
		return
	}

	if err := os.MkdirAll(h.cfg.UploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "准备上传目录失败"})
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(h.cfg.UploadDir, filename)

	cleanDst, err := filepath.Abs(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
		return
	}
	uploadAbs, err := filepath.Abs(h.cfg.UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
		return
	}
	if !strings.HasPrefix(cleanDst, uploadAbs+string(filepath.Separator)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "非法的文件路径"})
		return
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "/uploads/" + filename,
	})
}
