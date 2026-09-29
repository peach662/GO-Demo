package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const ContextUserIDKey = "userID"

func AuthMiddleware(jwt *JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"data":    nil,
			})
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "前缀不正确",
				"data":    nil,
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token为空",
				"data":    nil,
			})
			c.Abort()
			return
		}
		userID, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, userID)
		c.Next()

	}
}
