package handler

import (
	"fmt"
	"net/http"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type authPayload struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var payload authPayload
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, token, err := h.authService.Register(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload authPayload
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, token, err := h.authService.Login(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	user, err := h.authService.Me(authUser.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var payload struct {
		Avatar      string   `json:"avatar"`
		Headline    string   `json:"headline"`
		CurrentRole string   `json:"currentRole"`
		YearsLabel  string   `json:"yearsLabel"`
		Motto       string   `json:"motto"`
		Location    string   `json:"location"`
		Email       string   `json:"email"`
		ResumeURL   string   `json:"resumeUrl"`
		WebsiteURL  string   `json:"websiteUrl"`
		GithubURL   string   `json:"githubUrl"`
		GiteeURL    string   `json:"giteeUrl"`
		JuejinURL   string   `json:"juejinUrl"`
		CSDNURL     string   `json:"csdnUrl"`
		Skills      []string `json:"skills"`
		FocusAreas  []string `json:"focusAreas"`
		Bio         string   `json:"bio"`
	}
	if err := safeBindJSON(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	user, err := h.authService.UpdateProfile(authUser.ID, service.UpdateProfileInput{
		Avatar:      payload.Avatar,
		Headline:    payload.Headline,
		CurrentRole: payload.CurrentRole,
		YearsLabel:  payload.YearsLabel,
		Motto:       payload.Motto,
		Location:    payload.Location,
		Email:       payload.Email,
		ResumeURL:   payload.ResumeURL,
		WebsiteURL:  payload.WebsiteURL,
		GithubURL:   payload.GithubURL,
		GiteeURL:    payload.GiteeURL,
		JuejinURL:   payload.JuejinURL,
		CSDNURL:     payload.CSDNURL,
		Skills:      payload.Skills,
		FocusAreas:  payload.FocusAreas,
		Bio:         payload.Bio,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) AuthorProfile(c *gin.Context) {
	var userID uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户 ID"})
		return
	}

	user, articles, projects, stats, err := h.authService.AuthorProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "作者不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "articles": articles, "projects": projects, "stats": stats})
}

func (h *AuthHandler) RecommendedAuthors(c *gin.Context) {
	users, err := h.authService.GetRecommendedAuthors(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载推荐作者失败"})
		return
	}
	if users == nil {
		c.JSON(http.StatusOK, gin.H{"items": []struct{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": users})
}
