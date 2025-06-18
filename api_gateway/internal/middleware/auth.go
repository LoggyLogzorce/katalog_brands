package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthInfo struct {
	UserID uint64 `json:"userID"`
	Role   string `json:"role"`
	Name   string `json:"name"`
}

// AuthMiddleware Мидлвэр для проверки JWT и роли
func AuthMiddleware(requiredRole []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("access_token")
		if err != nil {
			log.Println("AuthMiddleware: Ошибка получения access_token")
		}
		if authCookie == "" || !strings.HasPrefix(authCookie, "Bearer ") {
			log.Println("AuthMiddleware: нет токена или ошибочный формат")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "нет токена"})
			return
		}
		tokenString := strings.TrimPrefix(authCookie, "Bearer ")

		req, err := http.NewRequest("GET", "http://localhost:8081/api/v1/validate", nil)
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

		userIDStr := strconv.FormatUint(info.UserID, 10)

		c.Set("userID", userIDStr)
		c.Set("role", info.Role)
		c.Set("name", info.Name)

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := "guest" // по-умолчанию
		name := "U"
		var userID uint64

		authCookie, err := c.Cookie("access_token")
		if err != nil {
			log.Println("Ошибка получения access_token")
		}

		if strings.HasPrefix(authCookie, "Bearer ") {
			token := strings.TrimPrefix(authCookie, "Bearer ")
			// запрос к auth-сервису
			req, _ := http.NewRequest("GET", "http://localhost:8081/api/v1/validate", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			if resp, err := http.DefaultClient.Do(req); err == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				var info AuthInfo
				if body, _ := io.ReadAll(resp.Body); json.Unmarshal(body, &info) == nil {
					role = info.Role
					userID = info.UserID
					name = info.Name
				}
			}
		}

		userIDStr := strconv.FormatUint(userID, 10)
		// сохраняем в контекст и пропускаем
		c.Set("role", role)
		c.Set("name", name)
		c.Set("userID", userIDStr)
		c.Next()
	}
}
