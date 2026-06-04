package handler

import (
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AIReviewHandler struct {
	svc *service.AIReviewService
}

func NewAIReviewHandler(svc *service.AIReviewService) *AIReviewHandler {
	return &AIReviewHandler{svc: svc}
}

func (h *AIReviewHandler) GetContent(c *gin.Context) {
	targetType := c.Param("type")
	if targetType != model.TargetTypeArticle && targetType != model.TargetTypeProject {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不支持的内容类型，仅支持 article 或 project"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}

	content, err := h.svc.GetContentForReview(targetType, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "内容不存在"})
		return
	}

	prompt := service.FormatContentForPrompt(content, targetType)

	c.JSON(http.StatusOK, gin.H{
		"item":   content,
		"prompt": prompt,
	})
}

func (h *AIReviewHandler) SaveResult(c *gin.Context) {
	targetType := c.Param("type")
	if targetType != model.TargetTypeArticle && targetType != model.TargetTypeProject {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不支持的内容类型，仅支持 article 或 project"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}

	var payload struct {
		RiskLevel          string   `json:"riskLevel"`
		RiskLabels         []string `json:"riskLabels"`
		Summary            string   `json:"summary"`
		SuspiciousSegments []string `json:"suspiciousSegments"`
		Suggestion         string   `json:"suggestion"`
		ModelName          string   `json:"modelName"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	authUser := middleware.GetAuthUser(c)
	record, err := h.svc.SaveResult(authUser.ID, targetType, uint(id), service.AIReviewResultInput{
		RiskLevel:          payload.RiskLevel,
		RiskLabels:         payload.RiskLabels,
		Summary:            payload.Summary,
		SuspiciousSegments: payload.SuspiciousSegments,
		Suggestion:         payload.Suggestion,
		ModelName:          payload.ModelName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": record})
}

func (h *AIReviewHandler) GetResult(c *gin.Context) {
	targetType := c.Param("type")
	if targetType != model.TargetTypeArticle && targetType != model.TargetTypeProject {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不支持的内容类型"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}

	record, err := h.svc.GetRecord(targetType, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "暂无 AI 审核记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": record})
}
