package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"blog/server/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	AuthCookieName       = "blog_token"
	CSRFCookieName       = "blog_csrf"
	AuthCookieMaxAge     = 7 * 24 * 3600
	RememberCookieMaxAge = 30 * 24 * 3600
)

type AuthUser struct {
	ID       uint
	Username string
	Role     string
	Status   string
}

func secureRequest(c *gin.Context) bool {
	return c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
}

func randomToken() string {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		}
		if token == "" {
			if cookie, err := c.Cookie(AuthCookieName); err == nil {
				token = cookie
			}
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "缺少登录凭证"})
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "登录凭证无效或已过期"})
			c.Abort()
			return
		}

		c.Set("authUser", AuthUser{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
			Status:   claims.Status,
		})
		if claims.Status == "banned" {
			c.JSON(http.StatusForbidden, gin.H{"message": "账户已被封禁"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetSessionCookies(c *gin.Context, token string, maxAge int) {
	secure := secureRequest(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AuthCookieName, token, maxAge, "/", "", secure, true)
	c.SetSameSite(http.SameSiteLaxMode)
	csrf := randomToken()
	c.SetCookie(CSRFCookieName, csrf, maxAge, "/", "", secure, false)
}

func ClearSessionCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AuthCookieName, "", -1, "/", "", secureRequest(c), true)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(CSRFCookieName, "", -1, "/", "", secureRequest(c), false)
}

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(CSRFCookieName)
		if err != nil || cookie == "" {
			cookie = randomToken()
			secure := secureRequest(c)
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie(CSRFCookieName, cookie, AuthCookieMaxAge, "/", "", secure, false)
			c.Next()
			return
		}

		method := c.Request.Method
		if method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions {
			if c.GetHeader("X-CSRF-Token") != cookie {
				c.JSON(http.StatusForbidden, gin.H{"message": "请求校验失败，请刷新页面后重试"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser := GetAuthUser(c)
		if authUser == nil || authUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "仅管理员可访问"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetAuthUser(c *gin.Context) *AuthUser {
	value, exists := c.Get("authUser")
	if !exists {
		return nil
	}
	authUser, ok := value.(AuthUser)
	if !ok {
		return nil
	}
	return &authUser
}
