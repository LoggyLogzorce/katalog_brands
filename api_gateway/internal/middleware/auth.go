package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type AuthInfo struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}

// AuthMiddleware Мидлвэр для проверки JWT и роли
func AuthMiddleware(requiredRole []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "нет токена"})
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		req, err := http.NewRequest("GET", "http://localhost:8081/validate", nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "ошибка запроса в auth-сервис"})
			return
		}

		req.Header.Set("Authorization", "Bearer "+tokenString)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth-сервис недоступен"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен", "detail": string(body)})
			return
		}

		var info AuthInfo
		if err = json.Unmarshal(body, &info); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ошибка разбора данных авторизации"})
			return
		}

		var ok bool

		for _, e := range requiredRole {
			if info.Role == e {
				ok = true
				break
			}
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
			return
		}

		c.Set("userID", info.UserID)
		c.Set("role", info.Role)

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		role := "guest" // по-умолчанию
		userID := ""

		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			// запрос к auth-сервису
			req, _ := http.NewRequest("GET", "http://localhost:8081/validate", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			if resp, err := http.DefaultClient.Do(req); err == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				var info AuthInfo
				if body, _ := io.ReadAll(resp.Body); json.Unmarshal(body, &info) == nil {
					role = info.Role
					userID = info.UserID
				}
			}
		}

		// сохраняем в контекст и пропускаем
		c.Set("role", role)
		c.Set("userID", userID)
		c.Next()
	}
}
