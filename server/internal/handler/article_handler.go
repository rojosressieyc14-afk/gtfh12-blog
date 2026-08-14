package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func safeBindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if strings.Contains(err.Error(), "invalid character") || strings.Contains(err.Error(), "unexpected end") {
			return errors.New("请求数据格式异常，请检查提交内容中是否有特殊字符或尝试刷新后重试")
		}
		return err
	}
	return nil
}

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

func (h *ArticleHandler) Feed(c *gin.Context) {
	articles, err := h.articleService.ListFeed(20)
	if err != nil {
		c.String(500, "Failed to generate feed")
		return
	}
	c.Header("Content-Type", "application/atom+xml; charset=utf-8")
	c.String(200, `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>PulseBlog</title>
  <link href="%s/feed" rel="self"/>
  <link href="%s"/>
  <updated>%s</updated>
  <id>urn:uuid:pulseblog</id>
`, h.siteURL(c), h.siteURL(c), time.Now().Format(time.RFC3339))
	for _, a := range articles {
		updated := a.UpdatedAt.Format(time.RFC3339)
		published := ""
		if a.PublishedAt != nil {
			published = a.PublishedAt.Format(time.RFC3339)
		}
		url := h.siteURL(c) + "/article/" + strconv.Itoa(int(a.ID))
		c.String(200, `  <entry>
    <title>%s</title>
    <link href="%s"/>
    <id>%s</id>
    <updated>%s</updated>
    <published>%s</published>
    <summary>%s</summary>
  </entry>
`, xmlEscape(a.Title), url, url, updated, published, xmlEscape(a.Summary))
	}
	c.String(200, `</feed>`)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func (h *ArticleHandler) siteURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
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
	if article.IsPrivate && (authUser == nil || (authUser.Role != model.RoleAdmin && authUser.ID != article.AuthorID)) {
		c.JSON(http.StatusForbidden, gin.H{"message": "文章暂未公开"})
		return
	}
	if article.Status != model.ArticlePublished && (authUser == nil || (authUser.Role != model.RoleAdmin && authUser.ID != article.AuthorID)) {
		c.JSON(http.StatusForbidden, gin.H{"message": "文章暂未公开"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": article})
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var payload service.ArticlePayload
	if err := safeBindJSON(c, &payload); err != nil {
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
	if err := safeBindJSON(c, &payload); err != nil {
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

func (h *ArticleHandler) Stats(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)
	items, err := h.articleService.GetArticleStats(uint(id), authUser.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func defaultQueryValue(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
