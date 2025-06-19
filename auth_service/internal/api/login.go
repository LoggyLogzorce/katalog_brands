package api

import (
	"auth_service/internal/models"
	"auth_service/internal/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strings"
)

func Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	user, err := storage.SelectUser(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(401, gin.H{"error": "не верные email или пароль"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка поиска пользователя"})
		return
	}

	if req.PasswordHash != user.PasswordHash {
		c.AbortWithStatusJSON(401, gin.H{"error": "не верные email или пароль"})
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка создания токена"})
		return
	}

	prefixToken := strings.Join([]string{"Bearer", token}, " ")

	c.SetCookie(
		"access_token",
		prefixToken,
		3600,
		"/",
		"localhost",
		true,
		true,
	)
	c.JSON(200, gin.H{})
}
