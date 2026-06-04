package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

type AIReviewContent struct {
	ID         uint     `json:"id"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	AuthorName string   `json:"authorName"`
}

type AIReviewResultInput struct {
	RiskLevel          string   `json:"riskLevel"`
	RiskLabels         []string `json:"riskLabels"`
	Summary            string   `json:"summary"`
	SuspiciousSegments []string `json:"suspiciousSegments"`
	Suggestion         string   `json:"suggestion"`
	ModelName          string   `json:"modelName"`
}

type AIReviewService struct {
	db *gorm.DB
}

func NewAIReviewService(db *gorm.DB) *AIReviewService {
	return &AIReviewService{db: db}
}

func (s *AIReviewService) GetContentForReview(targetType string, targetID uint) (*AIReviewContent, error) {
	switch targetType {
	case model.TargetTypeArticle:
		return s.getArticleContent(targetID)
	case model.TargetTypeProject:
		return s.getProjectContent(targetID)
	default:
		return nil, ErrUnsupportedTarget
	}
}

func (s *AIReviewService) getArticleContent(id uint) (*AIReviewContent, error) {
	var article model.Article
	if err := s.db.Preload("Author").Preload("Tags").First(&article, id).Error; err != nil {
		return nil, err
	}
	tags := make([]string, len(article.Tags))
	for i, tag := range article.Tags {
		tags[i] = tag.Name
	}
	authorName := ""
	if article.Author.Username != "" {
		authorName = article.Author.Username
	}
	return &AIReviewContent{
		ID:         article.ID,
		Title:      article.Title,
		Summary:    article.Summary,
		Content:    article.Content,
		Tags:       tags,
		AuthorName: authorName,
	}, nil
}

func (s *AIReviewService) getProjectContent(id uint) (*AIReviewContent, error) {
	var project model.Project
	if err := s.db.Preload("Author").First(&project, id).Error; err != nil {
		return nil, err
	}
	authorName := ""
	if project.Author.Username != "" {
		authorName = project.Author.Username
	}
	return &AIReviewContent{
		ID:         project.ID,
		Title:      project.Title,
		Summary:    project.Summary,
		Content:    project.Content,
		Tags:       project.TechStacks,
		AuthorName: authorName,
	}, nil
}

func (s *AIReviewService) SaveResult(operatorID uint, targetType string, targetID uint, input AIReviewResultInput) (*model.AIReviewRecord, error) {
	if targetType != model.TargetTypeArticle && targetType != model.TargetTypeProject {
		return nil, ErrUnsupportedTarget
	}
	if input.RiskLevel != model.RiskLevelLow && input.RiskLevel != model.RiskLevelMid && input.RiskLevel != model.RiskLevelHigh {
		input.RiskLevel = model.RiskLevelLow
	}

	var existing model.AIReviewRecord
	err := s.db.Where("target_type = ? AND target_id = ?", targetType, targetID).First(&existing).Error

	if err == nil {
		updates := map[string]any{
			"operator_id":         operatorID,
			"risk_level":          input.RiskLevel,
			"risk_labels":         normalizeStringList(input.RiskLabels),
			"summary":             strings.TrimSpace(input.Summary),
			"suspicious_segments": normalizeStringList(input.SuspiciousSegments),
			"suggestion":          strings.TrimSpace(input.Suggestion),
			"model_name":          strings.TrimSpace(input.ModelName),
			"updated_at":          time.Now(),
		}
		if err := s.db.Model(&existing).Updates(updates).Error; err != nil {
			return nil, err
		}
		return s.GetRecord(targetType, targetID)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	record := model.AIReviewRecord{
		TargetType:         targetType,
		TargetID:           targetID,
		OperatorID:         operatorID,
		RiskLevel:          input.RiskLevel,
		RiskLabels:         normalizeStringList(input.RiskLabels),
		Summary:            strings.TrimSpace(input.Summary),
		SuspiciousSegments: normalizeStringList(input.SuspiciousSegments),
		Suggestion:         strings.TrimSpace(input.Suggestion),
		ModelName:          strings.TrimSpace(input.ModelName),
	}
	if record.ModelName == "" {
		record.ModelName = "manual-review"
	}

	if err := s.db.Create(&record).Error; err != nil {
		return nil, err
	}
	return s.GetRecord(targetType, targetID)
}

func (s *AIReviewService) GetRecord(targetType string, targetID uint) (*model.AIReviewRecord, error) {
	var record model.AIReviewRecord
	if err := s.db.Preload("Operator").Where("target_type = ? AND target_id = ?", targetType, targetID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func FormatContentForPrompt(content *AIReviewContent, targetType string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s ===\n", targetTypeLabel(targetType)))
	if content.AuthorName != "" {
		sb.WriteString(fmt.Sprintf("作者: %s\n", content.AuthorName))
	}
	sb.WriteString(fmt.Sprintf("标题: %s\n", content.Title))
	if content.Summary != "" {
		sb.WriteString(fmt.Sprintf("摘要: %s\n", content.Summary))
	}
	if len(content.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("标签: %s\n", strings.Join(content.Tags, ", ")))
	}
	sb.WriteString(fmt.Sprintf("\n正文:\n%s\n", content.Content))
	sb.WriteString("\n---\n请对以上内容进行 AI 审核分析，返回风险等级、风险标签、可疑片段和审核建议。")
	return sb.String()
}

func targetTypeLabel(t string) string {
	switch t {
	case model.TargetTypeArticle:
		return "文章"
	case model.TargetTypeProject:
		return "项目"
	default:
		return "内容"
	}
}
