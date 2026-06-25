package pkgauth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"code": 401, "message": "Неавторизован"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "message": "Неверный токен"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "message": "Неверный id пользователя в токене"})
			c.Abort()
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}
