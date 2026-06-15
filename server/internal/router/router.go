package router

import (
	"net/url"
	"strings"
	"time"

	"blog/server/internal/config"
	"blog/server/internal/handler"
	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SecurityHeaders())

	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	registerLimiter := middleware.NewRateLimiter(5, time.Minute)
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	allowOrigins := resolveAllowedOrigins(cfg)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authService := service.NewAuthService(db, cfg)
	articleService := service.NewArticleService(db)
	projectService := service.NewProjectService(db)
	metaService := service.NewMetaService(db)
	commentService := service.NewCommentService(db)
	adminService := service.NewAdminService(db, cfg)
	notificationService := service.NewNotificationService(db)
	aiReviewService := service.NewAIReviewService(db)
	userCenterService := service.NewUserCenterService(db)

	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(articleService)
	projectHandler := handler.NewProjectHandler(projectService)
	adminHandler := handler.NewAdminHandler(articleService, projectService, adminService)
	metaHandler := handler.NewMetaHandler(metaService)
	commentHandler := handler.NewCommentHandler(commentService)
	uploadHandler := handler.NewUploadHandler(cfg)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	aiReviewHandler := handler.NewAIReviewHandler(aiReviewService)
	userCenterHandler := handler.NewUserCenterHandler(userCenterService)

	var apiKeyService *service.ApiKeyService
	var apiKeyHandler *handler.ApiKeyHandler
	if encService, err := service.NewEncryptionService(cfg.APIEncryptionKey); err == nil {
		apiKeyService = service.NewApiKeyService(db, encService)
		apiKeyHandler = handler.NewApiKeyHandler(apiKeyService)
	}

	interviewHandler := handler.NewInterviewHandler(service.NewInterviewService(db, cfg, apiKeyService))

	embeddingProvider := service.NewDeepSeekEmbedding(cfg.DeepSeekKey, cfg.DeepSeekURL)
	kbService := service.NewKnowledgeBaseService(db, cfg.QdrantAddr, cfg.QdrantAPIKey, embeddingProvider)
	kbHandler := handler.NewKnowledgeBaseHandler(kbService)

	r.Static("/uploads", cfg.UploadDir)

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/auth/register", middleware.RateLimit(registerLimiter), authHandler.Register)
		api.POST("/auth/login", middleware.RateLimit(loginLimiter), authHandler.Login)
		api.GET("/authors/recommended", authHandler.RecommendedAuthors)
		api.GET("/authors/:id", authHandler.AuthorProfile)
		api.GET("/kb-notes/:id", kbHandler.GetPublicNote)
		api.GET("/metadata", metaHandler.List)
		api.GET("/articles", articleHandler.ListPublished)
		api.GET("/articles/trending", articleHandler.Trending)
		api.GET("/projects", projectHandler.ListPublished)
		api.GET("/projects/:id", projectHandler.Detail)
		api.GET("/articles/:id", articleHandler.Detail)
		api.GET("/articles/:id/comments", commentHandler.List)

		api.GET("/auth/me", middleware.RequireAuth(), authHandler.Me)
		api.PUT("/auth/me", middleware.RequireAuth(), authHandler.UpdateProfile)
		api.GET("/notifications", middleware.RequireAuth(), notificationHandler.List)
		api.POST("/notifications/:id/read", middleware.RequireAuth(), notificationHandler.MarkRead)
		api.POST("/notifications/read-all", middleware.RequireAuth(), notificationHandler.MarkAllRead)
		api.GET("/my/articles", middleware.RequireAuth(), articleHandler.Mine)
		api.GET("/my/projects", middleware.RequireAuth(), projectHandler.Mine)
		api.GET("/my/projects/:id", middleware.RequireAuth(), projectHandler.MineDetail)
		api.GET("/my/likes", middleware.RequireAuth(), articleHandler.Liked)
		api.GET("/my/favorites", middleware.RequireAuth(), articleHandler.Favorited)
		api.POST("/articles", middleware.RequireAuth(), articleHandler.Create)
		api.POST("/projects", middleware.RequireAuth(), projectHandler.Create)
		api.PUT("/articles/:id", middleware.RequireAuth(), articleHandler.Update)
		api.PUT("/projects/:id", middleware.RequireAuth(), projectHandler.Update)
		api.DELETE("/articles/:id", middleware.RequireAuth(), articleHandler.Delete)
		api.DELETE("/projects/:id", middleware.RequireAuth(), projectHandler.Delete)
		api.POST("/articles/:id/submit", middleware.RequireAuth(), articleHandler.Submit)
		api.POST("/projects/:id/submit", middleware.RequireAuth(), projectHandler.Submit)
		api.POST("/articles/:id/like", middleware.RequireAuth(), articleHandler.ToggleLike)
		api.POST("/articles/:id/favorite", middleware.RequireAuth(), articleHandler.ToggleFavorite)
		api.POST("/articles/:id/comments", middleware.RequireAuth(), commentHandler.Create)
		api.POST("/upload", middleware.RequireAuth(), uploadHandler.UploadImage)
		api.POST("/interview/start", middleware.RequireAuth(), interviewHandler.Start)
		api.POST("/interview/:id/answer", middleware.RequireAuth(), interviewHandler.Answer)
		api.GET("/interview/:id", middleware.RequireAuth(), interviewHandler.Get)
		api.POST("/interview/:id/end", middleware.RequireAuth(), interviewHandler.End)

		api.POST("/knowledge-bases", middleware.RequireAuth(), kbHandler.Create)
		api.GET("/knowledge-bases", middleware.RequireAuth(), kbHandler.List)
		api.GET("/knowledge-bases/:id", middleware.RequireAuth(), kbHandler.Get)
		api.DELETE("/knowledge-bases/:id", middleware.RequireAuth(), kbHandler.Delete)
		api.POST("/knowledge-bases/:id/documents", middleware.RequireAuth(), kbHandler.AddDocument)
		api.PUT("/knowledge-bases/:id/documents/:docId", middleware.RequireAuth(), kbHandler.UpdateDocument)
		api.GET("/knowledge-bases/:id/documents", middleware.RequireAuth(), kbHandler.ListDocuments)
		api.DELETE("/knowledge-bases/:id/documents/:docId", middleware.RequireAuth(), kbHandler.DeleteDocument)
		api.POST("/knowledge-bases/:id/query", middleware.RequireAuth(), kbHandler.Query)

		if apiKeyHandler != nil {
			api.GET("/user/api-keys", middleware.RequireAuth(), apiKeyHandler.List)
			api.POST("/user/api-keys", middleware.RequireAuth(), apiKeyHandler.Create)
			api.DELETE("/user/api-keys/:id", middleware.RequireAuth(), apiKeyHandler.Delete)
		}

		api.GET("/user-center/stats", middleware.RequireAuth(), userCenterHandler.GetStats)
		api.GET("/user-center/recent-activity", middleware.RequireAuth(), userCenterHandler.GetRecentActivity)

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(), middleware.RequireAdmin())
		{
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/reviews", adminHandler.PendingArticles)
			admin.GET("/project-reviews", adminHandler.PendingProjects)
			admin.GET("/users", adminHandler.Users)
			admin.GET("/articles", adminHandler.Articles)
			admin.GET("/articles/:id", adminHandler.ArticleDetail)
			admin.PUT("/articles/:id/taxonomy", adminHandler.UpdateArticleTaxonomy)
			admin.PUT("/articles/taxonomy/bulk", adminHandler.BulkUpdateArticleTaxonomy)
			admin.POST("/articles/publish/bulk", adminHandler.BulkPublishArticles)
			admin.POST("/articles/delete/bulk", adminHandler.BulkDeleteArticles)
			admin.POST("/articles/reject/bulk", adminHandler.BulkRejectArticles)
			admin.POST("/articles/approve/bulk", adminHandler.BulkApproveArticles)
			admin.GET("/projects", adminHandler.Projects)
			admin.GET("/projects/:id", adminHandler.ProjectDetail)
			admin.POST("/projects/review/bulk", adminHandler.BulkReviewProjects)
			admin.GET("/comments", adminHandler.Comments)
			admin.GET("/logs", adminHandler.Logs)
			admin.GET("/moderation-hits", adminHandler.ModerationHits)
			admin.GET("/moderation-settings", adminHandler.GetModerationSettings)
			admin.GET("/sensitive-words", adminHandler.SensitiveWords)
			admin.GET("/uploads", adminHandler.Uploads)
			admin.DELETE("/uploads/:name", adminHandler.DeleteUpload)
			admin.POST("/sensitive-words", adminHandler.CreateSensitiveWord)
			admin.PUT("/moderation-settings", adminHandler.UpdateModerationSettings)
			admin.POST("/articles/:id/review", adminHandler.ReviewArticle)
			admin.POST("/articles/:id/publish", adminHandler.ForcePublishArticle)
			admin.POST("/projects/:id/review", adminHandler.ReviewProject)
			admin.POST("/projects/:id/publish", adminHandler.ForcePublishProject)
			admin.PUT("/projects/:id/meta", adminHandler.UpdateProjectMeta)
			admin.DELETE("/articles/:id", adminHandler.DeleteArticle)
			admin.DELETE("/projects/:id", adminHandler.DeleteProject)
			admin.DELETE("/comments/:id", adminHandler.DeleteComment)
			admin.DELETE("/sensitive-words/:id", adminHandler.DeleteSensitiveWord)
			admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
			admin.POST("/categories", metaHandler.CreateCategory)
			admin.PUT("/categories/:id", metaHandler.UpdateCategory)
			admin.DELETE("/categories/:id", metaHandler.DeleteCategory)
			admin.POST("/tags", metaHandler.CreateTag)
			admin.PUT("/tags/:id", metaHandler.UpdateTag)
			admin.DELETE("/tags/:id", metaHandler.DeleteTag)
			admin.GET("/projects/:id/review-detail", adminHandler.ProjectDetail)
			admin.GET("/ai-review/:type/:id/content", aiReviewHandler.GetContent)
			admin.GET("/ai-review/:type/:id", aiReviewHandler.GetResult)
			admin.POST("/ai-review/:type/:id", aiReviewHandler.SaveResult)
		}
	}

	return r
}

func resolveAllowedOrigins(cfg config.Config) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	add := func(origin string) {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			return
		}
		if _, ok := seen[origin]; ok {
			return
		}
		seen[origin] = struct{}{}
		result = append(result, origin)
	}

	for _, origin := range []string{cfg.WebOrigin, cfg.AdminOrigin} {
		add(origin)
		add(loopbackVariant(origin))
	}

	if cfg.CORSOrigins != "" {
		for _, origin := range strings.Split(cfg.CORSOrigins, ",") {
			add(origin)
		}
	}

	return result
}

func loopbackVariant(origin string) string {
	parsed, err := url.Parse(strings.TrimSpace(origin))
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	port := parsed.Port()
	switch host {
	case "localhost":
		parsed.Host = "127.0.0.1"
	case "127.0.0.1":
		parsed.Host = "localhost"
	default:
		return ""
	}

	if port != "" {
		parsed.Host += ":" + port
	}

	return parsed.String()
}
