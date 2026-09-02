package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	defer rl.Stop()

	if !rl.Allow("ip1") {
		t.Fatal("first request should be allowed")
	}
	if !rl.Allow("ip1") {
		t.Fatal("second request should be allowed")
	}
	if !rl.Allow("ip1") {
		t.Fatal("third request should be allowed")
	}
	if rl.Allow("ip1") {
		t.Fatal("fourth request should be blocked")
	}
}

func TestRateLimiterDifferentKeys(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	defer rl.Stop()

	if !rl.Allow("ip1") {
		t.Fatal("ip1 first request should be allowed")
	}
	if rl.Allow("ip1") {
		t.Fatal("ip1 second request should be blocked")
	}
	if !rl.Allow("ip2") {
		t.Fatal("ip2 should have its own limit")
	}
}

func TestRateLimiterWindowReset(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	defer rl.Stop()

	if !rl.Allow("ip1") {
		t.Fatal("first request should be allowed")
	}
	if rl.Allow("ip1") {
		t.Fatal("second request should be blocked")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("ip1") {
		t.Fatal("request after window reset should be allowed")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(2, time.Minute)
	defer rl.Stop()

	router := gin.New()
	router.Use(RateLimit(rl))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "1.2.3.4:1234"
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("request 1: status %d, want 200", w1.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "1.2.3.4:5678"
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("request 2: status %d, want 200", w2.Code)
	}

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "1.2.3.4:9012"
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusTooManyRequests {
		t.Fatalf("request 3: status %d, want 429", w3.Code)
	}
}

func TestRedisRateLimitFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RedisRateLimit("test", 5, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Redis rate limit with nil RDB should pass through, got %d", w.Code)
	}
}
