package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	articleService *service.ArticleService
	projectService *service.ProjectService
	adminService   *service.AdminService
}

func NewAdminHandler(articleService *service.ArticleService, projectService *service.ProjectService, adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{articleService: articleService, projectService: projectService, adminService: adminService}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	data, err := h.adminService.Dashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载仪表盘失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminHandler) PendingArticles(c *gin.Context) {
	items, err := h.articleService.ListPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载待审核文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AdminHandler) PendingProjects(c *gin.Context) {
	items, err := h.projectService.ListPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载待审核项目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AdminHandler) ArticleDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.adminService.GetArticleDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) ProjectDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.adminService.GetProjectDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "项目不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) ReviewArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload service.ReviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	item, err := h.articleService.Review(uint(id), authUser.ID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	title := "审核结果通知"
	content := "你的文章《" + item.Title + "》审核结果已更新。"
	if payload.Action == "approve" {
		content = "你的文章《" + item.Title + "》已审核通过并发布。"
	}
	if payload.Action == "reject" {
		content = "你的文章《" + item.Title + "》未通过审核。"
	}
	_ = h.adminService.CreateNotification(service.NotificationCreateInput{
		UserID:    item.AuthorID,
		Title:     title,
		Content:   content,
		Type:      model.NotificationTypeArticleReview,
		ActionURL: "/my-articles",
		Payload: map[string]any{
			"articleId": item.ID,
			"action":    payload.Action,
		},
	})
	_ = h.adminService.CreateOperationLog(authUser.ID, "review_article", "article", item.ID, content)

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) Users(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.adminService.ListUsers(service.UserFilter{
		Keyword:  c.Query("keyword"),
		Role:     c.Query("role"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载用户列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) Articles(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	var categoryID *uint
	var tagID *uint
	if raw := c.Query("categoryId"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			value := uint(parsed)
			categoryID = &value
		}
	}
	if raw := c.Query("tagId"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			value := uint(parsed)
			tagID = &value
		}
	}
	items, pagination, err := h.adminService.ListArticles(service.ArticleFilter{
		Status:     c.Query("status"),
		Keyword:    c.Query("keyword"),
		CategoryID: categoryID,
		TagID:      tagID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载文章列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) Projects(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.adminService.ListProjects(service.ProjectAdminFilter{
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载项目列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) Comments(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.adminService.ListComments(service.CommentFilter{
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载评论列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) Logs(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	var dateFrom *time.Time
	var dateTo *time.Time
	if raw := c.Query("dateFrom"); raw != "" {
		if parsed, err := time.Parse("2006-01-02", raw); err == nil {
			dateFrom = &parsed
		}
	}
	if raw := c.Query("dateTo"); raw != "" {
		if parsed, err := time.Parse("2006-01-02", raw); err == nil {
			next := parsed.Add(24 * time.Hour)
			dateTo = &next
		}
	}
	items, pagination, err := h.adminService.ListOperationLogs(service.LogFilter{
		Keyword:  c.Query("keyword"),
		Action:   c.Query("action"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载操作日志失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) Uploads(c *gin.Context) {
	items, err := h.adminService.ListUploads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载上传列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AdminHandler) DeleteUpload(c *gin.Context) {
	name := c.Param("name")
	if err := h.adminService.DeleteUpload(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "delete_upload", "upload", 0, "删除上传资源："+name)
	}
	c.JSON(http.StatusOK, gin.H{"message": "上传资源已删除"})
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item, err := h.adminService.UpdateUserRole(uint(id), payload.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "update_user_role", "user", item.ID, "修改用户角色为 "+payload.Role)
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item, err := h.adminService.UpdateUserStatus(uint(id), payload.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "update_user_status", "user", item.ID, "修改用户状态为 "+payload.Status)
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteArticle(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "delete_article", "article", uint(id), "删除文章")
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}

func (h *AdminHandler) UpdateArticleTaxonomy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload service.ArticleTaxonomyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item, err := h.adminService.UpdateArticleTaxonomy(uint(id), payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "update_article_taxonomy", "article", item.ID, "update article taxonomy")
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) BulkUpdateArticleTaxonomy(c *gin.Context) {
	var payload service.BulkArticleTaxonomyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updated, err := h.adminService.BulkUpdateArticleTaxonomy(payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_update_article_taxonomy", "article", 0, "bulk update article taxonomy")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) BulkPublishArticles(c *gin.Context) {
	var payload service.BulkArticleIDsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updated, err := h.adminService.BulkPublishArticles(payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_publish_articles", "article", 0, "bulk publish articles")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) BulkDeleteArticles(c *gin.Context) {
	var payload service.BulkArticleIDsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updated, err := h.adminService.BulkDeleteArticles(payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_delete_articles", "article", 0, "bulk delete articles")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) BulkRejectArticles(c *gin.Context) {
	var payload service.BulkArticleRejectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	updated, err := h.adminService.BulkRejectArticles(authUser.ID, payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_reject_articles", "article", 0, "bulk reject articles")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) BulkApproveArticles(c *gin.Context) {
	var payload service.BulkArticleIDsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	updated, err := h.adminService.BulkApproveArticles(authUser.ID, payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_approve_articles", "article", 0, "bulk approve articles")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "delete_project", "project", uint(id), "删除项目")
	}
	c.JSON(http.StatusOK, gin.H{"message": "项目已删除"})
}

func (h *AdminHandler) ForcePublishArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.adminService.ForcePublishArticle(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "publish_article", "article", item.ID, "直接发布文章")
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) ForcePublishProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.adminService.ForcePublishProject(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "publish_project", "project", item.ID, "直接发布项目")
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) ReviewProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload service.ProjectReviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	item, err := h.projectService.Review(uint(id), authUser.ID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	title := "审核结果通知"
	content := "你的项目《" + item.Title + "》审核结果已更新。"
	if payload.Action == "approve" {
		content = "你的项目《" + item.Title + "》已审核通过并发布。"
	}
	if payload.Action == "reject" {
		content = "你的项目《" + item.Title + "》未通过审核。"
	}
	_ = h.adminService.CreateNotification(service.NotificationCreateInput{
		UserID:    item.AuthorID,
		Title:     title,
		Content:   content,
		Type:      model.NotificationTypeProjectReview,
		ActionURL: "/my-projects",
		Payload: map[string]any{
			"projectId": item.ID,
			"action":    payload.Action,
		},
	})
	_ = h.adminService.CreateOperationLog(authUser.ID, "review_project", "project", item.ID, content)

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) BulkReviewProjects(c *gin.Context) {
	var payload service.BulkProjectReviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	updated, err := h.adminService.BulkReviewProjects(authUser.ID, payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	if authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "bulk_review_projects", "project", 0, "bulk review projects")
	}
	c.JSON(http.StatusOK, gin.H{"updated": updated})
}

func (h *AdminHandler) UpdateProjectMeta(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload service.ProjectMetaPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item, err := h.adminService.UpdateProjectMeta(uint(id), payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "update_project_meta", "project", item.ID, "更新项目精选状态和排序")
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *AdminHandler) DeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteComment(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "delete_comment", "comment", uint(id), "删除评论")
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}

func (h *AdminHandler) SensitiveWords(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.adminService.ListSensitiveWords(service.SensitiveWordFilter{
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载敏感词列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) CreateSensitiveWord(c *gin.Context) {
	var payload service.SensitiveWordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.adminService.CreateSensitiveWord(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "create_sensitive_word", "sensitive_word", item.ID, "新增敏感词："+item.Word)
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h *AdminHandler) DeleteSensitiveWord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.adminService.DeleteSensitiveWord(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "delete_sensitive_word", "sensitive_word", uint(id), "删除敏感词")
	}
	c.JSON(http.StatusOK, gin.H{"message": "敏感词已删除"})
}

func (h *AdminHandler) ModerationHits(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	items, pagination, err := h.adminService.ListModerationHits(service.ModerationHitFilter{
		Keyword:    c.Query("keyword"),
		Scene:      c.Query("scene"),
		AutoBanned: c.Query("autoBanned") == "true",
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载风控命中记录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *AdminHandler) GetModerationSettings(c *gin.Context) {
	data, err := h.adminService.GetModerationSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载风控设置失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminHandler) UpdateModerationSettings(c *gin.Context) {
	var payload service.ModerationSettingPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	data, err := h.adminService.UpdateModerationSettings(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		_ = h.adminService.CreateOperationLog(authUser.ID, "update_moderation_settings", "system_setting", 0, "更新自动封禁阈值")
	}
	c.JSON(http.StatusOK, data)
}
