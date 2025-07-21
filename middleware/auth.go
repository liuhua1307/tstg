package middleware

import (
	"strings"
	"tangsong-esports/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			utils.ErrorWithCode(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.ErrorWithCode(c, 401, "无效的认证令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("member_id", claims.MemberID)
		c.Set("account", claims.Account)
		c.Set("user_role", claims.UserRole)

		c.Next()
	}
}
