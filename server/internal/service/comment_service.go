package service

import (
	"strconv"
	"strings"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

type CommentPayload struct {
	Content  string `json:"content"`
	ParentID *uint  `json:"parentId"`
}

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) List(articleID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := s.db.
		Preload("User").
		Preload("Replies").
		Preload("Replies.User").
		Where("article_id = ? AND parent_id IS NULL", articleID).
		Order("created_at asc").
		Find(&comments).Error
	return comments, err
}

func (s *CommentService) Create(articleID, userID uint, payload CommentPayload) (*model.Comment, error) {
	content := strings.TrimSpace(payload.Content)
	if content == "" {
		return nil, ErrCommentEmpty
	}
	if err := validateModerationField("评论内容", content); err != nil {
		createModerationHit(s.db, userID, "comment_create", "评论内容", content, err)
		return nil, err
	}

	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	if user.Status == model.UserBanned {
		return nil, ErrUserBannedComment
	}

	var article model.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return nil, err
	}
	if article.Status != model.ArticlePublished {
		return nil, ErrArticleNotPublished
	}

	comment := model.Comment{
		Content:   content,
		ArticleID: articleID,
		UserID:    userID,
		ParentID:  payload.ParentID,
	}

	if payload.ParentID != nil {
		var parent model.Comment
		if err := s.db.First(&parent, *payload.ParentID).Error; err != nil {
			return nil, ErrParentCommentNotFound
		}
		if parent.ArticleID != articleID {
			return nil, ErrInvalidReplyComment
		}
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	articlePath := "/article/" + strconv.FormatUint(uint64(article.ID), 10)
	if article.AuthorID != userID {
		_ = createNotification(s.db, NotificationCreateInput{
			UserID:    article.AuthorID,
			Title:     "收到新评论",
			Content:   "你的文章收到了新的评论：" + article.Title,
			Type:      model.NotificationTypeArticleComment,
			ActionURL: articlePath,
			Payload: map[string]any{
				"articleId": article.ID,
				"commentId": comment.ID,
			},
		})
	}
	if payload.ParentID != nil {
		var parent model.Comment
		if err := s.db.First(&parent, *payload.ParentID).Error; err == nil && parent.UserID != userID {
			_ = createNotification(s.db, NotificationCreateInput{
				UserID:    parent.UserID,
				Title:     "收到评论回复",
				Content:   "你在文章下的评论收到了回复：" + article.Title,
				Type:      model.NotificationTypeCommentReply,
				ActionURL: articlePath,
				Payload: map[string]any{
					"articleId":      article.ID,
					"commentId":      parent.ID,
					"replyCommentId": comment.ID,
				},
			})
		}
	}

	var created model.Comment
	if err := s.db.Preload("User").First(&created, comment.ID).Error; err != nil {
		return nil, err
	}
	return &created, nil
}
