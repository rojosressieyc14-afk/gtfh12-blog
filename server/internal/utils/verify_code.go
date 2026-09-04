package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type verifyEntry struct {
	code      string
	expiresAt time.Time
}

var (
	verifyStore = make(map[string]*verifyEntry)
	verifyMu    sync.Mutex
)

func init() {
	go cleanupVerifyStore()
}

func cleanupVerifyStore() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		verifyMu.Lock()
		now := time.Now()
		for email, entry := range verifyStore {
			if now.After(entry.expiresAt) {
				delete(verifyStore, email)
			}
		}
		verifyMu.Unlock()
	}
}

func GenerateVerifyCode(email string) string {
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = '0' + byte(n.Int64())
	}
	codeStr := string(code)

	verifyMu.Lock()
	verifyStore[email] = &verifyEntry{
		code:      codeStr,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	verifyMu.Unlock()

	return codeStr
}

func VerifyCode(email, code string) bool {
	verifyMu.Lock()
	defer verifyMu.Unlock()

	entry, ok := verifyStore[email]
	if !ok {
		return false
	}
	delete(verifyStore, email)

	if time.Now().After(entry.expiresAt) {
		return false
	}
	return entry.code == code
}

func CanSendVerifyCode(email string) (bool, string) {
	verifyMu.Lock()
	defer verifyMu.Unlock()

	entry, ok := verifyStore[email]
	if !ok {
		return true, ""
	}

	remaining := time.Until(entry.expiresAt.Add(-5 * time.Minute))
	if remaining > 0 {
		secs := int(remaining.Seconds()) + 1
		return false, fmt.Sprintf("%ds", secs)
	}
	return true, ""
}
