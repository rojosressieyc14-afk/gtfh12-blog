package service

import (
	"errors"
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

var (
	ErrAPIKeyNotFound = errors.New("API Key 不存在")
	ErrAPIKeyAccess   = errors.New("无权访问该 API Key")
	ErrInvalidProvider = errors.New("不支持的 Provider，仅支持 deepseek 或 openai")
	ErrAPIKeyEmpty    = errors.New("API Key 不能为空")
)

type ApiKeyResponse struct {
	ID         uint       `json:"id"`
	Provider   string     `json:"provider"`
	KeyPrefix  string     `json:"keyPrefix"`
	BaseURL    string     `json:"baseURL"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type ApiKeyService struct {
	db     *gorm.DB
	crypto *EncryptionService
}

func NewApiKeyService(db *gorm.DB, crypto *EncryptionService) *ApiKeyService {
	return &ApiKeyService{db: db, crypto: crypto}
}

func (s *ApiKeyService) Create(userID uint, provider, apiKey, baseURL string) (*ApiKeyResponse, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	if provider != "deepseek" && provider != "openai" {
		return nil, ErrInvalidProvider
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	encrypted, err := s.crypto.Encrypt(apiKey)
	if err != nil {
		return nil, err
	}

	prefix := apiKey
	if len(apiKey) > 8 {
		prefix = apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
	} else {
		prefix = apiKey[:2] + "****"
	}

	if baseURL == "" {
		switch provider {
		case "deepseek":
			baseURL = "https://api.deepseek.com/v1"
		case "openai":
			baseURL = "https://api.openai.com/v1"
		}
	}

	item := model.UserApiKey{
		UserID:       userID,
		Provider:     provider,
		EncryptedKey: encrypted,
		KeyPrefix:    prefix,
		BaseURL:      strings.TrimRight(baseURL, "/"),
	}
	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}

	return s.toResponse(&item), nil
}

func (s *ApiKeyService) List(userID uint) ([]ApiKeyResponse, error) {
	var items []model.UserApiKey
	if err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}

	resp := make([]ApiKeyResponse, len(items))
	for i, item := range items {
		resp[i] = *s.toResponse(&item)
	}
	return resp, nil
}

func (s *ApiKeyService) Delete(keyID, userID uint) error {
	var item model.UserApiKey
	if err := s.db.First(&item, keyID).Error; err != nil {
		return ErrAPIKeyNotFound
	}
	if item.UserID != userID {
		return ErrAPIKeyAccess
	}
	return s.db.Delete(&item).Error
}

func (s *ApiKeyService) GetDecryptedKey(keyID, userID uint) (string, string, string, error) {
	var item model.UserApiKey
	if err := s.db.First(&item, keyID).Error; err != nil {
		return "", "", "", ErrAPIKeyNotFound
	}
	if item.UserID != userID {
		return "", "", "", ErrAPIKeyAccess
	}

	decrypted, err := s.crypto.Decrypt(item.EncryptedKey)
	if err != nil {
		return "", "", "", err
	}

	now := time.Now()
	s.db.Model(&item).UpdateColumn("last_used_at", &now)

	return decrypted, item.Provider, item.BaseURL, nil
}

func (s *ApiKeyService) toResponse(item *model.UserApiKey) *ApiKeyResponse {
	return &ApiKeyResponse{
		ID:         item.ID,
		Provider:   item.Provider,
		KeyPrefix:  item.KeyPrefix,
		BaseURL:    item.BaseURL,
		LastUsedAt: item.LastUsedAt,
		CreatedAt:  item.CreatedAt,
	}
}
