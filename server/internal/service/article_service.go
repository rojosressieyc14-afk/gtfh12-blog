package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

type ArticlePayload struct {
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	CoverImage string   `json:"coverImage"`
	CategoryID *uint    `json:"categoryId"`
	Tags       []string `json:"tags"`
}

type ReviewPayload struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type ReactionSummary struct {
	LikesCount     int64 `json:"likesCount"`
	FavoritesCount int64 `json:"favoritesCount"`
	IsLiked        bool  `json:"isLiked"`
	IsFavorited    bool  `json:"isFavorited"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type PublishedArticleFilter struct {
	Keyword    string
	CategoryID *uint
	Tag        string
	AuthorID   uint
	Sort       string
	Page       int
	PageSize   int
}

type ArticleService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{db: db}
}

func (s *ArticleService) Create(authorID uint, role string, payload ArticlePayload) (*model.Article, error) {
	var author model.User
	if err := s.db.First(&author, authorID).Error; err != nil {
		return nil, err
	}
	if author.Status == model.UserBanned {
		return nil, ErrUserBanned
	}
	if err := validateModerationField("文章标题", payload.Title); err != nil {
		createModerationHit(s.db, authorID, "article_create", "文章标题", payload.Title, err)
		return nil, err
	}
	if err := validateModerationField("文章摘要", payload.Summary); err != nil {
		createModerationHit(s.db, authorID, "article_create", "文章摘要", payload.Summary, err)
		return nil, err
	}
	if err := validateModerationField("文章正文", payload.Content); err != nil {
		createModerationHit(s.db, authorID, "article_create", "文章正文", payload.Content, err)
		return nil, err
	}
	if err := validateModerationList("文章标签", payload.Tags); err != nil {
		createModerationHit(s.db, authorID, "article_create", "文章标签", strings.Join(payload.Tags, ", "), err)
		return nil, err
	}

	article := model.Article{
		Title:      payload.Title,
		Summary:    payload.Summary,
		Content:    payload.Content,
		CoverImage: payload.CoverImage,
		Status:     model.ArticleDraft,
		AuthorID:   authorID,
		CategoryID: payload.CategoryID,
	}

	if err := s.db.Create(&article).Error; err != nil {
		return nil, err
	}
	if err := s.syncTags(&article, payload.Tags); err != nil {
		return nil, err
	}

	return s.GetByID(article.ID, authorID)
}

func (s *ArticleService) Update(articleID, userID uint, role string, payload ArticlePayload) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if role != model.RoleAdmin && article.AuthorID != userID {
		return nil, ErrArticleNoPermission
	}
	if err := validateModerationField("文章标题", payload.Title); err != nil {
		createModerationHit(s.db, userID, "article_update", "文章标题", payload.Title, err)
		return nil, err
	}
	if err := validateModerationField("文章摘要", payload.Summary); err != nil {
		createModerationHit(s.db, userID, "article_update", "文章摘要", payload.Summary, err)
		return nil, err
	}
	if err := validateModerationField("文章正文", payload.Content); err != nil {
		createModerationHit(s.db, userID, "article_update", "文章正文", payload.Content, err)
		return nil, err
	}
	if err := validateModerationList("文章标签", payload.Tags); err != nil {
		createModerationHit(s.db, userID, "article_update", "文章标签", strings.Join(payload.Tags, ", "), err)
		return nil, err
	}

	updates := map[string]interface{}{
		"title":         payload.Title,
		"summary":       payload.Summary,
		"content":       payload.Content,
		"cover_image":   payload.CoverImage,
		"category_id":   payload.CategoryID,
		"reject_reason": "",
	}
	if article.Status == model.ArticleRejected {
		updates["status"] = model.ArticleDraft
	}

	if err := s.db.Model(&article).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.syncTags(&article, payload.Tags); err != nil {
		return nil, err
	}

	return s.GetByID(article.ID, userID)
}

func (s *ArticleService) Submit(articleID, userID uint, role string) (*model.Article, error) {
	var author model.User
	if err := s.db.First(&author, userID).Error; err != nil {
		return nil, err
	}
	if author.Status == model.UserBanned {
		return nil, ErrUserBannedSubmit
	}

	var article model.Article
	if err := s.db.Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if role != model.RoleAdmin && article.AuthorID != userID {
		return nil, ErrArticleNoPermission
	}
	if err := validateModerationField("文章标题", article.Title); err != nil {
		createModerationHit(s.db, userID, "article_submit", "文章标题", article.Title, err)
		reason := "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
		_ = s.db.Model(&article).Updates(map[string]interface{}{"status": model.ArticleRejected, "reject_reason": reason}).Error
		return s.GetByID(articleID, userID)
	}
	if err := validateModerationField("文章摘要", article.Summary); err != nil {
		createModerationHit(s.db, userID, "article_submit", "文章摘要", article.Summary, err)
		reason := "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
		_ = s.db.Model(&article).Updates(map[string]interface{}{"status": model.ArticleRejected, "reject_reason": reason}).Error
		return s.GetByID(articleID, userID)
	}
	if err := validateModerationField("文章正文", article.Content); err != nil {
		createModerationHit(s.db, userID, "article_submit", "文章正文", article.Content, err)
		reason := "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
		_ = s.db.Model(&article).Updates(map[string]interface{}{"status": model.ArticleRejected, "reject_reason": reason}).Error
		return s.GetByID(articleID, userID)
	}
	tagNames := make([]string, 0, len(article.Tags))
	for _, tag := range article.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	if err := validateModerationList("文章标签", tagNames); err != nil {
		createModerationHit(s.db, userID, "article_submit", "文章标签", strings.Join(tagNames, ", "), err)
		reason := "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
		_ = s.db.Model(&article).Updates(map[string]interface{}{"status": model.ArticleRejected, "reject_reason": reason}).Error
		return s.GetByID(articleID, userID)
	}
	if role == model.RoleAdmin {
		now := time.Now()
		if err := s.db.Model(&article).Updates(map[string]interface{}{
			"status":        model.ArticlePublished,
			"published_at":  &now,
			"reject_reason": "",
		}).Error; err != nil {
			return nil, err
		}
		return s.GetByID(articleID, userID)
	}

	if err := s.db.Model(&article).Updates(map[string]interface{}{
		"status":        model.ArticlePending,
		"reject_reason": "",
	}).Error; err != nil {
		return nil, err
	}
	return s.GetByID(articleID, userID)
}

func (s *ArticleService) ListPublished(filter PublishedArticleFilter) ([]model.Article, Pagination, error) {
	var articles []model.Article
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)
	query := s.baseArticleQuery().
		Where("articles.status = ?", model.ArticlePublished)
	countQuery := s.db.Model(&model.Article{}).Where("articles.status = ?", model.ArticlePublished)
	if filter.Keyword != "" {
		query = query.Where("articles.title LIKE ? OR articles.summary LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		countQuery = countQuery.Where("articles.title LIKE ? OR articles.summary LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.CategoryID != nil {
		query = query.Where("articles.category_id = ?", *filter.CategoryID)
		countQuery = countQuery.Where("articles.category_id = ?", *filter.CategoryID)
	}
	if filter.Tag != "" {
		query = query.
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", filter.Tag).
			Select("articles.*").
			Distinct()
		countQuery = countQuery.
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", filter.Tag).
			Distinct("articles.id")
	}
	if filter.AuthorID > 0 {
		query = query.Where("articles.author_id = ?", filter.AuthorID)
		countQuery = countQuery.Where("articles.author_id = ?", filter.AuthorID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	orderClause := "articles.published_at desc, articles.created_at desc"
	switch filter.Sort {
	case "oldest":
		orderClause = "articles.published_at asc, articles.created_at asc"
	case "popular":
		orderClause = "articles.view_count desc, articles.published_at desc, articles.created_at desc"
	case "latest":
		orderClause = "articles.published_at desc, articles.created_at desc"
	}
	err := query.
		Order(orderClause).
		Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).
		Find(&articles).Error
	if err == nil {
		for index := range articles {
			s.fillReactionSummary(&articles[index], 0)
		}
	}
	pagination.Total = total
	return articles, pagination, err
}

func (s *ArticleService) ListTrending(limit int) ([]model.Article, error) {
	if limit <= 0 {
		limit = 6
	}
	var articles []model.Article
	err := s.baseArticleQuery().
		Where("status = ?", model.ArticlePublished).
		Order("published_at desc, created_at desc").
		Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	for index := range articles {
		s.fillReactionSummary(&articles[index], 0)
	}

	for i := 0; i < len(articles); i++ {
		for j := i + 1; j < len(articles); j++ {
			leftScore := articles[i].LikesCount*2 + articles[i].FavoritesCount
			rightScore := articles[j].LikesCount*2 + articles[j].FavoritesCount
			if rightScore > leftScore {
				articles[i], articles[j] = articles[j], articles[i]
			}
		}
	}

	return articles, nil
}

func (s *ArticleService) ListMine(userID uint, page, pageSize int) ([]model.Article, Pagination, error) {
	var articles []model.Article
	var total int64
	pagination := normalizePagination(page, pageSize)
	query := s.baseArticleQuery().Where("author_id = ?", userID).Order("updated_at desc")
	if err := query.Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&articles).Error
	if err == nil {
		for index := range articles {
			s.fillReactionSummary(&articles[index], userID)
		}
	}
	pagination.Total = total
	return articles, pagination, err
}

func (s *ArticleService) ListLiked(userID uint) ([]model.Article, error) {
	var articles []model.Article
	err := s.baseArticleQuery().
		Joins("JOIN article_likes ON article_likes.article_id = articles.id").
		Where("article_likes.user_id = ? AND articles.status = ?", userID, model.ArticlePublished).
		Order("article_likes.created_at desc").
		Find(&articles).Error
	if err == nil {
		for index := range articles {
			s.fillReactionSummary(&articles[index], userID)
		}
	}
	return articles, err
}

func (s *ArticleService) ListFavorited(userID uint) ([]model.Article, error) {
	var articles []model.Article
	err := s.baseArticleQuery().
		Joins("JOIN article_favorites ON article_favorites.article_id = articles.id").
		Where("article_favorites.user_id = ? AND articles.status = ?", userID, model.ArticlePublished).
		Order("article_favorites.created_at desc").
		Find(&articles).Error
	if err == nil {
		for index := range articles {
			s.fillReactionSummary(&articles[index], userID)
		}
	}
	return articles, err
}

func (s *ArticleService) ListPending() ([]model.Article, error) {
	var articles []model.Article
	err := s.baseArticleQuery().Where("status = ?", model.ArticlePending).Order("updated_at asc").Find(&articles).Error
	if err == nil {
		for index := range articles {
			s.fillReactionSummary(&articles[index], 0)
		}
	}
	return articles, err
}

func (s *ArticleService) GetByID(articleID uint, viewerID uint) (*model.Article, error) {
	var article model.Article
	if err := s.baseArticleQuery().Preload("Comments.User").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if article.Status == model.ArticlePublished {
		_ = s.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
		article.ViewCount++
	}
	s.fillReactionSummary(&article, viewerID)
	return &article, nil
}

func (s *ArticleService) Review(articleID, reviewerID uint, payload ReviewPayload) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Tags").First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if article.Status != model.ArticlePending {
		return nil, ErrNotPending
	}

	updates := map[string]interface{}{}
	review := model.ArticleReview{
		ArticleID:  articleID,
		ReviewerID: reviewerID,
		Action:     payload.Action,
		Reason:     payload.Reason,
	}

	switch payload.Action {
	case "approve":
		if err := validateModerationField("文章标题", article.Title); err != nil {
			createModerationHit(s.db, article.AuthorID, "article_review", "文章标题", article.Title, err)
			payload.Action = "reject"
			payload.Reason = "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
			review.Action = payload.Action
			review.Reason = payload.Reason
			updates["status"] = model.ArticleRejected
			updates["reject_reason"] = payload.Reason
			break
		}
		if err := validateModerationField("文章摘要", article.Summary); err != nil {
			createModerationHit(s.db, article.AuthorID, "article_review", "文章摘要", article.Summary, err)
			payload.Action = "reject"
			payload.Reason = "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
			review.Action = payload.Action
			review.Reason = payload.Reason
			updates["status"] = model.ArticleRejected
			updates["reject_reason"] = payload.Reason
			break
		}
		if err := validateModerationField("文章正文", article.Content); err != nil {
			createModerationHit(s.db, article.AuthorID, "article_review", "文章正文", article.Content, err)
			payload.Action = "reject"
			payload.Reason = "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
			review.Action = payload.Action
			review.Reason = payload.Reason
			updates["status"] = model.ArticleRejected
			updates["reject_reason"] = payload.Reason
			break
		}
		tagNames := make([]string, 0, len(article.Tags))
		for _, tag := range article.Tags {
			tagNames = append(tagNames, tag.Name)
		}
		if err := validateModerationList("文章标签", tagNames); err != nil {
			createModerationHit(s.db, article.AuthorID, "article_review", "文章标签", strings.Join(tagNames, ", "), err)
			payload.Action = "reject"
			payload.Reason = "内容命中敏感词，系统已自动驳回：" + err.(*ModerationError).Word
			review.Action = payload.Action
			review.Reason = payload.Reason
			updates["status"] = model.ArticleRejected
			updates["reject_reason"] = payload.Reason
			break
		}
		now := time.Now()
		updates["status"] = model.ArticlePublished
		updates["published_at"] = &now
		updates["reject_reason"] = ""
	case "reject":
		updates["status"] = model.ArticleRejected
		updates["reject_reason"] = payload.Reason
	default:
		return nil, ErrInvalidAction
	}

	return s.transactionReview(article, review, updates)
}

func (s *ArticleService) Dashboard() (map[string]int64, error) {
	result := map[string]int64{}
	statuses := []string{model.ArticleDraft, model.ArticlePending, model.ArticlePublished, model.ArticleRejected}
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
	return result, nil
}

func (s *ArticleService) transactionReview(article model.Article, review model.ArticleReview, updates map[string]interface{}) (*model.Article, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&article).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Create(&review).Error
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(article.ID, 0)
}

func (s *ArticleService) baseArticleQuery() *gorm.DB {
	return s.db.Preload("Author").Preload("Category").Preload("Tags")
}

func (s *ArticleService) syncTags(article *model.Article, names []string) error {
	unique := make(map[string]struct{})
	tags := make([]model.Tag, 0)

	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, exists := unique[name]; exists {
			continue
		}
		unique[name] = struct{}{}

		var tag model.Tag
		err := s.db.Where("name = ?", name).First(&tag).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = model.Tag{Name: name}
			if err := s.db.Create(&tag).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	return s.db.Model(article).Association("Tags").Replace(tags)
}

func (s *ArticleService) ToggleLike(articleID, userID uint) (*ReactionSummary, error) {
	return s.toggleReaction(articleID, userID, "like")
}

func (s *ArticleService) ToggleFavorite(articleID, userID uint) (*ReactionSummary, error) {
	return s.toggleReaction(articleID, userID, "favorite")
}

func (s *ArticleService) toggleReaction(articleID, userID uint, kind string) (*ReactionSummary, error) {
	var article model.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if article.Status != model.ArticlePublished {
		return nil, ErrNotPublishedArticle
	}

	if kind == "like" {
		var like model.ArticleLike
		err := s.db.Where("article_id = ? AND user_id = ?", articleID, userID).First(&like).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&model.ArticleLike{ArticleID: articleID, UserID: userID}).Error; err != nil {
				return nil, err
			}
		} else if err == nil {
			if err := s.db.Delete(&like).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if kind == "favorite" {
		var favorite model.ArticleFavorite
		err := s.db.Where("article_id = ? AND user_id = ?", articleID, userID).First(&favorite).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&model.ArticleFavorite{ArticleID: articleID, UserID: userID}).Error; err != nil {
				return nil, err
			}
		} else if err == nil {
			if err := s.db.Delete(&favorite).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	summary := &ReactionSummary{}
	if err := s.db.Model(&model.ArticleLike{}).Where("article_id = ?", articleID).Count(&summary.LikesCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.ArticleFavorite{}).Where("article_id = ?", articleID).Count(&summary.FavoritesCount).Error; err != nil {
		return nil, err
	}
	var count int64
	_ = s.db.Model(&model.ArticleLike{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	summary.IsLiked = count > 0
	count = 0
	_ = s.db.Model(&model.ArticleFavorite{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	summary.IsFavorited = count > 0
	return summary, nil
}

func (s *ArticleService) fillReactionSummary(article *model.Article, viewerID uint) {
	_ = s.db.Model(&model.ArticleLike{}).Where("article_id = ?", article.ID).Count(&article.LikesCount).Error
	_ = s.db.Model(&model.ArticleFavorite{}).Where("article_id = ?", article.ID).Count(&article.FavoritesCount).Error
	if viewerID == 0 {
		return
	}
	var count int64
	_ = s.db.Model(&model.ArticleLike{}).Where("article_id = ? AND user_id = ?", article.ID, viewerID).Count(&count).Error
	article.IsLiked = count > 0
	count = 0
	_ = s.db.Model(&model.ArticleFavorite{}).Where("article_id = ? AND user_id = ?", article.ID, viewerID).Count(&count).Error
	article.IsFavorited = count > 0
}

func normalizePagination(page, pageSize int) Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 9
	}
	if pageSize > 30 {
		pageSize = 30
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func ParsePagination(pageRaw, sizeRaw string) (int, int) {
	page, _ := strconv.Atoi(pageRaw)
	pageSize, _ := strconv.Atoi(sizeRaw)
	return page, pageSize
}
