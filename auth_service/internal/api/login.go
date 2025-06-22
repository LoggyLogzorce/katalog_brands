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

type AuthHandler struct {
	Repo storage.UserRepository
}

func NewAuthHandler(authRepo storage.UserRepository) *AuthHandler {
	return &AuthHandler{Repo: authRepo}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	user, err := h.Repo.SelectUser(c.Request.Context(), req)
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
