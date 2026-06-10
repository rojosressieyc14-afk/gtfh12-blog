package service

import (
	"encoding/json"
	"strings"

	"blog/server/internal/config"
	"blog/server/internal/model"
	"blog/server/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg config.Config
}

type UpdateProfileInput struct {
	Avatar      string
	Headline    string
	CurrentRole string
	YearsLabel  string
	Motto       string
	Location    string
	Email       string
	ResumeURL   string
	WebsiteURL  string
	GithubURL   string
	GiteeURL    string
	JuejinURL   string
	CSDNURL     string
	Skills      []string
	FocusAreas  []string
	Bio         string
}

func NewAuthService(db *gorm.DB, cfg config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

func (s *AuthService) Register(username, password string) (*model.User, string, error) {
	username = strings.TrimSpace(username)
	if err := validateModerationField("用户名", username); err != nil {
		createModerationHit(s.db, 0, "register", "用户名", username, err)
		return nil, "", err
	}

	var exists int64
	if err := s.db.Model(&model.User{}).Where("username = ?", username).Count(&exists).Error; err != nil {
		return nil, "", err
	}
	if exists > 0 {
		return nil, "", ErrUsernameExists
	}

	if err := utils.ValidatePassword(password, utils.DefaultPasswordPolicy); err != nil {
		return nil, "", err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := model.User{
		Username: username,
		Password: hash,
		Role:     model.RoleUser,
		Status:   model.UserActive,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	return &user, token, err
}

func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", ErrInvalidCredentials
	}
	if user.Status == model.UserBanned {
		return nil, "", ErrAccountBanned
	}
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	return &user, token, err
}

func (s *AuthService) Me(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateProfile(userID uint, input UpdateProfileInput) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	if err := validateModerationField("标题文案", input.Headline); err != nil {
		createModerationHit(s.db, userID, "profile", "标题文案", input.Headline, err)
		return nil, err
	}
	if err := validateModerationField("当前身份", input.CurrentRole); err != nil {
		createModerationHit(s.db, userID, "profile", "当前身份", input.CurrentRole, err)
		return nil, err
	}
	if err := validateModerationField("经验标签", input.YearsLabel); err != nil {
		createModerationHit(s.db, userID, "profile", "经验标签", input.YearsLabel, err)
		return nil, err
	}
	if err := validateModerationField("个人签名", input.Motto); err != nil {
		createModerationHit(s.db, userID, "profile", "个人签名", input.Motto, err)
		return nil, err
	}
	if err := validateModerationField("所在地", input.Location); err != nil {
		createModerationHit(s.db, userID, "profile", "所在地", input.Location, err)
		return nil, err
	}
	if err := validateModerationField("个人简介", input.Bio); err != nil {
		createModerationHit(s.db, userID, "profile", "个人简介", input.Bio, err)
		return nil, err
	}
	if err := validateModerationList("技能栏", input.Skills); err != nil {
		createModerationHit(s.db, userID, "profile", "技能栏", strings.Join(input.Skills, ", "), err)
		return nil, err
	}
	if err := validateModerationList("关注方向", input.FocusAreas); err != nil {
		createModerationHit(s.db, userID, "profile", "关注方向", strings.Join(input.FocusAreas, ", "), err)
		return nil, err
	}

	skillsJSON, _ := json.Marshal(normalizeSkills(input.Skills))
	focusJSON, _ := json.Marshal(normalizeSkills(input.FocusAreas))
	updates := map[string]any{
		"avatar":       strings.TrimSpace(input.Avatar),
		"headline":     strings.TrimSpace(input.Headline),
		"current_role": strings.TrimSpace(input.CurrentRole),
		"years_label":  strings.TrimSpace(input.YearsLabel),
		"motto":        strings.TrimSpace(input.Motto),
		"location":     strings.TrimSpace(input.Location),
		"email":        strings.TrimSpace(input.Email),
		"resume_url":   strings.TrimSpace(input.ResumeURL),
		"website_url":  strings.TrimSpace(input.WebsiteURL),
		"github_url":   strings.TrimSpace(input.GithubURL),
		"gitee_url":    strings.TrimSpace(input.GiteeURL),
		"juejin_url":   strings.TrimSpace(input.JuejinURL),
		"csdn_url":     strings.TrimSpace(input.CSDNURL),
		"skills":       string(skillsJSON),
		"focus_areas":  string(focusJSON),
		"bio":          strings.TrimSpace(input.Bio),
	}
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func normalizeSkills(skills []string) []string {
	if len(skills) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(skills))
	seen := map[string]struct{}{}
	for _, skill := range skills {
		item := strings.TrimSpace(skill)
		if item == "" {
			continue
		}
		key := strings.ToLower(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

func (s *AuthService) AuthorProfile(userID uint) (*model.User, []model.Article, []model.Project, map[string]int64, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var articles []model.Article
	if err := s.db.Preload("Category").Preload("Tags").Where("author_id = ? AND status = ?", userID, model.ArticlePublished).Order("published_at desc, created_at desc").Find(&articles).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var projects []model.Project
	if err := s.db.Where("author_id = ? AND status = ?", userID, model.ProjectPublished).Order("is_featured desc, sort_order desc, published_at desc, created_at desc").Find(&projects).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	stats := map[string]int64{}
	var count int64
	if err := s.db.Model(&model.Article{}).Where("author_id = ? AND status = ?", userID, model.ArticlePublished).Count(&count).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	stats["articles"] = count
	if err := s.db.Model(&model.Project{}).Where("author_id = ? AND status = ?", userID, model.ProjectPublished).Count(&count).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	stats["projects"] = count
	if err := s.db.Table("comments").Joins("JOIN articles ON articles.id = comments.article_id").Where("articles.author_id = ?", userID).Count(&count).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	stats["comments"] = count
	if err := s.db.Table("article_likes").Joins("JOIN articles ON articles.id = article_likes.article_id").Where("articles.author_id = ?", userID).Count(&count).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	stats["likes"] = count

	return &user, articles, projects, stats, nil
}

func (s *AuthService) GetRecommendedAuthors(limit int) ([]model.User, error) {
	if limit <= 0 {
		limit = 6
	}
	var authorIDs []uint
	if err := s.db.Table("articles").
		Select("author_id").
		Where("status = ?", model.ArticlePublished).
		Group("author_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("author_id", &authorIDs).Error; err != nil {
		return nil, err
	}
	if len(authorIDs) < limit {
		var projectIDs []uint
		if err := s.db.Table("projects").
			Select("author_id").
			Where("status = ? AND author_id NOT IN ?", model.ProjectPublished, authorIDs).
			Group("author_id").
			Order("COUNT(*) DESC").
			Limit(limit-len(authorIDs)).
			Pluck("author_id", &projectIDs).Error; err != nil {
			return nil, err
		}
		authorIDs = append(authorIDs, projectIDs...)
	}
	if len(authorIDs) == 0 {
		return nil, nil
	}
	var users []model.User
	if err := s.db.Where("id IN ? AND status = ?", authorIDs, "active").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
