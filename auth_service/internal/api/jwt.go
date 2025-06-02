package api

import (
	"auth_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// секрет для подписи токенов
var jwtSecret = []byte("secret_key_123123123")

// GenerateJWT создаёт JWT с полями user_id и role
func GenerateJWT(user *models.User) (string, error) {
	// задаём стандартные клеймы
	claims := jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(15 * time.Minute).Unix(), // срок действия 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
