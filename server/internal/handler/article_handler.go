package handler

import (
	"errors"
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) ListPublished(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	var categoryID *uint
	var authorID uint

	if raw := c.Query("categoryId"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			value := uint(parsed)
			categoryID = &value
		}
	}
	if raw := c.Query("authorId"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			authorID = uint(parsed)
		}
	}

	articles, pagination, err := h.articleService.ListPublished(service.PublishedArticleFilter{
		Keyword:    c.Query("keyword"),
		CategoryID: categoryID,
		Tag:        c.Query("tag"),
		AuthorID:   authorID,
		Sort:       defaultQueryValue(c.Query("sort"), "latest"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": articles, "pagination": pagination})
}

func (h *ArticleHandler) Trending(c *gin.Context) {
	items, err := h.articleService.ListTrending(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载推荐文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ArticleHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var viewerID uint
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		viewerID = authUser.ID
	}

	article, err := h.articleService.GetByID(uint(id), viewerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "文章不存在"})
		return
	}

	authUser := middleware.GetAuthUser(c)
	if article.Status != model.ArticlePublished && (authUser == nil || (authUser.Role != model.RoleAdmin && authUser.ID != article.AuthorID)) {
		c.JSON(http.StatusForbidden, gin.H{"message": "文章暂未公开"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": article})
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var payload service.ArticlePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	article, err := h.articleService.Create(authUser.ID, authUser.Role, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if service.IsModerationError(err) || errors.Is(err, service.ErrUserBanned) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": article})
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var payload service.ArticlePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	article, err := h.articleService.Update(uint(id), authUser.ID, authUser.Role, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrArticleNoPermission) {
			status = http.StatusForbidden
		} else if service.IsModerationError(err) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": article})
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	if err := h.articleService.Delete(uint(id), authUser.ID, authUser.Role); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrArticleNoPermission) {
			status = http.StatusForbidden
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}

func (h *ArticleHandler) Submit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	article, err := h.articleService.Submit(uint(id), authUser.ID, authUser.Role)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrArticleNoPermission) {
			status = http.StatusForbidden
		} else if service.IsModerationError(err) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": article})
}

func (h *ArticleHandler) Mine(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))

	articles, pagination, err := h.articleService.ListMine(authUser.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载我的文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": articles, "pagination": pagination})
}

func (h *ArticleHandler) Liked(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	items, err := h.articleService.ListLiked(authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载点赞列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ArticleHandler) Favorited(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	items, err := h.articleService.ListFavorited(authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载收藏列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ArticleHandler) ToggleLike(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	item, err := h.articleService.ToggleLike(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *ArticleHandler) ToggleFavorite(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	item, err := h.articleService.ToggleFavorite(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func defaultQueryValue(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
