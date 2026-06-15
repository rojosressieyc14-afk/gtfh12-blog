package service

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"blog/server/internal/config"
	"blog/server/internal/model"
	"gorm.io/gorm"
)

type AdminService struct {
	db  *gorm.DB
	cfg config.Config
}

type UserFilter struct {
	Keyword  string
	Role     string
	Page     int
	PageSize int
}

type ArticleFilter struct {
	Status     string
	Keyword    string
	CategoryID *uint
	TagID      *uint
	Page       int
	PageSize   int
}

type ProjectAdminFilter struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type ProjectMetaPayload struct {
	IsFeatured bool
	SortOrder  int
}

type ArticleTaxonomyPayload struct {
	CategoryID *uint  `json:"categoryId"`
	TagIDs     []uint `json:"tagIds"`
}

type BulkArticleTaxonomyPayload struct {
	ArticleIDs  []uint `json:"articleIds"`
	CategoryID  *uint  `json:"categoryId"`
	TagIDs      []uint `json:"tagIds"`
	ReplaceTags bool   `json:"replaceTags"`
}

type BulkArticleIDsPayload struct {
	ArticleIDs []uint `json:"articleIds"`
}

type BulkArticleRejectPayload struct {
	ArticleIDs []uint `json:"articleIds"`
	Reason     string `json:"reason"`
}

type BulkProjectReviewPayload struct {
	ProjectIDs []uint `json:"projectIds"`
	Action     string `json:"action"`
	Reason     string `json:"reason"`
}

type CommentFilter struct {
	Keyword  string
	Page     int
	PageSize int
}

type LogFilter struct {
	Keyword  string
	Action   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PageSize int
}

type SensitiveWordFilter struct {
	Keyword  string
	Page     int
	PageSize int
}

type SensitiveWordPayload struct {
	Word     string `json:"word"`
	Category string `json:"category"`
	Note     string `json:"note"`
}

type ModerationHitFilter struct {
	Keyword    string
	Scene      string
	AutoBanned bool
	Page       int
	PageSize   int
}

type ModerationSettingPayload struct {
	BanThreshold int `json:"banThreshold"`
}

type UploadAsset struct {
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Size int64     `json:"size"`
	Time time.Time `json:"time"`
}

type DashboardTrendPoint struct {
	Date              string `json:"date"`
	Label             string `json:"label"`
	NewUsers          int64  `json:"newUsers"`
	NewArticles       int64  `json:"newArticles"`
	PublishedArticles int64  `json:"publishedArticles"`
	NewComments       int64  `json:"newComments"`
	NewProjects       int64  `json:"newProjects"`
	PublishedProjects int64  `json:"publishedProjects"`
	ModerationHits    int64  `json:"moderationHits"`
}

type DashboardPayload struct {
	Stats  map[string]int64      `json:"stats"`
	Trends []DashboardTrendPoint `json:"trends"`
}

type groupedCount struct {
	Day   string `json:"day"`
	Total int64  `json:"total"`
}

func NewAdminService(db *gorm.DB, cfg config.Config) *AdminService {
	return &AdminService{db: db, cfg: cfg}
}

func (s *AdminService) Dashboard() (*DashboardPayload, error) {
	stats, err := s.dashboardStats()
	if err != nil {
		return nil, err
	}

	trends, err := s.dashboardTrends()
	if err != nil {
		return nil, err
	}

	return &DashboardPayload{Stats: stats, Trends: trends}, nil
}

func (s *AdminService) dashboardStats() (map[string]int64, error) {
	result := map[string]int64{}
	statuses := []string{
		model.ArticleDraft,
		model.ArticlePending,
		model.ArticlePublished,
		model.ArticleRejected,
	}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&model.Article{}).Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		result[status] = count
	}

	var users int64
	if err := s.db.Model(&model.User{}).Count(&users).Error; err != nil {
		return nil, err
	}
	result["users"] = users

	var comments int64
	if err := s.db.Model(&model.Comment{}).Count(&comments).Error; err != nil {
		return nil, err
	}
	result["comments"] = comments

	var categories int64
	if err := s.db.Model(&model.Category{}).Count(&categories).Error; err != nil {
		return nil, err
	}
	result["categories"] = categories

	var tags int64
	if err := s.db.Model(&model.Tag{}).Count(&tags).Error; err != nil {
		return nil, err
	}
	result["tags"] = tags

	var projects int64
	if err := s.db.Model(&model.Project{}).Where("status = ?", model.ProjectPublished).Count(&projects).Error; err != nil {
		return nil, err
	}
	result["projects"] = projects

	var pendingProjects int64
	if err := s.db.Model(&model.Project{}).Where("status = ?", model.ProjectPending).Count(&pendingProjects).Error; err != nil {
		return nil, err
	}
	result["pendingProjects"] = pendingProjects

	var draftProjects int64
	if err := s.db.Model(&model.Project{}).Where("status = ?", model.ProjectDraft).Count(&draftProjects).Error; err != nil {
		return nil, err
	}
	result["draftProjects"] = draftProjects

	var moderationHits int64
	if err := s.db.Model(&model.ModerationHit{}).Count(&moderationHits).Error; err != nil {
		return nil, err
	}
	result["moderationHits"] = moderationHits

	var autoBannedUsers int64
	if err := s.db.Model(&model.User{}).
		Where("status = ? AND ban_reason LIKE ?", model.UserBanned, autoBanReasonPrefix+"%").
		Count(&autoBannedUsers).Error; err != nil {
		return nil, err
	}
	result["autoBannedUsers"] = autoBannedUsers

	var sensitiveWords int64
	if err := s.db.Model(&model.SensitiveWord{}).Where("is_enabled = ?", true).Count(&sensitiveWords).Error; err != nil {
		return nil, err
	}
	result["sensitiveWords"] = sensitiveWords

	return result, nil
}

func (s *AdminService) dashboardTrends() ([]DashboardTrendPoint, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	userMap, err := s.countByDate(&model.User{}, "created_at", start, nil)
	if err != nil {
		return nil, err
	}
	articleMap, err := s.countByDate(&model.Article{}, "created_at", start, nil)
	if err != nil {
		return nil, err
	}
	commentMap, err := s.countByDate(&model.Comment{}, "created_at", start, nil)
	if err != nil {
		return nil, err
	}
	projectMap, err := s.countByDate(&model.Project{}, "created_at", start, nil)
	if err != nil {
		return nil, err
	}
	moderationHitMap, err := s.countByDate(&model.ModerationHit{}, "created_at", start, nil)
	if err != nil {
		return nil, err
	}
	publishedArticleMap, err := s.countByDate(&model.Article{}, "published_at", start, func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ? AND published_at IS NOT NULL", model.ArticlePublished)
	})
	if err != nil {
		return nil, err
	}
	publishedProjectMap, err := s.countByDate(&model.Project{}, "published_at", start, func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ? AND published_at IS NOT NULL", model.ProjectPublished)
	})
	if err != nil {
		return nil, err
	}

	trends := make([]DashboardTrendPoint, 0, 7)
	for offset := 0; offset < 7; offset++ {
		day := start.AddDate(0, 0, offset)
		key := day.Format("2006-01-02")
		trends = append(trends, DashboardTrendPoint{
			Date:              key,
			Label:             day.Format("01/02"),
			NewUsers:          userMap[key],
			NewArticles:       articleMap[key],
			PublishedArticles: publishedArticleMap[key],
			NewComments:       commentMap[key],
			NewProjects:       projectMap[key],
			PublishedProjects: publishedProjectMap[key],
			ModerationHits:    moderationHitMap[key],
		})
	}

	return trends, nil
}

func (s *AdminService) countByDate(target interface{}, column string, start time.Time, scope func(*gorm.DB) *gorm.DB) (map[string]int64, error) {
	query := s.db.Model(target).Where(column+" >= ?", start)
	if scope != nil {
		query = scope(query)
	}

	var rows []groupedCount
	if err := query.
		Select("DATE(" + column + ") AS day, COUNT(*) AS total").
		Group("DATE(" + column + ")").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]int64, len(rows))
	for _, row := range rows {
		result[row.Day] = row.Total
	}
	return result, nil
}

func (s *AdminService) ListUsers(filter UserFilter) ([]model.User, Pagination, error) {
	var users []model.User
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.User{}).Order("created_at desc")
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("username LIKE ? OR role LIKE ? OR status LIKE ?", like, like, like)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&users).Error
	pagination.Total = total
	return users, pagination, err
}

func (s *AdminService) ListArticles(filter ArticleFilter) ([]model.Article, Pagination, error) {
	var articles []model.Article
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.Article{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Joins("LEFT JOIN users ON users.id = articles.author_id").
		Order("articles.updated_at desc")

	if filter.Status != "" {
		query = query.Where("articles.status = ?", filter.Status)
	}
	if filter.CategoryID != nil {
		query = query.Where("articles.category_id = ?", *filter.CategoryID)
	}
	if filter.TagID != nil {
		query = query.
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", *filter.TagID).
			Distinct("articles.id")
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("articles.title LIKE ? OR articles.summary LIKE ? OR users.username LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&articles).Error
	pagination.Total = total
	return articles, pagination, err
}

func (s *AdminService) GetArticleDetail(articleID uint) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Author").Preload("Category").Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *AdminService) UpdateArticleTaxonomy(articleID uint, payload ArticleTaxonomyPayload) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}

	var category *model.Category
	if payload.CategoryID != nil {
		category = &model.Category{}
		if err := s.db.First(category, *payload.CategoryID).Error; err != nil {
			return nil, err
		}
	}

	tagIDs := uniqueUintIDs(payload.TagIDs)
	tags := make([]model.Tag, 0, len(tagIDs))
	if len(tagIDs) > 0 {
		if err := s.db.Where("id IN ?", tagIDs).Order("name asc").Find(&tags).Error; err != nil {
			return nil, err
		}
		if len(tags) != len(tagIDs) {
			return nil, gorm.ErrRecordNotFound
		}
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&article).Update("category_id", payload.CategoryID).Error; err != nil {
			return err
		}
		return tx.Model(&article).Association("Tags").Replace(tags)
	})
	if err != nil {
		return nil, err
	}

	return s.GetArticleDetail(articleID)
}

func (s *AdminService) BulkUpdateArticleTaxonomy(payload BulkArticleTaxonomyPayload) (int, error) {
	articleIDs := uniqueUintIDs(payload.ArticleIDs)
	if len(articleIDs) == 0 {
		return 0, ErrNoArticleIDs
	}

	var category *model.Category
	if payload.CategoryID != nil {
		category = &model.Category{}
		if err := s.db.First(category, *payload.CategoryID).Error; err != nil {
			return 0, err
		}
	}

	tagIDs := uniqueUintIDs(payload.TagIDs)
	tags := make([]model.Tag, 0, len(tagIDs))
	if payload.ReplaceTags && len(tagIDs) > 0 {
		if err := s.db.Where("id IN ?", tagIDs).Order("name asc").Find(&tags).Error; err != nil {
			return 0, err
		}
		if len(tags) != len(tagIDs) {
			return 0, gorm.ErrRecordNotFound
		}
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var articles []model.Article
		if err := tx.Preload("Tags").Where("id IN ?", articleIDs).Find(&articles).Error; err != nil {
			return err
		}
		if len(articles) != len(articleIDs) {
			return gorm.ErrRecordNotFound
		}

		for index := range articles {
			article := &articles[index]
			if payload.CategoryID != nil || payload.CategoryID == nil {
				if err := tx.Model(article).Update("category_id", payload.CategoryID).Error; err != nil {
					return err
				}
			}
			if payload.ReplaceTags {
				if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return len(articleIDs), nil
}

func (s *AdminService) BulkPublishArticles(payload BulkArticleIDsPayload) (int, error) {
	articleIDs := uniqueUintIDs(payload.ArticleIDs)
	if len(articleIDs) == 0 {
		return 0, ErrNoArticleIDs
	}

	updated := 0
	for _, articleID := range articleIDs {
		item, err := s.GetArticleDetail(articleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, err
			}
			return 0, err
		}
		if item.Status == model.ArticlePublished {
			continue
		}
		if _, err := s.ForcePublishArticle(articleID); err != nil {
			return 0, err
		}
		updated++
	}

	return updated, nil
}

func (s *AdminService) BulkDeleteArticles(payload BulkArticleIDsPayload) (int, error) {
	articleIDs := uniqueUintIDs(payload.ArticleIDs)
	if len(articleIDs) == 0 {
		return 0, ErrNoArticleIDs
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var articles []model.Article
		if err := tx.Preload("Tags").Where("id IN ?", articleIDs).Find(&articles).Error; err != nil {
			return err
		}
		if len(articles) != len(articleIDs) {
			return gorm.ErrRecordNotFound
		}

		for index := range articles {
			switch articles[index].Status {
			case model.ArticleDraft, model.ArticleRejected:
			default:
				return ErrDeleteNotDraft
			}
		}

		for index := range articles {
			if err := s.deleteArticleTx(tx, &articles[index]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return len(articleIDs), nil
}

func (s *AdminService) BulkRejectArticles(reviewerID uint, payload BulkArticleRejectPayload) (int, error) {
	articleIDs := uniqueUintIDs(payload.ArticleIDs)
	if len(articleIDs) == 0 {
		return 0, ErrNoArticleIDs
	}
	reason := strings.TrimSpace(payload.Reason)
	if reason == "" {
		return 0, ErrRejectReasonMissing
	}

	var articles []model.Article
	if err := s.db.Where("id IN ?", articleIDs).Find(&articles).Error; err != nil {
		return 0, err
	}
	if len(articles) != len(articleIDs) {
		return 0, gorm.ErrRecordNotFound
	}
	for index := range articles {
		if articles[index].Status != model.ArticlePending {
			return 0, ErrBulkNotPending
		}
	}

	articleService := NewArticleService(s.db)
	for _, articleID := range articleIDs {
		if _, err := articleService.Review(articleID, reviewerID, ReviewPayload{
			Action: "reject",
			Reason: reason,
		}); err != nil {
			return 0, err
		}
	}

	return len(articleIDs), nil
}

func (s *AdminService) BulkApproveArticles(reviewerID uint, payload BulkArticleIDsPayload) (int, error) {
	articleIDs := uniqueUintIDs(payload.ArticleIDs)
	if len(articleIDs) == 0 {
		return 0, ErrNoArticleIDs
	}

	var articles []model.Article
	if err := s.db.Where("id IN ?", articleIDs).Find(&articles).Error; err != nil {
		return 0, err
	}
	if len(articles) != len(articleIDs) {
		return 0, gorm.ErrRecordNotFound
	}
	for index := range articles {
		if articles[index].Status != model.ArticlePending {
			return 0, ErrBulkApproveNotPending
		}
	}

	articleService := NewArticleService(s.db)
	for _, articleID := range articleIDs {
		if _, err := articleService.Review(articleID, reviewerID, ReviewPayload{
			Action: "approve",
			Reason: "",
		}); err != nil {
			return 0, err
		}
	}

	return len(articleIDs), nil
}

func (s *AdminService) BulkReviewProjects(reviewerID uint, payload BulkProjectReviewPayload) (int, error) {
	projectIDs := uniqueUintIDs(payload.ProjectIDs)
	if len(projectIDs) == 0 {
		return 0, ErrNoProjectIDs
	}
	action := strings.TrimSpace(payload.Action)
	if action != "approve" && action != "reject" {
		return 0, ErrInvalidReviewAction
	}
	reason := strings.TrimSpace(payload.Reason)
	if action == "reject" && reason == "" {
		return 0, ErrProjectRejectReason
	}

	var projects []model.Project
	if err := s.db.Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
		return 0, err
	}
	if len(projects) != len(projectIDs) {
		return 0, gorm.ErrRecordNotFound
	}
	for index := range projects {
		if projects[index].Status != model.ProjectPending {
			return 0, ErrBulkProjectNotPending
		}
	}

	projectService := NewProjectService(s.db)
	for _, projectID := range projectIDs {
		if _, err := projectService.Review(projectID, reviewerID, ProjectReviewPayload{
			Action: action,
			Reason: reason,
		}); err != nil {
			return 0, err
		}
	}

	return len(projectIDs), nil
}

func (s *AdminService) ListProjects(filter ProjectAdminFilter) ([]model.Project, Pagination, error) {
	var items []model.Project
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.Project{}).Preload("Author").Order("updated_at desc")
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&items).Error
	pagination.Total = total
	return items, pagination, err
}

func (s *AdminService) GetProjectDetail(projectID uint) (*model.Project, error) {
	var item model.Project
	if err := s.db.Preload("Author").First(&item, projectID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *AdminService) ListComments(filter CommentFilter) ([]model.Comment, Pagination, error) {
	var comments []model.Comment
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.Comment{}).
		Preload("User").
		Preload("Article").
		Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at asc") }).
		Preload("Replies.User").
		Joins("LEFT JOIN users ON users.id = comments.user_id").
		Joins("LEFT JOIN articles ON articles.id = comments.article_id").
		Where("comments.parent_id IS NULL").
		Order("comments.created_at desc")

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("comments.content LIKE ? OR users.username LIKE ? OR articles.title LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&comments).Error
	pagination.Total = total
	return comments, pagination, err
}

func (s *AdminService) ListOperationLogs(filter LogFilter) ([]model.OperationLog, Pagination, error) {
	var items []model.OperationLog
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.OperationLog{}).
		Preload("Operator").
		Joins("LEFT JOIN users ON users.id = operation_logs.operator_id").
		Order("operation_logs.created_at desc")

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("operation_logs.action LIKE ? OR operation_logs.description LIKE ? OR users.username LIKE ?", like, like, like)
	}
	if filter.Action != "" {
		query = query.Where("operation_logs.action = ?", filter.Action)
	}
	if filter.DateFrom != nil {
		query = query.Where("operation_logs.created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("operation_logs.created_at < ?", *filter.DateTo)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&items).Error
	pagination.Total = total
	return items, pagination, err
}

func (s *AdminService) ListUploads() ([]UploadAsset, error) {
	if err := os.MkdirAll(s.cfg.UploadDir, 0o755); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.cfg.UploadDir)
	if err != nil {
		return nil, err
	}

	assets := make([]UploadAsset, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		assets = append(assets, UploadAsset{
			Name: entry.Name(),
			URL:  "/uploads/" + entry.Name(),
			Size: info.Size(),
			Time: info.ModTime(),
		})
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Time.After(assets[j].Time)
	})
	return assets, nil
}

func (s *AdminService) DeleteUpload(name string) error {
	if name == "" {
		return ErrUploadNameEmpty
	}
	if filepath.Base(name) != name {
		return ErrUploadNameInvalid
	}

	target := filepath.Join(s.cfg.UploadDir, name)
	if err := os.Remove(target); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrUploadNotFound
		}
		return err
	}
	return nil
}

func (s *AdminService) UpdateUserRole(userID uint, role string) (*model.User, error) {
	if role != model.RoleAdmin && role != model.RoleUser {
		return nil, ErrInvalidRole
	}

	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&user).Update("role", role).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AdminService) UpdateUserStatus(userID uint, status string) (*model.User, error) {
	if status != model.UserActive && status != model.UserBanned {
		return nil, ErrInvalidStatus
	}

	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	if user.Role == model.RoleAdmin && status == model.UserBanned {
		var adminCount int64
		if err := s.db.Model(&model.User{}).Where("role = ? AND status = ?", model.RoleAdmin, model.UserActive).Count(&adminCount).Error; err != nil {
			return nil, err
		}
		if adminCount <= 1 {
			return nil, ErrLastAdminBan
		}
	}

	updates := map[string]any{"status": status}
	if status == model.UserActive {
		updates["ban_reason"] = ""
	}
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AdminService) DeleteUser(userID uint) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}
	if user.Role == model.RoleAdmin {
		var adminCount int64
		if err := s.db.Model(&model.User{}).Where("role = ? AND status = ?", model.RoleAdmin, model.UserActive).Count(&adminCount).Error; err != nil {
			return err
		}
		if adminCount <= 1 {
			return ErrLastAdminDelete
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("author_id = ?", userID).Delete(&model.Article{})
		var commentIDs []uint
		tx.Model(&model.Comment{}).Where("user_id = ?", userID).Pluck("id", &commentIDs)
		if len(commentIDs) > 0 {
			tx.Where("parent_id IN ?", commentIDs).Delete(&model.Comment{})
		}
		tx.Where("user_id = ?", userID).Delete(&model.Comment{})
		tx.Where("user_id = ?", userID).Delete(&model.ArticleLike{})
		tx.Where("user_id = ?", userID).Delete(&model.ArticleFavorite{})
		tx.Where("user_id = ?", userID).Delete(&model.ModerationHit{})
		tx.Where("operator_id = ?", userID).Delete(&model.OperationLog{})
		tx.Where("user_id = ?", userID).Delete(&model.Notification{})
		tx.Where("author_id = ?", userID).Delete(&model.Project{})
		return tx.Delete(&user).Error
	})
}

func (s *AdminService) DeleteComment(commentID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.Comment{}).Where("id = ? OR parent_id = ?", commentID, commentID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Where("id = ? OR parent_id = ?", commentID, commentID).Delete(&model.Comment{}).Error
	})
}

func (s *AdminService) DeleteArticle(articleID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var article model.Article
		if err := tx.Preload("Tags").First(&article, articleID).Error; err != nil {
			return err
		}
		return s.deleteArticleTx(tx, &article)
	})
}

func (s *AdminService) ForcePublishArticle(articleID uint) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if err := validateModerationField("文章标题", article.Title); err != nil {
		reason := "内容命中敏感词，系统已自动驳回：" + moderationWord(err)
		if updateErr := s.db.Model(&article).Updates(map[string]any{
			"status":        model.ArticleRejected,
			"published_at":  nil,
			"reject_reason": reason,
		}).Error; updateErr != nil {
			return nil, updateErr
		}
		return s.GetArticleDetail(articleID)
	}

	now := time.Now()
	if err := s.db.Model(&article).Updates(map[string]any{
		"status":        model.ArticlePublished,
		"published_at":  &now,
		"reject_reason": "",
	}).Error; err != nil {
		return nil, err
	}

	_ = createNotification(s.db, NotificationCreateInput{
		UserID:    article.AuthorID,
		Title:     "文章已发布",
		Content:   "你的文章《" + article.Title + "》已被管理员直接发布。",
		Type:      model.NotificationTypeArticleReview,
		ActionURL: "/article/" + strconv.FormatUint(uint64(article.ID), 10),
		Payload: map[string]any{
			"articleId": article.ID,
			"action":    "approve",
		},
	})

	var refreshed model.Article
	if err := s.db.Preload("Author").Preload("Category").Preload("Tags").First(&refreshed, articleID).Error; err != nil {
		return nil, err
	}
	return &refreshed, nil
}

func (s *AdminService) DeleteProject(projectID uint) error {
	return s.db.Delete(&model.Project{}, projectID).Error
}

func (s *AdminService) ForcePublishProject(projectID uint) (*model.Project, error) {
	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if err := validateModerationField("项目标题", item.Title); err != nil {
		reason := "内容命中敏感词，系统已自动驳回：" + moderationWord(err)
		if updateErr := s.db.Model(&item).Updates(map[string]any{
			"status":        model.ProjectRejected,
			"published_at":  nil,
			"reject_reason": reason,
		}).Error; updateErr != nil {
			return nil, updateErr
		}
		return s.GetProjectDetail(projectID)
	}

	now := time.Now()
	if err := s.db.Model(&item).Updates(map[string]any{
		"status":        model.ProjectPublished,
		"published_at":  &now,
		"reject_reason": "",
	}).Error; err != nil {
		return nil, err
	}

	_ = createNotification(s.db, NotificationCreateInput{
		UserID:    item.AuthorID,
		Title:     "项目已发布",
		Content:   "你的项目《" + item.Title + "》已被管理员直接发布。",
		Type:      model.NotificationTypeProjectReview,
		ActionURL: "/projects/" + strconv.FormatUint(uint64(item.ID), 10),
		Payload: map[string]any{
			"projectId": item.ID,
			"action":    "approve",
		},
	})

	var refreshed model.Project
	if err := s.db.Preload("Author").First(&refreshed, projectID).Error; err != nil {
		return nil, err
	}
	return &refreshed, nil
}

func (s *AdminService) UpdateProjectMeta(projectID uint, payload ProjectMetaPayload) (*model.Project, error) {
	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if item.Status != model.ProjectPublished {
		return nil, ErrPublishedMetaOnly
	}

	if err := s.db.Model(&item).Updates(map[string]any{
		"is_featured": payload.IsFeatured,
		"sort_order":  payload.SortOrder,
	}).Error; err != nil {
		return nil, err
	}

	var refreshed model.Project
	if err := s.db.Preload("Author").First(&refreshed, projectID).Error; err != nil {
		return nil, err
	}
	return &refreshed, nil
}

func uniqueUintIDs(values []uint) []uint {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[uint]struct{}, len(values))
	result := make([]uint, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func (s *AdminService) deleteArticleTx(tx *gorm.DB, article *model.Article) error {
	if err := tx.Model(article).Association("Tags").Clear(); err != nil {
		return err
	}
	if err := tx.Where("article_id = ?", article.ID).Delete(&model.ArticleReview{}).Error; err != nil {
		return err
	}
	if err := tx.Where("article_id = ?", article.ID).Delete(&model.Comment{}).Error; err != nil {
		return err
	}
	if err := tx.Where("article_id = ?", article.ID).Delete(&model.ArticleLike{}).Error; err != nil {
		return err
	}
	if err := tx.Where("article_id = ?", article.ID).Delete(&model.ArticleFavorite{}).Error; err != nil {
		return err
	}
	return tx.Delete(article).Error
}

func (s *AdminService) CreateNotification(input NotificationCreateInput) error {
	return createNotification(s.db, input)
}

func (s *AdminService) CreateOperationLog(operatorID uint, action, targetType string, targetID uint, description string) error {
	return s.db.Create(&model.OperationLog{
		OperatorID:  operatorID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Description: description,
	}).Error
}

func (s *AdminService) UploadURL(name string) string {
	return filepath.ToSlash("/uploads/" + name)
}

func (s *AdminService) ListModerationHits(filter ModerationHitFilter) ([]model.ModerationHit, Pagination, error) {
	var items []model.ModerationHit
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.ModerationHit{}).Preload("User").Order("created_at desc")
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Joins("LEFT JOIN users ON users.id = moderation_hits.user_id").
			Where("moderation_hits.field LIKE ? OR moderation_hits.matched_word LIKE ? OR moderation_hits.snippet LIKE ? OR users.username LIKE ?", like, like, like, like)
	}
	if filter.Scene != "" {
		query = query.Where("scene = ?", filter.Scene)
	}
	if filter.AutoBanned {
		query = query.Joins("JOIN users banned_users ON banned_users.id = moderation_hits.user_id").
			Where("banned_users.status = ? AND banned_users.ban_reason LIKE ?", model.UserBanned, autoBanReasonPrefix+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	if err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&items).Error; err != nil {
		return nil, Pagination{}, err
	}
	pagination.Total = total
	return items, pagination, nil
}

func (s *AdminService) ListSensitiveWords(filter SensitiveWordFilter) ([]model.SensitiveWord, Pagination, error) {
	var items []model.SensitiveWord
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.SensitiveWord{}).Order("created_at desc")
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("word LIKE ? OR category LIKE ? OR note LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	if err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&items).Error; err != nil {
		return nil, Pagination{}, err
	}
	pagination.Total = total
	return items, pagination, nil
}

func (s *AdminService) CreateSensitiveWord(payload SensitiveWordPayload) (*model.SensitiveWord, error) {
	word := strings.TrimSpace(payload.Word)
	if word == "" {
		return nil, ErrSensitiveWordEmpty
	}

	normalized := normalizeModerationText(word)
	if normalized == "" {
		return nil, ErrSensitiveWordInvalid
	}

	item := model.SensitiveWord{
		Word:      normalized,
		Category:  strings.TrimSpace(payload.Category),
		Note:      strings.TrimSpace(payload.Note),
		IsEnabled: true,
	}
	if item.Category == "" {
		item.Category = "custom"
	}

	var exists model.SensitiveWord
	if err := s.db.Where("word = ?", item.Word).First(&exists).Error; err == nil {
		return nil, ErrSensitiveWordExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	if err := LoadSensitiveWords(s.db); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *AdminService) DeleteSensitiveWord(id uint) error {
	var item model.SensitiveWord
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&item).Error; err != nil {
		return err
	}
	return LoadSensitiveWords(s.db)
}

func (s *AdminService) GetModerationSettings() (map[string]int, error) {
	if err := LoadModerationSettings(s.db); err != nil {
		return nil, err
	}
	return map[string]int{
		"banThreshold": currentModerationBanThreshold(),
	}, nil
}

func (s *AdminService) UpdateModerationSettings(payload ModerationSettingPayload) (map[string]int, error) {
	if payload.BanThreshold <= 0 || payload.BanThreshold > 100 {
		return nil, ErrBanThresholdRange
	}

	var item model.SystemSetting
	err := s.db.Where("`key` = ?", moderationBanThresholdKey).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = model.SystemSetting{
			Key:   moderationBanThresholdKey,
			Value: strconv.Itoa(payload.BanThreshold),
			Note:  "24 小时内触发多少次违禁词后自动封禁普通用户",
		}
		if err := s.db.Create(&item).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		if err := s.db.Model(&item).Update("value", strconv.Itoa(payload.BanThreshold)).Error; err != nil {
			return nil, err
		}
	}

	setModerationBanThreshold(payload.BanThreshold)
	return map[string]int{
		"banThreshold": payload.BanThreshold,
	}, nil
}
