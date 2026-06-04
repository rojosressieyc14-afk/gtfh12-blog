package service

import (
	"blog/server/internal/model"
	"gorm.io/gorm"
)

type NotificationListFilter struct {
	UnreadOnly bool
	Type       string
	Page       int
	PageSize   int
}

type NotificationCreateInput struct {
	UserID    uint
	Title     string
	Content   string
	Type      string
	ActionURL string
	Payload   map[string]any
}

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) List(userID uint, filter NotificationListFilter) ([]model.Notification, Pagination, error) {
	var items []model.Notification
	var total int64
	pagination := normalizePagination(filter.Page, filter.PageSize)

	query := s.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if filter.UnreadOnly {
		query = query.Where("is_read = ?", false)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, pagination, err
	}
	err := query.Order("created_at desc").
		Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).
		Find(&items).Error
	pagination.Total = total
	return items, pagination, err
}

func (s *NotificationService) MarkRead(userID, notificationID uint) error {
	result := s.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func createNotification(db *gorm.DB, input NotificationCreateInput) error {
	if input.Type == "" {
		input.Type = model.NotificationTypeSystem
	}
	return db.Create(&model.Notification{
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		Type:      input.Type,
		ActionURL: input.ActionURL,
		Payload:   input.Payload,
	}).Error
}
