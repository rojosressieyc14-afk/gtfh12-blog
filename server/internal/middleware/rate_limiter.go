package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"blog/server/internal/cache"
	"github.com/gin-gonic/gin"
)

type rateEntry struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateEntry
	limit   int
	window  time.Duration
	ticker  *time.Ticker
	stop    chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateEntry),
		limit:   limit,
		window:  window,
		ticker:  time.NewTicker(window),
		stop:    make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	close(rl.stop)
	rl.ticker.Stop()
}

func (rl *RateLimiter) cleanup() {
	for {
		select {
		case <-rl.ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, entry := range rl.entries {
				if now.After(entry.resetAt) {
					delete(rl.entries, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stop:
			return
		}
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.entries[key]

	if !exists || now.After(entry.resetAt) {
		rl.entries[key] = &rateEntry{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RedisRateLimit uses Redis INCR + EXPIRE for distributed rate limiting.
// Falls back to in-memory if Redis is unavailable.
func RedisRateLimit(prefix string, limit int, window time.Duration) gin.HandlerFunc {
	fallback := NewRateLimiter(limit, window)
	return func(c *gin.Context) {
		if cache.RDB == nil {
			RateLimit(fallback)(c)
			return
		}

		key := "ratelimit:" + prefix + ":" + c.ClientIP()
		ctx := context.Background()

		count, err := cache.RDB.Incr(ctx, key).Result()
		if err != nil {
			RateLimit(fallback)(c)
			return
		}

		if count == 1 {
			cache.RDB.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
