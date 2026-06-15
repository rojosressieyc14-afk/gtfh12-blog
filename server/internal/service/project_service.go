package service

import (
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

type ProjectPayload struct {
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	RoleLabel  string   `json:"roleLabel"`
	Duration   string   `json:"duration"`
	TeamLabel  string   `json:"teamLabel"`
	Content    string   `json:"content"`
	CoverImage string   `json:"coverImage"`
	TechStacks []string `json:"techStacks"`
	Highlights []string `json:"highlights"`
	Process    []string `json:"process"`
	Challenges []string `json:"challenges"`
	Solutions  []string `json:"solutions"`
	Results    []string `json:"results"`
	DemoURL    string   `json:"demoUrl"`
	RepoURL    string   `json:"repoUrl"`
	IsPrivate  bool     `json:"isPrivate"`
	IsFeatured bool     `json:"isFeatured"`
	SortOrder  int      `json:"sortOrder"`
}

type ProjectFilter struct {
	FeaturedOnly bool
	Keyword      string
	Stack        string
	AuthorID     uint
	Sort         string
	Page         int
	PageSize     int
}

type ProjectReviewPayload struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) ListPublished(filter ProjectFilter) ([]model.Project, Pagination, error) {
	var items []model.Project
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.Project{}).Preload("Author").Where("status = ? AND is_private = ?", model.ProjectPublished, false)
	countQuery := s.db.Model(&model.Project{}).Where("status = ? AND is_private = ?", model.ProjectPublished, false)

	if filter.FeaturedOnly {
		query = query.Where("is_featured = ?", true)
		countQuery = countQuery.Where("is_featured = ?", true)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", like, like, like)
		countQuery = countQuery.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", like, like, like)
	}
	if filter.Stack != "" {
		like := "%" + filter.Stack + "%"
		query = query.Where("tech_stacks LIKE ?", like)
		countQuery = countQuery.Where("tech_stacks LIKE ?", like)
	}
	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
		countQuery = countQuery.Where("author_id = ?", filter.AuthorID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}

	orderClause := "is_featured desc, sort_order desc, published_at desc, created_at desc"
	switch filter.Sort {
	case "latest":
		orderClause = "published_at desc, created_at desc"
	case "oldest":
		orderClause = "published_at asc, created_at asc"
	}

	err := query.
		Order(orderClause).
		Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).
		Find(&items).Error
	pagination.Total = total
	return items, pagination, err
}

func (s *ProjectService) ListMine(userID uint, page, pageSize int) ([]model.Project, Pagination, error) {
	var items []model.Project
	var total int64
	pagination := normalizePagination(page, pageSize)

	query := s.db.Model(&model.Project{}).Preload("Author").Where("author_id = ?", userID).Order("updated_at desc")
	if err := query.Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&items).Error
	pagination.Total = total
	return items, pagination, err
}

func (s *ProjectService) ListPending() ([]model.Project, error) {
	var items []model.Project
	err := s.db.Preload("Author").
		Where("status = ?", model.ProjectPending).
		Order("updated_at asc").
		Find(&items).Error
	return items, err
}

func (s *ProjectService) GetByID(projectID, viewerID uint, role string) (*model.Project, error) {
	var item model.Project
	if err := s.db.Preload("Author").First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if item.IsPrivate && role != model.RoleAdmin && item.AuthorID != viewerID {
		return nil, ErrProjectNotPublished
	}
	if item.Status != model.ProjectPublished && role != model.RoleAdmin && item.AuthorID != viewerID {
		return nil, ErrProjectNotPublished
	}
	return &item, nil
}

func (s *ProjectService) Create(authorID uint, role string, payload ProjectPayload) (*model.Project, error) {
	var author model.User
	if err := s.db.First(&author, authorID).Error; err != nil {
		return nil, err
	}
	if author.Status == model.UserBanned {
		return nil, ErrUserBannedProject
	}
	item, err := s.validatePayload(payload, authorID, "project_create")
	if err != nil {
		return nil, err
	}

	item.Status = model.ProjectDraft
	item.IsPrivate = payload.IsPrivate
	item.IsFeatured = payload.IsFeatured && role == model.RoleAdmin
	item.SortOrder = payload.SortOrder
	item.AuthorID = authorID

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return s.GetByID(item.ID, authorID, role)
}

func (s *ProjectService) Update(projectID, userID uint, role string, payload ProjectPayload) (*model.Project, error) {
	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if err := s.ensureAuthorCanWrite(item, userID, role); err != nil {
		return nil, err
	}
	if _, err := s.validatePayload(payload, userID, "project_update"); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"title":         strings.TrimSpace(payload.Title),
		"summary":       strings.TrimSpace(payload.Summary),
		"role_label":    strings.TrimSpace(payload.RoleLabel),
		"duration":      strings.TrimSpace(payload.Duration),
		"team_label":    strings.TrimSpace(payload.TeamLabel),
		"content":       strings.TrimSpace(payload.Content),
		"cover_image":   strings.TrimSpace(payload.CoverImage),
		"tech_stacks":   normalizeStringList(payload.TechStacks),
		"highlights":    normalizeStringList(payload.Highlights),
		"process":       normalizeStringList(payload.Process),
		"challenges":    normalizeStringList(payload.Challenges),
		"solutions":     normalizeStringList(payload.Solutions),
		"results":       normalizeStringList(payload.Results),
		"demo_url":      strings.TrimSpace(payload.DemoURL),
		"repo_url":      strings.TrimSpace(payload.RepoURL),
		"is_private":    payload.IsPrivate,
		"reject_reason": "",
	}

	if role == model.RoleAdmin {
		updates["is_featured"] = payload.IsFeatured
		updates["sort_order"] = payload.SortOrder
	}

	if item.Status == model.ProjectPublished || item.Status == model.ProjectRejected {
		updates["status"] = model.ProjectDraft
		updates["published_at"] = nil
	}

	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(projectID, userID, role)
}

func (s *ProjectService) Submit(projectID, userID uint, role string) (*model.Project, error) {
	var author model.User
	if err := s.db.First(&author, userID).Error; err != nil {
		return nil, err
	}
	if author.Status == model.UserBanned {
		return nil, ErrUserBannedProject
	}

	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if err := s.ensureAuthorCanWrite(item, userID, role); err != nil {
		return nil, err
	}
	if err := s.validateProjectModel(item, item.AuthorID, "project_submit"); err != nil {
		reason := "内容命中敏感词，系统已自动驳回：" + moderationWord(err)
		if updateErr := s.db.Model(&item).Updates(map[string]any{
			"status":        model.ProjectRejected,
			"published_at":  nil,
			"reject_reason": reason,
		}).Error; updateErr != nil {
			return nil, updateErr
		}
		return s.GetByID(projectID, userID, role)
	}

	updates := map[string]any{
		"reject_reason": "",
	}
	now := time.Now()
	if role == model.RoleAdmin || item.IsPrivate {
		updates["status"] = model.ProjectPublished
		updates["published_at"] = &now
	} else {
		updates["status"] = model.ProjectPending
		updates["published_at"] = nil
	}

	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(projectID, userID, role)
}

func (s *ProjectService) Review(projectID, reviewerID uint, payload ProjectReviewPayload) (*model.Project, error) {
	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return nil, err
	}
	if item.Status != model.ProjectPending {
		return nil, ErrProjectNotPending
	}

	switch payload.Action {
	case "approve":
		if err := s.validateProjectModel(item, item.AuthorID, "project_review"); err != nil {
			reason := "内容命中敏感词，系统已自动驳回：" + moderationWord(err)
			if updateErr := s.db.Model(&item).Updates(map[string]any{
				"status":        model.ProjectRejected,
				"published_at":  nil,
				"reject_reason": reason,
			}).Error; updateErr != nil {
				return nil, updateErr
			}
			return s.GetByID(projectID, reviewerID, model.RoleAdmin)
		}

		now := time.Now()
		if err := s.db.Model(&item).Updates(map[string]any{
			"status":        model.ProjectPublished,
			"published_at":  &now,
			"reject_reason": "",
		}).Error; err != nil {
			return nil, err
		}
	case "reject":
		reason := strings.TrimSpace(payload.Reason)
		if reason == "" {
			return nil, ErrRejectReasonRequired
		}
		if err := s.db.Model(&item).Updates(map[string]any{
			"status":        model.ProjectRejected,
			"published_at":  nil,
			"reject_reason": reason,
		}).Error; err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidAction
	}

	return s.GetByID(projectID, reviewerID, model.RoleAdmin)
}

func (s *ProjectService) Delete(projectID, userID uint, role string) error {
	var item model.Project
	if err := s.db.First(&item, projectID).Error; err != nil {
		return err
	}
	if role != model.RoleAdmin && item.AuthorID != userID {
		return ErrProjectNoPermission
	}
	return s.db.Delete(&item).Error
}

func (s *ProjectService) ensureAuthorCanWrite(item model.Project, userID uint, role string) error {
	if role == model.RoleAdmin {
		return nil
	}
	if item.AuthorID != userID {
		return ErrProjectNoPermission
	}

	var author model.User
	if err := s.db.Select("id", "status").First(&author, userID).Error; err != nil {
		return err
	}
	if author.Status == model.UserBanned {
		return ErrUserBannedProject
	}
	return nil
}

func (s *ProjectService) validatePayload(payload ProjectPayload, userID uint, scene string) (model.Project, error) {
	item := model.Project{
		Title:      strings.TrimSpace(payload.Title),
		Summary:    strings.TrimSpace(payload.Summary),
		RoleLabel:  strings.TrimSpace(payload.RoleLabel),
		Duration:   strings.TrimSpace(payload.Duration),
		TeamLabel:  strings.TrimSpace(payload.TeamLabel),
		Content:    strings.TrimSpace(payload.Content),
		CoverImage: strings.TrimSpace(payload.CoverImage),
		TechStacks: normalizeStringList(payload.TechStacks),
		Highlights: normalizeStringList(payload.Highlights),
		Process:    normalizeStringList(payload.Process),
		Challenges: normalizeStringList(payload.Challenges),
		Solutions:  normalizeStringList(payload.Solutions),
		Results:    normalizeStringList(payload.Results),
		DemoURL:    strings.TrimSpace(payload.DemoURL),
		RepoURL:    strings.TrimSpace(payload.RepoURL),
	}

	if item.Title == "" || item.Content == "" {
		return model.Project{}, ErrProjectContentEmpty
	}
	if err := s.validateProjectModel(item, userID, scene); err != nil {
		return model.Project{}, err
	}
	return item, nil
}

func (s *ProjectService) validateProjectModel(item model.Project, userID uint, scene string) error {
	fields := []struct {
		label string
		value string
	}{
		{label: "项目标题", value: item.Title},
		{label: "项目摘要", value: item.Summary},
		{label: "项目角色", value: item.RoleLabel},
		{label: "项目周期", value: item.Duration},
		{label: "团队信息", value: item.TeamLabel},
		{label: "项目正文", value: item.Content},
	}
	for _, field := range fields {
		if err := validateModerationField(field.label, field.value); err != nil {
			createModerationHit(s.db, userID, scene, field.label, field.value, err)
			return err
		}
	}

	lists := []struct {
		label  string
		values []string
	}{
		{label: "技术栈", values: item.TechStacks},
		{label: "项目亮点", values: item.Highlights},
		{label: "推进过程", values: item.Process},
		{label: "遇到的挑战", values: item.Challenges},
		{label: "解决方案", values: item.Solutions},
		{label: "最终结果", values: item.Results},
	}
	for _, field := range lists {
		if err := validateModerationList(field.label, field.values); err != nil {
			createModerationHit(s.db, userID, scene, field.label, strings.Join(field.values, ", "), err)
			return err
		}
	}

	return nil
}

func normalizeStringList(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
