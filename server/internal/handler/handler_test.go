package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"blog/server/internal/config"
	"blog/server/internal/database"
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
	token := result["token"].(string)
	if token == "" {
		t.Fatal("register: expected token")
	}

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

	adminToken, err := utils.GenerateJWT(admin.ID, admin.Username, model.RoleAdmin)
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
	userToken, err := utils.GenerateJWT(user.ID, user.Username, model.RoleUser)
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
