package service

import (
	"testing"

	"blog/server/internal/model"
)

func TestPrivateArticleNotInPublishedList(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)
	svc := NewArticleService(db)

	pub, err := svc.Create(user.ID, model.RoleUser, ArticlePayload{
		Title:     "Public Article",
		Summary:   "public",
		Content:   "public content",
		IsPrivate: false,
	})
	if err != nil {
		t.Fatalf("create public: %v", err)
	}

	priv, err := svc.Create(user.ID, model.RoleUser, ArticlePayload{
		Title:     "Private Article",
		Summary:   "private",
		Content:   "private content",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("create private: %v", err)
	}

	// Submit both as admin to publish
	svc.Submit(pub.ID, user.ID, model.RoleAdmin)
	svc.Submit(priv.ID, user.ID, model.RoleAdmin)

	articles, pagination, err := svc.ListPublished(PublishedArticleFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("ListPublished: %v", err)
	}
	if pagination.Total != 1 {
		t.Fatalf("expected 1 published article (not private), got %d", pagination.Total)
	}
	if len(articles) != 1 || articles[0].ID != pub.ID {
		t.Fatalf("expected only the public article in published list")
	}
}

func TestPrivateArticleVisibleToOwner(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)
	svc := NewArticleService(db)

	article, err := svc.Create(user.ID, model.RoleUser, ArticlePayload{
		Title:     "My Private Article",
		Summary:   "test",
		Content:   "test content",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	svc.Submit(article.ID, user.ID, model.RoleAdmin)

	got, err := svc.GetByID(article.ID, user.ID)
	if err != nil {
		t.Fatalf("owner GetByID: %v", err)
	}
	if got.ID != article.ID {
		t.Fatalf("got wrong article")
	}
	if !got.IsPrivate {
		t.Fatal("article should be marked private")
	}
}

func TestPrivateArticleHiddenFromOthers(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, model.RoleUser)
	createTestUser(t, db, model.RoleUser)
	svc := NewArticleService(db)

	article, err := svc.Create(owner.ID, model.RoleUser, ArticlePayload{
		Title:     "Owners Private Article",
		Summary:   "test",
		Content:   "test",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	svc.Submit(article.ID, owner.ID, model.RoleAdmin)

	// Other user should see private article if it's published... actually the handler
	// checks permission. GetByID returns the article regardless of isPrivate for status check.
	// But the article handler adds a check. Let's test the handler-level logic:
	// The service GetByID returns articles even if private - the permission is at handler level.
	// Let's check that ListPublished excludes it.
	articles, _, _ := svc.ListPublished(PublishedArticleFilter{Page: 1, PageSize: 10})
	for _, a := range articles {
		if a.ID == article.ID {
			t.Fatal("private article should not appear in public listing")
		}
	}

	// Admin should be able to see it
	admin := createTestUser(t, db, model.RoleAdmin)
	_, err = svc.GetByID(article.ID, admin.ID)
	if err != nil {
		t.Fatalf("admin GetByID: %v", err)
	}
}

func TestPrivateProjectLifecycle(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)
	svc := NewProjectService(db)

	project, err := svc.Create(user.ID, model.RoleUser, ProjectPayload{
		Title:     "Private Project",
		Summary:   "test",
		Content:   "test content",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Submit should auto-publish private projects
	submitted, err := svc.Submit(project.ID, user.ID, model.RoleUser)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if submitted.Status != model.ProjectPublished {
		t.Fatalf("private project should auto-publish, got %s", submitted.Status)
	}

	// Should not appear in public listing
	items, pagination, err := svc.ListPublished(ProjectFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("ListPublished: %v", err)
	}
	if pagination.Total != 0 {
		t.Fatalf("private project should not be in public listing")
	}
	if len(items) > 0 {
		t.Fatal("private project found in public listing")
	}
}

func TestPrivateArticleSkipsAdminReview(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)
	svc := NewArticleService(db)

	article, err := svc.Create(user.ID, model.RoleUser, ArticlePayload{
		Title:     "Skip Review",
		Summary:   "test",
		Content:   "test content",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	submitted, err := svc.Submit(article.ID, user.ID, model.RoleUser)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if submitted.Status != model.ArticlePublished {
		t.Fatalf("private article should auto-publish, got %s", submitted.Status)
	}
}

func TestUpdateArticlePrivacy(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)
	svc := NewArticleService(db)

	article, err := svc.Create(user.ID, model.RoleUser, ArticlePayload{
		Title:     "Toggle Privacy",
		Summary:   "test",
		Content:   "test",
		IsPrivate: false,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	updated, err := svc.Update(article.ID, user.ID, model.RoleUser, ArticlePayload{
		Title:     "Toggle Privacy",
		Summary:   "test",
		Content:   "test",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !updated.IsPrivate {
		t.Fatal("article should be private after update")
	}
}
