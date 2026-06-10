package middleware

import (
	"net/http"
	"strings"

	"blog/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	ID       uint
	Username string
	Role     string
	Status   string
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "缺少登录凭证"})
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "登录凭证无效或已过期"})
			c.Abort()
			return
		}

		c.Set("authUser", AuthUser{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
			Status:   "active",
		})
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
