package service

import (
	"testing"

	"blog/server/internal/model"
)

func newTestEncryption(t *testing.T) *EncryptionService {
	t.Helper()
	svc, err := NewEncryptionService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("NewEncryptionService: %v", err)
	}
	return svc
}

func TestApiKeyCreateAndList(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, "user")
	enc := newTestEncryption(t)
	svc := NewApiKeyService(db, enc)

	resp, err := svc.Create(user.ID, "deepseek", "sk-test-key-12345", "")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if resp.Provider != "deepseek" {
		t.Fatalf("provider: got %q, want deepseek", resp.Provider)
	}
	if resp.KeyPrefix == "" {
		t.Fatal("keyPrefix should not be empty")
	}
	if resp.KeyPrefix == "sk-test-key-12345" {
		t.Fatal("keyPrefix should be masked, not the full key")
	}

	list, err := svc.List(user.ID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 key, got %d", len(list))
	}
	if list[0].KeyPrefix != resp.KeyPrefix {
		t.Fatalf("prefix mismatch: %q vs %q", list[0].KeyPrefix, resp.KeyPrefix)
	}
}

func TestApiKeyDecrypt(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, "user")
	enc := newTestEncryption(t)
	svc := NewApiKeyService(db, enc)

	originalKey := "sk-original-secret-key-999"
	resp, err := svc.Create(user.ID, "openai", originalKey, "https://api.openai.com/v1")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	decrypted, provider, baseURL, err := svc.GetDecryptedKey(resp.ID, user.ID)
	if err != nil {
		t.Fatalf("GetDecryptedKey: %v", err)
	}
	if decrypted != originalKey {
		t.Fatalf("decrypted key mismatch: got %q, want %q", decrypted, originalKey)
	}
	if provider != "openai" {
		t.Fatalf("provider: got %q", provider)
	}
	if baseURL != "https://api.openai.com/v1" {
		t.Fatalf("baseURL: got %q", baseURL)
	}
}

func TestApiKeyAccessControl(t *testing.T) {
	db := newTestDB(t)
	user1 := createTestUser(t, db, "user")
	user2 := createTestUser(t, db, "user")
	enc := newTestEncryption(t)
	svc := NewApiKeyService(db, enc)

	resp, _ := svc.Create(user1.ID, "deepseek", "sk-user1-key", "")

	// user2 should not be able to access user1's key
	_, _, _, err := svc.GetDecryptedKey(resp.ID, user2.ID)
	if err == nil {
		t.Fatal("expected access denied error")
	}

	err = svc.Delete(resp.ID, user2.ID)
	if err == nil {
		t.Fatal("expected access denied on delete")
	}

	err = svc.Delete(resp.ID, user1.ID)
	if err != nil {
		t.Fatalf("user1 should be able to delete own key: %v", err)
	}
}

func TestApiKeyInvalidProvider(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, "user")
	enc := newTestEncryption(t)
	svc := NewApiKeyService(db, enc)

	_, err := svc.Create(user.ID, "invalid-provider", "sk-key", "")
	if err == nil {
		t.Fatal("expected error for invalid provider")
	}

	_, err = svc.Create(user.ID, "deepseek", "", "")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApiKeyStoreEncrypted(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, "user")
	enc := newTestEncryption(t)
	svc := NewApiKeyService(db, enc)

	original := "sk-raw-secret-plaintext"
	svc.Create(user.ID, "deepseek", original, "")

	var stored model.UserApiKey
	if err := db.First(&stored).Error; err != nil {
		t.Fatalf("find stored key: %v", err)
	}
	if stored.EncryptedKey == original {
		t.Fatal("key should be encrypted in DB, not stored in plaintext")
	}

	decrypted, _ := enc.Decrypt(stored.EncryptedKey)
	if decrypted != original {
		t.Fatalf("stored encrypted key does not match original")
	}
}
