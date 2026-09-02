package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGenerated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	RequestID()(c)

	id := w.Header().Get("X-Request-ID")
	if id == "" {
		t.Fatal("X-Request-ID not generated")
	}
	if len(id) < 36 {
		t.Fatalf("X-Request-ID too short: %q", id)
	}
}

func TestRequestIDPreserved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("X-Request-ID", "my-custom-id")

	RequestID()(c)

	got := w.Header().Get("X-Request-ID")
	if got != "my-custom-id" {
		t.Fatalf("X-Request-ID = %q, want %q", got, "my-custom-id")
	}
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	RequestID()(c)

	id := GetRequestID(c)
	if id == "" {
		t.Fatal("GetRequestID returned empty")
	}
}

func TestGetRequestIDMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	id := GetRequestID(c)
	if id != "" {
		t.Fatalf("GetRequestID on fresh context = %q, want empty", id)
	}
}
