package middleware

import (
	"net/http"
	"soccer_manager_service/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userIDKey           = "user_id"
	emailKey            = "email"
)

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()

			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()

			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Set(emailKey, claims.Email)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return uuid.Nil, false
	}

	id, ok := userID.(uuid.UUID)

	return id, ok
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(emailKey)
	if !exists {
		return "", false
	}

	e, ok := email.(string)

	return e, ok
}
