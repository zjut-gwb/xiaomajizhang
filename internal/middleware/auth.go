package middleware

import (
	"strings"

	"github.com/asikeida/xiaomajizhang/internal/auth"
	"github.com/asikeida/xiaomajizhang/internal/response"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, 401, 40100, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, 401, 40100, "Token 格式错误")
			c.Abort()
			return
		}

		claims, err := auth.ParseAccessToken(secret, parts[1])
		if err != nil {
			response.Error(c, 401, 40100, "Token 无效")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
