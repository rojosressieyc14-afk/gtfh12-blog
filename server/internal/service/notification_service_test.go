package service

import (
	"fmt"
	"testing"

	"blog/server/internal/model"
)

func TestNotificationListSupportsPaginationAndUnreadFilter(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, model.RoleUser)

	for i := 0; i < 5; i++ {
		if err := createNotification(db, NotificationCreateInput{
			UserID:  user.ID,
			Title:   fmt.Sprintf("notice-%d", i),
			Content: "content",
			Type:    model.NotificationTypeSystem,
		}); err != nil {
			t.Fatalf("create notification %d: %v", i, err)
		}
	}

	var first model.Notification
	if err := db.Where("user_id = ?", user.ID).Order("id asc").First(&first).Error; err != nil {
		t.Fatalf("query first notification: %v", err)
	}
	if err := db.Model(&model.Notification{}).Where("id = ?", first.ID).Update("is_read", true).Error; err != nil {
		t.Fatalf("mark first notification read: %v", err)
	}

	service := NewNotificationService(db)

	items, pagination, err := service.List(user.ID, NotificationListFilter{
		Page:     1,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items on first page, got %d", len(items))
	}
	if pagination.Total != 5 {
		t.Fatalf("expected total 5, got %d", pagination.Total)
	}

	unreadItems, unreadPagination, err := service.List(user.ID, NotificationListFilter{
		UnreadOnly: true,
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("list unread notifications: %v", err)
	}
	if len(unreadItems) != 4 {
		t.Fatalf("expected 4 unread notifications, got %d", len(unreadItems))
	}
	if unreadPagination.Total != 4 {
		t.Fatalf("expected unread total 4, got %d", unreadPagination.Total)
	}
}
