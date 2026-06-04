package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

var testSeq uint64

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
	gin.SetMode(gin.TestMode)
	db := newTestDB(t)

	cfg := config.Config{
		ServerPort:  "8080",
		JWTSecret:   "test-secret-for-handler-tests",
		DBHost:      "localhost",
		DBPort:      "3306",
		DBUser:      "root",
		DBName:      "blog_test",
		WebOrigin:   "http://localhost:5173",
		AdminOrigin: "http://localhost:5174",
		GinMode:     gin.TestMode,
		UploadDir:   "./test-uploads",
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
	body := fmt.Sprintf(`{"username":"%s","password":"password123"}`, username)
	w := requestJSON(r, "POST", "/api/auth/register", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	regResult := parseResponse(t, w)
	token, ok := regResult["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected non-empty token, got %v", regResult["token"])
	}

	w2 := authenticatedRequest(r, "GET", "/api/auth/me", token, "")
	if w2.Code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d, body: %s", w2.Code, w2.Body.String())
	}
	meResult := parseResponse(t, w2)
	user, ok := meResult["user"].(map[string]any)
	if !ok {
		t.Fatalf("expected user object, got %v", meResult["user"])
	}
	if user["username"] != username {
		t.Fatalf("expected username %q, got %q", username, user["username"])
	}
	if user["role"] != "user" {
		t.Fatalf("expected role user, got %v", user["role"])
	}

	w3 := requestJSON(r, "GET", "/api/auth/me", "")
	if w3.Code != http.StatusUnauthorized {
		t.Fatalf("me without token: expected 401, got %d", w3.Code)
	}
}

func TestAdminAccessProtected(t *testing.T) {
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

	adminToken, err := utils.GenerateJWT("test-secret-for-handler-tests", admin.ID, admin.Username, model.RoleAdmin)
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
	userToken, err := utils.GenerateJWT("test-secret-for-handler-tests", user.ID, user.Username, model.RoleUser)
	if err != nil {
		t.Fatalf("generate user token: %v", err)
	}

	w := authenticatedRequest(r, "GET", "/api/admin/dashboard", adminToken, "")
	if w.Code != http.StatusOK {
		t.Fatalf("admin dashboard by admin: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	w2 := authenticatedRequest(r, "GET", "/api/admin/dashboard", userToken, "")
	if w2.Code != http.StatusForbidden {
		t.Fatalf("admin dashboard by user: expected 403, got %d, body: %s", w2.Code, w2.Body.String())
	}

	w3 := requestJSON(r, "GET", "/api/admin/dashboard", "")
	if w3.Code != http.StatusUnauthorized {
		t.Fatalf("admin dashboard without token: expected 401, got %d", w3.Code)
	}
}

func TestArticleCreateSubmitViaHandlers(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	authorUser := fmt.Sprintf("author_%d", seq)
	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"password123"}`, authorUser))
	if w.Code != http.StatusCreated {
		t.Fatalf("register author: expected 201, got %d", w.Code)
	}
	authorResult := parseResponse(t, w)
	authorToken := authorResult["token"].(string)

	w3 := authenticatedRequest(r, "POST", "/api/articles", authorToken, `{"title":"Test Article","summary":"A test","content":"This is test content","tags":["go","test"]}`)
	if w3.Code != http.StatusCreated && w3.Code != http.StatusOK {
		t.Fatalf("create article: expected 2xx, got %d, body: %s", w3.Code, w3.Body.String())
	}
	articleResult := parseResponse(t, w3)
	articleItem, ok := articleResult["item"].(map[string]any)
	if !ok {
		articleItem = articleResult
	}
	articleID := int(articleItem["id"].(float64))
	if articleItem["status"].(string) != string(model.ArticleDraft) {
		t.Fatalf("expected draft status, got %s", articleItem["status"])
	}

	w4 := authenticatedRequest(r, "POST", fmt.Sprintf("/api/articles/%d/submit", articleID), authorToken, "")
	if w4.Code != http.StatusOK {
		t.Fatalf("submit article: expected 200, got %d, body: %s", w4.Code, w4.Body.String())
	}
	submitResult := parseResponse(t, w4)
	submitItem := submitResult
	if item, ok := submitResult["item"].(map[string]any); ok {
		submitItem = item
	}
	if submitItem["status"].(string) != string(model.ArticlePending) {
		t.Fatalf("expected pending after submit, got %s", submitItem["status"])
	}

	w5 := authenticatedRequest(r, "GET", "/api/my/articles", authorToken, "")
	if w5.Code != http.StatusOK {
		t.Fatalf("my articles: expected 200, got %d, body: %s", w5.Code, w5.Body.String())
	}
}

func TestProjectCreateSubmitFlow(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	authorUser := fmt.Sprintf("projauthor_%d", seq)
	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"password123"}`, authorUser))
	if w.Code != http.StatusCreated {
		t.Fatalf("register author: expected 201, got %d", w.Code)
	}
	authorResult := parseResponse(t, w)
	authorToken := authorResult["token"].(string)

	w3 := authenticatedRequest(r, "POST", "/api/projects", authorToken, `{"title":"Test Project","summary":"A project","content":"Project description"}`)
	if w3.Code != http.StatusCreated && w3.Code != http.StatusOK {
		t.Fatalf("create project: expected 2xx, got %d, body: %s", w3.Code, w3.Body.String())
	}
	projectResult := parseResponse(t, w3)
	projItem, ok := projectResult["item"].(map[string]any)
	if !ok {
		projItem = projectResult
	}
	projectID := int(projItem["id"].(float64))
	if projItem["status"].(string) != string(model.ProjectDraft) {
		t.Fatalf("expected draft status, got %s", projItem["status"])
	}

	w4 := authenticatedRequest(r, "POST", fmt.Sprintf("/api/projects/%d/submit", projectID), authorToken, "")
	if w4.Code != http.StatusOK {
		t.Fatalf("submit project: expected 200, got %d, body: %s", w4.Code, w4.Body.String())
	}
	submitResult := parseResponse(t, w4)
	submitItem := submitResult
	if item, ok := submitResult["item"].(map[string]any); ok {
		submitItem = item
	}
	if submitItem["status"].(string) != string(model.ProjectPending) {
		t.Fatalf("expected pending after submit, got %s", submitItem["status"])
	}

	w5 := authenticatedRequest(r, "GET", "/api/my/projects", authorToken, "")
	if w5.Code != http.StatusOK {
		t.Fatalf("my projects: expected 200, got %d, body: %s", w5.Code, w5.Body.String())
	}
}

func TestPublicEndpointsAccessibleWithoutAuth(t *testing.T) {
	r, _ := setupTest(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/health"},
		{"GET", "/api/metadata"},
		{"GET", "/api/articles"},
		{"GET", "/api/projects"},
	}

	for _, ep := range endpoints {
		w := requestJSON(r, ep.method, ep.path, "")
		if w.Code != http.StatusOK {
			t.Fatalf("%s %s: expected 200, got %d, body: %s", ep.method, ep.path, w.Code, w.Body.String())
		}
	}
}

func TestRegisterRejectsDuplicateUsername(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	username := fmt.Sprintf("dupuser_%d", seq)
	body := fmt.Sprintf(`{"username":"%s","password":"password123"}`, username)

	w := requestJSON(r, "POST", "/api/auth/register", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("first register: expected 201, got %d", w.Code)
	}

	w2 := requestJSON(r, "POST", "/api/auth/register", body)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("duplicate register: expected 400, got %d, body: %s", w2.Code, w2.Body.String())
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	username := fmt.Sprintf("shortpw_%d", seq)
	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"abc"}`, username))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("short password: expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	r, _ := setupTest(t)

	w := requestJSON(r, "POST", "/api/auth/login", `{"username":"nonexistent","password":"wrongpass"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("invalid login: expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}
