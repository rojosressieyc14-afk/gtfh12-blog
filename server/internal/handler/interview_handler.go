package handler

import (
	"net/http"
	"strconv"
	"strings"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type InterviewHandler struct {
	svc *service.InterviewService
}

func NewInterviewHandler(svc *service.InterviewService) *InterviewHandler {
	return &InterviewHandler{svc: svc}
}

func (h *InterviewHandler) Start(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	var payload struct {
		Position       string `json:"position"`
		ResumeText     string `json:"resumeText"`
		TotalQuestions int    `json:"totalQuestions"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	payload.Position = strings.TrimSpace(payload.Position)
	if payload.Position == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写应聘职位"})
		return
	}
	if payload.TotalQuestions < 1 {
		payload.TotalQuestions = 5
	}

	payload.ResumeText = strings.TrimSpace(payload.ResumeText)

	session, err := h.svc.StartSession(authUser.ID, payload.Position, payload.ResumeText, payload.TotalQuestions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": session})
}

func (h *InterviewHandler) Answer(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil || sessionID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的会话 ID"})
		return
	}

	var payload struct {
		Answer string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数无效"})
		return
	}

	payload.Answer = strings.TrimSpace(payload.Answer)
	if payload.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "回答不能为空"})
		return
	}

	session, completed, err := h.svc.SubmitAnswer(uint(sessionID), authUser.ID, payload.Answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": session, "completed": completed})
}

func (h *InterviewHandler) Get(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil || sessionID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的会话 ID"})
		return
	}

	session, err := h.svc.GetSession(uint(sessionID), authUser.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": session})
}

func (h *InterviewHandler) End(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil || sessionID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的会话 ID"})
		return
	}

	session, err := h.svc.EndSession(uint(sessionID), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": session})
}
