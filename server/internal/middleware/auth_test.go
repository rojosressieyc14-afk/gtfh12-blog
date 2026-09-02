package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	user := GetAuthUser(c)
	if user != nil {
		t.Fatal("GetAuthUser should return nil on fresh context")
	}
}

func TestGetAuthUserWithType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("authUser", AuthUser{ID: 1, Username: "test", Role: "admin", Status: "active"})

	user := GetAuthUser(c)
	if user == nil {
		t.Fatal("GetAuthUser should return user")
	}
	if user.Username != "test" {
		t.Fatalf("Username = %q, want %q", user.Username, "test")
	}
	if user.Role != "admin" {
		t.Fatalf("Role = %q, want %q", user.Role, "admin")
	}
}

func TestGetAuthUserWithWrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("authUser", "not-an-AuthUser")

	user := GetAuthUser(c)
	if user != nil {
		t.Fatal("GetAuthUser should return nil for wrong type")
	}
}

func TestRequireAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authUser", AuthUser{ID: 1, Role: "admin"})
		c.Next()
	})
	router.Use(RequireAdmin())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("admin request: status %d, want 200", w.Code)
	}
}

func TestRequireAdminRejectsUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authUser", AuthUser{ID: 2, Role: "user"})
		c.Next()
	})
	router.Use(RequireAdmin())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("user request: status %d, want 403", w.Code)
	}
}

func TestRequireAdminRejectsNoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequireAdmin())
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("no-auth request: status %d, want 403", w.Code)
	}
}

func TestCSRFGetSkipsValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CSRF())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET should skip CSRF validation, got %d", w.Code)
	}
}

func TestCSRFPostRequiresToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CSRF())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "abc123"})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("POST without CSRF header: status %d, want 403", w.Code)
	}
}

func TestCSRFPostWithMatchingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CSRF())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "abc123"})
	req.Header.Set("X-CSRF-Token", "abc123")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST with matching CSRF token: status %d, want 200", w.Code)
	}
}
