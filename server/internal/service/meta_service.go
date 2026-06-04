package service

import (
	"errors"
	"strings"
	"unicode/utf8"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

const (
	categoryNameMaxRunes = 50
	tagNameMaxRunes      = 30
)

type MetaService struct {
	db *gorm.DB
}

type CategoryMetadata struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"articleCount"`
}

type TagMetadata struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"articleCount"`
}

func NewMetaService(db *gorm.DB) *MetaService {
	return &MetaService{db: db}
}

func (s *MetaService) List() (map[string]interface{}, error) {
	var categories []CategoryMetadata
	var tags []TagMetadata
	if err := s.db.Model(&model.Category{}).
		Select("categories.id, categories.name, categories.slug, COUNT(articles.id) AS article_count").
		Joins("LEFT JOIN articles ON articles.category_id = categories.id").
		Group("categories.id").
		Order("categories.id asc").
		Scan(&categories).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Tag{}).
		Select("tags.id, tags.name, COUNT(article_tags.article_id) AS article_count").
		Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
		Group("tags.id").
		Order("tags.name asc").
		Scan(&tags).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"categories": categories,
		"tags":       tags,
	}, nil
}

func (s *MetaService) CreateCategory(name string) (*model.Category, error) {
	name, slug, err := validateCategoryName(name)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCategoryAvailable(name, slug); err != nil {
		return nil, err
	}

	category := model.Category{Name: name, Slug: slug}
	if err := s.db.Create(&category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrCategoryExists
		}
		return nil, err
	}
	return &category, nil
}

func (s *MetaService) UpdateCategory(categoryID uint, name string) (*model.Category, error) {
	name, slug, err := validateCategoryName(name)
	if err != nil {
		return nil, err
	}

	var category model.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if err := s.ensureCategoryAvailable(name, slug, category.ID); err != nil {
		return nil, err
	}

	category.Name = name
	category.Slug = slug
	if err := s.db.Save(&category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrCategoryExists
		}
		return nil, err
	}
	return &category, nil
}

func (s *MetaService) DeleteCategory(categoryID uint) error {
	var category model.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	var articleCount int64
	if err := s.db.Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&articleCount).Error; err != nil {
		return err
	}
	if articleCount > 0 {
		return ErrCategoryInUse
	}

	return s.db.Delete(&category).Error
}

func (s *MetaService) CreateTag(name string) (*model.Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}
	if err := s.ensureTagAvailable(name); err != nil {
		return nil, err
	}

	tag := model.Tag{Name: name}
	if err := s.db.Create(&tag).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return &tag, nil
}

func (s *MetaService) UpdateTag(tagID uint, name string) (*model.Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}

	var tag model.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	if err := s.ensureTagAvailable(name, tag.ID); err != nil {
		return nil, err
	}

	tag.Name = name
	if err := s.db.Save(&tag).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return &tag, nil
}

func (s *MetaService) DeleteTag(tagID uint) error {
	var tag model.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTagNotFound
		}
		return err
	}

	var articleCount int64
	if err := s.db.Table("article_tags").Where("tag_id = ?", tagID).Count(&articleCount).Error; err != nil {
		return err
	}
	if articleCount > 0 {
		return ErrTagInUse
	}

	return s.db.Delete(&tag).Error
}

func validateCategoryName(name string) (string, string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", "", ErrCategoryNameEmpty
	}
	if utf8.RuneCountInString(name) > categoryNameMaxRunes {
		return "", "", ErrCategoryNameLong
	}

	slug := categorySlug(name)
	if slug == "" {
		return "", "", ErrCategoryNameInvalid
	}

	return name, slug, nil
}

func validateTagName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrTagNameEmpty
	}
	if utf8.RuneCountInString(name) > tagNameMaxRunes {
		return "", ErrTagNameLong
	}
	return name, nil
}

func (s *MetaService) ensureCategoryAvailable(name, slug string, excludeIDs ...uint) error {
	query := s.db.Where("lower(name) = ?", strings.ToLower(name)).Or("slug = ?", slug)
	if len(excludeIDs) > 0 {
		query = query.Not("id IN ?", excludeIDs)
	}

	var existing model.Category
	err := query.First(&existing).Error
	if err == nil {
		return ErrCategoryExists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (s *MetaService) ensureTagAvailable(name string, excludeIDs ...uint) error {
	query := s.db.Where("lower(name) = ?", strings.ToLower(name))
	if len(excludeIDs) > 0 {
		query = query.Not("id IN ?", excludeIDs)
	}

	var existing model.Tag
	err := query.First(&existing).Error
	if err == nil {
		return ErrTagExists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate entry") ||
		strings.Contains(message, "unique constraint failed")
}

func categorySlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	replacer := strings.NewReplacer(
		" ", "-",
		"_", "-",
		"/", "-",
		"\\", "-",
		".", "-",
		",", "-",
		":", "-",
		";", "-",
		"|", "-",
		"+", "-",
		"&", "-",
		"@", "-",
		"#", "-",
		"(", "-",
		")", "-",
		"[", "-",
		"]", "-",
		"{", "-",
		"}", "-",
		"?", "-",
		"!", "-",
		"'", "",
		"\"", "",
	)
	slug = replacer.Replace(slug)
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return slug
}
