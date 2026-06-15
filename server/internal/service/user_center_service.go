package service

import (
	"sort"
	"time"

	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserCenterService struct {
	db *gorm.DB
}

func NewUserCenterService(db *gorm.DB) *UserCenterService {
	return &UserCenterService{db: db}
}

type UserStats struct {
	ArticlesCount int   `json:"articlesCount"`
	ProjectsCount int   `json:"projectsCount"`
	KbDocsCount   int   `json:"kbDocsCount"`
	TotalViews    int64 `json:"totalViews"`
}

func (s *UserCenterService) GetStats(c *gin.Context) (*UserStats, error) {
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		return &UserStats{}, nil
	}
	userID := authUser.ID

	var articlesCount int64
	if err := s.db.Model(&model.Article{}).Where("author_id = ?", userID).Count(&articlesCount).Error; err != nil {
		return nil, err
	}

	var projectsCount int64
	if err := s.db.Model(&model.Project{}).Where("author_id = ?", userID).Count(&projectsCount).Error; err != nil {
		return nil, err
	}

	var kbDocsCount int64
	if err := s.db.Model(&model.KnowledgeDocument{}).Where("user_id = ?", userID).Count(&kbDocsCount).Error; err != nil {
		return nil, err
	}

	var totalViews int64
	if err := s.db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Where("author_id = ?", userID).Scan(&totalViews).Error; err != nil {
		return nil, err
	}

	return &UserStats{
		ArticlesCount: int(articlesCount),
		ProjectsCount: int(projectsCount),
		KbDocsCount:   int(kbDocsCount),
		TotalViews:    totalViews,
	}, nil
}

type ActivityItem struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *UserCenterService) GetRecentActivity(c *gin.Context) ([]ActivityItem, error) {
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		return nil, nil
	}
	userID := authUser.ID

	var articles []ActivityItem
	if err := s.db.Model(&model.Article{}).
		Select("id, 'article' as type, title, updated_at").
		Where("author_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&articles).Error; err != nil {
		return nil, err
	}

	var projects []ActivityItem
	if err := s.db.Model(&model.Project{}).
		Select("id, 'project' as type, title, updated_at").
		Where("author_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&projects).Error; err != nil {
		return nil, err
	}

	var kbDocs []ActivityItem
	if err := s.db.Model(&model.KnowledgeDocument{}).
		Select("id, 'kb_doc' as type, title, updated_at").
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Limit(10).
		Scan(&kbDocs).Error; err != nil {
		return nil, err
	}

	all := articles
	all = append(all, projects...)
	all = append(all, kbDocs...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].UpdatedAt.After(all[j].UpdatedAt)
	})

	if len(all) > 10 {
		all = all[:10]
	}

	return all, nil
}
