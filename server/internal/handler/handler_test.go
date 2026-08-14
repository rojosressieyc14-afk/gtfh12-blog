package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"blog/server/internal/config"
	"blog/server/internal/database"
	"blog/server/internal/middleware"
	"blog/server/internal/model"
	"blog/server/internal/router"
	"blog/server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testSeq       uint64
	initJWTOnce   sync.Once
	jwtPrivateKey string
	jwtPublicKey  string
)

const testKeyDir = "./test-keys"

func initTestKeys(t *testing.T) {
	t.Helper()
	initJWTOnce.Do(func() {
		if err := os.MkdirAll(testKeyDir, 0700); err != nil {
			t.Fatalf("mkdir test keys: %v", err)
		}

		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("generate test key: %v", err)
		}

		privPath := filepath.Join(testKeyDir, "private.pem")
		pubPath := filepath.Join(testKeyDir, "public.pem")

		privFile, err := os.Create(privPath)
		if err != nil {
			t.Fatalf("create test private key: %v", err)
		}
		defer privFile.Close()
		if err := pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
			t.Fatalf("encode test private key: %v", err)
		}

		pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			t.Fatalf("marshal test public key: %v", err)
		}
		pubFile, err := os.Create(pubPath)
		if err != nil {
			t.Fatalf("create test public key: %v", err)
		}
		defer pubFile.Close()
		if err := pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}); err != nil {
			t.Fatalf("encode test public key: %v", err)
		}

		jwtPrivateKey = privPath
		jwtPublicKey = pubPath

		if err := utils.InitKeys(jwtPrivateKey, jwtPublicKey); err != nil {
			t.Fatalf("init JWT keys: %v", err)
		}
	})
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func setupTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	initTestKeys(t)
	gin.SetMode(gin.TestMode)
	db := newTestDB(t)

	cfg := config.Config{
		ServerPort:        "8080",
		JWTPrivateKeyPath: jwtPrivateKey,
		JWTPublicKeyPath:  jwtPublicKey,
		DBHost:            "localhost",
		DBPort:            "3306",
		DBUser:            "root",
		DBName:            "blog_test",
		WebOrigin:         "http://localhost:5173",
		AdminOrigin:       "http://localhost:5174",
		GinMode:           gin.TestMode,
		UploadDir:         "./test-uploads",
	}
	r := router.New(cfg, db)
	return r, db
}

func requestJSON(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func authenticatedRequest(r *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("parse response: %v, body: %s", err, w.Body.String())
	}
	return result
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupTest(t)
	w := requestJSON(r, "GET", "/api/health", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	result := parseResponse(t, w)
	if result["message"] != "ok" {
		t.Fatalf("expected message ok, got %v", result["message"])
	}
}

func TestRegisterLoginMeFlow(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	username := fmt.Sprintf("testuser_%d", seq)
	password := "TestPass123!"

	// Register
	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	result := parseResponse(t, w)
	if _, ok := result["user"]; !ok {
		t.Fatalf("register: expected user in response, got %v", result)
	}
	cookies := w.Result().Cookies()
	var tokenCookie *http.Cookie
	for _, ck := range cookies {
		if ck.Name == middleware.AuthCookieName {
			tokenCookie = ck
			break
		}
	}
	if tokenCookie == nil || tokenCookie.Value == "" {
		t.Fatal("register: expected blog_token cookie")
	}
	token := tokenCookie.Value

	// Me
	w = authenticatedRequest(r, "GET", "/api/auth/me", token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// Login
	w = requestJSON(r, "POST", "/api/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestLoginRememberCookie(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	username := fmt.Sprintf("remember_%d", seq)
	password := "TestPass123!"

	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// 不勾选 -> 会话级 cookie（MaxAge == 0）
	w = requestJSON(r, "POST", "/api/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	authCookie := findCookie(t, w, middleware.AuthCookieName)
	if authCookie.MaxAge != 0 {
		t.Fatalf("login without remember: expected session cookie (MaxAge 0), got %d", authCookie.MaxAge)
	}
	nonRememberClaims, nonRememberErr := utils.ParseJWT(authCookie.Value)
	if nonRememberErr != nil {
		t.Fatalf("parse non-remember token: %v", nonRememberErr)
	}
	nonRememberRemaining := time.Until(nonRememberClaims.ExpiresAt.Time)
	if nonRememberRemaining > 25*time.Hour {
		t.Fatalf("non-remember token should expire within 24h, remaining %v", nonRememberRemaining)
	}

	// 勾选 -> 30 天
	w = requestJSON(r, "POST", "/api/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","remember":true}`, username, password))
	if w.Code != http.StatusOK {
		t.Fatalf("login remember: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	authCookie = findCookie(t, w, middleware.AuthCookieName)
	if authCookie.MaxAge != middleware.RememberCookieMaxAge {
		t.Fatalf("login with remember: expected MaxAge %d, got %d", middleware.RememberCookieMaxAge, authCookie.MaxAge)
	}
	tokenClaims, err := utils.ParseJWT(authCookie.Value)
	if err != nil {
		t.Fatalf("parse remember token: %v", err)
	}
	remaining := time.Until(tokenClaims.ExpiresAt.Time)
	if remaining < 29*24*time.Hour {
		t.Fatalf("remember token should last ~30 days, remaining %v", remaining)
	}
}

func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, ck := range w.Result().Cookies() {
		if ck.Name == name {
			return ck
		}
	}
	t.Fatalf("expected cookie %q", name)
	return nil
}

func TestAdminLogin(t *testing.T) {
	r, db := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	admin := model.User{
		Username: fmt.Sprintf("admin_%d", seq),
		Password: "hashed",
		Role:     model.RoleAdmin,
		Status:   model.UserActive,
	}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	adminToken, err := utils.GenerateJWT(admin.ID, admin.Username, model.RoleAdmin, utils.TokenTTLDefault)
	if err != nil {
		t.Fatalf("generate admin token: %v", err)
	}

	user := model.User{
		Username: fmt.Sprintf("user_%d", seq),
		Password: "hashed",
		Role:     model.RoleUser,
		Status:   model.UserActive,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	userToken, err := utils.GenerateJWT(user.ID, user.Username, model.RoleUser, utils.TokenTTLDefault)
	if err != nil {
		t.Fatalf("generate user token: %v", err)
	}

	w := authenticatedRequest(r, "GET", "/api/admin/dashboard", adminToken, "")
	if w.Code != http.StatusOK {
		t.Fatalf("admin dashboard by admin: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	w = authenticatedRequest(r, "GET", "/api/admin/dashboard", userToken, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("admin dashboard by user: expected 403, got %d, body: %s", w.Code, w.Body.String())
	}
}
