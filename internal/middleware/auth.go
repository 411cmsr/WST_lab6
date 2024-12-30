package middleware

import (
	"encoding/base64"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var validUsers = map[string]string{
	"root": "$2a$10$PbueWoNyctbsSD0b52FXvuDz4y2hDQ3z5HE.Sqi9eJIul6Mc7xnt2",
}

// Проверяем базовую аутентификацию
func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		payload := strings.TrimPrefix(auth, "Basic ")
		username, password := decodeBasicAuth(payload)

		if !isValidUser(username, password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Декодируем базовую аутентификацию в строку
func decodeBasicAuth(payload string) (string, string) {
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", ""
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// Проверяем данные пользователя
func isValidUser(username, password string) bool {
	if validPasswordHash, exists := validUsers[username]; exists {
		err := bcrypt.CompareHashAndPassword([]byte(validPasswordHash), []byte(password))
		return err == nil
	}
	return false

}
