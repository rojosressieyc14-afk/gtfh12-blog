package database

import (
	"fmt"

	"blog/server/internal/config"
	"blog/server/internal/model"
	"blog/server/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.ArticleReview{},
		&model.Comment{},
		&model.ArticleLike{},
		&model.ArticleFavorite{},
		&model.Project{},
		&model.Notification{},
		&model.OperationLog{},
		&model.SensitiveWord{},
		&model.ModerationHit{},
		&model.SystemSetting{},
		&model.AIReviewRecord{},
		&model.InterviewSession{},
		&model.InterviewRound{},
		&model.KnowledgeBase{},
		&model.KnowledgeDocument{},
		&model.KbDocumentTag{},
		&model.UserApiKey{},
	)
}

func SeedAdmin(db *gorm.DB, cfg config.Config) error {
	var count int64
	if err := db.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := utils.HashPassword(cfg.DefaultPass)
	if err != nil {
		return err
	}

	admin := model.User{
		Username: cfg.DefaultAdmin,
		Password: hash,
		Role:     model.RoleAdmin,
		Status:   model.UserActive,
		Headline: "Site owner / administrator",
		Bio:      "Built with Go, Gin, Vue 3 and a love for shipping.",
	}

	return db.Create(&admin).Error
}

func SeedCategories(db *gorm.DB) error {
	defaults := []model.Category{
		{Name: "General", Slug: "default"},
		{Name: "Tech Notes", Slug: "tech-notes"},
		{Name: "Product Design", Slug: "product-design"},
	}

	for _, item := range defaults {
		var category model.Category
		err := db.Where("slug = ?", item.Slug).First(&category).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&item).Error; err != nil {
				return err
			}
			continue
		}

		if category.Name != item.Name {
			if err := db.Model(&category).Update("name", item.Name).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
