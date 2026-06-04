package service

import (
	"errors"
	"strings"
	"testing"

	"blog/server/internal/model"
)

func TestMetaServiceCreateCategoryTrimsNameAndBuildsSlug(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	category, err := service.CreateCategory("  Platform Engineering  ")
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if category.Name != "Platform Engineering" {
		t.Fatalf("expected trimmed name, got %q", category.Name)
	}
	if category.Slug != "platform-engineering" {
		t.Fatalf("expected slug platform-engineering, got %q", category.Slug)
	}
}

func TestMetaServiceCreateCategoryRejectsEmptyName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.CreateCategory("   ")
	if err == nil || err.Error() != "分类名称不能为空" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateCategoryRejectsTooLongName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.CreateCategory(strings.Repeat("长", 51))
	if err == nil || err.Error() != "分类名称不能超过50个字符" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateCategoryRejectsDuplicateName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	if _, err := service.CreateCategory("Backend"); err != nil {
		t.Fatalf("seed category: %v", err)
	}
	_, err := service.CreateCategory("backend")
	if err == nil || err.Error() != "分类已存在" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateCategoryRejectsDuplicateSlug(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	if _, err := service.CreateCategory("Go Lang"); err != nil {
		t.Fatalf("seed category: %v", err)
	}
	_, err := service.CreateCategory("go-lang")
	if err == nil || err.Error() != "分类已存在" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateCategoryRejectsInvalidSlug(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.CreateCategory("!!!")
	if err == nil || err.Error() != "分类名称无效" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceUpdateCategoryTrimsNameAndBuildsSlug(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	category, err := service.CreateCategory("Backend")
	if err != nil {
		t.Fatalf("seed category: %v", err)
	}

	updated, err := service.UpdateCategory(category.ID, "  API Design  ")
	if err != nil {
		t.Fatalf("update category: %v", err)
	}
	if updated.Name != "API Design" || updated.Slug != "api-design" {
		t.Fatalf("unexpected updated category: %+v", updated)
	}
}

func TestMetaServiceUpdateCategoryRejectsDuplicateName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	first, _ := service.CreateCategory("Backend")
	if _, err := service.CreateCategory("Frontend"); err != nil {
		t.Fatalf("seed second category: %v", err)
	}

	_, err := service.UpdateCategory(first.ID, "frontend")
	if err == nil || err.Error() != "分类已存在" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceUpdateCategoryRejectsInvalidName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	category, _ := service.CreateCategory("Backend")
	_, err := service.UpdateCategory(category.ID, "!!!")
	if err == nil || err.Error() != "分类名称无效" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceUpdateCategoryReturnsNotFound(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.UpdateCategory(999, "Backend")
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestMetaServiceDeleteCategoryRemovesUnusedCategory(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	category, _ := service.CreateCategory("Backend")
	if err := service.DeleteCategory(category.ID); err != nil {
		t.Fatalf("delete category: %v", err)
	}

	var count int64
	if err := db.Model(&model.Category{}).Where("id = ?", category.ID).Count(&count).Error; err != nil {
		t.Fatalf("count categories: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected category deleted, count=%d", count)
	}
}

func TestMetaServiceDeleteCategoryRejectsReferencedCategory(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	category, _ := service.CreateCategory("Backend")
	author := createTestUser(t, db, model.RoleUser)
	article := model.Article{
		Title:      "Referenced article",
		Summary:    "summary",
		Content:    "content",
		Status:     model.ArticleDraft,
		AuthorID:   author.ID,
		CategoryID: &category.ID,
	}
	if err := db.Create(&article).Error; err != nil {
		t.Fatalf("seed article: %v", err)
	}

	err := service.DeleteCategory(category.ID)
	if !errors.Is(err, ErrCategoryInUse) {
		t.Fatalf("expected ErrCategoryInUse, got %v", err)
	}
}

func TestMetaServiceDeleteCategoryReturnsNotFound(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	err := service.DeleteCategory(999)
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestMetaServiceCreateTagTrimsName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	tag, err := service.CreateTag("  Go  ")
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	if tag.Name != "Go" {
		t.Fatalf("expected trimmed name, got %q", tag.Name)
	}
}

func TestMetaServiceCreateTagRejectsEmptyName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.CreateTag("   ")
	if err == nil || err.Error() != "标签名称不能为空" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateTagRejectsTooLongName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.CreateTag(strings.Repeat("长", 31))
	if err == nil || err.Error() != "标签名称不能超过30个字符" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceCreateTagRejectsDuplicateName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	if _, err := service.CreateTag("Go"); err != nil {
		t.Fatalf("seed tag: %v", err)
	}
	_, err := service.CreateTag("go")
	if err == nil || err.Error() != "标签已存在" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceUpdateTagRenamesTag(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	tag, _ := service.CreateTag("Go")
	updated, err := service.UpdateTag(tag.ID, "  Gin  ")
	if err != nil {
		t.Fatalf("update tag: %v", err)
	}
	if updated.Name != "Gin" {
		t.Fatalf("expected updated name Gin, got %q", updated.Name)
	}
}

func TestMetaServiceUpdateTagRejectsDuplicateName(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	first, _ := service.CreateTag("Go")
	if _, err := service.CreateTag("Gin"); err != nil {
		t.Fatalf("seed second tag: %v", err)
	}
	_, err := service.UpdateTag(first.ID, "gin")
	if err == nil || err.Error() != "标签已存在" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetaServiceUpdateTagReturnsNotFound(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	_, err := service.UpdateTag(999, "Go")
	if !errors.Is(err, ErrTagNotFound) {
		t.Fatalf("expected ErrTagNotFound, got %v", err)
	}
}

func TestMetaServiceDeleteTagRemovesUnusedTag(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	tag, _ := service.CreateTag("Go")
	if err := service.DeleteTag(tag.ID); err != nil {
		t.Fatalf("delete tag: %v", err)
	}

	var count int64
	if err := db.Model(&model.Tag{}).Where("id = ?", tag.ID).Count(&count).Error; err != nil {
		t.Fatalf("count tags: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected tag deleted, count=%d", count)
	}
}

func TestMetaServiceDeleteTagRejectsReferencedTag(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	author := createTestUser(t, db, model.RoleUser)
	articleService := NewArticleService(db)
	article, err := articleService.Create(author.ID, model.RoleUser, ArticlePayload{
		Title:   "Tagged article",
		Summary: "summary",
		Content: "content",
		Tags:    []string{"Go"},
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}

	if len(article.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(article.Tags))
	}

	err = service.DeleteTag(article.Tags[0].ID)
	if !errors.Is(err, ErrTagInUse) {
		t.Fatalf("expected ErrTagInUse, got %v", err)
	}
}

func TestMetaServiceDeleteTagReturnsNotFound(t *testing.T) {
	db := newTestDB(t)
	service := NewMetaService(db)

	err := service.DeleteTag(999)
	if !errors.Is(err, ErrTagNotFound) {
		t.Fatalf("expected ErrTagNotFound, got %v", err)
	}
}
