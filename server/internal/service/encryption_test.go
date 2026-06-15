package service

import (
	"testing"
)

func TestEncryptionRoundTrip(t *testing.T) {
	svc, err := NewEncryptionService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("NewEncryptionService: %v", err)
	}

	plaintext := "sk-abc123def456ghi789jkl"
	cipherHex, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if cipherHex == "" {
		t.Fatal("cipherHex is empty")
	}

	decrypted, err := svc.Decrypt(cipherHex)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("decrypted != plaintext: %q vs %q", decrypted, plaintext)
	}
}

func TestEncryptionInvalidKeyLength(t *testing.T) {
	_, err := NewEncryptionService("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}

	_, err = NewEncryptionService("abc")
	if err == nil {
		t.Fatal("expected error for short key")
	}

	_, err = NewEncryptionService("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	if err == nil {
		t.Fatal("expected error for invalid hex key")
	}
}

func TestEncryptionDifferentKeys(t *testing.T) {
	plaintext := "sk-test-key-12345"

	svc1, _ := NewEncryptionService("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	svc2, _ := NewEncryptionService("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	c1, _ := svc1.Encrypt(plaintext)
	c2, _ := svc2.Encrypt(plaintext)

	if c1 == c2 {
		t.Fatal("different keys should produce different ciphertexts")
	}

	// svc2 should not be able to decrypt svc1's ciphertext
	_, err := svc2.Decrypt(c1)
	if err == nil {
		t.Fatal("expected decryption with wrong key to fail")
	}
}

func TestEncryptionUniqueNonces(t *testing.T) {
	svc, _ := NewEncryptionService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	plaintext := "same-text"
	c1, _ := svc.Encrypt(plaintext)
	c2, _ := svc.Encrypt(plaintext)

	if c1 == c2 {
		t.Fatal("encrypted outputs should differ due to unique nonces")
	}
}

func TestEncryptionTamperedCiphertext(t *testing.T) {
	svc, _ := NewEncryptionService("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	cipherHex, _ := svc.Encrypt("test-data")

	// tamper with the hex string
	tampered := cipherHex[:len(cipherHex)-4] + "ffff"
	_, err := svc.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected decryption of tampered data to fail")
	}
}
