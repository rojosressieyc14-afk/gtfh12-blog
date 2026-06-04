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

	loginLimiter := middleware.NewRateLimiter(10, time.Minute)
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

	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(articleService)
	projectHandler := handler.NewProjectHandler(projectService)
	adminHandler := handler.NewAdminHandler(articleService, projectService, adminService)
	metaHandler := handler.NewMetaHandler(metaService)
	commentHandler := handler.NewCommentHandler(commentService)
	uploadHandler := handler.NewUploadHandler(cfg)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	aiReviewHandler := handler.NewAIReviewHandler(aiReviewService)
	interviewHandler := handler.NewInterviewHandler(service.NewInterviewService(db, cfg))

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
		api.GET("/metadata", metaHandler.List)
		api.GET("/articles", articleHandler.ListPublished)
		api.GET("/articles/trending", articleHandler.Trending)
		api.GET("/projects", projectHandler.ListPublished)
		api.GET("/projects/:id", projectHandler.Detail)
		api.GET("/articles/:id", articleHandler.Detail)
		api.GET("/articles/:id/comments", commentHandler.List)

		api.GET("/auth/me", middleware.RequireAuth(cfg), authHandler.Me)
		api.PUT("/auth/me", middleware.RequireAuth(cfg), authHandler.UpdateProfile)
		api.GET("/notifications", middleware.RequireAuth(cfg), notificationHandler.List)
		api.POST("/notifications/:id/read", middleware.RequireAuth(cfg), notificationHandler.MarkRead)
		api.POST("/notifications/read-all", middleware.RequireAuth(cfg), notificationHandler.MarkAllRead)
		api.GET("/my/articles", middleware.RequireAuth(cfg), articleHandler.Mine)
		api.GET("/my/projects", middleware.RequireAuth(cfg), projectHandler.Mine)
		api.GET("/my/projects/:id", middleware.RequireAuth(cfg), projectHandler.MineDetail)
		api.GET("/my/likes", middleware.RequireAuth(cfg), articleHandler.Liked)
		api.GET("/my/favorites", middleware.RequireAuth(cfg), articleHandler.Favorited)
		api.POST("/articles", middleware.RequireAuth(cfg), articleHandler.Create)
		api.POST("/projects", middleware.RequireAuth(cfg), projectHandler.Create)
		api.PUT("/articles/:id", middleware.RequireAuth(cfg), articleHandler.Update)
		api.PUT("/projects/:id", middleware.RequireAuth(cfg), projectHandler.Update)
		api.POST("/articles/:id/submit", middleware.RequireAuth(cfg), articleHandler.Submit)
		api.POST("/projects/:id/submit", middleware.RequireAuth(cfg), projectHandler.Submit)
		api.POST("/articles/:id/like", middleware.RequireAuth(cfg), articleHandler.ToggleLike)
		api.POST("/articles/:id/favorite", middleware.RequireAuth(cfg), articleHandler.ToggleFavorite)
		api.POST("/articles/:id/comments", middleware.RequireAuth(cfg), commentHandler.Create)
		api.POST("/upload", middleware.RequireAuth(cfg), uploadHandler.UploadImage)
		api.POST("/interview/start", middleware.RequireAuth(cfg), interviewHandler.Start)
		api.POST("/interview/:id/answer", middleware.RequireAuth(cfg), interviewHandler.Answer)
		api.GET("/interview/:id", middleware.RequireAuth(cfg), interviewHandler.Get)
		api.POST("/interview/:id/end", middleware.RequireAuth(cfg), interviewHandler.End)

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(cfg), middleware.RequireAdmin())
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
