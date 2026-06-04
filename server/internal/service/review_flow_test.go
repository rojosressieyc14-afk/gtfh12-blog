package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"blog/server/internal/config"
	"blog/server/internal/database"
	"blog/server/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testUserSeq uint64

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=private", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func createTestUser(t *testing.T, db *gorm.DB, role string) model.User {
	t.Helper()
	seq := atomic.AddUint64(&testUserSeq, 1)
	user := model.User{
		Username: fmt.Sprintf("user_%s_%s_%d", role, strings.ReplaceAll(t.Name(), "/", "_"), seq),
		Password: "hashed",
		Role:     role,
		Status:   model.UserActive,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func TestArticleReviewApproveFlow(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	articleService := NewArticleService(db)
	article, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Go review flow",
		Summary: "summary",
		Content: "content",
		Tags:    []string{"go", "review"},
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	if _, err := articleService.Submit(article.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit article: %v", err)
	}

	reviewed, err := articleService.Review(article.ID, reviewer.ID, ReviewPayload{Action: "approve"})
	if err != nil {
		t.Fatalf("review article: %v", err)
	}
	if reviewed.Status != model.ArticlePublished {
		t.Fatalf("expected published, got %s", reviewed.Status)
	}
	if reviewed.RejectReason != "" {
		t.Fatalf("expected empty reject reason, got %q", reviewed.RejectReason)
	}
	if reviewed.PublishedAt == nil {
		t.Fatal("expected publishedAt to be set")
	}
}

func TestProjectReviewRejectFlow(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	projectService := NewProjectService(db)
	project, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Portfolio project",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	if _, err := projectService.Submit(project.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit project: %v", err)
	}

	reviewed, err := projectService.Review(project.ID, reviewer.ID, ProjectReviewPayload{
		Action: "reject",
		Reason: "内容还不完整",
	})
	if err != nil {
		t.Fatalf("review project: %v", err)
	}
	if reviewed.Status != model.ProjectRejected {
		t.Fatalf("expected rejected, got %s", reviewed.Status)
	}
	if reviewed.RejectReason != "内容还不完整" {
		t.Fatalf("unexpected reject reason: %q", reviewed.RejectReason)
	}
	if reviewed.PublishedAt != nil {
		t.Fatal("expected publishedAt to be nil after reject")
	}
}

func TestCommentCreatesNotificationsForAuthorAndParent(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	replyTarget := createTestUser(t, db, model.RoleUser)
	commenter := createTestUser(t, db, model.RoleUser)

	articleService := NewArticleService(db)
	commentService := NewCommentService(db)

	article, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Published article",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}

	parent, err := commentService.Create(article.ID, replyTarget.ID, CommentPayload{
		Content: "first comment",
	})
	if err != nil {
		t.Fatalf("create parent comment: %v", err)
	}

	_, err = commentService.Create(article.ID, commenter.ID, CommentPayload{
		Content:  "reply comment",
		ParentID: &parent.ID,
	})
	if err != nil {
		t.Fatalf("create reply comment: %v", err)
	}

	var notifications []model.Notification
	if err := db.Order("id asc").Find(&notifications).Error; err != nil {
		t.Fatalf("query notifications: %v", err)
	}
	if len(notifications) != 3 {
		t.Fatalf("expected 3 notifications, got %d", len(notifications))
	}

	var authorNotice, parentNotice *model.Notification
	for i := range notifications {
		item := &notifications[i]
		if item.UserID == author.ID && item.Type == model.NotificationTypeArticleComment {
			authorNotice = item
		}
		if item.UserID == replyTarget.ID && item.Type == model.NotificationTypeCommentReply {
			parentNotice = item
		}
	}

	if authorNotice == nil {
		t.Fatal("expected article author notification")
	}
	expectedURL := "/article/" + strconv.FormatUint(uint64(article.ID), 10)
	if authorNotice.ActionURL != expectedURL {
		t.Fatalf("unexpected author action url: %s", authorNotice.ActionURL)
	}
	if !strings.Contains(authorNotice.Content, article.Title) {
		t.Fatalf("expected author notification to mention article title, got %q", authorNotice.Content)
	}

	if parentNotice == nil {
		t.Fatal("expected parent comment reply notification")
	}
	if parentNotice.ActionURL != expectedURL {
		t.Fatalf("unexpected parent action url: %s", parentNotice.ActionURL)
	}
}

func TestCannotBanLastActiveAdmin(t *testing.T) {
	db := newTestDB(t)
	admin := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	_, err := adminService.UpdateUserStatus(admin.ID, model.UserBanned)
	if err == nil {
		t.Fatal("expected banning last active admin to fail")
	}
}

func TestCanBanAdminWhenAnotherActiveAdminExists(t *testing.T) {
	db := newTestDB(t)
	adminA := createTestUser(t, db, model.RoleAdmin)
	backup := model.User{
		Username: "backup_admin_" + strings.ReplaceAll(t.Name(), "/", "_"),
		Password: "hashed",
		Role:     model.RoleAdmin,
		Status:   model.UserActive,
	}
	if err := db.Create(&backup).Error; err != nil {
		t.Fatalf("create backup admin: %v", err)
	}

	adminService := NewAdminService(db, config.Config{})
	updated, err := adminService.UpdateUserStatus(adminA.ID, model.UserBanned)
	if err != nil {
		t.Fatalf("ban admin with backup present: %v", err)
	}
	if updated.Status != model.UserBanned {
		t.Fatalf("expected banned status, got %s", updated.Status)
	}
}

func TestForcePublishArticleCreatesReviewNotification(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	articleService := NewArticleService(db)
	adminService := NewAdminService(db, config.Config{})

	article, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Admin publish article",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}

	published, err := adminService.ForcePublishArticle(article.ID)
	if err != nil {
		t.Fatalf("force publish article: %v", err)
	}
	if published.Status != model.ArticlePublished {
		t.Fatalf("expected published, got %s", published.Status)
	}

	var notices []model.Notification
	if err := db.Where("user_id = ?", author.ID).Find(&notices).Error; err != nil {
		t.Fatalf("query notifications: %v", err)
	}
	if len(notices) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notices))
	}
	if notices[0].Type != model.NotificationTypeArticleReview {
		t.Fatalf("unexpected notification type: %s", notices[0].Type)
	}
	expectedURL := "/article/" + strconv.FormatUint(uint64(article.ID), 10)
	if notices[0].ActionURL != expectedURL {
		t.Fatalf("unexpected action url: %s", notices[0].ActionURL)
	}
}

func TestForcePublishProjectCreatesReviewNotification(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	projectService := NewProjectService(db)
	adminService := NewAdminService(db, config.Config{})

	project, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Admin publish project",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	published, err := adminService.ForcePublishProject(project.ID)
	if err != nil {
		t.Fatalf("force publish project: %v", err)
	}
	if published.Status != model.ProjectPublished {
		t.Fatalf("expected published, got %s", published.Status)
	}

	var notices []model.Notification
	if err := db.Where("user_id = ?", author.ID).Find(&notices).Error; err != nil {
		t.Fatalf("query notifications: %v", err)
	}
	if len(notices) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notices))
	}
	if notices[0].Type != model.NotificationTypeProjectReview {
		t.Fatalf("unexpected notification type: %s", notices[0].Type)
	}
	expectedURL := "/projects/" + strconv.FormatUint(uint64(project.ID), 10)
	if notices[0].ActionURL != expectedURL {
		t.Fatalf("unexpected action url: %s", notices[0].ActionURL)
	}
}

func TestUpdateProjectMetaUpdatesPublishedProject(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	projectService := NewProjectService(db)
	adminService := NewAdminService(db, config.Config{})

	project, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Published project meta",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	published, err := adminService.ForcePublishProject(project.ID)
	if err != nil {
		t.Fatalf("force publish project: %v", err)
	}

	updated, err := adminService.UpdateProjectMeta(published.ID, ProjectMetaPayload{
		IsFeatured: true,
		SortOrder:  9,
	})
	if err != nil {
		t.Fatalf("update project meta: %v", err)
	}
	if !updated.IsFeatured {
		t.Fatal("expected project to be featured")
	}
	if updated.SortOrder != 9 {
		t.Fatalf("expected sort order 9, got %d", updated.SortOrder)
	}
	if updated.Status != model.ProjectPublished {
		t.Fatalf("expected project to stay published, got %s", updated.Status)
	}
}

func TestUpdateProjectMetaRejectsUnpublishedProject(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	projectService := NewProjectService(db)
	adminService := NewAdminService(db, config.Config{})

	project, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Draft project meta",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	if _, err := adminService.UpdateProjectMeta(project.ID, ProjectMetaPayload{
		IsFeatured: true,
		SortOrder:  12,
	}); err == nil {
		t.Fatal("expected update project meta to fail for draft project")
	}

	reloaded, err := projectService.GetByID(project.ID, author.ID, model.RoleUser)
	if err != nil {
		t.Fatalf("reload project: %v", err)
	}
	if reloaded.IsFeatured {
		t.Fatal("expected draft project to remain unfeatured")
	}
	if reloaded.SortOrder != 0 {
		t.Fatalf("expected sort order 0, got %d", reloaded.SortOrder)
	}
}

func TestForcePublishProjectClearsRejectReason(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	projectService := NewProjectService(db)
	adminService := NewAdminService(db, config.Config{})

	project, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Rejected project republish",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	if _, err := projectService.Submit(project.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit project: %v", err)
	}
	if _, err := projectService.Review(project.ID, reviewer.ID, ProjectReviewPayload{
		Action: "reject",
		Reason: "needs work",
	}); err != nil {
		t.Fatalf("reject project: %v", err)
	}

	published, err := adminService.ForcePublishProject(project.ID)
	if err != nil {
		t.Fatalf("force publish project: %v", err)
	}
	if published.Status != model.ProjectPublished {
		t.Fatalf("expected published, got %s", published.Status)
	}
	if published.RejectReason != "" {
		t.Fatalf("expected reject reason to clear, got %q", published.RejectReason)
	}
	if published.PublishedAt == nil {
		t.Fatal("expected publishedAt to be set")
	}
}

func TestAdminListArticlesFiltersByTag(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	articleService := NewArticleService(db)
	adminService := NewAdminService(db, config.Config{})

	first, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Go article",
		Summary: "summary",
		Content: "content",
		Tags:    []string{"Go", "Backend"},
	})
	if err != nil {
		t.Fatalf("create first article: %v", err)
	}
	_, err = articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Vue article",
		Summary: "summary",
		Content: "content",
		Tags:    []string{"Vue"},
	})
	if err != nil {
		t.Fatalf("create second article: %v", err)
	}

	var goTag model.Tag
	if err := db.Where("name = ?", "Go").First(&goTag).Error; err != nil {
		t.Fatalf("load go tag: %v", err)
	}

	items, pagination, err := adminService.ListArticles(ArticleFilter{
		TagID:    &goTag.ID,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list articles by tag: %v", err)
	}
	if pagination.Total != 1 {
		t.Fatalf("expected 1 matching article, got %d", pagination.Total)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 returned article, got %d", len(items))
	}
	if items[0].ID != first.ID {
		t.Fatalf("expected article %d, got %d", first.ID, items[0].ID)
	}
}

func TestMetaServiceListIncludesTaxonomyUsageCounts(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	metaService := NewMetaService(db)
	articleService := NewArticleService(db)

	category, err := metaService.CreateCategory("后端实践")
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := metaService.CreateTag("Go"); err != nil {
		t.Fatalf("create tag: %v", err)
	}
	if _, err := metaService.CreateTag("Unused"); err != nil {
		t.Fatalf("create unused tag: %v", err)
	}

	if _, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:      "Tagged article",
		Summary:    "summary",
		Content:    "content",
		CategoryID: &category.ID,
		Tags:       []string{"Go"},
	}); err != nil {
		t.Fatalf("create article: %v", err)
	}

	data, err := metaService.List()
	if err != nil {
		t.Fatalf("list metadata: %v", err)
	}

	categories, ok := data["categories"].([]CategoryMetadata)
	if !ok {
		t.Fatalf("expected []CategoryMetadata, got %T", data["categories"])
	}
	tags, ok := data["tags"].([]TagMetadata)
	if !ok {
		t.Fatalf("expected []TagMetadata, got %T", data["tags"])
	}

	var categoryCount int64 = -1
	for _, item := range categories {
		if item.ID == category.ID {
			categoryCount = item.ArticleCount
			break
		}
	}
	if categoryCount != 1 {
		t.Fatalf("expected category article count 1, got %d", categoryCount)
	}

	counts := map[string]int64{}
	for _, item := range tags {
		counts[item.Name] = item.ArticleCount
	}
	if counts["Go"] != 1 {
		t.Fatalf("expected Go tag count 1, got %d", counts["Go"])
	}
	if counts["Unused"] != 0 {
		t.Fatalf("expected Unused tag count 0, got %d", counts["Unused"])
	}
}

func TestAdminUpdateArticleTaxonomyReassignsCategoryAndTags(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	metaService := NewMetaService(db)
	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	oldCategory, err := metaService.CreateCategory("后端")
	if err != nil {
		t.Fatalf("create old category: %v", err)
	}
	newCategory, err := metaService.CreateCategory("架构")
	if err != nil {
		t.Fatalf("create new category: %v", err)
	}
	if _, err := metaService.CreateTag("Go"); err != nil {
		t.Fatalf("create tag Go: %v", err)
	}
	if _, err := metaService.CreateTag("Gin"); err != nil {
		t.Fatalf("create tag Gin: %v", err)
	}

	article, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:      "Reassign taxonomy",
		Summary:    "summary",
		Content:    "content",
		CategoryID: &oldCategory.ID,
		Tags:       []string{"Go"},
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}

	var goTag model.Tag
	if err := db.Where("name = ?", "Go").First(&goTag).Error; err != nil {
		t.Fatalf("load Go tag: %v", err)
	}
	var ginTag model.Tag
	if err := db.Where("name = ?", "Gin").First(&ginTag).Error; err != nil {
		t.Fatalf("load Gin tag: %v", err)
	}

	updated, err := adminService.UpdateArticleTaxonomy(article.ID, ArticleTaxonomyPayload{
		CategoryID: &newCategory.ID,
		TagIDs:     []uint{ginTag.ID, ginTag.ID},
	})
	if err != nil {
		t.Fatalf("update taxonomy: %v", err)
	}
	if updated.Status != model.ArticlePublished {
		t.Fatalf("expected status to remain published, got %s", updated.Status)
	}
	if updated.Title != article.Title {
		t.Fatalf("expected title to stay %q, got %q", article.Title, updated.Title)
	}
	if updated.Category == nil || updated.Category.ID != newCategory.ID {
		t.Fatalf("expected category %d, got %+v", newCategory.ID, updated.Category)
	}
	if len(updated.Tags) != 1 || updated.Tags[0].ID != ginTag.ID {
		t.Fatalf("expected only Gin tag, got %+v", updated.Tags)
	}
	if len(updated.Tags) == 1 && updated.Tags[0].ID == goTag.ID {
		t.Fatalf("expected Go tag to be removed")
	}
}

func TestAdminUpdateArticleTaxonomyKeepsPendingStatusAndAllowsClearing(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	metaService := NewMetaService(db)
	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	category, err := metaService.CreateCategory("后端")
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := metaService.CreateTag("Go"); err != nil {
		t.Fatalf("create tag: %v", err)
	}

	article, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:      "Pending taxonomy",
		Summary:    "summary",
		Content:    "content",
		CategoryID: &category.ID,
		Tags:       []string{"Go"},
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	if _, err := articleService.Submit(article.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit article: %v", err)
	}

	updated, err := adminService.UpdateArticleTaxonomy(article.ID, ArticleTaxonomyPayload{
		CategoryID: nil,
		TagIDs:     nil,
	})
	if err != nil {
		t.Fatalf("clear taxonomy: %v", err)
	}
	if updated.Status != model.ArticlePending {
		t.Fatalf("expected status to remain pending, got %s", updated.Status)
	}
	if updated.Category != nil || updated.CategoryID != nil {
		t.Fatalf("expected category cleared, got %+v", updated.Category)
	}
	if len(updated.Tags) != 0 {
		t.Fatalf("expected tags cleared, got %+v", updated.Tags)
	}
}

func TestAdminBulkUpdateArticleTaxonomyUpdatesMultipleArticles(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	metaService := NewMetaService(db)
	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	category, err := metaService.CreateCategory("批量分类")
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := metaService.CreateTag("Go"); err != nil {
		t.Fatalf("create Go tag: %v", err)
	}
	if _, err := metaService.CreateTag("Backend"); err != nil {
		t.Fatalf("create Backend tag: %v", err)
	}

	first, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Bulk first",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create first article: %v", err)
	}
	second, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk second",
		Summary: "summary",
		Content: "content",
		Tags:    []string{"Go"},
	})
	if err != nil {
		t.Fatalf("create second article: %v", err)
	}
	if _, err := articleService.Submit(second.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit second article: %v", err)
	}

	var goTag model.Tag
	if err := db.Where("name = ?", "Go").First(&goTag).Error; err != nil {
		t.Fatalf("load Go tag: %v", err)
	}
	var backendTag model.Tag
	if err := db.Where("name = ?", "Backend").First(&backendTag).Error; err != nil {
		t.Fatalf("load Backend tag: %v", err)
	}

	updated, err := adminService.BulkUpdateArticleTaxonomy(BulkArticleTaxonomyPayload{
		ArticleIDs:  []uint{first.ID, second.ID, first.ID},
		CategoryID:  &category.ID,
		TagIDs:      []uint{goTag.ID, backendTag.ID},
		ReplaceTags: true,
	})
	if err != nil {
		t.Fatalf("bulk update taxonomy: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 updated articles, got %d", updated)
	}

	reloadedFirst, err := adminService.GetArticleDetail(first.ID)
	if err != nil {
		t.Fatalf("reload first article: %v", err)
	}
	reloadedSecond, err := adminService.GetArticleDetail(second.ID)
	if err != nil {
		t.Fatalf("reload second article: %v", err)
	}

	for _, item := range []*model.Article{reloadedFirst, reloadedSecond} {
		if item.Category == nil || item.Category.ID != category.ID {
			t.Fatalf("expected category %d, got %+v", category.ID, item.Category)
		}
		if len(item.Tags) != 2 {
			t.Fatalf("expected 2 tags, got %+v", item.Tags)
		}
	}
	if reloadedFirst.Status != model.ArticlePublished {
		t.Fatalf("expected first to stay published, got %s", reloadedFirst.Status)
	}
	if reloadedSecond.Status != model.ArticlePending {
		t.Fatalf("expected second to stay pending, got %s", reloadedSecond.Status)
	}
}

func TestAdminBulkPublishArticlesPublishesOnlyNonPublishedItems(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	draft, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk draft",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create draft article: %v", err)
	}
	pending, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk pending",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create pending article: %v", err)
	}
	if _, err := articleService.Submit(pending.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit pending article: %v", err)
	}
	published, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Bulk published",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create published article: %v", err)
	}

	updated, err := adminService.BulkPublishArticles(BulkArticleIDsPayload{
		ArticleIDs: []uint{draft.ID, pending.ID, published.ID, pending.ID},
	})
	if err != nil {
		t.Fatalf("bulk publish articles: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 updated articles, got %d", updated)
	}

	reloadedDraft, err := adminService.GetArticleDetail(draft.ID)
	if err != nil {
		t.Fatalf("reload draft article: %v", err)
	}
	reloadedPending, err := adminService.GetArticleDetail(pending.ID)
	if err != nil {
		t.Fatalf("reload pending article: %v", err)
	}
	reloadedPublished, err := adminService.GetArticleDetail(published.ID)
	if err != nil {
		t.Fatalf("reload published article: %v", err)
	}

	if reloadedDraft.Status != model.ArticlePublished {
		t.Fatalf("expected draft to become published, got %s", reloadedDraft.Status)
	}
	if reloadedPending.Status != model.ArticlePublished {
		t.Fatalf("expected pending to become published, got %s", reloadedPending.Status)
	}
	if reloadedPublished.Status != model.ArticlePublished {
		t.Fatalf("expected published to stay published, got %s", reloadedPublished.Status)
	}
}

func TestAdminBulkDeleteArticlesDeletesOnlyEligibleItems(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	draft, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk delete draft",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create draft article: %v", err)
	}
	rejected, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk delete rejected",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create rejected article: %v", err)
	}
	if _, err := articleService.Submit(rejected.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit rejected article: %v", err)
	}
	if _, err := articleService.Review(rejected.ID, reviewer.ID, ReviewPayload{Action: "reject", Reason: "needs work"}); err != nil {
		t.Fatalf("reject article: %v", err)
	}

	updated, err := adminService.BulkDeleteArticles(BulkArticleIDsPayload{
		ArticleIDs: []uint{draft.ID, rejected.ID, draft.ID},
	})
	if err != nil {
		t.Fatalf("bulk delete articles: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 deleted articles, got %d", updated)
	}
	if _, err := adminService.GetArticleDetail(draft.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected draft article to be deleted, got %v", err)
	}
	if _, err := adminService.GetArticleDetail(rejected.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected rejected article to be deleted, got %v", err)
	}
}

func TestAdminBulkDeleteArticlesRejectsUnsafeStatuses(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	draft, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Unsafe delete draft",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create draft article: %v", err)
	}
	pending, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Unsafe delete pending",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create pending article: %v", err)
	}
	if _, err := articleService.Submit(pending.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit pending article: %v", err)
	}

	if _, err := adminService.BulkDeleteArticles(BulkArticleIDsPayload{
		ArticleIDs: []uint{draft.ID, pending.ID},
	}); err == nil {
		t.Fatal("expected bulk delete to reject pending articles")
	}

	if _, err := adminService.GetArticleDetail(draft.ID); err != nil {
		t.Fatalf("expected draft article to remain after rejected bulk delete, got %v", err)
	}
	if item, err := adminService.GetArticleDetail(pending.ID); err != nil {
		t.Fatalf("expected pending article to remain after rejected bulk delete, got %v", err)
	} else if item.Status != model.ArticlePending {
		t.Fatalf("expected pending article to remain pending, got %s", item.Status)
	}
}

func TestAdminBulkRejectArticlesRejectsPendingItemsWithSharedReason(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	first, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk reject first",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create first article: %v", err)
	}
	second, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk reject second",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create second article: %v", err)
	}
	if _, err := articleService.Submit(first.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit first article: %v", err)
	}
	if _, err := articleService.Submit(second.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit second article: %v", err)
	}

	updated, err := adminService.BulkRejectArticles(reviewer.ID, BulkArticleRejectPayload{
		ArticleIDs: []uint{first.ID, second.ID, first.ID},
		Reason:     "内容需要统一补充案例",
	})
	if err != nil {
		t.Fatalf("bulk reject articles: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 rejected articles, got %d", updated)
	}

	reloadedFirst, err := adminService.GetArticleDetail(first.ID)
	if err != nil {
		t.Fatalf("reload first article: %v", err)
	}
	reloadedSecond, err := adminService.GetArticleDetail(second.ID)
	if err != nil {
		t.Fatalf("reload second article: %v", err)
	}
	for _, item := range []*model.Article{reloadedFirst, reloadedSecond} {
		if item.Status != model.ArticleRejected {
			t.Fatalf("expected rejected status, got %s", item.Status)
		}
		if item.RejectReason != "内容需要统一补充案例" {
			t.Fatalf("expected shared reject reason, got %q", item.RejectReason)
		}
	}
}

func TestAdminBulkRejectArticlesRejectsUnsafeStatuses(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	pending, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Pending bulk reject",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create pending article: %v", err)
	}
	published, err := articleService.Create(author.ID, model.RoleAdmin, ArticlePayload{
		Title:   "Published bulk reject",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create published article: %v", err)
	}
	if _, err := articleService.Submit(pending.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit pending article: %v", err)
	}

	if _, err := adminService.BulkRejectArticles(reviewer.ID, BulkArticleRejectPayload{
		ArticleIDs: []uint{pending.ID, published.ID},
		Reason:     "统一驳回",
	}); err == nil {
		t.Fatal("expected bulk reject to reject non-pending articles")
	}

	reloadedPending, err := adminService.GetArticleDetail(pending.ID)
	if err != nil {
		t.Fatalf("reload pending article: %v", err)
	}
	if reloadedPending.Status != model.ArticlePending {
		t.Fatalf("expected pending article to remain pending, got %s", reloadedPending.Status)
	}
	reloadedPublished, err := adminService.GetArticleDetail(published.ID)
	if err != nil {
		t.Fatalf("reload published article: %v", err)
	}
	if reloadedPublished.Status != model.ArticlePublished {
		t.Fatalf("expected published article to remain published, got %s", reloadedPublished.Status)
	}
}

func TestAdminBulkApproveArticlesApprovesPendingItems(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	first, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk approve first",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create first article: %v", err)
	}
	second, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Bulk approve second",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create second article: %v", err)
	}
	if _, err := articleService.Submit(first.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit first article: %v", err)
	}
	if _, err := articleService.Submit(second.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit second article: %v", err)
	}

	updated, err := adminService.BulkApproveArticles(reviewer.ID, BulkArticleIDsPayload{
		ArticleIDs: []uint{first.ID, second.ID, first.ID},
	})
	if err != nil {
		t.Fatalf("bulk approve articles: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 approved articles, got %d", updated)
	}

	reloadedFirst, err := adminService.GetArticleDetail(first.ID)
	if err != nil {
		t.Fatalf("reload first article: %v", err)
	}
	reloadedSecond, err := adminService.GetArticleDetail(second.ID)
	if err != nil {
		t.Fatalf("reload second article: %v", err)
	}
	for _, item := range []*model.Article{reloadedFirst, reloadedSecond} {
		if item.Status != model.ArticlePublished {
			t.Fatalf("expected published status, got %s", item.Status)
		}
		if item.PublishedAt == nil {
			t.Fatal("expected publishedAt to be set after bulk approve")
		}
	}
}

func TestAdminBulkApproveArticlesRejectsUnsafeStatuses(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	articleService := NewArticleService(db)

	pending, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Pending bulk approve",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create pending article: %v", err)
	}
	draft, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Draft bulk approve",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create draft article: %v", err)
	}
	if _, err := articleService.Submit(pending.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit pending article: %v", err)
	}

	if _, err := adminService.BulkApproveArticles(reviewer.ID, BulkArticleIDsPayload{
		ArticleIDs: []uint{pending.ID, draft.ID},
	}); err == nil {
		t.Fatal("expected bulk approve to reject non-pending articles")
	}

	reloadedPending, err := adminService.GetArticleDetail(pending.ID)
	if err != nil {
		t.Fatalf("reload pending article: %v", err)
	}
	if reloadedPending.Status != model.ArticlePending {
		t.Fatalf("expected pending article to remain pending, got %s", reloadedPending.Status)
	}
	reloadedDraft, err := adminService.GetArticleDetail(draft.ID)
	if err != nil {
		t.Fatalf("reload draft article: %v", err)
	}
	if reloadedDraft.Status != model.ArticleDraft {
		t.Fatalf("expected draft article to remain draft, got %s", reloadedDraft.Status)
	}
}

func TestAdminBulkReviewProjectsApprovesPendingItems(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	projectService := NewProjectService(db)

	first, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Bulk approve project first",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create first project: %v", err)
	}
	second, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Bulk approve project second",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create second project: %v", err)
	}
	if _, err := projectService.Submit(first.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit first project: %v", err)
	}
	if _, err := projectService.Submit(second.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit second project: %v", err)
	}

	updated, err := adminService.BulkReviewProjects(reviewer.ID, BulkProjectReviewPayload{
		ProjectIDs: []uint{first.ID, second.ID, first.ID},
		Action:     "approve",
	})
	if err != nil {
		t.Fatalf("bulk approve projects: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 approved projects, got %d", updated)
	}

	reloadedFirst, err := adminService.GetProjectDetail(first.ID)
	if err != nil {
		t.Fatalf("reload first project: %v", err)
	}
	reloadedSecond, err := adminService.GetProjectDetail(second.ID)
	if err != nil {
		t.Fatalf("reload second project: %v", err)
	}
	for _, item := range []*model.Project{reloadedFirst, reloadedSecond} {
		if item.Status != model.ProjectPublished {
			t.Fatalf("expected published status, got %s", item.Status)
		}
		if item.PublishedAt == nil {
			t.Fatal("expected publishedAt to be set after bulk approve")
		}
	}
}

func TestAdminBulkReviewProjectsRejectsPendingItemsWithSharedReason(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	projectService := NewProjectService(db)

	first, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Bulk reject project first",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create first project: %v", err)
	}
	second, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Bulk reject project second",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create second project: %v", err)
	}
	if _, err := projectService.Submit(first.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit first project: %v", err)
	}
	if _, err := projectService.Submit(second.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit second project: %v", err)
	}

	updated, err := adminService.BulkReviewProjects(reviewer.ID, BulkProjectReviewPayload{
		ProjectIDs: []uint{first.ID, second.ID},
		Action:     "reject",
		Reason:     "请补充项目成果与截图",
	})
	if err != nil {
		t.Fatalf("bulk reject projects: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 rejected projects, got %d", updated)
	}

	reloadedFirst, err := adminService.GetProjectDetail(first.ID)
	if err != nil {
		t.Fatalf("reload first project: %v", err)
	}
	reloadedSecond, err := adminService.GetProjectDetail(second.ID)
	if err != nil {
		t.Fatalf("reload second project: %v", err)
	}
	for _, item := range []*model.Project{reloadedFirst, reloadedSecond} {
		if item.Status != model.ProjectRejected {
			t.Fatalf("expected rejected status, got %s", item.Status)
		}
		if item.RejectReason != "请补充项目成果与截图" {
			t.Fatalf("expected shared reject reason, got %q", item.RejectReason)
		}
	}
}

func TestAdminBulkReviewProjectsRejectsUnsafeStatuses(t *testing.T) {
	db := newTestDB(t)
	author := createTestUser(t, db, model.RoleUser)
	reviewer := createTestUser(t, db, model.RoleAdmin)

	adminService := NewAdminService(db, config.Config{})
	projectService := NewProjectService(db)

	pending, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Pending bulk review project",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create pending project: %v", err)
	}
	draft, err := projectService.Create(author.ID, model.RoleUser, ProjectPayload{
		Title:   "Draft bulk review project",
		Summary: "summary",
		Content: "content",
	})
	if err != nil {
		t.Fatalf("create draft project: %v", err)
	}
	if _, err := projectService.Submit(pending.ID, author.ID, model.RoleUser); err != nil {
		t.Fatalf("submit pending project: %v", err)
	}

	if _, err := adminService.BulkReviewProjects(reviewer.ID, BulkProjectReviewPayload{
		ProjectIDs: []uint{pending.ID, draft.ID},
		Action:     "approve",
	}); err == nil {
		t.Fatal("expected bulk project review to reject non-pending projects")
	}

	reloadedPending, err := adminService.GetProjectDetail(pending.ID)
	if err != nil {
		t.Fatalf("reload pending project: %v", err)
	}
	if reloadedPending.Status != model.ProjectPending {
		t.Fatalf("expected pending project to remain pending, got %s", reloadedPending.Status)
	}
	reloadedDraft, err := adminService.GetProjectDetail(draft.ID)
	if err != nil {
		t.Fatalf("reload draft project: %v", err)
	}
	if reloadedDraft.Status != model.ProjectDraft {
		t.Fatalf("expected draft project to remain draft, got %s", reloadedDraft.Status)
	}
}
